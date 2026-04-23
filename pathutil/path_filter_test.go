package pathutil

import (
	"errors"
	"io/fs"
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/require"
)

func buildTestFS() fstest.MapFS {
	return fstest.MapFS{
		"src/main.go":             {Data: []byte("package main\nfunc main(){}\n")},
		"src/util/helper.go":      {Data: []byte("package util\n// TODO: refactor\n")},
		"src/util/helper_test.go": {Data: []byte("package util\n")},
		"docs/readme.md":          {Data: []byte("# readme\n")},
		"vendor/pkg/lib.go":       {Data: []byte("package pkg\n")},
		"Info.plist":              {Data: []byte("plist")},
	}
}

func TestBaseFilter(t *testing.T) {
	fsys := buildTestFS()

	got, err := FilterFS(fsys, BaseFilter("main.go", true), IsDirectoryFilter(false))
	require.NoError(t, err)
	require.Equal(t, []string{"src/main.go"}, got)

	got, err = FilterFS(fsys, BaseFilter("MAIN.GO", true), IsDirectoryFilter(false))
	require.NoError(t, err, "case-insensitive match expected")
	require.Equal(t, []string{"src/main.go"}, got)

	got, err = FilterFS(fsys, BaseFilter("main.go", false), IsDirectoryFilter(false))
	require.NoError(t, err)
	require.NotContains(t, got, "src/main.go")
}

func TestExtensionFilter(t *testing.T) {
	fsys := buildTestFS()

	got, err := FilterFS(fsys, ExtensionFilter(".md", true), IsDirectoryFilter(false))
	require.NoError(t, err)
	require.Equal(t, []string{"docs/readme.md"}, got)

	got, err = FilterFS(fsys, ExtensionFilter(".MD", true), IsDirectoryFilter(false))
	require.NoError(t, err)
	require.Equal(t, []string{"docs/readme.md"}, got)
}

func TestRegexpFilter(t *testing.T) {
	fsys := buildTestFS()

	got, err := FilterFS(fsys, RegexpFilter(`.*_test\.go$`, true), IsDirectoryFilter(false))
	require.NoError(t, err)
	require.Equal(t, []string{"src/util/helper_test.go"}, got)

	got, err = FilterFS(fsys, RegexpFilter(`.*_test\.go$`, false), IsDirectoryFilter(false))
	require.NoError(t, err)
	require.NotContains(t, got, "src/util/helper_test.go")
}

func TestComponentFilter(t *testing.T) {
	fsys := buildTestFS()

	got, err := FilterFS(fsys, ComponentFilter("vendor", true))
	require.NoError(t, err)
	require.Equal(t, []string{"vendor", "vendor/pkg", "vendor/pkg/lib.go"}, got)

	got, err = FilterFS(fsys, ComponentFilter("vendor", false), IsDirectoryFilter(false))
	require.NoError(t, err)
	require.NotContains(t, got, "vendor/pkg/lib.go")
}

func TestComponentWithExtensionFilter(t *testing.T) {
	fsys := fstest.MapFS{
		"proj.xcodeproj/project.pbxproj": {Data: []byte("x")},
		"src/main.go":                    {Data: []byte("x")},
	}

	got, err := FilterFS(fsys, ComponentWithExtensionFilter(".xcodeproj", true), IsDirectoryFilter(false))
	require.NoError(t, err)
	require.Equal(t, []string{"proj.xcodeproj/project.pbxproj"}, got)
}

func TestIsDirectoryFilter(t *testing.T) {
	fsys := fstest.MapFS{
		"a/b.txt": {Data: []byte("x")},
	}

	dirs, err := FilterFS(fsys, IsDirectoryFilter(true))
	require.NoError(t, err)
	require.ElementsMatch(t, []string{".", "a"}, dirs)

	files, err := FilterFS(fsys, IsDirectoryFilter(false))
	require.NoError(t, err)
	require.Equal(t, []string{"a/b.txt"}, files)
}

func TestInDirectoryFilter(t *testing.T) {
	fsys := buildTestFS()

	got, err := FilterFS(fsys, InDirectoryFilter("src", true), IsDirectoryFilter(false))
	require.NoError(t, err)
	require.Equal(t, []string{"src/main.go"}, got)

	got, err = FilterFS(fsys, InDirectoryFilter("src/util", true), IsDirectoryFilter(false))
	require.NoError(t, err)
	require.ElementsMatch(t, []string{"src/util/helper.go", "src/util/helper_test.go"}, got)
}

func TestSkipDirectoryNameFilter(t *testing.T) {
	fsys := buildTestFS()

	got, err := FilterFS(fsys, SkipDirectoryNameFilter("vendor"), IsDirectoryFilter(false))
	require.NoError(t, err)
	require.NotContains(t, got, "vendor/pkg/lib.go", "vendor subtree must be skipped")
	require.Contains(t, got, "src/main.go")
}

func TestDirectoryContainsFileFilter(t *testing.T) {
	fsys := fstest.MapFS{
		"proj-a/Podfile":        {Data: []byte("pod 'X'")},
		"proj-a/main.m":         {Data: []byte("x")},
		"proj-b/main.m":         {Data: []byte("x")},
		"proj-c/Podfile/nested": {Data: []byte("x")}, // "Podfile" is a dir here
	}

	got, err := FilterFS(fsys, DirectoryContainsFileFilter(fsys, "Podfile"), IsDirectoryFilter(true))
	require.NoError(t, err)
	require.Equal(t, []string{"proj-a"}, got)
}

func TestFileContainsFilter(t *testing.T) {
	fsys := buildTestFS()

	got, err := FilterFS(fsys, FileContainsFilter(fsys, "TODO"), IsDirectoryFilter(false))
	require.NoError(t, err)
	require.Equal(t, []string{"src/util/helper.go"}, got)
}

func TestFilterFS_propagatesError(t *testing.T) {
	sentinel := errors.New("boom")
	bad := func(string, fs.DirEntry) (bool, error) {
		return false, sentinel
	}

	_, err := FilterFS(fstest.MapFS{"a": {Data: []byte("x")}}, bad)
	require.ErrorIs(t, err, sentinel)
}

func TestFilterFS_composite(t *testing.T) {
	fsys := buildTestFS()

	got, err := FilterFS(
		fsys,
		SkipDirectoryNameFilter("vendor"),
		IsDirectoryFilter(false),
		ExtensionFilter(".go", true),
		RegexpFilter(`.*_test\.go$`, false),
	)
	require.NoError(t, err)
	require.ElementsMatch(t, []string{"src/main.go", "src/util/helper.go"}, got)
}
