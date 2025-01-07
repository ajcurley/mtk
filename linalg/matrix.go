package linalg

import (
	"math"
	"slices"
)

// Two-dimensional matrix
type Matrix struct {
	shape [2]int
	data  []float64
}

// Construct a two-dimensional matrix
func NewMatrix(rows, cols int, data []float64) *Matrix {
	if rows <= 0 || cols <= 0 {
		panic("matrix dimensions must be non-zero")
	}

	if data == nil {
		data = make([]float64, rows*cols)
	}

	if rows*cols != len(data) {
		panic("matrix dimensions and data size must match")
	}

	return &Matrix{
		shape: [2]int{rows, cols},
		data:  data,
	}
}

// Construct an identity matrix of size (n, n)
func NewIdentityMatrix(size int) *Matrix {
	m := NewMatrix(size, size, nil)

	for i := 0; i < size; i++ {
		m.Set(i, i, 1)
	}

	return m
}

// Get the shape
func (m *Matrix) Shape() (int, int) {
	return m.shape[0], m.shape[1]
}

// Get the size
func (m *Matrix) Size() int {
	return len(m.data)
}

// Copy the matrix
func (m *Matrix) Copy() *Matrix {
	rows, cols := m.Shape()
	data := make([]float64, rows*cols)
	copy(data, m.data)
	return NewMatrix(rows, cols, data)
}

// Convert the (i, j) index to the flattened array index
func (m *Matrix) index(i, j int) int {
	rows, cols := m.Shape()

	if i < 0 || i >= rows || j < 0 || j > cols {
		panic("index out of range")
	}

	return i*cols + j
}

// Get the value at an index
func (m *Matrix) At(i, j int) float64 {
	return m.data[m.index(i, j)]
}

// Get the row values
func (m *Matrix) Row(i int) *Vector {
	_, cols := m.Shape()
	data := make([]float64, cols)

	for j := 0; j < cols; j++ {
		data[j] = m.At(i, j)
	}

	return NewVector(cols, data)
}

// Get the column values
func (m *Matrix) Col(i int) *Vector {
	rows, _ := m.Shape()
	data := make([]float64, rows)

	for j := 0; j < rows; j++ {
		data[j] = m.At(j, i)
	}

	return NewVector(rows, data)
}

// Set the value at (i, j)
func (m *Matrix) Set(i, j int, v float64) {
	m.data[m.index(i, j)] = v
}

// Fill the matrix with a value
func (m *Matrix) Fill(value float64) {
	for i := 0; i < len(m.data); i++ {
		m.data[i] = value
	}
}

// Check if the matrix is square
func (m *Matrix) IsSquare() bool {
	rows, cols := m.Shape()
	return rows == cols
}

// Check if the matrix is symmetric
func (m *Matrix) IsSymmetric() bool {
	rows, cols := m.Shape()

	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			if m.At(i, j) != m.At(j, i) {
				return false
			}
		}
	}

	return true
}

// Get the trace (diagonal) of the matrix
func (m *Matrix) Trace() *Vector {
	if !m.IsSquare() {
		panic("matrix must be square")
	}

	rows, _ := m.Shape()
	data := make([]float64, rows)

	for i := 0; i < rows; i++ {
		data[i] = m.At(i, i)
	}

	return NewVector(rows, data)
}

// Compute the transpose of the matrix
func (m *Matrix) Transpose() *Matrix {
	rows, cols := m.Shape()
	n := NewMatrix(rows, cols, nil)

	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			v := m.At(j, i)
			n.Set(i, j, v)
		}
	}

	return n
}

// Compute the matrix multiplication m * n
func (m *Matrix) Dot(n *Matrix) *Matrix {
	mRows, mCols := m.shape[0], m.shape[1]
	nRows, nCols := n.shape[0], n.shape[1]

	if mCols != nRows {
		panic("matrix inner dimensions must match")
	}

	d := NewMatrix(mRows, nCols, nil)

	for i := 0; i < mRows; i++ {
		for j := 0; j < nCols; j++ {
			var value float64

			for k := 0; k < nRows; k++ {
				value += m.At(i, k) * n.At(k, j)
			}

			d.Set(i, j, value)
		}
	}

	return d
}

// Compute the covariance matrix
func (m *Matrix) Covariance() *Matrix {
	rows, _ := m.Shape()
	result := NewMatrix(rows, rows, nil)

	for i := 0; i < rows; i++ {
		x := m.Col(i)

		for j := i; j < rows; j++ {
			y := m.Col(j)
			c := Covariance(x, y)

			result.Set(i, j, c)
			result.Set(j, i, c)
		}
	}

	return result
}

// Compute the eigenvalue and eigenvector pairs for a real, symmetric
// matrix using the QR iteration. The returned pairs are sorted in descending
// order of the eigenvalue magnitude.
func (m *Matrix) SymmetricEigen() (*Vector, *Matrix) {
	if !m.IsSquare() || !m.IsSymmetric() {
		panic("matrix must be square and symmetric")
	}

	rows, _ := m.Shape()
	a := m.Copy()
	qq := NewIdentityMatrix(rows)

	for i := 0; i < 100; i++ {
		q, r := a.QR()
		a = r.Dot(q)
		qq = qq.Dot(q)
	}

	index := make([]int, rows)
	trace := a.Trace()

	for i := 0; i < rows; i++ {
		index[i] = i
	}

	slices.SortStableFunc(index, func(i, j int) int {
		return -int(math.Abs(trace.At(i)) - math.Abs(trace.At(j)))
	})

	e := NewVector(rows, nil)
	v := NewMatrix(rows, rows, nil)

	for j, k := range index {
		e.Set(j, a.At(k, k))

		for i := 0; i < rows; i++ {
			v.Set(i, j, qq.At(i, k))
		}
	}

	return e, v
}

// Compute the QR decomposition of the matrix using the Gram-Schmidt process
func (m *Matrix) QR() (*Matrix, *Matrix) {
	rows, cols := m.Shape()
	q := NewMatrix(rows, cols, nil)
	r := NewMatrix(cols, cols, nil)

	for j := 0; j < cols; j++ {
		v := m.Col(j)

		for i := 0; i < j; i++ {
			rij := q.Col(i).Dot(m.Col(j))
			r.Set(i, j, rij)

			for k := 0; k < v.Size(); k++ {
				vk := v.At(k) - rij*q.At(k, i)
				v.Set(k, vk)
			}
		}

		rjj := v.Magnitude()
		r.Set(j, j, rjj)

		for k := 0; k < rows; k++ {
			q.Set(k, j, v.At(k)/rjj)
		}
	}

	return q, r
}

// Compute the orthogonal axes in descending order of their eigenvalue
// magnitude using principal component analysis (PCA).
func (m *Matrix) PCA() []*Vector {
	_, cols := m.Shape()
	axes := make([]*Vector, cols)
	_, eigenValues := m.Covariance().SymmetricEigen()

	for i := 0; i < cols; i++ {
		axes[i] = eigenValues.Col(i)
	}

	return axes
}
