package command

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/bitrise-io/go-utils/v2/env"
)

// ErrorFinder ...
type ErrorFinder func(out string) []string

// Opts ...
type Opts struct {
	Stdout      io.Writer
	Stderr      io.Writer
	Stdin       io.Reader
	Env         []string
	Dir         string
	ErrorFinder ErrorFinder

	// NewOutputObserver builds an observer for the command's output. When the command fails, what the
	// observer reports becomes the reason of the error. It takes the place of ErrorFinder when both are set.
	//
	// The command decides how the tool's output is piped, so that the observer reads it in a known order:
	//
	//   - Run, RunAndReturnExitCode and Start with Wait give the tool one pipe for stdout and stderr. One
	//     observer reads both in the order the tool wrote them. Everything is forwarded to Stdout, or to
	//     Stderr when only that is set, and nothing is written to the other one: two pipes would lose the
	//     order between the streams.
	//   - RunAndReturnTrimmedCombinedOutput works the same way, into the output it returns.
	//   - RunAndReturnTrimmedOutput has to keep stdout apart, so each stream gets an observer of its own
	//     and neither sees the other stream. The reasons of the stderr observer are reported first.
	//
	// The function is called once per observer and has to return a new one each time. An Opts used for
	// several commands, a retry among them, then never carries what one run read into another.
	NewOutputObserver func() OutputObserver
}

// Factory ...
type Factory interface {
	Create(name string, args []string, opts *Opts) Command
}

type factory struct {
	envRepository env.Repository
}

// NewFactory ...
func NewFactory(envRepository env.Repository) Factory {
	return factory{envRepository: envRepository}
}

// Create ...
func (f factory) Create(name string, args []string, opts *Opts) Command {
	cmd := exec.Command(name, args...)
	var collector *errorCollector
	var newObserver func() OutputObserver

	if opts != nil {
		newObserver = opts.NewOutputObserver
		if newObserver == nil && opts.ErrorFinder != nil {
			collector = &errorCollector{errorFinder: opts.ErrorFinder}
		}

		cmd.Stdout = opts.Stdout
		cmd.Stderr = opts.Stderr
		cmd.Stdin = opts.Stdin

		// If Env is nil, the new process uses the current process's
		// environment.
		// If we pass env vars we want to append them to the
		// current process's environment.
		cmd.Env = append(f.envRepository.List(), opts.Env...)
		cmd.Dir = opts.Dir
	}
	return &command{
		cmd:            cmd,
		errorCollector: collector,
		newObserver:    newObserver,
	}
}

// Command ...
type Command interface {
	PrintableCommandArgs() string
	Run() error
	RunAndReturnExitCode() (int, error)
	RunAndReturnTrimmedOutput() (string, error)
	RunAndReturnTrimmedCombinedOutput() (string, error)
	Start() error
	Wait() error
	Signal(sig os.Signal) error
	Kill() error
}

type command struct {
	cmd            *exec.Cmd
	errorCollector *errorCollector

	newObserver func() OutputObserver
	// observers are built when the outputs are wired, one per pipe. The stderr observer comes first.
	observers []OutputObserver
	closed    bool
}

// ErrProcessNotStarted ...
var ErrProcessNotStarted = errors.New("command has not been started")

// PrintableCommandArgs ...
func (c *command) PrintableCommandArgs() string {
	return printableCommandArgs(false, c.cmd.Args)
}

// Run ...
func (c *command) Run() error {
	c.wrapOutputs()

	err := c.cmd.Run()
	c.closeObservers()

	if err != nil {
		return c.wrapError(err)
	}

	return nil
}

// RunAndReturnExitCode ...
func (c *command) RunAndReturnExitCode() (int, error) {
	c.wrapOutputs()

	err := c.cmd.Run()
	c.closeObservers()

	if err != nil {
		err = c.wrapError(err)
	}

	exitCode := c.cmd.ProcessState.ExitCode()
	return exitCode, err
}

// RunAndReturnTrimmedOutput ...
func (c *command) RunAndReturnTrimmedOutput() (string, error) {
	// A Stdout the caller set is left to cmd.Output, which refuses it.
	if c.newObserver != nil && c.cmd.Stdout == nil {
		return c.observedOutput()
	}

	outBytes, err := c.cmd.Output()
	outStr := string(outBytes)
	if err != nil {
		if c.errorCollector != nil {
			c.errorCollector.collectErrors(outStr)
		}
		err = c.wrapError(err)
	}

	return strings.TrimSpace(outStr), err
}

// RunAndReturnTrimmedCombinedOutput ...
func (c *command) RunAndReturnTrimmedCombinedOutput() (string, error) {
	// A writer the caller set is left to cmd.CombinedOutput, which refuses it.
	if c.newObserver != nil && c.cmd.Stdout == nil && c.cmd.Stderr == nil {
		return c.observedCombinedOutput()
	}

	outBytes, err := c.cmd.CombinedOutput()
	outStr := string(outBytes)
	if err != nil {
		if c.errorCollector != nil {
			c.errorCollector.collectErrors(outStr)
		}
		err = c.wrapError(err)
	}

	return strings.TrimSpace(outStr), err
}

// observedOutput is RunAndReturnTrimmedOutput with observers. Stdout is what the caller gets back, so the two
// streams cannot share a pipe and each has an observer of its own.
func (c *command) observedOutput() (string, error) {
	stderrObserver := c.addObserver()
	stdoutObserver := c.addObserver()

	// cmd.Output keeps what the command wrote to stderr for the error when the caller did not ask for it.
	// The same is kept here, so that the error falls back to it when the observers found nothing.
	var captured *headTailBuffer
	stderr := c.cmd.Stderr
	if stderr == nil {
		captured = &headTailBuffer{limit: stderrCaptureLimit}
		stderr = captured
	}

	var stdout bytes.Buffer
	c.cmd.Stdout = &tap{dest: &stdout, observer: stdoutObserver}
	c.cmd.Stderr = &tap{dest: stderr, observer: stderrObserver}

	err := c.cmd.Run()
	c.closeObservers()

	if err != nil {
		var exitErr *exec.ExitError
		if captured != nil && errors.As(err, &exitErr) {
			exitErr.Stderr = captured.Bytes()
		}
		err = c.wrapError(err)
	}

	return strings.TrimSpace(stdout.String()), err
}

// observedCombinedOutput is RunAndReturnTrimmedCombinedOutput with an observer. Both streams go into one buffer
// through one pipe, as they do in cmd.CombinedOutput.
func (c *command) observedCombinedOutput() (string, error) {
	var output bytes.Buffer
	combined := &tap{dest: &output, observer: c.addObserver()}
	c.cmd.Stdout = combined
	c.cmd.Stderr = combined

	err := c.cmd.Run()
	c.closeObservers()

	if err != nil {
		err = c.wrapError(err)
	}

	return strings.TrimSpace(output.String()), err
}

// Start ...
func (c *command) Start() error {
	c.wrapOutputs()
	return c.cmd.Start()
}

// Wait ...
func (c *command) Wait() error {
	err := c.cmd.Wait()
	c.closeObservers()

	if err != nil {
		err = c.wrapError(err)
	}

	return err
}

// Signal ...
func (c *command) Signal(sig os.Signal) error {
	if c.cmd.Process == nil {
		return ErrProcessNotStarted
	}
	// os.Process tracks completion itself, so signalling after Wait reports ErrProcessDone
	// instead of reaching a process that reused the PID. Reading cmd.ProcessState here would
	// race with Wait.
	if err := c.cmd.Process.Signal(sig); err != nil {
		return fmt.Errorf("signalling command failed (%s): %w", c.PrintableCommandArgs(), err)
	}

	return nil
}

// Kill ...
func (c *command) Kill() error {
	return c.Signal(os.Kill)
}

func printableCommandArgs(isQuoteFirst bool, fullCommandArgs []string) string {
	var cmdArgsDecorated []string
	for idx, anArg := range fullCommandArgs {
		quotedArg := fmt.Sprintf("\"%s\"", anArg)
		if idx == 0 && !isQuoteFirst {
			quotedArg = anArg
		}
		cmdArgsDecorated = append(cmdArgsDecorated, quotedArg)
	}

	return strings.Join(cmdArgsDecorated, " ")
}

func (c *command) wrapError(err error) error {
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		return NewExitStatusError(c.PrintableCommandArgs(), exitErr, c.reasons())
	}

	return fmt.Errorf("executing command failed (%s): %w", c.PrintableCommandArgs(), err)
}

// reasons returns what was found in the output about why the command failed.
func (c *command) reasons() []string {
	reasons := []string{}

	if c.newObserver != nil {
		for _, observer := range c.observers {
			reasons = append(reasons, observer.Reasons()...)
		}

		// An observer can report as much as it likes. The cap is applied here, so it holds for every observer.
		return capReasons(reasons)
	}

	if c.errorCollector != nil {
		reasons = c.errorCollector.collectedErrorLines()
	}

	return reasons
}

// wrapOutputs puts a tap in front of the caller's writers so that the observer or the ErrorFinder sees the output.
func (c *command) wrapOutputs() {
	if c.newObserver != nil {
		// One tap for both streams: os/exec then gives the tool one pipe, and the observer reads stdout and
		// stderr in the order the tool wrote them.
		merged := &tap{dest: c.mergedOutput(), observer: c.addObserver()}
		c.cmd.Stdout = merged
		c.cmd.Stderr = merged
		return
	}

	if c.errorCollector == nil {
		return
	}

	// os/exec gives the child a single pipe and a single copy goroutine for stdout and stderr only while
	// cmd.Stdout and cmd.Stderr compare equal. Two separate taps would break that and hand a writer the
	// caller shares between the two streams (a bytes.Buffer, a filter) to two goroutines at once. Tap once instead.
	if interfaceEqual(c.cmd.Stdout, c.cmd.Stderr) {
		shared := c.tapForFinder(c.cmd.Stdout)
		c.cmd.Stdout = shared
		c.cmd.Stderr = shared
		return
	}

	c.cmd.Stdout = c.tapForFinder(c.cmd.Stdout)
	c.cmd.Stderr = c.tapForFinder(c.cmd.Stderr)
}

func (c *command) tapForFinder(w io.Writer) io.Writer {
	if w == nil {
		w = io.Discard
	}

	return &tap{dest: w, observer: c.errorCollector}
}

// mergedOutput returns the one writer both streams are forwarded to when an observer reads them.
func (c *command) mergedOutput() io.Writer {
	if c.cmd.Stdout != nil {
		return c.cmd.Stdout
	}
	if c.cmd.Stderr != nil {
		return c.cmd.Stderr
	}

	return io.Discard
}

// addObserver builds the observer of one pipe.
func (c *command) addObserver() OutputObserver {
	observer := c.newObserver()
	c.observers = append(c.observers, observer)

	return observer
}

// closeObservers tells the observers that the command has ended. The copy goroutines of os/exec have finished by
// then, so nothing is observed afterwards.
func (c *command) closeObservers() {
	if c.closed {
		return
	}
	c.closed = true

	for _, observer := range c.observers {
		observer.Close()
	}
}

// interfaceEqual compares the two writers the way os/exec does, without panicking on uncomparable dynamic types.
func interfaceEqual(a, b any) bool {
	defer func() { _ = recover() }()
	return a == b
}
