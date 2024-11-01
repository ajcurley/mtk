package mtk

import (
	"math"
)

// Three-dimensional Cartesian vector
type Vector3 [3]float64

// Get the x-component
func (v Vector3) X() float64 {
	return v[0]
}

// Get the y-component
func (v Vector3) Y() float64 {
	return v[1]
}

// Get the z-component
func (v Vector3) Z() float64 {
	return v[2]
}

// Get the magnitude (L2-norm)
func (v Vector3) Mag() float64 {
	return math.Sqrt(v.Dot(v))
}

// Get the unit vector
func (v Vector3) Unit() Vector3 {
	return v.DivScalar(v.Mag())
}

// Elementwise vector addition v + u
func (v Vector3) Add(u Vector3) Vector3 {
	return Vector3{
		v[0] + u[0],
		v[1] + u[1],
		v[2] + u[2],
	}
}

// Elementwise vector/scalar addition v + s
func (v Vector3) AddScalar(s float64) Vector3 {
	return Vector3{
		v[0] + s,
		v[1] + s,
		v[2] + s,
	}
}

// Elementwise vector subtraction v - u
func (v Vector3) Sub(u Vector3) Vector3 {
	return Vector3{
		v[0] - u[0],
		v[1] - u[1],
		v[2] - u[2],
	}
}

// Elementwise vector/scalar subtraction v - s
func (v Vector3) SubScalar(s float64) Vector3 {
	return Vector3{
		v[0] - s,
		v[1] - s,
		v[2] - s,
	}
}

// Elementwise vector multiplication v * u
func (v Vector3) Mul(u Vector3) Vector3 {
	return Vector3{
		v[0] * u[0],
		v[1] * u[1],
		v[2] * u[2],
	}
}

// Elementwise vector/scalar multiplication v * s
func (v Vector3) MulScalar(s float64) Vector3 {
	return Vector3{
		v[0] * s,
		v[1] * s,
		v[2] * s,
	}
}

// Elementwise vector division v / u
func (v Vector3) Div(u Vector3) Vector3 {
	return Vector3{
		v[0] / u[0],
		v[1] / u[1],
		v[2] / u[2],
	}
}

// Elementwise vector/scalar division
func (v Vector3) DivScalar(s float64) Vector3 {
	return Vector3{
		v[0] / s,
		v[1] / s,
		v[2] / s,
	}
}

// Get the dot product v * u
func (v Vector3) Dot(u Vector3) float64 {
	return u[0]*v[0] + u[1]*v[1] + u[2]*v[2]
}

// Get the cross product v x
func (v Vector3) Cross(u Vector3) Vector3 {
	return Vector3{
		v[1]*u[2] - v[2]*u[1],
		v[2]*u[0] - v[0]*u[2],
		v[0]*u[1] - v[1]*u[0],
	}
}

// Get the angle (in radians) between the vectors
func (v Vector3) AngleTo(u Vector3) float64 {
	arg := v.Dot(u) / (v.Mag() * u.Mag())
	arg = min(max(arg, -1), 1)
	return math.Acos(arg)
}
