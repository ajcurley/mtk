package surface

import (
	"bytes"
	"compress/gzip"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ajcurley/mtk/geometry"
)

// Read an OBJ file from path.
func TestOBJReaderReadFile(t *testing.T) {
	path := "../testdata/box.obj"

	objReader := NewOBJReader()
	soup, err := objReader.ReadFile(path)

	assert.Empty(t, err)
	assert.Equal(t, 8, soup.NumberOfVertices())
	assert.Equal(t, 12, soup.NumberOfFaces())
	assert.Equal(t, 0, soup.NumberOfPatches())
}

// Read an OBJ file from path (gzip).
func TestOBJReaderReadFileGZIP(t *testing.T) {
	path := "../testdata/box.obj.gz"

	objReader := NewOBJReader()
	soup, err := objReader.ReadFile(path)

	assert.Empty(t, err)
	assert.Equal(t, 8, soup.NumberOfVertices())
	assert.Equal(t, 12, soup.NumberOfFaces())
	assert.Equal(t, 0, soup.NumberOfPatches())
}

// Read an OBJ file from path with mixed elements and patches.
func TestOBJReaderReadFileGroups(t *testing.T) {
	path := "../testdata/box.groups.obj"

	objReader := NewOBJReader()
	soup, err := objReader.ReadFile(path)

	assert.Empty(t, err)
	assert.Equal(t, 8, soup.NumberOfVertices())
	assert.Equal(t, 7, soup.NumberOfFaces())
	assert.Equal(t, 6, soup.NumberOfPatches())
}

// Write an OBJ file.
func TestOBJWriterWrite(t *testing.T) {
	vertices := []geometry.Vector3{
		geometry.Vector3{0, 0, 0},
		geometry.Vector3{0, 1, 0},
		geometry.Vector3{1, 1, 0},
		geometry.Vector3{1, 1, 1},
	}

	faces := [][]int{
		[]int{0, 1, 2},
		[]int{2, 0, 3},
	}

	groups := []string{
		"testFace",
	}

	faceGroups := []int{
		0,
		-1,
	}

	lines := [][]int{
		[]int{0, 1},
		[]int{2, 3},
	}

	var expected string
	expected += "v 0.000000 0.000000 0.000000\n"
	expected += "v 0.000000 1.000000 0.000000\n"
	expected += "v 1.000000 1.000000 0.000000\n"
	expected += "v 1.000000 1.000000 1.000000\n"
	expected += "l 1 2\n"
	expected += "l 3 4\n"
	expected += "f 3 1 4\n"
	expected += "g testFace\n"
	expected += "f 1 2 3\n"

	objWriter := NewOBJWriter()
	objWriter.SetVertices(vertices)
	objWriter.SetLines(lines)
	objWriter.SetFaces(faces)
	objWriter.SetFaceGroups(faceGroups)
	objWriter.SetGroups(groups)

	var writer bytes.Buffer
	err := objWriter.Write(&writer)

	assert.Empty(t, err)
	assert.Equal(t, expected, writer.String())
}

// Write an OBJ file (gzip).
func TestOBJWriterWriteGZIP(t *testing.T) {
	vertices := []geometry.Vector3{
		geometry.Vector3{0, 0, 0},
		geometry.Vector3{0, 1, 0},
		geometry.Vector3{1, 1, 0},
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

	objWriter := NewOBJWriter()
	objWriter.SetVertices(vertices)
	objWriter.SetFaces(faces)

	var writer bytes.Buffer
	gzipWriter := gzip.NewWriter(&writer)
	err := objWriter.Write(gzipWriter)

	assert.Empty(t, err)
	assert.Equal(t, expectedBuf.String(), writer.String())
}
