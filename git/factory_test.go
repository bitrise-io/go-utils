package git

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/bitrise-io/go-utils/v2/command"
	"github.com/bitrise-io/go-utils/v2/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type nonTemplateAction struct{}

func (nonTemplateAction) Create(stdOut, stdErr io.Writer, envs []string) command.Command {
	return nil
}

func TestNewFactory_CreatesDir(t *testing.T) {
	tmp := t.TempDir()
	target := filepath.Join(tmp, "repo")

	factoryMock := mocks.NewFactory(t)
	f, err := NewFactory(target, factoryMock, []string{"BASE=1"})
	require.NoError(t, err)
	require.NotNil(t, f)

	st, err := os.Stat(target)
	require.NoError(t, err)
	assert.True(t, st.IsDir())
}

func TestNewFactory_ReturnsErrorWhenDirIsAFile(t *testing.T) {
	tmp := t.TempDir()
	target := filepath.Join(tmp, "not_a_dir")

	require.NoError(t, os.WriteFile(target, []byte("x"), 0o644))

	factoryMock := mocks.NewFactory(t)
	_, err := NewFactory(target, factoryMock, []string{"BASE=1"})
	require.Error(t, err)
}

func TestDefaultFactory_SetsAskPassEnv(t *testing.T) {
	tmp := t.TempDir()
	factoryMock := mocks.NewFactory(t)
	cmdMock := mocks.NewCommand(t)

	f, err := DefaultFactory(tmp, factoryMock)
	require.NoError(t, err)

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	factoryMock.
		On("Create", "git", []string{"init"}, mock.MatchedBy(func(opts *command.Opts) bool {
			if opts == nil {
				return false
			}
			if opts.Dir != tmp {
				return false
			}
			if opts.Stdout != stdout || opts.Stderr != stderr {
				return false
			}
			return reflect.DeepEqual(opts.Env, []string{"GIT_ASKPASS=echo", "RUNTIME=2"})
		})).
		Return(cmdMock).
		Once()

	got := f.Init().Create(stdout, stderr, []string{"RUNTIME=2"})
	require.Same(t, cmdMock, got)
}

func TestFactory_TemplatesBuildExpectedArgsAndOpts(t *testing.T) {
	tmp := t.TempDir()

	tests := []struct {
		name     string
		template func(Factory) Template
		wantArgs []string
	}{
		{name: "init", template: func(f Factory) Template { return f.Init() }, wantArgs: []string{"init"}},
		{name: "clone", template: func(f Factory) Template { return f.Clone("https://example.com/repo.git", "--depth", "1") }, wantArgs: []string{"clone", "--depth", "1", "https://example.com/repo.git", "."}},
		{name: "clone tag/branch", template: func(f Factory) Template {
			return f.CloneTagOrBranch("https://example.com/repo.git", "main", "--depth", "1")
		}, wantArgs: []string{"clone", "--recursive", "--branch", "main", "--depth", "1", "https://example.com/repo.git", "."}},

		{name: "remote branches", template: func(f Factory) Template { return f.RemoteBranches() }, wantArgs: []string{"ls-remote", "-b"}},
		{name: "remote list", template: func(f Factory) Template { return f.RemoteList() }, wantArgs: []string{"remote", "-v"}},
		{name: "remote add", template: func(f Factory) Template { return f.RemoteAdd("origin", "https://example.com/repo.git") }, wantArgs: []string{"remote", "add", "origin", "https://example.com/repo.git"}},
		{name: "fetch", template: func(f Factory) Template { return f.Fetch("--tags") }, wantArgs: []string{"fetch", "--tags"}},
		{name: "pull", template: func(f Factory) Template { return f.Pull() }, wantArgs: []string{"pull"}},
		{name: "push", template: func(f Factory) Template { return f.Push("main") }, wantArgs: []string{"push", "-u", "origin", "main"}},

		{name: "checkout", template: func(f Factory) Template { return f.Checkout("main") }, wantArgs: []string{"checkout", "main"}},
		{name: "merge", template: func(f Factory) Template { return f.Merge("main") }, wantArgs: []string{"merge", "main"}},
		{name: "branch", template: func(f Factory) Template { return f.Branch("-a") }, wantArgs: []string{"branch", "-a"}},
		{name: "new branch", template: func(f Factory) Template { return f.NewBranch("feature") }, wantArgs: []string{"checkout", "-b", "feature"}},

		{name: "submodule update", template: func(f Factory) Template { return f.SubmoduleUpdate("--force") }, wantArgs: []string{"submodule", "update", "--init", "--recursive", "--force"}},
		{name: "submodule foreach (template action)", template: func(f Factory) Template { return f.SubmoduleForeach(f.Checkout("main")) }, wantArgs: []string{"submodule", "foreach", "checkout", "main"}},
		{name: "submodule foreach (non-template action)", template: func(f Factory) Template { return f.SubmoduleForeach(nonTemplateAction{}) }, wantArgs: []string{"submodule", "foreach"}},

		{name: "log", template: func(f Factory) Template { return f.Log("%H", "--no-color") }, wantArgs: []string{"log", "-1", "--format=%H", "--no-color"}},
		{name: "rev-list", template: func(f Factory) Template { return f.RevList("HEAD", "--max-count", "1") }, wantArgs: []string{"rev-list", "HEAD", "--max-count", "1"}},
		{name: "rev-parse", template: func(f Factory) Template { return f.RevParse("HEAD") }, wantArgs: []string{"rev-parse", "HEAD"}},
		{name: "update-ref", template: func(f Factory) Template { return f.UpdateRef("refs/heads/main", "HEAD") }, wantArgs: []string{"update-ref", "refs/heads/main", "HEAD"}},

		{name: "config", template: func(f Factory) Template { return f.Config("user.name", "Jane", "--local") }, wantArgs: []string{"config", "user.name", "Jane", "--local"}},
		{name: "sparse-checkout init", template: func(f Factory) Template { return f.SparseCheckoutInit(false) }, wantArgs: []string{"sparse-checkout", "init"}},
		{name: "sparse-checkout init cone", template: func(f Factory) Template { return f.SparseCheckoutInit(true) }, wantArgs: []string{"sparse-checkout", "init", "--cone"}},
		{name: "sparse-checkout set", template: func(f Factory) Template { return f.SparseCheckoutSet("path/") }, wantArgs: []string{"sparse-checkout", "set", "path/"}},

		{name: "reset", template: func(f Factory) Template { return f.Reset("--hard", "HEAD") }, wantArgs: []string{"reset", "--hard", "HEAD"}},
		{name: "clean", template: func(f Factory) Template { return f.Clean("-fdx") }, wantArgs: []string{"clean", "-fdx"}},
		{name: "add", template: func(f Factory) Template { return f.Add(".") }, wantArgs: []string{"add", "."}},
		{name: "apply", template: func(f Factory) Template { return f.Apply("/tmp/patch.diff") }, wantArgs: []string{"apply", "--index", "/tmp/patch.diff"}},
		{name: "commit", template: func(f Factory) Template { return f.Commit("msg") }, wantArgs: []string{"commit", "-m", "msg"}},
		{name: "status", template: func(f Factory) Template { return f.Status("--porcelain") }, wantArgs: []string{"status", "--porcelain"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			factoryMock := mocks.NewFactory(t)
			cmdMock := mocks.NewCommand(t)

			f, err := NewFactory(tmp, factoryMock, []string{"BASE=1"})
			require.NoError(t, err)

			stdout := &bytes.Buffer{}
			stderr := &bytes.Buffer{}

			factoryMock.
				On("Create", "git", tt.wantArgs, mock.MatchedBy(func(opts *command.Opts) bool {
					if opts == nil {
						return false
					}
					if opts.Dir != tmp {
						return false
					}
					if opts.Stdout != stdout || opts.Stderr != stderr {
						return false
					}
					return reflect.DeepEqual(opts.Env, []string{"BASE=1", "RUNTIME=2"})
				})).
				Return(cmdMock).
				Once()

			tmpl := tt.template(f)
			got := tmpl.Create(stdout, stderr, []string{"RUNTIME=2"})
			require.Same(t, cmdMock, got)
			factoryMock.AssertExpectations(t)
		})
	}
}
