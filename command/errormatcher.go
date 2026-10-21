package command

import (
	"io"
	"strings"
	"sync"
)

const (
	// maxReasonBytes leaves room for the layers a step wraps around the reason inside the 1 KB cap of
	// BITRISE_FAILED_STEP_ERROR_MESSAGE.
	maxReasonBytes = 800

	truncatedMarker = "... (truncated)"
)

// ErrorMatcher reads the output of one tool and reports why it failed. It is given complete lines, in
// the order the tool wrote them, so that a matcher holds the rules of its tool and no stream handling
// of its own.
type ErrorMatcher interface {
	// ConsumeLine is given one line of output without its newline, as the tool wrote it. A tool that
	// colors its output leaves the escape sequences in, so a matcher for one has to take them out.
	ConsumeLine(line string)
	// Reasons returns the lines that say why the command failed, most telling first, or nothing when
	// the output held none. Nothing is the right answer for a failure the matcher does not recognise:
	// the command then reports its generic message rather than a guess.
	Reasons() []string
}

// matcherCollector feeds an ErrorMatcher the lines of a command's output.
//
// When Stdout and Stderr are different writers, os/exec copies the two pipes on two goroutines, so
// the matcher runs under a lock and is never entered twice at once.
type matcherCollector struct {
	mu      sync.Mutex
	matcher ErrorMatcher
}

func (c *matcherCollector) consumeLine(line []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.matcher.ConsumeLine(string(line))
}

// reasons returns why the command failed, cut to maxReasonBytes. Whole lines are dropped at the cap;
// only a first line too long on its own is cut, because a reason that says nothing is worse than one
// that stops mid-sentence.
func (c *matcherCollector) reasons() []string {
	c.mu.Lock()
	defer c.mu.Unlock()

	var kept []string
	remaining := maxReasonBytes

	for _, line := range c.matcher.Reasons() {
		if len(line) > remaining {
			if len(kept) == 0 {
				kept = append(kept, strings.ToValidUTF8(line[:remaining], ""))
			}
			kept = append(kept, truncatedMarker)

			break
		}

		kept = append(kept, line)
		remaining -= len(line) + 1
	}

	return kept
}

// collectFromWriter wraps the stream's writer so that the matcher sees its lines.
func (c *matcherCollector) collectFromWriter(dest io.Writer) *lineWriter {
	return &lineWriter{dest: dest, collector: c}
}
