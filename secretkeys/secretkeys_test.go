package secretkeys

import (
	"github.com/bitrise-io/go-utils/v2/mocks"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFormat(t *testing.T) {
	tests := []struct {
		name string
		keys []string
		want string
	}{
		{
			name: "Empty keys",
			keys: []string{},
			want: "",
		},
		{
			name: "Some keys",
			keys: []string{"ABC", "DEF", "GHI"},
			want: "ABC,DEF,GHI",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewManager().Format(tt.keys)
			if tt.want != got {
				t.Errorf("got formatted keys = %s, want %s", got, tt.want)
			}
		})
	}
}

func TestLoad(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []string
	}{
		{
			name:  "No saved secret keys",
			input: "",
			want:  []string{""},
		},
		{
			name:  "Has saved secret keys",
			input: "ABC,DEF,GHI",
			want:  []string{"ABC", "DEF", "GHI"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			envMock := mocks.NewRepository(t)
			envMock.On("Get", EnvKey).Return(tt.input)
			got := NewManager().Load(envMock)
			assert.Equal(t, tt.want, got, "got formatted keys = %s, want %s", got, tt.want)
		})
	}
}
