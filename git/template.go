package git

import (
	"io"

	"github.com/bitrise-io/go-utils/v2/command"
)

// Template ...
type Template interface {
	Create(stdOut, stdErr io.Writer, envs []string) command.Command
}

type template struct {
	cmdFactory command.Factory
	args       []string
	envs       []string
	dir        string
}

// Create ...
func (t *template) Create(stdOut, stdErr io.Writer, envs []string) command.Command {
	opts := &command.Opts{
		Stdout: stdOut,
		Stderr: stdErr,
		Dir:    t.dir,
		Env:    append(t.envs, envs...),
	}

	return t.cmdFactory.Create("git", t.args, opts)
}
