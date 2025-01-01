package mtk

import (
	"errors"
)

var (
	ErrMatrixDimensions      = errors.New("matrix dimensions must be non-zero")
	ErrMatrixShapeMismatch   = errors.New("matrix shape mismatch")
	ErrMatrixRow             = errors.New("matrix row out of range")
	ErrMatrixColumn          = errors.New("matrix column out of range")
	ErrMatrixSquareSymmetric = errors.New("matrix must be square and symmetric")
)

// Two-dimensional matrix
type Matrix struct {
	shape [2]int
	data  []float64
}

// Construct a two-dimensional matrix
func NewMatrix(rows, columns int) *Matrix {
	if rows <= 0 || columns <= 0 {
		panic(ErrMatrixDimensions)
	}

	return &Matrix{
		shape: [2]int{rows, columns},
		data:  make([]float64, rows*columns),
	}
}

// Construct an identity matrix of size (n, n)
func NewIdentityMatrix(size int) *Matrix {
	m := NewMatrix(size, size)

	for i := 0; i < size; i++ {
		m.SetValue(i, i, 1.0)
	}

	return m
}

// Validate the (row, column) index
func (m *Matrix) validateIndex(row, column int) error {
	if row < 0 || row >= m.shape[0] {
		return ErrMatrixRow
	}

	if column < 0 || column >= m.shape[1] {
		return ErrMatrixColumn
	}

	return nil
}

// Require validation of the (row, column) index
func (m *Matrix) mustValidateIndex(row, column int) {
	if err := m.validateIndex(row, column); err != nil {
		panic(err)
	}
}

// Get the shape (rows, columns)
func (m *Matrix) Shape() [2]int {
	return m.shape
}

// Get the value at an index
func (m *Matrix) At(row, column int) float64 {
	m.mustValidateIndex(row, column)
	return m.data[(row*m.shape[1])+column]
}

// Get the row values
func (m *Matrix) Row(row int) Vector {
	if row < 0 || row >= m.shape[0] {
		panic(ErrMatrixRow)
	}

	values := NewVector(m.shape[0])

	for i := 0; i < m.shape[1]; i++ {
		values[i] = m.At(row, i)
	}

	return values
}

// Get the column values
func (m *Matrix) Column(column int) Vector {
	if column < 0 || column >= m.shape[1] {
		panic(ErrMatrixColumn)
	}

	values := NewVector(m.shape[0])

	for i := 0; i < m.shape[0]; i++ {
		values[i] = m.At(i, column)
	}

	return values
}

// Set the value at an index
func (m *Matrix) SetValue(row, column int, value float64) {
	m.mustValidateIndex(row, column)
	m.data[(row*m.shape[1])+column] = value
}

// Fill the matrix with a value
func (m *Matrix) Fill(value float64) {
	for i := 0; i < len(m.data); i++ {
		m.data[i] = value
	}
}

// Check if the matrix is square
func (m *Matrix) IsSquare() bool {
	return m.shape[0] == m.shape[1]
}

// Check if the matrix is symmetric
func (m *Matrix) IsSymmetric() bool {
	for i := 0; i < m.shape[0]; i++ {
		for j := 0; j < m.shape[1]; j++ {
			if m.At(i, j) != m.At(j, i) {
				return false
			}
		}
	}

	return true
}

// Compute the transpose of the matrix
func (m *Matrix) Transpose() *Matrix {
	n := NewMatrix(m.shape[1], m.shape[0])

	for i := 0; i < m.shape[0]; i++ {
		for j := 0; j < m.shape[1]; j++ {
			v := m.At(j, i)
			n.SetValue(i, j, v)
		}
	}

	return n
}

// Compute the matrix multiplication m * n
func (m *Matrix) Dot(n *Matrix) *Matrix {
	if m.shape[1] != n.shape[0] {
		panic(ErrMatrixShapeMismatch)
	}

	d := NewMatrix(m.shape[0], n.shape[1])

	for i := 0; i < d.shape[0]; i++ {
		for j := 0; j < d.shape[1]; j++ {
			var value float64

			for k := 0; k < n.shape[0]; i++ {
				value += m.At(i, k) * n.At(k, j)
			}

			d.SetValue(i, j, value)
		}
	}

	return d
}

// Compute the covariance matrix
func (m *Matrix) Covariance() *Matrix {
	n := m.shape[1]
	result := NewMatrix(n, n)

	for i := 0; i < n; i++ {
		x := m.Column(i)

		for j := i; j < n; j++ {
			y := m.Column(j)
			c := Covariance(x, y)

			result.SetValue(i, j, c)
			result.SetValue(j, i, c)
		}
	}

	return result
}

// Compute the QR decomposition of the matrix using the Gram-Schmidt process
func (m *Matrix) QR() (*Matrix, *Matrix) {
	rows := m.shape[0]
	cols := m.shape[1]

	q := NewMatrix(rows, cols)
	r := NewMatrix(cols, cols)

	for j := 0; j < cols; j++ {
		v := m.Column(j)

		for i := 0; i < j; i++ {
			rij := q.Column(i).Dot(m.Column(j))
			r.SetValue(i, j, rij)

			for k := 0; k < v.Size(); k++ {
				v[k] -= rij * q.At(k, i)
			}
		}

		rjj := v.Magnitude()
		r.SetValue(j, j, rjj)

		for k := 0; k < rows; k++ {
			q.SetValue(k, j, v[k]/rjj)
		}
	}

	return q, r
}
