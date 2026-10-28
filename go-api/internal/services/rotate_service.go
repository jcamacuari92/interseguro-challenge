package services

import "fmt"

// Rotate90CW rota la matriz 90 grados en sentido horario.
// Regla: nueva[i][j] = original[m-1-j][i]
func Rotate90CW(matrix [][]float64) ([][]float64, error) {
	if len(matrix) == 0 {
		return nil, fmt.Errorf("la matriz está vacía")
	}
	m := len(matrix)
	n := len(matrix[0])

	for i, row := range matrix {
		if len(row) != n {
			return nil, fmt.Errorf("fila %d tiene %d columnas, se esperaban %d", i, len(row), n)
		}
	}

	result := make([][]float64, n)
	for i := 0; i < n; i++ {
		result[i] = make([]float64, m)
		for j := 0; j < m; j++ {
			result[i][j] = matrix[m-1-j][i]
		}
	}
	return result, nil
}