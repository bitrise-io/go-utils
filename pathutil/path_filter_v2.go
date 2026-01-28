package pathutil

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func FilterPathsV2(fsys fs.FS, filters ...FilterFuncV2) ([]string, error) {
	var filtered []string

	fs.WalkDir(fsys, ".", func(pth string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		allowed := true
		for _, filter := range filters {
			matches, err := filter(pth, d)
			if err != nil {
				return err
			}
			if !matches {
				allowed = false
				break
			}
		}
		if allowed {
			filtered = append(filtered, pth)
		}

		return nil
	})

	return filtered, nil
}

// pth string, d fs.DirEntry
type FilterFuncV2 func(string, fs.DirEntry) (bool, error)

func BaseFilterV2(base string, _ fs.DirEntry, allowed bool) FilterFuncV2 {
	return func(pth string, _ fs.DirEntry) (bool, error) {
		b := filepath.Base(pth)
		return allowed == strings.EqualFold(base, b), nil
	}
}

func ExtensionFilterV2(ext string, _ fs.DirEntry, allowed bool) FilterFuncV2 {
	return func(pth string, _ fs.DirEntry) (bool, error) {
		e := filepath.Ext(pth)
		return allowed == strings.EqualFold(ext, e), nil
	}
}

func RegexpFilterV2(pattern string, _ fs.DirEntry, allowed bool) FilterFuncV2 {
	return func(pth string, _ fs.DirEntry) (bool, error) {
		re := regexp.MustCompile(pattern)
		found := re.FindString(pth) != ""
		return allowed == found, nil
	}
}

func ComponentFilterV2(component string, _ fs.DirEntry, allowed bool) FilterFuncV2 {
	return func(pth string, _ fs.DirEntry) (bool, error) {
		found := false
		pathComponents := strings.Split(pth, string(filepath.Separator))
		for _, c := range pathComponents {
			if c == component {
				found = true
			}
		}
		return allowed == found, nil
	}
}

func ComponentWithExtensionFilterV2(ext string, allowed bool) FilterFunc {
	return func(pth string) (bool, error) {
		found := false
		pathComponents := strings.Split(pth, string(filepath.Separator))
		for _, c := range pathComponents {
			e := filepath.Ext(c)
			if e == ext {
				found = true
			}
		}
		return allowed == found, nil
	}
}

func IsDirectoryFilterV2(allowed bool) FilterFunc {
	return func(pth string) (bool, error) {
		fileInf, err := os.Lstat(pth)
		if err != nil {
			return false, err
		}
		if fileInf == nil {
			return false, errors.New("no file info available")
		}
		return allowed == fileInf.IsDir(), nil
	}
}

func InDirectoryFilterV2(dir string, allowed bool) FilterFunc {
	return func(pth string) (bool, error) {
		in := filepath.Dir(pth) == dir
		return allowed == in, nil
	}
}

// DirectoryContainsFileFilterV2 returns a FilterFunc that checks if a directory contains a file
func DirectoryContainsFileFilterV2(fileName string) FilterFunc {
	return func(pth string) (bool, error) {
		isDir, err := IsDirectoryFilter(true)(pth)
		if err != nil {
			return false, err
		}
		if !isDir {
			return false, nil
		}

		absPath := filepath.Join(pth, fileName)
		if _, err := os.Lstat(absPath); err != nil {
			if !os.IsNotExist(err) {
				return false, err
			}
			return false, nil
		}
		return true, nil
	}
}

func FileContainsFilterV2(pth string, _ fs.DirEntry, str string) (bool, error) {
	bytes, err := os.ReadFile(pth)
	if err != nil {
		return false, err
	}

	return strings.Contains(string(bytes), str), nil
}
