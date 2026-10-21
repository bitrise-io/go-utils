package command

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLineWriter(t *testing.T) {
	tests := []struct {
		name   string
		writes []string
		want   []string
	}{
		{
			name:   "one write per line",
			writes: []string{"first\n", "second\n"},
			want:   []string{"first", "second"},
		},
		{
			name:   "many lines in one write",
			writes: []string{"first\nsecond\nthird\n"},
			want:   []string{"first", "second", "third"},
		},
		{
			name:   "a line split across writes",
			writes: []string{"Execution failed for ", "task ':app:compileDebugKotlin'.\n"},
			want:   []string{"Execution failed for task ':app:compileDebugKotlin'."},
		},
		{
			name:   "an escape sequence split across writes",
			writes: []string{"\x1b[3", "1mred\x1b[0m\n"},
			want:   []string{"red"},
		},
		{
			name:   "a last line with no newline is flushed on close",
			writes: []string{"first\n", "no newline here"},
			want:   []string{"first", "no newline here"},
		},
		{
			name:   "empty lines are lines",
			writes: []string{"\n\nfirst\n"},
			want:   []string{"", "", "first"},
		},
		{
			name:   "CRLF",
			writes: []string{"first\r\nsecond\r\n"},
			want:   []string{"first", "second"},
		},
		{
			name:   "nothing written, nothing consumed",
			writes: nil,
			want:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var dest bytes.Buffer
			matcher := &recordingMatcher{}
			w := (&matcherCollector{matcher: matcher}).collectFromWriter(&dest)

			for _, chunk := range tt.writes {
				n, err := w.Write([]byte(chunk))
				require.NoError(t, err)
				require.Equal(t, len(chunk), n)
			}
			require.NoError(t, w.Close())

			require.Equal(t, tt.want, matcher.lines)
			require.Equal(t, strings.Join(tt.writes, ""), dest.String())
		})
	}
}

func TestLineWriterForwardsBeforeTheLineEnds(t *testing.T) {
	var dest bytes.Buffer
	matcher := &recordingMatcher{}
	w := (&matcherCollector{matcher: matcher}).collectFromWriter(&dest)

	_, err := w.Write([]byte("Downloading... 40%"))
	require.NoError(t, err)

	require.Equal(t, "Downloading... 40%", dest.String())
	require.Empty(t, matcher.lines)
}

func TestLineWriterCutsAnOverlongLine(t *testing.T) {
	long := strings.Repeat("a", maxLineLength+100)

	var dest bytes.Buffer
	matcher := &recordingMatcher{}
	w := (&matcherCollector{matcher: matcher}).collectFromWriter(&dest)

	_, err := w.Write([]byte(long + "\nshort\n"))
	require.NoError(t, err)

	require.Equal(t, []string{strings.Repeat("a", maxLineLength), "short"}, matcher.lines)
	require.Equal(t, long+"\nshort\n", dest.String(), "the caller's writer still gets the line whole")
}

func TestLineWriterCutsOnARuneBoundary(t *testing.T) {
	long := strings.Repeat("a", maxLineLength-1) + "ő"

	matcher := &recordingMatcher{}
	w := (&matcherCollector{matcher: matcher}).collectFromWriter(&bytes.Buffer{})

	_, err := w.Write([]byte(long + "\n"))
	require.NoError(t, err)

	require.Equal(t, []string{strings.Repeat("a", maxLineLength-1)}, matcher.lines)
}

func TestLineWriterReportsADestinationError(t *testing.T) {
	failure := errors.New("disk full")
	matcher := &recordingMatcher{}
	w := (&matcherCollector{matcher: matcher}).collectFromWriter(failingWriter{err: failure})

	n, err := w.Write([]byte("first\n"))

	require.ErrorIs(t, err, failure)
	require.Zero(t, n)
	require.Empty(t, matcher.lines, "a line that never reached the destination is not consumed")
}

func TestLineWriterCloseIsIdempotent(t *testing.T) {
	matcher := &recordingMatcher{}
	w := (&matcherCollector{matcher: matcher}).collectFromWriter(&bytes.Buffer{})

	_, err := w.Write([]byte("no newline"))
	require.NoError(t, err)

	require.NoError(t, w.Close())
	require.NoError(t, w.Close())

	require.Equal(t, []string{"no newline"}, matcher.lines)
}

type recordingMatcher struct {
	lines []string
}

func (m *recordingMatcher) ConsumeLine(line string) { m.lines = append(m.lines, line) }
func (m *recordingMatcher) Reasons() []string       { return m.lines }

type failingWriter struct {
	err error
}

func (w failingWriter) Write([]byte) (int, error) { return 0, w.err }
