package command

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/bitrise-io/go-utils/errorutil"
)

// ----------

// Model ...
type Model struct {
	cmd *exec.Cmd
}

// New ...
func New(name string, args ...string) *Model {
	return &Model{
		cmd: exec.Command(name, args...),
	}
}

// NewWithStandardOuts - same as NewCommand, but sets the command's
// stdout and stderr to the standard (OS) out (os.Stdout) and err (os.Stderr)
func NewWithStandardOuts(name string, args ...string) *Model {
	return New(name, args...).SetStdout(os.Stdout).SetStderr(os.Stderr)
}

// NewFromSlice ...
func NewFromSlice(slice ...string) (*Model, error) {
	if len(slice) == 0 {
		return nil, errors.New("no command provided")
	} else if len(slice) == 1 {
		return New(slice[0]), nil
	}

	return New(slice[0], slice[1:]...), nil
}

// NewCmd ...
func NewCmd(cmd *exec.Cmd) *Model {
	return &Model{
		cmd: cmd,
	}
}

// GetCmd ...
func (command *Model) GetCmd() *exec.Cmd {
	return command.cmd
}

// SetDir ...
func (command *Model) SetDir(dir string) *Model {
	command.cmd.Dir = dir
	return command
}

// SetEnvs ...
func (command *Model) SetEnvs(envs ...string) *Model {
	command.cmd.Env = envs
	return command
}

// AppendEnvs - appends the envs to the current os.Environ()
// Calling this multiple times will NOT appens the envs one by one,
// only the last "envs" set will be appended to os.Environ()!
func (command *Model) AppendEnvs(envs ...string) *Model {
	return command.SetEnvs(append(os.Environ(), envs...)...)
}

// SetStdin ...
func (command *Model) SetStdin(in io.Reader) *Model {
	command.cmd.Stdin = in
	return command
}

// SetStdout ...
func (command *Model) SetStdout(out io.Writer) *Model {
	command.cmd.Stdout = out
	return command
}

// SetStderr ...
func (command *Model) SetStderr(err io.Writer) *Model {
	command.cmd.Stderr = err
	return command
}

// Run ...
func (command Model) Run() error {
	return command.cmd.Run()
}

// RunAndReturnExitCode ...
func (command Model) RunAndReturnExitCode() (int, error) {
	return RunCmdAndReturnExitCode(command.cmd)
}

// RunAndReturnTrimmedOutput ...
func (command Model) RunAndReturnTrimmedOutput() (string, error) {
	return RunCmdAndReturnTrimmedOutput(command.cmd)
}

// RunAndReturnTrimmedCombinedOutput ...
func (command Model) RunAndReturnTrimmedCombinedOutput() (string, error) {
	return RunCmdAndReturnTrimmedCombinedOutput(command.cmd)
}

// PrintableCommandArgs ...
func (command Model) PrintableCommandArgs() string {
	return PrintableCommandArgs(false, command.cmd.Args)
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
func RunCmdAndReturnExitCode(cmd *exec.Cmd) (int, error) {
	err := cmd.Run()
	if err != nil {
		exitCode, castErr := errorutil.CmdExitCodeFromError(err)
		if castErr != nil {
			return 1, fmt.Errorf("failed get exit code from error: %s, error: %s", err, castErr)
		}

		return exitCode, err
	}

	return 0, nil
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
