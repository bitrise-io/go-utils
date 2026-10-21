// Package ansi reads the terminal control sequences in a command's output, so that a caller can tell
// what a sequence means and remove what a reader would not see on screen.
package ansi

import "strings"

// Strip removes every escape sequence and control byte from s, keeping newlines, tabs and the text a
// terminal would leave on screen.
//
// Within a line the text before the last carriage return is dropped, because a redraw overwrites what
// came before it. A carriage return directly before a newline only ends the line and keeps its text.
func Strip(s string) string {
	if !needsStripping(s) {
		return s
	}

	var b strings.Builder
	b.Grow(len(s))

	for i := 0; i < len(s); {
		if s[i] == escape {
			seq, res := scanAt(s, i)
			switch res {
			case scanComplete:
				i = seq.End
			case scanUnterminated:
				i = len(s)
			case scanInvalid:
				i++
			}

			continue
		}

		if keepsByte(s[i]) {
			b.WriteByte(s[i])
		}
		i++
	}

	return applyCarriageReturns(b.String())
}

// Sequences reports the escape sequences in s with their positions, so that a caller can find where a
// color block starts and ends instead of matching text.
//
// pending is the index of a sequence that starts in s but is not terminated in it, or -1. A caller
// reading a stream carries s[pending:] over to its next chunk.
func Sequences(s string) (seqs []Sequence, pending int) {
	for i := 0; i < len(s); {
		if s[i] != escape {
			i++
			continue
		}

		seq, res := scanAt(s, i)
		switch res {
		case scanComplete:
			seqs = append(seqs, seq)
			i = seq.End
		case scanUnterminated:
			return seqs, i
		case scanInvalid:
			i++
		}
	}

	return seqs, -1
}

func needsStripping(s string) bool {
	for i := 0; i < len(s); i++ {
		if !keepsByte(s[i]) || s[i] == '\r' {
			return true
		}
	}

	return false
}

// keepsByte reports whether b survives stripping on its own: every printable byte, every byte of a
// multi-byte rune, and the two whitespace controls a log line is built from.
func keepsByte(b byte) bool {
	return b >= 0x20 && b != deleteByte || b == '\n' || b == '\t' || b == '\r'
}

func applyCarriageReturns(s string) string {
	if !strings.Contains(s, "\r") {
		return s
	}

	lines := strings.Split(s, "\n")
	for i, line := range lines {
		line = strings.TrimSuffix(line, "\r")
		if idx := strings.LastIndex(line, "\r"); idx >= 0 {
			line = line[idx+1:]
		}
		lines[i] = line
	}

	return strings.Join(lines, "\n")
}
