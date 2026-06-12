package ziputil

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/bitrise-io/go-utils/command"
	"github.com/bitrise-io/go-utils/pathutil"
)

// ZipDir zips sourceDirPth into destinationZipPth.
// isContentOnly=false: archive contains the directory under its own basename (e.g. "mydir/file.txt").
// isContentOnly=true: archive contains the directory's contents directly (e.g. "file.txt").
// Directory modification times are preserved. Symlinks are stored as symlinks (not followed).
func ZipDir(sourceDirPth, destinationZipPth string, isContentOnly bool) error {
	if exist, err := pathutil.IsDirExists(sourceDirPth); err != nil {
		return err
	} else if !exist {
		return fmt.Errorf("dir (%s) not exist", sourceDirPth)
	}

	workDir := filepath.Dir(sourceDirPth)
	zipTarget := filepath.Base(sourceDirPth)

	if isContentOnly {
		workDir = sourceDirPth
		zipTarget = "."
	}

	return internalZipDir(destinationZipPth, zipTarget, workDir)

}

// ZipDirs zips multiple directories into a single archive, each under its own basename.
// When two entries in sourceDirPths share the same basename (e.g. "/a/shared" and
// "/b/shared"), their contents are merged; files unique to either are preserved,
// and conflicting filenames are overwritten by entries later in sourceDirPths.
// If the temporary directory cleanup fails, the process is terminated via log.Fatal.
func ZipDirs(sourceDirPths []string, destinationZipPth string) error {
	for _, path := range sourceDirPths {
		if exist, err := pathutil.IsDirExists(path); err != nil {
			return err
		} else if !exist {
			return fmt.Errorf("directory (%s) not exist", path)
		}
	}

	tempDir, err := pathutil.NormalizedOSTempDirPath("zip")
	if err != nil {
		return err
	}

	defer func() {
		if err = os.RemoveAll(tempDir); err != nil {
			log.Fatal(err)
		}
	}()

	for _, path := range sourceDirPths {
		err := command.CopyDir(path, tempDir, false)
		if err != nil {
			return err
		}
	}

	return internalZipDir(destinationZipPth, ".", tempDir)
}

func internalZipDir(destinationZipPth, zipTarget, workDir string) error {
	// -r: recurse into directories
	// -T: test the temp archive with unzip before overwriting the destination;
	//     on failure the destination is left in its prior state
	// -y: store symlinks as symlinks instead of following them
	cmd := command.New("/usr/bin/zip", "-rTy", destinationZipPth, zipTarget)
	cmd.SetDir(workDir)
	if out, err := cmd.RunAndReturnTrimmedCombinedOutput(); err != nil {
		return fmt.Errorf("command: (%s) failed, output: %s, error: %s", cmd.PrintableCommandArgs(), out, err)
	}

	return nil
}

// ZipFile zips a single file into destinationZipPth.
func ZipFile(sourceFilePth, destinationZipPth string) error {
	return ZipFiles([]string{sourceFilePth}, destinationZipPth)
}

// ZipFiles zips multiple files into destinationZipPth without preserving directory structure.
// Each file is stored under its base name only (-j junk paths).
// If two source files share the same base name, the command fails and the destination
// file is not created.
func ZipFiles(sourceFilePths []string, destinationZipPth string) error {
	for _, path := range sourceFilePths {
		if exist, err := pathutil.IsPathExists(path); err != nil {
			return err
		} else if !exist {
			return fmt.Errorf("file (%s) not exist", path)
		}
	}

	// -T: test the temp archive with unzip before overwriting the destination;
	//     on failure the destination is left in its prior state
	// -y: store symlinks as symlinks instead of following them
	// -j: junk directory paths, store files under their base names only
	parameters := []string{"-Tyj", destinationZipPth}
	parameters = append(parameters, sourceFilePths...)
	cmd := command.New("/usr/bin/zip", parameters...)
	if out, err := cmd.RunAndReturnTrimmedCombinedOutput(); err != nil {
		return fmt.Errorf("command: (%s) failed, output: %s, error: %s", cmd.PrintableCommandArgs(), out, err)
	}

	return nil
}

// UnZip extracts the zip archive at zipPth into intoDir.
// Entries with path-traversal components (e.g. "../../escape") are not rejected;
// instead the system unzip strips those components, extracts the file inside intoDir,
// and returns a non-zero exit code.
func UnZip(zip, intoDir string) error {
	cmd := command.New("/usr/bin/unzip", zip, "-d", intoDir)
	if out, err := cmd.RunAndReturnTrimmedCombinedOutput(); err != nil {
		return fmt.Errorf("command: (%s) failed, output: %s, error: %s", cmd.PrintableCommandArgs(), out, err)
	}

	return nil
}
