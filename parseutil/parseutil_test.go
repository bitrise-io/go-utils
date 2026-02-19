package parseutil

import (
	"strings"
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

	t.Run("pointer prints address format", func(t *testing.T) {
		s := "test"
		result := StringFrom(&s)
		require.True(t, strings.HasPrefix(result, "0x"),
			"pointer should print address starting with 0x, got: %s", result)
	})
}

func TestStringPtrFrom(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected string
	}{
		{"string", "hello", "hello"},
		{"int", 42, "42"},
		{"bool", true, "true"},
		{"float", 3.14, "3.14"},
		{"nil", nil, "<nil>"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := StringPtrFrom(tt.input)
			require.NotNil(t, result)
			require.Equal(t, tt.expected, *result)
		})
	}
}

func TestBoolFrom(t *testing.T) {
	tests := []struct {
		name       string
		input      interface{}
		expectedOk bool
		expected   bool
	}{
		// Direct bool values
		{"bool true", true, true, true},
		{"bool false", false, true, false},

		// String values that can be parsed
		{"string yes", "yes", true, true},
		{"string no", "no", true, false},
		{"string y", "y", true, true},
		{"string n", "n", true, false},
		{"string true", "true", true, true},
		{"string false", "false", true, false},
		{"string 1", "1", true, true},
		{"string 0", "0", true, false},
		{"string t", "t", true, true},
		{"string f", "f", true, false},

		// Integer values converted to strings
		{"int 1", 1, true, true},
		{"int 0", 0, true, false},

		// Values that cannot be parsed
		{"string invalid", "invalid", false, false},
		{"int 2", 2, false, false},
		{"int 42", 42, false, false},
		{"float 3.14", 3.14, false, false},
		{"nil", nil, false, false},
		{"empty string", "", false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, ok := BoolFrom(tt.input)
			require.Equal(t, tt.expectedOk, ok)
			if tt.expectedOk {
				require.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestBoolPtrFrom(t *testing.T) {
	t.Run("successful conversions return pointer", func(t *testing.T) {
		tests := []struct {
			name     string
			input    interface{}
			expected bool
		}{
			{"bool true", true, true},
			{"bool false", false, false},
			{"string yes", "yes", true},
			{"string no", "no", false},
			{"int 1", 1, true},
			{"int 0", 0, false},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result, ok := BoolPtrFrom(tt.input)
				require.True(t, ok)
				require.NotNil(t, result)
				require.Equal(t, tt.expected, *result)
			})
		}
	})

	t.Run("failed conversions return nil", func(t *testing.T) {
		tests := []struct {
			name  string
			input interface{}
		}{
			{"invalid string", "invalid"},
			{"int 42", 42},
			{"float", 3.14},
			{"nil", nil},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result, ok := BoolPtrFrom(tt.input)
				require.False(t, ok)
				require.Nil(t, result)
			})
		}
	})
}
