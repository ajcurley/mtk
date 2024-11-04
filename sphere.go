package mtk

// Three-dimensional Cartesian sphere
type Sphere struct {
	Center Vector3
	Radius float64
}

// Construct a Sphere from its center and radius
func NewSphere(center Vector3, radius float64) Sphere {
	return Sphere{Center: center, Radius: radius}
}

// Check for an intersection with an AABB
func (s Sphere) IntersectsAABB(a AABB) bool {
	var d float64
	minBound := a.Min()
	maxBound := a.Max()

	for i := 0; i < 3; i++ {
		if s.Center[i] < minBound[i] {
			t := s.Center[i] - minBound[i]
			d += t * t
		} else if s.Center[i] > maxBound[i] {
			t := s.Center[i] - maxBound[i]
			d += t * t
		}
	}

	return d <= s.Radius*s.Radius
}

// Check for an intersection with a Vector
func (s Sphere) IntersectsVector3(v Vector3) bool {
	t := v.Sub(s.Center)
	return t.Dot(t) <= s.Radius*s.Radius
}
