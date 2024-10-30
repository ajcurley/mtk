package mtk

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test reading from file
func TestNewHEMeshFromOBJFile(t *testing.T) {
	path := "testdata/box.obj"
	mesh, err := NewHEMeshFromOBJFile(path)

	assert.Empty(t, err)
	assert.Equal(t, 8, mesh.GetNumberOfVertices())
	assert.Equal(t, 12, mesh.GetNumberOfFaces())
	assert.Equal(t, 36, mesh.GetNumberOfHalfEdges())
}

// Test for a non-manifold mesh
func TestNewHEMeshFromOBJFileNonManifold(t *testing.T) {
	path := "testdata/box.nonmanifold.obj"
	_, err := NewHEMeshFromOBJFile(path)

	assert.True(t, errors.Is(err, ErrNonManifoldMesh))
}

// Test for a closed mesh
func TestHEMeshIsClosedTrue(t *testing.T) {
	path := "testdata/box.obj"
	mesh, _ := NewHEMeshFromOBJFile(path)

	assert.True(t, mesh.IsClosed())
}

// Test for an open mesh
func TestHEMeshIsClosedFalse(t *testing.T) {
	path := "testdata/box.open.obj"
	mesh, _ := NewHEMeshFromOBJFile(path)

	assert.False(t, mesh.IsClosed())
}

// Test for a consistently oriented mesh
func TestHEMeshIsConsistentTrue(t *testing.T) {
	path := "testdata/box.obj"
	mesh, _ := NewHEMeshFromOBJFile(path)

	assert.True(t, mesh.IsConsistent())
}

// Test for an inconsistently oriented mesh
func TestHEMeshIsConsistentFalse(t *testing.T) {
	path := "testdata/box.inconsistent.obj"
	mesh, _ := NewHEMeshFromOBJFile(path)

	assert.False(t, mesh.IsConsistent())
}

// Test computing the bounding box
func TestHEMeshAABB(t *testing.T) {
	path := "testdata/box.obj"
	mesh, _ := NewHEMeshFromOBJFile(path)

	aabb := mesh.GetAABB()

	assert.Equal(t, aabb.Min, Vector3{-0.5, -0.5, -0.5})
	assert.Equal(t, aabb.Max, Vector3{0.5, 0.5, 0.5})
}

// Test for the vertex neighbors of a consistently oriented mesh
func TestHEMeshVertexNeighborsConsistent(t *testing.T) {
	path := "testdata/box.obj"
	mesh, _ := NewHEMeshFromOBJFile(path)

	neighbors := mesh.GetVertexNeighbors(1)

	assert.Equal(t, 5, len(neighbors))
	assert.Equal(t, 5, neighbors[0])
	assert.Equal(t, 4, neighbors[1])
	assert.Equal(t, 0, neighbors[2])
	assert.Equal(t, 2, neighbors[3])
	assert.Equal(t, 3, neighbors[4])
}

// Test for the vertex neighbors of an inconsistently oriented mesh
func TestHEMeshVertexNeighborsInconsistent(t *testing.T) {
	path := "testdata/box.inconsistent.obj"
	mesh, _ := NewHEMeshFromOBJFile(path)

	neighbors := mesh.GetVertexNeighbors(1)

	assert.Equal(t, 5, len(neighbors))
	assert.Equal(t, 3, neighbors[0])
	assert.Equal(t, 2, neighbors[1])
	assert.Equal(t, 0, neighbors[2])
	assert.Equal(t, 4, neighbors[3])
	assert.Equal(t, 5, neighbors[4])
}

// Test for the vertex faces of a consistently oriented mesh
func TestHEMeshVertexFacesConsistent(t *testing.T) {
	path := "testdata/box.obj"
	mesh, _ := NewHEMeshFromOBJFile(path)

	faces := mesh.GetVertexFaces(1)

	assert.Equal(t, 5, len(faces))
	assert.Equal(t, 10, faces[0])
	assert.Equal(t, 5, faces[1])
	assert.Equal(t, 4, faces[2])
	assert.Equal(t, 0, faces[3])
	assert.Equal(t, 1, faces[4])
}

// Test for the vertex faces of an inconsistently oriented mesh
func TestHEMeshVertexFacesInconsistent(t *testing.T) {
	path := "testdata/box.inconsistent.obj"
	mesh, _ := NewHEMeshFromOBJFile(path)

	faces := mesh.GetVertexFaces(1)

	assert.Equal(t, 5, len(faces))
	assert.Equal(t, 10, faces[0])
	assert.Equal(t, 1, faces[1])
	assert.Equal(t, 0, faces[2])
	assert.Equal(t, 4, faces[3])
	assert.Equal(t, 5, faces[4])
}

// Test for the face vertices
func TestHEMeshFaceVertices(t *testing.T) {
	path := "testdata/box.obj"
	mesh, _ := NewHEMeshFromOBJFile(path)

	vertices := mesh.GetFaceVertices(1)

	assert.Equal(t, 3, len(vertices))
	assert.Equal(t, vertices[0], 1)
	assert.Equal(t, vertices[1], 3)
	assert.Equal(t, vertices[2], 2)
}

// Test for the face neighbors
func TestHEMeshFaceNeighbors(t *testing.T) {
	path := "testdata/box.obj"
	mesh, _ := NewHEMeshFromOBJFile(path)

	neighbors := mesh.GetFaceNeighbors(1)

	assert.Equal(t, 3, len(neighbors))
	assert.Equal(t, 10, neighbors[0])
	assert.Equal(t, 6, neighbors[1])
	assert.Equal(t, 0, neighbors[2])
}
