package fileutil

import (
	"os"
	"path/filepath"
)

// FileManager ...
type FileManager interface {
	Open(path string) (*os.File, error)
	Remove(path string) error
	RemoveAll(path string) error
	Write(path string, value string, perm os.FileMode) error
	WriteBytes(path string, value []byte) error
}

type fileManager struct {
}

// NewFileManager ...
func NewFileManager() FileManager {
	return fileManager{}
}

// Open ...
func (fileManager) Open(path string) (*os.File, error) {
	return os.Open(path)
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
func (f fileManager) Write(path string, value string, mode os.FileMode) error {
	if err := f.ensureSavePath(path); err != nil {
		return err
	}
	if err := os.WriteFile(path, []byte(value), mode); err != nil {
		return err
	}
	return os.Chmod(path, mode)
}

func (fileManager) ensureSavePath(savePath string) error {
	dirPath := filepath.Dir(savePath)
	return os.MkdirAll(dirPath, 0700)
}

// WriteBytes ...
func (f fileManager) WriteBytes(path string, value []byte) error {
	return os.WriteFile(path, value, 0600)
}
