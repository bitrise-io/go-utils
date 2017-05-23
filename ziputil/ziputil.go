package ziputil

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/bitrise-io/go-utils/command"
	"github.com/bitrise-io/go-utils/pathutil"
)

// Zip ...
func Zip(sourcePth, destinationZipPth string) error {
	parentDir := filepath.Dir(sourcePth)
	dirName := filepath.Base(sourcePth)
	cmd := command.New("/usr/bin/zip", "-rTy", destinationZipPth, dirName)
	cmd.SetDir(parentDir)
	out, err := cmd.RunAndReturnTrimmedCombinedOutput()
	if err != nil {
		return fmt.Errorf("Failed to zip: %s, output: %s, error: %s", sourcePth, out, err)
	}
	return nil
}

// UnZip ...
func UnZip(sourceZipPth string) (string, error) {
	// copy zip to tmp dir
	tmpDir, err := pathutil.NormalizedOSTempDirPath("__unzip__")
	if err != nil {
		return "", err
	}

	tmpSourcePth := filepath.Join(tmpDir, filepath.Base(sourceZipPth))
	if err := command.CopyFile(sourceZipPth, tmpSourcePth); err != nil {
		return "", err
	}
	// ---

	// unzip
	cmd := command.New("/usr/bin/unzip", sourceZipPth)
	cmd.SetDir(tmpDir)
	out, err := cmd.RunAndReturnTrimmedCombinedOutput()
	if err != nil {
		return "", fmt.Errorf("Failed to unzip: %s, output: %s, error: %s", tmpSourcePth, out, err)
	}
	// ---

	// serch for the unzipped path
	if err := os.RemoveAll(tmpSourcePth); err != nil {
		return "", err
	}

	pattern := filepath.Join(tmpDir, "*")
	pths, err := filepath.Glob(pattern)
	if err != nil {
		return "", err
	}
	if len(pths) == 0 {
		return "", fmt.Errorf("unzipped file not found")
	} else if len(pths) > 1 {
		return "", fmt.Errorf("multiple file generated")
	}
	// ---

	return pths[0], nil
}
