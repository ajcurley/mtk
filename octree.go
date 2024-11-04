package mtk

import (
	"runtime"
	"sync"
)

const (
	OctreeMaxDepth        int = 21
	OctreeMaxItemsPerNode int = 100
)

// Linear octree implementation
type Octree struct {
	nodes map[uint64]*octreeNode
	items []IntersectsAABB
}

// Construct an Octree indexing items
func NewOctree(bounds AABB) *Octree {
	return &Octree{
		nodes: map[uint64]*octreeNode{1: newOctreeNode(1, bounds)},
		items: make([]IntersectsAABB, 0),
	}
}

// Get the number of indexed items
func (o *Octree) GetNumberOfItems() int {
	return len(o.items)
}

// Get an item by ID
func (o *Octree) GetItem(id int) IntersectsAABB {
	return o.items[id]
}

// Insert an item into the octree
func (o *Octree) Insert(item IntersectsAABB) (int, bool) {
	var code uint64
	index := len(o.items)
	queue := []uint64{1}
	codes := make([]uint64, 0)

	for len(queue) > 0 {
		code, queue = queue[0], queue[1:]
		node := o.nodes[code]

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
func (o *Octree) Split(code uint64) {
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

// Query the octree for intersecting items
func (o *Octree) Query(query IntersectsAABB) []int {
	var code uint64
	items := make(map[int]struct{})
	queue := []uint64{1}

	for len(queue) > 0 {
		code, queue = queue[0], queue[1:]
		node := o.nodes[code]

		if query.IntersectsAABB(node.bounds) {
			if node.isLeaf {
				for _, index := range node.items {
					if _, ok := items[index]; !ok {
						var intersects bool

						switch value := query.(type) {
						case AABB:
							if item, ok := o.items[index].(IntersectsAABB); ok {
								intersects = item.IntersectsAABB(value)
							}
						case *AABB:
							if item, ok := o.items[index].(IntersectsAABB); ok {
								intersects = item.IntersectsAABB(*value)
							}
						case Ray:
							if item, ok := o.items[index].(IntersectsRay); ok {
								intersects = item.IntersectsRay(value)
							}
						case *Ray:
							if item, ok := o.items[index].(IntersectsRay); ok {
								intersects = item.IntersectsRay(*value)
							}
						case Triangle:
							if item, ok := o.items[index].(IntersectsTriangle); ok {
								intersects = item.IntersectsTriangle(value)
							}
						case *Triangle:
							if item, ok := o.items[index].(IntersectsTriangle); ok {
								intersects = item.IntersectsTriangle(*value)
							}
						case Vector3:
							if item, ok := o.items[index].(IntersectsVector3); ok {
								intersects = item.IntersectsVector3(value)
							}
						case *Vector3:
							if item, ok := o.items[index].(IntersectsVector3); ok {
								intersects = item.IntersectsVector3(*value)
							}
						}

						if intersects {
							items[index] = struct{}{}
						}
					}
				}
			} else {
				childrenCodes := node.childrenCodes()
				queue = append(queue, childrenCodes...)
			}
		}
	}

	results := make([]int, 0, len(items))

	for index := range items {
		results = append(results, index)
	}

	return results
}

// Query the octree for many intersecting items in parallel using the available
// number of processors.
func (o *Octree) QueryMany(queries []IntersectsAABB) [][]int {
	var wg sync.WaitGroup
	queue := make(chan int, len(queries))
	items := make([][]int, len(queries))

	for i := range queries {
		queue <- i
	}

	for i := 0; i < runtime.NumCPU(); i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for i := range queue {
				items[i] = o.Query(queries[i])
			}
		}()
	}

	close(queue)
	wg.Wait()

	return items
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
