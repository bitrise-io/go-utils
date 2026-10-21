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

// scanResult tells a caller what to do with the bytes at an ESC: take the sequence, carry the rest of
// the input over to the next chunk, or treat the ESC as a stray byte and move on one byte.
type scanResult int

const (
	scanComplete scanResult = iota
	scanUnterminated
	scanInvalid
)

// scanAt reads the escape sequence starting at the ESC byte s[start].
func scanAt(s string, start int) (Sequence, scanResult) {
	if start+1 >= len(s) {
		return Sequence{}, scanUnterminated
	}

	switch b := s[start+1]; {
	case b == '[':
		return scanCSI(s, start)
	case b == ']':
		return scanString(s, start, KindOSC)
	case b == 'P' || b == 'X' || b == '^' || b == '_':
		return scanString(s, start, KindString)
	case b >= intermediateFirst && b <= intermediateLast:
		return scanEscape(s, start)
	case b >= escFinalFirst && b <= escFinalLast:
		return Sequence{Start: start, End: start + 2, Kind: KindEscape}, scanComplete
	default:
		return Sequence{}, scanInvalid
	}
}

func scanCSI(s string, start int) (Sequence, scanResult) {
	i := start + 2
	for i < len(s) && s[i] >= paramFirst && s[i] <= paramLast {
		i++
	}
	params := s[start+2 : i]

	for i < len(s) && s[i] >= intermediateFirst && s[i] <= intermediateLast {
		i++
	}

	if i >= len(s) {
		return Sequence{}, scanUnterminated
	}
	if s[i] < csiFinalFirst || s[i] > csiFinalLast {
		return Sequence{}, scanInvalid
	}

	seq := Sequence{Start: start, End: i + 1, Kind: KindCSI, Final: s[i]}
	if seq.IsSGR() {
		seq.Params = parseParams(params)
	}

	return seq, scanComplete
}

// scanString reads a sequence whose payload runs until a terminator. BEL ends an OSC, and ESC \ ends
// any of them; BEL is taken for the others too, because a terminal that accepts it is what a writer
// of one is likely to have tested against.
func scanString(s string, start int, kind Kind) (Sequence, scanResult) {
	for i := start + 2; i < len(s); i++ {
		switch s[i] {
		case bell:
			return Sequence{Start: start, End: i + 1, Kind: kind}, scanComplete
		case escape:
			if i+1 >= len(s) {
				return Sequence{}, scanUnterminated
			}
			if s[i+1] == '\\' {
				return Sequence{Start: start, End: i + 2, Kind: kind}, scanComplete
			}

			return Sequence{}, scanInvalid
		}
	}

	return Sequence{}, scanUnterminated
}

func scanEscape(s string, start int) (Sequence, scanResult) {
	i := start + 1
	for i < len(s) && s[i] >= intermediateFirst && s[i] <= intermediateLast {
		i++
	}

	if i >= len(s) {
		return Sequence{}, scanUnterminated
	}
	if s[i] < escFinalFirst || s[i] > escFinalLast {
		return Sequence{}, scanInvalid
	}

	return Sequence{Start: start, End: i + 1, Kind: KindEscape}, scanComplete
}
