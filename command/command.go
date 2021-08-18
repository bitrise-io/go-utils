package command

import (
	"github.com/bitrise-io/go-utils/env"
	"io"
	"os/exec"
	"strconv"
	"strings"
)

// Opts ...
type Opts struct {
	Stdout io.Writer
	Stderr io.Writer
	Stdin  io.Reader
	Env    []string
	Dir    string
}

// Factory ...
type Factory interface {
	Create(name string, args []string, opts *Opts) Command
}

type defaultFactory struct {
	envRepository env.Repository
}

// NewFactory ...
func NewFactory(envRepository env.Repository) Factory {
	return defaultFactory{envRepository: envRepository}
}

// Create ...
func (f defaultFactory) Create(name string, args []string, opts *Opts) Command {
	cmd := exec.Command(name, args...)
	if opts != nil {
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
	return defaultCommand{cmd}
}

// Command ...
type Command interface {
	PrintableCommandArgs() string
	Run() error
	RunAndReturnExitCode() (int, error)
	RunAndReturnTrimmedOutput() (string, error)
	RunAndReturnTrimmedCombinedOutput() (string, error)
}

type defaultCommand struct {
	cmd *exec.Cmd
}

// GetCmd ...
func (c defaultCommand) GetCmd() *exec.Cmd {
	return c.cmd
}

// Run ...
func (c defaultCommand) Run() error {
	return c.cmd.Run()
}

// RunAndReturnExitCode ...
func (c defaultCommand) RunAndReturnExitCode() (int, error) {
	err := c.cmd.Run()
	exitCode := c.cmd.ProcessState.ExitCode()
	return exitCode, err
}

// RunAndReturnTrimmedOutput ...
func (c defaultCommand) RunAndReturnTrimmedOutput() (string, error) {
	outBytes, err := c.cmd.Output()
	outStr := string(outBytes)
	return strings.TrimSpace(outStr), err
}

// RunAndReturnTrimmedCombinedOutput ...
func (c defaultCommand) RunAndReturnTrimmedCombinedOutput() (string, error) {
	outBytes, err := c.cmd.CombinedOutput()
	outStr := string(outBytes)
	return strings.TrimSpace(outStr), err
}

// PrintableCommandArgs ...
func (c defaultCommand) PrintableCommandArgs() string {
	return printableCommandArgs(false, c.cmd.Args)
}

// Args ...
func (c defaultCommand) Args() []string {
	return c.cmd.Args
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
