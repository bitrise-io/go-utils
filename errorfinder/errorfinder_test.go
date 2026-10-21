package errorfinder

import (
	"io"
	"strings"
	"sync"
	"testing"

	"github.com/bitrise-io/go-utils/v2/linepipe"
	"github.com/stretchr/testify/require"
)

func TestCollectorGivesTheMatcherCleanedLines(t *testing.T) {
	matcher := &recordingMatcher{}

	writeThrough(t, matcher, "\x1b[31;1mCompilation error\x1b[0m\r\n10%\rdone\nno newline")

	require.Equal(t, []string{"Compilation error", "done", "no newline"}, matcher.lines)
}

func TestCollectorReasons(t *testing.T) {
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
			name:    "repeated lines are reported once",
			reasons: []string{"e: Unresolved reference 'bar'", "e: Unresolved reference 'bar'", "> Task :app:compileDebugKotlin FAILED"},
			want:    []string{"e: Unresolved reference 'bar'", "> Task :app:compileDebugKotlin FAILED"},
		},
		{
			name:    "whole lines past the cap are dropped",
			reasons: []string{strings.Repeat("a", MaxReasonBytes-10), "this line no longer fits"},
			want:    []string{strings.Repeat("a", MaxReasonBytes-10), truncatedMarker},
		},
		{
			name:    "a first line over the cap on its own is cut",
			reasons: []string{strings.Repeat("a", MaxReasonBytes+10)},
			want:    []string{strings.Repeat("a", MaxReasonBytes), truncatedMarker},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, NewCollector(&fixedMatcher{reasons: tt.reasons}).Reasons())
		})
	}
}

func TestCollectorReasonsStayUnderTheCap(t *testing.T) {
	reasons := make([]string, 50)
	for i := range reasons {
		reasons[i] = strings.Repeat("a", 100)
	}

	got := NewCollector(&fixedMatcher{reasons: reasons}).Reasons()

	require.LessOrEqual(t, len(strings.Join(got, "\n")), MaxReasonBytes+len(truncatedMarker)+1)
}

func TestCollectorServesBothStreamsAtOnce(t *testing.T) {
	collector := NewCollector(&recordingMatcher{})

	var wg sync.WaitGroup
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			pipe := linepipe.New(nil, linepipe.WithTap(collector))
			for j := 0; j < 200; j++ {
				_, _ = io.WriteString(pipe, "a line\n")
			}
			_ = pipe.Close()
		}()
	}
	wg.Wait()

	require.NotEmpty(t, collector.Reasons())
}

func writeThrough(t *testing.T, matcher Matcher, output string) {
	t.Helper()

	pipe := linepipe.New(nil, linepipe.WithTap(NewCollector(matcher)))
	_, err := io.WriteString(pipe, output)
	require.NoError(t, err)
	require.NoError(t, pipe.Close())
}

type recordingMatcher struct {
	lines []string
}

func (m *recordingMatcher) ConsumeLine(line string) { m.lines = append(m.lines, line) }
func (m *recordingMatcher) Reasons() []string       { return m.lines }

type fixedMatcher struct {
	reasons []string
}

func (m *fixedMatcher) ConsumeLine(string) {}
func (m *fixedMatcher) Reasons() []string  { return m.reasons }
