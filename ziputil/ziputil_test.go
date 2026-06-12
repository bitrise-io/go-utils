package ziputil

import (
	"archive/zip"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/bitrise-io/go-utils/fileutil"
	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/stretchr/testify/require"
)

func TestZipFiles(t *testing.T) {
	t.Log("files from different directories are stored under their base names (junk paths)")
	{
		tmpDir, err := pathutil.NormalizedOSTempDirPath("test")
		require.NoError(t, err)

		var sourceFilePaths []string
		for _, name := range []string{"A", "B", "C"} {
			baseDir := filepath.Join(tmpDir, name)
			require.NoError(t, pathutil.EnsureDirExist(baseDir))

			sourceFile := filepath.Join(baseDir, "sourceFile"+name)
			require.NoError(t, fileutil.WriteStringToFile(sourceFile, name))
			sourceFilePaths = append(sourceFilePaths, sourceFile)
		}

		destinationZip := filepath.Join(tmpDir, "dest.zip")
		require.NoError(t, ZipFiles(sourceFilePaths, destinationZip))

		unzipDir, err := pathutil.NormalizedOSTempDirPath("unzip")
		require.NoError(t, err)
		require.NoError(t, UnZip(destinationZip, unzipDir))

		for _, name := range []string{"A", "B", "C"} {
			content, err := fileutil.ReadStringFromFile(filepath.Join(unzipDir, "sourceFile"+name))
			require.NoError(t, err)
			require.Equal(t, name, content)
		}
	}

	t.Log("duplicate base names cause an error")
	{
		tmpDir, err := pathutil.NormalizedOSTempDirPath("test")
		require.NoError(t, err)

		var sourceFilePaths []string
		for _, name := range []string{"A", "B"} {
			baseDir := filepath.Join(tmpDir, name)
			require.NoError(t, pathutil.EnsureDirExist(baseDir))

			sourceFile := filepath.Join(baseDir, "sourceFile")
			require.NoError(t, fileutil.WriteStringToFile(sourceFile, name))
			sourceFilePaths = append(sourceFilePaths, sourceFile)
		}

		destinationZip := filepath.Join(tmpDir, "destinationFile.zip")
		require.Error(t, ZipFiles(sourceFilePaths, destinationZip))
	}
}

func TestZipDirectories(t *testing.T) {
	t.Log("dirs at different parent paths are each stored under their own basename")
	{
		tmpDir, err := pathutil.NormalizedOSTempDirPath("test")
		require.NoError(t, err)

		dirA := filepath.Join(tmpDir, "main", "A")
		require.NoError(t, os.MkdirAll(dirA, 0755))
		require.NoError(t, fileutil.WriteStringToFile(filepath.Join(dirA, "fileA.txt"), "from_a"))

		subDir := filepath.Join(tmpDir, "sub")
		require.NoError(t, os.MkdirAll(subDir, 0755))
		dirB := filepath.Join(subDir, "B")
		require.NoError(t, os.MkdirAll(dirB, 0755))
		require.NoError(t, fileutil.WriteStringToFile(filepath.Join(dirB, "fileB.txt"), "from_b"))

		destinationZip := filepath.Join(tmpDir, "dest.zip")
		require.NoError(t, ZipDirs([]string{dirA, dirB}, destinationZip))

		unzipDir, err := pathutil.NormalizedOSTempDirPath("unzip")
		require.NoError(t, err)
		require.NoError(t, UnZip(destinationZip, unzipDir))

		contentA, err := fileutil.ReadStringFromFile(filepath.Join(unzipDir, "A", "fileA.txt"))
		require.NoError(t, err)
		require.Equal(t, "from_a", contentA)

		contentB, err := fileutil.ReadStringFromFile(filepath.Join(unzipDir, "B", "fileB.txt"))
		require.NoError(t, err)
		require.Equal(t, "from_b", contentB)
	}
}

func TestUnZipFile(t *testing.T) {
	t.Log("unzip zipped file")
	{
		tmpDir, err := pathutil.NormalizedOSTempDirPath("test")
		require.NoError(t, err)

		sourceFile := filepath.Join(tmpDir, "source", "sourceFile")
		require.NoError(t, os.MkdirAll(filepath.Dir(sourceFile), 0755))
		require.NoError(t, fileutil.WriteStringToFile(sourceFile, ""))

		destinationZip := filepath.Join(tmpDir, "destinationFile.zip")
		require.NoError(t, ZipFile(sourceFile, destinationZip))

		require.NoError(t, UnZip(destinationZip, tmpDir))

		content, err := fileutil.ReadStringFromFile(filepath.Join(tmpDir, "sourceFile"))
		require.NoError(t, err)
		require.Equal(t, "", content)
	}

	t.Log("unzip zipped files")
	{
		tmpDir, err := pathutil.NormalizedOSTempDirPath("test")
		require.NoError(t, err)

		sourceDir := filepath.Join(tmpDir, "source")
		require.NoError(t, os.MkdirAll(sourceDir, 0755))

		var sourceFilePaths []string
		for _, name := range []string{"A", "B", "C"} {
			sourceFile := filepath.Join(sourceDir, "sourceFile"+name)
			require.NoError(t, fileutil.WriteStringToFile(sourceFile, ""))
			sourceFilePaths = append(sourceFilePaths, sourceFile)
		}

		destinationZip := filepath.Join(tmpDir, "destinationFile.zip")
		require.NoError(t, ZipFiles(sourceFilePaths, destinationZip))

		require.NoError(t, UnZip(destinationZip, tmpDir))

		for _, path := range sourceFilePaths {
			content, err := fileutil.ReadStringFromFile(filepath.Join(tmpDir, filepath.Base(path)))
			require.NoError(t, err)
			require.Equal(t, "", content)
		}
	}
}

func TestUnZipDirectory(t *testing.T) {
	t.Log("unzip zipped dir")
	{
		tmpDir, err := pathutil.NormalizedOSTempDirPath("test")
		require.NoError(t, err)

		sourceDir := filepath.Join(tmpDir, "sourceDir")
		require.NoError(t, os.MkdirAll(sourceDir, 0777))

		destinationZip := filepath.Join(tmpDir, "destinationDir.zip")
		require.NoError(t, ZipDir(sourceDir, destinationZip, false))

		require.NoError(t, UnZip(destinationZip, tmpDir))

		isDir, err := pathutil.IsDirExists(filepath.Join(tmpDir, "sourceDir"))
		require.NoError(t, err)
		require.Equal(t, true, isDir)
	}

	t.Log("unzip zipped directories")
	{
		tmpDir, err := pathutil.NormalizedOSTempDirPath("zip")
		require.NoError(t, err)

		mainDir := filepath.Join(tmpDir, "main")
		dirA := filepath.Join(mainDir, "A")
		require.NoError(t, os.MkdirAll(dirA, 0777))
		dirB := filepath.Join(mainDir, "B")
		require.NoError(t, os.MkdirAll(dirB, 0777))

		destinationZip := filepath.Join(tmpDir, "destinationDir.zip")
		require.NoError(t, ZipDirs([]string{dirA, dirB}, destinationZip))

		destTmpDir, err := pathutil.NormalizedOSTempDirPath("unzip")
		require.NoError(t, err)

		require.NoError(t, UnZip(destinationZip, destTmpDir))

		exist, err := pathutil.IsPathExists(filepath.Join(destTmpDir, "A"))
		require.NoError(t, err)
		require.Equal(t, true, exist)

		exist, err = pathutil.IsPathExists(filepath.Join(destTmpDir, "B"))
		require.NoError(t, err)
		require.Equal(t, true, exist)
	}

	t.Log("unzip zipped content of dir")
	{
		tmpDir, err := pathutil.NormalizedOSTempDirPath("test")
		require.NoError(t, err)

		contentOfDirToZip := filepath.Join(tmpDir, "source")
		require.NoError(t, os.MkdirAll(contentOfDirToZip, 0777))

		sourceDir := filepath.Join(contentOfDirToZip, "sourceDir")
		require.NoError(t, os.MkdirAll(sourceDir, 0777))

		sourceFile := filepath.Join(contentOfDirToZip, "sourceFile")
		require.NoError(t, fileutil.WriteStringToFile(sourceFile, ""))

		destinationZip := filepath.Join(tmpDir, "destinationFile.zip")
		require.NoError(t, ZipDir(contentOfDirToZip, destinationZip, true))

		require.NoError(t, os.RemoveAll(contentOfDirToZip))

		require.NoError(t, UnZip(destinationZip, tmpDir))

		isDir, err := pathutil.IsDirExists(filepath.Join(tmpDir, "sourceDir"))
		require.NoError(t, err)
		require.Equal(t, true, isDir)

		content, err := fileutil.ReadStringFromFile(filepath.Join(tmpDir, "sourceFile"))
		require.NoError(t, err)
		require.Equal(t, "", content)
	}

	t.Log("unzip into different dir")
	{
		sourceTmpDir, err := pathutil.NormalizedOSTempDirPath("__1__")
		require.NoError(t, err)

		sourceFile := filepath.Join(sourceTmpDir, "sourceFile")
		require.NoError(t, fileutil.WriteStringToFile(sourceFile, ""))

		destinationZip := filepath.Join(sourceTmpDir, "destinationFile.zip")
		require.NoError(t, ZipFile(sourceFile, destinationZip))

		destTmpDir, err := pathutil.NormalizedOSTempDirPath("__2__")
		require.NoError(t, err)

		require.NoError(t, UnZip(destinationZip, destTmpDir))
		exist, err := pathutil.IsPathExists(filepath.Join(destTmpDir, "sourceFile"))
		require.NoError(t, err)
		require.Equal(t, true, exist)
	}

	t.Log("relative path")
	{
		sourceTmpDir, err := pathutil.NormalizedOSTempDirPath("__1__")
		require.NoError(t, err)

		revokeFn, err := pathutil.RevokableChangeDir(sourceTmpDir)
		require.NoError(t, err)
		defer func() {
			require.NoError(t, revokeFn())
		}()

		sourceFile := filepath.Join(sourceTmpDir, "sourceFile")
		require.NoError(t, fileutil.WriteStringToFile(sourceFile, ""))

		require.NoError(t, ZipFile("./sourceFile", "./destinationFile.zip"))

		require.NoError(t, UnZip("./destinationFile.zip", "./unzipped"))
		exist, err := pathutil.IsPathExists("./unzipped/sourceFile")
		require.NoError(t, err)
		require.Equal(t, true, exist)

		destTmpDir, err := pathutil.NormalizedOSTempDirPath("__2__")
		require.NoError(t, err)

		require.NoError(t, UnZip("./destinationFile.zip", destTmpDir))
		exist, err = pathutil.IsPathExists(filepath.Join(destTmpDir, "sourceFile"))
		require.NoError(t, err)
		require.Equal(t, true, exist)

		require.NoError(t, revokeFn())
	}
}

func TestZipDirPreservesDirMtime(t *testing.T) {
	tmpDir, err := pathutil.NormalizedOSTempDirPath("test")
	require.NoError(t, err)

	sourceDir := filepath.Join(tmpDir, "sourceDir")
	require.NoError(t, os.MkdirAll(sourceDir, 0755))

	knownTime := time.Date(2020, 1, 1, 12, 0, 0, 0, time.UTC)
	require.NoError(t, os.Chtimes(sourceDir, knownTime, knownTime))

	destinationZip := filepath.Join(tmpDir, "out.zip")
	require.NoError(t, ZipDir(sourceDir, destinationZip, false))

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
	tmpDir, err := pathutil.NormalizedOSTempDirPath("test")
	require.NoError(t, err)

	// /a/shared and /b/shared share the basename "shared".
	dirA := filepath.Join(tmpDir, "a", "shared")
	dirB := filepath.Join(tmpDir, "b", "shared")
	require.NoError(t, os.MkdirAll(dirA, 0755))
	require.NoError(t, os.MkdirAll(dirB, 0755))

	require.NoError(t, fileutil.WriteStringToFile(filepath.Join(dirA, "only_in_a.txt"), "from_a"))
	require.NoError(t, fileutil.WriteStringToFile(filepath.Join(dirB, "only_in_b.txt"), "from_b"))
	require.NoError(t, fileutil.WriteStringToFile(filepath.Join(dirA, "conflict.txt"), "first"))
	require.NoError(t, fileutil.WriteStringToFile(filepath.Join(dirB, "conflict.txt"), "second"))

	destinationZip := filepath.Join(tmpDir, "out.zip")
	require.NoError(t, ZipDirs([]string{dirA, dirB}, destinationZip))

	unzipDir, err := pathutil.NormalizedOSTempDirPath("unzip")
	require.NoError(t, err)
	require.NoError(t, UnZip(destinationZip, unzipDir))

	// Files unique to each directory are both present.
	onlyA, err := fileutil.ReadStringFromFile(filepath.Join(unzipDir, "shared", "only_in_a.txt"))
	require.NoError(t, err)
	require.Equal(t, "from_a", onlyA)

	onlyB, err := fileutil.ReadStringFromFile(filepath.Join(unzipDir, "shared", "only_in_b.txt"))
	require.NoError(t, err)
	require.Equal(t, "from_b", onlyB)

	// Conflicting file resolves to the last directory in the list.
	conflict, err := fileutil.ReadStringFromFile(filepath.Join(unzipDir, "shared", "conflict.txt"))
	require.NoError(t, err)
	require.Equal(t, "second", conflict)
}

// TestZipFilesDuplicateBasenameNoOutputCreated verifies that when two source files share
// the same base name, ZipFiles returns an error and does not create the destination file.
func TestZipFilesDuplicateBasenameNoOutputCreated(t *testing.T) {
	tmpDir, err := pathutil.NormalizedOSTempDirPath("test")
	require.NoError(t, err)

	dirA := filepath.Join(tmpDir, "a")
	dirB := filepath.Join(tmpDir, "b")
	require.NoError(t, os.MkdirAll(dirA, 0755))
	require.NoError(t, os.MkdirAll(dirB, 0755))

	require.NoError(t, fileutil.WriteStringToFile(filepath.Join(dirA, "file.txt"), "a"))
	require.NoError(t, fileutil.WriteStringToFile(filepath.Join(dirB, "file.txt"), "b"))

	destinationZip := filepath.Join(tmpDir, "out.zip")
	require.Error(t, ZipFiles([]string{
		filepath.Join(dirA, "file.txt"),
		filepath.Join(dirB, "file.txt"),
	}, destinationZip))

	// The destination file must not be left behind after the failure.
	exist, err := pathutil.IsPathExists(destinationZip)
	require.NoError(t, err)
	require.False(t, exist, "destination zip must not be created when duplicate base names are detected")
}

// TestUnZipZipSlipHandling verifies the behaviour for archives with path-traversal entries.
// The system unzip strips the "../" components, extracts the file inside intoDir,
// and returns a non-zero exit code (which UnZip surfaces as an error).
func TestUnZipZipSlipHandling(t *testing.T) {
	tmpDir, err := pathutil.NormalizedOSTempDirPath("test")
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

	// UnZip returns an error because unzip exits non-zero for path-traversal warnings.
	err = UnZip(evilZip, destDir)
	require.Error(t, err)

	// Despite the error, the file is extracted inside destDir with the traversal stripped.
	exist, statErr := pathutil.IsPathExists(filepath.Join(destDir, "escape.txt"))
	require.NoError(t, statErr)
	require.True(t, exist, "system unzip strips '../' and extracts inside dest")

	// The file must not have escaped outside destDir. The entry has two "../" components
	// and destDir is one level below tmpDir, so a real escape would land at the parent of
	// tmpDir; that is the path to check.
	escapeTarget := filepath.Join(filepath.Dir(tmpDir), "escape.txt")
	escaped, statErr := pathutil.IsPathExists(escapeTarget)
	require.NoError(t, statErr)
	require.False(t, escaped, "file must not escape outside intoDir")
}
