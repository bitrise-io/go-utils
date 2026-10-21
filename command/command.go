package command

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/bitrise-io/go-utils/v2/env"
	"github.com/bitrise-io/go-utils/v2/errorfinder"
	"github.com/bitrise-io/go-utils/v2/linepipe"
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

	// Deprecated: ErrorFinder is called once per written chunk, so a line that lands in two chunks is
	// cut in half and matches nothing. Use ErrorMatcher.
	ErrorFinder ErrorFinder

	// ErrorMatcher is given each complete line of the command's output, escape sequences removed, and
	// says why the command failed. What it reports becomes the reason of the error Run returns, so a
	// caller reports what the tool said rather than its exit status. It takes the place of
	// ErrorFinder when both are set.
	ErrorMatcher errorfinder.Matcher
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
	var reasons *errorfinder.Collector

	if opts != nil {
		switch {
		case opts.ErrorMatcher != nil:
			reasons = errorfinder.NewCollector(opts.ErrorMatcher)
		case opts.ErrorFinder != nil:
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
		cmd:             cmd,
		errorCollector:  collector,
		reasonCollector: reasons,
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
	cmd             *exec.Cmd
	errorCollector  *errorCollector
	reasonCollector *errorfinder.Collector
	outputs         []io.Closer
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
	c.closeOutputs()

	if err != nil {
		return c.wrapError(err)
	}

	return nil
}

// RunAndReturnExitCode ...
func (c *command) RunAndReturnExitCode() (int, error) {
	c.wrapOutputs()

	err := c.cmd.Run()
	c.closeOutputs()

	if err != nil {
		err = c.wrapError(err)
	}

	exitCode := c.cmd.ProcessState.ExitCode()
	return exitCode, err
}

// RunAndReturnTrimmedOutput ...
func (c *command) RunAndReturnTrimmedOutput() (string, error) {
	outBytes, err := c.cmd.Output()
	outStr := string(outBytes)
	if err != nil {
		c.collectFromOutput(outStr)
		err = c.wrapError(err)
	}

	return strings.TrimSpace(outStr), err
}

// RunAndReturnTrimmedCombinedOutput ...
func (c *command) RunAndReturnTrimmedCombinedOutput() (string, error) {
	outBytes, err := c.cmd.CombinedOutput()
	outStr := string(outBytes)
	if err != nil {
		c.collectFromOutput(outStr)
		err = c.wrapError(err)
	}

	return strings.TrimSpace(outStr), err
}

// Start ...
func (c *command) Start() error {
	c.wrapOutputs()
	return c.cmd.Start()
}

// Wait ...
func (c *command) Wait() error {
	err := c.cmd.Wait()
	c.closeOutputs()

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
		return NewExitStatusError(c.PrintableCommandArgs(), exitErr, c.failureReason())
	}

	return fmt.Errorf("executing command failed (%s): %w", c.PrintableCommandArgs(), err)
}

func (c *command) failureReason() []string {
	switch {
	case c.reasonCollector != nil:
		return c.reasonCollector.Reasons()
	case c.errorCollector != nil:
		return c.errorCollector.collectedErrorLines()
	default:
		return []string{}
	}
}

func (c *command) wrapOutputs() {
	if c.reasonCollector != nil {
		c.wrapOutputsIntoPipes()
		return
	}

	if c.errorCollector == nil {
		return
	}

	// os/exec gives the child a single pipe and a single copy goroutine for stdout and stderr only while
	// cmd.Stdout and cmd.Stderr compare equal. Two separate MultiWriters would break that and hand a writer the
	// caller shares between the two streams (a bytes.Buffer, a filter) to two goroutines at once. Wrap once instead.
	if c.cmd.Stderr != nil && interfaceEqual(c.cmd.Stdout, c.cmd.Stderr) {
		shared := c.collectFrom(c.cmd.Stdout)
		c.cmd.Stdout = shared
		c.cmd.Stderr = shared
		return
	}

	c.cmd.Stdout = c.collectFrom(c.cmd.Stdout)
	c.cmd.Stderr = c.collectFrom(c.cmd.Stderr)
}

func (c *command) collectFrom(w io.Writer) io.Writer {
	if w == nil {
		return c.errorCollector
	}
	return io.MultiWriter(c.errorCollector, w)
}

// wrapOutputsIntoPipes gives each stream a line pipe that forwards to the caller's writer and feeds
// the matcher. The streams share one pipe when the caller shares one writer, for the reason
// wrapOutputs gives.
func (c *command) wrapOutputsIntoPipes() {
	tap := linepipe.WithTap(c.reasonCollector)

	if c.cmd.Stderr != nil && interfaceEqual(c.cmd.Stdout, c.cmd.Stderr) {
		shared := linepipe.New(c.cmd.Stdout, tap)
		c.cmd.Stdout = shared
		c.cmd.Stderr = shared
		c.outputs = []io.Closer{shared}

		return
	}

	stdout := linepipe.New(c.cmd.Stdout, tap)
	stderr := linepipe.New(c.cmd.Stderr, tap)
	c.cmd.Stdout = stdout
	c.cmd.Stderr = stderr
	c.outputs = []io.Closer{stdout, stderr}
}

// closeOutputs hands the matcher a last line the command left without a newline. The caller's writers
// already have every byte, so this costs nothing when there is no such line.
func (c *command) closeOutputs() {
	for _, output := range c.outputs {
		_ = output.Close()
	}
	c.outputs = nil
}

// collectFromOutput feeds output that was captured whole rather than streamed, so that the two
// RunAndReturn...Output calls report the same reason as a streamed run.
func (c *command) collectFromOutput(output string) {
	switch {
	case c.reasonCollector != nil:
		pipe := linepipe.New(nil, linepipe.WithTap(c.reasonCollector))
		_, _ = io.WriteString(pipe, output)
		_ = pipe.Close()
	case c.errorCollector != nil:
		c.errorCollector.collectErrors(output)
	}
}

// interfaceEqual compares the two writers the way os/exec does, without panicking on uncomparable dynamic types.
func interfaceEqual(a, b any) bool {
	defer func() { _ = recover() }()
	return a == b
}
