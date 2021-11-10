package ziputil

import (
	"fmt"
	"github.com/bitrise-io/go-utils/env"
	"path/filepath"

	"github.com/bitrise-io/go-utils/command"
	"github.com/bitrise-io/go-utils/pathutil"
)

// ZipDir ...
func ZipDir(sourceDirPth, destinationZipPth string, isContentOnly bool) error {
	if exist, err := pathutil.IsDirExists(sourceDirPth); err != nil {
		return err
	} else if !exist {
		return fmt.Errorf("dir (%s) not exist", sourceDirPth)
	}

	workDir := filepath.Dir(sourceDirPth)
	if isContentOnly {
		workDir = sourceDirPth
	}

	zipTarget := filepath.Base(sourceDirPth)
	if isContentOnly {
		zipTarget = "."
	}

	// -r - Travel the directory structure recursively
	// -T - Test the integrity of the new zip file
	// -y - Store symbolic links as such in the zip archive, instead of compressing and storing the file referred to by the link
	opts := &command.Opts{Dir: workDir}
	factory := command.NewFactory(env.NewRepository())
	cmd := factory.Create("/usr/bin/zip", []string{"-rTy", destinationZipPth, zipTarget}, opts)
	if out, err := cmd.RunAndReturnTrimmedCombinedOutput(); err != nil {
		return fmt.Errorf("command: (%s) failed, output: %s, error: %s", cmd.PrintableCommandArgs(), out, err)
	}

	return nil
}

// ZipFile ...
func ZipFile(sourceFilePth, destinationZipPth string) error {
	return ZipFiles([]string{sourceFilePth}, destinationZipPth)
}

// ZipFiles ...
func ZipFiles(sourceFilePths []string, destinationZipPth string) error {
	for _, path := range sourceFilePths {
		if exist, err := pathutil.IsPathExists(path); err != nil {
			return err
		} else if !exist {
			return fmt.Errorf("file (%s) not exist", path)
		}
	}

	factory := command.NewFactory(env.NewRepository())

	// -T - Test the integrity of the new zip file
	// -y - Store symbolic links as such in the zip archive, instead of compressing and storing the file referred to by the link
	// -j - Do not recreate the directory structure inside the zip. Kind of equivalent of copying all the files in one folder and zipping it.
	parameters := []string{"-Tyj", destinationZipPth}
	parameters = append(parameters, sourceFilePths...)

	cmd := factory.Create("/usr/bin/zip", parameters, nil)
	if out, err := cmd.RunAndReturnTrimmedCombinedOutput(); err != nil {
		return fmt.Errorf("command: (%s) failed, output: %s, error: %s", cmd.PrintableCommandArgs(), out, err)
	}

	return nil
}

// UnZip ...
func UnZip(zip, intoDir string) error {
	factory := command.NewFactory(env.NewRepository())
	cmd := factory.Create("/usr/bin/unzip", []string{zip, "-d", intoDir}, nil)
	if out, err := cmd.RunAndReturnTrimmedCombinedOutput(); err != nil {
		return fmt.Errorf("command: (%s) failed, output: %s, error: %s", cmd.PrintableCommandArgs(), out, err)
	}

	return nil
}
