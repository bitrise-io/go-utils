package readerutil

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadLongLine(t *testing.T) {
	t.Log("Empty string")
	{
		reader := bufio.NewReader(strings.NewReader(``))
		line, err := ReadLongLine(reader)
		require.Equal(t, io.EOF, err)
		require.Equal(t, "", line)
	}

	t.Log("Single line")
	{
		reader := bufio.NewReader(strings.NewReader(`a single line`))
		line, err := ReadLongLine(reader)
		require.NoError(t, err)
		require.Equal(t, "a single line", line)
		// read one more
		line, err = ReadLongLine(reader)
		require.Equal(t, io.EOF, err)
		require.Equal(t, "", line)
	}

	t.Log("Two lines")
	{
		reader := bufio.NewReader(strings.NewReader(`first line
second line`))
		// first line
		line, readErr := ReadLongLine(reader)
		require.NoError(t, readErr)
		require.Equal(t, "first line", line)
		// second line
		line, readErr = ReadLongLine(reader)
		require.NoError(t, readErr)
		require.Equal(t, "second line", line)
		// read one more
		line, readErr = ReadLongLine(reader)
		require.Equal(t, io.EOF, readErr)
		require.Equal(t, "", line)
	}

	t.Log("Multi line, with long line")
	{
		inputStr := fmt.Sprintf(`first line
second line
third, really long line: %s
  fourth line
`, strings.Repeat("-", 1000000))

		reader := bufio.NewReader(strings.NewReader(inputStr))
		//
		lines := []string{}
		line, readErr := ReadLongLine(reader)
		for ; readErr == nil; line, readErr = ReadLongLine(reader) {
			lines = append(lines, line)
		}
		// ideally the error will be io.EOF
		require.Equal(t, io.EOF, readErr)
		//
		require.Equal(t, 4, len(lines))
		require.Equal(t, "first line", lines[0])
		require.Equal(t, "second line", lines[1])
		// check the start of the long line
		require.Equal(t, "third, really long line: ---", lines[2][0:28])
		require.Equal(t, "  fourth line", lines[3])
	}
}
