package pathutil

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func writeTempTree(t *testing.T) string {
	t.Helper()

	root := t.TempDir()

	files := map[string]string{
		"README.md":            "# readme",
		"main.go":              "package main",
		"main_test.go":         "package main",
		"sub/note.txt":         "todo",
		"sub/nested/deep.txt":  "x",
		"vendor/pkg/lib.go":    "package pkg",
	}
	for rel, content := range files {
		full := filepath.Join(root, rel)
		require.NoError(t, os.MkdirAll(filepath.Dir(full), 0o755))
		require.NoError(t, os.WriteFile(full, []byte(content), 0o644))
	}
	return root
}

func TestFilterPaths(t *testing.T) {
	root := writeTempTree(t)
	paths := []string{
		filepath.Join(root, "README.md"),
		filepath.Join(root, "main.go"),
		filepath.Join(root, "main_test.go"),
		filepath.Join(root, "sub"),
	}

	got, err := FilterPaths(paths, ExtensionFilter(".go", true), IsDirectoryFilter(false))
	require.NoError(t, err)
	require.ElementsMatch(t, []string{
		filepath.Join(root, "main.go"),
		filepath.Join(root, "main_test.go"),
	}, got)
}

func TestFilterPaths_missingPath_pathOnlyFilters(t *testing.T) {
	// v1 parity: purely lexical filters work against paths that do not
	// exist on disk. IsDirectoryFilter is not used here so no stat is needed.
	got, err := FilterPaths(
		[]string{"./Podfile", "/nowhere/Podfile.lock"},
		BaseFilter("Podfile", true),
	)
	require.NoError(t, err)
	require.Equal(t, []string{"./Podfile"}, got)
}

func TestFilterPaths_missingPath_dirEntryFilterErrors(t *testing.T) {
	// Filters that consult the DirEntry (IsDirectoryFilter) still error on
	// missing paths because there is nothing to inspect.
	_, err := FilterPaths([]string{"/does/not/exist/ever"}, IsDirectoryFilter(false))
	require.Error(t, err)
}

func TestFilterPaths_empty(t *testing.T) {
	got, err := FilterPaths(nil, ExtensionFilter(".go", true))
	require.NoError(t, err)
	require.Empty(t, got)
}

func TestListEntries(t *testing.T) {
	root := writeTempTree(t)

	got, err := ListEntries(root, IsDirectoryFilter(false))
	require.NoError(t, err)
	require.ElementsMatch(t, []string{
		filepath.Join(root, "README.md"),
		filepath.Join(root, "main.go"),
		filepath.Join(root, "main_test.go"),
	}, got, "ListEntries must be non-recursive")
}

func TestListEntries_onlyDirs(t *testing.T) {
	root := writeTempTree(t)

	got, err := ListEntries(root, IsDirectoryFilter(true))
	require.NoError(t, err)
	require.ElementsMatch(t, []string{
		filepath.Join(root, "sub"),
		filepath.Join(root, "vendor"),
	}, got)
}
