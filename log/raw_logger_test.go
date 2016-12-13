package log

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type TestFormattable struct {
	a string
	b string
}

// String ...
func (f TestFormattable) String() string {
	return fmt.Sprintf("%s %s", f.a, f.b)
}

// JSON ...
func (f TestFormattable) JSON() string {
	return fmt.Sprintf(`{"a":"%s","b":"%s"}`, f.a, f.b)
}

func TestRawPrint(t *testing.T) {
	t.Log("Default Formattable (Message)")
	{
		var b bytes.Buffer
		logger := NewRawLogger(&b)

		logger.Printd(Message{Content: "test"})
		require.Equal(t, "test\n", b.String())
	}

	t.Log("Custom Formattable")
	{
		var b bytes.Buffer
		logger := NewRawLogger(&b)

		test := TestFormattable{
			a: "log",
			b: "test",
		}

		logger.Printd(test)
		require.Equal(t, "log test\n", b.String())
	}
}

func TestRawPrintf(t *testing.T) {
	t.Log("string")
	{
		var b bytes.Buffer
		logger := NewRawLogger(&b)

		logger.Printf("test")
		require.Equal(t, "test\n", b.String())
	}

	t.Log("format")
	{
		var b bytes.Buffer
		logger := NewRawLogger(&b)

		logger.Printf("%s", "test")
		require.Equal(t, "test\n", b.String())
	}

	t.Log("complex format")
	{
		var b bytes.Buffer
		logger := NewRawLogger(&b)

		logger.Printf("%s %s", "log", "test")
		require.Equal(t, "log test\n", b.String())
	}
}
