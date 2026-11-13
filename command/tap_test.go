package command

import (
	"bytes"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTapForwardsBeforeItObserves(t *testing.T) {
	observer := &recordingObserver{}
	dest := &orderCheckingWriter{t: t, observer: observer}
	w := &tap{dest: dest, observer: observer}

	for _, chunk := range []string{"first ", "second\n", "third"} {
		n, err := w.Write([]byte(chunk))
		require.NoError(t, err)
		require.Equal(t, len(chunk), n)
	}

	require.Equal(t, "first second\nthird", dest.buf.String())
	require.Equal(t, []string{"first ", "second\n", "third"}, observer.seen)
}

func TestTapObservesOnlyWhatTheDestinationTook(t *testing.T) {
	observer := &recordingObserver{}
	w := &tap{dest: shortWriter{take: 3}, observer: observer}

	n, err := w.Write([]byte("first\n"))

	require.ErrorIs(t, err, io.ErrShortWrite)
	require.Equal(t, 3, n)
	require.Equal(t, []string{"fir"}, observer.seen)
}

func TestTapReportsADestinationError(t *testing.T) {
	failure := errors.New("disk full")
	observer := &recordingObserver{}
	w := &tap{dest: failingWriter{err: failure}, observer: observer}

	n, err := w.Write([]byte("first\n"))

	require.ErrorIs(t, err, failure)
	require.Zero(t, n)
	require.Empty(t, observer.seen, "what never reached the destination is not observed")
}

func TestTapPassesBackAnInvalidWriteCount(t *testing.T) {
	for _, count := range []int{-1, len("first\n") + 1} {
		observer := &recordingObserver{}
		w := &tap{dest: miscountingWriter{count: count}, observer: observer}

		_, err := io.Copy(w, io.LimitReader(strings.NewReader("first\n"), 1024))

		require.ErrorContains(t, err, "invalid write result")
		require.Empty(t, observer.seen)
	}
}

type recordingObserver struct {
	seen []string
}

func (o *recordingObserver) Observe(p []byte)  { o.seen = append(o.seen, string(p)) }
func (o *recordingObserver) Close()            {}
func (o *recordingObserver) Reasons() []string { return nil }

// orderCheckingWriter fails the test if the observer was shown a chunk before the writer got it.
type orderCheckingWriter struct {
	t        *testing.T
	observer *recordingObserver
	buf      bytes.Buffer
	writes   int
}

func (w *orderCheckingWriter) Write(p []byte) (int, error) {
	require.Len(w.t, w.observer.seen, w.writes, "the observer is shown a chunk only after the writer took it")
	w.writes++

	return w.buf.Write(p)
}

type shortWriter struct {
	take int
}

func (w shortWriter) Write([]byte) (int, error) { return w.take, io.ErrShortWrite }

type failingWriter struct {
	err error
}

func (w failingWriter) Write([]byte) (int, error) { return 0, w.err }

type miscountingWriter struct {
	count int
}

func (w miscountingWriter) Write([]byte) (int, error) { return w.count, nil }
