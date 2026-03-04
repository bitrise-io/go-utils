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

func TestLastNLines(t *testing.T) {
	manager := NewFileManager()

	tests := []struct {
		name string
		s    string
		n    int
		want string
	}{
		{
			name: "n=0 returns empty string",
			s:    "line1\nline2\nline3",
			n:    0,
			want: "",
		},
		{
			name: "negative n returns empty string",
			s:    "line1\nline2\nline3",
			n:    -1,
			want: "",
		},
		{
			name: "empty string returns empty string",
			s:    "",
			n:    3,
			want: "",
		},
		{
			name: "all-newlines string returns empty string",
			s:    "\n\n\n",
			n:    2,
			want: "",
		},
		{
			name: "last 1 line",
			s:    "line1\nline2\nline3",
			n:    1,
			want: "line3",
		},
		{
			name: "last 2 lines",
			s:    "line1\nline2\nline3",
			n:    2,
			want: "line2\nline3",
		},
		{
			name: "n equals total number of lines",
			s:    "line1\nline2\nline3",
			n:    3,
			want: "line1\nline2\nline3",
		},
		{
			name: "n greater than total lines returns full string",
			s:    "line1\nline2\nline3",
			n:    10,
			want: "line1\nline2\nline3",
		},
		{
			name: "trailing newline is ignored",
			s:    "line1\nline2\nline3\n",
			n:    2,
			want: "line2\nline3",
		},
		{
			name: "multiple trailing newlines are ignored",
			s:    "line1\nline2\nline3\n\n\n",
			n:    2,
			want: "line2\nline3",
		},
		{
			name: "CRLF line endings are normalized",
			s:    "line1\r\nline2\r\nline3\r\n",
			n:    2,
			want: "line2\nline3",
		},
		{
			name: "single line no newline",
			s:    "only line",
			n:    1,
			want: "only line",
		},
		{
			name: "single line with trailing spaces trimmed",
			s:    "line1\nline2   ",
			n:    1,
			want: "line2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := manager.LastNLines(tt.s, tt.n)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestFileSizeInBytes(t *testing.T) {
	provider := pathutil.NewPathProvider()
	tmpDirPath, err := provider.CreateTempDir("go-utils-test-")
	require.NoError(t, err)
	manager := NewFileManager()
	t.Run("Success when existing file provided", func(t *testing.T) {
		const content = "test string"
		{
			tmpFilePath := filepath.Join(tmpDirPath, "FileSizeInBytes-success.txt")
			require.NoError(t, manager.Write(tmpFilePath, content, 0600))

			fileSize, err := manager.FileSizeInBytes(tmpFilePath)
			require.NoError(t, err)
			require.Equal(t, int64(len([]byte(content))), fileSize)
		}
	})

	t.Run("Failure when non-existing path", func(t *testing.T) {
		tmpFilePath := filepath.Join(tmpDirPath, "dir-does-not-exist-2", "FileSizeInBytes-error.txt")
		_, err := manager.FileSizeInBytes(tmpFilePath)
		require.EqualError(t, err, "stat "+tmpFilePath+": no such file or directory")
	})

	t.Run("Failure when path is empty string", func(t *testing.T) {
		_, err := manager.FileSizeInBytes("")
		require.EqualError(t, err, "No path provided")
	})
}
