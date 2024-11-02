package mtk

// Three dimensional Cartesian axis-aligned bounding box
type AABB struct {
	Min Vector3
	Max Vector3
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

// Check for an intersection with an AABB
func (a AABB) IntersectsAABB(b AABB) bool {
	return a.Min[0] <= b.Max[0] &&
		a.Max[0] >= b.Min[0] &&
		a.Min[1] <= b.Max[1] &&
		a.Max[1] >= b.Min[1] &&
		a.Min[2] <= b.Max[2] &&
		a.Max[2] >= b.Min[2]
}
