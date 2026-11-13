package command

import (
	"bytes"
	"strings"
	"sync"
	"testing"
	"time"

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

func TestLineObserver_ReadsALineSplitAcrossWrites(t *testing.T) {
	err := runScript(script(splitLineScript, withMatcher(nil, matchPrefix("Execution failed"))))

	require.Equal(t, "Execution failed for task ':app:compileDebugKotlin'.", collectedErrorOutput(t, err))
}

func TestLineObserver_ReadsALineWithNoTrailingNewline(t *testing.T) {
	err := runScript(script(`printf 'error: no newline after me'; exit 1`, withMatcher(nil, matchPrefix("error:"))))

	require.Equal(t, "error: no newline after me", collectedErrorOutput(t, err))
}

// TestLineObserver_GivesTheLineAsWritten pins what a matcher is handed: the bytes the tool wrote,
// escape sequences included. A tool that colors its output is the matcher's problem, not this layer's.
func TestLineObserver_GivesTheLineAsWritten(t *testing.T) {
	var out bytes.Buffer
	matcher := &recordingMatcher{}
	const colored = `printf '\033[31;1merror: Compilation error\033[0m\n'; exit 1`

	err := runScript(script(colored, withMatcher(&Opts{Stdout: &out}, matcher)))

	require.Error(t, err)
	require.Equal(t, []string{"\x1b[31;1merror: Compilation error\x1b[0m"}, matcher.lines)
	require.Equal(t, "\x1b[31;1merror: Compilation error\x1b[0m\n", out.String(), "the log keeps its colors")
}

// TestLineObserver_ReadsBothStreamsInOrder is the case a Gradle matcher needs: lines of stdout and
// stderr in the order the tool wrote them, with the caller passing two different writers.
func TestLineObserver_ReadsBothStreamsInOrder(t *testing.T) {
	alternating, want := alternatingScript(20)

	for run := 0; run < 20; run++ {
		var stdout, stderr bytes.Buffer
		matcher := &recordingMatcher{}

		err := runScript(script(alternating, withMatcher(&Opts{Stdout: &stdout, Stderr: &stderr}, matcher)))

		require.Error(t, err)
		require.Equal(t, want, matcher.lines, "run %d", run)
	}
}

// TestLineObserver_ForwardsWhileTheCommandRuns pins the two routes: the caller's writer has an
// unfinished line while the command is still running, and the matcher is not given it until it ends.
func TestLineObserver_ForwardsWhileTheCommandRuns(t *testing.T) {
	out := &signallingWriter{first: make(chan struct{})}
	matcher := &lockedMatcher{}
	cmd := script(`printf 'Downloading... 40%%'; sleep 1; printf ' done\n'`, withMatcher(&Opts{Stdout: out}, matcher))

	done := make(chan error, 1)
	go func() { done <- cmd.Run() }()

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

func TestLineObserver_NothingMatchedKeepsTheGenericMessage(t *testing.T) {
	err := runScript(script(`printf 'all quiet\n'; exit 1`, withMatcher(nil, matchPrefix("error:"))))

	require.Equal(t, "check the command's output for details", collectedErrorOutput(t, err))
}

// TestLineObserver_EachCommandGetsItsOwnMatcher is the retry case: one Opts, two commands.
func TestLineObserver_EachCommandGetsItsOwnMatcher(t *testing.T) {
	opts := &Opts{NewOutputObserver: func() OutputObserver { return NewLineObserver(matchPrefix("error:")) }}

	first := runScript(script(`printf 'error: first attempt\n'; exit 1`, opts))
	second := runScript(script(`printf 'error: second attempt\n'; exit 1`, opts))

	require.Equal(t, "error: first attempt", collectedErrorOutput(t, first))
	require.Equal(t, "error: second attempt", collectedErrorOutput(t, second), "the first run's reason is not carried over")
}

// TestLineObserver_RunAndReturnTrimmedOutput has two pipes, so two matchers: each reads the lines of
// one stream, and what the stderr matcher found is reported first.
func TestLineObserver_RunAndReturnTrimmedOutput(t *testing.T) {
	opts := &Opts{NewOutputObserver: func() OutputObserver { return NewLineObserver(matchPrefix("error:")) }}

	out, err := script(bothStreamsScript, opts).RunAndReturnTrimmedOutput()

	require.Equal(t, "noise on stdout\nerror: on stdout", out)
	require.Equal(t, "error: on stderr\nerror: on stdout", collectedErrorOutput(t, err))
}

func TestLineObserver_StartAndWait(t *testing.T) {
	cmd := script(`printf 'error: started and waited'; exit 1`, withMatcher(nil, matchPrefix("error:")))

	require.NoError(t, cmd.Start())
	err := cmd.Wait()

	require.Equal(t, "error: started and waited", collectedErrorOutput(t, err))
}

// TestObserver_ReasonsAreCapped pins where the cap lives: in the command, so it holds for whatever
// an observer reports.
func TestObserver_ReasonsAreCapped(t *testing.T) {
	long := strings.Repeat("a", maxReasonBytes)
	opts := &Opts{NewOutputObserver: func() OutputObserver {
		return NewLineObserver(&fixedMatcher{reasons: []string{"first", long, "last"}})
	}}

	err := runScript(script(`exit 1`, opts))

	require.Equal(t, "first\n"+truncatedMarker, collectedErrorOutput(t, err))
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

// withMatcher returns opts with a line observer that feeds matcher. The matcher is shared by every
// observer the command builds, so use it only where the command has one pipe.
func withMatcher(opts *Opts, matcher ErrorMatcher) *Opts {
	if opts == nil {
		opts = &Opts{}
	}
	opts.NewOutputObserver = func() OutputObserver { return NewLineObserver(matcher) }

	return opts
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

type fixedMatcher struct {
	reasons []string
}

func (m *fixedMatcher) ConsumeLine(string) {}
func (m *fixedMatcher) Reasons() []string  { return m.reasons }

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
