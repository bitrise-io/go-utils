package cmdex

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_LogPrintableCommandArgs(t *testing.T) {
	origCmdArgs := []string{
		"tool",
		// simples
		"-simple", "simple",
		"another", "second-simple",
		// arg with spaces
		"-with-spaces", "arg with spaces",
		// arg with quotation
		`-with-quotation`, `extra with "quotation" included`,
	}
	printableStr := LogPrintableCommandArgs(origCmdArgs)
	expectedStr := `tool "-simple" "simple" "another" "second-simple" "-with-spaces" "arg with spaces" "-with-quotation" "extra with \"quotation\" included"`

	require.Equal(t, expectedStr, printableStr)
}
