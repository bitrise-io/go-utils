package progress

import (
	"fmt"
	"testing"

	"github.com/bitrise-io/go-utils/v2/progress/mocks"
	"github.com/stretchr/testify/require"
)

func TestSimpleProgress_Run(t *testing.T) {
	ticker := mocks.NewTicker()
	progress := NewSimpleDotsWithTicker(NewFmtPrinter(), ticker)

	// Start the progress and run a dummy action
	called := make(chan bool)
	err := progress.Run(func() error {
		require.True(t, progress.stopChan != nil, "stopChan should be initialized")
		ticker.DoTicks(10)
		close(called)
		return nil
	})
	require.NoError(t, err)

	_, ok := <-called
	require.False(t, ok, "action should have been called and channel closed")

	// check progress.stopChan is closed
	_, ok = <-progress.stopChan
	require.False(t, ok, "stopChan should be closed")
}

func TestSimpleDots_RunTwice(t *testing.T) {
	progress := NewDefaultSimpleDots(NewFmtPrinter())

	err := progress.Run(func() error {
		return nil
	})
	require.NoError(t, err)

	// Running again should return an error
	err = progress.Run(func() error {
		return nil
	})
	require.Error(t, err)
	require.Equal(t, "progress can only be run once", err.Error())
}

func TestSimpleDots_RunError(t *testing.T) {
	progress := NewDefaultSimpleDots(NewFmtPrinter())

	err := progress.Run(func() error {
		return fmt.Errorf("an error occurred")
	})

	require.EqualError(t, err, "an error occurred")
}

func TestSimpleDots_RunPanic(t *testing.T) {
	progress := NewDefaultSimpleDots(NewFmtPrinter())

	require.Panics(t, func() {
		_ = progress.Run(func() error {
			panic("something went wrong")
		})
	})

	_, ok := <-progress.stopChan
	require.False(t, ok, "stopChan should be closed")
}
