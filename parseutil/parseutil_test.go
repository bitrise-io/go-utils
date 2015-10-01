package parseutil

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseBool(t *testing.T) {
	testUserInput := "y"
	isYes, err := ParseBool("YeS")
	require.Equal(t, nil, err)
	require.Equal(t, true, isYes)

	testUserInput = "no"
	isYes, err = ParseBool("n")
	require.Equal(t, nil, err)
	require.Equal(t, false, isYes)

	testUserInput = `
 yes
`
	isYes, err = ParseBool(testUserInput)
	require.Equal(t, nil, err)
	require.Equal(t, true, isYes)
}

func TestCastToString(t *testing.T) {
	require.Equal(t, "1", CastToString(1))
	require.Equal(t, "1.1", CastToString(1.1))
	require.Equal(t, "true", CastToString(true))
	require.Equal(t, "false", CastToString("false"))
}
