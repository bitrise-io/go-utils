package command

import (
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
	Stdout io.Writer
	Stderr io.Writer
	Stdin  io.Reader
	Env    []string
	Dir    string

	// ErrorFinder is called once per written chunk, so a line that lands in two chunks is cut in half
	// and matches nothing. For a tool that writes its lines in pieces, use Command.RunWithMatcher.
	ErrorFinder ErrorFinder
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

	if opts != nil {
		if opts.ErrorFinder != nil {
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
	}
}

// Command ...
type Command interface {
	PrintableCommandArgs() string
	Run() error
	// RunWithMatcher runs the command like Run and, when it fails, reports what matcher found as the reason.
	RunWithMatcher(matcher ErrorMatcher) error
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
}

// ErrProcessNotStarted ...
var ErrProcessNotStarted = errors.New("command has not been started")

// PrintableCommandArgs ...
func (c command) PrintableCommandArgs() string {
	return printableCommandArgs(false, c.cmd.Args)
}

// Run ...
func (c *command) Run() error {
	c.wrapOutputs()

	if err := c.cmd.Run(); err != nil {
		return c.wrapError(err)
	}

	return nil
}

// RunWithMatcher runs the command and, when it fails, reports what matcher found as the reason.
//
// The tool gets one pipe for stdout and stderr, so the matcher is given the lines of both in the order the tool
// wrote them. Two pipes cannot keep that order, so everything the tool prints is forwarded to Opts.Stdout
// (to Opts.Stderr when only that is set) and nothing is written to the other one. Each byte is forwarded as it
// arrives. The matcher is shown it afterwards and never holds output back.
//
// The pipe stays inside the command: there is no writer for a caller to share between commands or to give to one
// stream only. The matcher is given for this one run, so nothing it read reaches another run unless the caller
// passes the same matcher again. An ErrorFinder in Opts is not used by this method.
func (c *command) RunWithMatcher(matcher ErrorMatcher) error {
	lines := &lineSplitter{matcher: matcher}
	output := &tap{dest: c.mergedOutput(), observer: lines}
	c.cmd.Stdout = output
	c.cmd.Stderr = output

	err := c.cmd.Run()
	lines.close()

	if err != nil {
		return c.exitError(err, capReasons(matcher.Reasons()))
	}

	return nil
}

// RunAndReturnExitCode ...
func (c command) RunAndReturnExitCode() (int, error) {
	c.wrapOutputs()
	err := c.cmd.Run()
	if err != nil {
		err = c.wrapError(err)
	}

	exitCode := c.cmd.ProcessState.ExitCode()
	return exitCode, err
}

// RunAndReturnTrimmedOutput ...
func (c command) RunAndReturnTrimmedOutput() (string, error) {
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
func (c command) RunAndReturnTrimmedCombinedOutput() (string, error) {
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

// Start ...
func (c command) Start() error {
	c.wrapOutputs()
	return c.cmd.Start()
}

// Wait ...
func (c command) Wait() error {
	err := c.cmd.Wait()
	if err != nil {
		err = c.wrapError(err)
	}

	return err
}

// Signal ...
func (c command) Signal(sig os.Signal) error {
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
func (c command) Kill() error {
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

func (c command) wrapError(err error) error {
	errorLines := []string{}
	if c.errorCollector != nil {
		errorLines = c.errorCollector.collectedErrorLines()
	}

	return c.exitError(err, errorLines)
}

// exitError builds the error of a failed run from the reasons found in its output.
func (c command) exitError(err error, reasons []string) error {
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		return NewExitStatusError(c.PrintableCommandArgs(), exitErr, reasons)
	}

	return fmt.Errorf("executing command failed (%s): %w", c.PrintableCommandArgs(), err)
}

// mergedOutput returns the one writer RunWithMatcher forwards both streams to.
func (c command) mergedOutput() io.Writer {
	if c.cmd.Stdout != nil {
		return c.cmd.Stdout
	}
	if c.cmd.Stderr != nil {
		return c.cmd.Stderr
	}

	return io.Discard
}

// wrapOutputs puts a tap in front of the caller's writers so that the ErrorFinder sees the output.
func (c command) wrapOutputs() {
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

func (c command) tapForFinder(w io.Writer) io.Writer {
	if w == nil {
		w = io.Discard
	}

	return &tap{dest: w, observer: c.errorCollector}
}

// interfaceEqual compares the two writers the way os/exec does, without panicking on uncomparable dynamic types.
func interfaceEqual(a, b any) bool {
	defer func() { _ = recover() }()
	return a == b
}
