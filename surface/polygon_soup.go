package surface

import (
	"github.com/ajcurley/mtk/geometry"
)

type PolygonSoup struct {
	vertices     []geometry.Vector3
	faceOffsets  []int
	faceVertices []int
	facePatches  []int
	patches      []string
}

func NewPolygonSoup() *PolygonSoup {
	return &PolygonSoup{
		vertices:     make([]geometry.Vector3, 0),
		faceOffsets:  make([]int, 0),
		faceVertices: make([]int, 0),
		facePatches:  make([]int, 0),
		patches:      make([]string, 0),
	}
}

// Get the number of vertices
func (m *PolygonSoup) NumberOfVertices() int {
	return len(m.vertices)
}

// Get a vertex by ID
func (m *PolygonSoup) Vertex(id int) geometry.Vector3 {
	return m.vertices[id]
}

// Insert a vertex
func (m *PolygonSoup) InsertVertex(vertex geometry.Vector3) int {
	m.vertices = append(m.vertices, vertex)
	return m.NumberOfVertices() - 1
}

// Get the number of faces
func (m *PolygonSoup) NumberOfFaces() int {
	return len(m.faceOffsets)
}

// Get a face's ordered set of vertices by face ID
func (m *PolygonSoup) Face(id int) []int {
	nFaces := m.NumberOfFaces()
	offset := m.faceOffsets[id]

	if id < nFaces-1 {
		nextOffset := m.faceOffsets[id+1]
		return m.faceVertices[offset:nextOffset]
	}

	return m.faceVertices[offset:]
}

// Get a face's patch by face ID
func (m *PolygonSoup) FacePatch(id int) int {
	return m.facePatches[id]
}

// Insert a face. By default, the patch is empty.
func (m *PolygonSoup) InsertFace(vertices []int) int {
	m.faceOffsets = append(m.faceOffsets, len(m.faceVertices))
	m.faceVertices = append(m.faceVertices, vertices...)
	m.facePatches = append(m.facePatches, -1)
	return m.NumberOfFaces() - 1
}

// Insert a face with a patch
func (m *PolygonSoup) InsertFaceWithPatch(vertices []int, patch int) int {
	id := m.InsertFace(vertices)
	m.facePatches[id] = patch
	return id
}

// Get the number of patches
func (m *PolygonSoup) NumberOfPatches() int {
	return len(m.patches)
}

// Get a patch by ID
func (m *PolygonSoup) Patch(id int) string {
	return m.patches[id]
}

// Insert a patch
func (m *PolygonSoup) InsertPatch(name string) int {
	m.patches = append(m.patches, name)
	return m.NumberOfPatches() - 1
}
