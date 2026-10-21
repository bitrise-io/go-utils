package ansi

import (
	"testing"

	"github.com/bitrise-io/go-utils/v2/log/colorstring"
	"github.com/stretchr/testify/require"
)

func TestSequenceIsSGR(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want bool
	}{
		{"color", "\x1b[31;1m", true},
		{"reset", "\x1b[0m", true},
		{"no parameters", "\x1b[m", true},
		{"cursor hiding", "\x1b[?25l", false},
		{"erase line", "\x1b[2K", false},
		{"window title", "\x1b]0;title\x07", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, onlySequence(t, tt.in).IsSGR())
		})
	}
}

func TestSequenceSetsRedForeground(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want bool
	}{
		{"red", "\x1b[31m", true},
		{"bold red", "\x1b[31;1m", true},
		{"bright red", "\x1b[91m", true},
		{"red among other parameters", "\x1b[1;4;31;7m", true},
		{"red background", "\x1b[41m", false},
		{"reset", "\x1b[0m", false},
		{"green", "\x1b[32m", false},
		{"a 256 color whose index is 31", "\x1b[38;5;31m", false},
		{"a 256 color whose index is 31, colon form", "\x1b[38:5:31m", false},
		{"a truecolor selection containing 31", "\x1b[38;2;31;31;31m", false},
		{"parameter 131", "\x1b[131m", false},
		{"red after an extended color", "\x1b[38;5;1;31m", true},
		{"not a color at all", "\x1b[2K", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, onlySequence(t, tt.in).SetsRedForeground())
		})
	}
}

func TestSequenceReadsTheColorTheLoggerWrites(t *testing.T) {
	seqs, pending := Sequences(colorstring.Red("build task failed"))
	require.Equal(t, -1, pending)
	require.Len(t, seqs, 2)

	require.True(t, seqs[0].SetsRedForeground())
	require.True(t, seqs[1].IsSGR())
	require.False(t, seqs[1].SetsRedForeground())
}

func onlySequence(t *testing.T, s string) Sequence {
	t.Helper()

	seqs, pending := Sequences(s)
	require.Equal(t, -1, pending)
	require.Len(t, seqs, 1)

	return seqs[0]
}
