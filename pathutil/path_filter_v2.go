package pathutil

import (
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

// ToDo: separate Allow and Deny filter func types
type AllowFilterFuncV2 func(string, fs.DirEntry) (bool, error)
type DenyFilterFuncV2 func(string, fs.DirEntry) (bool, error)

// NOTE: NEW Completely skip the directory if matched
func SkipDirectoryFilterV2(dir string) FilterFuncV2 {
	return func(pth string, d fs.DirEntry) (bool, error) {
		if !d.IsDir() {
			return true, nil
		}

		dirPath := filepath.Dir(pth)
		isPathMatch := strings.EqualFold(dir, dirPath)
		if isPathMatch {
			return false, fs.SkipDir
		}

		return true, nil
	}
}

func BaseFilterV2(base string, allowed bool) FilterFuncV2 {
	return func(pth string, _ fs.DirEntry) (bool, error) {
		b := filepath.Base(pth)
		return allowed == strings.EqualFold(base, b), nil
	}
}

func ExtensionFilterV2(ext string, allowed bool) FilterFuncV2 {
	return func(pth string, _ fs.DirEntry) (bool, error) {
		e := filepath.Ext(pth)
		return allowed == strings.EqualFold(ext, e), nil
	}
}

func RegexpFilterV2(pattern string, allowed bool) FilterFuncV2 {
	return func(pth string, _ fs.DirEntry) (bool, error) {
		re := regexp.MustCompile(pattern)
		found := re.FindString(pth) != ""
		return allowed == found, nil
	}
}

// func GlobFilterV2(pattern string, allowed bool) FilterFuncV2 {
// 	return func(pth string, _ fs.DirEntry) (bool, error) {
// 		matched, err := filepath.Match(pattern, pth)
// 		if err != nil {
// 			return false, err
// 		}
// 		return allowed == matched, nil
// 	}
// }

func ComponentFilterV2(component string, allowed bool) FilterFuncV2 {
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

func ComponentWithExtensionFilterV2(ext string, allowed bool) FilterFuncV2 {
	return func(pth string, _ fs.DirEntry) (bool, error) {
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

// NOTE: simplifed relative to V1 version, as fs.DirEntry provides IsDir info
func IsDirectoryFilterV2(allowed bool) FilterFuncV2 {
	return func(pth string, d fs.DirEntry) (bool, error) {
		return allowed == d.IsDir(), nil
	}
}

// NOTE: simplifed relative to V1 version, as fs.DirEntry provides IsDir info
func InDirectoryFilterV2(dir string, allowed bool) FilterFuncV2 {
	return func(pth string, _ fs.DirEntry) (bool, error) {
		in := filepath.Dir(pth) == dir
		return allowed == in, nil
	}
}

// NOTE: simplified checking for directory using fs.DirEntry
// DirectoryContainsFileFilterV2 returns a FilterFunc that checks if a directory contains a file
func DirectoryContainsFileFilterV2(fileName string) FilterFuncV2 {
	return func(pth string, d fs.DirEntry) (bool, error) {
		if !d.IsDir() {
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

func FileContainsFilterV2(pth string, d fs.DirEntry, str string) (bool, error) {
	if d.IsDir() { // Note: added check for directory
		return false, nil
	}

	bytes, err := os.ReadFile(pth)
	if err != nil {
		return false, err
	}

	return strings.Contains(string(bytes), str), nil
}
