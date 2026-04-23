package urlutil

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestJoin(t *testing.T) {
	tests := []struct {
		name     string
		parts    []string
		expected string
	}{
		{
			name:     "single segment",
			parts:    []string{"https://bitrise.io", "something"},
			expected: "https://bitrise.io/something",
		},
		{
			name:     "nested segment",
			parts:    []string{"https://bitrise.io", "something/a"},
			expected: "https://bitrise.io/something/a",
		},
		{
			name:     "deeper nested segment",
			parts:    []string{"https://bitrise.io", "something/a/b"},
			expected: "https://bitrise.io/something/a/b",
		},
		{
			name:     "preserves trailing slash on last segment",
			parts:    []string{"https://bitrise.io", "something/a/b/"},
			expected: "https://bitrise.io/something/a/b/",
		},

		{
			name:     "strips leading slash on second element",
			parts:    []string{"https://bitrise.io", "/something"},
			expected: "https://bitrise.io/something",
		},
		{
			name:     "strips leading slash, keeps inner",
			parts:    []string{"https://bitrise.io", "/something/a"},
			expected: "https://bitrise.io/something/a",
		},
		{
			name:     "strips leading slash, keeps depth",
			parts:    []string{"https://bitrise.io", "/something/a/b"},
			expected: "https://bitrise.io/something/a/b",
		},
		{
			name:     "strips leading slash, preserves trailing",
			parts:    []string{"https://bitrise.io", "/something/a/b/"},
			expected: "https://bitrise.io/something/a/b/",
		},

		{
			name:     "trailing slash on host",
			parts:    []string{"https://bitrise.io/", "/something"},
			expected: "https://bitrise.io/something",
		},
		{
			name:     "trailing slash on host with nested segment",
			parts:    []string{"https://bitrise.io/", "/something/a"},
			expected: "https://bitrise.io/something/a",
		},
		{
			name:     "trailing slash on host with deeper segment",
			parts:    []string{"https://bitrise.io/", "/something/a/b"},
			expected: "https://bitrise.io/something/a/b",
		},
		{
			name:     "trailing slash on host preserves trailing slash on last",
			parts:    []string{"https://bitrise.io/", "/something/a/b/"},
			expected: "https://bitrise.io/something/a/b/",
		},

		{
			name:     "double slashes are collapsed",
			parts:    []string{"https://bitrise.io//", "//something"},
			expected: "https://bitrise.io/something",
		},
		{
			name:     "double slashes collapsed, nested segment",
			parts:    []string{"https://bitrise.io//", "//something/a"},
			expected: "https://bitrise.io/something/a",
		},
		{
			name:     "double slashes collapsed, deeper segment",
			parts:    []string{"https://bitrise.io//", "//something/a/b"},
			expected: "https://bitrise.io/something/a/b",
		},
		{
			name:     "double slashes collapsed, trailing slash preserved",
			parts:    []string{"https://bitrise.io//", "//something/a/b/"},
			expected: "https://bitrise.io/something/a/b/",
		},

		{
			name:     "multiple segments with host path",
			parts:    []string{"https://bitrise-steplib-collection.s3.amazonaws.com/steps", "activate-ssh-key", "assets", "icon.svg"},
			expected: "https://bitrise-steplib-collection.s3.amazonaws.com/steps/activate-ssh-key/assets/icon.svg",
		},
		{
			name:     "multiple segments with trailing slash on host",
			parts:    []string{"https://bitrise-steplib-collection.s3.amazonaws.com/steps/", "activate-ssh-key", "assets", "icon.svg"},
			expected: "https://bitrise-steplib-collection.s3.amazonaws.com/steps/activate-ssh-key/assets/icon.svg",
		},
		{
			name:     "multiple segments, one leading slash",
			parts:    []string{"https://bitrise-steplib-collection.s3.amazonaws.com/steps/", "/activate-ssh-key", "assets", "icon.svg"},
			expected: "https://bitrise-steplib-collection.s3.amazonaws.com/steps/activate-ssh-key/assets/icon.svg",
		},
		{
			name:     "multiple segments, two leading slashes",
			parts:    []string{"https://bitrise-steplib-collection.s3.amazonaws.com/steps/", "/activate-ssh-key", "/assets", "icon.svg"},
			expected: "https://bitrise-steplib-collection.s3.amazonaws.com/steps/activate-ssh-key/assets/icon.svg",
		},
		{
			name:     "multiple segments, all leading slashes",
			parts:    []string{"https://bitrise-steplib-collection.s3.amazonaws.com/steps/", "/activate-ssh-key", "/assets", "/icon.svg"},
			expected: "https://bitrise-steplib-collection.s3.amazonaws.com/steps/activate-ssh-key/assets/icon.svg",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := Join(tc.parts...)
			require.NoError(t, err)
			require.Equal(t, tc.expected, got)
		})
	}
}

func TestJoin_errors(t *testing.T) {
	tests := []struct {
		name    string
		parts   []string
		wantErr string
	}{
		{
			name:    "no elements",
			parts:   []string{},
			wantErr: "No elements defined to Join",
		},
		{
			name:    "missing host",
			parts:   []string{"https://", "bitrise.io"},
			wantErr: "No Host defined",
		},
		{
			name:    "missing scheme",
			parts:   []string{"bitrise.io", "something"},
			wantErr: "No Scheme defined",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := Join(tc.parts...)
			require.Error(t, err)
			require.Equal(t, tc.wantErr, err.Error())
			require.Equal(t, "", got)
		})
	}
}
