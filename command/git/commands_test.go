package git

import (
	"testing"

	"github.com/bitrise-io/go-utils/command"
	"github.com/stretchr/testify/assert"
)

func TestGitCommands(t *testing.T) {
	type testCase struct {
		command *command.Model
		want    string
	}

	testCases := []testCase{
		// SparseCheckout
		{
			command: (&Git{}).SparseCheckoutInit(false),
			want:    `git "sparse-checkout" "init"`,
		},
		{
			command: (&Git{}).SparseCheckoutInit(true),
			want:    `git "sparse-checkout" "init" "--cone"`,
		},
		{
			command: (&Git{}).SparseCheckoutSet("client/android"),
			want:    `git "sparse-checkout" "set" "client/android"`,
		},
		{
			command: (&Git{}).SparseCheckoutSet("client/android", "client/ios"),
			want:    `git "sparse-checkout" "set" "client/android" "client/ios"`,
		},
	}

	for _, testCase := range testCases {
		assertPrintableCommandArgs(t, testCase.want, testCase.command)
	}
}

func assertPrintableCommandArgs(t *testing.T, expectedArgs string, gitCommand *command.Model) {
	actualArgs := gitCommand.PrintableCommandArgs()
	assert.Equal(t, expectedArgs, actualArgs)
}
