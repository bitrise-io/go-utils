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

	// Deprecated: ErrorFinder is called once per written chunk, so a line that lands in two chunks is
	// cut in half and matches nothing. Use ErrorMatcher.
	ErrorFinder ErrorFinder

	// ErrorMatcher is given each complete line of the command's output and says why the command
	// failed. What it reports becomes the reason of the error Run returns, so a caller reports what
	// the tool said rather than its exit status. It takes the place of ErrorFinder when both are set.
	ErrorMatcher ErrorMatcher
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
	var matcher *matcherCollector

	if opts != nil {
		switch {
		case opts.ErrorMatcher != nil:
			matcher = &matcherCollector{matcher: opts.ErrorMatcher}
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
		cmd:              cmd,
		errorCollector:   collector,
		matcherCollector: matcher,
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
	cmd              *exec.Cmd
	errorCollector   *errorCollector
	matcherCollector *matcherCollector
	lineWriters      []*lineWriter
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
	case c.matcherCollector != nil:
		return c.matcherCollector.reasons()
	case c.errorCollector != nil:
		return c.errorCollector.collectedErrorLines()
	default:
		return []string{}
	}
}

func (c *command) wrapOutputs() {
	wrap := c.outputWrapper()
	if wrap == nil {
		return
	}

	// os/exec gives the child a single pipe and a single copy goroutine for stdout and stderr only while
	// cmd.Stdout and cmd.Stderr compare equal. Wrapping the two separately would break that and hand a writer the
	// caller shares between the two streams (a bytes.Buffer, a filter) to two goroutines at once. Wrap once instead.
	if c.cmd.Stderr != nil && interfaceEqual(c.cmd.Stdout, c.cmd.Stderr) {
		shared := wrap(c.cmd.Stdout)
		c.cmd.Stdout = shared
		c.cmd.Stderr = shared

		return
	}

	c.cmd.Stdout = wrap(c.cmd.Stdout)
	c.cmd.Stderr = wrap(c.cmd.Stderr)
}

// outputWrapper returns how one stream's writer is wrapped so that what the command prints reaches
// the collector as well, or nil when nothing collects.
func (c *command) outputWrapper() func(io.Writer) io.Writer {
	switch {
	case c.matcherCollector != nil:
		return c.collectLinesFrom
	case c.errorCollector != nil:
		return c.collectFrom
	default:
		return nil
	}
}

// collectLinesFrom keeps the writer it builds, because the matcher is owed a last line without a
// newline once the process ends.
func (c *command) collectLinesFrom(w io.Writer) io.Writer {
	lines := c.matcherCollector.collectFromWriter(orDiscard(w))
	c.lineWriters = append(c.lineWriters, lines)

	return lines
}

func (c *command) collectFrom(w io.Writer) io.Writer {
	if w == nil {
		return c.errorCollector
	}

	return io.MultiWriter(c.errorCollector, w)
}

// closeOutputs hands the matcher a last line the command left without a newline. The caller's writers
// already have every byte by then.
func (c *command) closeOutputs() {
	for _, w := range c.lineWriters {
		_ = w.Close()
	}
	c.lineWriters = nil
}

// collectFromOutput feeds output that was captured whole rather than streamed, so that the two
// RunAndReturn...Output calls report the same reason as a streamed run.
func (c *command) collectFromOutput(output string) {
	switch {
	case c.matcherCollector != nil:
		_, _ = io.WriteString(c.collectLinesFrom(nil), output)
		c.closeOutputs()
	case c.errorCollector != nil:
		c.errorCollector.collectErrors(output)
	}
}

func orDiscard(w io.Writer) io.Writer {
	if w == nil {
		return io.Discard
	}

	return w
}

// interfaceEqual compares the two writers the way os/exec does, without panicking on uncomparable dynamic types.
func interfaceEqual(a, b any) bool {
	defer func() { _ = recover() }()
	return a == b
}
