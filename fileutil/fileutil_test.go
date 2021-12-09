package fileutil

import (
	"io/ioutil"
	"math"
	"path/filepath"
	"testing"

	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/stretchr/testify/require"
)

func TestWriteStringToFile(t *testing.T) {
	tmpDirPath, err := pathutil.NormalizedOSTempDirPath("go-utils-test-")
	require.NoError(t, err)

	t.Log("success test")
	{
		tmpFilePath := filepath.Join(tmpDirPath, "WriteStringToFile-success.txt")
		require.NoError(t, WriteStringToFile(tmpFilePath, "test string"))

		fileContent, err := ioutil.ReadFile(tmpFilePath)
		require.NoError(t, err)
		require.Equal(t, "test string", string(fileContent))
	}

	t.Log("error test")
	{
		tmpFilePath := filepath.Join(tmpDirPath, "dir-does-not-exist", "WriteStringToFile-error.txt")
		require.Error(t, WriteStringToFile(tmpFilePath, "test string"), "open "+tmpFilePath+": no such file or directory")
	}
}

func TestWriteJSONToFile(t *testing.T) {
	tmpDirPath, err := pathutil.NormalizedOSTempDirPath("go-utils-test-")
	require.NoError(t, err)

	t.Log("success test")
	{
		testContent := map[string]interface{}{
			"root": map[string]string{
				"key1": "value1",
			},
		}
		tmpFilePath := filepath.Join(tmpDirPath, "WriteJSONToFile-success.json")
		require.NoError(t, WriteJSONToFile(tmpFilePath, testContent))

		fileContent, err := ioutil.ReadFile(tmpFilePath)
		require.NoError(t, err)
		require.Equal(t, `{"root":{"key1":"value1"}}`, string(fileContent))
	}

	t.Log("error test")
	{
		testContent := map[string]interface{}{
			"root": math.Inf(1),
		}
		tmpFilePath := filepath.Join(tmpDirPath, "WriteJSONToFile-error.json")
		require.Error(t, WriteJSONToFile(tmpFilePath, testContent), "failed to JSON marshal the provided object: json: unsupported value: +Inf")

		t.Log("File should not exist if a JSON marshaling error happened")
		exists, err := pathutil.IsPathExists(tmpFilePath)
		require.NoError(t, err)
		require.Equal(t, false, exists, "file should not exist at: %s", tmpFilePath)
	}
}
