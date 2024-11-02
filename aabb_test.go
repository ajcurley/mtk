package mtk

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test an AABB/AABB intersection hit with overlap
func TestAABBIntersectsAABBHitOverlap(t *testing.T) {
	a := AABB{Min: Vector3{0, 0, 0}, Max: Vector3{2, 2, 2}}
	b := AABB{Min: Vector3{1, 1, 1}, Max: Vector3{3, 3, 3}}

	assert.True(t, a.IntersectsAABB(b))
	assert.True(t, b.IntersectsAABB(a))
}

// Test an AABB/AABB intersection hit with full containment
func TestAABBIntersectsAABBHitContained(t *testing.T) {
	a := AABB{Min: Vector3{0, 0, 0}, Max: Vector3{3, 3, 3}}
	b := AABB{Min: Vector3{1, 1, 1}, Max: Vector3{2, 2, 2}}

	assert.True(t, a.IntersectsAABB(b))
	assert.True(t, b.IntersectsAABB(a))
}

// Test an AABB/AABB intersect miss
func TestAABBIntersectsAABBMiss(t *testing.T) {
	a := AABB{Min: Vector3{0, 0, 0}, Max: Vector3{1, 1, 1}}
	b := AABB{Min: Vector3{2, 2, 2}, Max: Vector3{3, 3, 3}}

	assert.False(t, a.IntersectsAABB(b))
	assert.False(t, b.IntersectsAABB(a))
}
