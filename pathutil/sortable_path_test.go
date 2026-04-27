package pathutil

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"
)

// These tests assume POSIX-style absolute paths (no drive letter, "/" as
// separator). filepath.Abs and component splitting differ on Windows, so
// the bare-string assertions would not hold there. CI for go-utils runs on
// Linux and macOS only; skip on Windows to keep the test honest.
func skipOnWindows(t *testing.T) {
	t.Helper()
	if runtime.GOOS == "windows" {
		t.Skip("POSIX-only: filepath.Abs and component split differ on Windows")
	}
}

func TestNewSortablePath(t *testing.T) {
	skipOnWindows(t)

	sp, err := NewSortablePath("/a/b/c.txt")
	require.NoError(t, err)
	require.Equal(t, "/a/b/c.txt", sp.Pth)
	require.Equal(t, "/a/b/c.txt", sp.AbsPth)
	require.Equal(t, []string{"a", "b", "c.txt"}, sp.Components)
}

func TestSortPathsByComponents(t *testing.T) {
	skipOnWindows(t)

	input := []string{
		"/a/b/c/deep.txt",
		"/x.txt",
		"/a/mid.txt",
		"/a/b/other.txt",
	}

	got, err := SortPathsByComponents(input)
	require.NoError(t, err)
	require.Equal(t, []string{
		"/x.txt",
		"/a/mid.txt",
		"/a/b/other.txt",
		"/a/b/c/deep.txt",
	}, got)
}

func TestSortPathsByComponents_tiebreakAlphabetic(t *testing.T) {
	skipOnWindows(t)

	input := []string{"/a/zeta.txt", "/a/alpha.txt", "/a/mu.txt"}

	got, err := SortPathsByComponents(input)
	require.NoError(t, err)
	require.Equal(t, []string{"/a/alpha.txt", "/a/mu.txt", "/a/zeta.txt"}, got)
}

func TestSortPathsByComponents_tiebreakSameBaseDifferentDir(t *testing.T) {
	skipOnWindows(t)

	// Same depth, same base name — without the AbsPth tie-breaker the
	// resulting order would be nondeterministic (sort.Sort is not stable).
	input := []string{"/b/x.txt", "/a/x.txt"}

	got, err := SortPathsByComponents(input)
	require.NoError(t, err)
	require.Equal(t, []string{"/a/x.txt", "/b/x.txt"}, got)
}

func TestListPathInDirSortedByComponents(t *testing.T) {
	root := t.TempDir()
	mustWrite(t, filepath.Join(root, "top.txt"), "x")
	mustWrite(t, filepath.Join(root, "sub", "mid.txt"), "x")
	mustWrite(t, filepath.Join(root, "sub", "deep", "leaf.txt"), "x")

	got, err := ListPathInDirSortedByComponents(root, true)
	require.NoError(t, err)

	// Expect shallowest entries first; root itself ("." after Rel) is included.
	require.Equal(t, ".", got[0])
	require.Contains(t, got, "top.txt")
	require.Contains(t, got, filepath.Join("sub", "deep", "leaf.txt"))
	// Deepest entry must come last.
	require.Equal(t, filepath.Join("sub", "deep", "leaf.txt"), got[len(got)-1])
}

func mustWrite(t *testing.T, pth, content string) {
	t.Helper()
	require.NoError(t, os.MkdirAll(filepath.Dir(pth), 0o755))
	require.NoError(t, os.WriteFile(pth, []byte(content), 0o644))
}
