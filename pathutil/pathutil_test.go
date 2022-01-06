package pathutil

import (
	"os"
	"os/user"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

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
