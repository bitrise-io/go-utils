package command

import (
	"syscall"
	"testing"
	"time"

	"github.com/bitrise-io/go-utils/v2/env"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSignalBeforeStart(t *testing.T) {
	factory := NewFactory(env.NewRepository())
	cmd := factory.Create("sleep", []string{"30"}, nil)

	assert.ErrorIs(t, cmd.Kill(), ErrProcessNotStarted)
	assert.ErrorIs(t, cmd.Signal(syscall.SIGTERM), ErrProcessNotStarted)
}

func TestKillStopsRunningCommand(t *testing.T) {
	factory := NewFactory(env.NewRepository())
	cmd := factory.Create("sleep", []string{"30"}, nil)

	require.NoError(t, cmd.Start())

	waitErr := make(chan error, 1)
	go func() { waitErr <- cmd.Wait() }()

	require.NoError(t, cmd.Kill())

	select {
	case err := <-waitErr:
		assert.Error(t, err, "a killed command should report a non-nil wait error")
	case <-time.After(10 * time.Second):
		t.Fatal("killed command did not exit")
	}
}

func TestSignalAfterExitReportsProcessFinished(t *testing.T) {
	factory := NewFactory(env.NewRepository())
	cmd := factory.Create("true", nil, nil)

	require.NoError(t, cmd.Start())
	require.NoError(t, cmd.Wait())

	assert.ErrorIs(t, cmd.Kill(), ErrProcessFinished)
}

func TestSignalDeliversTermToCommand(t *testing.T) {
	factory := NewFactory(env.NewRepository())
	cmd := factory.Create("bash", []string{"testdata/trap_term.sh"}, nil)

	require.NoError(t, cmd.Start())

	waitErr := make(chan error, 1)
	go func() { waitErr <- cmd.Wait() }()

	require.Eventually(t, func() bool {
		return cmd.Signal(syscall.SIGTERM) == nil
	}, 10*time.Second, 50*time.Millisecond, "could not deliver SIGTERM")

	select {
	case err := <-waitErr:
		// The script traps SIGTERM and exits 0, so a clean exit proves the signal was handled.
		assert.NoError(t, err)
	case <-time.After(10 * time.Second):
		t.Fatal("command did not exit after SIGTERM")
	}
}
