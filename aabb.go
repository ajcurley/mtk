package mtk

// Three dimensional Cartesian axis-aligned bounding box
type AABB struct {
	Min Vector3
	Max Vector3
}

// Construct an AABB from its min/max bounds
func NewAABB(minBound, maxBound Vector3) AABB {
	return AABB{Min: minBound, Max: maxBound}
}

// Get the center
func (a AABB) Center() Vector3 {
	return a.Max.Add(a.Min).MulScalar(0.5)
}

// Get the size vector
func (a AABB) Size() Vector3 {
	return a.Max.Sub(a.Min)
}

// Get the halfsize vector
func (a AABB) HalfSize() Vector3 {
	return a.Size().MulScalar(0.5)
}

// Get the AABB representing the octant
func (a AABB) Octant(octant int) AABB {
	halfSize := a.HalfSize()
	minBound := a.Min

	if octant&4 != 0 {
		minBound[0] += halfSize[0]
	}

	if octant&2 != 0 {
		minBound[1] += halfSize[1]
	}

	if octant&1 != 0 {
		minBound[2] += halfSize[2]
	}

	return AABB{
		Min: minBound,
		Max: minBound.Add(halfSize),
	}
}

// Check for an intersection with an AABB
func (a AABB) IntersectsAABB(b AABB) bool {
	return a.Min[0] <= b.Max[0] &&
		a.Max[0] >= b.Min[0] &&
		a.Min[1] <= b.Max[1] &&
		a.Max[1] >= b.Min[1] &&
		a.Min[2] <= b.Max[2] &&
		a.Max[2] >= b.Min[2]
}
