package env

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// Tests mutate real process env. t.Setenv restores on test end.

func TestGetOrDefault(t *testing.T) {
	repo := NewRepository()

	t.Setenv("KEY_A", "value")
	require.Equal(t, "value", GetOrDefault(repo, "KEY_A", "fallback"))

	t.Setenv("KEY_B", "")
	require.Equal(t, "fallback", GetOrDefault(repo, "KEY_B", "fallback"))

	require.Equal(t, "fallback", GetOrDefault(repo, "KEY_NOT_SET", "fallback"))
}

func TestRequired(t *testing.T) {
	repo := NewRepository()

	t.Setenv("REQ_A", "value")
	got, err := Required(repo, "REQ_A")
	require.NoError(t, err)
	require.Equal(t, "value", got)

	t.Setenv("REQ_B", "")
	_, err = Required(repo, "REQ_B")
	require.EqualError(t, err, "required environment variable (REQ_B) not provided")

	_, err = Required(repo, "REQ_NOT_SET")
	require.EqualError(t, err, "required environment variable (REQ_NOT_SET) not provided")
}

func TestFlagOrEnv(t *testing.T) {
	repo := NewRepository()
	t.Setenv("FLAG_KEY", "env-value")

	flagVal := "flag-value"
	empty := ""

	require.Equal(t, "flag-value", FlagOrEnv(repo, &flagVal, "FLAG_KEY"))
	require.Equal(t, "env-value", FlagOrEnv(repo, &empty, "FLAG_KEY"))
	require.Equal(t, "env-value", FlagOrEnv(repo, nil, "FLAG_KEY"))
}

func TestRevokable(t *testing.T) {
	repo := NewRepository()
	t.Setenv("REV_A", "orig")

	revoke, err := Revokable(repo, "REV_A", "new")
	require.NoError(t, err)
	require.Equal(t, "new", repo.Get("REV_A"))

	require.NoError(t, revoke())
	require.Equal(t, "orig", repo.Get("REV_A"))
}

func TestRevokable_previouslyUnset(t *testing.T) {
	repo := NewRepository()
	// Not calling t.Setenv: the key is unset at test start.
	const key = "REV_NEW_ONLY"

	revoke, err := Revokable(repo, key, "tmp")
	require.NoError(t, err)
	require.Equal(t, "tmp", repo.Get(key))

	require.NoError(t, revoke())
	require.Equal(t, "", repo.Get(key))
}

func TestRevokableMany(t *testing.T) {
	repo := NewRepository()
	t.Setenv("REV_M_A", "origA")
	t.Setenv("REV_M_B", "origB")

	revoke, err := RevokableMany(repo, map[string]string{
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

func TestRevokableMany_atomicOnError(t *testing.T) {
	repo := NewRepository()
	t.Setenv("REV_M_ATOMIC_A", "origA")
	t.Setenv("REV_M_ATOMIC_B", "origB")

	// Empty key forces os.Setenv to fail ("setenv: invalid argument").
	// Whatever the map-iteration order, every valid key must end at its
	// original value and the returned revoke must be a no-op.
	revoke, err := RevokableMany(repo, map[string]string{
		"REV_M_ATOMIC_A": "newA",
		"REV_M_ATOMIC_B": "newB",
		"":               "boom",
	})
	require.Error(t, err)
	require.Equal(t, "origA", repo.Get("REV_M_ATOMIC_A"))
	require.Equal(t, "origB", repo.Get("REV_M_ATOMIC_B"))

	require.NoError(t, revoke())
	require.Equal(t, "origA", repo.Get("REV_M_ATOMIC_A"))
	require.Equal(t, "origB", repo.Get("REV_M_ATOMIC_B"))
}

func TestScoped(t *testing.T) {
	repo := NewRepository()
	t.Setenv("SC_A", "orig")

	var seenInside string
	err := Scoped(repo, "SC_A", "temp", func() {
		seenInside = repo.Get("SC_A")
	})
	require.NoError(t, err)
	require.Equal(t, "temp", seenInside)
	require.Equal(t, "orig", repo.Get("SC_A"))
}

func TestScopedMany(t *testing.T) {
	repo := NewRepository()
	t.Setenv("SC_M_A", "origA")
	t.Setenv("SC_M_B", "origB")

	seen := map[string]string{}
	err := ScopedMany(repo, map[string]string{
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

func TestScoped_restoresOnPanic(t *testing.T) {
	repo := NewRepository()
	t.Setenv("SC_P", "orig")

	defer func() {
		_ = recover()
		require.Equal(t, "orig", repo.Get("SC_P"))
	}()

	_ = Scoped(repo, "SC_P", "temp", func() {
		panic("boom")
	})
}
