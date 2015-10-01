package stringutil

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenericTrim(t *testing.T) {
	require.Equal(t, "", genericTrim("", 4, false, false))
	require.Equal(t, "", genericTrim("", 4, false, true))
	require.Equal(t, "", genericTrim("", 4, true, false))
	require.Equal(t, "", genericTrim("", 4, true, true))

	require.Equal(t, "1234", genericTrim("123456789", 4, false, false))
	require.Equal(t, "1...", genericTrim("123456789", 4, false, true))
	require.Equal(t, "6789", genericTrim("123456789", 4, true, false))
	require.Equal(t, "...9", genericTrim("123456789", 4, true, true))
}

func TestMaxLastChars(t *testing.T) {
	require.Equal(t, "", MaxLastChars("", 10))
	require.Equal(t, "a", MaxLastChars("a", 1))
	require.Equal(t, "a", MaxLastChars("ba", 1))
	require.Equal(t, "ba", MaxLastChars("ba", 10))
	require.Equal(t, "a", MaxLastChars("cba", 1))
	require.Equal(t, "cba", MaxLastChars("cba", 10))

	require.Equal(t, "llo world!", MaxLastChars("hello world!", 10))
}

func TestMaxLastCharsWithDots(t *testing.T) {
	require.Equal(t, "", MaxLastCharsWithDots("", 10))
	require.Equal(t, "", MaxLastCharsWithDots("1234", 1))
	require.Equal(t, "...56", MaxLastCharsWithDots("123456", 5))
	require.Equal(t, "123456", MaxFirstCharsWithDots("123456", 6))
	require.Equal(t, "123456", MaxLastCharsWithDots("123456", 10))

	require.Equal(t, "... world!", MaxLastCharsWithDots("hello world!", 10))
}

func TestMaxFirstChars(t *testing.T) {
	require.Equal(t, "", MaxFirstChars("", 10))
	require.Equal(t, "a", MaxFirstChars("a", 1))
	require.Equal(t, "b", MaxFirstChars("ba", 1))
	require.Equal(t, "ba", MaxFirstChars("ba", 10))
	require.Equal(t, "c", MaxFirstChars("cba", 1))
	require.Equal(t, "cba", MaxFirstChars("cba", 10))

	require.Equal(t, "hello worl", MaxFirstChars("hello world!", 10))
}

func TestMaxFirstCharsWithDots(t *testing.T) {
	require.Equal(t, "", MaxFirstCharsWithDots("", 10))
	require.Equal(t, "", MaxFirstCharsWithDots("1234", 1))
	require.Equal(t, "12...", MaxFirstCharsWithDots("123456", 5))
	require.Equal(t, "123456", MaxFirstCharsWithDots("123456", 6))
	require.Equal(t, "123456", MaxFirstCharsWithDots("123456", 10))

	require.Equal(t, "hello w...", MaxFirstCharsWithDots("hello world!", 10))
}
