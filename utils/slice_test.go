package utils

import (
	"slices"
	"testing"
)

func TestRemoveElement(t *testing.T) {
	tests := []struct {
		name     string
		input    []int
		element  int
		expected []int
	}{
		{
			name:     "Remove existing element",
			input:    []int{1, 2, 3, 4, 5},
			element:  3,
			expected: []int{1, 2, 4, 5},
		},
		{
			name:     "Remove non-existing element",
			input:    []int{1, 2, 3, 4, 5},
			element:  6,
			expected: []int{1, 2, 3, 4, 5},
		},
		{
			name:     "Remove element from empty slice",
			input:    []int{},
			element:  1,
			expected: []int{},
		},
		{
			name:     "Remove first element",
			input:    []int{1, 2, 3, 4, 5},
			element:  1,
			expected: []int{2, 3, 4, 5},
		},
		{
			name:     "Remove last element",
			input:    []int{1, 2, 3, 4, 5},
			element:  5,
			expected: []int{1, 2, 3, 4},
		},
		{
			name:     "Remove element from single-element slice",
			input:    []int{1},
			element:  1,
			expected: []int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RemoveElement(tt.input, tt.element)
			if !slices.Equal(result, tt.expected) {
				t.Fatalf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}
