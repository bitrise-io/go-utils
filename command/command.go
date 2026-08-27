package command

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/bitrise-io/go-utils/v2/env"
)

// ErrorFinder ...
type ErrorFinder func(out string) []string

// Opts ...
type Opts struct {
	Stdout       io.Writer
	Stderr       io.Writer
	Stdin        io.Reader
	Env          []string
	Dir          string
	ErrorFinder  ErrorFinder
	ProcessGroup bool
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

		if opts.ProcessGroup {
			cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
		}
	}
	return &command{
		cmd:            cmd,
		errorCollector: collector,
		processGroup:   opts != nil && opts.ProcessGroup,
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
	processGroup   bool
}

// ErrProcessNotStarted ...
var ErrProcessNotStarted = errors.New("command has not been started")

// ErrProcessFinished is returned when the process exited before the signal could be delivered.
// A process can exit on its own at any point, so this is an expected outcome, not a failure.
var ErrProcessFinished = errors.New("process has already finished")

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
	// Wait releases the process once it exits, after which the PID can be reused by an
	// unrelated process.
	if c.cmd.ProcessState != nil {
		return ErrProcessFinished
	}

	var err error
	if c.processGroup {
		err = signalGroup(c.cmd.Process.Pid, sig)
	} else {
		err = c.cmd.Process.Signal(sig)
	}

	if errors.Is(err, os.ErrProcessDone) || errors.Is(err, syscall.ESRCH) {
		return ErrProcessFinished
	}
	if err != nil {
		return fmt.Errorf("signalling command failed (%s): %w", c.PrintableCommandArgs(), err)
	}

	return nil
}

// Kill ...
func (c command) Kill() error {
	return c.Signal(os.Kill)
}

func signalGroup(pid int, sig os.Signal) error {
	sysSig, ok := sig.(syscall.Signal)
	if !ok {
		return fmt.Errorf("unsupported signal: %v", sig)
	}

	return syscall.Kill(-pid, sysSig)
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
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		errorLines := []string{}
		if c.errorCollector != nil {
			errorLines = c.errorCollector.errorLines
		}

		return NewExitStatusError(c.PrintableCommandArgs(), exitErr, errorLines)
	}

	return fmt.Errorf("executing command failed (%s): %w", c.PrintableCommandArgs(), err)
}

func (c command) wrapOutputs() {
	if c.errorCollector == nil {
		return
	}

	if c.cmd.Stdout != nil {
		outWriter := io.MultiWriter(c.errorCollector, c.cmd.Stdout)
		c.cmd.Stdout = outWriter
	} else {
		c.cmd.Stdout = c.errorCollector
	}

	if c.cmd.Stderr != nil {
		errWriter := io.MultiWriter(c.errorCollector, c.cmd.Stderr)
		c.cmd.Stderr = errWriter
	} else {
		c.cmd.Stderr = c.errorCollector
	}
}
