//go:build integration

package ziputil_test

import (
	"archive/zip"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/bitrise-io/go-utils/v2/pathutil"
	"github.com/stretchr/testify/require"
)

// TestZipFilesContents verifies that ZipFiles stores files under their base names and that
// the contents survive a round-trip through UnZip.
func TestZipFilesContents(t *testing.T) {
	t.Run("files from different directories are stored under their base names", func(t *testing.T) {
		tmpDir, err := pathutil.NewPathProvider().CreateTempDir("test")
		require.NoError(t, err)

		var sourceFilePaths []string
		for _, name := range []string{"A", "B", "C"} {
			baseDir := filepath.Join(tmpDir, name)
			require.NoError(t, os.MkdirAll(baseDir, 0755))
			sourceFile := filepath.Join(baseDir, "sourceFile"+name)
			require.NoError(t, os.WriteFile(sourceFile, []byte(name), 0644))
			sourceFilePaths = append(sourceFilePaths, sourceFile)
		}

		destinationZip := filepath.Join(tmpDir, "dest.zip")
		require.NoError(t, newManager().ZipFiles(sourceFilePaths, destinationZip))

		unzipDir, err := pathutil.NewPathProvider().CreateTempDir("unzip")
		require.NoError(t, err)
		require.NoError(t, newManager().UnZip(destinationZip, unzipDir))

		for _, name := range []string{"A", "B", "C"} {
			b, err := os.ReadFile(filepath.Join(unzipDir, "sourceFile"+name))
			require.NoError(t, err)
			require.Equal(t, name, string(b))
		}
	})

	t.Run("duplicate base names cause an error", func(t *testing.T) {
		tmpDir, err := pathutil.NewPathProvider().CreateTempDir("test")
		require.NoError(t, err)

		var sourceFilePaths []string
		for _, name := range []string{"A", "B"} {
			baseDir := filepath.Join(tmpDir, name)
			require.NoError(t, os.MkdirAll(baseDir, 0755))
			sourceFile := filepath.Join(baseDir, "sourceFile")
			require.NoError(t, os.WriteFile(sourceFile, []byte(name), 0644))
			sourceFilePaths = append(sourceFilePaths, sourceFile)
		}

		require.Error(t, newManager().ZipFiles(sourceFilePaths, filepath.Join(tmpDir, "dest.zip")))
	})
}

// TestZipDirsContents verifies that ZipDirs places each directory under its own basename and
// that file contents survive a round-trip through UnZip.
func TestZipDirsContents(t *testing.T) {
	t.Run("dirs at different parent paths are each stored under their own basename", func(t *testing.T) {
		tmpDir, err := pathutil.NewPathProvider().CreateTempDir("test")
		require.NoError(t, err)

		dirA := filepath.Join(tmpDir, "main", "A")
		require.NoError(t, os.MkdirAll(dirA, 0755))
		require.NoError(t, os.WriteFile(filepath.Join(dirA, "fileA.txt"), []byte("from_a"), 0644))

		dirB := filepath.Join(tmpDir, "sub", "B")
		require.NoError(t, os.MkdirAll(dirB, 0755))
		require.NoError(t, os.WriteFile(filepath.Join(dirB, "fileB.txt"), []byte("from_b"), 0644))

		destinationZip := filepath.Join(tmpDir, "dest.zip")
		require.NoError(t, newManager().ZipDirs([]string{dirA, dirB}, destinationZip))

		unzipDir, err := pathutil.NewPathProvider().CreateTempDir("unzip")
		require.NoError(t, err)
		require.NoError(t, newManager().UnZip(destinationZip, unzipDir))

		contentA, err := os.ReadFile(filepath.Join(unzipDir, "A", "fileA.txt"))
		require.NoError(t, err)
		require.Equal(t, "from_a", string(contentA))

		contentB, err := os.ReadFile(filepath.Join(unzipDir, "B", "fileB.txt"))
		require.NoError(t, err)
		require.Equal(t, "from_b", string(contentB))
	})
}

func TestZipDirPreservesDirMtime(t *testing.T) {
	tmpDir, err := pathutil.NewPathProvider().CreateTempDir("test")
	require.NoError(t, err)

	sourceDir := filepath.Join(tmpDir, "sourceDir")
	require.NoError(t, os.MkdirAll(sourceDir, 0755))

	knownTime := time.Date(2020, 1, 1, 12, 0, 0, 0, time.UTC)
	require.NoError(t, os.Chtimes(sourceDir, knownTime, knownTime))

	destinationZip := filepath.Join(tmpDir, "out.zip")
	require.NoError(t, newManager().ZipDir(sourceDir, destinationZip, false))

	r, err := zip.OpenReader(destinationZip)
	require.NoError(t, err)
	defer r.Close() //nolint:errcheck

	var found bool
	for _, f := range r.File {
		if f.Name == "sourceDir/" {
			found = true
			stored := f.Modified.UTC()
			require.Equal(t, knownTime.Year(), stored.Year())
			require.Equal(t, knownTime.Month(), stored.Month())
			require.Equal(t, knownTime.Day(), stored.Day())
			break
		}
	}
	require.True(t, found, "sourceDir/ entry not found in archive")
}

// TestZipDirsSameBasenamesMerge verifies the merge behaviour when two entries in
// sourceDirPths share the same basename. Files unique to either directory are preserved;
// conflicting files (same name) are resolved in favour of the last directory in the list.
func TestZipDirsSameBasenamesMerge(t *testing.T) {
	tmpDir, err := pathutil.NewPathProvider().CreateTempDir("test")
	require.NoError(t, err)

	// /a/shared and /b/shared share the basename "shared".
	dirA := filepath.Join(tmpDir, "a", "shared")
	dirB := filepath.Join(tmpDir, "b", "shared")
	require.NoError(t, os.MkdirAll(dirA, 0755))
	require.NoError(t, os.MkdirAll(dirB, 0755))

	require.NoError(t, os.WriteFile(filepath.Join(dirA, "only_in_a.txt"), []byte("from_a"), 0644))
	require.NoError(t, os.WriteFile(filepath.Join(dirB, "only_in_b.txt"), []byte("from_b"), 0644))
	require.NoError(t, os.WriteFile(filepath.Join(dirA, "conflict.txt"), []byte("first"), 0644))
	require.NoError(t, os.WriteFile(filepath.Join(dirB, "conflict.txt"), []byte("second"), 0644))

	destinationZip := filepath.Join(tmpDir, "out.zip")
	require.NoError(t, newManager().ZipDirs([]string{dirA, dirB}, destinationZip))

	unzipDir, err := pathutil.NewPathProvider().CreateTempDir("unzip")
	require.NoError(t, err)
	require.NoError(t, newManager().UnZip(destinationZip, unzipDir))

	onlyA, err := os.ReadFile(filepath.Join(unzipDir, "shared", "only_in_a.txt"))
	require.NoError(t, err)
	require.Equal(t, "from_a", string(onlyA))

	onlyB, err := os.ReadFile(filepath.Join(unzipDir, "shared", "only_in_b.txt"))
	require.NoError(t, err)
	require.Equal(t, "from_b", string(onlyB))

	// Conflicting file resolves to the last directory in the list.
	conflict, err := os.ReadFile(filepath.Join(unzipDir, "shared", "conflict.txt"))
	require.NoError(t, err)
	require.Equal(t, "second", string(conflict))
}

// TestZipFilesDuplicateBasenameNoOutputCreated verifies that when two source files share
// the same base name, ZipFiles returns an error and does not create the destination file.
func TestZipFilesDuplicateBasenameNoOutputCreated(t *testing.T) {
	tmpDir, err := pathutil.NewPathProvider().CreateTempDir("test")
	require.NoError(t, err)

	dirA := filepath.Join(tmpDir, "a")
	dirB := filepath.Join(tmpDir, "b")
	require.NoError(t, os.MkdirAll(dirA, 0755))
	require.NoError(t, os.MkdirAll(dirB, 0755))

	require.NoError(t, os.WriteFile(filepath.Join(dirA, "file.txt"), []byte("a"), 0644))
	require.NoError(t, os.WriteFile(filepath.Join(dirB, "file.txt"), []byte("b"), 0644))

	destinationZip := filepath.Join(tmpDir, "out.zip")
	require.Error(t, newManager().ZipFiles([]string{
		filepath.Join(dirA, "file.txt"),
		filepath.Join(dirB, "file.txt"),
	}, destinationZip))

	// The destination file must not be left behind after the failure.
	exist, err := pathutil.NewPathChecker().IsPathExists(destinationZip)
	require.NoError(t, err)
	require.False(t, exist, "destination zip must not be created when duplicate base names are detected")
}

// TestUnZipZipSlipHandling verifies the behaviour for archives with path-traversal entries.
// UnZip returns an error immediately on the first traversal entry; no file is extracted.
func TestUnZipZipSlipHandling(t *testing.T) {
	tmpDir, err := pathutil.NewPathProvider().CreateTempDir("test")
	require.NoError(t, err)

	evilZip := filepath.Join(tmpDir, "evil.zip")
	zf, err := os.Create(evilZip)
	require.NoError(t, err)
	zw := zip.NewWriter(zf)
	w, err := zw.CreateHeader(&zip.FileHeader{Name: "../../escape.txt", Method: zip.Store})
	require.NoError(t, err)
	_, err = w.Write([]byte("escaped"))
	require.NoError(t, err)
	require.NoError(t, zw.Close())
	require.NoError(t, zf.Close())

	destDir := filepath.Join(tmpDir, "dest")
	require.NoError(t, os.MkdirAll(destDir, 0755))

	err = newManager().UnZip(evilZip, destDir)
	require.Error(t, err)

	// Extraction aborted: no file is created inside destDir.
	exist, statErr := pathutil.NewPathChecker().IsPathExists(filepath.Join(destDir, "escape.txt"))
	require.NoError(t, statErr)
	require.False(t, exist, "traversal entry must not be extracted")

	// The file must not have escaped outside destDir either. The entry has two "../"
	// components and destDir is one level below tmpDir, so a real escape would land at the
	// parent of tmpDir; that is the path to check.
	escapeTarget := filepath.Join(filepath.Dir(tmpDir), "escape.txt")
	escaped, statErr := pathutil.NewPathChecker().IsPathExists(escapeTarget)
	require.NoError(t, statErr)
	require.False(t, escaped, "file must not escape outside intoDir")
}

// TestZipPreservesExistingDestinationOnFailure verifies that when zipping fails mid-write, a
// pre-existing destination archive is left untouched. This mirrors v1, where `zip -T` validated
// a temporary archive before overwriting the destination, leaving it in its prior state on failure.
func TestZipPreservesExistingDestinationOnFailure(t *testing.T) {
	if os.Geteuid() == 0 {
		t.Skip("an unreadable file cannot trigger a read failure when running as root")
	}

	tmpDir, err := pathutil.NewPathProvider().CreateTempDir("test")
	require.NoError(t, err)

	// A pre-existing destination archive with known content.
	destinationZip := filepath.Join(tmpDir, "out.zip")
	require.NoError(t, os.WriteFile(destinationZip, []byte("original"), 0644))

	// An unreadable source file forces ZipFiles to fail while writing the archive.
	unreadable := filepath.Join(tmpDir, "unreadable.txt")
	require.NoError(t, os.WriteFile(unreadable, []byte("secret"), 0000))
	defer os.Chmod(unreadable, 0644) //nolint:errcheck

	require.Error(t, newManager().ZipFiles([]string{unreadable}, destinationZip))

	// The pre-existing destination must be intact, and no temporary archive left behind.
	content, err := os.ReadFile(destinationZip)
	require.NoError(t, err)
	require.Equal(t, "original", string(content), "pre-existing destination must be preserved on failure")

	exist, err := pathutil.NewPathChecker().IsPathExists(destinationZip + ".tmp")
	require.NoError(t, err)
	require.False(t, exist, "temporary archive must be cleaned up on failure")
}
