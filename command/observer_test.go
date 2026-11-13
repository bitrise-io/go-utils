package command

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/bitrise-io/go-utils/v2/env"
	"github.com/stretchr/testify/require"
)

// bothStreamsScript prints one error line on each stream and fails.
const bothStreamsScript = `echo "noise on stdout"; echo "error: on stdout"; echo "noise on stderr" >&2; echo "error: on stderr" >&2; exit 1`

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

func TestObserver_Run(t *testing.T) {
	var out bytes.Buffer
	observers := &observerFactory{}

	err := runScript(script(bothStreamsScript, &Opts{Stdout: &out, NewOutputObserver: observers.new}))

	require.Equal(t, "error: on stdout\nerror: on stderr", collectedErrorOutput(t, err))
	require.Len(t, observers.built, 1, "one pipe, one observer")
	require.Equal(t, out.String(), observers.built[0].seen.String(), "the observer is shown what was forwarded")
	require.Equal(t, 1, observers.built[0].closed)
}

// TestObserver_RunKeepsTheOrderAcrossStreams is why the command pipes the output itself. The caller
// passes two different writers, which without an observer means two pipes and no order between them.
// Neither writer nor the observer is goroutine-safe, so the race detector also proves there is one
// copy goroutine.
func TestObserver_RunKeepsTheOrderAcrossStreams(t *testing.T) {
	alternating, want := alternatingScript(20)

	for run := 0; run < 20; run++ {
		var stdout, stderr bytes.Buffer
		observers := &observerFactory{}

		err := runScript(script(alternating, &Opts{Stdout: &stdout, Stderr: &stderr, NewOutputObserver: observers.new}))

		require.Error(t, err)
		require.Equal(t, want, strings.Fields(observers.built[0].seen.String()), "what the observer read, run %d", run)
		require.Equal(t, want, strings.Fields(stdout.String()), "what was forwarded, run %d", run)
		require.Empty(t, stderr.String(), "both streams are forwarded to Stdout")
	}
}

func TestObserver_RunForwardsToStderrWhenOnlyThatIsSet(t *testing.T) {
	var stderr bytes.Buffer
	observers := &observerFactory{}

	err := runScript(script(`echo out; echo err >&2; exit 1`, &Opts{Stderr: &stderr, NewOutputObserver: observers.new}))

	require.Error(t, err)
	require.Equal(t, "out\nerr\n", stderr.String())
}

func TestObserver_RunWithoutWriters(t *testing.T) {
	observers := &observerFactory{}

	err := runScript(script(bothStreamsScript, &Opts{NewOutputObserver: observers.new}))

	require.Equal(t, "error: on stdout\nerror: on stderr", collectedErrorOutput(t, err))
}

func TestObserver_SuccessfulRunReportsNothing(t *testing.T) {
	var out bytes.Buffer
	observers := &observerFactory{}

	err := runScript(script(`echo "error: but the tool went on"; printf 'no newline'`, &Opts{Stdout: &out, NewOutputObserver: observers.new}))

	require.NoError(t, err)
	require.Equal(t, "error: but the tool went on\nno newline", out.String())
	require.Equal(t, 1, observers.built[0].closed)
}

func TestObserver_NothingFoundKeepsTheGenericMessage(t *testing.T) {
	observers := &observerFactory{}

	err := runScript(script(`echo "all quiet"; exit 1`, &Opts{NewOutputObserver: observers.new}))

	require.Equal(t, "check the command's output for details", collectedErrorOutput(t, err))
}

func TestObserver_RunAndReturnExitCode(t *testing.T) {
	observers := &observerFactory{}
	cmd := script(`echo "error: exit code run"; exit 3`, &Opts{NewOutputObserver: observers.new})

	code, err := cmd.RunAndReturnExitCode()

	require.Equal(t, 3, code)
	require.Equal(t, "error: exit code run", collectedErrorOutput(t, err))
	require.Equal(t, 1, observers.built[0].closed)
}

func TestObserver_StartAndWait(t *testing.T) {
	var out bytes.Buffer
	observers := &observerFactory{}
	cmd := script(bothStreamsScript, &Opts{Stdout: &out, NewOutputObserver: observers.new})

	require.NoError(t, cmd.Start())
	err := cmd.Wait()

	require.Equal(t, "error: on stdout\nerror: on stderr", collectedErrorOutput(t, err))
	require.Len(t, observers.built, 1)
	require.Equal(t, out.String(), observers.built[0].seen.String())
	require.Equal(t, 1, observers.built[0].closed)
}

func TestObserver_RunAndReturnTrimmedCombinedOutput(t *testing.T) {
	alternating, want := alternatingScript(20)
	observers := &observerFactory{}
	cmd := script(alternating, &Opts{NewOutputObserver: observers.new})

	out, err := cmd.RunAndReturnTrimmedCombinedOutput()

	require.Error(t, err)
	require.Len(t, observers.built, 1, "one pipe, one observer")
	require.Equal(t, want, strings.Fields(out), "the output is returned in the order it was written")
	require.Equal(t, want, strings.Fields(observers.built[0].seen.String()))
	require.Equal(t, 1, observers.built[0].closed)
}

func TestObserver_RunAndReturnTrimmedCombinedOutputReportsReasons(t *testing.T) {
	observers := &observerFactory{}
	cmd := script(bothStreamsScript, &Opts{NewOutputObserver: observers.new})

	out, err := cmd.RunAndReturnTrimmedCombinedOutput()

	require.Equal(t, "noise on stdout\nerror: on stdout\nnoise on stderr\nerror: on stderr", out)
	require.Equal(t, "error: on stdout\nerror: on stderr", collectedErrorOutput(t, err))
}

// TestObserver_RunAndReturnTrimmedOutput pins the one method that cannot have one pipe: stdout is the
// value the caller asked for, so each stream is observed on its own.
func TestObserver_RunAndReturnTrimmedOutput(t *testing.T) {
	t.Run("each stream has its own observer, stderr reported first", func(t *testing.T) {
		observers := &observerFactory{}
		cmd := script(bothStreamsScript, &Opts{NewOutputObserver: observers.new})

		out, err := cmd.RunAndReturnTrimmedOutput()

		require.Equal(t, "noise on stdout\nerror: on stdout", out, "stdout is returned without stderr")
		require.Equal(t, "error: on stderr\nerror: on stdout", collectedErrorOutput(t, err))
		require.Len(t, observers.built, 2)
		require.Equal(t, "noise on stderr\nerror: on stderr\n", observers.built[0].seen.String())
		require.Equal(t, "noise on stdout\nerror: on stdout\n", observers.built[1].seen.String())
		require.Equal(t, 1, observers.built[0].closed)
		require.Equal(t, 1, observers.built[1].closed)
	})

	t.Run("stderr still reaches the caller's writer", func(t *testing.T) {
		var stderr bytes.Buffer
		observers := &observerFactory{}
		cmd := script(bothStreamsScript, &Opts{Stderr: &stderr, NewOutputObserver: observers.new})

		out, err := cmd.RunAndReturnTrimmedOutput()

		require.Equal(t, "noise on stdout\nerror: on stdout", out)
		require.Equal(t, "error: on stderr\nerror: on stdout", collectedErrorOutput(t, err))
		require.Equal(t, "noise on stderr\nerror: on stderr\n", stderr.String())
	})

	t.Run("nothing found falls back to the captured stderr, as it does without an observer", func(t *testing.T) {
		const quiet = `echo "all quiet"; echo "noise on stderr" >&2; exit 1`

		_, withObserver := script(quiet, &Opts{NewOutputObserver: (&observerFactory{}).new}).RunAndReturnTrimmedOutput()
		_, without := script(quiet, nil).RunAndReturnTrimmedOutput()

		require.Equal(t, "noise on stderr\n", collectedErrorOutput(t, withObserver))
		require.Equal(t, without.Error(), withObserver.Error())
	})

	t.Run("a Stdout the caller set is refused, as it is without an observer", func(t *testing.T) {
		var stdout bytes.Buffer

		_, withObserver := script(`exit 1`, &Opts{Stdout: &stdout, NewOutputObserver: (&observerFactory{}).new}).RunAndReturnTrimmedOutput()
		_, without := script(`exit 1`, &Opts{Stdout: &stdout}).RunAndReturnTrimmedOutput()

		require.ErrorContains(t, withObserver, "Stdout already set")
		require.Equal(t, without.Error(), withObserver.Error())
	})
}

// TestObserver_EachCommandGetsItsOwn shows why Opts holds a function and not an observer: an Opts
// reused for a retry carries nothing from one run into the next.
func TestObserver_EachCommandGetsItsOwn(t *testing.T) {
	observers := &observerFactory{}
	opts := &Opts{NewOutputObserver: observers.new}

	first := runScript(script(`echo "error: first attempt"; exit 1`, opts))
	second := runScript(script(`echo "error: second attempt"; exit 1`, opts))

	require.Equal(t, "error: first attempt", collectedErrorOutput(t, first))
	require.Equal(t, "error: second attempt", collectedErrorOutput(t, second))
	require.Len(t, observers.built, 2)
}

func TestObserver_TakesThePlaceOfTheErrorFinder(t *testing.T) {
	observers := &observerFactory{}
	opts := &Opts{
		NewOutputObserver: observers.new,
		ErrorFinder:       func(string) []string { return []string{"from the chunk finder"} },
	}

	err := runScript(script(`echo "error: the real one"; exit 1`, opts))

	require.Equal(t, "error: the real one", collectedErrorOutput(t, err))
}

func TestObserver_ACommandThatCannotStart(t *testing.T) {
	observers := &observerFactory{}
	cmd := NewFactory(env.NewRepository()).Create("/no/such/tool", nil, &Opts{NewOutputObserver: observers.new})

	err := cmd.Run()

	require.ErrorContains(t, err, "executing command failed")
	require.Equal(t, 1, observers.built[0].closed)
}

func TestHeadTailBuffer(t *testing.T) {
	t.Run("keeps everything that fits", func(t *testing.T) {
		b := &headTailBuffer{limit: 8}
		_, _ = b.Write([]byte("0123"))
		_, _ = b.Write([]byte("456789ab"))

		require.Equal(t, "0123456789ab", string(b.Bytes()))
	})

	t.Run("drops the middle and says how much", func(t *testing.T) {
		b := &headTailBuffer{limit: 4}
		for _, chunk := range []string{"01", "2345", "6789abcdef", "gh", "i"} {
			n, err := b.Write([]byte(chunk))
			require.NoError(t, err)
			require.Equal(t, len(chunk), n)
		}

		require.Equal(t, "0123\n... omitting 11 bytes ...\nfghi", string(b.Bytes()))
	})

	t.Run("stays bounded", func(t *testing.T) {
		b := &headTailBuffer{limit: 16}
		for i := 0; i < 10000; i++ {
			_, _ = b.Write([]byte("some stderr output\n"))
		}

		require.LessOrEqual(t, len(b.head)+len(b.tail), 32)
	})
}

func script(source string, opts *Opts) Command {
	return NewFactory(env.NewRepository()).Create("sh", []string{"-c", source}, opts)
}

func runScript(cmd Command) error {
	return cmd.Run()
}

// observerFactory hands out a new observer on every call and remembers them in the order they were
// asked for.
type observerFactory struct {
	built []*prefixObserver
}

func (f *observerFactory) new() OutputObserver {
	observer := &prefixObserver{}
	f.built = append(f.built, observer)

	return observer
}

// prefixObserver reports the lines that start with "error:". It is not goroutine-safe.
type prefixObserver struct {
	seen   strings.Builder
	closed int
}

func (o *prefixObserver) Observe(p []byte) { o.seen.Write(p) }
func (o *prefixObserver) Close()           { o.closed++ }

func (o *prefixObserver) Reasons() []string {
	var reasons []string
	for _, line := range strings.Split(o.seen.String(), "\n") {
		if strings.HasPrefix(line, "error:") {
			reasons = append(reasons, line)
		}
	}

	return reasons
}
