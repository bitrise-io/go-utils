package ansi

import (
	"strconv"
	"strings"
)

const (
	redForeground       = 31
	brightRedForeground = 91

	extendedForeground = 38
	extendedBackground = 48
	extendedUnderline  = 58

	extendedColorIndexed = 5
	extendedColorRGB     = 2
)

// Kind tells which escape sequence family a Sequence belongs to.
type Kind int

const (
	// KindCSI is a control sequence introducer: ESC [ parameters intermediates final. Colors, cursor
	// movement and erasing are all CSI.
	KindCSI Kind = iota + 1
	// KindOSC is an operating system command: ESC ] text BEL, or ESC ] text ESC \. Window titles and
	// hyperlinks are OSC.
	KindOSC
	// KindString is one of the other string sequences, DCS, SOS, PM and APC, which run until ESC \.
	KindString
	// KindEscape is a short escape with no string payload, such as ESC M or ESC ( B.
	KindEscape
)

// Sequence is one escape sequence: its span in the input, its family, and, for a CSI sequence, the
// final byte that names it.
type Sequence struct {
	Start, End int
	Kind       Kind
	// Final is the last byte of a CSI sequence, 'm' for a color, and 0 for every other kind.
	Final byte
	// Params holds the parameters of an SGR sequence, flattened over both separators, so that
	// 38;5;31 and 38:5:31 read the same. It is nil for every other sequence.
	Params []int
}

// IsSGR reports whether the sequence selects graphic rendition, so colors and resets.
func (s Sequence) IsSGR() bool {
	return s.Kind == KindCSI && s.Final == 'm'
}

// SetsRedForeground reports whether the sequence sets a red foreground, so parameter 31 or 91. An
// index inside a 256 color or truecolor selection (38;5;31, 38;2;r;g;b) is not a color of its own and
// does not count.
func (s Sequence) SetsRedForeground() bool {
	if !s.IsSGR() {
		return false
	}

	for i := 0; i < len(s.Params); {
		switch s.Params[i] {
		case redForeground, brightRedForeground:
			return true
		case extendedForeground, extendedBackground, extendedUnderline:
			i += extendedColorLength(s.Params[i+1:])
		default:
			i++
		}
	}

	return false
}

// parseParams reads the parameter bytes of an SGR sequence. Sub-parameters are flattened into the
// list, so that the colon form of an extended color is read like the semicolon form. A parameter that
// is empty or not a number is 0, which is what a terminal makes of it.
func parseParams(params string) []int {
	if params == "" {
		return nil
	}

	parts := strings.Split(strings.ReplaceAll(params, ":", ";"), ";")

	values := make([]int, len(parts))
	for i, part := range parts {
		values[i], _ = strconv.Atoi(part)
	}

	return values
}

// extendedColorLength reports how many parameters the color selection starting at 38, 48 or 58 spans,
// counting that parameter itself, so that its index parameters are not read as colors of their own.
func extendedColorLength(rest []int) int {
	if len(rest) == 0 {
		return 1
	}

	switch rest[0] {
	case extendedColorIndexed:
		return min(3, len(rest)+1)
	case extendedColorRGB:
		return min(5, len(rest)+1)
	default:
		return 1
	}
}
