package system

import (
	"strings"

	"golang.org/x/sys/unix"
)

type cpu struct{}

// CPU is a singleton that returns information about the CPU
var CPU cpu

// Architecture returns the CPU's architecture name as a string
//
// Some example return values:
//   - x86_64
//   - arm64
func (cpu) Architecture() (string, error) {
	var utsname unix.Utsname
	err := unix.Uname(&utsname)
	if err != nil {
		return "", err
	}

	architecture := string(utsname.Machine[:])
	architecture = strings.Trim(architecture, "\x00")
	return architecture, nil
}

// IsARM returns true if the CPU's architecture is ARM
//
// Specifically, it returns true if the architecture's name contains the string 'arm'
func (cpu) IsARM() (bool, error) {
	cpuType, err := CPU.Architecture()

	if err != nil {
		return false, err
	}

	return strings.Contains(cpuType, "arm"), nil
}
