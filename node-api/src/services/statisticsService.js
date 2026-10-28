// Umbral para comparar flotantes con 0 (tolerancia numérica).
const EPSILON = 1e-10;

/**
 * Valida que el payload sea un array no vacío de matrices rectangulares
 * de números finitos.
 * @param {unknown} matrices
 */
function validateMatrices(matrices) {
  if (!Array.isArray(matrices) || matrices.length === 0) {
    throw new Error('Se esperaba un array no vacío de matrices');
  }

  for (let k = 0; k < matrices.length; k++) {
    const m = matrices[k];
    if (!Array.isArray(m) || m.length === 0) {
      throw new Error(`Matriz ${k} inválida: debe ser un array no vacío`);
    }

    const cols = m[0].length;
    if (!Array.isArray(m[0]) || cols === 0) {
      throw new Error(`Matriz ${k} inválida: las filas deben ser arrays no vacíos`);
    }

    for (let i = 0; i < m.length; i++) {
      const row = m[i];
      if (!Array.isArray(row) || row.length !== cols) {
        throw new Error(`Matriz ${k} fila ${i} inválida: columnas inconsistentes`);
      }
      for (const v of row) {
        if (typeof v !== 'number' || !Number.isFinite(v)) {
          throw new Error(`Matriz ${k} contiene un valor no numérico en fila ${i}`);
        }
      }
    }
  }
}

/**
 * Determina si una matriz es diagonal: cuadrada y con todos los
 * elementos fuera de la diagonal principal ≈ 0.
 */
function isDiagonal(matrix) {
  const rows = matrix.length;
  const cols = matrix[0].length;
  if (rows !== cols) return false;

  for (let i = 0; i < rows; i++) {
    for (let j = 0; j < cols; j++) {
      if (i !== j && Math.abs(matrix[i][j]) > EPSILON) {
        return false;
      }
    }
  }
  return true;
}

/**
 * Calcula estadísticas sobre un conjunto de matrices.
 * Estrategia: aplana todo en una pasada → O(N) donde N = total de elementos.
 *
 * @param {number[][][]} matrices
 * @returns {{
 *   max: number, min: number, average: number,
 *   totalSum: number, count: number,
 *   diagonalInfo: { anyDiagonal: boolean, matrices: boolean[] }
 * }}
 */
function computeStatistics(matrices) {
  let max = -Infinity;
  let min = Infinity;
  let sum = 0;
  let count = 0;

  // Pasada única sobre todos los elementos de todas las matrices.
  for (const matrix of matrices) {
    for (const row of matrix) {
      for (const v of row) {
        if (v > max) max = v;
        if (v < min) min = v;
        sum += v;
        count++;
      }
    }
  }

  if (count === 0) {
    throw new Error('No hay elementos para calcular estadísticas');
  }

  const diagonalChecks = matrices.map(isDiagonal);

  return {
    max,
    min,
    average: sum / count,
    totalSum: sum,
    count,
    diagonalInfo: {
      anyDiagonal: diagonalChecks.some(Boolean),
      matrices: diagonalChecks,
    },
  };
}

module.exports = { validateMatrices, computeStatistics, isDiagonal };