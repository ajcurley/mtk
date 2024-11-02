package mtk

import ()

const (
	OctreeMaxDepth        int = 21
	OctreeMaxItemsPerNode int = 50
)

// Generic octree implementation
type Octree[T IntersectsAABB] struct {
	nodes map[uint64]*octreeNode
	items []T
}

// Construct an Octree indexing items of type T
func NewOctree[T IntersectsAABB](bounds AABB) *Octree[T] {
	return &Octree[T]{
		nodes: map[uint64]*octreeNode{1: newOctreeNode(1, bounds)},
		items: make([]T, 0),
	}
}

// Insert an item into the octree
func (o *Octree[T]) Insert(item T) (int, bool) {
	index := len(o.items)
	queue := []uint64{1}
	codes := make([]uint64, 0)

	for len(queue) > 0 {
		code := queue[0]
		queue = queue[1:]

		if node, ok := o.nodes[code]; ok {
			if item.IntersectsAABB(node.bounds) {
				if node.isLeaf {
					node.items = append(node.items, index)
					codes = append(codes, code)
				} else {
					childrenCodes := node.childrenCodes()
					queue = append(queue, childrenCodes...)
				}
			}
		}
	}

	if len(codes) == 0 {
		return -1, false
	}

	o.items = append(o.items, item)

	for _, code := range codes {
		if o.nodes[code].shouldSplit() {
			o.Split(code)
		}
	}

	return index, true
}

// Split an octree node
func (o *Octree[T]) Split(code uint64) {
	if node, ok := o.nodes[code]; ok && node.canSplit() {
		for octant, childCode := range node.childrenCodes() {
			bounds := node.bounds.Octant(octant)
			childNode := newOctreeNode(childCode, bounds)

			for _, index := range node.items {
				if o.items[index].IntersectsAABB(bounds) {
					childNode.items = append(childNode.items, index)
				}
			}

			o.nodes[childCode] = childNode
		}

		node.isLeaf = false
		clear(node.items)
	}
}

// Node within an Octree
type octreeNode struct {
	code   uint64
	bounds AABB
	isLeaf bool
	items  []int
}

// Construct an octree node from an AABB and location code
func newOctreeNode(code uint64, bounds AABB) *octreeNode {
	return &octreeNode{
		code:   code,
		bounds: bounds,
		isLeaf: true,
		items:  make([]int, 0),
	}
}

// Get the depth in the octree
func (o *octreeNode) depth() int {
	for i := 0; i < OctreeMaxDepth+1; i++ {
		if o.code>>(3*uint64(i)) == 1 {
			return i
		}
	}

	panic("invalid octree node code")
}

// Get the eight codes of the children nodes
func (o *octreeNode) childrenCodes() []uint64 {
	codes := make([]uint64, 8)

	for i := 0; i < 8; i++ {
		codes[i] = (o.code << 3) | uint64(i)
	}

	return codes
}

// Check if the node can be split
func (o *octreeNode) canSplit() bool {
	return o.isLeaf && o.depth() < OctreeMaxDepth
}

// Check if the node should be split
func (o *octreeNode) shouldSplit() bool {
	return o.canSplit() && len(o.items) > OctreeMaxItemsPerNode
}
