package spatial

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ajcurley/mtk/geometry"
)

// Test splitting an octree node
func TestOctreeSplit(t *testing.T) {
	bounds := geometry.NewAABB(geometry.Vector3{1, 1, 1}, geometry.Vector3{1, 1, 1})
	octree := NewOctree(bounds)

	assert.Equal(t, 1, len(octree.nodes))
	assert.True(t, octree.nodes[1].isLeaf)

	octree.Split(1)

	assert.Equal(t, 9, len(octree.nodes))
	assert.False(t, octree.nodes[1].isLeaf)
}

// Test inserting items into an octree
func TestOctreeInsert(t *testing.T) {
	bounds := geometry.NewAABB(geometry.Vector3{0.5, 0.5, 0.5}, geometry.Vector3{0.5, 0.5, 0.5})
	octree := NewOctree(bounds)
	count := OctreeMaxItemsPerNode + 1

	for i := 0; i < count; i++ {
		point := geometry.Vector3{
			float64(i) / float64(count),
			float64(i) / float64(count),
			float64(i) / float64(count),
		}
		octree.Insert(point)
	}

	assert.Equal(t, count, octree.NumberOfItems())
	assert.Equal(t, 9, len(octree.nodes))
}

// Test querying an octree with an AABB
func TestOctreeQueryAABB(t *testing.T) {
	bounds := geometry.NewAABB(geometry.Vector3{0.5, 0.5, 0.5}, geometry.Vector3{0.5, 0.5, 0.5})
	octree := NewOctree(bounds)
	count := OctreeMaxItemsPerNode * 2

	for i := 0; i < count; i++ {
		point := geometry.Vector3{
			float64(i) / float64(count),
			float64(i) / float64(count),
			float64(i) / float64(count),
		}
		octree.Insert(point)
	}

	assert.Equal(t, count, octree.NumberOfItems())

	query := geometry.NewAABB(geometry.Vector3{0.2, 0.2, 0.2}, geometry.Vector3{0.05, 0.05, 0.05})
	results := octree.Query(query)

	assert.Equal(t, count/10, len(results)) // because of floating point math
}

// Test querying an octree with multiple AABB in parallel
func TestOctreeQueryManyAABB(t *testing.T) {
	bounds := geometry.NewAABB(geometry.Vector3{0.5, 0.5, 0.5}, geometry.Vector3{0.5, 0.5, 0.5})
	octree := NewOctree(bounds)
	count := OctreeMaxItemsPerNode * 2

	for i := 0; i < count; i++ {
		point := geometry.Vector3{
			float64(i) / float64(count),
			float64(i) / float64(count),
			float64(i) / float64(count),
		}
		octree.Insert(point)
	}

	assert.Equal(t, count, octree.NumberOfItems())

	queries := []geometry.IntersectsAABB{
		geometry.NewAABB(geometry.Vector3{0.2, 0.2, 0.2}, geometry.Vector3{0.05, 0.05, 0.05}),
		geometry.NewAABB(geometry.Vector3{0.275, 0.275, 0.275}, geometry.Vector3{0.025, 0.025, 0.025}),
		geometry.NewAABB(geometry.Vector3{0.3, 0.3, 0.3}, geometry.Vector3{0.05, 0.05, 0.05}),
	}

	results := octree.QueryMany(queries)

	assert.Equal(t, 3, len(results))
	assert.Equal(t, count/10, len(results[0])) // because of floating point math
	assert.Equal(t, count/20+1, len(results[1]))
	assert.Equal(t, count/10+1, len(results[2]))
}
