package command

import (
	"slices"
	"sync"
)

// errorCollector runs the ErrorFinder over the output. When Stdout and Stderr are different writers, os/exec copies
// the two pipes on two goroutines, so Observe must be goroutine-safe. The finder runs under the same lock, so a finder
// that carries state across calls is never entered twice at once.
type errorCollector struct {
	mu          sync.Mutex
	errorLines  []string
	errorFinder ErrorFinder
}

// Observe makes the collector the observer of a tap: the finder is run over each chunk as it was written.
func (e *errorCollector) Observe(p []byte) {
	e.collectErrors(string(p))
}

// Close does nothing: the finder is called per chunk and holds nothing back.
func (e *errorCollector) Close() {}

// Reasons returns the lines the finder collected.
func (e *errorCollector) Reasons() []string {
	return e.collectedErrorLines()
}

func (e *errorCollector) collectErrors(output string) {
	e.mu.Lock()
	defer e.mu.Unlock()

	lines := e.errorFinder(output)
	if len(lines) == 0 {
		return
	}

	e.errorLines = append(e.errorLines, lines...)
}

// collectedErrorLines returns a copy: the collected lines outlive the lock in the returned error, and the slice keeps
// growing while the copy goroutines run.
func (e *errorCollector) collectedErrorLines() []string {
	e.mu.Lock()
	defer e.mu.Unlock()

	return slices.Clone(e.errorLines)
}
