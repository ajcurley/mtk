package linalg

import (
	"math"

	"github.com/ajcurley/mtk/geometry"
)

// Compute the rotation matrix from u to v using the Rodriguez rotation
// matrix forumla.
func RotationMatrix(u, v geometry.Vector3) *Matrix {
	un := u.Unit()
	vn := v.Unit()

	wn := un.Cross(vn).Unit()
	angle := un.AngleTo(vn)

	cos := math.Cos(angle)
	sin := math.Sin(angle)
	c := 1 - cos

	data := []float64{
		wn.X()*wn.Y()*c + cos,
		wn.Y()*wn.Y()*c - wn.Z()*sin,
		wn.X()*wn.Y()*c + wn.Y()*sin,
		wn.Y()*wn.X()*c + wn.Z()*sin,
		wn.Y()*wn.Y()*c + cos,
		wn.Y()*wn.Z()*c - wn.X()*sin,
		wn.Z()*wn.X()*c - wn.Y()*sin,
		wn.Z()*wn.Y()*c + wn.X()*sin,
		wn.Z()*wn.Z()*c + cos,
	}

	return NewMatrix(3, 3, data)
}
