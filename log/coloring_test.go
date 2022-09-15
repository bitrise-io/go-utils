package log

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_format_with_severity_color(t *testing.T) {
	testcases := []struct {
		title string

		expectedOutput string
		format         string
		params         string
		severity       Severity
	}{
		{
			title: "Error",

			expectedOutput: "\x1b[31;1mtest log\x1b[0m",
			format:         "test %s",
			params:         "log",
			severity:       errorSeverity,
		},
		{
			title: "Warn",

			expectedOutput: "\x1b[33;1mtest log\x1b[0m",
			format:         "test %s",
			params:         "log",
			severity:       warnSeverity,
		},
		{
			title: "Debug",

			expectedOutput: "\x1b[35;1mtest log\x1b[0m",
			format:         "test %s",
			params:         "log",
			severity:       debugSeverity,
		},
		{
			title: "Info",

			expectedOutput: "\x1b[34;1mtest log\x1b[0m",
			format:         "test %s",
			params:         "log",
			severity:       infoSeverity,
		},
		{
			title: "Done",

			expectedOutput: "\x1b[32;1mtest log\x1b[0m",
			format:         "test %s",
			params:         "log",
			severity:       doneSeverity,
		},
	}

	for _, testCase := range testcases {
		t.Run(testCase.title, func(t *testing.T) {
			require.Equal(t, testCase.expectedOutput, FormatWithSeverityColor(testCase.severity, testCase.format, testCase.params))
		})
	}
}
