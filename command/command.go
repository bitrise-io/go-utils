package command

import (
	"errors"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Command ...
type Command interface {
	SetDir(dir string) Command
	SetStdout(stdout io.Writer) Command
	SetStderr(stdout io.Writer) Command
	SetStdoutAndStderr() Command
	SetEnvs(envs ...string) Command
	AppendEnvs(envs ...string) Command
	Args() []string
	PrintableCommandArgs() string
	Run() error
	RunAndReturnTrimmedOutput() (string, error)
	RunAndReturnExitCode() (int, error)
	RunAndReturnTrimmedCombinedOutput() (string, error)
}

// Factory ...
type Factory func(name string, args ...string) Command

type cmdWrapper struct {
	cmd *exec.Cmd
}

func newCmdWrapper(name string, args ...string) *cmdWrapper {
	return &cmdWrapper{
		cmd: exec.Command(name, args...),
	}
}

// New ...
func New(name string, args ...string) Command {
	return newCmdWrapper(name, args...)
}

// NewWithStandardOuts - same as New, but sets the command's
// stdout and stderr to the standard (OS) out (os.Stdout) and err (os.Stderr)
// Deprecated: Use New and SetStdoutAndStderr instead.
func NewWithStandardOuts(name string, args ...string) Command {
	return newCmdWrapper(name, args...).SetStdout(os.Stdout).SetStderr(os.Stderr)
}

// NewWithParams ...
// Deprecated.
func NewWithParams(params ...string) (Command, error) {
	if len(params) == 0 {
		return nil, errors.New("no command provided")
	} else if len(params) == 1 {
		return newCmdWrapper(params[0]), nil
	}

	return newCmdWrapper(params[0], params[1:]...), nil
}

// NewFromSlice ...
// Deprecated.
func NewFromSlice(slice []string) (Command, error) {
	return NewWithParams(slice...)
}

// NewWithCmd ...
// Deprecated.
func NewWithCmd(cmd *exec.Cmd) Command {
	return &cmdWrapper{cmd: cmd}
}

// GetCmd ...
func (cmd *cmdWrapper) GetCmd() *exec.Cmd {
	return cmd.cmd
}

// SetDir ...
func (cmd *cmdWrapper) SetDir(dir string) Command {
	cmd.cmd.Dir = dir
	return cmd
}

// SetEnvs ...
func (cmd *cmdWrapper) SetEnvs(envs ...string) Command {
	cmd.cmd.Env = envs
	return cmd
}

// AppendEnvs - appends the envs to the current os.Environ()
// Calling this multiple times will NOT append the envs one by one,
// only the last "envs" set will be appended to os.Environ()!
func (cmd *cmdWrapper) AppendEnvs(envs ...string) Command {
	return cmd.SetEnvs(append(os.Environ(), envs...)...)
}

// SetStdin ...
func (cmd *cmdWrapper) SetStdin(in io.Reader) Command {
	cmd.cmd.Stdin = in
	return cmd
}

// SetStdout ...
func (cmd *cmdWrapper) SetStdout(out io.Writer) Command {
	cmd.cmd.Stdout = out
	return cmd
}

// SetStderr ...
func (cmd *cmdWrapper) SetStderr(err io.Writer) Command {
	cmd.cmd.Stderr = err
	return cmd
}

// SetStdoutAndStderr ...
func (cmd *cmdWrapper) SetStdoutAndStderr() Command {
	cmd.SetStdout(os.Stdout).SetStderr(os.Stderr)
	return cmd
}

// Run ...
func (cmd *cmdWrapper) Run() error {
	return cmd.cmd.Run()
}

// RunAndReturnExitCode ...
func (cmd *cmdWrapper) RunAndReturnExitCode() (int, error) {
	return runCmdAndReturnExitCode(cmd.cmd)
}

// RunAndReturnTrimmedOutput ...
func (cmd *cmdWrapper) RunAndReturnTrimmedOutput() (string, error) {
	return runCmdAndReturnTrimmedOutput(cmd.cmd)
}

// RunAndReturnTrimmedCombinedOutput ...
func (cmd *cmdWrapper) RunAndReturnTrimmedCombinedOutput() (string, error) {
	outBytes, err := cmd.cmd.CombinedOutput()
	outStr := string(outBytes)
	return strings.TrimSpace(outStr), err
}

// PrintableCommandArgs ...
func (cmd *cmdWrapper) PrintableCommandArgs() string {
	return PrintableCommandArgs(false, cmd.cmd.Args)
}

// Args ...
func (cmd *cmdWrapper) Args() []string {
	return cmd.cmd.Args
}

// PrintableCommandArgs ...
func PrintableCommandArgs(isQuoteFirst bool, fullCommandArgs []string) string {
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

// Deprecated: Use Command instead.
func runCmdAndReturnExitCode(cmd *exec.Cmd) (exitCode int, err error) {
	err = cmd.Run()
	exitCode = cmd.ProcessState.ExitCode()
	return
}

// Deprecated: Use Command instead.
func runCmdAndReturnTrimmedOutput(cmd *exec.Cmd) (string, error) {
	outBytes, err := cmd.Output()
	outStr := string(outBytes)
	return strings.TrimSpace(outStr), err
}

// Deprecated: Use Command instead.
func runCmdAndReturnTrimmedCombinedOutput(cmd *exec.Cmd) (string, error) {
	outBytes, err := cmd.CombinedOutput()
	outStr := string(outBytes)
	return strings.TrimSpace(outStr), err
}

// RunCommandWithWriters ...
// Deprecated: Use Command instead.
func RunCommandWithWriters(outWriter, errWriter io.Writer, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = outWriter
	cmd.Stderr = errWriter
	return cmd.Run()
}

// RunCommandWithEnvsAndReturnExitCode ...
// Deprecated: Use Command instead.
func RunCommandWithEnvsAndReturnExitCode(envs []string, name string, args ...string) (int, error) {
	cmd := exec.Command(name, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if len(envs) > 0 {
		cmd.Env = envs
	}

	return runCmdAndReturnExitCode(cmd)
}

// RunCommand ...
// Deprecated: Use Command instead.
func RunCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// RunCommandAndReturnStdout ...
// Deprecated: Use Command instead.
func RunCommandAndReturnStdout(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	return runCmdAndReturnTrimmedOutput(cmd)
}

// RunCommandAndReturnCombinedStdoutAndStderr ...
// Deprecated: Use Command instead.
func RunCommandAndReturnCombinedStdoutAndStderr(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	return runCmdAndReturnTrimmedCombinedOutput(cmd)
}
