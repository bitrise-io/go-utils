package progress

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestSimpleProgress_Run(t *testing.T) {
	var buf bytes.Buffer
	ticker := newMockTicker()
	progress := NewSimpleDotsWithTicker(&buf, ticker)

	go func() {
		ticker.doTicks(3)
	}()
	// Start the progress and run a dummy action
	called := false
	err := progress.Run(func() error {
		// Simulate some work, allow ticker channel to be drained
		called = true
		time.Sleep(10 * time.Millisecond)
		return nil
	})
	require.NoError(t, err)
	require.True(t, called, "action should have been called")
	<-ticker.C() // ticker should be stopped after action completes

	// Expected output: dots + newline
	output := buf.String()
	require.True(t, strings.Contains(output, "."))
	require.True(t, strings.HasSuffix(output, "\n"))
}

func TestSimpleDots_RunTwice(t *testing.T) {
	var buf bytes.Buffer
	ticker := NewDefaultSimpleDots(&buf)

	err := ticker.Run(func() error {
		return nil
	})
	require.NoError(t, err)

	// Running again should return an error
	err = ticker.Run(func() error {
		return nil
	})
	require.Error(t, err)
	require.Equal(t, "progress can only be run once", err.Error())
}

func TestSimpleDots_RunError(t *testing.T) {
	var buf bytes.Buffer
	ticker := NewDefaultSimpleDots(&buf)

	err := ticker.Run(func() error {
		return fmt.Errorf("an error occurred")
	})
	require.Equal(t, fmt.Errorf("an error occurred"), err)
}

func TestSimpleDots_RunPanic(t *testing.T) {
	var buf bytes.Buffer
	ticker := NewDefaultSimpleDots(&buf)

	require.Panics(t, func() {
		err := ticker.Run(func() error {
			panic("something went wrong")
		})
		_ = err
	})
}
