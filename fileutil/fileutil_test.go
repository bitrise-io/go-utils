package fileutil

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/stretchr/testify/require"
)

func TestWriteStringToFile(t *testing.T) {
	tmpDirPath, err := pathutil.NormalizedOSTempDirPath("go-utils-test-")
	require.NoError(t, err)

	t.Log("simple string test")
	{
		tmpFilePath := filepath.Join(tmpDirPath, "WriteStringToFile.txt")
		require.NoError(t, WriteStringToFile(tmpFilePath, "test string"))

		fileContent, err := ioutil.ReadFile(tmpFilePath)
		require.NoError(t, err)
		require.Equal(t, "test string", string(fileContent))
	}
}

func TestWriteJSONToFile(t *testing.T) {
	tmpDirPath, err := pathutil.NormalizedOSTempDirPath("go-utils-test-")
	require.NoError(t, err)

	t.Log("map test")
	{
		testContent := map[string]interface{}{
			"root": map[string]string{
				"key1": "value1",
			},
		}
		tmpFilePath := filepath.Join(tmpDirPath, "WriteJSONToFile.json")
		require.NoError(t, WriteJSONToFile(tmpFilePath, testContent))

		fileContent, err := ioutil.ReadFile(tmpFilePath)
		require.NoError(t, err)
		require.Equal(t, `{"root":{"key1":"value1"}}`, string(fileContent))
	}
}
