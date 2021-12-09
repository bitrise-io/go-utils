package pathutil

import (
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
