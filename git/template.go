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
	mergedEnv := make([]string, 0, len(t.envs)+len(envs))
	mergedEnv = append(mergedEnv, t.envs...)
	mergedEnv = append(mergedEnv, envs...)

	opts := &command.Opts{
		Stdout: stdOut,
		Stderr: stdErr,
		Dir:    t.dir,
		Env:    mergedEnv,
	}

	return t.cmdFactory.Create("git", t.args, opts)
}
