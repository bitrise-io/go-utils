package command

import (
	"bytes"
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"

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

func TestRunWithMatcher_ReadsALineSplitAcrossWrites(t *testing.T) {
	err := runWithMatcher(t, splitLineScript, nil, matchPrefix("Execution failed"))

	require.Equal(t, "Execution failed for task ':app:compileDebugKotlin'.", collectedErrorOutput(t, err))
}

func TestRunWithMatcher_ReadsALineWithNoTrailingNewline(t *testing.T) {
	err := runWithMatcher(t, `printf 'error: no newline after me'; exit 1`, nil, matchPrefix("error:"))

	require.Equal(t, "error: no newline after me", collectedErrorOutput(t, err))
}

// TestRunWithMatcher_GivesTheLineAsWritten pins what a matcher is handed: the bytes the tool wrote,
// escape sequences included. A tool that colors its output is the matcher's problem, not this layer's.
func TestRunWithMatcher_GivesTheLineAsWritten(t *testing.T) {
	var out bytes.Buffer
	matcher := &recordingMatcher{}
	const colored = `printf '\033[31;1merror: Compilation error\033[0m\n'; exit 1`

	err := runWithMatcher(t, colored, &Opts{Stdout: &out}, matcher)

	require.Error(t, err)
	require.Equal(t, []string{"\x1b[31;1merror: Compilation error\x1b[0m"}, matcher.lines)
	require.Equal(t, "\x1b[31;1merror: Compilation error\x1b[0m\n", out.String(), "the log keeps its colors")
}

// TestRunWithMatcher_KeepsTheOrderAcrossStreams is the reason the method owns the pipe. The caller
// passes two different writers, which with Run means two pipes and no order between them. Neither
// writer nor the matcher is goroutine-safe, so the race detector also proves there is one copy goroutine.
func TestRunWithMatcher_KeepsTheOrderAcrossStreams(t *testing.T) {
	var script strings.Builder
	var want []string
	for i := 0; i < 20; i++ {
		fmt.Fprintf(&script, "echo o%d; echo e%d >&2; ", i, i)
		want = append(want, fmt.Sprintf("o%d", i), fmt.Sprintf("e%d", i))
	}
	script.WriteString("exit 1")

	for run := 0; run < 20; run++ {
		var stdout, stderr bytes.Buffer
		matcher := &recordingMatcher{}

		err := runWithMatcher(t, script.String(), &Opts{Stdout: &stdout, Stderr: &stderr}, matcher)

		require.Error(t, err)
		require.Equal(t, want, matcher.lines, "what the matcher read, run %d", run)
		require.Equal(t, want, strings.Fields(stdout.String()), "what was forwarded, run %d", run)
		require.Empty(t, stderr.String(), "both streams are forwarded to Stdout")
	}
}

func TestRunWithMatcher_ForwardsToStderrWhenOnlyThatIsSet(t *testing.T) {
	var stderr bytes.Buffer

	err := runWithMatcher(t, `echo out; echo err >&2; exit 1`, &Opts{Stderr: &stderr}, matchPrefix("error:"))

	require.Error(t, err)
	require.Equal(t, "out\nerr\n", stderr.String())
}

// TestRunWithMatcher_ForwardsWhileTheCommandRuns pins the two routes: the caller's writer has an
// unfinished line while the command is still running, and the matcher is not given it until it ends.
func TestRunWithMatcher_ForwardsWhileTheCommandRuns(t *testing.T) {
	out := &signallingWriter{first: make(chan struct{})}
	matcher := &lockedMatcher{}
	cmd := NewFactory(env.NewRepository()).Create("sh", []string{"-c", `printf 'Downloading... 40%%'; sleep 1; printf ' done\n'`}, &Opts{Stdout: out})

	done := make(chan error, 1)
	go func() { done <- cmd.RunWithMatcher(matcher) }()

	select {
	case <-out.first:
	case <-time.After(10 * time.Second):
		t.Fatal("nothing was forwarded while the command ran")
	}
	require.Equal(t, "Downloading... 40%", out.String())
	require.Empty(t, matcher.Reasons(), "the line has not ended yet")

	require.NoError(t, <-done)
	require.Equal(t, "Downloading... 40% done\n", out.String())
	require.Equal(t, []string{"Downloading... 40% done"}, matcher.Reasons())
}

func TestRunWithMatcher_NothingMatchedKeepsTheGenericMessage(t *testing.T) {
	err := runWithMatcher(t, `printf 'all quiet\n'; exit 1`, nil, matchPrefix("error:"))

	require.Equal(t, "check the command's output for details", collectedErrorOutput(t, err))
}

func TestRunWithMatcher_DoesNotUseTheErrorFinder(t *testing.T) {
	opts := &Opts{ErrorFinder: func(string) []string { return []string{"from the chunk finder"} }}

	err := runWithMatcher(t, `printf 'error: the real one\n'; exit 1`, opts, matchPrefix("error:"))

	require.Equal(t, "error: the real one", collectedErrorOutput(t, err))
}

// TestRunWithMatcher_TheMatcherIsGivenPerRun shows why the matcher is an argument and not part of
// Opts: an Opts reused for a retry carries nothing from one run into the next.
func TestRunWithMatcher_TheMatcherIsGivenPerRun(t *testing.T) {
	opts := &Opts{}

	first := runWithMatcher(t, `printf 'error: first attempt\n'; exit 1`, opts, matchPrefix("error:"))
	second := runWithMatcher(t, `printf 'error: second attempt\n'; exit 1`, opts, matchPrefix("error:"))

	require.Equal(t, "error: first attempt", collectedErrorOutput(t, first))
	require.Equal(t, "error: second attempt", collectedErrorOutput(t, second))
}

func TestRunWithMatcher_SuccessfulRunForwardsEverything(t *testing.T) {
	var out bytes.Buffer

	err := runWithMatcher(t, `printf 'first\nsecond without a newline'`, &Opts{Stdout: &out}, matchPrefix("error:"))

	require.NoError(t, err)
	require.Equal(t, "first\nsecond without a newline", out.String())
}

func TestRunWithMatcher_ReportsACommandThatCannotStart(t *testing.T) {
	err := NewFactory(env.NewRepository()).Create("/no/such/tool", nil, nil).RunWithMatcher(matchPrefix("error:"))

	require.ErrorContains(t, err, "executing command failed")
}

func TestCapReasons(t *testing.T) {
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
			got := capReasons(tt.reasons)

			require.Equal(t, tt.want, got)
			require.LessOrEqual(t, len(strings.Join(got, "\n")), maxReasonBytes)
		})
	}
}

// runWithMatcher runs a shell script through RunWithMatcher.
func runWithMatcher(t *testing.T, script string, opts *Opts, matcher ErrorMatcher) error {
	t.Helper()

	return NewFactory(env.NewRepository()).Create("sh", []string{"-c", script}, opts).RunWithMatcher(matcher)
}

func matchPrefix(prefix string) *prefixMatcher {
	return &prefixMatcher{prefix: prefix}
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

// lockedMatcher can be read while the command is still running.
type lockedMatcher struct {
	mu    sync.Mutex
	lines []string
}

func (m *lockedMatcher) ConsumeLine(line string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.lines = append(m.lines, line)
}

func (m *lockedMatcher) Reasons() []string {
	m.mu.Lock()
	defer m.mu.Unlock()

	return append([]string(nil), m.lines...)
}

// signallingWriter tells the test when the first bytes have been forwarded.
type signallingWriter struct {
	mu       sync.Mutex
	buf      bytes.Buffer
	first    chan struct{}
	signaled bool
}

func (w *signallingWriter) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.buf.Write(p)
	if !w.signaled {
		w.signaled = true
		close(w.first)
	}

	return len(p), nil
}

func (w *signallingWriter) String() string {
	w.mu.Lock()
	defer w.mu.Unlock()

	return w.buf.String()
}
