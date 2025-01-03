package mesh

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"unicode"
	"unicode/utf8"

	"github.com/ajcurley/mtk/geometry"
)

const (
	prefixVertex = "v"
	prefixFace   = "f"
	prefixGroup  = "g"
)

var (
	ErrInvalidVertex = errors.New("invalid vertex")
	ErrInvalidFace   = errors.New("invalid face")
)

type OBJReader struct {
	polygonSoup *PolygonSoup
}

func NewOBJReader() *OBJReader {
	return &OBJReader{
		polygonSoup: NewPolygonSoup(),
	}
}

// Read an OBJ file from an io.Reader interface
func (r *OBJReader) Read(reader io.Reader) (*PolygonSoup, error) {
	count := 1
	buffer := bufio.NewReader(reader)

	// Check if the file is gzip compressed.
	testBytes, err := buffer.Peek(2)
	if err != nil {
		return nil, err
	}

	if testBytes[0] == 31 && testBytes[1] == 139 {
		gzipFile, err := gzip.NewReader(buffer)
		if err != nil {
			return nil, err
		}
		defer gzipFile.Close()

		buffer = bufio.NewReader(gzipFile)
	}

	for {
		data, err := buffer.ReadBytes('\n')

		if errors.Is(err, io.EOF) {
			break
		}

		data = bytes.TrimSpace(data)
		prefix := r.parsePrefix(data)

		switch string(prefix) {
		case prefixVertex:
			err = r.parseVertex(data)
		case prefixFace:
			err = r.parseFace(data)
		case prefixGroup:
			r.parseGroup(data)
		}

		if err != nil {
			return nil, fmt.Errorf("line %d: %v", count, err)
		}

		count++
	}

	return r.polygonSoup, nil
}

// Read an OBJ file from path
func (r *OBJReader) ReadFile(path string) (*PolygonSoup, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return r.Read(file)
}

// Parse a prefix to determine the data type to read
func (r *OBJReader) parsePrefix(data []byte) []byte {
	for i := 0; i < len(data); i++ {
		text, _ := utf8.DecodeRune(data[i : i+1])

		if unicode.IsSpace(text) {
			return data[:i]
		}
	}

	return data
}

// Parse a vertex from a line
func (r *OBJReader) parseVertex(data []byte) error {
	fields := bytes.Fields(data[len(prefixVertex):])

	if len(fields) != 3 {
		return ErrInvalidVertex
	}

	var values [3]float64

	for i := 0; i < len(values); i++ {
		text := string(fields[i])
		value, err := strconv.ParseFloat(text, 64)

		if err != nil {
			return err
		}

		values[i] = value
	}

	vertex := geometry.Vector3(values)
	r.polygonSoup.InsertVertex(vertex)

	return nil
}

// Parse a face from a line
func (r *OBJReader) parseFace(data []byte) error {
	fields := bytes.Fields(data[len(prefixFace):])

	if len(fields) <= 2 {
		return ErrInvalidFace
	}

	face := make([]int, len(fields))

	for i := 0; i < len(fields); i++ {
		// Split by a forward slash and keep the first token.
		if index := bytes.IndexByte(fields[i], byte('/')); index != -1 {
			fields[i] = fields[i][:index]
		}

		text := string(fields[i])
		value, err := strconv.Atoi(text)

		if err != nil || value <= 0 {
			return ErrInvalidFace
		}

		face[i] = value - 1
	}

	patch := r.polygonSoup.NumberOfPatches() - 1
	r.polygonSoup.InsertFaceWithPatch(face, patch)

	return nil
}

// Parse a group from a line
func (r *OBJReader) parseGroup(data []byte) {
	group := bytes.TrimSpace(data[len(prefixGroup):])
	r.polygonSoup.InsertPatch(string(group))
}

// Write an OBJ file to an io.Writer interface
type OBJWriter struct {
	vertices   []geometry.Vector3
	faces      [][]int
	faceGroups []int
	lines      [][]int
	groups     []string
}

func NewOBJWriter() *OBJWriter {
	return &OBJWriter{
		vertices:   make([]geometry.Vector3, 0),
		faces:      make([][]int, 0),
		faceGroups: make([]int, 0),
		lines:      make([][]int, 0),
		groups:     make([]string, 0),
	}
}

// Set the vertices to write
func (w *OBJWriter) SetVertices(vertices []geometry.Vector3) {
	w.vertices = vertices
}

// Set the faces to write
func (w *OBJWriter) SetFaces(faces [][]int) {
	w.faces = faces
}

// Set the groups of each face to write. This must be the same length
// as the faces. Any faces that are not assigned to a group must specify -1.
func (w *OBJWriter) SetFaceGroups(faceGroups []int) {
	w.faceGroups = faceGroups
}

// Set the lines to write
func (w *OBJWriter) SetLines(lines [][]int) {
	w.lines = lines
}

// Set the groups to write
func (w *OBJWriter) SetGroups(groups []string) {
	w.groups = groups
}

// Write the mesh to the io.Writer interface
func (w *OBJWriter) Write(writer io.Writer) error {
	buffer := bufio.NewWriter(writer)

	if err := w.writeVertices(buffer); err != nil {
		return err
	}

	if err := w.writeLines(buffer); err != nil {
		return err
	}

	if err := w.writeFaces(buffer); err != nil {
		return err
	}

	return buffer.Flush()
}

// Write the vertices to the buffer
func (w *OBJWriter) writeVertices(buffer *bufio.Writer) error {
	for _, v := range w.vertices {
		entry := fmt.Sprintf("v %f %f %f\n", v[0], v[1], v[2])

		if _, err := buffer.WriteString(entry); err != nil {
			return err
		}
	}

	return nil
}

// Write the lines to the buffer
func (w *OBJWriter) writeLines(buffer *bufio.Writer) error {
	for _, l := range w.lines {
		if _, err := buffer.WriteString("l"); err != nil {
			return err
		}

		for _, v := range l {
			entry := fmt.Sprintf(" %d", v+1)

			if _, err := buffer.WriteString(entry); err != nil {
				return err
			}
		}

		if _, err := buffer.WriteString("\n"); err != nil {
			return err
		}
	}

	return nil
}

// Write the faces to the buffer
func (w *OBJWriter) writeFaces(buffer *bufio.Writer) error {
	// Index the faces for each group
	groupFaces := make(map[int][]int)

	for i, group := range w.faceGroups {
		if faces, ok := groupFaces[group]; ok {
			groupFaces[group] = append(faces, i)
		} else {
			groupFaces[group] = []int{i}
		}
	}

	// If no faces were indexed, then no faces are assigned to groups. In this case,
	// assign all faces to group "-1"
	if len(groupFaces) == 0 {
		faceGroups := make([]int, len(w.faces))

		for i := 0; i < len(w.faces); i++ {
			faceGroups[i] = i
		}

		groupFaces[-1] = faceGroups
	}

	// Write the faces for each group starting with group "-1" which is all
	// of the faces that are not assigned to a group.
	for i := -1; i < len(w.groups); i++ {
		if faces, ok := groupFaces[i]; ok {
			if i >= 0 {
				entry := fmt.Sprintf("g %s\n", w.groups[i])

				if _, err := buffer.WriteString(entry); err != nil {
					return err
				}
			}

			for _, j := range faces {
				if _, err := buffer.WriteString("f"); err != nil {
					return err
				}

				for _, v := range w.faces[j] {
					entry := fmt.Sprintf(" %d", v+1)

					if _, err := buffer.WriteString(entry); err != nil {
						return err
					}
				}

				if _, err := buffer.WriteString("\n"); err != nil {
					return err
				}
			}
		}
	}

	return nil
}
