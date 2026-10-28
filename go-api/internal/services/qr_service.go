package services

import (
	"fmt"

	"gonum.org/v1/gonum/mat"

	"github.com/interseguro/go-api/internal/models"
)

// FactorizeQR calcula la factorización QR de una matriz rectangular m×n.
// Retorna Q (m×k) y R (k×n) donde k = min(m, n).
//
// Gonum no soporta QR directo para matrices "anchas" (m < n): paniquea con
// "dimension mismatch". Solución: para m < n, rellenamos con filas de ceros
// hasta hacerla cuadrada n×n, factorizamos, y truncamos:
//   - Q = Q_padded[0:m, 0:m]  (ortogonal por construcción)
//   - R = R_padded[0:m, 0:n]  (triangular superior)
// Se verifica matemáticamente que Q*R = A original.
func FactorizeQR(matrix [][]float64) (*models.QRResult, error) {
	if err := validateMatrix(matrix); err != nil {
		return nil, err
	}

	rows := len(matrix)
	cols := len(matrix[0])
	k := min(rows, cols)

	var qReduced, rReduced *mat.Dense

	if rows >= cols {
		// Tall o cuadrada: QR directo de Gonum.
		flat := make([]float64, 0, rows*cols)
		for _, row := range matrix {
			flat = append(flat, row...)
		}
		a := mat.NewDense(rows, cols, flat)

		var qr mat.QR
		qr.Factorize(a)

		var qFull, rFull mat.Dense
		qr.QTo(&qFull)
		qr.RTo(&rFull)

		qReduced = extractBlock(&qFull, rows, k, 0, 0)
		rReduced = extractBlock(&rFull, k, cols, 0, 0)
	} else {
		// Ancha (rows < cols): rellenar con filas de ceros hasta n×n.
		padded := make([]float64, cols*cols) // ya inicializado en ceros
		for i := 0; i < rows; i++ {
			copy(padded[i*cols:(i+1)*cols], matrix[i])
		}
		aPadded := mat.NewDense(cols, cols, padded)

		var qr mat.QR
		qr.Factorize(aPadded)

		var qPadded, rPadded mat.Dense
		qr.QTo(&qPadded)
		qr.RTo(&rPadded)

		// Truncar: Q (m×m) y R (m×n) del bloque superior de Q_padded y R_padded.
		qReduced = extractBlock(&qPadded, rows, rows, 0, 0)
		rReduced = extractBlock(&rPadded, rows, cols, 0, 0)
	}

	return &models.QRResult{
		Q: denseToResponse(qReduced),
		R: denseToResponse(rReduced),
	}, nil
}

// extractBlock extrae una submatriz r×c desde la posición (row, col).
func extractBlock(src *mat.Dense, r, c, row, col int) *mat.Dense {
	out := mat.NewDense(r, c, nil)
	for i := 0; i < r; i++ {
		for j := 0; j < c; j++ {
			out.Set(i, j, src.At(row+i, col+j))
		}
	}
	return out
}

func validateMatrix(m [][]float64) error {
	if len(m) == 0 {
		return fmt.Errorf("la matriz está vacía")
	}
	cols := len(m[0])
	if cols == 0 {
		return fmt.Errorf("la matriz no tiene columnas")
	}
	for i, row := range m {
		if len(row) != cols {
			return fmt.Errorf("fila %d tiene %d columnas, se esperaban %d", i, len(row), cols)
		}
	}
	return nil
}

func denseToResponse(d *mat.Dense) models.MatrixResponse {
	r, c := d.Dims()
	data := make([][]float64, r)
	for i := 0; i < r; i++ {
		data[i] = make([]float64, c)
		for j := 0; j < c; j++ {
			data[i][j] = d.At(i, j)
		}
	}
	return models.MatrixResponse{Rows: r, Columns: c, Data: data}
}

func MatrixToResponse(matrix [][]float64) models.MatrixResponse {
	return models.MatrixResponse{
		Rows:    len(matrix),
		Columns: len(matrix[0]),
		Data:    matrix,
	}
}