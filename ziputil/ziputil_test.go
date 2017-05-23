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
	tmpDir, err := pathutil.NormalizedOSTempDirPath("test")
	require.NoError(t, err)

	t.Log("create zip from file")
	{
		sourceFile := filepath.Join(tmpDir, "sourceFile")
		require.NoError(t, fileutil.WriteStringToFile(sourceFile, ""))

		destinationZip := filepath.Join(tmpDir, "destinationFile.zip")
		require.NoError(t, Zip(sourceFile, destinationZip))

		exist, err := pathutil.IsPathExists(destinationZip)
		require.NoError(t, err)
		require.Equal(t, true, exist)
	}

	t.Log("create zip from dir")
	{
		sourceDir := filepath.Join(tmpDir, "sourceDir")
		require.NoError(t, os.MkdirAll(sourceDir, 0777))

		destinationZip := filepath.Join(tmpDir, "destinationDir.zip")
		require.NoError(t, Zip(sourceDir, destinationZip))

		exist, err := pathutil.IsPathExists(destinationZip)
		require.NoError(t, err)
		require.Equal(t, true, exist, destinationZip)
	}
}

func TestUnZip(t *testing.T) {
	tmpDir, err := pathutil.NormalizedOSTempDirPath("test")
	require.NoError(t, err)

	t.Log("unzip zipped file")
	{
		sourceFile := filepath.Join(tmpDir, "sourceFile")
		require.NoError(t, fileutil.WriteStringToFile(sourceFile, ""))

		destinationZip := filepath.Join(tmpDir, "destinationFile.zip")
		require.NoError(t, Zip(sourceFile, destinationZip))

		unziped, err := UnZip(destinationZip)
		require.NoError(t, err)

		content, err := fileutil.ReadStringFromFile(unziped)
		require.NoError(t, err)
		require.Equal(t, "", content)
	}

	t.Log("unzip zipped dir")
	{
		sourceDir := filepath.Join(tmpDir, "sourceDir")
		require.NoError(t, os.MkdirAll(sourceDir, 0777))

		destinationZip := filepath.Join(tmpDir, "destinationDir.zip")
		require.NoError(t, Zip(sourceDir, destinationZip))

		unziped, err := UnZip(destinationZip)
		require.NoError(t, err)

		isDir, err := pathutil.IsDirExists(unziped)
		require.NoError(t, err)
		require.Equal(t, true, isDir)
	}
}
