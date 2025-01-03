package geometry

const (
	GeometricTolerance float64 = 1e-8
)

// Interface for an AABB intersection test
type IntersectsAABB interface {
	IntersectsAABB(AABB) bool
}

// Interface for a Ray intersection test
type IntersectsRay interface {
	IntersectsRay(Ray) bool
}

// Interface for a Sphere intersection test
type IntersectsSphere interface {
	IntersectsSphere(Sphere) bool
}

// Interface for a Triangle intersection test
type IntersectsTriangle interface {
	IntersectsTriangle(Triangle) bool
}

// Interface for a Vector3 intersection test
type IntersectsVector3 interface {
	IntersectsVector3(Vector3) bool
}
