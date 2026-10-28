package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRotate90CW(t *testing.T) {
	tests := []struct {
		name     string
		input    [][]float64
		expected [][]float64
	}{
		{
			name:     "2x2",
			input:    [][]float64{{1, 2}, {3, 4}},
			expected: [][]float64{{3, 1}, {4, 2}},
		},
		{
			name:     "2x3",
			input:    [][]float64{{1, 2, 3}, {4, 5, 6}},
			expected: [][]float64{{4, 1}, {5, 2}, {6, 3}},
		},
		{
			name:     "3x2",
			input:    [][]float64{{1, 2}, {3, 4}, {5, 6}},
			expected: [][]float64{{5, 3, 1}, {6, 4, 2}},
		},
		{
			name:     "1x3",
			input:    [][]float64{{1, 2, 3}},
			expected: [][]float64{{1}, {2}, {3}},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := Rotate90CW(tc.input)
			require.NoError(t, err)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestRotate90CW_Invalid(t *testing.T) {
	tests := []struct {
		name  string
		input [][]float64
	}{
		{"empty", [][]float64{}},
		{"ragged", [][]float64{{1, 2}, {3}}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := Rotate90CW(tc.input)
			assert.Error(t, err)
		})
	}
}