package envutil

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

// These tests exercise the deprecated package-level wrappers to confirm
// they delegate to env.Repository correctly. The wrappers mutate real
// process env; t.Setenv restores each key when the test ends.

func TestRevokableSetenv(t *testing.T) {
	t.Setenv("KEY_RevokableSetenv", "orig")

	revoke, err := RevokableSetenv("KEY_RevokableSetenv", "new")
	require.NoError(t, err)
	require.Equal(t, "new", os.Getenv("KEY_RevokableSetenv"))

	require.NoError(t, revoke())
	require.Equal(t, "orig", os.Getenv("KEY_RevokableSetenv"))
}

func TestRevokableSetenvs(t *testing.T) {
	t.Setenv("KEY_M1", "o1")
	t.Setenv("KEY_M2", "o2")

	revoke, err := RevokableSetenvs(map[string]string{
		"KEY_M1": "n1",
		"KEY_M2": "n2",
	})
	require.NoError(t, err)
	require.Equal(t, "n1", os.Getenv("KEY_M1"))
	require.Equal(t, "n2", os.Getenv("KEY_M2"))

	require.NoError(t, revoke())
	require.Equal(t, "o1", os.Getenv("KEY_M1"))
	require.Equal(t, "o2", os.Getenv("KEY_M2"))
}

func TestSetenvForFunction(t *testing.T) {
	t.Setenv("KEY_SF", "orig")

	var inside string
	require.NoError(t, SetenvForFunction("KEY_SF", "temp", func() {
		inside = os.Getenv("KEY_SF")
	}))
	require.Equal(t, "temp", inside)
	require.Equal(t, "orig", os.Getenv("KEY_SF"))
}

func TestSetenvsForFunction(t *testing.T) {
	t.Setenv("KEY_SFM1", "o1")
	t.Setenv("KEY_SFM2", "o2")

	seen := map[string]string{}
	require.NoError(t, SetenvsForFunction(map[string]string{
		"KEY_SFM1": "t1",
		"KEY_SFM2": "t2",
	}, func() {
		seen["KEY_SFM1"] = os.Getenv("KEY_SFM1")
		seen["KEY_SFM2"] = os.Getenv("KEY_SFM2")
	}))
	require.Equal(t, "t1", seen["KEY_SFM1"])
	require.Equal(t, "t2", seen["KEY_SFM2"])
	require.Equal(t, "o1", os.Getenv("KEY_SFM1"))
	require.Equal(t, "o2", os.Getenv("KEY_SFM2"))
}

func TestStringFlagOrEnv(t *testing.T) {
	t.Setenv("FLAG_ENV", "env")

	flagVal := "flag"
	empty := ""
	require.Equal(t, "flag", StringFlagOrEnv(&flagVal, "FLAG_ENV"))
	require.Equal(t, "env", StringFlagOrEnv(&empty, "FLAG_ENV"))
	require.Equal(t, "env", StringFlagOrEnv(nil, "FLAG_ENV"))
}

func TestGetenvWithDefault(t *testing.T) {
	t.Setenv("DEF_KEY_SET", "set")
	require.Equal(t, "set", GetenvWithDefault("DEF_KEY_SET", "fallback"))
	require.Equal(t, "fallback", GetenvWithDefault("DEF_KEY_UNSET", "fallback"))
}

func TestRequiredEnv(t *testing.T) {
	t.Setenv("REQ_KEY_SET", "set")
	got, err := RequiredEnv("REQ_KEY_SET")
	require.NoError(t, err)
	require.Equal(t, "set", got)

	_, err = RequiredEnv("REQ_KEY_UNSET")
	require.EqualError(t, err, "required environment variable (REQ_KEY_UNSET) not provided")
}
