package ziputil

import (
	"fmt"
	"path/filepath"

	"github.com/bitrise-io/go-utils/command"
	"github.com/bitrise-io/go-utils/pathutil"
)

// Zip ...
func Zip(sourcePth, destinationZipPth string, isContentOnly bool) error {
	workDir := ""
	isDir := false

	if isContentOnly {
		if exist, err := pathutil.IsDirExists(sourcePth); err != nil {
			return err
		} else if exist {
			workDir = sourcePth
			isDir = true
		}
	}

	if workDir == "" {
		workDir = filepath.Dir(sourcePth)
	}

	targetName := ""
	if isContentOnly && isDir {
		targetName = "."
	} else {
		targetName = filepath.Base(sourcePth)
	}

	cmd := command.New("/usr/bin/zip", "-rTy", destinationZipPth, targetName)
	cmd.SetDir(workDir)
	out, err := cmd.RunAndReturnTrimmedCombinedOutput()
	if err != nil {
		return fmt.Errorf("$ %s\nFailed - output: %s, error: %s", cmd.PrintableCommandArgs(), out, err)
	}

	return nil
}

// UnZip ...
func UnZip(zip, intoDir string) error {
	workDir := filepath.Dir(intoDir)

	cmd := command.New("/usr/bin/unzip", zip, "-d", intoDir)
	cmd.SetDir(workDir)
	out, err := cmd.RunAndReturnTrimmedCombinedOutput()
	if err != nil {
		return fmt.Errorf("$ %s\nFailed - output: %s, error: %s", cmd.PrintableCommandArgs(), out, err)
	}

	return nil
}
