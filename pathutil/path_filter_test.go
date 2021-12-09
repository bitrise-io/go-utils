package pathutil_test

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bitrise-io/go-utils/fileutil"
	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/stretchr/testify/require"
)

func TestFilterPaths(t *testing.T) {
	t.Log("without any filter")
	{
		paths := []string{
			"/Users/bitrise/test",
			"/Users/vagrant/test",
		}
		filtered, err := pathutil.FilterPaths(paths)
		require.NoError(t, err)
		require.Equal(t, paths, filtered)
	}

	t.Log("with filter")
	{
		paths := []string{
			"/Users/bitrise/test",
			"/Users/vagrant/test",
		}
		filter := func(pth string) (bool, error) {
			return strings.Contains(pth, "vagrant"), nil
		}
		filtered, err := pathutil.FilterPaths(paths, filter)
		require.NoError(t, err)
		require.Equal(t, []string{"/Users/vagrant/test"}, filtered)
	}
}

func TestBaseFilter(t *testing.T) {
	t.Log("allow")
	{
		paths := []string{
			"path/to/my/gradlew",
			"path/to/my/gradlew/file",
		}
		filter := pathutil.BaseFilter("gradlew", true)
		filtered, err := pathutil.FilterPaths(paths, filter)
		require.NoError(t, err)
		require.Equal(t, []string{"path/to/my/gradlew"}, filtered)
	}

	t.Log("forbid")
	{
		paths := []string{
			"path/to/my/gradlew",
			"path/to/my/gradlew/file",
		}
		filter := pathutil.BaseFilter("gradlew", false)
		filtered, err := pathutil.FilterPaths(paths, filter)
		require.NoError(t, err)
		require.Equal(t, []string{"path/to/my/gradlew/file"}, filtered)
	}
}

func TestExtensionFilter(t *testing.T) {
	t.Log("allow")
	{
		paths := []string{
			"path/to/my/project.xcodeproj",
			"path/to/my/project.xcworkspace",
		}
		filter := pathutil.ExtensionFilter(".xcodeproj", true)
		filtered, err := pathutil.FilterPaths(paths, filter)
		require.NoError(t, err)
		require.Equal(t, []string{"path/to/my/project.xcodeproj"}, filtered)
	}

	t.Log("forbid")
	{
		paths := []string{
			"path/to/my/project.xcodeproj",
			"path/to/my/project.xcworkspace",
		}
		filter := pathutil.ExtensionFilter(".xcodeproj", false)
		filtered, err := pathutil.FilterPaths(paths, filter)
		require.NoError(t, err)
		require.Equal(t, []string{"path/to/my/project.xcworkspace"}, filtered)
	}
}

func TestRegexpFilter(t *testing.T) {
	t.Log("allow")
	{
		paths := []string{
			"path/to/my/project.xcodeproj",
			"path/to/my/project.xcworkspace",
		}
		filter := pathutil.RegexpFilter(".*.xcodeproj", true)
		filtered, err := pathutil.FilterPaths(paths, filter)
		require.NoError(t, err)
		require.Equal(t, []string{"path/to/my/project.xcodeproj"}, filtered)
	}

	t.Log("forbid")
	{
		paths := []string{
			"path/to/my/project.xcodeproj",
			"path/to/my/project.xcworkspace",
		}
		filter := pathutil.RegexpFilter(".*.xcodeproj", false)
		filtered, err := pathutil.FilterPaths(paths, filter)
		require.NoError(t, err)
		require.Equal(t, []string{"path/to/my/project.xcworkspace"}, filtered)
	}
}

func TestComponentFilter(t *testing.T) {
	t.Log("allow")
	{
		paths := []string{
			"/Users/bitrise/test",
			"/Users/vagrant/test",
		}
		filter := pathutil.ComponentFilter("bitrise", true)
		filtered, err := pathutil.FilterPaths(paths, filter)
		require.NoError(t, err)
		require.Equal(t, []string{"/Users/bitrise/test"}, filtered)
	}

	t.Log("forbid")
	{
		paths := []string{
			"/Users/bitrise/test",
			"/Users/vagrant/test",
		}
		filter := pathutil.ComponentFilter("bitrise", false)
		filtered, err := pathutil.FilterPaths(paths, filter)
		require.NoError(t, err)
		require.Equal(t, []string{"/Users/vagrant/test"}, filtered)
	}
}

func TestComponentWithExtensionFilter(t *testing.T) {
	t.Log("allow")
	{
		paths := []string{
			"/Users/bitrise.framework/test",
			"/Users/vagrant/test",
		}
		filter := pathutil.ComponentWithExtensionFilter(".framework", true)
		filtered, err := pathutil.FilterPaths(paths, filter)
		require.NoError(t, err)
		require.Equal(t, []string{"/Users/bitrise.framework/test"}, filtered)
	}

	t.Log("forbid")
	{
		paths := []string{
			"/Users/bitrise.framework/test",
			"/Users/vagrant/test",
		}
		filter := pathutil.ComponentWithExtensionFilter(".framework", false)
		filtered, err := pathutil.FilterPaths(paths, filter)
		require.NoError(t, err)
		require.Equal(t, []string{"/Users/vagrant/test"}, filtered)
	}
}

func TestIsDirectoryFilter(t *testing.T) {
	tmpDir, err := pathutil.NormalizedOSTempDirPath("__bitrise-init__")
	require.NoError(t, err)
	defer func() {
		require.NoError(t, os.RemoveAll(tmpDir))
	}()

	tmpFile := filepath.Join(tmpDir, "file.txt")
	require.NoError(t, fileutil.WriteStringToFile(tmpFile, ""))

	t.Log("allow")
	{
		paths := []string{
			tmpDir,
			tmpFile,
		}
		filter := pathutil.IsDirectoryFilter(true)
		filtered, err := pathutil.FilterPaths(paths, filter)
		require.NoError(t, err)
		require.Equal(t, []string{tmpDir}, filtered)
	}

	t.Log("forbid")
	{
		paths := []string{
			tmpDir,
			tmpFile,
		}
		filter := pathutil.IsDirectoryFilter(false)
		filtered, err := pathutil.FilterPaths(paths, filter)
		require.NoError(t, err)
		require.Equal(t, []string{tmpFile}, filtered)
	}
}

func TestInDirectoryFilter(t *testing.T) {
	t.Log("allow")
	{
		paths := []string{
			"/Users/bitrise/test",
			"/Users/vagrant/test",
		}
		filter := pathutil.InDirectoryFilter("/Users/bitrise", true)
		filtered, err := pathutil.FilterPaths(paths, filter)
		require.NoError(t, err)
		require.Equal(t, []string{"/Users/bitrise/test"}, filtered)
	}

	t.Log("forbid")
	{
		paths := []string{
			"/Users/bitrise/test",
			"/Users/vagrant/test",
		}
		filter := pathutil.InDirectoryFilter("/Users/bitrise", false)
		filtered, err := pathutil.FilterPaths(paths, filter)
		require.NoError(t, err)
		require.Equal(t, []string{"/Users/vagrant/test"}, filtered)
	}
}

func TestDirectoryContainsFileFilter(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "*.xcworkspace")
	if err != nil {
		t.Errorf("failed to create temp dir, error: %s", err)
	}

	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Errorf("failed to remove temp dir, error: %s", err)
		}
	}()

	const filterFileName = "contents.xcworkspacedata"
	tempFile, err := os.Create(path.Join(tempDir, filterFileName))
	if err != nil {
		t.Errorf("failed to create temp file, error: %s", err)
	}
	if err := tempFile.Close(); err != nil {
		t.Errorf("failed to close file, error: %s", err)
	}

	tests := []struct {
		name           string
		path           string
		filterFileName string
		want           bool
		wantErr        bool
	}{
		{
			name:           "contains file",
			path:           tempDir,
			filterFileName: filterFileName,
			want:           true,
			wantErr:        false,
		},
		{
			name:           "does not contain file",
			filterFileName: filterFileName + "asd",
			path:           tempDir,
			want:           false,
			wantErr:        false,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				got, err := pathutil.DirectoryContainsFileFilter(tt.filterFileName)(tt.path)
				if (err != nil) != tt.wantErr {
					t.Errorf("DirectoryContainsFile() returned error: %v, wantErr: %v", err, tt.wantErr)
				}
				if got != tt.want {
					t.Errorf("DirectoryContainsFile() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}
