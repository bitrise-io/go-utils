package env

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// Tests mutate real process env. t.Setenv restores on test end.

func TestRepository_GetOrDefault(t *testing.T) {
	repo := NewRepository()

	t.Setenv("KEY_A", "value")
	require.Equal(t, "value", repo.GetOrDefault("KEY_A", "fallback"))

	t.Setenv("KEY_B", "")
	require.Equal(t, "fallback", repo.GetOrDefault("KEY_B", "fallback"))

	require.Equal(t, "fallback", repo.GetOrDefault("KEY_NOT_SET", "fallback"))
}

func TestRepository_Required(t *testing.T) {
	repo := NewRepository()

	t.Setenv("REQ_A", "value")
	got, err := repo.Required("REQ_A")
	require.NoError(t, err)
	require.Equal(t, "value", got)

	t.Setenv("REQ_B", "")
	_, err = repo.Required("REQ_B")
	require.EqualError(t, err, "required environment variable (REQ_B) not provided")

	_, err = repo.Required("REQ_NOT_SET")
	require.EqualError(t, err, "required environment variable (REQ_NOT_SET) not provided")
}

func TestRepository_FlagOrEnv(t *testing.T) {
	repo := NewRepository()
	t.Setenv("FLAG_KEY", "env-value")

	flagVal := "flag-value"
	empty := ""

	require.Equal(t, "flag-value", repo.FlagOrEnv(&flagVal, "FLAG_KEY"))
	require.Equal(t, "env-value", repo.FlagOrEnv(&empty, "FLAG_KEY"))
	require.Equal(t, "env-value", repo.FlagOrEnv(nil, "FLAG_KEY"))
}

func TestRepository_Revokable(t *testing.T) {
	repo := NewRepository()
	t.Setenv("REV_A", "orig")

	revoke, err := repo.Revokable("REV_A", "new")
	require.NoError(t, err)
	require.Equal(t, "new", repo.Get("REV_A"))

	require.NoError(t, revoke())
	require.Equal(t, "orig", repo.Get("REV_A"))
}

func TestRepository_Revokable_previouslyUnset(t *testing.T) {
	repo := NewRepository()
	// Not calling t.Setenv: the key is unset at test start.
	const key = "REV_NEW_ONLY"

	revoke, err := repo.Revokable(key, "tmp")
	require.NoError(t, err)
	require.Equal(t, "tmp", repo.Get(key))

	require.NoError(t, revoke())
	require.Equal(t, "", repo.Get(key))
}

func TestRepository_RevokableMany(t *testing.T) {
	repo := NewRepository()
	t.Setenv("REV_M_A", "origA")
	t.Setenv("REV_M_B", "origB")

	revoke, err := repo.RevokableMany(map[string]string{
		"REV_M_A": "newA",
		"REV_M_B": "newB",
	})
	require.NoError(t, err)
	require.Equal(t, "newA", repo.Get("REV_M_A"))
	require.Equal(t, "newB", repo.Get("REV_M_B"))

	require.NoError(t, revoke())
	require.Equal(t, "origA", repo.Get("REV_M_A"))
	require.Equal(t, "origB", repo.Get("REV_M_B"))
}

func TestRepository_Scoped(t *testing.T) {
	repo := NewRepository()
	t.Setenv("SC_A", "orig")

	var seenInside string
	err := repo.Scoped("SC_A", "temp", func() {
		seenInside = repo.Get("SC_A")
	})
	require.NoError(t, err)
	require.Equal(t, "temp", seenInside)
	require.Equal(t, "orig", repo.Get("SC_A"))
}

func TestRepository_ScopedMany(t *testing.T) {
	repo := NewRepository()
	t.Setenv("SC_M_A", "origA")
	t.Setenv("SC_M_B", "origB")

	seen := map[string]string{}
	err := repo.ScopedMany(map[string]string{
		"SC_M_A": "tempA",
		"SC_M_B": "tempB",
	}, func() {
		seen["SC_M_A"] = repo.Get("SC_M_A")
		seen["SC_M_B"] = repo.Get("SC_M_B")
	})
	require.NoError(t, err)
	require.Equal(t, "tempA", seen["SC_M_A"])
	require.Equal(t, "tempB", seen["SC_M_B"])
	require.Equal(t, "origA", repo.Get("SC_M_A"))
	require.Equal(t, "origB", repo.Get("SC_M_B"))
}

func TestRepository_Scoped_restoresOnPanic(t *testing.T) {
	repo := NewRepository()
	t.Setenv("SC_P", "orig")

	defer func() {
		recover()
		require.Equal(t, "orig", repo.Get("SC_P"))
	}()

	_ = repo.Scoped("SC_P", "temp", func() {
		panic("boom")
	})
}
