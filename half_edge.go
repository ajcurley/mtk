package mtk

import (
	"errors"
	"io"
	"math"
)

var (
	ErrNonManifoldMesh = errors.New("non-manifold mesh")
)

type HEMesh struct {
	vertices  []HEVertex
	faces     []HEFace
	halfEdges []HEHalfEdge
}

// Construct a half edge mesh from a PolygonSoup
func NewHEMeshFromPolygonSoup(soup *PolygonSoup) (*HEMesh, error) {
	nVertices := soup.GetNumberOfVertices()
	nFaces := soup.GetNumberOfFaces()

	mesh := HEMesh{
		vertices:  make([]HEVertex, 0, nVertices),
		faces:     make([]HEFace, 0, nFaces),
		halfEdges: make([]HEHalfEdge, 0, 3*nFaces),
	}

	// Index the vertices without reference to their originating half edge
	// which will be indexed later during the face insertion.
	indexHalfEdges := make(map[[2]int][]int)

	for vi := 0; vi < nVertices; vi++ {
		vertex := HEVertex{
			Origin: soup.GetVertex(vi),
		}
		mesh.vertices = append(mesh.vertices, vertex)
	}

	// Index the faces and all edges of each face as half edges. Each half
	// edge will not have a reference to its twin until later.
	for fi := 0; fi < nFaces; fi++ {
		faceVertices := soup.GetFace(fi)
		nFaceVertices := len(faceVertices)
		nHalfEdges := len(mesh.halfEdges)

		face := HEFace{HalfEdge: nHalfEdges}
		mesh.faces = append(mesh.faces, face)

		for hi := 0; hi < nFaceVertices; hi++ {
			halfEdge := HEHalfEdge{
				Origin: faceVertices[hi],
				Face:   fi,
				Prev:   nHalfEdges + (hi+nFaceVertices-1)%nFaceVertices,
				Next:   nHalfEdges + (hi+nFaceVertices+1)%nFaceVertices,
				Twin:   -1,
			}

			mesh.halfEdges = append(mesh.halfEdges, halfEdge)

			// Update the originating vertex's half edge
			vertex := mesh.vertices[halfEdge.Origin]
			vertex.HalfEdge = nHalfEdges + hi
			mesh.vertices[halfEdge.Origin] = vertex

			// Index the half edges by the vertices defining the edge. The vertices
			// are sorted since the face orientation does not matter when assigning
			// twin half edges.
			x := faceVertices[hi]
			y := faceVertices[(hi+1)%nFaceVertices]
			k := [2]int{min(x, y), max(x, y)}

			if shared, ok := indexHalfEdges[k]; ok {
				if len(shared) == 2 {
					return nil, ErrNonManifoldMesh
				}

				indexHalfEdges[k] = append(shared, hi+nHalfEdges)
			} else {
				indexHalfEdges[k] = []int{hi + nHalfEdges}
			}
		}
	}

	for _, shared := range indexHalfEdges {
		if len(shared) == 2 {
			for i, index := range shared {
				halfEdge := mesh.halfEdges[index]
				halfEdge.Twin = shared[(i+1)%2]
				mesh.halfEdges[index] = halfEdge
			}
		}
	}

	return &mesh, nil
}

// Construct a half edge mesh from an OBJ file reader
func NewHEMeshFromOBJ(reader io.Reader) (*HEMesh, error) {
	objReader := NewOBJReader()
	soup, err := objReader.Read(reader)

	if err != nil {
		return nil, err
	}

	return NewHEMeshFromPolygonSoup(soup)
}

// Construct a half edge mesh from an OBJ file
func NewHEMeshFromOBJFile(path string) (*HEMesh, error) {
	objReader := NewOBJReader()
	soup, err := objReader.ReadFile(path)

	if err != nil {
		return nil, err
	}

	return NewHEMeshFromPolygonSoup(soup)
}

// Compute the axis-aligned bounding box
func (m *HEMesh) GetAABB() AABB {
	minPosition := Vector3{1, 1, 1}.MulScalar(math.Inf(1))
	maxPosition := Vector3{1, 1, 1}.MulScalar(math.Inf(-1))

	for _, vertex := range m.vertices {
		for i := 0; i < 3; i++ {
			if vertex.Origin[i] < minPosition[i] {
				minPosition[i] = vertex.Origin[i]
			}

			if vertex.Origin[i] > maxPosition[i] {
				maxPosition[i] = vertex.Origin[i]
			}
		}
	}

	return AABB{Min: minPosition, Max: maxPosition}
}

// Check if the half edge mesh is closed (no open boundaries)
func (m *HEMesh) IsClosed() bool {
	for _, halfEdge := range m.halfEdges {
		if halfEdge.IsBoundary() {
			return false
		}
	}

	return true
}

// Check if the half edge mesh has consistently oriented faces
func (m *HEMesh) IsConsistent() bool {
	for _, halfEdge := range m.halfEdges {
		if !halfEdge.IsBoundary() {
			twin := m.GetHalfEdge(halfEdge.Twin)

			if twin.Origin == halfEdge.Origin {
				return false
			}
		}
	}

	return true
}

// Get the number of vertices
func (m *HEMesh) GetNumberOfVertices() int {
	return len(m.vertices)
}

// Get the vertex by ID
func (m *HEMesh) GetVertex(id int) HEVertex {
	return m.vertices[id]
}

// Get the neighboring vertices for a vertex by ID
func (m *HEMesh) GetVertexNeighbors(id int) []int {
	neighbors := make([]int, 0)
	vertex := m.GetVertex(id)
	current := vertex.HalfEdge

	for {
		halfEdge := m.GetHalfEdge(current)

		if halfEdge.IsBoundary() {
			panic("vertex neighbors requires a closed mesh")
		}

		if halfEdge.Origin == id {
			neighbor := m.GetHalfEdge(halfEdge.Next).Origin
			neighbors = append(neighbors, neighbor)
		} else {
			neighbors = append(neighbors, halfEdge.Origin)
		}

		twin := m.GetHalfEdge(halfEdge.Twin)

		if twin.Origin != id {
			current = twin.Next
		} else {
			current = twin.Prev
		}

		if current == vertex.HalfEdge {
			break
		}
	}

	return neighbors
}

// Get the faces using the vertex by ID
func (m *HEMesh) GetVertexFaces(id int) []int {
	faces := make([]int, 0)
	vertex := m.GetVertex(id)
	current := vertex.HalfEdge

	for {
		halfEdge := m.GetHalfEdge(current)

		if halfEdge.IsBoundary() {
			panic("vertex faces requires a closed mesh")
		}

		faces = append(faces, halfEdge.Face)
		twin := m.GetHalfEdge(halfEdge.Twin)

		if twin.Origin != id {
			current = twin.Next
		} else {
			current = twin.Prev
		}

		if current == vertex.HalfEdge {
			break
		}
	}

	return faces
}

// Get the number of faces
func (m *HEMesh) GetNumberOfFaces() int {
	return len(m.faces)
}

// Get the face by ID
func (m *HEMesh) GetFace(id int) HEFace {
	return m.faces[id]
}

// Get the vertices defining the face by ID
func (m *HEMesh) GetFaceVertices(id int) []int {
	faceHalfEdges := m.GetFaceHalfEdges(id)
	vertices := make([]int, len(faceHalfEdges))

	for i, faceHalfEdge := range faceHalfEdges {
		vertices[i] = m.GetHalfEdge(faceHalfEdge).Origin
	}

	return vertices
}

// Get the neighboring faces for a face by ID
func (m *HEMesh) GetFaceNeighbors(id int) []int {
	faceHalfEdges := m.GetFaceHalfEdges(id)
	neighbors := make([]int, 0, len(faceHalfEdges))

	for _, faceHalfEdge := range faceHalfEdges {
		halfEdge := m.GetHalfEdge(faceHalfEdge)

		if !halfEdge.IsBoundary() {
			twin := m.GetHalfEdge(halfEdge.Twin)
			neighbors = append(neighbors, twin.Face)
		}
	}

	return neighbors
}

// Get the half edges of the face by ID
func (m *HEMesh) GetFaceHalfEdges(id int) []int {
	halfEdges := make([]int, 0)

	face := m.GetFace(id)
	current := face.HalfEdge

	for {
		halfEdges = append(halfEdges, current)
		current = m.GetHalfEdge(current).Next

		if current == face.HalfEdge {
			break
		}
	}

	return halfEdges
}

// Get the number of half edges
func (m *HEMesh) GetNumberOfHalfEdges() int {
	return len(m.halfEdges)
}

// Get the half edge by ID
func (m *HEMesh) GetHalfEdge(id int) HEHalfEdge {
	return m.halfEdges[id]
}

// Get the distinct components (connected faces). Each component is
// defined by the indices of the faces.
func (m *HEMesh) GetComponents() [][]int {
	components := make([][]int, 0)
	visited := make([]bool, m.GetNumberOfFaces())

	for next, isVisited := range visited {
		if !isVisited {
			component := make([]int, 0)
			queue := []int{next}

			for len(queue) > 0 {
				current := queue[0]
				queue = queue[1:]

				if !visited[current] {
					visited[current] = true
					component = append(component, current)

					for _, neighbor := range m.GetFaceNeighbors(current) {
						if !visited[neighbor] {
							queue = append(queue, neighbor)
						}
					}
				}
			}

			components = append(components, component)
		}
	}

	return components
}

// Naively copy another half edge mesh into the current. This does not
// merge any duplicate vertices or faces.
func (m *HEMesh) Merge(other *HEMesh) {
	offsetVertices := m.GetNumberOfVertices()
	offsetFaces := m.GetNumberOfFaces()
	offsetHalfEdges := m.GetNumberOfHalfEdges()

	m.vertices = append(m.vertices, other.vertices...)
	m.faces = append(m.faces, other.faces...)
	m.halfEdges = append(m.halfEdges, other.halfEdges...)

	for i, vertex := range m.vertices[offsetVertices:] {
		vertex.HalfEdge += offsetHalfEdges
		m.vertices[i] = vertex
	}

	for i, face := range m.faces[offsetFaces:] {
		face.HalfEdge += offsetHalfEdges
		m.faces[i] = face
	}

	for i, halfEdge := range m.halfEdges[offsetHalfEdges:] {
		halfEdge.Origin += offsetVertices
		halfEdge.Face += offsetFaces
		halfEdge.Prev += offsetHalfEdges
		halfEdge.Next += offsetHalfEdges

		if !halfEdge.IsBoundary() {
			halfEdge.Twin += offsetHalfEdges
		}

		m.halfEdges[i] = halfEdge
	}
}

type HEVertex struct {
	Origin   Vector3
	HalfEdge int
}

type HEFace struct {
	HalfEdge int
}

type HEHalfEdge struct {
	Origin int
	Face   int
	Prev   int
	Next   int
	Twin   int
}

// Get if the half edge is a boundary (no twin)
func (e HEHalfEdge) IsBoundary() bool {
	return e.Twin < 0
}
