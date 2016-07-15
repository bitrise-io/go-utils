package cmdex

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
)

// ----------

const (
	systemRubyPth = "/usr/bin/ruby"
	brewRubyPth   = "/usr/local/bin/ruby"
)

// BitriseDeveloperModeEnvKey ...
const BitriseDeveloperModeEnvKey = "BITRISE_DEVELOPER_MODE"

// CommandModel ...
type CommandModel struct {
	cmd *exec.Cmd
}

// NewCommand ...
func NewCommand(name string, args ...string) *CommandModel {
	return &CommandModel{
		cmd: exec.Command(name, args...),
	}
}

// NewBashCommand ...
func NewBashCommand(name string, args ...string) *CommandModel {
	args = append([]string{"-c", name}, args...)
	return &CommandModel{
		cmd: exec.Command("bash", args...),
	}
}

// NewRubyCommand ...
func NewRubyCommand(whichRuby, name string, args ...string) *CommandModel {
	if whichRuby == systemRubyPth {
		args = append([]string{name}, args...)
		name = "sudo"
	}

	return &CommandModel{
		cmd: exec.Command(name, args...),
	}
}

// NewBundleCommandModel ...
func NewBundleCommandModel(name string, args ...string) *CommandModel {
	args = append([]string{name}, args...)
	return &CommandModel{
		cmd: exec.Command("bundle", args...),
	}
}

// Dir ...
func (command *CommandModel) Dir(dir string) *CommandModel {
	command.cmd.Dir = dir
	return command
}

// Envs ...
func (command *CommandModel) Envs(envs []string) *CommandModel {
	command.cmd.Env = envs
	return command
}

// Stdin ...
func (command *CommandModel) Stdin(in io.Reader) *CommandModel {
	command.cmd.Stdin = in
	return command
}

// Stdout ...
func (command *CommandModel) Stdout(out io.Writer) *CommandModel {
	command.cmd.Stdout = out
	return command
}

// Stderr ...
func (command *CommandModel) Stderr(err io.Writer) *CommandModel {
	command.cmd.Stderr = err
	return command
}

// Run ...
func (command CommandModel) Run() error {
	return command.cmd.Run()
}

// RunForExitCode ...
func (command CommandModel) RunForExitCode() (int, error) {
	cmdExitCode := 0
	if err := command.cmd.Run(); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			waitStatus, ok := exitError.Sys().(syscall.WaitStatus)
			if !ok {
				return 1, errors.New("Failed to cast exit status")
			}
			cmdExitCode = waitStatus.ExitStatus()
		}
		return cmdExitCode, err
	}

	return 0, nil
}

// RunForOutput ...
func (command CommandModel) RunForOutput() (string, error) {
	outBytes, err := command.cmd.Output()
	outStr := string(outBytes)
	return strings.TrimSpace(outStr), err
}

// RunForCombinedOutput ...
func (command CommandModel) RunForCombinedOutput() (string, error) {
	outBytes, err := command.cmd.CombinedOutput()
	outStr := string(outBytes)
	return strings.TrimSpace(outStr), err
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

// Replacement - RunCommandWithReaderAndWriters ...
func runCommandWithReaderAndWriters(inReader io.Reader, outWriter, errWriter io.Writer, name string, args ...string) error {
	return NewCommand(name, args...).Stdin(inReader).Stdout(outWriter).Stderr(errWriter).Run()
}

// RunCommandWithReaderAndWriters ...
func RunCommandWithReaderAndWriters(inReader io.Reader, outWriter, errWriter io.Writer, name string, args ...string) error {
	if os.Getenv(BitriseDeveloperModeEnvKey) == "1" {
		deprecationMsg := `DEPRECATED METHOD! - RunCommandWithReaderAndWriters(inReader io.Reader, outWriter, errWriter io.Writer, name string, args ...string) error
USE INSTEAD! - NewCommand(name, args...).Stdin(inReader).Stdout(outWriter).Stderr(errWriter).Run()
`
		fmt.Println(deprecationMsg)
	}

	cmd := exec.Command(name, args...)
	cmd.Stdin = inReader
	cmd.Stdout = outWriter
	cmd.Stderr = errWriter
	return cmd.Run()
}

// Replacement - RunCommandWithWriters ...
func runCommandWithWriters(outWriter, errWriter io.Writer, name string, args ...string) error {
	return NewCommand(name, args...).Stdout(outWriter).Stderr(errWriter).Run()
}

// RunCommandWithWriters ...
func RunCommandWithWriters(outWriter, errWriter io.Writer, name string, args ...string) error {
	if os.Getenv(BitriseDeveloperModeEnvKey) == "1" {
		deprecationMsg := `DEPRECATED METHOD! - RunCommandWithWriters(outWriter, errWriter io.Writer, name string, args ...string) error
USE INSTEAD! - NewCommand(name, args...).Stdout(outWriter).Stderr(errWriter).Run()
`
		fmt.Println(deprecationMsg)
	}

	cmd := exec.Command(name, args...)
	cmd.Stdout = outWriter
	cmd.Stderr = errWriter
	return cmd.Run()
}

// Replacement - RunCommandInDirWithEnvsAndReturnExitCode ...
func runCommandInDirWithEnvsAndReturnExitCode(envs []string, dir, name string, args ...string) (int, error) {
	return NewCommand(name, args...).Envs(envs).Dir(dir).RunForExitCode()
}

// RunCommandInDirWithEnvsAndReturnExitCode ...
func RunCommandInDirWithEnvsAndReturnExitCode(envs []string, dir, name string, args ...string) (int, error) {
	if os.Getenv(BitriseDeveloperModeEnvKey) == "1" {
		deprecationMsg := `DEPRECATED METHOD! - RunCommandInDirWithEnvsAndReturnExitCode(envs []string, dir, name string, args ...string) (int, error)
USE INSTEAD! - NewCommand(name, args...).Envs(envs).Dir(dir).RunForExitCode()
`
		fmt.Println(deprecationMsg)
	}

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

	cmdExitCode := 0
	if err := cmd.Run(); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			waitStatus, ok := exitError.Sys().(syscall.WaitStatus)
			if !ok {
				return 1, errors.New("Failed to cast exit status")
			}
			cmdExitCode = waitStatus.ExitStatus()
		}
		return cmdExitCode, err
	}

	return 0, nil
}

// Replacement - RunCommandInDirAndReturnExitCode ...
func runCommandInDirAndReturnExitCode(dir, name string, args ...string) (int, error) {
	return NewCommand(name, args...).Dir(dir).RunForExitCode()
}

// RunCommandInDirAndReturnExitCode ...
func RunCommandInDirAndReturnExitCode(dir, name string, args ...string) (int, error) {
	if os.Getenv(BitriseDeveloperModeEnvKey) == "1" {
		deprecationMsg := `DEPRECATED METHOD! - RunCommandInDirAndReturnExitCode(dir, name string, args ...string) (int, error)
USE INSTEAD! - NewCommand(name, args...).Dir(dir).RunForExitCode()
`
		fmt.Println(deprecationMsg)
	}

	return RunCommandInDirWithEnvsAndReturnExitCode([]string{}, dir, name, args...)
}

// Replacement - runCommandWithEnvsAndReturnExitCode ...
func runCommandWithEnvsAndReturnExitCode(envs []string, name string, args ...string) (int, error) {
	return NewCommand(name, args...).Envs(envs).RunForExitCode()
}

// RunCommandWithEnvsAndReturnExitCode ...
func RunCommandWithEnvsAndReturnExitCode(envs []string, name string, args ...string) (int, error) {
	if os.Getenv(BitriseDeveloperModeEnvKey) == "1" {
		deprecationMsg := `DEPRECATED METHOD! - RunCommandWithEnvsAndReturnExitCode(envs []string, name string, args ...string) (int, error)
USE INSTEAD! - NewCommand(name, args...).Envs(envs).RunForExitCode()
`
		fmt.Println(deprecationMsg)
	}

	return RunCommandInDirWithEnvsAndReturnExitCode(envs, "", name, args...)
}

// Replacement - RunCommandInDir ...
func runCommandInDir(dir, name string, args ...string) error {
	return NewCommand(name, args...).Dir(dir).Run()
}

// RunCommandInDir ...
func RunCommandInDir(dir, name string, args ...string) error {
	if os.Getenv(BitriseDeveloperModeEnvKey) == "1" {
		deprecationMsg := `DEPRECATED METHOD! - RunCommandInDir(dir, name string, args ...string) error
USE INSTEAD! - NewCommand(name, args...).Dir(dir).Run()
`
		fmt.Println(deprecationMsg)
	}

	cmd := exec.Command(name, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if dir != "" {
		cmd.Dir = dir
	}
	return cmd.Run()
}

// Replacement - runCommand ...
func runCommand(name string, args ...string) error {
	return NewCommand(name, args...).Run()
}

// RunCommand ...
func RunCommand(name string, args ...string) error {
	if os.Getenv(BitriseDeveloperModeEnvKey) == "1" {
		deprecationMsg := `DEPRECATED METHOD! - RunCommand(name string, args ...string) error
USE INSTEAD! - NewCommand(name, args...).Run()
`
		fmt.Println(deprecationMsg)
	}

	return RunCommandInDir("", name, args...)
}

// Replacement - RunCommandAndReturnStdout ...
func runCommandAndReturnStdout(name string, args ...string) (string, error) {
	return NewCommand(name, args...).RunForOutput()
}

// RunCommandAndReturnStdout ..
func RunCommandAndReturnStdout(name string, args ...string) (string, error) {
	if os.Getenv(BitriseDeveloperModeEnvKey) == "1" {
		deprecationMsg := `DEPRECATED METHOD! - RunCommandAndReturnStdout(name string, args ...string) (string, error)
USE INSTEAD! - NewCommand(name, args...).RunForOutput()
`
		fmt.Println(deprecationMsg)
	}

	cmd := exec.Command(name, args...)
	outBytes, err := cmd.Output()
	outStr := string(outBytes)
	return strings.TrimSpace(outStr), err
}

// Replacement - runCommandInDirAndReturnCombinedStdoutAndStderr ...
func runCommandInDirAndReturnCombinedStdoutAndStderr(dir, name string, args ...string) (string, error) {
	return NewCommand(name, args...).Dir(dir).RunForCombinedOutput()
}

// RunCommandInDirAndReturnCombinedStdoutAndStderr ...
func RunCommandInDirAndReturnCombinedStdoutAndStderr(dir, name string, args ...string) (string, error) {
	if os.Getenv(BitriseDeveloperModeEnvKey) == "1" {
		deprecationMsg := `DEPRECATED METHOD! - RunCommandInDirAndReturnCombinedStdoutAndStderr(dir, name string, args ...string) (string, error)
USE INSTEAD! - NewCommand(name, args...).Dir(dir).RunForCombinedOutput()
`
		fmt.Println(deprecationMsg)
	}

	cmd := exec.Command(name, args...)
	if dir != "" {
		cmd.Dir = dir
	}
	outBytes, err := cmd.CombinedOutput()
	outStr := string(outBytes)
	return strings.TrimSpace(outStr), err
}

// Replacement - RunCommandAndReturnCombinedStdoutAndStderr ...
func runCommandAndReturnCombinedStdoutAndStderr(name string, args ...string) (string, error) {
	return NewCommand(name, args...).RunForCombinedOutput()
}

// RunCommandAndReturnCombinedStdoutAndStderr ..
func RunCommandAndReturnCombinedStdoutAndStderr(name string, args ...string) (string, error) {
	if os.Getenv(BitriseDeveloperModeEnvKey) == "1" {
		deprecationMsg := `DEPRECATED METHOD! - RunCommandAndReturnCombinedStdoutAndStderr(name string, args ...string) (string, error)
USE INSTEAD! - NewCommand(name, args...).RunForCombinedOutput()
`
		fmt.Println(deprecationMsg)
	}

	return RunCommandInDirAndReturnCombinedStdoutAndStderr("", name, args...)
}

// Replacement - runBashCommand ...
func runBashCommand(cmdStr string) error {
	return NewBashCommand(cmdStr).Run()
}

// Or mutch more better
// Replacement - runBashCommand ...
func betterRunBashCommand(name string, args ...string) error {
	return NewBashCommand(name, args...).Run()
}

// RunBashCommand ...
func RunBashCommand(cmdStr string) error {
	if os.Getenv(BitriseDeveloperModeEnvKey) == "1" {
		deprecationMsg := `DEPRECATED METHOD! - RunBashCommand(cmdStr string) error
USE INSTEAD! - NewBashCommand(name, args...).Run()
`
		fmt.Println(deprecationMsg)
	}

	return RunCommand("bash", "-c", cmdStr)
}

// Replacement - runBashCommandLines ...
func runBashCommandLines(cmdLines []string) error {
	for _, aLine := range cmdLines {
		if err := NewBashCommand(aLine).Run(); err != nil {
			return err
		}
	}
	return nil
}

// RunBashCommandLines ...
func RunBashCommandLines(cmdLines []string) error {
	if os.Getenv(BitriseDeveloperModeEnvKey) == "1" {
		deprecationMsg := `DEPRECATED METHOD! - RunBashCommandLines(cmdLines []string) error
USE INSTEAD! - NewBashCommand(aLine).Run()
`
		fmt.Println(deprecationMsg)
	}

	for _, aLine := range cmdLines {
		if err := RunCommand("bash", "-c", aLine); err != nil {
			return err
		}
	}
	return nil
}
