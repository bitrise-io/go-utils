package git

import (
	"os"

	"github.com/bitrise-io/go-utils/v2/command"
)

// Factory ...
type Factory interface {
	Init() Template
	Clone(repo string, opts ...string) Template
	CloneTagOrBranch(repo, tagOrBranch string, opts ...string) Template

	RemoteBranches() Template
	RemoteList() Template
	RemoteAdd(name, url string) Template
	Fetch(opts ...string) Template
	Pull() Template
	Push(branch string) Template

	Checkout(args ...string) Template
	Merge(arg string) Template
	Branch(opts ...string) Template
	NewBranch(branch string) Template

	SubmoduleUpdate(opts ...string) Template
	SubmoduleForeach(action Template) Template

	Log(format string, opts ...string) Template
	RevList(commit string, opts ...string) Template
	RevParse(arg string) Template
	UpdateRef(opts ...string) Template

	Config(key string, value string, opts ...string) Template
	SparseCheckoutInit(cone bool) Template
	SparseCheckoutSet(opts ...string) Template

	Reset(mode, commit string) Template
	Clean(options ...string) Template
	Add(pathspec string) Template
	Apply(patch string) Template
	Commit(message string) Template
	Status(opts ...string) Template
}

// NewFactory ...
func NewFactory(dir string, cmdFactory command.Factory, envs []string) (Factory, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, err
	}

	return &factory{
		cmdFactory: cmdFactory,
		dir:        dir,
		envs:       envs,
	}, nil
}

// DefaultFactory ...
func DefaultFactory(dir string, cmdFactory command.Factory) (Factory, error) {
	return NewFactory(
		dir,
		cmdFactory,
		[]string{"GIT_ASKPASS=echo"},
	)
}

type factory struct {
	cmdFactory command.Factory
	envs       []string
	dir        string
}

func (f *factory) template(args ...string) Template {
	return &template{
		cmdFactory: f.cmdFactory,
		args:       args,
		envs:       f.envs,
		dir:        f.dir,
	}
}

// Status ...
func (f *factory) Status(opts ...string) Template {
	args := append([]string{"status"}, opts...)
	return f.template(args...)
}

// Commit ...
func (f *factory) Commit(message string) Template {
	return f.template("commit", "-m", message)
}

// Apply ...
func (f *factory) Apply(patch string) Template {
	return f.template("apply", "--index", patch)
}

// Add ...
func (f *factory) Add(pathspec string) Template {
	return f.template("add", pathspec)
}

// Reset ...
func (f *factory) Reset(mode, commit string) Template {
	return f.template("reset", mode, commit)
}

// SparseCheckoutSet ...
func (f *factory) SparseCheckoutSet(opts ...string) Template {
	args := append([]string{"sparse-checkout", "set"}, opts...)
	return f.template(args...)
}

// SparseCheckoutInit ...
func (f *factory) SparseCheckoutInit(cone bool) Template {
	args := []string{"sparse-checkout", "init"}
	if cone {
		args = append(args, "--cone")
	}
	return f.template(args...)
}

// Config ...
func (f *factory) Config(key string, value string, opts ...string) Template {
	args := []string{"config", key, value}
	args = append(args, opts...)
	return f.template(args...)
}

// UpdateRef ...
func (f *factory) UpdateRef(opts ...string) Template {
	args := append([]string{"update-ref"}, opts...)
	return f.template(args...)
}

// RevParse ...
func (f *factory) RevParse(arg string) Template {
	return f.template("rev-parse", arg)
}

// RevList ...
func (f *factory) RevList(commit string, opts ...string) Template {
	args := []string{"rev-list", commit}
	args = append(args, opts...)
	return f.template(args...)
}

// Log ...
func (f *factory) Log(format string, opts ...string) Template {
	args := []string{"log", "-1", "--format=" + format}
	args = append(args, opts...)
	return f.template(args...)
}

// SubmoduleForeach ...
func (f *factory) SubmoduleForeach(action Template) Template {
	args := []string{"submodule", "foreach"}
	if t, ok := action.(*template); ok {
		args = append(args, t.args...)
	}
	return f.template(args...)
}

// SubmoduleUpdate ...
func (f *factory) SubmoduleUpdate(opts ...string) Template {
	args := []string{"submodule", "update", "--init", "--recursive"}
	args = append(args, opts...)
	return f.template(args...)
}

// NewBranch ...
func (f *factory) NewBranch(branch string) Template {
	return f.template("checkout", "-b", branch)
}

// Branch returns a template for `git branch`.
func (f *factory) Branch(opts ...string) Template {
	args := append([]string{"branch"}, opts...)
	return f.template(args...)
}

// Merge ...
func (f *factory) Merge(arg string) Template {
	return f.template("merge", arg)
}

// Checkout ...
func (f *factory) Checkout(args ...string) Template {
	a := append([]string{"checkout"}, args...)
	return f.template(a...)
}

// Fetch ...
func (f *factory) Fetch(opts ...string) Template {
	args := append([]string{"fetch"}, opts...)
	return f.template(args...)
}

// RemoteAdd ...
func (f *factory) RemoteAdd(name, url string) Template {
	return f.template("remote", "add", name, url)
}

// RemoteList ...
func (f *factory) RemoteList() Template {
	return f.template("remote", "-v")
}

// RemoteBranches ...
func (f *factory) RemoteBranches() Template {
	return f.template("ls-remote", "-b")
}

// CloneTagOrBranch ...
func (f *factory) CloneTagOrBranch(repo, tagOrBranch string, opts ...string) Template {
	args := []string{"clone", "--recursive", "--branch", tagOrBranch}
	args = append(args, opts...)
	args = append(args, repo, ".")
	return f.template(args...)
}

// Clone ...
func (f *factory) Clone(repo string, opts ...string) Template {
	args := []string{"clone"}
	args = append(args, opts...)
	args = append(args, repo, ".")
	return f.template(args...)
}

// Init ...
func (f *factory) Init() Template {
	return f.template("init")
}

// Push ...
func (f *factory) Push(branch string) Template {
	return f.template("push", "-u", "origin", branch)
}

// Pull ...
func (f *factory) Pull() Template {
	return f.template("pull")
}

// Clean ...
func (f *factory) Clean(options ...string) Template {
	args := append([]string{"clean"}, options...)
	return f.template(args...)
}
