package command

import (
	"os"
	"strings"

	"github.com/bitrise-io/go-utils/pathutil"
)

// RunCopy ...
func RunCopy(src, dst string) error {
	args := []string{src, dst}
	return RunCommand("rsync", args...)
}

// RunCopyFile ...
func RunCopyFile(src, dst string) error {
	args := []string{src, dst}
	return RunCommand("rsync", args...)
}

// RunCopyDir ...
func RunCopyDir(src, dst string, isOnlyContent bool) error {
	if isOnlyContent && !strings.HasSuffix(src, "/") {
		src = src + "/"
	}
	args := []string{"-r", src, dst}
	return RunCommand("rsync", args...)
}

// RemoveDir ...
func RemoveDir(dirPth string) error {
	if exist, err := pathutil.IsPathExists(dirPth); err != nil {
		return err
	} else if exist {
		if err := os.RemoveAll(dirPth); err != nil {
			return err
		}
	}
	return nil
}
