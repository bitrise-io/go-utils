package command

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/bitrise-io/go-utils/v2/env"
	"github.com/stretchr/testify/require"
)

const interleavedScript = `i=0
while [ $i -lt 2000 ]; do
  echo "out $i"
  echo "err $i" >&2
  i=$((i+1))
done
echo "error: on stderr" >&2
echo "error: on stdout"
exit 1`

func findErrorLines(out string) []string {
	var lines []string
	for _, line := range strings.Split(out, "\n") {
		if strings.HasPrefix(line, "error:") {
			lines = append(lines, line)
		}
	}
	return lines
}

func TestErrorFinder_SharedStdoutStderrHasOneWriter(t *testing.T) {
	var out bytes.Buffer // not goroutine-safe: only valid while os/exec keeps a single copy goroutine
	cmd := NewFactory(env.NewRepository()).Create("sh", []string{"-c", interleavedScript}, &Opts{
		Stdout:      &out,
		Stderr:      &out,
		ErrorFinder: findErrorLines,
	})

	err := cmd.Run()

	require.ErrorContains(t, err, "error: on stderr")
	require.ErrorContains(t, err, "error: on stdout")

	var want strings.Builder
	for i := 0; i < 2000; i++ {
		fmt.Fprintf(&want, "out %d\nerr %d\n", i, i)
	}
	want.WriteString("error: on stderr\nerror: on stdout\n")
	require.Equal(t, want.String(), out.String(), "shared writer must see every byte, in the child's write order")
}

func TestErrorFinder_SeparateStdoutStderrCollectsBoth(t *testing.T) {
	var stdout, stderr bytes.Buffer
	cmd := NewFactory(env.NewRepository()).Create("sh", []string{"-c", interleavedScript}, &Opts{
		Stdout:      &stdout,
		Stderr:      &stderr,
		ErrorFinder: findErrorLines,
	})

	err := cmd.Run()

	require.ErrorContains(t, err, "error: on stderr")
	require.ErrorContains(t, err, "error: on stdout")
	require.Equal(t, 2001, strings.Count(stdout.String(), "\n"))
	require.Equal(t, 2001, strings.Count(stderr.String(), "\n"))
}
