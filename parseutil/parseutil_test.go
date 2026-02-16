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

func TestStringFrom(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected string
	}{
		// String input (fast path)
		{"string", "hello", "hello"},
		{"empty string", "", ""},
		{"string with spaces", "  hello  ", "  hello  "},

		// Integer types
		{"int", 42, "42"},
		{"int zero", 0, "0"},
		{"int negative", -123, "-123"},
		{"int8", int8(127), "127"},
		{"int16", int16(32767), "32767"},
		{"int32", int32(2147483647), "2147483647"},
		{"int64", int64(9223372036854775807), "9223372036854775807"},

		// Unsigned integer types
		{"uint", uint(42), "42"},
		{"uint8", uint8(255), "255"},
		{"uint16", uint16(65535), "65535"},
		{"uint32", uint32(4294967295), "4294967295"},
		{"uint64", uint64(18446744073709551615), "18446744073709551615"},

		// Float types
		{"float32", float32(3.14), "3.14"},
		{"float64", 3.14159, "3.14159"},
		{"float64 zero", 0.0, "0"},
		{"float64 negative", -2.5, "-2.5"},

		// Boolean types
		{"bool true", true, "true"},
		{"bool false", false, "false"},

		// Nil
		{"nil", nil, "<nil>"},

		// Pointer
		{"string pointer", stringPtr("test"), "test"},

		// Struct
		{"struct", struct{ Name string }{Name: "test"}, "{test}"},

		// Slice
		{"slice", []int{1, 2, 3}, "[1 2 3]"},
		{"empty slice", []int{}, "[]"},

		// Map
		{"map", map[string]int{"a": 1}, "map[a:1]"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := StringFrom(tt.input)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestStringFrom_ExplicitInterface(t *testing.T) {

	t.Run("interface containing map", func(t *testing.T) {
		var v interface{} = map[string]string{"key": "value"}
		result := StringFrom(v)
		require.Equal(t, "map[key:value]", result)
	})
}

// Helper function for tests
func stringPtr(s string) *string {
	return &s
}
