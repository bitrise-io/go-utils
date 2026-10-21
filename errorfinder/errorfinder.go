// Package errorfinder collects the lines of a failing command's output that say why it failed, so
// that a tool's finder carries its rules and nothing else.
package errorfinder

import (
	"slices"
	"strings"
	"sync"

	"github.com/bitrise-io/go-utils/v2/linepipe"
)

const (
	// MaxReasonBytes leaves room for the layers a step wraps around the reason inside the 1 KB cap of
	// BITRISE_FAILED_STEP_ERROR_MESSAGE.
	MaxReasonBytes = 800

	truncatedMarker = "... (truncated)"
)

// Matcher reads the output of one tool and reports why it failed. It is given complete lines, in the
// order the tool wrote them, with escape sequences and control bytes already removed, so a matcher is
// the rules of its tool and no stream handling.
type Matcher interface {
	// ConsumeLine is given one line of output without its newline.
	ConsumeLine(line string)
	// Reasons returns the lines that say why the command failed, most telling first, or nothing when
	// the output held none. Nothing is the right answer for a tool whose failure was not recognised:
	// the caller then keeps its own message instead of showing a guess.
	Reasons() []string
}

// Collector feeds a Matcher the lines of a command's output and shortens what it reports.
//
// A command's two output streams are copied by two goroutines, so the matcher runs under a lock and
// is never entered twice at once.
type Collector struct {
	mu      sync.Mutex
	matcher Matcher
}

// NewCollector ...
func NewCollector(matcher Matcher) *Collector {
	return &Collector{matcher: matcher}
}

// Observe implements linepipe.Tap.
func (c *Collector) Observe(line *linepipe.Line) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.matcher.ConsumeLine(line.Text())
}

// Reasons returns why the command failed, with repeated lines left out and the whole cut to
// MaxReasonBytes.
func (c *Collector) Reasons() []string {
	c.mu.Lock()
	defer c.mu.Unlock()

	return capLines(dedupe(c.matcher.Reasons()))
}

func dedupe(lines []string) []string {
	seen := make(map[string]bool, len(lines))

	return slices.DeleteFunc(slices.Clone(lines), func(line string) bool {
		if seen[line] {
			return true
		}
		seen[line] = true

		return false
	})
}

// capLines keeps as many whole lines as fit. A first line that is too long on its own is cut, because a
// reason that says nothing is worse than one that stops mid-sentence.
func capLines(lines []string) []string {
	var kept []string
	remaining := MaxReasonBytes

	for _, line := range lines {
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
