package fileutil

import (
	fileutilV1 "github.com/bitrise-io/go-utils/fileutil"
	"os"
	"path/filepath"
)

// FileManager ...
type FileManager interface {
	Remove(path string) error
	RemoveAll(path string) error
	Write(path string, value string, mode os.FileMode) error
}

type fileManager struct{}

// NewFileManager ...
func NewFileManager() FileManager {
	return fileManager{}
}

// Remove ...
func (fileManager) Remove(path string) error {
	return os.Remove(path)
}

// RemoveAll ...
func (fileManager) RemoveAll(path string) error {
	return os.RemoveAll(path)
}

// Write ...
func (fileManager) Write(path string, value string, mode os.FileMode) error {
	if err := ensureSavePath(path); err != nil {
		return err
	}

	if err := fileutilV1.WriteStringToFile(path, value); err != nil {
		return err
	}

	if err := os.Chmod(path, mode); err != nil {
		return err
	}
	return nil
}

func ensureSavePath(savePath string) error {
	dirPath := filepath.Dir(savePath)
	return os.MkdirAll(dirPath, 0700)
}
