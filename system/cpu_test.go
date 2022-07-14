package system

import (
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// These tests aren't perfectly isolated, because they depend on the environment
// they're being run in, but they are better than having nothing.

func TestCpu_Architecture(t *testing.T) {
	// Given
	unameCmd := exec.Command("uname", "-m")
	unameOut, err := unameCmd.Output()
	if err != nil {
		t.Errorf("uname failed: %v", err)
	}
	unameOutString := string(unameOut)
	unameOutString = strings.TrimSpace(unameOutString)

	// When
	architecture, err := CPU.Architecture()
	if err != nil {
		t.Errorf("CPU.Architecture() failed: %v", err)
	}

	// Then
	require.Equal(t, unameOutString, architecture)
}

func TestCpu_IsARM(t *testing.T) {
	// Given
	unameCmd := exec.Command("uname", "-m")
	unameOut, err := unameCmd.Output()
	if err != nil {
		t.Errorf("uname failed: %v", err)
	}
	unameOutString := string(unameOut)
	unameOutString = strings.TrimSpace(unameOutString)

	// When
	isARM, err := CPU.IsARM()
	if err != nil {
		t.Errorf("CPU.IsARM() failed: %v", err)
	}

	// Then
	require.Equal(t, strings.Contains(unameOutString, "arm"), isARM)
}
