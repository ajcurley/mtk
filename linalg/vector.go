package linalg

import (
	"errors"
	"math"
)

var (
	ErrVectorShapeMismatch = errors.New("vector shape mismatch")
)

// One-dimension vector
type Vector struct {
	data []float64
}

// Construct a Vector of size n
func NewVector(n int, data []float64) *Vector {
	if data == nil {
		data = make([]float64, n)
	}

	if n != len(data) {
		panic("vector length and data length must match")
	}

	return &Vector{data: data}
}

// Copy into a new vector
func (v *Vector) Copy() *Vector {
	data := make([]float64, v.Size())
	copy(data, v.data)
	return NewVector(v.Size(), data)
}

// Get the size
func (v *Vector) Size() int {
	return len(v.data)
}

// Get the value at i
func (v *Vector) At(i int) float64 {
	return v.data[i]
}

// Set the value at i
func (v *Vector) Set(i int, value float64) {
	v.data[i] = value
}

// Compute the L2-normalized Vector
func (v *Vector) Unit() *Vector {
	return v.MulScalar(1. / v.Magnitude())
}

// Compute the L2-norm
func (v *Vector) Magnitude() float64 {
	return math.Sqrt(v.Dot(v))
}

// Compute the mean of values in the Vector
func (v *Vector) Mean() float64 {
	var mean float64

	for i := 0; i < v.Size(); i++ {
		mean += v.At(i)
	}

	return mean / float64(v.Size())
}

// Element-wise addition
func (v *Vector) Add(w *Vector) *Vector {
	if v.Size() != w.Size() {
		panic("vector lengths must match")
	}

	u := v.Copy()

	for i := 0; i < u.Size(); i++ {
		u.Set(i, v.At(i)+w.At(i))
	}

	return u
}

// Scalar addition
func (v *Vector) AddScalar(s float64) *Vector {
	u := v.Copy()

	for i := 0; i < u.Size(); i++ {
		u.Set(i, v.At(i)+s)
	}

	return u
}

// Element-wise subtraction
func (v *Vector) Sub(w *Vector) *Vector {
	if v.Size() != w.Size() {
		panic("vector lengths must match")
	}

	u := v.Copy()

	for i := 0; i < u.Size(); i++ {
		u.Set(i, v.At(i)-w.At(i))
	}

	return u
}

// Scalar subtraction
func (v *Vector) SubScalar(s float64) *Vector {
	u := v.Copy()

	for i := 0; i < u.Size(); i++ {
		u.Set(i, v.At(i)-s)
	}

	return u
}

// Element-wise multiplication
func (v *Vector) Mul(w *Vector) *Vector {
	if v.Size() != w.Size() {
		panic("vector lengths must match")
	}

	u := v.Copy()

	for i := 0; i < u.Size(); i++ {
		u.Set(i, v.At(i)*w.At(i))
	}

	return u
}

// Scalar multiplication
func (v *Vector) MulScalar(s float64) *Vector {
	u := v.Copy()

	for i := 0; i < u.Size(); i++ {
		u.Set(i, v.At(i)*s)
	}

	return u
}

// Element-wise division
func (v *Vector) Div(w *Vector) *Vector {
	if v.Size() != w.Size() {
		panic("vector lengths must match")
	}

	u := v.Copy()

	for i := 0; i < u.Size(); i++ {
		u.Set(i, v.At(i)/w.At(i))
	}

	return u
}

// Scalar division
func (v *Vector) DivScalar(s float64) *Vector {
	u := v.Copy()

	for i := 0; i < u.Size(); i++ {
		u.Set(i, v.At(i)/s)
	}

	return u
}

// Compute the Vector dot product v * w
func (v *Vector) Dot(w *Vector) float64 {
	if v.Size() != w.Size() {
		panic("vector lengths must match")
	}

	var d float64

	for i := 0; i < v.Size(); i++ {
		d += v.At(i) * w.At(i)
	}

	return d
}

// Compute the covariance of two Vectors
func Covariance(x, y *Vector) float64 {
	if x.Size() != y.Size() {
		panic("vector lengths must match")
	}

	var value float64
	xMean := x.Mean()
	yMean := y.Mean()

	for i := 0; i < x.Size(); i++ {
		value += (x.At(i) - xMean) * (y.At(i) - yMean)
	}

	return value / (float64(x.Size() - 1))
}
