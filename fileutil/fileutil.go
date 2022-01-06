package fileutil

import (
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
)

// FileManager ...
type FileManager interface {
	Remove(path string) error
	RemoveAll(path string) error
	Write(path string, value string, perm fs.FileMode) error
}

type fileManager struct {
}

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
func (f fileManager) Write(path string, value string, mode fs.FileMode) error {
	if err := f.ensureSavePath(path); err != nil {
		return err
	}
	return ioutil.WriteFile(path, []byte(value), mode)
}

func (fileManager) ensureSavePath(savePath string) error {
	dirPath := filepath.Dir(savePath)
	return os.MkdirAll(dirPath, 0700)
}
