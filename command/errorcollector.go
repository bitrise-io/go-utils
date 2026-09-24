package command

import "sync"

// errorCollector runs the ErrorFinder over the output. When Stdout and Stderr are different writers, os/exec copies
// the two pipes on two goroutines, so Write must be goroutine-safe.
type errorCollector struct {
	mu          sync.Mutex
	errorLines  []string
	errorFinder ErrorFinder
}

func (e *errorCollector) Write(p []byte) (n int, err error) {
	e.collectErrors(string(p))
	return len(p), nil
}

func (e *errorCollector) collectErrors(output string) {
	lines := e.errorFinder(output)
	if len(lines) == 0 {
		return
	}

	e.mu.Lock()
	defer e.mu.Unlock()
	e.errorLines = append(e.errorLines, lines...)
}

func (e *errorCollector) collectedErrorLines() []string {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.errorLines
}
