package git

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/bitrise-io/go-utils/v2/command"
	"github.com/bitrise-io/go-utils/v2/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestTemplate_Create_MergesEnvsAndForwardsWriters(t *testing.T) {
	factoryMock := mocks.NewFactory(t)
	cmdMock := mocks.NewCommand(t)

	tmpl := &template{
		cmdFactory: factoryMock,
		args:       []string{"status", "--porcelain"},
		envs:       []string{"A=B"},
		dir:        "/workdir",
	}

	stdoutWriter := &bytes.Buffer{}
	stderrWriter := &bytes.Buffer{}

	factoryMock.
		On("Create", "git", []string{"status", "--porcelain"}, mock.MatchedBy(func(opts *command.Opts) bool {
			if opts == nil {
				return false
			}
			if opts.Stdout != stdoutWriter || opts.Stderr != stderrWriter {
				return false
			}
			if opts.Dir != "/workdir" {
				return false
			}
			return reflect.DeepEqual(opts.Env, []string{"A=B", "C=D"})
		})).
		Return(cmdMock).
		Once()

	got := tmpl.Create(stdoutWriter, stderrWriter, []string{"C=D"})
	require.Same(t, cmdMock, got)
}
