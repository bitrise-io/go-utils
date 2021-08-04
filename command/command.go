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
	PrintableCommandArgs() string
	Run() error
	RunAndReturnTrimmedOutput() (string, error)
	RunAndReturnExitCode() (int, error)
	RunAndReturnTrimmedCombinedOutput() (string, error)
}

// Factory ...
type Factory func(name string, args ...string) Command

// NewCommand ...
func NewCommand(name string, args ...string) Command {
	return newCommand(name, args...)
}

type cmdWrapper struct {
	cmd *exec.Cmd
}

func newCommand(name string, args ...string) Command {
	return &cmdWrapper{
		cmd: exec.Command(name, args...),
	}
}

// NewWithStandardOuts - same as NewCommand, but sets the command's
// stdout and stderr to the standard (OS) out (os.Stdout) and err (os.Stderr)
func NewWithStandardOuts(name string, args ...string) Command {
	return newCommand(name, args...).SetStdout(os.Stdout).SetStderr(os.Stderr)
}

// NewWithParams ...
func NewWithParams(params ...string) (Command, error) {
	if len(params) == 0 {
		return nil, errors.New("no command provided")
	} else if len(params) == 1 {
		return newCommand(params[0]), nil
	}

	return newCommand(params[0], params[1:]...), nil
}

// NewFromSlice ...
func NewFromSlice(slice []string) (Command, error) {
	return NewWithParams(slice...)
}

// NewWithCmd ...
func NewWithCmd(cmd *exec.Cmd) Command {
	return &cmdWrapper{
		cmd: cmd,
	}
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
	return RunCmdAndReturnExitCode(m.cmd)
}

// RunAndReturnTrimmedOutput ...
func (m *cmdWrapper) RunAndReturnTrimmedOutput() (string, error) {
	return RunCmdAndReturnTrimmedOutput(m.cmd)
}

// RunAndReturnTrimmedCombinedOutput ...
func (m *cmdWrapper) RunAndReturnTrimmedCombinedOutput() (string, error) {
	return RunCmdAndReturnTrimmedCombinedOutput(m.cmd)
}

// PrintableCommandArgs ...
func (m *cmdWrapper) PrintableCommandArgs() string {
	return PrintableCommandArgs(false, m.cmd.Args)
}

// ----------

// PrintableCommandArgs ...
func PrintableCommandArgs(isQuoteFirst bool, fullCommandArgs []string) string {
	cmdArgsDecorated := []string{}
	for idx, anArg := range fullCommandArgs {
		quotedArg := strconv.Quote(anArg)
		if idx == 0 && !isQuoteFirst {
			quotedArg = anArg
		}
		cmdArgsDecorated = append(cmdArgsDecorated, quotedArg)
	}

	return strings.Join(cmdArgsDecorated, " ")
}

// RunCmdAndReturnExitCode ...
func RunCmdAndReturnExitCode(cmd *exec.Cmd) (exitCode int, err error) {
	err = cmd.Run()
	exitCode = cmd.ProcessState.ExitCode()
	return
}

// RunCmdAndReturnTrimmedOutput ...
func RunCmdAndReturnTrimmedOutput(cmd *exec.Cmd) (string, error) {
	outBytes, err := cmd.Output()
	outStr := string(outBytes)
	return strings.TrimSpace(outStr), err
}

// RunCmdAndReturnTrimmedCombinedOutput ...
func RunCmdAndReturnTrimmedCombinedOutput(cmd *exec.Cmd) (string, error) {
	outBytes, err := cmd.CombinedOutput()
	outStr := string(outBytes)
	return strings.TrimSpace(outStr), err
}

// RunCommandWithReaderAndWriters ...
func RunCommandWithReaderAndWriters(inReader io.Reader, outWriter, errWriter io.Writer, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdin = inReader
	cmd.Stdout = outWriter
	cmd.Stderr = errWriter
	return cmd.Run()
}

// RunCommandWithWriters ...
func RunCommandWithWriters(outWriter, errWriter io.Writer, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = outWriter
	cmd.Stderr = errWriter
	return cmd.Run()
}

// RunCommandInDirWithEnvsAndReturnExitCode ...
func RunCommandInDirWithEnvsAndReturnExitCode(envs []string, dir, name string, args ...string) (int, error) {
	cmd := exec.Command(name, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if dir != "" {
		cmd.Dir = dir
	}
	if len(envs) > 0 {
		cmd.Env = envs
	}

	return RunCmdAndReturnExitCode(cmd)
}

// RunCommandInDirAndReturnExitCode ...
func RunCommandInDirAndReturnExitCode(dir, name string, args ...string) (int, error) {
	return RunCommandInDirWithEnvsAndReturnExitCode([]string{}, dir, name, args...)
}

// RunCommandWithEnvsAndReturnExitCode ...
func RunCommandWithEnvsAndReturnExitCode(envs []string, name string, args ...string) (int, error) {
	return RunCommandInDirWithEnvsAndReturnExitCode(envs, "", name, args...)
}

// RunCommandInDir ...
func RunCommandInDir(dir, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if dir != "" {
		cmd.Dir = dir
	}
	return cmd.Run()
}

// RunCommand ...
func RunCommand(name string, args ...string) error {
	return RunCommandInDir("", name, args...)
}

// RunCommandAndReturnStdout ..
func RunCommandAndReturnStdout(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	return RunCmdAndReturnTrimmedOutput(cmd)
}

// RunCommandInDirAndReturnCombinedStdoutAndStderr ...
func RunCommandInDirAndReturnCombinedStdoutAndStderr(dir, name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	if dir != "" {
		cmd.Dir = dir
	}
	return RunCmdAndReturnTrimmedCombinedOutput(cmd)
}

// RunCommandAndReturnCombinedStdoutAndStderr ..
func RunCommandAndReturnCombinedStdoutAndStderr(name string, args ...string) (string, error) {
	return RunCommandInDirAndReturnCombinedStdoutAndStderr("", name, args...)
}

// RunBashCommand ...
func RunBashCommand(cmdStr string) error {
	return RunCommand("bash", "-c", cmdStr)
}

// RunBashCommandLines ...
func RunBashCommandLines(cmdLines []string) error {
	for _, aLine := range cmdLines {
		if err := RunCommand("bash", "-c", aLine); err != nil {
			return err
		}
	}
	return nil
}
