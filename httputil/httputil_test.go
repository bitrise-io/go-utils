package httputil

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetSingleValueFromHeader(t *testing.T) {
	t.Log("Simple OK")
	{
		val, err := GetSingleValueFromHeader("Key", http.Header{"Key": {"value1"}})
		require.NoError(t, err)
		require.Equal(t, "value1", val)
	}

	t.Log("Not found - empty header")
	{
		_, err := GetSingleValueFromHeader("No-Key", http.Header{})
		require.EqualError(t, err, "No value found in HEADER for the key: No-Key")
	}

	t.Log("Not found - different key")
	{
		_, err := GetSingleValueFromHeader("No-Key", http.Header{"Key": {"value1"}})
		require.EqualError(t, err, "No value found in HEADER for the key: No-Key")
	}

	t.Log("Multiple values")
	{
		_, err := GetSingleValueFromHeader("Key", http.Header{"Key": {"value1", "value2"}})
		require.EqualError(t, err, "Multiple values found in HEADER for the key: Key")
	}
}
