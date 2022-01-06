package pathutil

import (
	"errors"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	pathutilV1 "github.com/bitrise-io/go-utils/pathutil"
)

// PathProvider ...
type PathProvider interface {
	CreateTempDir(prefix string) (string, error)
}

type pathProvider struct{}

// NewPathProvider ...
func NewPathProvider() PathProvider {
	return pathProvider{}
}

// CreateTempDir ...
func (pathProvider) CreateTempDir(prefix string) (string, error) {
	return pathutilV1.NormalizedOSTempDirPath(prefix)
}

// PathChecker ...
type PathChecker interface {
	IsPathExists(pth string) (bool, error)
}

type pathChecker struct{}

// NewPathChecker ...
func NewPathChecker() PathChecker {
	return pathChecker{}
}

// IsPathExists ...
func (c pathChecker) IsPathExists(pth string) (bool, error) {
	return pathutilV1.IsPathExists(pth)
}

// PathModifier ...
type PathModifier interface {
	AbsPath(pth string) (string, error)
}

type pathModifier struct{}

// NewPathModifier ...
func NewPathModifier() PathModifier {
	return pathModifier{}
}

// AbsPath expands ENV vars and the ~ character then calls Go's Abs
func (p pathModifier) AbsPath(pth string) (string, error) {
	if pth == "" {
		return "", errors.New("No Path provided")
	}

	pth, err := p.expandTilde(pth)
	if err != nil {
		return "", err
	}

	return filepath.Abs(os.ExpandEnv(pth))
}

func (pathModifier) expandTilde(pth string) (string, error) {
	if pth == "" {
		return "", errors.New("No Path provided")
	}

	if strings.HasPrefix(pth, "~") {
		pth = strings.TrimPrefix(pth, "~")

		if len(pth) == 0 || strings.HasPrefix(pth, "/") {
			return os.ExpandEnv("$HOME" + pth), nil
		}

		splitPth := strings.Split(pth, "/")
		username := splitPth[0]

		usr, err := user.Lookup(username)
		if err != nil {
			return "", err
		}

		pathInUsrHome := strings.Join(splitPth[1:], "/")

		return filepath.Join(usr.HomeDir, pathInUsrHome), nil
	}

	return pth, nil
}
