package linepipe

import (
	"bytes"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/bitrise-io/go-utils/v2/ansi"
	"github.com/stretchr/testify/require"
)

func TestWriterAssemblesLines(t *testing.T) {
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
			name:   "nothing written, nothing observed",
			writes: nil,
			want:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var dest bytes.Buffer
			tap := &recordingTap{}

			w := New(&dest, WithTap(tap))
			for _, chunk := range tt.writes {
				n, err := w.Write([]byte(chunk))
				require.NoError(t, err)
				require.Equal(t, len(chunk), n)
			}
			require.NoError(t, w.Close())

			require.Equal(t, tt.want, tap.texts)
			require.Equal(t, strings.Join(tt.writes, ""), dest.String())
		})
	}
}

func TestWriterForwardsBeforeTheLineEnds(t *testing.T) {
	var dest bytes.Buffer
	tap := &recordingTap{}

	w := New(&dest, WithTap(tap))
	_, err := w.Write([]byte("Downloading... 40%"))
	require.NoError(t, err)

	require.Equal(t, "Downloading... 40%", dest.String())
	require.Empty(t, tap.texts)
}

func TestWriterRunsStagesInOrder(t *testing.T) {
	var order []string
	first := tapFunc(func(*Line) { order = append(order, "first") })
	second := tapFunc(func(*Line) { order = append(order, "second") })

	w := New(nil, WithTap(first), WithTap(second))
	_, err := w.Write([]byte("line\n"))
	require.NoError(t, err)

	require.Equal(t, []string{"first", "second"}, order)
}

func TestWriterParsesEachLineOnce(t *testing.T) {
	var texts []string
	var seqs [][]ansi.Sequence

	readText := tapFunc(func(l *Line) { texts = append(texts, l.Text()) })
	readSeqs := tapFunc(func(l *Line) { seqs = append(seqs, l.Sequences()) })

	w := New(nil, WithTap(readText), WithTap(readSeqs))
	_, err := w.Write([]byte("\x1b[31;1mCompilation error\x1b[0m\n"))
	require.NoError(t, err)

	require.Equal(t, []string{"Compilation error"}, texts)
	require.Len(t, seqs, 1)
	require.Len(t, seqs[0], 2)
	require.True(t, seqs[0][0].SetsRedForeground())
}

func TestWriterCutsAnOverlongLine(t *testing.T) {
	long := strings.Repeat("a", MaxLineLength+100)

	var dest bytes.Buffer
	tap := &recordingTap{}

	w := New(&dest, WithTap(tap))
	_, err := w.Write([]byte(long + "\nshort\n"))
	require.NoError(t, err)
	require.NoError(t, w.Close())

	require.Equal(t, []string{strings.Repeat("a", MaxLineLength), "short"}, tap.texts)
	require.Equal(t, long+"\nshort\n", dest.String(), "the destination still gets the line whole")
}

func TestWriterCutsOnARuneBoundary(t *testing.T) {
	const multiByte = "ő"

	long := strings.Repeat("a", MaxLineLength-1) + multiByte

	tap := &recordingTap{}
	w := New(nil, WithTap(tap))
	_, err := w.Write([]byte(long + "\n"))
	require.NoError(t, err)

	require.Equal(t, []string{strings.Repeat("a", MaxLineLength-1)}, tap.texts)
}

func TestWriterReportsADestinationError(t *testing.T) {
	failure := errors.New("disk full")
	tap := &recordingTap{}

	w := New(failingWriter{err: failure}, WithTap(tap))
	n, err := w.Write([]byte("first\n"))

	require.ErrorIs(t, err, failure)
	require.Zero(t, n)
	require.Empty(t, tap.texts, "a line that never reached the destination is not observed")
}

func TestWriterCloseIsIdempotent(t *testing.T) {
	tap := &recordingTap{}

	w := New(nil, WithTap(tap))
	_, err := w.Write([]byte("no newline"))
	require.NoError(t, err)

	require.NoError(t, w.Close())
	require.NoError(t, w.Close())

	require.Equal(t, []string{"no newline"}, tap.texts)
}

func TestWriterWithoutStagesIsATunnel(t *testing.T) {
	var dest bytes.Buffer

	w := New(&dest)
	_, err := io.WriteString(w, "first\nsecond")
	require.NoError(t, err)
	require.NoError(t, w.Close())

	require.Equal(t, "first\nsecond", dest.String())
}

type recordingTap struct {
	texts []string
}

func (r *recordingTap) Observe(line *Line) {
	r.texts = append(r.texts, line.Text())
}

type tapFunc func(*Line)

func (f tapFunc) Observe(line *Line) { f(line) }

type failingWriter struct {
	err error
}

func (w failingWriter) Write([]byte) (int, error) { return 0, w.err }
