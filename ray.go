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
