package command

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"strconv"
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
	var errorFinder ErrorFinder

	if opts != nil {
		errorFinder = opts.ErrorFinder

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
		cmd:         cmd,
		errorFinder: errorFinder,
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
}

type command struct {
	cmd         *exec.Cmd
	errorFinder ErrorFinder
}

// PrintableCommandArgs ...
func (c command) PrintableCommandArgs() string {
	return printableCommandArgs(false, c.cmd.Args)
}

// Run ...
func (c *command) Run() error {
	outBuffer, errBuffer := c.wrappedOutputBuffers()

	if err := c.cmd.Run(); err != nil {
		return c.wrapError(err, outBuffer.String(), errBuffer.String())
	}

	return nil
}

// RunAndReturnExitCode ...
func (c command) RunAndReturnExitCode() (int, error) {
	outBuffer, errBuffer := c.wrappedOutputBuffers()
	err := c.cmd.Run()
	if err != nil {
		err = c.wrapError(err, outBuffer.String(), errBuffer.String())
	}

	exitCode := c.cmd.ProcessState.ExitCode()
	return exitCode, err
}

// RunAndReturnTrimmedOutput ...
func (c command) RunAndReturnTrimmedOutput() (string, error) {
	outBytes, err := c.cmd.Output()
	outStr := string(outBytes)
	if err != nil {
		err = c.wrapError(err, outStr, "")
	}

	return strings.TrimSpace(outStr), err
}

// RunAndReturnTrimmedCombinedOutput ...
func (c command) RunAndReturnTrimmedCombinedOutput() (string, error) {
	outBytes, err := c.cmd.CombinedOutput()
	outStr := string(outBytes)

	if err != nil {
		err = c.wrapError(err, outStr, "")
	}

	return strings.TrimSpace(outStr), err
}

// Start ...
func (c command) Start() error {
	return c.cmd.Start()
}

// Wait ...
func (c command) Wait() error {
	outBuffer, errBuffer := c.wrappedOutputBuffers()
	err := c.cmd.Wait()
	if err != nil {
		err = c.wrapError(err, outBuffer.String(), errBuffer.String())
	}

	return err
}

func printableCommandArgs(isQuoteFirst bool, fullCommandArgs []string) string {
	var cmdArgsDecorated []string
	for idx, anArg := range fullCommandArgs {
		quotedArg := strconv.Quote(anArg)
		if idx == 0 && !isQuoteFirst {
			quotedArg = anArg
		}
		cmdArgsDecorated = append(cmdArgsDecorated, quotedArg)
	}

	return strings.Join(cmdArgsDecorated, " ")
}

func (c command) wrapError(err error, stdout, stderr string) error {
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		if c.errorFinder != nil {
			reasonsO := c.errorFinder(stdout)
			reasonsE := c.errorFinder(stderr)
			reasons := append(reasonsO, reasonsE...)
			if len(reasons) > 0 {
				return fmt.Errorf("command failed with exit status %d (%s): %w", exitErr.ExitCode(), c.PrintableCommandArgs(), errors.New(strings.Join(reasons, "\n")))
			}
		}
		return fmt.Errorf("command failed with exit status %d (%s)", exitErr.ExitCode(), c.PrintableCommandArgs())
	}
	return fmt.Errorf("executing command failed (%s): %w", c.PrintableCommandArgs(), err)
}

func (c command) wrappedOutputBuffers() (*bytes.Buffer, *bytes.Buffer) {
	var outBuffer, errBuffer bytes.Buffer
	if c.errorFinder != nil {
		if c.cmd.Stdout != nil {
			outWriter := io.MultiWriter(&outBuffer, c.cmd.Stdout)
			c.cmd.Stdout = outWriter
		} else {
			c.cmd.Stdout = &outBuffer
		}

		if c.cmd.Stderr != nil {
			errWriter := io.MultiWriter(&errBuffer, c.cmd.Stderr)
			c.cmd.Stderr = errWriter
		} else {
			c.cmd.Stderr = &errBuffer
		}
	}
	return &outBuffer, &errBuffer
}
