const {
  validateMatrices,
  computeStatistics,
  isDiagonal,
} = require('../src/services/statisticsService');

describe('validateMatrices', () => {
  test('acepta matrices validas', () => {
    expect(() => validateMatrices([[[1, 2], [3, 4]]])).not.toThrow();
  });

  test('rechaza entrada no-array', () => {
    expect(() => validateMatrices(null)).toThrow();
    expect(() => validateMatrices('abc')).toThrow();
  });

  test('rechaza array vacio', () => {
    expect(() => validateMatrices([])).toThrow();
  });

  test('rechaza matriz vacia', () => {
    expect(() => validateMatrices([[]])).toThrow();
  });

  test('rechaza matriz irregular', () => {
    expect(() => validateMatrices([[[1, 2], [3]]])).toThrow();
  });

  test('rechaza valores no numericos', () => {
    expect(() => validateMatrices([[[1, 'x']]])).toThrow();
    expect(() => validateMatrices([[[NaN]]])).toThrow();
    expect(() => validateMatrices([[[Infinity]]])).toThrow();
  });
});

describe('isDiagonal', () => {
  test('true para matriz diagonal', () => {
    expect(isDiagonal([[1, 0], [0, 2]])).toBe(true);
  });

  test('true para identidad', () => {
    expect(isDiagonal([[1, 0, 0], [0, 1, 0], [0, 0, 1]])).toBe(true);
  });

  test('false para no cuadrada', () => {
    expect(isDiagonal([[1, 0, 0], [0, 1, 0]])).toBe(false);
  });

  test('false con off-diagonal no cero', () => {
    expect(isDiagonal([[0, 5], [0, 0]])).toBe(false);
  });

  test('true con off-diagonal casi cero (epsilon)', () => {
    expect(isDiagonal([[1, 1e-12], [1e-12, 2]])).toBe(true);
  });
});

describe('computeStatistics', () => {
  test('calcula correctamente para dos matrices', () => {
    const result = computeStatistics([[[1, 0], [0, 2]], [[0, 5], [0, 0]]]);
    expect(result.max).toBe(5);
    expect(result.min).toBe(0);
    expect(result.totalSum).toBe(8);
    expect(result.count).toBe(8);
    expect(result.average).toBe(1);
    expect(result.diagonalInfo.anyDiagonal).toBe(true);
    expect(result.diagonalInfo.matrices).toEqual([true, false]);
  });

  test('valores negativos', () => {
    const result = computeStatistics([[[-1, -2], [-3, -4]]]);
    expect(result.max).toBe(-1);
    expect(result.min).toBe(-4);
    expect(result.totalSum).toBe(-10);
    expect(result.average).toBe(-2.5);
  });

  test('un solo elemento', () => {
    const result = computeStatistics([[[42]]]);
    expect(result.max).toBe(42);
    expect(result.min).toBe(42);
    expect(result.count).toBe(1);
    expect(result.average).toBe(42);
  });
});