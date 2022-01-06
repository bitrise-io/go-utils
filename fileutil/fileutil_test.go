package fileutil

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/bitrise-io/go-utils/v2/pathutil"
	"github.com/stretchr/testify/require"
)

func TestWriteStringToFile(t *testing.T) {
	provider := pathutil.NewPathProvider()
	tmpDirPath, err := provider.CreateTempDir("go-utils-test-")
	require.NoError(t, err)
	manager := fileManager{}

	t.Log("success test")
	{
		tmpFilePath := filepath.Join(tmpDirPath, "WriteStringToFile-success.txt")
		require.NoError(t, manager.writeStringToFile(tmpFilePath, "test string"))

		fileContent, err := ioutil.ReadFile(tmpFilePath)
		require.NoError(t, err)
		require.Equal(t, "test string", string(fileContent))
	}

	t.Log("error test")
	{
		tmpFilePath := filepath.Join(tmpDirPath, "dir-does-not-exist", "WriteStringToFile-error.txt")
		require.Error(t, manager.writeStringToFile(tmpFilePath, "test string"), "open "+tmpFilePath+": no such file or directory")
	}
}
