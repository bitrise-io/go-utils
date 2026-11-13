//go:build unix

package command

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"syscall"
	"testing"

	"github.com/bitrise-io/go-utils/v2/env"
	"github.com/stretchr/testify/require"
)

// alternatingScript writes to stdout and stderr in turn with nothing in between, which is the case
// two pipes cannot keep in order.
func alternatingScript(lines int) (script string, want []string) {
	var b strings.Builder
	for i := 0; i < lines; i++ {
		fmt.Fprintf(&b, "echo o%d; echo e%d >&2; ", i, i)
		want = append(want, fmt.Sprintf("o%d", i), fmt.Sprintf("e%d", i))
	}
	b.WriteString("exit 1")

	return b.String(), want
}

// samePipeTwice returns two *os.File for the write end of one pipe, the way a step holds os.Stdout and
// os.Stderr when the CLI gave it a single pipe for both.
func samePipeTwice(t *testing.T) (r, stdout, stderr *os.File) {
	t.Helper()

	r, w, err := os.Pipe()
	require.NoError(t, err)

	fd, err := syscall.Dup(int(w.Fd()))
	require.NoError(t, err)

	return r, w, os.NewFile(uintptr(fd), "stderr")
}

func TestErrorMatcher_TwoFilesOfOnePipeKeepTheirOrder(t *testing.T) {
	script, want := alternatingScript(20)

	for run := 0; run < 20; run++ {
		r, stdout, stderr := samePipeTwice(t)
		matcher := &prefixMatcher{}

		err := NewFactory(env.NewRepository()).Create("sh", []string{"-c", script}, &Opts{
			Stdout:          stdout,
			Stderr:          stderr,
			NewErrorMatcher: func() ErrorMatcher { return matcher },
		}).Run()
		require.Error(t, err)

		require.NoError(t, stdout.Close())
		require.NoError(t, stderr.Close())
		forwarded, readErr := io.ReadAll(r)
		require.NoError(t, readErr)
		require.NoError(t, r.Close())

		require.Equal(t, want, matcher.matched, "what the matcher read, run %d", run)
		require.Equal(t, want, strings.Fields(string(forwarded)), "what the caller's pipe received, run %d", run)
	}
}

func TestSameDestination(t *testing.T) {
	r, stdout, stderr := samePipeTwice(t)
	otherR, other, err := os.Pipe()
	require.NoError(t, err)
	t.Cleanup(func() {
		for _, f := range []*os.File{r, stdout, stderr, otherR, other} {
			_ = f.Close()
		}
	})

	path := t.TempDir() + "/log"
	fileOnce, err := os.Create(path)
	require.NoError(t, err)
	fileTwice, err := os.OpenFile(path, os.O_WRONLY, 0)
	require.NoError(t, err)
	t.Cleanup(func() { _ = fileOnce.Close(); _ = fileTwice.Close() })

	buffer := &bytes.Buffer{}
	var closed *os.File

	tests := []struct {
		name string
		a, b io.Writer
		want bool
	}{
		{name: "the same writer", a: buffer, b: buffer, want: true},
		{name: "one pipe held as two files", a: stdout, b: stderr, want: true},
		{name: "one file opened twice", a: fileOnce, b: fileTwice, want: true},
		{name: "two pipes", a: stdout, b: other, want: false},
		{name: "a pipe and a file", a: stdout, b: fileOnce, want: false},
		{name: "a file and a buffer", a: stdout, b: buffer, want: false},
		{name: "two buffers", a: buffer, b: &bytes.Buffer{}, want: false},
		{name: "a nil file", a: stdout, b: closed, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, sameDestination(tt.a, tt.b))
		})
	}
}
