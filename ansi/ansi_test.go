package ansi

import (
	"math/rand"
	"testing"

	"github.com/bitrise-io/go-utils/v2/log/colorstring"
	"github.com/stretchr/testify/require"
)

func TestStrip(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{"plain text is untouched", "Execution failed for task ':app:compileDebugKotlin'.", "Execution failed for task ':app:compileDebugKotlin'."},
		{"tabs and newlines survive", "a\tb\nc\n", "a\tb\nc\n"},
		{"non-ASCII text survives", "árvíztűrő tükörfúrógép", "árvíztűrő tükörfúrógép"},
		{"color", "\x1b[31;1mCompilation error\x1b[0m", "Compilation error"},
		{"private parameter", "\x1b[?25lCompilation error", "Compilation error"},
		{"intermediate byte", "\x1b[6 qCompilation error", "Compilation error"},
		{"cursor movement", "\x1b[2K\x1b[1Gprogress", "progress"},
		{"OSC ended by BEL", "\x1b]0;window title\x07text", "text"},
		{"OSC ended by ST", "\x1b]8;;https://example.com\x1b\\link", "link"},
		{"DCS", "\x1bPq payload \x1b\\text", "text"},
		{"two byte escape", "\x1b(Btext", "text"},
		{"one byte escape", "\x1bMtext", "text"},
		{"lone escape at the end", "text\x1b", "text"},
		{"unterminated CSI swallows its own tail only", "text\x1b[31", "text"},
		{"stray escape before a control byte", "a\x1b\x01b", "ab"},
		{"other control bytes", "a\x00b\x07c\x7fd", "abcd"},
		{"CRLF keeps its text", "line\r\nnext\r\n", "line\nnext\n"},
		{"redraw keeps what is on screen", "10%\r50%\r100% done", "100% done"},
		{"redraw per line", "10%\r100%\nnext\r\n", "100%\nnext\n"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, Strip(tt.in))
		})
	}
}

// TestStripKeepsGradleErrorLines runs the shapes a Gradle run puts a failure line in, with the
// matching the step does, since an escape left at the start of a line stops an anchored pattern from
// matching at all.
func TestStripKeepsGradleErrorLines(t *testing.T) {
	const want = "e: file:///src/A.kt:1:1 Unresolved reference 'a'."

	tests := map[string]string{
		"colored":                      "\x1b[31;1m" + want + "\x1b[0m",
		"cursor hidden around it":      "\x1b[?25l\x1b[31m" + want + "\x1b[0m\x1b[?25h",
		"redrawn over a progress line": "<-------> 0% EXECUTING [2s]\r" + want,
		"behind a hyperlink":           "\x1b]8;;file:///src/A.kt\x1b\\" + want,
	}

	for name, in := range tests {
		t.Run(name, func(t *testing.T) {
			require.Equal(t, want, Strip(in))
		})
	}
}

func TestStripKeepsWhatColorstringWrites(t *testing.T) {
	const message = "build task failed"

	for name, colored := range map[string]string{
		"red":    colorstring.Red(message),
		"yellow": colorstring.Yellow(message),
		"nested": colorstring.Red(colorstring.Green(message)),
	} {
		t.Run(name, func(t *testing.T) {
			require.Equal(t, message, Strip(colored))
		})
	}
}

func TestStripLeavesNoControlBytes(t *testing.T) {
	random := rand.New(rand.NewSource(1))

	for i := 0; i < 2000; i++ {
		in := randomControlHeavyString(random)

		for _, b := range []byte(Strip(in)) {
			switch {
			case b == '\n' || b == '\t':
			case b < 0x20 || b == deleteByte:
				require.Failf(t, "control byte survived", "%q in %q", b, in)
			}
		}
	}
}

// randomControlHeavyString builds input out of the bytes that drive the scanner, so that half written
// sequences and stray control bytes turn up far more often than they would in random text.
func randomControlHeavyString(random *rand.Rand) string {
	alphabet := []byte("\x1b[]()P;0123456789mqABZ \t\n\r\x00\x07\x7fabc")

	b := make([]byte, random.Intn(40))
	for i := range b {
		b[i] = alphabet[random.Intn(len(alphabet))]
	}

	return string(b)
}
