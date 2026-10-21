package ansi

import (
	"math/rand"
	"strings"
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

func TestStripKeepsWhatColorstringWrites(t *testing.T) {
	const message = "build task failed"

	for name, colored := range map[string]string{
		"red":    colorstring.Red(message),
		"yellow": colorstring.Yellow(message),
		"bold":   colorstring.Red(colorstring.Green(message)),
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

func TestSequences(t *testing.T) {
	tests := []struct {
		name        string
		in          string
		want        []Sequence
		wantPending int
	}{
		{
			name:        "no sequences",
			in:          "plain text",
			wantPending: -1,
		},
		{
			name:        "SGR parameters are parsed",
			in:          "\x1b[31;1mred\x1b[0m",
			want:        []Sequence{{Start: 0, End: 7, Kind: KindCSI, Final: 'm', Params: []int{31, 1}}, {Start: 10, End: 14, Kind: KindCSI, Final: 'm', Params: []int{0}}},
			wantPending: -1,
		},
		{
			name:        "a private parameter does not end at the first m in the text",
			in:          "\x1b[?25lCompilation error",
			want:        []Sequence{{Start: 0, End: 6, Kind: KindCSI, Final: 'l'}},
			wantPending: -1,
		},
		{
			name:        "OSC",
			in:          "\x1b]0;title\x07",
			want:        []Sequence{{Start: 0, End: 10, Kind: KindOSC}},
			wantPending: -1,
		},
		{
			name:        "unterminated CSI is pending",
			in:          "text\x1b[31",
			wantPending: 4,
		},
		{
			name:        "unterminated OSC is pending",
			in:          "\x1b]0;title",
			wantPending: 0,
		},
		{
			name:        "a stray escape does not hide what follows it",
			in:          "\x1b\x01\x1b[0m",
			want:        []Sequence{{Start: 2, End: 6, Kind: KindCSI, Final: 'm', Params: []int{0}}},
			wantPending: -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			seqs, pending := Sequences(tt.in)
			require.Equal(t, tt.want, seqs)
			require.Equal(t, tt.wantPending, pending)
		})
	}
}

func TestSequencesPendingCarriesOverExactly(t *testing.T) {
	const (
		first  = "error: \x1b[38;5;"
		second = "31mteal\x1b[0m"
	)

	seqs, pending := Sequences(first)
	require.Empty(t, seqs)
	require.Equal(t, strings.Index(first, "\x1b"), pending)

	seqs, pending = Sequences(first[pending:] + second)
	require.Equal(t, -1, pending)
	require.Len(t, seqs, 2)
	require.Equal(t, []int{38, 5, 31}, seqs[0].Params)
}
