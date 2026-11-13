package command

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLineSplitter(t *testing.T) {
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
			name:   "an escape sequence split across writes is handed over whole",
			writes: []string{"\x1b[3", "1mred\x1b[0m\n"},
			want:   []string{"\x1b[31mred\x1b[0m"},
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
			w, lines := tapLines(matcher, &dest)

			for _, chunk := range tt.writes {
				n, err := w.Write([]byte(chunk))
				require.NoError(t, err)
				require.Equal(t, len(chunk), n)
			}
			lines.Close()

			require.Equal(t, tt.want, matcher.lines)
			require.Equal(t, strings.Join(tt.writes, ""), dest.String())
		})
	}
}

func TestLineSplitterForwardsBeforeTheLineEnds(t *testing.T) {
	var dest bytes.Buffer
	matcher := &recordingMatcher{}
	w, _ := tapLines(matcher, &dest)

	_, err := w.Write([]byte("Downloading... 40%"))
	require.NoError(t, err)

	require.Equal(t, "Downloading... 40%", dest.String())
	require.Empty(t, matcher.lines)
}

func TestLineSplitterDropsAnOverlongLine(t *testing.T) {
	long := strings.Repeat("a", maxLineLength+100) + " error: not a reason"

	var dest bytes.Buffer
	matcher := &recordingMatcher{}
	w, _ := tapLines(matcher, &dest)

	_, err := w.Write([]byte(long + "\nshort\n"))
	require.NoError(t, err)

	require.Equal(t, []string{"", "short"}, matcher.lines,
		"the line is offered empty, so nothing in it can match and a multi-line block still ends here")
	require.Equal(t, long+"\nshort\n", dest.String(), "the caller's writer still gets the line whole")
}

func TestLineSplitterDropsAnOverlongLastLine(t *testing.T) {
	matcher := &recordingMatcher{}
	w, lines := tapLines(matcher, &bytes.Buffer{})

	_, err := w.Write([]byte(strings.Repeat("a", maxLineLength+100)))
	require.NoError(t, err)
	lines.Close()

	require.Equal(t, []string{""}, matcher.lines)
}

func TestLineSplitterRecoversAfterAnOverlongLine(t *testing.T) {
	matcher := &recordingMatcher{}
	w, _ := tapLines(matcher, &bytes.Buffer{})

	_, err := w.Write([]byte(strings.Repeat("a", maxLineLength+100) + "\n> Task :app:x FAILED\n"))
	require.NoError(t, err)

	require.Equal(t, []string{"", "> Task :app:x FAILED"}, matcher.lines)
}

func TestLineSplitterCloseIsIdempotent(t *testing.T) {
	matcher := &recordingMatcher{}
	w, lines := tapLines(matcher, &bytes.Buffer{})

	_, err := w.Write([]byte("no newline"))
	require.NoError(t, err)

	lines.Close()
	lines.Close()

	require.Equal(t, []string{"no newline"}, matcher.lines)
}

// tapLines builds what the command puts between the tool and the caller's writer when a line
// observer is set.
func tapLines(matcher ErrorMatcher, dest io.Writer) (*tap, OutputObserver) {
	lines := NewLineObserver(matcher)

	return &tap{dest: dest, observer: lines}, lines
}

type recordingMatcher struct {
	lines []string
}

func (m *recordingMatcher) ConsumeLine(line string) { m.lines = append(m.lines, line) }
func (m *recordingMatcher) Reasons() []string       { return m.lines }
