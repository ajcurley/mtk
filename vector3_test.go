package mtk

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test calculating the Vector3 magnitude
func TestVector3Mag(t *testing.T) {
	v := Vector3{-3, 0, 4}
	assert.Equal(t, v.Mag(), 5.)
}

// Test calculating the Vector3 unit vector
func TestVector3Unit(t *testing.T) {
	v := Vector3{-3, 0, 4}
	u := v.Unit()
	assert.Equal(t, u.X(), -3./5.)
	assert.Equal(t, u.Y(), 0.)
	assert.Equal(t, u.Z(), 4./5.)
}

// Test adding two vectors
func TestVector3Add(t *testing.T) {
	v := Vector3{-3, 0, 4}
	u := Vector3{21, 9, 7}
	w := v.Add(u)
	assert.Equal(t, w.X(), 18.)
	assert.Equal(t, w.Y(), 9.)
	assert.Equal(t, w.Z(), 11.)
}

// Test adding a vector and scalar
func TestVector3AddScalar(t *testing.T) {
	v := Vector3{-3, 0, 4}
	w := v.AddScalar(6)
	assert.Equal(t, w.X(), 3.)
	assert.Equal(t, w.Y(), 6.)
	assert.Equal(t, w.Z(), 10.)
}

// Test subtracting two vectors
func TestVector3Sub(t *testing.T) {
	v := Vector3{-3, 0, 4}
	u := Vector3{21, 9, 7}
	w := v.Sub(u)
	assert.Equal(t, w.X(), -24.)
	assert.Equal(t, w.Y(), -9.)
	assert.Equal(t, w.Z(), -3.)
}

// Test subtracting a vector and scalar
func TestVector3SubScalar(t *testing.T) {
	v := Vector3{-3, 0, 4}
	w := v.SubScalar(6)
	assert.Equal(t, w.X(), -9.)
	assert.Equal(t, w.Y(), -6.)
	assert.Equal(t, w.Z(), -2.)
}

// Test multiplying two vectors
func TestVector3Mul(t *testing.T) {
	v := Vector3{-3, 0, 4}
	u := Vector3{21, 9, 7}
	w := v.Mul(u)
	assert.Equal(t, w.X(), -63.)
	assert.Equal(t, w.Y(), 0.)
	assert.Equal(t, w.Z(), 28.)
}

// Test multiplying a vector and scalar
func TestVector3MulScalar(t *testing.T) {
	v := Vector3{-3, 0, 4}
	w := v.MulScalar(6)
	assert.Equal(t, w.X(), -18.)
	assert.Equal(t, w.Y(), 0.)
	assert.Equal(t, w.Z(), 24.)
}

// Test dividing two vectors
func TestVector3Div(t *testing.T) {
	v := Vector3{-3, 0, 4}
	u := Vector3{21, 9, 7}
	w := v.Div(u)
	assert.Equal(t, w.X(), -3./21.)
	assert.Equal(t, w.Y(), 0.)
	assert.Equal(t, w.Z(), 4./7.)
}

// Test dividing two vectors with a zero in the denominator
func TestVector3DivZero(t *testing.T) {
	v := Vector3{-3, 0, 4}
	u := Vector3{0, 1, 2}
	w := v.Div(u)
	assert.True(t, math.IsInf(w.X(), -1))
	assert.Equal(t, w.Y(), 0.)
	assert.Equal(t, w.Z(), 2.)
}

// Test dividing a vector and scalar
func TestVector3DivScalar(t *testing.T) {
	v := Vector3{-3, 0, 4}
	w := v.DivScalar(6)
	assert.Equal(t, w.X(), -1./2.)
	assert.Equal(t, w.Y(), 0.)
	assert.Equal(t, w.Z(), 2./3.)
}

// Test dividing a vector and a zero scalar
func TestVector3DivScalarZero(t *testing.T) {
	v := Vector3{-3, 0, 4}
	w := v.DivScalar(0)
	assert.True(t, math.IsInf(w.X(), -1))
	assert.True(t, math.IsNaN(w.Y()))
	assert.True(t, math.IsInf(w.Z(), 1))
}

// Test calculating the vector dot product
func TestVector3Dot(t *testing.T) {
	v := Vector3{-3, 0, 4}
	u := Vector3{9, 1, 6}
	d := v.Dot(u)
	assert.Equal(t, d, -3.)
}

// Test calculating the vector cross product
func TestVector3Cross(t *testing.T) {
	v := Vector3{-3, 8, 4}
	u := Vector3{9, 7, 2}
	w := v.Cross(u)
	assert.Equal(t, w.X(), -12.)
	assert.Equal(t, w.Y(), 42.)
	assert.Equal(t, w.Z(), -93.)
}
