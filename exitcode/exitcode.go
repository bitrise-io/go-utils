//go:generate stringer -type=ExitCode

package exitcode

// ExitCode is a simple integer type that represents an exit code.
// It can be used to provide a more semantically meaningful exit code than a simple integer.
type ExitCode int

const (
	// Success indicates that the program exited successfully.
	Success ExitCode = 0

	// Failure indicates that the program exited unsuccessfully.
	Failure ExitCode = 1
)
