package fileutil

import (
	"io/ioutil"
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

	t.Log("success test")
	{
		tmpFilePath := filepath.Join(tmpDirPath, "WriteStringToFile-success.txt")
		require.NoError(t, manager.Write(tmpFilePath, "test string", 0600))

		fileContent, err := ioutil.ReadFile(tmpFilePath)
		require.NoError(t, err)
		require.Equal(t, "test string", string(fileContent))
	}

	t.Log("error test")
	{
		tmpFilePath := filepath.Join(tmpDirPath, "dir-does-not-exist", "WriteStringToFile-error.txt")
		require.Error(t, manager.Write(tmpFilePath, "test string", 0600), "open "+tmpFilePath+": no such file or directory")
	}
}
