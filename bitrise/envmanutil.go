package bitrise

import (
	"os"
	"os/exec"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/bitrise-io/go-utils/command"
)

// RunEnvmanInit ...
func RunEnvmanInit() error {
	logLevel := log.GetLevel().String()
	args := []string{"--loglevel", logLevel, "init"}
	return command.RunCommand("envman", args...)
}

// RunEnvmanAdd ...
func RunEnvmanAdd(key, value string, expand bool) error {
	logLevel := log.GetLevel().String()
	args := []string{"--loglevel", logLevel, "add", "--key", key, "--append"}
	if !expand {
		args = []string{"--loglevel", logLevel, "add", "--key", key, "--no-expand", "--append"}
	}

	envman := exec.Command("envman", args...)
	envman.Stdin = strings.NewReader(value)
	envman.Stdout = os.Stdout
	envman.Stderr = os.Stderr
	return envman.Run()
}

// RunEnvmanRunInDir ...
func RunEnvmanRunInDir(dir string, cmd []string, logLevel string) (int, error) {
	if logLevel == "" {
		logLevel = log.GetLevel().String()
	}
	args := []string{"--loglevel", logLevel, "run"}
	args = append(args, cmd...)
	return command.RunCommandInDirWithExitCode(dir, "envman", args...)
}

// RunEnvmanRun ...
func RunEnvmanRun(cmd []string) (int, error) {
	return RunEnvmanRunInDir("", cmd, "")
}

// RunEnvmanEnvstoreTest ...
func RunEnvmanEnvstoreTest(pth string) error {
	logLevel := log.GetLevel().String()
	args := []string{"--loglevel", logLevel, "--path", pth, "print"}
	cmd := exec.Command("envman", args...)
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
