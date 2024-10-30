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
	visited := make([]bool, m.GetNumberOfHalfEdges())

	for i, halfEdge := range m.halfEdges {
		if !visited[i] && !halfEdge.IsBoundary() {
			twin := m.GetHalfEdge(halfEdge.Twin)

			if twin.Origin == halfEdge.Origin {
				return false
			}

			visited[i] = true
			visited[halfEdge.Twin] = true
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

// Get the number of faces
func (m *HEMesh) GetNumberOfFaces() int {
	return len(m.faces)
}

// Get the face by ID
func (m *HEMesh) GetFace(id int) HEFace {
	return m.faces[id]
}

// Get the number of half edges
func (m *HEMesh) GetNumberOfHalfEdges() int {
	return len(m.halfEdges)
}

// Get the half edge by ID
func (m *HEMesh) GetHalfEdge(id int) HEHalfEdge {
	return m.halfEdges[id]
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
