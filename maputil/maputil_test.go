package maputil

import (
	"testing"

	"github.com/bitrise-io/go-utils/testutil"
	"github.com/stretchr/testify/require"
)

func TestKeysOfStringStringMap(t *testing.T) {
	t.Log("Empty map")
	keys := KeysOfStringStringMap(map[string]string{})
	require.Equal(t, []string{}, keys)
	require.Equal(t, 0, len(keys))

	t.Log("Single key")
	keys = KeysOfStringStringMap(map[string]string{"a": "value"})
	require.Equal(t, 1, len(keys))
	require.Equal(t, []string{"a"}, keys)

	t.Log("Multiple keys")
	keys = KeysOfStringStringMap(map[string]string{"a": "value 1", "b": "value 2"})
	require.Equal(t, 2, len(keys))
	testutil.EqualSlicesWithoutOrder(t, []string{"a", "b"}, keys)
}
