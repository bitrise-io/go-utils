package log

import (
	"bytes"
	"regexp"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSetOutWriter(t *testing.T) {
	var b bytes.Buffer
	deprecatedLogger.stdout = &b
	Printf("test %s", "log")
	require.Equal(t, "test log\n", b.String())
}

func TestSetEnableDebugLog(t *testing.T) {
	t.Log("enable debug log")
	{
		SetEnableDebugLog(true)
		var b bytes.Buffer
		deprecatedLogger.stdout = &b
		Debugf("test %s", "log")
		require.Equal(t, "\x1b[35;1mtest log\x1b[0m\n", b.String())
	}

	t.Log("disable debug log")
	{
		SetEnableDebugLog(false)
		var b bytes.Buffer
		deprecatedLogger.stdout = &b
		Debugf("test %s", "log")
		require.Equal(t, "", b.String())
	}
}

func TestSetTimestampLayout(t *testing.T) {
	var b bytes.Buffer
	deprecatedLogger.stdout = &b
	SetTimestampLayout("15-04-05")
	TPrintf("test %s", "log")
	re := regexp.MustCompile(`\[.+-.+-.+\] test log`)
	require.True(t, re.MatchString(b.String()), b.String())
}

func Test_printf_with_time(t *testing.T) {
	var b bytes.Buffer
	logger := defaultLogger{
		enableDebugLog:  false,
		timestampLayout: "15.04.05",
		stdout:          &b,
	}
	logger.TPrintf("test %s", "log")
	re := regexp.MustCompile(`\[.+\..+\..+\] test log`)
	require.True(t, re.MatchString(b.String()), b.String())
}

func Test_printf_severity(t *testing.T) {
	t.Log("error")
	{
		var b bytes.Buffer
		logger := defaultLogger{
			enableDebugLog:  false,
			timestampLayout: "",
			stdout:          &b,
		}
		logger.Errorf("test %s", "log")
		require.Equal(t, "\x1b[31;1mtest log\x1b[0m\n", b.String())
	}

	t.Log("warn")
	{
		var b bytes.Buffer
		logger := defaultLogger{
			enableDebugLog:  false,
			timestampLayout: "",
			stdout:          &b,
		}
		logger.Warnf("test %s", "log")
		require.Equal(t, "\x1b[33;1mtest log\x1b[0m\n", b.String())
	}

	t.Log("debug")
	{
		var b bytes.Buffer
		logger := defaultLogger{
			enableDebugLog:  true,
			timestampLayout: "",
			stdout:          &b,
		}
		logger.Debugf("test %s", "log")
		require.Equal(t, "\x1b[35;1mtest log\x1b[0m\n", b.String())
	}

	t.Log("normal")
	{
		var b bytes.Buffer
		logger := defaultLogger{
			enableDebugLog:  false,
			timestampLayout: "",
			stdout:          &b,
		}
		logger.Printf("test %s", "log")
		require.Equal(t, "test log\n", b.String())
	}

	t.Log("info")
	{
		var b bytes.Buffer
		logger := defaultLogger{
			enableDebugLog:  false,
			timestampLayout: "",
			stdout:          &b,
		}
		logger.Infof("test %s", "log")
		require.Equal(t, "\x1b[34;1mtest log\x1b[0m\n", b.String())
	}

	t.Log("success")
	{
		var b bytes.Buffer
		logger := defaultLogger{
			enableDebugLog:  false,
			timestampLayout: "",
			stdout:          &b,
		}
		logger.Donef("test %s", "log")
		require.Equal(t, "\x1b[32;1mtest log\x1b[0m\n", b.String())
	}
}
