package mtk

import (
	"math"
)

// Three dimensional Cartesian triangle
// type Triangle [3]Vector3
type Triangle struct {
	P  Vector3
	Q  Vector3
	R  Vector3
	E0 Vector3
	E1 Vector3
	E2 Vector3
}

// Construct a Triangle from its points
func NewTriangle(p, q, r Vector3) Triangle {
	return Triangle{
		P:  p,
		Q:  q,
		R:  r,
		E0: q.Sub(p),
		E1: r.Sub(q),
		E2: p.Sub(r),
	}
}

// Get the normal vector (not necessarily a unit vector)
func (t Triangle) Normal() Vector3 {
	pq := t.Q.Sub(t.P)
	pr := t.R.Sub(t.P)
	return pq.Cross(pr)
}

// Get the unit normal vector
func (t Triangle) UnitNormal() Vector3 {
	return t.Normal().Unit()
}

// Get the area
func (t Triangle) Area() float64 {
	return 0.5 * t.Normal().Mag()
}

// Check for an intersection with a Ray
func (t Triangle) IntersectsRay(r Ray) bool {
	return r.IntersectsTriangle(t)
}

// Check for an intersection with an AABB
func (t Triangle) IntersectsAABB(a AABB) bool {
	// Shift the system such that the AABB is centered at the origin
	v0 := t.P.Sub(a.Center)
	v1 := t.Q.Sub(a.Center)
	v2 := t.R.Sub(a.Center)

	// Compute the triangle edges
	e0 := t.E0
	e1 := t.E1
	e2 := t.E2

	// Bullet #1 - Test the AABB against the minimum AABB of the triangle
	for i := 0; i < 3; i++ {
		vMin := min(v0[i], v1[i], v2[i])
		vMax := max(v0[i], v1[i], v2[i])

		if vMin > a.HalfSize[i] || vMax < -a.HalfSize[i] {
			return false
		}
	}

	// Bullet #2 - Test the triangle plane against the AABB
	normal := e0.Cross(e1)

	if !planeBoxOverlap(normal, v0, a.HalfSize) {
		return false
	}

	// Bullet #3 - 9 tests
	var fex, fey, fez float64
	fex = math.Abs(e0[0])
	fey = math.Abs(e0[1])
	fez = math.Abs(e0[2])

	if !axisTestX01(e0[2], e0[1], fez, fey, v0, v2, a.HalfSize) {
		return false
	}

	if !axisTestY02(e0[2], e0[0], fez, fex, v0, v2, a.HalfSize) {
		return false
	}

	if !axisTestZ12(e0[1], e0[0], fey, fex, v1, v2, a.HalfSize) {
		return false
	}

	fex = math.Abs(e1[0])
	fey = math.Abs(e1[1])
	fez = math.Abs(e1[2])

	if !axisTestX01(e1[2], e1[1], fez, fey, v0, v2, a.HalfSize) {
		return false
	}

	if !axisTestY02(e1[2], e1[0], fez, fex, v0, v2, a.HalfSize) {
		return false
	}

	if !axisTestZ0(e1[1], e1[0], fey, fex, v0, v1, a.HalfSize) {
		return false
	}

	fex = math.Abs(e2[0])
	fey = math.Abs(e2[1])
	fez = math.Abs(e2[2])

	if !axisTestX2(e2[2], e2[1], fez, fey, v0, v1, a.HalfSize) {
		return false
	}

	if !axisTestY1(e2[2], e2[0], fez, fex, v0, v1, a.HalfSize) {
		return false
	}

	if !axisTestZ12(e2[1], e2[0], fey, fex, v1, v2, a.HalfSize) {
		return false
	}

	return true
}

func axisTestX01(a, b, fa, fb float64, v0, v2, h Vector3) bool {
	p0 := a*v0[1] - b*v0[2]
	p2 := a*v2[1] - b*v2[2]
	pMin := min(p0, p2)
	pMax := max(p0, p2)
	rad := fa*h[1] + fb*h[2]
	return !(pMin > rad || pMax < -rad)
}

func axisTestX2(a, b, fa, fb float64, v0, v1, h Vector3) bool {
	p0 := a*v0[1] - b*v0[2]
	p1 := a*v1[1] - b*v1[2]
	pMin := min(p0, p1)
	pMax := max(p0, p1)
	rad := fa*h[1] + fb*h[2]
	return !(pMin > rad || pMax < -rad)
}

func axisTestY02(a, b, fa, fb float64, v0, v2, h Vector3) bool {
	p0 := -a*v0[0] + b*v0[2]
	p2 := -a*v2[0] + b*v2[2]
	pMin := min(p0, p2)
	pMax := max(p0, p2)
	rad := fa*h[0] + fb*h[2]
	return !(pMin > rad || pMax < -rad)
}

func axisTestY1(a, b, fa, fb float64, v0, v1, h Vector3) bool {
	p0 := -a*v0[0] + b*v0[2]
	p1 := -a*v1[0] + b*v1[2]
	pMin := min(p0, p1)
	pMax := max(p0, p1)
	rad := fa*h[0] + fb*h[2]
	return !(pMin > rad || pMax < -rad)
}

func axisTestZ12(a, b, fa, fb float64, v1, v2, h Vector3) bool {
	p1 := a*v1[0] - b*v1[1]
	p2 := a*v2[0] - b*v2[1]
	pMin := min(p1, p2)
	pMax := max(p1, p2)
	rad := fa*h[0] + fb*h[1]
	return !(pMin > rad || pMax < -rad)
}

func axisTestZ0(a, b, fa, fb float64, v0, v1, h Vector3) bool {
	p0 := a*v0[0] - b*v0[1]
	p1 := a*v1[0] - b*v1[1]
	pMin := min(p0, p1)
	pMax := max(p0, p1)
	rad := fa*h[0] + fb*h[1]
	return !(pMin > rad || pMax < -rad)
}

func planeBoxOverlap(normal, vertex, boxMax Vector3) bool {
	vMin := Vector3{0, 0, 0}
	vMax := Vector3{0, 0, 0}

	for i := 0; i < 3; i++ {
		if normal[i] > 0 {
			vMin[i] = -boxMax[i] - vertex[i]
			vMax[i] = boxMax[i] - vertex[i]
		} else {
			vMin[i] = boxMax[i] - vertex[i]
			vMax[i] = -boxMax[i] - vertex[i]
		}
	}

	return normal.Dot(vMin) <= 0 && normal.Dot(vMax) >= 0
}
