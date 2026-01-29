package sliceutil

import (
	"slices"
	"strings"
	"testing"
)

func TestUnique(t *testing.T) {
	tests := []struct {
		name  string
		input []string
		want  []string
	}{
		{"empty", []string{}, []string{}},
		{"no duplicates", []string{"a", "b", "c"}, []string{"a", "b", "c"}},
		{"with duplicates", []string{"a", "b", "a", "c", "b"}, []string{"a", "b", "c"}},
		{"all same", []string{"x", "x", "x"}, []string{"x"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Unique(tt.input)
			if len(got) != len(tt.want) {
				t.Errorf("Unique() len = %d, want %d", len(got), len(tt.want))
			}
			for _, item := range tt.want {
				if !slices.Contains(got, item) {
					t.Errorf("Unique() missing expected item %q", item)
				}
			}
		})
	}
}

func TestUniqueInt(t *testing.T) {
	got := Unique([]int{1, 2, 2, 3, 1})
	want := []int{1, 2, 3}

	if len(got) != len(want) {
		t.Fatalf("len = %d, want %d", len(got), len(want))
	}
	for _, item := range want {
		if !slices.Contains(got, item) {
			t.Errorf("missing expected item %d", item)
		}
	}
}

func TestFilter(t *testing.T) {
	tests := []struct {
		name       string
		input      []int
		filterFunc func(int) bool
		want       []int
	}{
		{"empty", []int{}, func(i int) bool { return true }, []int{}},
		{"none match", []int{1, 2, 3}, func(i int) bool { return i > 10 }, []int{}},
		{"all match", []int{1, 2, 3}, func(i int) bool { return i > 0 }, []int{1, 2, 3}},
		{"even numbers", []int{1, 2, 3, 4, 5, 6}, func(i int) bool { return i%2 == 0 }, []int{2, 4, 6}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Filter(tt.input, tt.filterFunc)
			if !slices.Equal(got, tt.want) {
				t.Errorf("Filter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFilterStrings(t *testing.T) {
	input := []string{"", "a", "", "b", "   "}
	got := Filter(input, func(s string) bool { return s != "" })
	want := []string{"a", "b", "   "}

	if !slices.Equal(got, want) {
		t.Errorf("Filter() = %v, want %v", got, want)
	}
}

func TestMap(t *testing.T) {
	input := []string{"", " a", "", " b ", "   "}
	got := Map(input, func(s string) string { return strings.TrimSpace(s) })
	want := []string{"", "a", "", "b", ""}

	if !slices.Equal(got, want) {
		t.Errorf("Filter() = %v, want %v", got, want)
	}
}
