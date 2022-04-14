package analytics

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type simpleTestStruct struct {
	Name             string `json:"name"`
	Address          string `json:"address,omitempty"`
	Age              int    `json:"age"`
	Happy            bool   `json:"happy"`
	NotJSONField     string
	IgnoredJSONField string `json:"-"`
}

type complexTestStruct struct {
	Name            string          `json:"name"`
	Address         string          `json:"address,omitempty"`
	City            string          `json:"city"`
	Age             int             `json:"age"`
	Happy           bool            `json:"happy"`
	Cars            []string        `json:"cars"`
	ColorCollection map[string]bool `json:"color_collection"`
}

type emptyTestStruct struct {
}

func TestToProperty(t *testing.T) {
	tests := []struct {
		name string
		item interface{}
		want Properties
	}{
		{
			name: "Simple Test",
			item: simpleTestStruct{
				Name:         "John",
				Age:          42,
				Happy:        true,
				NotJSONField: "Not a json field",
			},
			want: Properties{
				"name":         "John",
				"age":          42,
				"happy":        true,
				"NotJSONField": "Not a json field",
			},
		},
		{
			name: "Complex Test",
			item: complexTestStruct{
				Name:            "John",
				Age:             42,
				Happy:           true,
				Cars:            []string{"Ford", "Chevy"},
				ColorCollection: map[string]bool{"red": true, "blue": false},
			},
			want: Properties{
				"name":             "John",
				"city":             "",
				"age":              42,
				"happy":            true,
				"cars":             []string{"Ford", "Chevy"},
				"color_collection": map[string]bool{"red": true, "blue": false},
			},
		},
		{
			name: "Empty Test",
			item: emptyTestStruct{},
			want: Properties{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, ToProperty(tt.item), "ToProperty(%v)", tt.item)
		})
	}
}
