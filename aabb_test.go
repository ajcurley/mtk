package mtk

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAABBOctant(t *testing.T) {
	a := NewAABB(Vector3{2, 2, 2}, Vector3{2, 2, 2})

	octant0 := a.Octant(0)
	octant1 := a.Octant(1)
	octant2 := a.Octant(2)
	octant3 := a.Octant(3)
	octant4 := a.Octant(4)
	octant5 := a.Octant(5)
	octant6 := a.Octant(6)
	octant7 := a.Octant(7)

	assert.Equal(t, Vector3{1, 1, 1}, octant0.Center)
	assert.Equal(t, Vector3{1, 1, 3}, octant1.Center)
	assert.Equal(t, Vector3{1, 3, 1}, octant2.Center)
	assert.Equal(t, Vector3{1, 3, 3}, octant3.Center)
	assert.Equal(t, Vector3{3, 1, 1}, octant4.Center)
	assert.Equal(t, Vector3{3, 1, 3}, octant5.Center)
	assert.Equal(t, Vector3{3, 3, 1}, octant6.Center)
	assert.Equal(t, Vector3{3, 3, 3}, octant7.Center)
}

// Test an AABB/AABB intersection hit with overlap
func TestAABBIntersectsAABBHitOverlap(t *testing.T) {
	a := NewAABB(Vector3{1, 1, 1}, Vector3{1, 1, 1})
	b := NewAABB(Vector3{2, 2, 2}, Vector3{1, 1, 1})

	assert.True(t, a.IntersectsAABB(b))
	assert.True(t, b.IntersectsAABB(a))
}

// Test an AABB/AABB intersection hit with full containment
func TestAABBIntersectsAABBHitContained(t *testing.T) {
	a := NewAABB(Vector3{1.5, 1.5, 1.5}, Vector3{1.5, 1.5, 1.5})
	b := NewAABB(Vector3{1.5, 1.5, 1.5}, Vector3{0.5, 0.5, 0.5})

	assert.True(t, a.IntersectsAABB(b))
	assert.True(t, b.IntersectsAABB(a))
}

// Test an AABB/AABB intersect miss
func TestAABBIntersectsAABBMiss(t *testing.T) {
	a := NewAABB(Vector3{0.5, 0.5, 0.5}, Vector3{0.5, 0.5, 0.5})
	b := NewAABB(Vector3{2.5, 2.5, 2.5}, Vector3{0.5, 0.5, 0.5})

	assert.False(t, a.IntersectsAABB(b))
	assert.False(t, b.IntersectsAABB(a))
}
