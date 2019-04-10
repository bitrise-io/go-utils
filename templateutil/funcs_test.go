package templateutil

import (
	"testing"

	"github.com/bitrise-io/go-utils/envutil"
	"github.com/stretchr/testify/require"
)

func TestGetenv(t *testing.T) {
	testEnvKey := "KEY_TestGetenv"

	// no env set yet, empty value should be returned
	require.Equal(t, "", Getenv(testEnvKey))

	// set the env
	revokeFn, err := envutil.RevokableSetenv(testEnvKey, "env value")
	require.NoError(t, err)
	defer func() {
		require.NoError(t, revokeFn())
	}()

	// env set - value should be the env's value
	require.Equal(t, "env value", Getenv(testEnvKey))
}

func TestGetenvRequired(t *testing.T) {
	testEnvKey := "KEY_TestGetenvRequired"

	t.Log("no env set yet, empty value should be returned")
	{
		val, err := GetenvRequired(testEnvKey)
		require.Equal(t, "", val)
		require.Error(t, err, "no environment variable value found for key: KEY_TestGetenvRequired")
	}

	// set the env
	revokeFn, err := envutil.RevokableSetenv(testEnvKey, "env value")
	require.NoError(t, err)
	defer func() {
		require.NoError(t, revokeFn())
	}()

	t.Log("env set - value should be the env's value")
	{
		val, err := GetenvRequired(testEnvKey)
		require.NoError(t, err)
		require.Equal(t, "env value", val)
	}
}
