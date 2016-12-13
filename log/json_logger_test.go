package log

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestJSONPrint(t *testing.T) {
	t.Log("Default Formattable (Message)")
	{
		var b bytes.Buffer
		logger := NewJSONLoger(&b)

		logger.PrintO(Message{Content: "test"})
		require.Equal(t, `{"content":"test"}`+"\n", b.String())
	}

	t.Log("Custom Formattable")
	{
		var b bytes.Buffer
		logger := NewJSONLoger(&b)

		test := TestFormattable{
			a: "log",
			b: "test",
		}

		logger.PrintO(test)
		require.Equal(t, `{"a":"log","b":"test"}`+"\n", b.String())
	}
}

func TestJSONPrintf(t *testing.T) {
	t.Log("string")
	{
		var b bytes.Buffer
		logger := NewJSONLoger(&b)

		logger.Printf("test")
		require.Equal(t, `{"content":"test"}`+"\n", b.String())
	}

	t.Log("format")
	{
		var b bytes.Buffer
		logger := NewJSONLoger(&b)

		logger.Printf("%s", "test")
		require.Equal(t, `{"content":"test"}`+"\n", b.String())
	}

	t.Log("complex format")
	{
		var b bytes.Buffer
		logger := NewJSONLoger(&b)

		logger.Printf("%s %s", "log", "test")
		require.Equal(t, `{"content":"log test"}`+"\n", b.String())
	}
}
