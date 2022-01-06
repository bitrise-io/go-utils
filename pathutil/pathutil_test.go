package pathutil

import (
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_pathProvider_CreateTempDir(t *testing.T) {
	tests := []struct {
		name   string
		prefix string
	}{
		{
			name:   "prefix provided",
			prefix: "some-test",
		},
		{
			name:   "empty prefix",
			prefix: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := pathProvider{}
			gotDir, err := p.CreateTempDir(tt.prefix)

			require.NoError(t, err)
			require.True(t, len(gotDir) != 0)
			require.True(t, strings.HasPrefix(filepath.Base(gotDir), tt.prefix))
			// returned temp dir path should not have a / at it's end
			require.False(t, strings.HasSuffix(gotDir, "/"))
			// directory is created
			info, err := os.Lstat(gotDir)
			require.NoError(t, err)
			require.True(t, info.IsDir())
		})
	}
}

func Test_pathChecker_IsPathExists(t *testing.T) {
	tests := []struct {
		name    string
		pth     string
		want    bool
		wantErr bool
	}{
		{
			name: "path does not exists",
			pth:  filepath.Join("this", "should", "not", "exist"),
			want: false,
		},
		{
			name: "current directory",
			pth:  ".",
			want: true,
		},
		{
			name:    "empty path",
			pth:     "",
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := pathChecker{}
			got, err := c.IsPathExists(tt.pth)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, tt.want, got)
		})
	}
}

func Test_pathModifier_AbsPath(t *testing.T) {
	currDirPath, err := filepath.Abs(".")
	require.NoError(t, err)
	require.NotEqual(t, "", currDirPath)
	require.NotEqual(t, ".", currDirPath)

	currentUser, err := user.Current()
	require.NoError(t, err)

	sep := string(os.PathSeparator)
	homePathEnv := filepath.Join(sep, "path", "home", "test-user")
	require.Equal(t, nil, os.Setenv("HOME", homePathEnv))

	testFileRelPathFromHome := filepath.Join("some", "file.ext")

	tests := []struct {
		name    string
		pth     string
		want    string
		wantErr bool
	}{
		{
			name:    "Empty path",
			pth:     "",
			want:    "",
			wantErr: true,
		},
		{
			name: "Current dir",
			pth:  ".",
			want: currDirPath,
		},
		{
			name: "Relative dir",
			pth:  filepath.Join(homePathEnv, "..", "test-user"),
			want: homePathEnv,
		},
		{
			name: "Environment variable",
			pth:  filepath.Join("$HOME", testFileRelPathFromHome),
			want: filepath.Join(homePathEnv, testFileRelPathFromHome),
		},
		{
			name: "Tilde with path",
			pth:  filepath.Join("~", testFileRelPathFromHome),
			want: filepath.Join(homePathEnv, testFileRelPathFromHome),
		},
		{
			name: "Tilde with relative path",
			pth:  "~" + sep + ".." + sep + "test-user",
			want: homePathEnv,
		},
		{
			name: "Tilde with slash suffix",
			pth:  "~" + sep,
			want: homePathEnv,
		},
		{
			name: "Tilde only",
			pth:  "~",
			want: homePathEnv,
		},
		{
			name: "Tilde + username",
			pth:  "~" + currentUser.Name,
			want: currentUser.HomeDir,
		},
		{
			name:    "Tilde with nonexistent username",
			pth:     filepath.Join("~testaccnotexist", "folder"),
			wantErr: true,
		},
		{
			name: "Tilde + username, slash suffix",
			pth:  "~" + currentUser.Name + sep,
			want: currentUser.HomeDir,
		},
		{
			name: "Tilde + username with path",
			pth:  filepath.Join("~"+currentUser.Name, "folder"),
			want: filepath.Join(currentUser.HomeDir, "folder"),
		},
		{
			name: "Tilde as directory name",
			pth:  filepath.Join(sep, "test", "~", "in", "name"),
			want: filepath.Join(sep, "test", "~", "in", "name"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := pathModifier{}
			got, err := p.AbsPath(tt.pth)

			if (err != nil) != tt.wantErr {
				t.Errorf("pathModifier.AbsPath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			require.Equal(t, tt.want, got)
		})
	}
}
