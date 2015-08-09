package bitrise

import (
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/bitrise-io/go-utils/command"
)

// EnvmanInit ...
func EnvmanInit(logLevel log.Level, clear bool) error {
	if clear {
		return command.RunCommand("envman", "--loglevel", logLevel.String(), "init", "--clear")
	}
	return command.RunCommand("envman", "--loglevel", logLevel.String(), "init", "")
}

// EnvmanInitWithoutClear ...
func EnvmanInitWithoutClear(logLevel log.Level) error {
	return EnvmanInit(logLevel, false)
}

// EnvmanAdd ...
func EnvmanAdd(logLevel log.Level, key, value string, isExpand, isAppend bool) error {
	args := []string{"--loglevel", logLevel.String(), "add", "--key", key}
	if !isExpand {
		args = append(args, "--no-expand")
	}
	if isAppend {
		args = append(args, "--append")
	}

	envman := exec.Command("envman", args...)
	envman.Stdin = strings.NewReader(value)
	envman.Stdout = os.Stdout
	envman.Stderr = os.Stderr
	return envman.Run()
}

// EnvmanAddModeAppend ...
func EnvmanAddModeAppend(logLevel log.Level, key, value string, isExpand bool) error {
	return EnvmanAdd(logLevel, key, value, isExpand, true)
}

// EnvmanAddFromFile ...
func EnvmanAddFromFile(logLevel log.Level, key, valueFilePth string, isExpand, isAppend bool) error {
	bytes, err := ioutil.ReadFile(valueFilePth)
	if err != nil {
		return err
	}
	return EnvmanAdd(logLevel, key, string(bytes), isExpand, isAppend)
}

// EnvmanRunInDir ...
func EnvmanRunInDir(logLevel log.Level, dir string, cmd []string) (int, error) {
	args := []string{"--loglevel", logLevel.String(), "run"}
	for _, command := range cmd {
		args = append(args, command)
	}
	return command.RunCommandInDirWithExitCode(dir, "envman", args...)
}

// EnvmanRun ...
func EnvmanRun(logLevel log.Level, cmd []string) (int, error) {
	return EnvmanRunInDir(logLevel, "", cmd)
}

// EnvmanPrint ...
func EnvmanPrint(logLevel log.Level, pth string) error {
	return command.RunCommand("envman", "--loglevel", logLevel.String(), "--path", pth, "print")
}

// EnvmanClear ...
func EnvmanClear(logLevel log.Level, pth string) error {
	return command.RunCommand("envman", "--loglevel", logLevel.String(), "--path", pth, "clear")
}
