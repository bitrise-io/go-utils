//go:build !race
// +build !race

package progress

import (
	"bytes"
	"fmt"
	"io"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewWrapper(t *testing.T) {
	message := "loading"
	spinner := NewDefaultSpinner(message)

	isInteractiveMode := true
	NewWrapper(spinner, isInteractiveMode).WrapAction(func() {
		time.Sleep(2 * time.Second)
	})
}

func TestNewDefaultWrapper(t *testing.T) {
	message := "loading"
	NewDefaultWrapper(message).WrapAction(func() {
		time.Sleep(2 * time.Second)
	})
}

func TestNewDefaultWrapperWithOutput(t *testing.T) {
	message := "loading"

	var b bytes.Buffer
	NewDefaultWrapperWithOutput(message, io.Writer(&b)).WrapAction(func() {
		time.Sleep(2 * time.Second)
	})

	expected := fmt.Sprintf("%s...\n", message)
	got := b.String()
	assert.Equal(t, expected, got)
}
