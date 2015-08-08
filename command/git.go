package command

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"github.com/bitrise-io/go-utils/pathutil"
)

// DoGitClone ...
func DoGitClone(uri, pth string) (err error) {
	if uri == "" {
		return errors.New("Git Clone 'uri' missing")
	}
	if pth == "" {
		return errors.New("Git Clone 'pth' missing")
	}
	if err = RunCommand("git", "clone", "--recursive", uri, pth); err != nil {
		fmt.Errorf("Failed to git clone from (%s) to (%s)", uri, pth)
		return
	}
	return
}

// RunGitCloneAndCheckoutTagOrBranch ...
func RunGitCloneAndCheckoutTagOrBranch(uri, pth, tagOrBranch string) error {
	if uri == "" {
		return errors.New("Git Clone 'uri' missing")
	}
	if pth == "" {
		return errors.New("Git Clone 'path' missing")
	}
	if tagOrBranch == "" {
		return errors.New("Git Clone 'tag or branch' missing")
	}
	return RunCommand("git", []string{"clone", "--recursive", "--branch", tagOrBranch, uri, pth}...)
}

// RunGitCloneAndCheckoutTagOrBranchAndCheckHeadState ...
func RunGitCloneAndCheckoutTagOrBranchAndCheckHeadState(uri, pth, version, commithash string) (err error) {
	if uri == "" {
		return errors.New("Git Clone 'uri' missing")
	}
	if pth == "" {
		return errors.New("Git Clone 'pth' missing")
	}
	if version == "" {
		return errors.New("Git Clone 'version' missing")
	}
	if commithash == "" {
		return errors.New("Git Clone 'commithash' missing")
	}
	if err = RunCommand("git", "clone", "--recursive", uri, pth, "--branch", version); err != nil {
		return
	}

	// cleanup
	defer func() {
		if err != nil {
			if err := RemoveDir(pth); err != nil {
				fmt.Errorf("Failed to cleanup path (%s) error: (%v) ", pth, err)
			}
		}
	}()

	latestCommit, err := GetLatestGitCommitHashOnHead(pth)
	if err != nil {
		return
	}
	if commithash != latestCommit {
		return fmt.Errorf("Commit hash doesn't match the one specified for the version tag. (version tag: %s) (expected: %s) (got: %s)", version, latestCommit, commithash)
	}

	return
}

// DoGitPull ...
func DoGitPull(pth string) error {
	err := RunCommandInDir(pth, "git", []string{"pull"}...)
	if err != nil {
		fmt.Errorf("Git pull failed, error (%v)", err)
		return err
	}
	return nil
}

// DoGitUpdate ...
func DoGitUpdate(git, pth string) error {
	if exists, err := pathutil.IsPathExists(pth); err != nil {
		return err
	} else if !exists {
		fmt.Println("[STEPMAN] - Git path does not exist, do clone")
		return DoGitClone(git, pth)
	}

	fmt.Println("[STEPMAN] - Git path exist, do pull")
	return DoGitPull(pth)
}

// DoGitCheckout ...
func DoGitCheckout(dir, commithash string) error {
	if commithash == "" {
		return errors.New("Git Clone 'hash' missing")
	}
	return RunCommandInDir(dir, "git", "checkout", commithash)
}

// DoGitCheckoutBranch ...
func DoGitCheckoutBranch(repoPath, branch string) error {
	if branch == "" {
		return errors.New("Git checkout 'branch' missing")
	}
	if err := DoGitCheckout(repoPath, branch); err != nil {
		return RunCommandInDir(repoPath, "git", "checkout", "-b", branch)
	}
	return nil
}

// DoGitAddFile ...
func DoGitAddFile(repoPath, filePath string) error {
	if filePath == "" {
		return errors.New("Git add 'file' missing")
	}
	return RunCommandInDir(repoPath, "git", "add", filePath)
}

// DoGitPushToOrigin ...
func DoGitPushToOrigin(repoPath, branch string) error {
	return RunCommandInDir(repoPath, "git", "push", "-u", "origin", branch)
}

// CheckIsNoGitChanges ...
func CheckIsNoGitChanges(repoPath string) error {
	return RunCommandInDir(repoPath, "git", "diff", "--cached", "--exit-code", "--quiet")
}

// DoGitCommit ...
func DoGitCommit(repoPath string, message string) error {
	if message == "" {
		return errors.New("Git commit 'message' missing")
	}
	if err := CheckIsNoGitChanges(repoPath); err != nil {
		return RunCommandInDir(repoPath, "git", "commit", "-m", message)
	}
	return nil
}

// GetLatestGitCommitHashOnHead ...
func GetLatestGitCommitHashOnHead(pth string) (string, error) {
	cmd := exec.Command("git", "rev-parse", "HEAD")
	cmd.Dir = pth
	bytes, err := cmd.CombinedOutput()
	cmdOutput := string(bytes)
	if err != nil {
		fmt.Errorf(cmdOutput)
		return "", err
	}
	return strings.TrimSpace(cmdOutput), nil
}

// DoGitGetCommit ...
func DoGitGetCommit(pth string) (string, error) {
	cmd := exec.Command("git", "rev-parse", "HEAD")
	cmd.Dir = pth
	bytes, err := cmd.CombinedOutput()
	cmdOutput := string(bytes)
	if err != nil {
		fmt.Errorf(cmdOutput)
		return "", err
	}
	return cmdOutput, nil
}
