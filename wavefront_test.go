package mtk

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Read an OBJ file from path.
func TestOBJReaderReadFile(t *testing.T) {
	path := "testdata/box.obj"

	objReader := NewOBJReader()
	soup, err := objReader.ReadFile(path)

	assert.Empty(t, err)
	assert.Equal(t, 24, soup.GetNumberOfVertices())
	assert.Equal(t, 12, soup.GetNumberOfFaces())
	assert.Equal(t, 0, soup.GetNumberOfPatches())
}

// Read an OBJ file from path (gzip).
func TestOBJReaderReadFileGZIP(t *testing.T) {
	path := "testdata/box.obj.gz"

	objReader := NewOBJReader()
	soup, err := objReader.ReadFile(path)

	assert.Empty(t, err)
	assert.Equal(t, 24, soup.GetNumberOfVertices())
	assert.Equal(t, 12, soup.GetNumberOfFaces())
	assert.Equal(t, 0, soup.GetNumberOfPatches())
}

// Read an OBJ file from path with mixed elements and patches.
func TestOBJReaderReadFileGroups(t *testing.T) {
	path := "testdata/box.groups.obj"

	objReader := NewOBJReader()
	soup, err := objReader.ReadFile(path)

	assert.Empty(t, err)
	assert.Equal(t, 8, soup.GetNumberOfVertices())
	assert.Equal(t, 7, soup.GetNumberOfFaces())
	assert.Equal(t, 6, soup.GetNumberOfPatches())
}

/*
// Write an OBJ file.
func TestWriteOBJ(t *testing.T) {
	vertices := []Vector{
		NewVector(0, 0, 0),
		NewVector(0, 1, 0),
		NewVector(1, 1, 0),
	}

	faces := [][]int{
		[]int{0, 1, 2},
	}

	var expected string
	expected += "v 0.000000 0.000000 0.000000\n"
	expected += "v 0.000000 1.000000 0.000000\n"
	expected += "v 1.000000 1.000000 0.000000\n"
	expected += "f 1 2 3\n"

	var writer bytes.Buffer
	objWriter := NewOBJWriter(&writer)
	objWriter.SetVertices(vertices)
	objWriter.SetFaces(faces)

	err := objWriter.Write()
	assert.Empty(t, err)
	assert.Equal(t, expected, writer.String())
}

// Write an OBJ file (gzip).
func TestWriteOBJGZIP(t *testing.T) {
	vertices := []Vector{
		NewVector(0, 0, 0),
		NewVector(0, 1, 0),
		NewVector(1, 1, 0),
	}

	faces := [][]int{
		[]int{0, 1, 2},
	}

	var expected string
	expected += "v 0.000000 0.000000 0.000000\n"
	expected += "v 0.000000 1.000000 0.000000\n"
	expected += "v 1.000000 1.000000 0.000000\n"
	expected += "f 1 2 3\n"

	var expectedBuf bytes.Buffer
	expectedWriter := gzip.NewWriter(&expectedBuf)
	expectedWriter.Write([]byte(expected))

	var writer bytes.Buffer
	gzipWriter := gzip.NewWriter(&writer)
	objWriter := NewOBJWriter(gzipWriter)
	objWriter.SetVertices(vertices)
	objWriter.SetFaces(faces)

	err := objWriter.Write()
	assert.Empty(t, err)
	assert.Equal(t, expectedBuf.String(), writer.String())
}
*/
