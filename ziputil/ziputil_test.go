package ziputil_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/bitrise-io/go-utils/v2/pathutil"
	"github.com/bitrise-io/go-utils/v2/ziputil"
	"github.com/stretchr/testify/require"
)

func newManager() *ziputil.ZipManager {
	return ziputil.NewZipManager(pathutil.NewPathChecker())
}

func TestZipFile(t *testing.T) {
	provider := pathutil.NewPathProvider()
	tmpDir, err := provider.CreateTempDir("test")
	require.NoError(t, err)

	sourceFile := filepath.Join(tmpDir, "sourceFile")
	require.NoError(t, os.WriteFile(sourceFile, []byte("hello"), 0644))

	destinationZip := filepath.Join(tmpDir, "dest.zip")
	require.NoError(t, newManager().ZipFile(sourceFile, destinationZip))

	checker := pathutil.NewPathChecker()
	exist, err := checker.IsPathExists(destinationZip)
	require.NoError(t, err)
	require.True(t, exist)
}

func TestZipFiles(t *testing.T) {
	t.Run("files in multiple directories", func(t *testing.T) {
		provider := pathutil.NewPathProvider()
		tmpDir, err := provider.CreateTempDir("test")
		require.NoError(t, err)

		var sourceFilePths []string
		for _, name := range []string{"A", "B", "C"} {
			baseDir := filepath.Join(tmpDir, name)
			require.NoError(t, os.MkdirAll(baseDir, 0755))

			sourceFile := filepath.Join(baseDir, "sourceFile"+name)
			require.NoError(t, os.WriteFile(sourceFile, []byte(name), 0644))
			sourceFilePths = append(sourceFilePths, sourceFile)
		}

		destinationZip := filepath.Join(tmpDir, "dest.zip")
		require.NoError(t, newManager().ZipFiles(sourceFilePths, destinationZip))

		checker := pathutil.NewPathChecker()
		exist, err := checker.IsPathExists(destinationZip)
		require.NoError(t, err)
		require.True(t, exist)
	})

	t.Run("files in the same directory", func(t *testing.T) {
		provider := pathutil.NewPathProvider()
		tmpDir, err := provider.CreateTempDir("test")
		require.NoError(t, err)

		var sourceFilePths []string
		for _, name := range []string{"A", "B", "C"} {
			sourceFile := filepath.Join(tmpDir, "sourceFile"+name)
			require.NoError(t, os.WriteFile(sourceFile, []byte(name), 0644))
			sourceFilePths = append(sourceFilePths, sourceFile)
		}

		destinationZip := filepath.Join(tmpDir, "dest.zip")
		require.NoError(t, newManager().ZipFiles(sourceFilePths, destinationZip))

		checker := pathutil.NewPathChecker()
		exist, err := checker.IsPathExists(destinationZip)
		require.NoError(t, err)
		require.True(t, exist)
	})

	t.Run("duplicate base names return error", func(t *testing.T) {
		provider := pathutil.NewPathProvider()
		tmpDir, err := provider.CreateTempDir("test")
		require.NoError(t, err)

		var sourceFilePths []string
		for _, name := range []string{"A", "B"} {
			baseDir := filepath.Join(tmpDir, name)
			require.NoError(t, os.MkdirAll(baseDir, 0755))

			// Same base name in different directories
			sourceFile := filepath.Join(baseDir, "sourceFile")
			require.NoError(t, os.WriteFile(sourceFile, []byte(name), 0644))
			sourceFilePths = append(sourceFilePths, sourceFile)
		}

		destinationZip := filepath.Join(tmpDir, "dest.zip")
		require.Error(t, newManager().ZipFiles(sourceFilePths, destinationZip))
	})

	t.Run("non-existent file returns error", func(t *testing.T) {
		provider := pathutil.NewPathProvider()
		tmpDir, err := provider.CreateTempDir("test")
		require.NoError(t, err)

		err = newManager().ZipFiles([]string{filepath.Join(tmpDir, "nonexistent")}, filepath.Join(tmpDir, "dest.zip"))
		require.Error(t, err)
	})
}

func TestZipDir(t *testing.T) {
	t.Run("zip dir with directory entry", func(t *testing.T) {
		provider := pathutil.NewPathProvider()
		tmpDir, err := provider.CreateTempDir("test")
		require.NoError(t, err)

		sourceDir := filepath.Join(tmpDir, "sourceDir")
		require.NoError(t, os.MkdirAll(sourceDir, 0755))

		destinationZip := filepath.Join(tmpDir, "dest.zip")
		require.NoError(t, newManager().ZipDir(sourceDir, destinationZip, false))

		checker := pathutil.NewPathChecker()
		exist, err := checker.IsPathExists(destinationZip)
		require.NoError(t, err)
		require.True(t, exist)
	})

	t.Run("zip content of dir", func(t *testing.T) {
		provider := pathutil.NewPathProvider()
		tmpDir, err := provider.CreateTempDir("test")
		require.NoError(t, err)

		sourceDir := filepath.Join(tmpDir, "source")
		require.NoError(t, os.MkdirAll(sourceDir, 0755))

		subDir := filepath.Join(sourceDir, "subDir")
		require.NoError(t, os.MkdirAll(subDir, 0755))
		require.NoError(t, os.WriteFile(filepath.Join(subDir, "file.txt"), []byte("content"), 0644))

		destinationZip := filepath.Join(tmpDir, "dest.zip")
		require.NoError(t, newManager().ZipDir(sourceDir, destinationZip, true))

		checker := pathutil.NewPathChecker()
		exist, err := checker.IsPathExists(destinationZip)
		require.NoError(t, err)
		require.True(t, exist)
	})

	t.Run("non-existent dir returns error", func(t *testing.T) {
		provider := pathutil.NewPathProvider()
		tmpDir, err := provider.CreateTempDir("test")
		require.NoError(t, err)

		err = newManager().ZipDir(filepath.Join(tmpDir, "nonexistent"), filepath.Join(tmpDir, "dest.zip"), false)
		require.Error(t, err)
	})
}

func TestZipDirs(t *testing.T) {
	t.Run("dirs in different parent directories", func(t *testing.T) {
		provider := pathutil.NewPathProvider()
		tmpDir, err := provider.CreateTempDir("test")
		require.NoError(t, err)

		dirA := filepath.Join(tmpDir, "main", "A")
		require.NoError(t, os.MkdirAll(dirA, 0755))

		dirB := filepath.Join(tmpDir, "sub", "B")
		require.NoError(t, os.MkdirAll(dirB, 0755))

		destinationZip := filepath.Join(tmpDir, "dest.zip")
		require.NoError(t, newManager().ZipDirs([]string{dirA, dirB}, destinationZip))

		checker := pathutil.NewPathChecker()
		exist, err := checker.IsPathExists(destinationZip)
		require.NoError(t, err)
		require.True(t, exist)
	})

	t.Run("dirs in the same parent directory", func(t *testing.T) {
		provider := pathutil.NewPathProvider()
		tmpDir, err := provider.CreateTempDir("test")
		require.NoError(t, err)

		mainDir := filepath.Join(tmpDir, "main")
		dirA := filepath.Join(mainDir, "A")
		require.NoError(t, os.MkdirAll(dirA, 0755))
		dirB := filepath.Join(mainDir, "B")
		require.NoError(t, os.MkdirAll(dirB, 0755))

		destinationZip := filepath.Join(tmpDir, "dest.zip")
		require.NoError(t, newManager().ZipDirs([]string{dirA, dirB}, destinationZip))

		checker := pathutil.NewPathChecker()
		exist, err := checker.IsPathExists(destinationZip)
		require.NoError(t, err)
		require.True(t, exist)
	})

	t.Run("non-existent dir returns error", func(t *testing.T) {
		provider := pathutil.NewPathProvider()
		tmpDir, err := provider.CreateTempDir("test")
		require.NoError(t, err)

		err = newManager().ZipDirs([]string{filepath.Join(tmpDir, "nonexistent")}, filepath.Join(tmpDir, "dest.zip"))
		require.Error(t, err)
	})
}

func TestUnZipFile(t *testing.T) {
	t.Run("round-trip single file", func(t *testing.T) {
		provider := pathutil.NewPathProvider()
		tmpDir, err := provider.CreateTempDir("test")
		require.NoError(t, err)

		sourceFile := filepath.Join(tmpDir, "sourceFile")
		require.NoError(t, os.WriteFile(sourceFile, []byte("hello"), 0644))

		destinationZip := filepath.Join(tmpDir, "dest.zip")
		require.NoError(t, newManager().ZipFile(sourceFile, destinationZip))

		unzipDir := filepath.Join(tmpDir, "unzipped")
		require.NoError(t, newManager().UnZip(destinationZip, unzipDir))

		content, err := os.ReadFile(filepath.Join(unzipDir, "sourceFile"))
		require.NoError(t, err)
		require.Equal(t, "hello", string(content))
	})

	t.Run("round-trip multiple files", func(t *testing.T) {
		provider := pathutil.NewPathProvider()
		tmpDir, err := provider.CreateTempDir("test")
		require.NoError(t, err)

		var sourceFilePths []string
		for _, name := range []string{"A", "B", "C"} {
			sourceFile := filepath.Join(tmpDir, "sourceFile"+name)
			require.NoError(t, os.WriteFile(sourceFile, []byte(name), 0644))
			sourceFilePths = append(sourceFilePths, sourceFile)
		}

		destinationZip := filepath.Join(tmpDir, "dest.zip")
		require.NoError(t, newManager().ZipFiles(sourceFilePths, destinationZip))

		unzipDir := filepath.Join(tmpDir, "unzipped")
		require.NoError(t, newManager().UnZip(destinationZip, unzipDir))

		for _, pth := range sourceFilePths {
			content, err := os.ReadFile(filepath.Join(unzipDir, filepath.Base(pth)))
			require.NoError(t, err)
			require.Equal(t, filepath.Base(pth)[len("sourceFile"):], string(content))
		}
	})
}

func TestUnZipDirectory(t *testing.T) {
	t.Run("round-trip directory", func(t *testing.T) {
		provider := pathutil.NewPathProvider()
		tmpDir, err := provider.CreateTempDir("test")
		require.NoError(t, err)

		sourceDir := filepath.Join(tmpDir, "sourceDir")
		require.NoError(t, os.MkdirAll(sourceDir, 0755))

		destinationZip := filepath.Join(tmpDir, "dest.zip")
		require.NoError(t, newManager().ZipDir(sourceDir, destinationZip, false))

		unzipDir := filepath.Join(tmpDir, "unzipped")
		require.NoError(t, newManager().UnZip(destinationZip, unzipDir))

		checker := pathutil.NewPathChecker()
		isDir, err := checker.IsDirExists(filepath.Join(unzipDir, "sourceDir"))
		require.NoError(t, err)
		require.True(t, isDir)
	})

	t.Run("round-trip ZipDirs — each dir under its basename", func(t *testing.T) {
		provider := pathutil.NewPathProvider()
		tmpDir, err := provider.CreateTempDir("test")
		require.NoError(t, err)

		mainDir := filepath.Join(tmpDir, "main")
		dirA := filepath.Join(mainDir, "A")
		require.NoError(t, os.MkdirAll(dirA, 0755))
		dirB := filepath.Join(mainDir, "B")
		require.NoError(t, os.MkdirAll(dirB, 0755))

		destinationZip := filepath.Join(tmpDir, "dest.zip")
		require.NoError(t, newManager().ZipDirs([]string{dirA, dirB}, destinationZip))

		unzipDir := filepath.Join(tmpDir, "unzipped")
		require.NoError(t, newManager().UnZip(destinationZip, unzipDir))

		checker := pathutil.NewPathChecker()
		for _, name := range []string{"A", "B"} {
			exist, err := checker.IsPathExists(filepath.Join(unzipDir, name))
			require.NoError(t, err)
			require.True(t, exist, "expected %s to exist after unzip", name)
		}
	})

	t.Run("round-trip content-only zip", func(t *testing.T) {
		provider := pathutil.NewPathProvider()
		tmpDir, err := provider.CreateTempDir("test")
		require.NoError(t, err)

		contentDir := filepath.Join(tmpDir, "source")
		require.NoError(t, os.MkdirAll(contentDir, 0755))
		subDir := filepath.Join(contentDir, "subDir")
		require.NoError(t, os.MkdirAll(subDir, 0755))
		require.NoError(t, os.WriteFile(filepath.Join(contentDir, "topFile"), []byte("top"), 0644))

		destinationZip := filepath.Join(tmpDir, "dest.zip")
		require.NoError(t, newManager().ZipDir(contentDir, destinationZip, true))

		require.NoError(t, os.RemoveAll(contentDir))

		unzipDir := tmpDir
		require.NoError(t, newManager().UnZip(destinationZip, unzipDir))

		checker := pathutil.NewPathChecker()
		isDir, err := checker.IsDirExists(filepath.Join(unzipDir, "subDir"))
		require.NoError(t, err)
		require.True(t, isDir)

		content, err := os.ReadFile(filepath.Join(unzipDir, "topFile"))
		require.NoError(t, err)
		require.Equal(t, "top", string(content))
	})

	t.Run("unzip into different dir", func(t *testing.T) {
		provider := pathutil.NewPathProvider()
		srcTmpDir, err := provider.CreateTempDir("src")
		require.NoError(t, err)
		dstTmpDir, err := provider.CreateTempDir("dst")
		require.NoError(t, err)

		sourceFile := filepath.Join(srcTmpDir, "sourceFile")
		require.NoError(t, os.WriteFile(sourceFile, []byte(""), 0644))

		destinationZip := filepath.Join(srcTmpDir, "dest.zip")
		require.NoError(t, newManager().ZipFile(sourceFile, destinationZip))

		require.NoError(t, newManager().UnZip(destinationZip, dstTmpDir))

		checker := pathutil.NewPathChecker()
		exist, err := checker.IsPathExists(filepath.Join(dstTmpDir, "sourceFile"))
		require.NoError(t, err)
		require.True(t, exist)
	})
}

func TestSymlinks(t *testing.T) {
	t.Run("symlink is stored and restored", func(t *testing.T) {
		provider := pathutil.NewPathProvider()
		tmpDir, err := provider.CreateTempDir("test")
		require.NoError(t, err)

		sourceDir := filepath.Join(tmpDir, "source")
		require.NoError(t, os.MkdirAll(sourceDir, 0755))

		realFile := filepath.Join(sourceDir, "real.txt")
		require.NoError(t, os.WriteFile(realFile, []byte("real content"), 0644))

		linkFile := filepath.Join(sourceDir, "link.txt")
		require.NoError(t, os.Symlink("real.txt", linkFile))

		destinationZip := filepath.Join(tmpDir, "dest.zip")
		require.NoError(t, newManager().ZipDir(sourceDir, destinationZip, false))

		unzipDir := filepath.Join(tmpDir, "unzipped")
		require.NoError(t, newManager().UnZip(destinationZip, unzipDir))

		restoredLink := filepath.Join(unzipDir, "source", "link.txt")
		target, err := os.Readlink(restoredLink)
		require.NoError(t, err)
		require.Equal(t, "real.txt", target)
	})
}

func TestRelativePath(t *testing.T) {
	provider := pathutil.NewPathProvider()
	tmpDir, err := provider.CreateTempDir("test")
	require.NoError(t, err)

	require.NoError(t, os.Chdir(tmpDir))

	sourceFile := filepath.Join(tmpDir, "sourceFile")
	require.NoError(t, os.WriteFile(sourceFile, []byte(""), 0644))

	require.NoError(t, newManager().ZipFile("./sourceFile", "./dest.zip"))

	unzipDir := filepath.Join(tmpDir, "unzipped")
	require.NoError(t, newManager().UnZip("./dest.zip", unzipDir))

	checker := pathutil.NewPathChecker()
	exist, err := checker.IsPathExists(filepath.Join(unzipDir, "sourceFile"))
	require.NoError(t, err)
	require.True(t, exist)
}
