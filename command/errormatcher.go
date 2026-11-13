package command

import "strings"

const (
	// maxReasonBytes leaves room for the layers a step wraps around the reason inside the 1 KB cap of
	// BITRISE_FAILED_STEP_ERROR_MESSAGE.
	maxReasonBytes = 800

	truncatedMarker = "... (truncated)"

	// reasonBudget is held back for the marker up front, so a cut reason stays within maxReasonBytes.
	reasonBudget = maxReasonBytes - len("\n"+truncatedMarker)
)

// ErrorMatcher reads the output of one tool and reports why it failed. It is given complete lines, in
// the order the tool wrote them, so that a matcher holds the rules of its tool and no stream handling
// of its own. NewLineObserver turns one into an OutputObserver.
type ErrorMatcher interface {
	// ConsumeLine is given one line of output without its newline, as the tool wrote it. A tool that
	// colors its output leaves the escape sequences in, so a matcher for one has to take them out.
	ConsumeLine(line string)
	// Reasons returns the lines that say why the command failed, most telling first, or nothing when
	// the output held none. Nothing is the right answer for a failure the matcher does not recognise:
	// the command then reports its generic message rather than a guess.
	Reasons() []string
}

// capReasons cuts what the observers reported to maxReasonBytes. Whole lines are dropped at the cap;
// only a first line too long on its own is cut, because a reason that says nothing is worse than one
// that stops mid-sentence.
func capReasons(reasons []string) []string {
	var kept []string
	remaining := reasonBudget

	for _, line := range reasons {
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
