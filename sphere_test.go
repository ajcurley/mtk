package mtk

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test a Sphere/AABB intersection for a sphere contained inside an AABB
func TestSphereIntersectsAABBHitContained(t *testing.T) {
	s := NewSphere(Vector3{1, 1, 1}, 1)
	a := NewAABB(Vector3{1, 1, 1}, Vector3{2, 2, 2})

	assert.True(t, s.IntersectsAABB(a))
}

// Test a Sphere/AABB intersection hit for an overlapping sphere
func TestSphereIntersectsAABBHitOverlap(t *testing.T) {
	s := NewSphere(Vector3{3, 1, 1}, 2)
	a := NewAABB(Vector3{1, 1, 1}, Vector3{1, 1, 1})

	assert.True(t, s.IntersectsAABB(a))
}

// Test a Sphere/AABB intersection miss
func TestSphereIntersectsAABBMiss(t *testing.T) {
	s := NewSphere(Vector3{0, 0, 0}, 1)
	a := NewAABB(Vector3{5, 5, 5}, Vector3{1, 1, 1})

	assert.False(t, s.IntersectsAABB(a))
}

// Test a Sphere/Vector3 intersection hit
func TestSphereIntersectsVector3Hit(t *testing.T) {
	s := NewSphere(Vector3{0, 0, 0}, 1)
	v := NewVector3(0.5, 0.5, 0.5)

	assert.True(t, s.IntersectsVector3(v))
}

// Test a Sphere/Vector3 intersection miss
func TestSphereIntersectsVector3Miss(t *testing.T) {
	s := NewSphere(Vector3{0, 0, 0}, 1)
	v := NewVector3(2, 2, 2)

	assert.False(t, s.IntersectsVector3(v))
}
