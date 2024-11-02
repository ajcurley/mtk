package mtk

import (
	"math/rand/v2"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test splitting an octree node
func TestOctreeSplit(t *testing.T) {
	bounds := AABB{Min: Vector3{0, 0, 0}, Max: Vector3{2, 2, 2}}
	octree := NewOctree[Vector3](bounds)

	assert.Equal(t, 1, len(octree.nodes))
	assert.True(t, octree.nodes[1].isLeaf)

	octree.Split(1)

	assert.Equal(t, 9, len(octree.nodes))
	assert.False(t, octree.nodes[1].isLeaf)
}

// Test inserting items into an octree
func TestOctreeInsert(t *testing.T) {
	bounds := AABB{Min: Vector3{0, 0, 0}, Max: Vector3{1, 1, 1}}
	octree := NewOctree[Vector3](bounds)

	for i := 0; i < OctreeMaxItemsPerNode+1; i++ {
		point := Vector3{
			rand.Float64(),
			rand.Float64(),
			rand.Float64(),
		}
		octree.Insert(point)
	}

	assert.Equal(t, OctreeMaxItemsPerNode+1, len(octree.items))
	assert.Equal(t, 9, len(octree.nodes))
}
