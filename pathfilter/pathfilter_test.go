package pathfilter

import (
	"github.com/bitrise-io/go-utils/fileutil"
	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"
)

func TestFilterPaths(t *testing.T) {
	t.Log("without any filter")
	{
		paths := []string{
			"/Users/bitrise/test",
			"/Users/vagrant/test",
		}
		filtered, err := FilterPaths(paths)
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
		filtered, err := FilterPaths(paths, filter)
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
		filter := BaseFilter("gradlew", true)
		filtered, err := FilterPaths(paths, filter)
		require.NoError(t, err)
		require.Equal(t, []string{"path/to/my/gradlew"}, filtered)
	}

	t.Log("forbid")
	{
		paths := []string{
			"path/to/my/gradlew",
			"path/to/my/gradlew/file",
		}
		filter := BaseFilter("gradlew", false)
		filtered, err := FilterPaths(paths, filter)
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
		filter := ExtensionFilter(".xcodeproj", true)
		filtered, err := FilterPaths(paths, filter)
		require.NoError(t, err)
		require.Equal(t, []string{"path/to/my/project.xcodeproj"}, filtered)
	}

	t.Log("forbid")
	{
		paths := []string{
			"path/to/my/project.xcodeproj",
			"path/to/my/project.xcworkspace",
		}
		filter := ExtensionFilter(".xcodeproj", false)
		filtered, err := FilterPaths(paths, filter)
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
		filter := RegexpFilter(".*.xcodeproj", true)
		filtered, err := FilterPaths(paths, filter)
		require.NoError(t, err)
		require.Equal(t, []string{"path/to/my/project.xcodeproj"}, filtered)
	}

	t.Log("forbid")
	{
		paths := []string{
			"path/to/my/project.xcodeproj",
			"path/to/my/project.xcworkspace",
		}
		filter := RegexpFilter(".*.xcodeproj", false)
		filtered, err := FilterPaths(paths, filter)
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
		filter := ComponentFilter("bitrise", true)
		filtered, err := FilterPaths(paths, filter)
		require.NoError(t, err)
		require.Equal(t, []string{"/Users/bitrise/test"}, filtered)
	}

	t.Log("forbid")
	{
		paths := []string{
			"/Users/bitrise/test",
			"/Users/vagrant/test",
		}
		filter := ComponentFilter("bitrise", false)
		filtered, err := FilterPaths(paths, filter)
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
		filter := ComponentWithExtensionFilter(".framework", true)
		filtered, err := FilterPaths(paths, filter)
		require.NoError(t, err)
		require.Equal(t, []string{"/Users/bitrise.framework/test"}, filtered)
	}

	t.Log("forbid")
	{
		paths := []string{
			"/Users/bitrise.framework/test",
			"/Users/vagrant/test",
		}
		filter := ComponentWithExtensionFilter(".framework", false)
		filtered, err := FilterPaths(paths, filter)
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
		filter := IsDirectoryFilter(true)
		filtered, err := FilterPaths(paths, filter)
		require.NoError(t, err)
		require.Equal(t, []string{tmpDir}, filtered)
	}

	t.Log("forbid")
	{
		paths := []string{
			tmpDir,
			tmpFile,
		}
		filter := IsDirectoryFilter(false)
		filtered, err := FilterPaths(paths, filter)
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
		filter := InDirectoryFilter("/Users/bitrise", true)
		filtered, err := FilterPaths(paths, filter)
		require.NoError(t, err)
		require.Equal(t, []string{"/Users/bitrise/test"}, filtered)
	}

	t.Log("forbid")
	{
		paths := []string{
			"/Users/bitrise/test",
			"/Users/vagrant/test",
		}
		filter := InDirectoryFilter("/Users/bitrise", false)
		filtered, err := FilterPaths(paths, filter)
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
				got, err := DirectoryContainsFile(tt.filterFileName)(tt.path)
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
