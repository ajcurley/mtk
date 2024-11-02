package mtk

// Three dimension Cartesian ray
type Ray struct {
	Origin    Vector3
	Direction Vector3
}

// Construct a Ray from its origin and direction
func NewRay(origin, direction Vector3) Ray {
	return Ray{
		Origin:    origin,
		Direction: direction,
	}
}

// Check for an intersection with an AABB
func (r Ray) IntersectsAABB(a AABB) bool {
	inv := r.Direction.Inv()
	tx0 := (a.Min[0] - r.Origin[0]) * inv[0]
	tx1 := (a.Max[0] - r.Origin[0]) * inv[0]
	tMin := min(tx0, tx1)
	tMax := max(tx0, tx1)

	ty0 := (a.Min[1] - r.Origin[1]) * inv[1]
	ty1 := (a.Max[1] - r.Origin[1]) * inv[1]
	tMin = max(tMin, min(ty0, ty1))
	tMax = min(tMax, max(ty0, ty1))

	tz0 := (a.Min[2] - r.Origin[2]) * inv[2]
	tz1 := (a.Max[2] - r.Origin[2]) * inv[2]
	tMin = max(tMin, min(tz0, tz1))
	tMax = min(tMax, max(tz0, tz1))

	return tMax >= max(tMin, 0)
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
