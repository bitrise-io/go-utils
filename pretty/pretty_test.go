package pretty

import (
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestObject_map(t *testing.T) {
	o := map[string]any{
		"key": "value",
		"slice_key": []string{
			"item1",
			"item2",
		},
	}
	const want = `{
	"key": "value",
	"slice_key": [
		"item1",
		"item2"
	]
}`
	require.Equal(t, want, Object(o))
}

func TestObject_primitive(t *testing.T) {
	require.Equal(t, `"hello"`, Object("hello"))
	require.Equal(t, `42`, Object(42))
	require.Equal(t, `null`, Object(nil))
}

func TestObject_struct(t *testing.T) {
	type payload struct {
		Name  string `json:"name"`
		Count int    `json:"count"`
	}
	const want = `{
	"name": "bitrise",
	"count": 3
}`
	require.Equal(t, want, Object(payload{Name: "bitrise", Count: 3}))
}

func TestObject_unmarshalableFallsBackToFmt(t *testing.T) {
	// NaN cannot be marshaled as JSON; Object must fall back to fmt.Sprint.
	got := Object(math.NaN())
	require.Equal(t, "NaN", got)
}
