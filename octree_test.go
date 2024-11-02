package mtk

import (
	"math/rand/v2"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test splitting an octree node
func TestOctreeSplit(t *testing.T) {
	bounds := AABB{Min: Vector3{0, 0, 0}, Max: Vector3{2, 2, 2}}
	octree := NewOctree(bounds)

	assert.Equal(t, 1, len(octree.nodes))
	assert.True(t, octree.nodes[1].isLeaf)

	octree.Split(1)

	assert.Equal(t, 9, len(octree.nodes))
	assert.False(t, octree.nodes[1].isLeaf)
}

// Test inserting items into an octree
func TestOctreeInsert(t *testing.T) {
	bounds := AABB{Min: Vector3{0, 0, 0}, Max: Vector3{1, 1, 1}}
	octree := NewOctree(bounds)

	for i := 0; i < OctreeMaxItemsPerNode+1; i++ {
		point := Vector3{
			rand.Float64(),
			rand.Float64(),
			rand.Float64(),
		}
		octree.Insert(point)
	}

	assert.Equal(t, OctreeMaxItemsPerNode+1, octree.GetNumberOfItems())
	assert.Equal(t, 9, len(octree.nodes))
}

// Test querying an octree with an AABB
func TestOctreeQueryAABB(t *testing.T) {
	bounds := AABB{Min: Vector3{0, 0, 0}, Max: Vector3{1, 1, 1}}
	octree := NewOctree(bounds)
	count := OctreeMaxItemsPerNode * 2

	for i := 0; i < count; i++ {
		point := Vector3{
			float64(i) / float64(count),
			float64(i) / float64(count),
			float64(i) / float64(count),
		}
		octree.Insert(point)
	}

	assert.Equal(t, count, octree.GetNumberOfItems())

	query := AABB{
		Min: Vector3{0.15, 0.15, 0.15},
		Max: Vector3{0.25, 0.25, 0.25},
	}
	results := octree.Query(query)

	assert.Equal(t, 11, len(results))
}
