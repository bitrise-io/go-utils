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

// splitBlockScript writes a two-line error block on stderr with one stdout line in between, which is what a Step
// looks like when xcodebuild reports an error while its progress output keeps coming.
const splitBlockScript = `printf 'xcodebuild: error: Signing requires a development team\n' >&2
printf 'note: Building targets in dependency order\n'
printf '    Reason: No profiles matched\n' >&2
exit 1`

// findErrorBlocks collects an "error:" line together with the indented lines that follow it, and stops at the first
// line that belongs to neither. It is the shape of go-xcode's FindXcodebuildErrors: a block is only recognised while
// its lines are adjacent within one chunk.
func findErrorBlocks(out string) []string {
	var lines []string
	inBlock := false

	for _, line := range strings.Split(out, "\n") {
		switch {
		case strings.Contains(line, "error:"):
			lines = append(lines, line)
			inBlock = true
		case inBlock && strings.HasPrefix(line, "    "):
			lines = append(lines, line)
		default:
			inBlock = false
		}
	}

	return lines
}

// TestErrorFinder_SharedStreamSplitsErrorBlocks pins the trade-off the shared writer brings: one pipe means the
// child's stdout and stderr bytes interleave in the order it wrote them, so output from the other stream can land
// inside a multi-line error block and terminate it early. The block's first line is still reported -- the tail is
// not. Keeping the streams apart avoids it, at the cost of the ordering and the data race a shared writer had.
func TestErrorFinder_SharedStreamSplitsErrorBlocks(t *testing.T) {
	var out bytes.Buffer
	cmd := NewFactory(env.NewRepository()).Create("sh", []string{"-c", splitBlockScript}, &Opts{
		Stdout:      &out,
		Stderr:      &out,
		ErrorFinder: findErrorBlocks,
	})

	err := cmd.Run()

	require.Equal(t, "xcodebuild: error: Signing requires a development team", collectedErrorOutput(t, err),
		"the interleaved stdout line ends the block, so its tail is lost")
	require.Equal(t, "xcodebuild: error: Signing requires a development team\n"+
		"note: Building targets in dependency order\n"+
		"    Reason: No profiles matched\n", out.String())
}

// contiguousBlockScript writes the same error block in one stderr write, with the stdout line before it rather
// than inside it.
const contiguousBlockScript = `printf 'note: Building targets in dependency order\n'
printf 'xcodebuild: error: Signing requires a development team\n    Reason: No profiles matched\n' >&2
exit 1`

// TestErrorFinder_SharedStreamKeepsContiguousErrorBlocks is the same wiring and the same finder as above, with the
// block written in one go: it survives intact. Sharing the writer is not what loses the tail -- other-stream bytes
// landing between the block's lines is, which only a shared stream can do.
func TestErrorFinder_SharedStreamKeepsContiguousErrorBlocks(t *testing.T) {
	var out bytes.Buffer
	cmd := NewFactory(env.NewRepository()).Create("sh", []string{"-c", contiguousBlockScript}, &Opts{
		Stdout:      &out,
		Stderr:      &out,
		ErrorFinder: findErrorBlocks,
	})

	err := cmd.Run()

	require.Equal(t, "xcodebuild: error: Signing requires a development team\n    Reason: No profiles matched",
		collectedErrorOutput(t, err))
}

// collectedErrorOutput returns the lines the ErrorFinder collected, which ExitStatusError appends after the
// printable command. The command itself is quoted in the message, so matching on the whole string would also match
// the script that produced the output.
func collectedErrorOutput(t *testing.T, err error) string {
	t.Helper()

	var exitErr *ExitStatusError
	require.ErrorAs(t, err, &exitErr)

	msg := exitErr.Error()
	i := strings.LastIndex(msg, "): ")
	require.NotEqual(t, -1, i, "unexpected error format: %s", msg)

	return msg[i+len("): "):]
}
