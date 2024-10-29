package mtk

import (
	"io"
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
				Prev:   nHalfEdges + (hi-1)%nFaceVertices,
				Next:   nHalfEdges + (hi+1)%nFaceVertices,
				Twin:   -1,
			}

			mesh.halfEdges = append(mesh.halfEdges, halfEdge)
		}
	}

	// TODO: index half edges to get twins

	return &mesh, nil
}

// Construct a half edge mesh from an OBJ file
func NewHEMeshFromOBJ(reader io.Reader) (*HEMesh, error) {
	objReader := NewOBJReader()
	soup, err := objReader.Read(reader)

	if err != nil {
		return nil, err
	}

	return NewHEMeshFromPolygonSoup(soup)
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
