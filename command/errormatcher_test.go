package command

import (
	"bytes"
	"strings"
	"testing"

	"github.com/bitrise-io/go-utils/v2/env"
	"github.com/stretchr/testify/require"
)

// splitLineScript writes one error line in two pieces, which is how Gradle writes nearly every line
// it cares about: its output layer writes each styled span and the newline as separate writes, so a
// real `FAILURE: Build failed with an exception.` arrives as "FAILURE: ", then the message, then the
// newline.
const splitLineScript = `printf "Execution failed for "
sleep 0.2
printf "task ':app:compileDebugKotlin'.\n"
exit 1`

func TestErrorMatcher_ReadsALineSplitAcrossWrites(t *testing.T) {
	err := runScript(t, splitLineScript, &Opts{ErrorMatcher: &prefixMatcher{prefix: "Execution failed"}})

	require.Equal(t, "Execution failed for task ':app:compileDebugKotlin'.", collectedErrorOutput(t, err))
}

func TestErrorMatcher_ReadsALineWithNoTrailingNewline(t *testing.T) {
	err := runScript(t, `printf 'error: no newline after me'; exit 1`, &Opts{ErrorMatcher: &prefixMatcher{prefix: "error:"}})

	require.Equal(t, "error: no newline after me", collectedErrorOutput(t, err))
}

// TestErrorMatcher_GivesTheLineAsWritten pins what a matcher is handed: the bytes the tool wrote,
// escape sequences included. A tool that colors its output is the matcher's problem, not this layer's.
func TestErrorMatcher_GivesTheLineAsWritten(t *testing.T) {
	var out bytes.Buffer
	matcher := &recordingMatcher{}
	const colored = `printf '\033[31;1merror: Compilation error\033[0m\n'; exit 1`

	err := runScript(t, colored, &Opts{Stdout: &out, ErrorMatcher: matcher})

	require.Error(t, err)
	require.Equal(t, []string{"\x1b[31;1merror: Compilation error\x1b[0m"}, matcher.lines)
	require.Equal(t, "\x1b[31;1merror: Compilation error\x1b[0m\n", out.String(), "the log keeps its colors")
}

func TestErrorMatcher_ReadsBothStreams(t *testing.T) {
	var stdout, stderr bytes.Buffer
	const script = `printf 'error: on stdout\n'; printf 'error: on stderr\n' >&2; exit 1`

	err := runScript(t, script, &Opts{Stdout: &stdout, Stderr: &stderr, ErrorMatcher: &prefixMatcher{prefix: "error:"}})

	require.Contains(t, collectedErrorOutput(t, err), "error: on stdout")
	require.Contains(t, collectedErrorOutput(t, err), "error: on stderr")
	require.Equal(t, "error: on stdout\n", stdout.String())
	require.Equal(t, "error: on stderr\n", stderr.String())
}

func TestErrorMatcher_SharedStdoutStderrHasOneWriter(t *testing.T) {
	var out bytes.Buffer // not goroutine-safe: only valid while os/exec keeps a single copy goroutine

	err := runScript(t, interleavedScript, &Opts{Stdout: &out, Stderr: &out, ErrorMatcher: &prefixMatcher{prefix: "error:"}})

	require.Contains(t, collectedErrorOutput(t, err), "error: on stderr")
	require.Equal(t, 4002, strings.Count(out.String(), "\n"))
	require.True(t, strings.HasSuffix(out.String(), "error: on stderr\nerror: on stdout\n"))
}

func TestErrorMatcher_NothingMatchedKeepsTheGenericMessage(t *testing.T) {
	err := runScript(t, `printf 'all quiet\n'; exit 1`, &Opts{ErrorMatcher: &prefixMatcher{prefix: "error:"}})

	require.Equal(t, "check the command's output for details", collectedErrorOutput(t, err))
}

func TestErrorMatcher_TakesPrecedenceOverErrorFinder(t *testing.T) {
	err := runScript(t, `printf 'error: the real one\n'; exit 1`, &Opts{
		ErrorMatcher: &prefixMatcher{prefix: "error:"},
		ErrorFinder:  func(string) []string { return []string{"from the chunk finder"} },
	})

	require.Equal(t, "error: the real one", collectedErrorOutput(t, err))
}

func TestErrorMatcher_SuccessfulRunForwardsEverything(t *testing.T) {
	var out bytes.Buffer

	err := runScript(t, `printf 'first\nsecond without a newline'`, &Opts{Stdout: &out, ErrorMatcher: &prefixMatcher{prefix: "error:"}})

	require.NoError(t, err)
	require.Equal(t, "first\nsecond without a newline", out.String())
}

func TestErrorMatcher_RunAndReturnTrimmedCombinedOutput(t *testing.T) {
	cmd := NewFactory(env.NewRepository()).Create("sh", []string{"-c", `printf 'error: from captured output\n'; exit 1`}, &Opts{
		ErrorMatcher: &prefixMatcher{prefix: "error:"},
	})

	out, err := cmd.RunAndReturnTrimmedCombinedOutput()

	require.Equal(t, "error: from captured output", out)
	require.Equal(t, "error: from captured output", collectedErrorOutput(t, err))
}

func TestErrorMatcher_RunAndReturnTrimmedOutput(t *testing.T) {
	const script = `printf 'all quiet\n'; printf 'error: on stderr\n' >&2; exit 1`

	t.Run("stderr not captured by the caller", func(t *testing.T) {
		cmd := NewFactory(env.NewRepository()).Create("sh", []string{"-c", script}, &Opts{
			ErrorMatcher: &prefixMatcher{prefix: "error:"},
		})

		out, err := cmd.RunAndReturnTrimmedOutput()

		require.Equal(t, "all quiet", out)
		require.Equal(t, "error: on stderr", collectedErrorOutput(t, err))
	})

	t.Run("stderr written to the caller's writer", func(t *testing.T) {
		var stderr bytes.Buffer
		cmd := NewFactory(env.NewRepository()).Create("sh", []string{"-c", script}, &Opts{
			Stderr:       &stderr,
			ErrorMatcher: &prefixMatcher{prefix: "error:"},
		})

		out, err := cmd.RunAndReturnTrimmedOutput()

		require.Equal(t, "all quiet", out)
		require.Equal(t, "error: on stderr", collectedErrorOutput(t, err))
		require.Equal(t, "error: on stderr\n", stderr.String())
	})

	t.Run("nothing matched does not fall back to raw stderr", func(t *testing.T) {
		cmd := NewFactory(env.NewRepository()).Create("sh", []string{"-c", `printf 'noise on stderr\n' >&2; exit 1`}, &Opts{
			ErrorMatcher: &prefixMatcher{prefix: "error:"},
		})

		_, err := cmd.RunAndReturnTrimmedOutput()

		require.Equal(t, "check the command's output for details", collectedErrorOutput(t, err))
	})
}

func TestErrorMatcher_RunAndReturnExitCode(t *testing.T) {
	cmd := NewFactory(env.NewRepository()).Create("sh", []string{"-c", `printf 'error: exit code run\n'; exit 3`}, &Opts{
		ErrorMatcher: &prefixMatcher{prefix: "error:"},
	})

	code, err := cmd.RunAndReturnExitCode()

	require.Equal(t, 3, code)
	require.Equal(t, "error: exit code run", collectedErrorOutput(t, err))
}

func TestErrorMatcher_StartAndWait(t *testing.T) {
	cmd := NewFactory(env.NewRepository()).Create("sh", []string{"-c", `printf 'error: started and waited'; exit 1`}, &Opts{
		ErrorMatcher: &prefixMatcher{prefix: "error:"},
	})

	require.NoError(t, cmd.Start())
	err := cmd.Wait()

	require.Equal(t, "error: started and waited", collectedErrorOutput(t, err))
}

func TestMatcherCollectorReasons(t *testing.T) {
	tests := []struct {
		name    string
		reasons []string
		want    []string
	}{
		{
			name:    "nothing found stays nothing",
			reasons: nil,
			want:    nil,
		},
		{
			name:    "what fits is kept as it is",
			reasons: []string{"Execution failed for task ':app:compileDebugKotlin'.", "e: Unresolved reference 'bar'"},
			want:    []string{"Execution failed for task ':app:compileDebugKotlin'.", "e: Unresolved reference 'bar'"},
		},
		{
			name:    "whole lines past the cap are dropped",
			reasons: []string{strings.Repeat("a", reasonBudget-10), "this line no longer fits"},
			want:    []string{strings.Repeat("a", reasonBudget-10), truncatedMarker},
		},
		{
			name:    "a first line over the cap on its own is cut",
			reasons: []string{strings.Repeat("a", maxReasonBytes+10)},
			want:    []string{strings.Repeat("a", reasonBudget), truncatedMarker},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			collector := &matcherCollector{matcher: &fixedMatcher{reasons: tt.reasons}}

			got := collector.reasons()

			require.Equal(t, tt.want, got)
			require.LessOrEqual(t, len(strings.Join(got, "\n")), maxReasonBytes)
		})
	}
}

func runScript(t *testing.T, script string, opts *Opts) error {
	t.Helper()

	return NewFactory(env.NewRepository()).Create("sh", []string{"-c", script}, opts).Run()
}

// prefixMatcher keeps every line starting with its prefix. It carries no stream handling at all,
// which is the point of the interface.
type prefixMatcher struct {
	prefix  string
	matched []string
}

func (m *prefixMatcher) ConsumeLine(line string) {
	if strings.HasPrefix(line, m.prefix) {
		m.matched = append(m.matched, line)
	}
}

func (m *prefixMatcher) Reasons() []string {
	return m.matched
}

type fixedMatcher struct {
	reasons []string
}

func (m *fixedMatcher) ConsumeLine(string) {}
func (m *fixedMatcher) Reasons() []string  { return m.reasons }
