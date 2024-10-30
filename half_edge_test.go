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
