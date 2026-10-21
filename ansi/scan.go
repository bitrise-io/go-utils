package ansi

const (
	escape     = 0x1b
	bell       = 0x07
	deleteByte = 0x7f

	paramFirst, paramLast               = 0x30, 0x3f
	intermediateFirst, intermediateLast = 0x20, 0x2f
	csiFinalFirst, csiFinalLast         = 0x40, 0x7e
	escFinalFirst, escFinalLast         = 0x30, 0x7e
)

// scanResult tells the caller what to do with the bytes at an ESC: skip the sequence, drop the rest
// of the input with the sequence it never terminates, or treat the ESC as a stray byte and move on
// one byte.
type scanResult int

const (
	scanComplete scanResult = iota
	scanUnterminated
	scanInvalid
)

// scanAt reads the escape sequence starting at the ESC byte s[start] and reports where it ends.
func scanAt(s string, start int) (end int, res scanResult) {
	if start+1 >= len(s) {
		return 0, scanUnterminated
	}

	switch b := s[start+1]; {
	case b == '[':
		return scanCSI(s, start)
	case b == ']' || b == 'P' || b == 'X' || b == '^' || b == '_':
		return scanString(s, start)
	case b >= intermediateFirst && b <= intermediateLast:
		return scanEscape(s, start)
	case b >= escFinalFirst && b <= escFinalLast:
		return start + 2, scanComplete
	default:
		return 0, scanInvalid
	}
}

// scanCSI reads ESC [ followed by parameter bytes, intermediate bytes and a final byte. Colors,
// cursor movement and erasing are all CSI, as are the private forms such as the cursor hiding
// \x1b[?25l that a pattern over digits alone does not match.
func scanCSI(s string, start int) (int, scanResult) {
	i := start + 2
	for i < len(s) && s[i] >= paramFirst && s[i] <= paramLast {
		i++
	}
	for i < len(s) && s[i] >= intermediateFirst && s[i] <= intermediateLast {
		i++
	}

	if i >= len(s) {
		return 0, scanUnterminated
	}
	if s[i] < csiFinalFirst || s[i] > csiFinalLast {
		return 0, scanInvalid
	}

	return i + 1, scanComplete
}

// scanString reads a sequence whose payload runs until a terminator: OSC, which carries a window
// title or a hyperlink, and the rarer DCS, SOS, PM and APC. Left in, an unterminated one of these
// would swallow the rest of the line, which is why they are read rather than ignored.
func scanString(s string, start int) (int, scanResult) {
	for i := start + 2; i < len(s); i++ {
		switch s[i] {
		case bell:
			return i + 1, scanComplete
		case escape:
			if i+1 >= len(s) {
				return 0, scanUnterminated
			}
			if s[i+1] == '\\' {
				return i + 2, scanComplete
			}

			return 0, scanInvalid
		}
	}

	return 0, scanUnterminated
}

// scanEscape reads a short escape with no payload, such as the ESC ( B that selects a charset.
func scanEscape(s string, start int) (int, scanResult) {
	i := start + 1
	for i < len(s) && s[i] >= intermediateFirst && s[i] <= intermediateLast {
		i++
	}

	if i >= len(s) {
		return 0, scanUnterminated
	}
	if s[i] < escFinalFirst || s[i] > escFinalLast {
		return 0, scanInvalid
	}

	return i + 1, scanComplete
}
