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

func TestCastToStringPtr(t *testing.T) {
	require.Equal(t, "1", *CastToStringPtr(1))
	require.Equal(t, "0.1", *CastToStringPtr(0.1))
	require.Equal(t, "true", *CastToStringPtr(true))
	require.Equal(t, "false", *CastToStringPtr(false))
	require.Equal(t, "yes", *CastToStringPtr("yes"))
}

func TestCastToBoolPtr(t *testing.T) {
	casted, ok := CastToBoolPtr(true)
	require.Equal(t, true, ok)
	require.Equal(t, true, *casted)

	casted, ok = CastToBoolPtr("true")
	require.Equal(t, true, ok)
	require.Equal(t, true, *casted)

	casted, ok = CastToBoolPtr(false)
	require.Equal(t, true, ok)
	require.Equal(t, false, *casted)

	casted, ok = CastToBoolPtr("false")
	require.Equal(t, true, ok)
	require.Equal(t, false, *casted)

	casted, ok = CastToBoolPtr("yes")
	require.Equal(t, true, ok)
	require.Equal(t, true, *casted)

	casted, ok = CastToBoolPtr("no")
	require.Equal(t, true, ok)
	require.Equal(t, false, *casted)

	casted, ok = CastToBoolPtr(1)
	require.Equal(t, true, ok)
	require.Equal(t, true, *casted)

	casted, ok = CastToBoolPtr("1")
	require.Equal(t, true, ok)
	require.Equal(t, true, *casted)

	casted, ok = CastToBoolPtr(0)
	require.Equal(t, true, ok)
	require.Equal(t, false, *casted)

	casted, ok = CastToBoolPtr("0")
	require.Equal(t, true, ok)
	require.Equal(t, false, *casted)

	casted, ok = CastToBoolPtr(0.1)
	require.Equal(t, false, ok)
	require.Equal(t, (*bool)(nil), casted)

	casted, ok = CastToBoolPtr("0.1")
	require.Equal(t, false, ok)
	require.Equal(t, (*bool)(nil), casted)

	casted, ok = CastToBoolPtr("test")
	require.Equal(t, false, ok)
	require.Equal(t, (*bool)(nil), casted)
}
