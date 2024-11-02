package mtk

import (
	"math"
)

// Three dimension Cartesian ray
type Ray struct {
	Origin    Vector3
	Direction Vector3
}

// Check for an intersection with an AABB
func (r Ray) IntersectsAABB(a AABB) bool {
	inv := r.Direction.Inv()
	tMin := math.Inf(1)
	tMax := math.Inf(-1)

	for i := 0; i < 3; i++ {
		t1 := (a.Min[i] - r.Origin[i]) * inv[i]
		t2 := (a.Max[i] - r.Origin[i]) * inv[i]
		tMin = max(tMin, min(t1, t2))
		tMax = min(tMax, max(t1, t2))
	}

	return tMax >= max(tMin, 0.)
}

// Check for an intersection with a Triangle
func (r Ray) IntersectsTriangle(t Triangle) bool {
	e0 := t[1].Sub(t[0])
	e1 := t[2].Sub(t[0])

	p := r.Direction.Cross(e1)
	d := e0.Dot(p)

	if d < GeometricTolerance {
		return false
	}

	dInv := 1. / d
	s := r.Origin.Sub(t[0])
	u := dInv * s.Dot(p)

	if u < 0. || u > 1. {
		return false
	}

	q := s.Cross(e0)
	v := dInv * r.Direction.Dot(q)

	if v < 0. || u+v > 1. {
		return false
	}

	return (dInv * e1.Dot(q)) > GeometricTolerance
}
