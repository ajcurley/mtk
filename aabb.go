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
