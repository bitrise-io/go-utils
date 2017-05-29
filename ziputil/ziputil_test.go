package ziputil

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/bitrise-io/go-utils/fileutil"
	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/stretchr/testify/require"
)

func TestZip(t *testing.T) {
	t.Log("create zip from file")
	{
		tmpDir, err := pathutil.NormalizedOSTempDirPath("test")
		require.NoError(t, err)

		sourceFile := filepath.Join(tmpDir, "sourceFile")
		require.NoError(t, fileutil.WriteStringToFile(sourceFile, ""))

		destinationZip := filepath.Join(tmpDir, "destinationFile.zip")
		require.NoError(t, ZipFile(sourceFile, destinationZip))

		exist, err := pathutil.IsPathExists(destinationZip)
		require.NoError(t, err)
		require.Equal(t, true, exist)
	}

	t.Log("create zip from dir")
	{
		tmpDir, err := pathutil.NormalizedOSTempDirPath("test")
		require.NoError(t, err)

		sourceDir := filepath.Join(tmpDir, "sourceDir")
		require.NoError(t, os.MkdirAll(sourceDir, 0777))

		destinationZip := filepath.Join(tmpDir, "destinationDir.zip")
		require.NoError(t, ZipDir(sourceDir, destinationZip, false))

		exist, err := pathutil.IsPathExists(destinationZip)
		require.NoError(t, err)
		require.Equal(t, true, exist, destinationZip)
	}

	t.Log("zip content of dir")
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
	}
}

func TestUnZip(t *testing.T) {
	t.Log("unzip zipped file")
	{
		tmpDir, err := pathutil.NormalizedOSTempDirPath("test")
		require.NoError(t, err)

		// create file to zip tmp/source/sourceFile
		contentOfDirToZip := filepath.Join(tmpDir, "source")
		require.NoError(t, os.MkdirAll(contentOfDirToZip, 0777))

		sourceFile := filepath.Join(contentOfDirToZip, "sourceFile")
		require.NoError(t, fileutil.WriteStringToFile(sourceFile, ""))

		// create zip at tmp/destinationFile.zip
		destinationZip := filepath.Join(tmpDir, "destinationFile.zip")
		require.NoError(t, ZipFile(sourceFile, destinationZip))

		// unzip into tmp/
		require.NoError(t, UnZip(destinationZip, tmpDir))

		// tmp/sourceFile should exist
		content, err := fileutil.ReadStringFromFile(filepath.Join(tmpDir, "sourceFile"))
		require.NoError(t, err)
		require.Equal(t, "", content)
	}

	t.Log("unzip zipped dir")
	{
		tmpDir, err := pathutil.NormalizedOSTempDirPath("test")
		require.NoError(t, err)

		// create dir to zip tmp/source/sourceDir
		contentOfDirToZip := filepath.Join(tmpDir, "source")
		require.NoError(t, os.MkdirAll(contentOfDirToZip, 0777))

		sourceDir := filepath.Join(tmpDir, "sourceDir")
		require.NoError(t, os.MkdirAll(sourceDir, 0777))

		// create zip at tmp/destinationDir.zip
		destinationZip := filepath.Join(contentOfDirToZip, "destinationDir.zip")
		require.NoError(t, ZipDir(sourceDir, destinationZip, false))

		// unzip into tmp/
		require.NoError(t, UnZip(destinationZip, tmpDir))

		// tmp/sourceDir should exist
		isDir, err := pathutil.IsDirExists(filepath.Join(tmpDir, "sourceDir"))
		require.NoError(t, err)
		require.Equal(t, true, isDir)
	}

	t.Log("unzip zipped content of dir")
	{
		tmpDir, err := pathutil.NormalizedOSTempDirPath("test")
		require.NoError(t, err)

		// create dir to zip tmp/source/sourceDir
		contentOfDirToZip := filepath.Join(tmpDir, "source")
		require.NoError(t, os.MkdirAll(contentOfDirToZip, 0777))

		sourceDir := filepath.Join(contentOfDirToZip, "sourceDir")
		require.NoError(t, os.MkdirAll(sourceDir, 0777))

		// create file to zip tmp/source/sourceFile
		sourceFile := filepath.Join(contentOfDirToZip, "sourceFile")
		require.NoError(t, fileutil.WriteStringToFile(sourceFile, ""))

		// create zip at tmp/destinationDir.zip
		destinationZip := filepath.Join(tmpDir, "destinationFile.zip")
		require.NoError(t, ZipDir(contentOfDirToZip, destinationZip, true))

		// remove tmp/source, since this path would be the unzip destination
		require.NoError(t, os.RemoveAll(contentOfDirToZip))

		// unzip into tmp/
		require.NoError(t, UnZip(destinationZip, tmpDir))

		// tmp/sourceDir should exist
		isDir, err := pathutil.IsDirExists(filepath.Join(tmpDir, "sourceDir"))
		require.NoError(t, err)
		require.Equal(t, true, isDir)

		// tmp/sourceFile should exist
		content, err := fileutil.ReadStringFromFile(filepath.Join(tmpDir, "sourceFile"))
		require.NoError(t, err)
		require.Equal(t, "", content)
	}
}
