// Package ansi removes the terminal control sequences from a command's output, so that what is left
// is the text a reader would see on screen.
package ansi

import "strings"

// Strip removes every escape sequence and control byte from s, keeping newlines, tabs and every
// printable rune, so a path, a message or a non-ASCII identifier is never touched.
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
			end, res := scanAt(s, i)
			switch res {
			case scanComplete:
				i = end
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

func needsStripping(s string) bool {
	for i := 0; i < len(s); i++ {
		if !keepsByte(s[i]) || s[i] == '\r' {
			return true
		}
	}

	return false
}

// keepsByte reports whether b survives stripping on its own: every printable byte, every byte of a
// multi-byte rune, and the whitespace controls a log line is built from.
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
