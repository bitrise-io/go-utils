package bitrise

import (
	log "github.com/Sirupsen/logrus"
	"github.com/bitrise-io/go-utils/command"
)

// RunStepmanSetup ...
func RunStepmanSetup(collection string) error {
	logLevel := log.GetLevel().String()
	args := []string{"--debug", "--loglevel", logLevel, "setup", "--collection", collection}
	return command.RunCommand("stepman", args...)
}

// RunStepmanActivate ...
func RunStepmanActivate(collection, stepID, stepVersion, dir, ymlPth string) error {
	logLevel := log.GetLevel().String()
	args := []string{"--debug", "--loglevel", logLevel, "activate", "--collection", collection,
		"--id", stepID, "--version", stepVersion, "--path", dir, "--copyyml", ymlPth, "--update"}
	return command.RunCommand("stepman", args...)
}
