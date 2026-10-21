package linepipe

import "github.com/bitrise-io/go-utils/v2/ansi"

// Tap reads a line without changing what is forwarded: error finding, log level detection.
//
// A tap that is given to the writers of both output streams is called from the two copy goroutines
// os/exec runs, so it has to guard its own state.
type Tap interface {
	Observe(line *Line)
}

// Line is one line of a stream, without its newline. Its views are parsed on first use and kept, so
// every stage over the same line shares one parse of it.
type Line struct {
	raw  []byte
	text string
	seqs []ansi.Sequence

	hasText bool
	hasSeqs bool
}

// Raw returns the bytes of the line as they were written, escape sequences included. The writer
// reuses the slice for the next line, so a stage that keeps it has to copy it.
func (l *Line) Raw() []byte {
	return l.raw
}

// Text returns the line with its escape sequences and control bytes removed, so the words a reader
// would see on screen.
func (l *Line) Text() string {
	if !l.hasText {
		l.text = ansi.Strip(string(l.raw))
		l.hasText = true
	}

	return l.text
}

// Sequences returns the escape sequences in the line with their positions, for a stage that reads
// what a color means. A sequence the line does not terminate is left out.
func (l *Line) Sequences() []ansi.Sequence {
	if !l.hasSeqs {
		l.seqs, _ = ansi.Sequences(string(l.raw))
		l.hasSeqs = true
	}

	return l.seqs
}

func (l *Line) reset(raw []byte) {
	l.raw = raw
	l.text = ""
	l.seqs = nil
	l.hasText = false
	l.hasSeqs = false
}
