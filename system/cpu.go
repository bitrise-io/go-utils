package system

import (
	"strings"

	"golang.org/x/sys/unix"
)

type cpu struct{}

var CPU cpu

// Architecture returns the CPU's architecture name as a string
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
	} else {
		return strings.Contains(cpuType, "arm"), nil
	}
}