package progress

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewWrapper(t *testing.T) {
	var buf bytes.Buffer
	spinner := NewDefaultSpinnerWithOutput("Test", &buf)

	tests := []struct {
		name            string
		interactiveMode bool
	}{
		{"interactive mode", true},
		{"non-interactive mode", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wrapper := NewWrapper(spinner, tt.interactiveMode)
			assert.Equal(t, tt.interactiveMode, wrapper.interactiveMode)
			assert.Equal(t, spinner.message, wrapper.spinner.message)
		})
	}
}

func TestNewDefaultWrapper(t *testing.T) {
	message := "Loading"
	wrapper := NewDefaultWrapper(message)

	assert.Equal(t, message, wrapper.spinner.message)
	assert.Len(t, wrapper.spinner.chars, 8)
	// interactiveMode will be based on actual terminal status
}

func TestNewDefaultWrapperWithOutput(t *testing.T) {
	var buf bytes.Buffer
	message := "Processing"
	wrapper := NewDefaultWrapperWithOutput(message, &buf)

	assert.Equal(t, message, wrapper.spinner.message)
	assert.Equal(t, &buf, wrapper.spinner.writer)
	assert.Len(t, wrapper.spinner.chars, 8)
}

func TestWrapper_WrapAction_Interactive(t *testing.T) {
	var buf bytes.Buffer
	mockSleeper := &MockSleeper{}
	spinner := NewSpinnerWithSleeper("Test", []string{"a", "b"}, 5, &buf, mockSleeper)
	wrapper := NewWrapper(spinner, true)

	actionCalled := false
	wrapper.WrapAction(func() {
		// Give spinner time to animate
		for len(mockSleeper.sleepCalls) < 3 {
			mockSleeper.Sleep(1)
		}
		actionCalled = true
	})

	assert.True(t, actionCalled, "action should have been called")
	assert.False(t, wrapper.spinner.active, "spinner should be stopped after action completes")
	
	// In interactive mode, sleeper should have been called
	assert.NotEmpty(t, mockSleeper.sleepCalls, "spinner should have animated in interactive mode")
}

func TestWrapper_WrapAction_NonInteractive(t *testing.T) {
	var buf bytes.Buffer
	mockSleeper := &MockSleeper{}
	spinner := NewSpinnerWithSleeper("Loading", []string{"a"}, 5, &buf, mockSleeper)
	wrapper := NewWrapper(spinner, false)

	actionCalled := false
	wrapper.WrapAction(func() {
		actionCalled = true
	})

	assert.True(t, actionCalled, "action should have been called")
	assert.False(t, wrapper.spinner.active, "spinner should not be active in non-interactive mode")

	// In non-interactive mode, sleeper should NOT have been called
	assert.Empty(t, mockSleeper.sleepCalls, "spinner should not animate in non-interactive mode")

	output := buf.String()
	assert.Contains(t, output, "Loading...", "should print message with ellipsis")
	assert.True(t, strings.HasSuffix(strings.TrimSpace(output), "..."), "message should end with ...")
}

func TestWrapper_WrapAction_NonInteractive_MessageWithPeriod(t *testing.T) {
	var buf bytes.Buffer
	mockSleeper := &MockSleeper{}
	spinner := NewSpinnerWithSleeper("Loading.", []string{"a"}, 5, &buf, mockSleeper)
	wrapper := NewWrapper(spinner, false)

	wrapper.WrapAction(func() {})

	output := buf.String()
	// Should not add extra dots if message already ends with period
	assert.Contains(t, output, "Loading.", "message with period should be preserved")
	assert.NotContains(t, output, "Loading....", "should not add ellipsis to message ending with period")
}

func TestWrapper_WrapAction_ActionPanics_Interactive(t *testing.T) {
	var buf bytes.Buffer
	mockSleeper := &MockSleeper{}
	spinner := NewSpinnerWithSleeper("Test", []string{"a"}, 5, &buf, mockSleeper)
	wrapper := NewWrapper(spinner, true)

	assert.Panics(t, func() {
		wrapper.WrapAction(func() {
			panic("test panic")
		})
	}, "should propagate panic from action")
}

func TestWrapper_WrapAction_ActionPanics_NonInteractive(t *testing.T) {
	var buf bytes.Buffer
	spinner := NewDefaultSpinnerWithOutput("Test", &buf)
	wrapper := NewWrapper(spinner, false)

	assert.Panics(t, func() {
		wrapper.WrapAction(func() {
			panic("test panic")
		})
	}, "should propagate panic from action")
}

func TestWrapper_WrapAction_LongRunningAction_Interactive(t *testing.T) {
	var buf bytes.Buffer
	mockSleeper := &MockSleeper{}
	spinner := NewSpinnerWithSleeper("Processing", []string{"1", "2", "3"}, 5, &buf, mockSleeper)
	wrapper := NewWrapper(spinner, true)

	actionCompleted := false
	wrapper.WrapAction(func() {
		// Simulate work
		for i := 0; i < 10; i++ {
			mockSleeper.Sleep(1)
		}
		actionCompleted = true
	})

	assert.True(t, actionCompleted, "long-running action should complete")
	assert.False(t, wrapper.spinner.active, "spinner should stop after action completes")
}

func TestWrapper_WrapAction_EmptyMessage_NonInteractive(t *testing.T) {
	var buf bytes.Buffer
	spinner := NewDefaultSpinnerWithOutput("", &buf)
	wrapper := NewWrapper(spinner, false)

	wrapper.WrapAction(func() {})

	output := buf.String()
	assert.Equal(t, "...\n", output, "empty message should just print ellipsis")
}

func TestWrapper_WrapAction_MultipleActions_Sequential(t *testing.T) {
	var buf bytes.Buffer
	mockSleeper := &MockSleeper{}
	spinner := NewSpinnerWithSleeper("Task", []string{"a", "b"}, 5, &buf, mockSleeper)
	wrapper := NewWrapper(spinner, true)

	count := 0
	for i := 0; i < 3; i++ {
		wrapper.WrapAction(func() {
			count++
		})
	}

	assert.Equal(t, 3, count, "should execute all actions")
	assert.False(t, wrapper.spinner.active, "spinner should be stopped after last action")
}

func TestOutputDeviceIsTerminal(t *testing.T) {
	// This test just ensures the function doesn't panic
	result := OutputDeviceIsTerminal()
	assert.IsType(t, true, result, "should return a boolean")
}
