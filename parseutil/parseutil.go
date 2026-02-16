package parseutil

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// ParseBool parses a string representation of a boolean value.
//
// It accepts the following inputs (case-insensitive, whitespace trimmed):
//   - Custom values: "yes", "y" (true), "no", "n" (false)
//   - Standard values: "true", "t", "1" (true), "false", "f", "0" (false)
//
// Returns an error if the input is empty or cannot be parsed as a boolean.
func ParseBool(input string) (bool, error) {
	// Validation
	if input == "" {
		return false, errors.New("no string to parse")
	}

	// Normalization
	normalized := strings.ToLower(strings.TrimSpace(input))

	// Custom parsing
	switch normalized {
	case "yes", "y":
		return true, nil
	case "no", "n":
		return false, nil
	}

	// Delegate to stdlib
	return strconv.ParseBool(normalized)
}

// StringFrom converts any value to its string representation.
//
// If the value is already a string, it returns it directly for efficiency.
// Otherwise, it uses fmt.Sprintf with the %v verb to convert the value.
//
// This function always returns a string and never returns an error.
//
// Example:
//
//	StringFrom("hello")    // "hello"
//	StringFrom(42)         // "42"
//	StringFrom(true)       // "true"
//	StringFrom(3.14)       // "3.14"
func StringFrom(value interface{}) string {
	if str, ok := value.(string); ok {
		return str
	}
	return fmt.Sprintf("%v", value)
}
