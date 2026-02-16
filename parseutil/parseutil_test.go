package parseutil

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseBool(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		expected  bool
		expectErr bool
	}{
		// Validation - empty string
		{"empty string", "", false, true},

		// Custom parser - yes/y (case-insensitive)
		{"yes lowercase", "yes", true, false},
		{"yes uppercase", "YES", true, false},
		{"yes mixed case", "YeS", true, false},
		{"yes with leading space", "  yes", true, false},
		{"yes with trailing space", "yes  ", true, false},
		{"yes with surrounding spaces", "  yes  ", true, false},
		{"y lowercase", "y", true, false},
		{"y uppercase", "Y", true, false},

		// Custom parser - no/n (case-insensitive)
		{"no lowercase", "no", false, false},
		{"no uppercase", "NO", false, false},
		{"no mixed case", "No", false, false},
		{"no with spaces", "  no  ", false, false},
		{"n lowercase", "n", false, false},
		{"n uppercase", "N", false, false},

		// Stdlib parser - true
		{"true lowercase", "true", true, false},
		{"true uppercase", "TRUE", true, false},
		{"true mixed case", "True", true, false},
		{"true with spaces", "  true  ", true, false},
		{"t lowercase", "t", true, false},
		{"t uppercase", "T", true, false},
		{"1 string", "1", true, false},

		// Stdlib parser - false
		{"false lowercase", "false", false, false},
		{"false uppercase", "FALSE", false, false},
		{"false mixed case", "False", false, false},
		{"false with spaces", "  false  ", false, false},
		{"f lowercase", "f", false, false},
		{"f uppercase", "F", false, false},
		{"0 string", "0", false, false},

		// Error cases
		{"invalid word", "maybe", false, true},
		{"invalid number", "2", false, true},
		{"invalid number", "123", false, true},
		{"random string", "random", false, true},
		{"whitespace only", "   ", false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseBool(tt.input)

			if tt.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expected, result)
			}
		})
	}
}
