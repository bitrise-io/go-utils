package command

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNestedReason(t *testing.T) {
	tests := []struct {
		name  string
		lines []string
		want  string
	}{
		{
			name:  "no lines",
			lines: nil,
			want:  "",
		},
		{
			name:  "lines that are not indented are one link",
			lines: []string{"first error", "second error"},
			want:  "first error\nsecond error",
		},
		{
			name: "a line indented further is a link under the one above it",
			lines: []string{
				"Execution failed for task ':app:compileDebugKotlin'.",
				"  A failure occurred while executing the Kotlin compiler",
				"    Foo.kt:4:30: Unresolved reference 'missingHelper'.",
			},
			want: "Execution failed for task ':app:compileDebugKotlin': " +
				"A failure occurred while executing the Kotlin compiler: " +
				"Foo.kt:4:30: Unresolved reference 'missingHelper'.",
		},
		{
			name: "lines at the same indent share the link",
			lines: []string{
				"Execution failed for task ':app:compileDebugKotlin'.",
				"  Foo.kt:4:30: Unresolved reference 'missingHelper'.",
				"  Foo.kt:6:30: Unresolved reference 'anotherMissingHelper'.",
			},
			want: "Execution failed for task ':app:compileDebugKotlin': " +
				"Foo.kt:4:30: Unresolved reference 'missingHelper'.\n" +
				"Foo.kt:6:30: Unresolved reference 'anotherMissingHelper'.",
		},
		{
			name: "a line that steps back out joins the deepest link, because a chain holds no branches",
			lines: []string{
				"first test FAILED",
				"  its detail",
				"second test FAILED",
			},
			want: "first test FAILED: its detail\nsecond test FAILED",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reason := nestedReason(tt.lines)

			if tt.want == "" {
				require.Nil(t, reason)

				return
			}
			require.Equal(t, tt.want, reason.Error())
		})
	}
}
