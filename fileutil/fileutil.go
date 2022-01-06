package fileutil

import (
	"errors"
	"os"
	"path/filepath"
)

// FileManager ...
type FileManager interface {
	Remove(path string) error
	RemoveAll(path string) error
	Write(path string, value string, mode os.FileMode) error
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
func (f fileManager) Write(path string, value string, mode os.FileMode) error {
	if err := f.ensureSavePath(path); err != nil {
		return err
	}

	if err := f.writeStringToFile(path, value); err != nil {
		return err
	}

	if err := os.Chmod(path, mode); err != nil {
		return err
	}
	return nil
}

func (fileManager) ensureSavePath(savePath string) error {
	dirPath := filepath.Dir(savePath)
	return os.MkdirAll(dirPath, 0700)
}

func (f fileManager) writeStringToFile(pth string, fileCont string) (err error) {
	fc := []byte(fileCont)
	if pth == "" {
		return errors.New("no path provided")
	}

	var file *os.File
	file, err = os.Create(pth)
	if err != nil {
		return
	}

	defer func() {
		err2 := file.Close()
		if err == nil {
			err = err2
		}
	}()

	if _, err = file.Write(fc); err != nil {
		return
	}

	if err = file.Close(); err != nil {
		return
	}

	return
}
