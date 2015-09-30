package builtinutil

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCastInterfaceToInterfaceSlice(t *testing.T) {
	// string
	casted, err := CastInterfaceToInterfaceSlice([]string{"a", "b", "c"})
	require.NoError(t, err)
	require.Equal(t, "a", casted[0])
	require.Equal(t, "b", casted[1])
	require.Equal(t, "c", casted[2])

	// int
	casted, err = CastInterfaceToInterfaceSlice([]int{3, 2, 1})
	require.NoError(t, err)
	require.Equal(t, 3, casted[0])
	require.Equal(t, 2, casted[1])
	require.Equal(t, 1, casted[2])

	// empty
	casted, err = CastInterfaceToInterfaceSlice([]string{})
	require.NoError(t, err)
	require.Equal(t, []interface{}{}, casted)
}

func TestDeepEqualSlices(t *testing.T) {
	s1 := []string{"a", "b"}
	interfaceS1, _ := CastInterfaceToInterfaceSlice(s1)
	s2 := []string{"b", "a"}
	interfaceS2, _ := CastInterfaceToInterfaceSlice(s2)
	require.True(t, DeepEqualSlices(interfaceS1, interfaceS2))

	s3 := []string{"b", "a", "c"}
	interfaceS3, _ := CastInterfaceToInterfaceSlice(s3)
	require.False(t, DeepEqualSlices(interfaceS1, interfaceS3))

	// empty
	require.True(t, DeepEqualSlices([]interface{}{}, []interface{}{}))
}
