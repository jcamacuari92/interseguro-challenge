package models

// MatrixResponse representa una matriz con metadatos de dimensiones.
type MatrixResponse struct {
	Rows    int         `json:"rows"`
	Columns int         `json:"columns"`
	Data    [][]float64 `json:"data"`
}

// MatrixRequest es la entrada del cliente.
type MatrixRequest struct {
	Matrix [][]float64 `json:"matrix"`
}

// QRResult agrupa las matrices Q y R.
type QRResult struct {
	Q MatrixResponse `json:"q"`
	R MatrixResponse `json:"r"`
}

// StatsRequest es el payload que Go envía a Node.js.
// Enviamos las matrices como un slice de matrices (Q y R).
type StatsRequest struct {
	Matrices [][][]float64 `json:"matrices"`
}

// DiagonalInfo describe el chequeo de matrices diagonales.
type DiagonalInfo struct {
	AnyDiagonal bool   `json:"anyDiagonal"`
	Matrices    []bool `json:"matrices"`
}

// StatsResponse es el resultado calculado por Node.js.
type StatsResponse struct {
	Max          float64      `json:"max"`
	Min          float64      `json:"min"`
	Average      float64      `json:"average"`
	TotalSum     float64      `json:"totalSum"`
	Count        int          `json:"count"`
	DiagonalInfo DiagonalInfo `json:"diagonalInfo"`
}

// QRResponse es la respuesta final del endpoint /api/v1/qr.
type QRResponse struct {
	Operation  string         `json:"operation"`
	Original   MatrixResponse `json:"original"`
	QR         QRResult       `json:"qr"`
	Statistics StatsResponse  `json:"statistics"`
}

// RotateResponse es la respuesta final del endpoint /api/v1/rotate.
type RotateResponse struct {
	Operation  string         `json:"operation"`
	Original   MatrixResponse `json:"original"`
	Rotated    MatrixResponse `json:"rotated"`
	Statistics StatsResponse  `json:"statistics"`
}
