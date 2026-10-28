package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFactorizeQR_RoundTrip(t *testing.T) {
	tests := []struct {
		name   string
		matrix [][]float64
	}{
		{"2x2", [][]float64{{1, 2}, {3, 4}}},
		{"3x3", [][]float64{{12, -51, 4}, {6, 167, -68}, {-4, 24, -41}}},
		{"3x2 tall", [][]float64{{1, 2}, {3, 4}, {5, 6}}},
		{"2x3 wide", [][]float64{{1, 2, 3}, {4, 5, 6}}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			res, err := FactorizeQR(tc.matrix)
			require.NoError(t, err)

			reconstructed := multiplyQandR(res.Q.Data, res.R.Data)
			for i := range tc.matrix {
				for j := range tc.matrix[i] {
					assert.InDelta(t, tc.matrix[i][j], reconstructed[i][j], 1e-9,
						"Q*R mismatch at [%d][%d]", i, j)
				}
			}
		})
	}
}

func TestFactorizeQR_InvalidInput(t *testing.T) {
	tests := []struct {
		name   string
		matrix [][]float64
	}{
		{"empty", [][]float64{}},
		{"empty row", [][]float64{{}}},
		{"ragged", [][]float64{{1, 2}, {3}}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := FactorizeQR(tc.matrix)
			assert.Error(t, err)
		})
	}
}

// multiplyQandR multiplica Q (m x k) por R (k x n) para verificar Q*R == A.
func multiplyQandR(q, r [][]float64) [][]float64 {
	if len(q) == 0 || len(r) == 0 {
		return nil
	}
	rows := len(q)
	k := len(q[0])
	cols := len(r[0])
	result := make([][]float64, rows)
	for i := 0; i < rows; i++ {
		result[i] = make([]float64, cols)
		for j := 0; j < cols; j++ {
			sum := 0.0
			for x := 0; x < k; x++ {
				sum += q[i][x] * r[x][j]
			}
			result[i][j] = sum
		}
	}
	return result
}