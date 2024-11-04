package mtk

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test an AABB/Triangle intersection for a fully contained triangle
func TestTriangleIntersectsAABBHitContained(t *testing.T) {
	aabb := NewAABB(
		Vector3{0.5, 0.5, 0.5},
		Vector3{0.5, 0.5, 0.5},
	)
	triangle := NewTriangle(
		Vector3{0.1, 0.1, 0.1},
		Vector3{0.1, 0.1, 0.3},
		Vector3{0.1, 0.3, 0.1},
	)

	assert.True(t, triangle.IntersectsAABB(aabb))
}

// Test an AABB/Triangle intersection for a triangle crossing a face
func TestTriangleIntersectsAABBHitCrossFace(t *testing.T) {
	aabb := NewAABB(
		Vector3{0.5, 0.5, 0.5},
		Vector3{0.5, 0.5, 0.5},
	)
	triangle := NewTriangle(
		Vector3{0.5, 0.5, 0.5},
		Vector3{1.25, 1.75, 0.5},
		Vector3{1.25, 0.25, 0.5},
	)

	assert.True(t, triangle.IntersectsAABB(aabb))
}

// Test an AABB/Triangle intersection for a triangle crossing an edge
func TestTriangleIntersectsAABBHitCrossEdge(t *testing.T) {
	aabb := NewAABB(
		Vector3{0.5, 0.5, 0.5},
		Vector3{0.5, 0.5, 0.5},
	)
	triangle := NewTriangle(
		Vector3{0.25, -0.25, 0.5},
		Vector3{1.25, 0.75, 0.5},
		Vector3{1.25, -0.25, 0.5},
	)

	assert.True(t, triangle.IntersectsAABB(aabb))
}

// Test an AABB/Triangle intersection for a triangle crossing all four faces
func TestTriangleIntersectsAABBHitCrossFull(t *testing.T) {
	aabb := NewAABB(
		Vector3{0.5, 0.5, 0.5},
		Vector3{0.5, 0.5, 0.5},
	)
	triangle := NewTriangle(
		Vector3{-2, -1, 0.5},
		Vector3{1.5, 3, 0.5},
		Vector3{1.5, -1, 0.5},
	)

	assert.True(t, triangle.IntersectsAABB(aabb))
}

// Test an AABB/Triangle intersection miss bounds test
func TestTriangleIntersectsAABBMissAABB(t *testing.T) {
	aabb := NewAABB(
		Vector3{0.5, 0.5, 0.5},
		Vector3{0.5, 0.5, 0.5},
	)
	triangle := NewTriangle(
		Vector3{0, 0, 2},
		Vector3{1, 0, 2},
		Vector3{1, 1, 2},
	)

	assert.False(t, triangle.IntersectsAABB(aabb))
}

// Test an AABB/Triangle intersection miss plane test
func TestTriangleIntersectsAABBMissPlane(t *testing.T) {
	aabb := NewAABB(
		Vector3{0.5, 0.5, 0.5},
		Vector3{0.5, 0.5, 0.5},
	)
	triangle := NewTriangle(
		Vector3{1.1, 1.1, 0.9},
		Vector3{0.5, 0.8, 1.5},
		Vector3{0.9, 1.1, 0.9},
	)

	assert.False(t, triangle.IntersectsAABB(aabb))
}

// Test an AABB/Triangle intersection miss e0 x01 test
func TestTriangleIntersectsAABBMissE0X01(t *testing.T) {
	aabb := NewAABB(
		Vector3{0.5, 0.5, 0.5},
		Vector3{0.5, 0.5, 0.5},
	)
	triangle := NewTriangle(
		Vector3{0.5, 1.1, 0.9},
		Vector3{0.5, 0.8, 1.5},
		Vector3{0.5, 1.3, 1.2},
	)

	assert.False(t, triangle.IntersectsAABB(aabb))
}

// Test an AABB/Triangle intersection miss e0 y02 test
func TestTriangleIntersectsAABBMissE0Y02(t *testing.T) {
	aabb := NewAABB(
		Vector3{0.5, 0.5, 0.5},
		Vector3{0.5, 0.5, 0.5},
	)
	triangle := NewTriangle(
		Vector3{1.1, 0.5, 0.9},
		Vector3{0.8, 0.5, 1.5},
		Vector3{1.3, 0.5, 1.2},
	)

	assert.False(t, triangle.IntersectsAABB(aabb))
}

// Test an AABB/Triangle intersection miss e0 z12 test
func TestTriangleIntersectsAABBMissE0Z12(t *testing.T) {
	aabb := NewAABB(
		Vector3{0.5, 0.5, 0.5},
		Vector3{0.5, 0.5, 0.5},
	)
	triangle := NewTriangle(
		Vector3{1.1, 0.9, 0.5},
		Vector3{0.8, 1.5, 0.5},
		Vector3{1.3, 1.2, 0.5},
	)

	assert.False(t, triangle.IntersectsAABB(aabb))
}

// Test an AABB/Triangle intersection miss e1 x01 test
func TestTriangleIntersectsAABBMissE1X01(t *testing.T) {
	aabb := NewAABB(
		Vector3{0.5, 0.5, 0.5},
		Vector3{0.5, 0.5, 0.5},
	)
	triangle := NewTriangle(
		Vector3{0.5, 1.3, 1.2},
		Vector3{0.5, 1.1, 0.9},
		Vector3{0.5, 0.8, 1.5},
	)

	assert.False(t, triangle.IntersectsAABB(aabb))
}

// Test an AABB/Triangle intersection miss e1 y02 test
func TestTriangleIntersectsAABBMissE1Y02(t *testing.T) {
	aabb := NewAABB(
		Vector3{0.5, 0.5, 0.5},
		Vector3{0.5, 0.5, 0.5},
	)
	triangle := NewTriangle(
		Vector3{1.3, 0.5, 1.2},
		Vector3{1.1, 0.5, 0.9},
		Vector3{0.8, 0.5, 1.5},
	)

	assert.False(t, triangle.IntersectsAABB(aabb))
}

// Test an AABB/Triangle intersection miss e1 z0 test
func TestTriangleIntersectsAABBMissE1Z0(t *testing.T) {
	aabb := NewAABB(
		Vector3{0.5, 0.5, 0.5},
		Vector3{0.5, 0.5, 0.5},
	)
	triangle := NewTriangle(
		Vector3{1.3, 1.2, 0.5},
		Vector3{1.1, 0.9, 0.5},
		Vector3{0.8, 1.5, 0.5},
	)

	assert.False(t, triangle.IntersectsAABB(aabb))
}

// Test an AABB/Triangle intersection miss e2 x2 test
func TestTriangleIntersectsAABBMissE2X2(t *testing.T) {
	aabb := NewAABB(
		Vector3{0.5, 0.5, 0.5},
		Vector3{0.5, 0.5, 0.5},
	)
	triangle := NewTriangle(
		Vector3{0.5, 0.8, 1.5},
		Vector3{0.5, 1.3, 1.2},
		Vector3{0.5, 1.1, 0.9},
	)

	assert.False(t, triangle.IntersectsAABB(aabb))
}

// Test an AABB/Triangle intersection miss e2 y1 test
func TestTriangleIntersectsAABBMissE2Y1(t *testing.T) {
	aabb := NewAABB(
		Vector3{0.5, 0.5, 0.5},
		Vector3{0.5, 0.5, 0.5},
	)
	triangle := NewTriangle(
		Vector3{0.8, 0.5, 1.5},
		Vector3{1.3, 0.5, 1.2},
		Vector3{1.1, 0.5, 0.9},
	)

	assert.False(t, triangle.IntersectsAABB(aabb))
}

// Test an AABB/Triangle intersection miss e2 z12 test
func TestTriangleIntersectsAABBMissE2Z12(t *testing.T) {
	aabb := NewAABB(
		Vector3{0.5, 0.5, 0.5},
		Vector3{0.5, 0.5, 0.5},
	)
	triangle := NewTriangle(
		Vector3{0.8, 1.5, 0.5},
		Vector3{1.3, 1.2, 0.5},
		Vector3{1.1, 0.9, 0.5},
	)

	assert.False(t, triangle.IntersectsAABB(aabb))
}
