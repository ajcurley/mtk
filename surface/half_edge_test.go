package surface

import (
	"errors"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ajcurley/mtk/geometry"
)

// Test reading from file
func TestNewHEMeshFromOBJFile(t *testing.T) {
	path := "../testdata/box.obj"
	mesh, err := NewHEMeshFromOBJFile(path)

	assert.Empty(t, err)
	assert.Equal(t, 8, mesh.NumberOfVertices())
	assert.Equal(t, 12, mesh.NumberOfFaces())
	assert.Equal(t, 36, mesh.NumberOfHalfEdges())
}

// Test reading from file with patches
func TestNewHEMeshFromOBJFilePatches(t *testing.T) {
	path := "../testdata/box.groups.obj"
	mesh, err := NewHEMeshFromOBJFile(path)

	assert.Empty(t, err)
	assert.Equal(t, 6, mesh.NumberOfPatches())
	assert.Equal(t, 0, mesh.Face(0).Patch)
	assert.Equal(t, 1, mesh.Face(1).Patch)
	assert.Equal(t, 1, mesh.Face(2).Patch)
	assert.Equal(t, 2, mesh.Face(3).Patch)
	assert.Equal(t, 3, mesh.Face(4).Patch)
	assert.Equal(t, 4, mesh.Face(5).Patch)
	assert.Equal(t, 5, mesh.Face(6).Patch)
}

// Test for a non-manifold mesh
func TestNewHEMeshFromOBJFileNonManifold(t *testing.T) {
	path := "../testdata/box.nonmanifold.obj"
	_, err := NewHEMeshFromOBJFile(path)

	assert.True(t, errors.Is(err, ErrNonManifoldMesh))
}

// Test for a closed mesh
func TestHEMeshIsClosedTrue(t *testing.T) {
	path := "../testdata/box.obj"
	mesh, _ := NewHEMeshFromOBJFile(path)

	assert.True(t, mesh.IsClosed())
}

// Test for an open mesh
func TestHEMeshIsClosedFalse(t *testing.T) {
	path := "../testdata/box.open.obj"
	mesh, _ := NewHEMeshFromOBJFile(path)

	assert.False(t, mesh.IsClosed())
}

// Test for a consistently oriented mesh
func TestHEMeshIsConsistentTrue(t *testing.T) {
	path := "../testdata/box.obj"
	mesh, _ := NewHEMeshFromOBJFile(path)

	assert.True(t, mesh.IsConsistent())
}

// Test for an inconsistently oriented mesh
func TestHEMeshIsConsistentFalse(t *testing.T) {
	path := "../testdata/box.inconsistent.obj"
	mesh, _ := NewHEMeshFromOBJFile(path)

	assert.False(t, mesh.IsConsistent())
}

// Test computing the bounding box
func TestHEMeshBounds(t *testing.T) {
	path := "../testdata/box.obj"
	mesh, _ := NewHEMeshFromOBJFile(path)

	aabb := mesh.Bounds()

	assert.Equal(t, aabb.Min(), geometry.Vector3{-0.5, -0.5, -0.5})
	assert.Equal(t, aabb.Max(), geometry.Vector3{0.5, 0.5, 0.5})
}

// Test for the vertex neighbors of a consistently oriented mesh
func TestHEMeshVertexNeighborsConsistent(t *testing.T) {
	path := "../testdata/box.obj"
	mesh, _ := NewHEMeshFromOBJFile(path)

	neighbors := mesh.VertexNeighbors(1)

	assert.Equal(t, 5, len(neighbors))
	assert.Equal(t, 5, neighbors[0])
	assert.Equal(t, 4, neighbors[1])
	assert.Equal(t, 0, neighbors[2])
	assert.Equal(t, 2, neighbors[3])
	assert.Equal(t, 3, neighbors[4])
}

// Test for the vertex neighbors of an inconsistently oriented mesh
func TestHEMeshVertexNeighborsInconsistent(t *testing.T) {
	path := "../testdata/box.inconsistent.obj"
	mesh, _ := NewHEMeshFromOBJFile(path)

	neighbors := mesh.VertexNeighbors(1)

	assert.Equal(t, 5, len(neighbors))
	assert.Equal(t, 3, neighbors[0])
	assert.Equal(t, 2, neighbors[1])
	assert.Equal(t, 0, neighbors[2])
	assert.Equal(t, 4, neighbors[3])
	assert.Equal(t, 5, neighbors[4])
}

// Test for the vertex faces of a consistently oriented mesh
func TestHEMeshVertexFacesConsistent(t *testing.T) {
	path := "../testdata/box.obj"
	mesh, _ := NewHEMeshFromOBJFile(path)

	faces := mesh.VertexFaces(1)

	assert.Equal(t, 5, len(faces))
	assert.Equal(t, 10, faces[0])
	assert.Equal(t, 5, faces[1])
	assert.Equal(t, 4, faces[2])
	assert.Equal(t, 0, faces[3])
	assert.Equal(t, 1, faces[4])
}

// Test for the vertex faces of an inconsistently oriented mesh
func TestHEMeshVertexFacesInconsistent(t *testing.T) {
	path := "../testdata/box.inconsistent.obj"
	mesh, _ := NewHEMeshFromOBJFile(path)

	faces := mesh.VertexFaces(1)

	assert.Equal(t, 5, len(faces))
	assert.Equal(t, 10, faces[0])
	assert.Equal(t, 1, faces[1])
	assert.Equal(t, 0, faces[2])
	assert.Equal(t, 4, faces[3])
	assert.Equal(t, 5, faces[4])
}

// Test for the face vertices
func TestHEMeshFaceVertices(t *testing.T) {
	path := "../testdata/box.obj"
	mesh, _ := NewHEMeshFromOBJFile(path)

	vertices := mesh.FaceVertices(1)

	assert.Equal(t, 3, len(vertices))
	assert.Equal(t, vertices[0], 1)
	assert.Equal(t, vertices[1], 3)
	assert.Equal(t, vertices[2], 2)
}

// Test for the face neighbors
func TestHEMeshFaceNeighbors(t *testing.T) {
	path := "../testdata/box.obj"
	mesh, _ := NewHEMeshFromOBJFile(path)

	neighbors := mesh.FaceNeighbors(1)

	assert.Equal(t, 3, len(neighbors))
	assert.Equal(t, 10, neighbors[0])
	assert.Equal(t, 6, neighbors[1])
	assert.Equal(t, 0, neighbors[2])
}

// Test for the distinct components for a single element mesh
func TestHEMeshComponentsSingle(t *testing.T) {
	path := "../testdata/box.obj"
	mesh, _ := NewHEMeshFromOBJFile(path)

	components := mesh.Components()

	assert.Equal(t, 1, len(components))
	assert.Equal(t, mesh.NumberOfFaces(), len(components[0]))
}

// Test for the distinct components for a multiple element mesh
func TestHEMeshComponentsMultiple(t *testing.T) {
	boxPath := "../testdata/box.obj"
	mesh, _ := NewHEMeshFromOBJFile(boxPath)

	spherePath := "../testdata/sphere.obj"
	meshSphere, _ := NewHEMeshFromOBJFile(spherePath)

	mesh.Merge(meshSphere)
	components := mesh.Components()

	assert.Equal(t, 2, len(components))
	assert.Equal(t, 12, len(components[0]))
	assert.Equal(t, 96, len(components[1]))
}

// Test getting the shared vertices between two faces
func TestHEMeshSharedVertices(t *testing.T) {
	path := "../testdata/box.obj"
	mesh, _ := NewHEMeshFromOBJFile(path)

	shared := mesh.SharedVertices(0, 1)

	assert.Equal(t, 2, len(shared))
}

// Test merging two meshes
func TestHEMeshMerge(t *testing.T) {
	path := "../testdata/box.obj"
	mesh, _ := NewHEMeshFromOBJFile(path)
	other, _ := NewHEMeshFromOBJFile(path)

	mesh.Merge(other)

	assert.Equal(t, 16, mesh.NumberOfVertices())
	assert.Equal(t, 24, mesh.NumberOfFaces())
	assert.Equal(t, 72, mesh.NumberOfHalfEdges())
	assert.True(t, mesh.IsClosed())
	assert.True(t, mesh.IsConsistent())
}

// Test merging two meshes with unique patch names
func TestHEMeshMergeUniquePatches(t *testing.T) {
	boxPath := "../testdata/box.groups.obj"
	mesh, _ := NewHEMeshFromOBJFile(boxPath)

	spherePath := "../testdata/sphere.groups.obj"
	meshSphere, _ := NewHEMeshFromOBJFile(spherePath)

	mesh.Merge(meshSphere)

	assert.Equal(t, 103, mesh.NumberOfFaces())
	assert.Equal(t, 7, mesh.NumberOfPatches())
}

// Test merging two meshes with overlapping patch names
func TestHEMeshMergeSharedPatches(t *testing.T) {
	path := "../testdata/box.groups.obj"
	mesh, _ := NewHEMeshFromOBJFile(path)
	other, _ := NewHEMeshFromOBJFile(path)

	mesh.Merge(other)

	assert.Equal(t, 16, mesh.NumberOfVertices())
	assert.Equal(t, 14, mesh.NumberOfFaces())
	assert.Equal(t, 52, mesh.NumberOfHalfEdges())
	assert.Equal(t, 6, mesh.NumberOfPatches())
}

// Test orienting a consistently oriented mesh
func TestHEMeshOrientConsistent(t *testing.T) {
	path := "../testdata/box.obj"
	mesh, _ := NewHEMeshFromOBJFile(path)

	assert.True(t, mesh.IsConsistent())

	mesh.Orient()

	assert.True(t, mesh.IsConsistent())
}

// Test orienting an consistently oriented mesh
func TestHEMeshOrientInconsistent(t *testing.T) {
	path := "../testdata/box.inconsistent.obj"
	mesh, _ := NewHEMeshFromOBJFile(path)

	assert.False(t, mesh.IsConsistent())

	mesh.Orient()

	assert.True(t, mesh.IsConsistent())
}

// Test for a face normal
func TestHEMeshFaceNormal(t *testing.T) {
	path := "../testdata/box.obj"
	mesh, _ := NewHEMeshFromOBJFile(path)

	normal := mesh.FaceNormal(0)

	assert.Equal(t, geometry.Vector3{-1, 0, 0}, normal)
}

// Test extract faces from a mesh
func TestHEMeshExtractFaces(t *testing.T) {
	path := "../testdata/box.obj"
	mesh, _ := NewHEMeshFromOBJFile(path)

	subset, err := mesh.ExtractFaces([]int{0, 1, 7})

	assert.Empty(t, err)
	assert.Equal(t, 6, subset.NumberOfVertices())
	assert.Equal(t, 3, subset.NumberOfFaces())
	assert.Equal(t, 9, subset.NumberOfHalfEdges())
}

// Test extract patch names from a mesh
func TestHEMeshExtractPatchNames(t *testing.T) {
	path := "../testdata/box.groups.obj"
	mesh, _ := NewHEMeshFromOBJFile(path)

	subset, err := mesh.ExtractPatchNames([]string{"bottom", "back"})

	assert.Empty(t, err)
	assert.Equal(t, 6, subset.NumberOfVertices())
	assert.Equal(t, 3, subset.NumberOfFaces())
	assert.Equal(t, 10, subset.NumberOfHalfEdges())
}

// Test zipping edges for a mesh with some open and some closed edges
func TestHEMeshMergeZipEdgesPartial(t *testing.T) {
	path := "../testdata/box.duplicates-partial.obj"
	mesh, _ := NewHEMeshFromOBJFile(path)

	assert.Equal(t, 22, mesh.NumberOfVertices())
	assert.False(t, mesh.IsClosed())

	err := mesh.ZipEdges()

	assert.Empty(t, err)
	assert.Equal(t, 8, mesh.NumberOfVertices())
	assert.True(t, mesh.IsConsistent())
	assert.True(t, mesh.IsClosed())
}

// Test zipping edges for a mesh that results in a non-manifold edge
func TestHEMeshMergeZipEdgesNonManifold(t *testing.T) {
	path := "../testdata/box.duplicates-nonmanifold.obj"
	mesh, _ := NewHEMeshFromOBJFile(path)

	err := mesh.ZipEdges()

	assert.ErrorContains(t, err, "non-manifold mesh: near [0 0 0.5]")
}

// Test computing the feature edges
func TestHEMeshFeatureEdges(t *testing.T) {
	path := "../testdata/box.obj"
	mesh, _ := NewHEMeshFromOBJFile(path)

	threshold := math.Pi * 30. / 180.
	featureEdges := mesh.FeatureEdges(threshold)

	assert.Equal(t, 12, len(featureEdges))
}

// Test computing the principal axes
func TestHEMeshPrincipalAxes(t *testing.T) {
	path := "../testdata/box.obj"
	mesh, _ := NewHEMeshFromOBJFile(path)

	axes := mesh.PrincipalAxes()

	assert.Equal(t, len(axes), 3)
	assert.Equal(t, geometry.NewVector3(1, 0, 0), axes[0])
	assert.Equal(t, geometry.NewVector3(0, 1, 0), axes[1])
	assert.Equal(t, geometry.NewVector3(0, 0, 1), axes[2])
}
