package mtk

// Three dimensional Cartesian axis-aligned bounding box
type AABB struct {
	Center   Vector3
	HalfSize Vector3
}

// Construct an AABB from its min/max bounds
func NewAABB(center, halfSize Vector3) AABB {
	return AABB{Center: center, HalfSize: halfSize}
}

// Get the min bounds
func (a AABB) Min() Vector3 {
	return a.Center.Sub(a.HalfSize)
}

// Get the max bound
func (a AABB) Max() Vector3 {
	return a.Center.Add(a.HalfSize)
}

// Get the buffered AABB
func (a AABB) Buffer(r float64) AABB {
	halfSize := a.HalfSize.AddScalar(r)
	return NewAABB(a.Center, halfSize)
}

// Get the AABB representing the octant
func (a AABB) Octant(octant int) AABB {
	center := a.Center
	halfSize := a.HalfSize.MulScalar(0.5)

	if octant&4 != 0 {
		center[0] += halfSize[0]
	} else {
		center[0] -= halfSize[0]
	}

	if octant&2 != 0 {
		center[1] += halfSize[1]
	} else {
		center[1] -= halfSize[1]
	}

	if octant&1 != 0 {
		center[2] += halfSize[2]
	} else {
		center[2] -= halfSize[2]
	}

	return NewAABB(center, halfSize)
}

// Check for an intersection with an AABB
func (a AABB) IntersectsAABB(b AABB) bool {
	return a.Center[0]-a.HalfSize[0] <= b.Center[0]+b.HalfSize[0] &&
		a.Center[0]+a.HalfSize[0] >= b.Center[0]-b.HalfSize[0] &&
		a.Center[1]-a.HalfSize[1] <= b.Center[1]+b.HalfSize[1] &&
		a.Center[1]+a.HalfSize[1] >= b.Center[1]-b.HalfSize[1] &&
		a.Center[2]-a.HalfSize[2] <= b.Center[2]+b.HalfSize[2] &&
		a.Center[2]+a.HalfSize[2] >= b.Center[2]-b.HalfSize[2]
}

// Check for an intersection with a Sphere
func (a AABB) IntersectsSphere(s Sphere) bool {
	return s.IntersectsAABB(a)
}
