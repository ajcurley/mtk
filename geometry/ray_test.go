package geometry

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test a Ray/AABB intersection hit
func TestRayAABBIntersectsHit(t *testing.T) {
	ray := Ray{
		Origin:    Vector3{-1, 0.5, 0.5},
		Direction: Vector3{1, 0, 0},
	}

	aabb := NewAABB(
		Vector3{0.5, 0.5, 0.5},
		Vector3{0.5, 0.5, 0.5},
	)

	assert.True(t, ray.IntersectsAABB(aabb))
}

// Test a Ray/AABB intersection miss
func TestRayAABBIntersectsMiss(t *testing.T) {
	ray := Ray{
		Origin:    Vector3{-1, 2, 2},
		Direction: Vector3{1, 0, 0},
	}

	aabb := NewAABB(
		Vector3{0.5, 0.5, 0.5},
		Vector3{0.5, 0.5, 0.5},
	)

	assert.False(t, ray.IntersectsAABB(aabb))
}

// Test a Ray/Triangle intersection hit
func TestRayTriangleIntersectsHit(t *testing.T) {
	ray := Ray{
		Origin:    Vector3{0.5, 0.5, 0},
		Direction: Vector3{0, 0, 1},
	}

	triangle := NewTriangle(
		Vector3{0, 0, 1},
		Vector3{0, 1, 1},
		Vector3{1, 0, 1},
	)

	assert.True(t, ray.IntersectsTriangle(triangle))
}

// Test a Ray/Triangle intersection miss due to back-face culling
func TestRayTriangleIntersectsCulled(t *testing.T) {
	ray := Ray{
		Origin:    Vector3{0.5, 0.5, 0},
		Direction: Vector3{0, 0, 1},
	}

	triangle := NewTriangle(
		Vector3{0, 0, 1},
		Vector3{1, 0, 1},
		Vector3{0, 1, 1},
	)

	assert.False(t, ray.IntersectsTriangle(triangle))
}

// Test a Ray/Triangle intersection miss
func TestRayTriangleIntersectsMiss(t *testing.T) {
	ray := Ray{
		Origin:    Vector3{2, 2, 0},
		Direction: Vector3{0, 0, 1},
	}

	triangle := NewTriangle(
		Vector3{0, 0, 1},
		Vector3{0, 1, 1},
		Vector3{1, 0, 1},
	)

	assert.False(t, ray.IntersectsTriangle(triangle))
}
