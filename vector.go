package mtk

import (
	"errors"
	"math"
)

var (
	ErrVectorShapeMismatch = errors.New("vector shape mismatch")
)

// One-dimension vector
type Vector []float64

// Construct a Vector of size n
func NewVector(n int) Vector {
	return Vector(make([]float64, n))
}

// Get the value at the index
func (v Vector) At(i int) float64 {
	return v[i]
}

// Get the size
func (v Vector) Size() int {
	return len(v)
}

// Compute the L2-normalized Vector
func (v Vector) Normalize() Vector {
	magnitude := v.Magnitude()
	vn := make(Vector, len(v))

	for i := 0; i < len(v); i++ {
		vn[i] = v[i] / magnitude
	}

	return vn
}

// Compute the L2-norm
func (v Vector) Magnitude() float64 {
	return math.Sqrt(v.Dot(v))
}

// Compute the mean of values in the Vector
func (v Vector) Mean() float64 {
	var mean float64

	for _, value := range v {
		mean += value
	}

	return mean / float64(v.Size())
}

// Compute the Vector dot product u * v
func (v Vector) Dot(u Vector) float64 {
	if len(v) != len(u) {
		panic(ErrVectorShapeMismatch)
	}

	var d float64

	for i := 0; i < len(u); i++ {
		d += v[i] * u[i]
	}

	return d
}

// Compute the covariance of two Vectors
func Covariance(x, y Vector) float64 {
	if x.Size() != y.Size() {
		panic(ErrVectorShapeMismatch)
	}

	var value float64
	xMean := x.Mean()
	yMean := y.Mean()

	for i := 0; i < x.Size(); i++ {
		value += (x[i] - xMean) * (y[i] - yMean)
	}

	return value / (float64(x.Size() - 1))
}
