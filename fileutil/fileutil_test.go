package fileutil

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/bitrise-io/go-utils/v2/pathutil"
	"github.com/stretchr/testify/require"
)

func TestWrite(t *testing.T) {
	provider := pathutil.NewPathProvider()
	tmpDirPath, err := provider.CreateTempDir("go-utils-test-")
	require.NoError(t, err)
	manager := NewFileManager()

	t.Log("Success when writing string to existing dir")
	const content = "test string"
	{
		tmpFilePath := filepath.Join(tmpDirPath, "WriteStringToFile-success.txt")
		require.NoError(t, manager.Write(tmpFilePath, content, 0600))

		fileContent, err := os.ReadFile(tmpFilePath)
		require.NoError(t, err)
		require.Equal(t, content, string(fileContent))
	}

	t.Log("Success when writing string to non-existing dir")
	{
		tmpFilePath := filepath.Join(tmpDirPath, "dir-does-not-exist", "WriteStringToFile-success.txt")
		require.NoError(t, manager.Write(tmpFilePath, content, 0600))

		fileContent, err := os.ReadFile(tmpFilePath)
		require.NoError(t, err)
		require.Equal(t, content, string(fileContent))
	}
	t.Log("Success when writing bytes to existing dir")

	{
		tmpFilePath := filepath.Join(tmpDirPath, "WriteBytesToFile-success.txt")
		require.NoError(t, manager.WriteBytes(tmpFilePath, []byte("test string")))

		fileContent, err := os.ReadFile(tmpFilePath)
		require.NoError(t, err)
		require.Equal(t, "test string", string(fileContent))
	}

	t.Log("Failure when writing bytes to non-existing dir")
	{
		tmpFilePath := filepath.Join(tmpDirPath, "dir-does-not-exist-2", "WriteBytesToFile-error.txt")
		require.Error(t, manager.WriteBytes(tmpFilePath, []byte("test string")), "open "+tmpFilePath+": no such file or directory")
	}
}
