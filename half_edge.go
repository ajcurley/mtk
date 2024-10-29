package mtk

import (
	"io"
)

type HEMesh struct {
	vertices  []HEVertex
	faces     []HEFace
	halfEdges []HEHalfEdge
}

// Construct a half edge mesh from an OBJ file
func NewHEMeshFromOBJ(reader io.Reader) (*HEMesh, error) {
	objReader := NewOBJReader()
	if err := objReader.Read(reader); err != nil {
		return nil, err
	}

	mesh := HEMesh{
		vertices:  make([]HEVertex, 0, objReader.GetNumberOfVertices()),
		faces:     make([]HEFace, 0, objReader.GetNumberOfFaces()),
		halfEdges: make([]HEHalfEdge, 0, 3*objReader.GetNumberOfFaces()),
	}

	// Index the vertices without reference to their originating half edge. These
	// will be indexed later.
	for vertexID := 0; vertexID < len(mesh.vertices); vertexID++ {
		vertex := HEVertex{
			Origin:     objReader.vertices[vertexID],
			HalfEdgeID: -1,
		}
		mesh.vertices = append(mesh.vertices, vertex)
	}

	// Index the faces and their half edges. For each face, use the first half
	// edge as the reference.
	for faceID = 0; faceID < len(mesh.faces); faceID++ {
		// TODO: implement
	}

	return &mesh, nil
}

type HEVertex struct {
	Origin     Vector3
	HalfEdgeID int
}

type HEFace struct {
	HalfEdgeID int
}

type HEHalfEdge struct {
	OriginID   int
	FaceID     int
	PreviousID int
	NextID     int
	TwinID     int
}

// Get if the half edge is a boundary (no twin)
func (e HEHalfEdge) IsBoundary() bool {
	return e.TwinID < 0
}
