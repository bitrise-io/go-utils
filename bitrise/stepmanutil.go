package bitrise

import (
	log "github.com/Sirupsen/logrus"
	"github.com/bitrise-io/go-utils/command"
)

// StepmanSetup ...
func StepmanSetup(logLevel log.Level, collection, copySpecPth string) error {
	args := []string{"--debug", "--loglevel", logLevel.String(), "setup", "--collection", collection}
	if copySpecPth != "" {
		args = append(args, "--copy-spec-json", copySpecPth)
	}
	return command.RunCommand("stepman", args...)
}

// RunStepmanSetupWithoutCopySpec ...
func RunStepmanSetupWithoutCopySpec(logLevel log.Level, collection string) error {
	return StepmanSetup(logLevel, collection, "")
}

// StepmanUpdate ...
func StepmanUpdate(logLevel log.Level, collection string) error {
	return command.RunCommand("stepman", "--debug", "--loglevel", logLevel.String(), "update", "--collection", collection)
}

// StepmanDownload ...
func StepmanDownload(logLevel log.Level, collection, stepID, stepVersion, dir string, update bool) error {
	args := []string{"--debug", "--loglevel", logLevel.String(), "activate", "--collection", collection,
		"--id", stepID, "--version", stepVersion, "--path", dir}
	if update {
		args = append(args, "--update")
	}

	return command.RunCommand("stepman", args...)
}

// StepmanActivate ...
func StepmanActivate(logLevel log.Level, collection, stepID, stepVersion, dir, ymlPth string, update bool) error {
	args := []string{"--debug", "--loglevel", logLevel.String(), "activate", "--collection", collection,
		"--id", stepID, "--version", stepVersion, "--path", dir, "--copyyml", ymlPth}
	if update {
		args = append(args, "--update")
	}
	return command.RunCommand("stepman", args...)
}
