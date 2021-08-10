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
func NewWithStandardOuts(name string, args ...string) Command {
	return newCmdWrapper(name, args...).SetStdout(os.Stdout).SetStderr(os.Stderr)
}

// NewWithParams ...
func NewWithParams(params ...string) (Command, error) {
	if len(params) == 0 {
		return nil, errors.New("no command provided")
	} else if len(params) == 1 {
		return newCmdWrapper(params[0]), nil
	}

	return newCmdWrapper(params[0], params[1:]...), nil
}

// NewFromSlice ...
func NewFromSlice(slice []string) (Command, error) {
	return NewWithParams(slice...)
}

// NewWithCmd ...
func NewWithCmd(cmd *exec.Cmd) Command {
	return &cmdWrapper{cmd: cmd}
}

// GetCmd ...
func (m *cmdWrapper) GetCmd() *exec.Cmd {
	return m.cmd
}

// SetDir ...
func (m *cmdWrapper) SetDir(dir string) Command {
	m.cmd.Dir = dir
	return m
}

// SetEnvs ...
func (m *cmdWrapper) SetEnvs(envs ...string) Command {
	m.cmd.Env = envs
	return m
}

// AppendEnvs - appends the envs to the current os.Environ()
// Calling this multiple times will NOT append the envs one by one,
// only the last "envs" set will be appended to os.Environ()!
func (m *cmdWrapper) AppendEnvs(envs ...string) Command {
	return m.SetEnvs(append(os.Environ(), envs...)...)
}

// SetStdin ...
func (m *cmdWrapper) SetStdin(in io.Reader) Command {
	m.cmd.Stdin = in
	return m
}

// SetStdout ...
func (m *cmdWrapper) SetStdout(out io.Writer) Command {
	m.cmd.Stdout = out
	return m
}

// SetStderr ...
func (m *cmdWrapper) SetStderr(err io.Writer) Command {
	m.cmd.Stderr = err
	return m
}

// Run ...
func (m *cmdWrapper) Run() error {
	return m.cmd.Run()
}

// RunAndReturnExitCode ...
func (m *cmdWrapper) RunAndReturnExitCode() (int, error) {
	return runCmdAndReturnExitCode(m.cmd)
}

// RunAndReturnTrimmedOutput ...
func (m *cmdWrapper) RunAndReturnTrimmedOutput() (string, error) {
	return runCmdAndReturnTrimmedOutput(m.cmd)
}

// RunAndReturnTrimmedCombinedOutput ...
func (m *cmdWrapper) RunAndReturnTrimmedCombinedOutput() (string, error) {
	outBytes, err := m.cmd.CombinedOutput()
	outStr := string(outBytes)
	return strings.TrimSpace(outStr), err
}

// PrintableCommandArgs ...
func (m *cmdWrapper) PrintableCommandArgs() string {
	return PrintableCommandArgs(false, m.cmd.Args)
}

// Args ...
func (m *cmdWrapper) Args() []string {
	return m.GetCmd().Args
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
