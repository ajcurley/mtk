package mtk

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

	if IsGZIP(buffer) {
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

	vertex := Vector3(values)
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

	patch := r.polygonSoup.GetNumberOfPatches() - 1
	r.polygonSoup.InsertFaceWithPatch(face, patch)

	return nil
}

// Parse a group from a line
func (r *OBJReader) parseGroup(data []byte) {
	group := bytes.TrimSpace(data[len(prefixGroup):])
	r.polygonSoup.InsertPatch(string(group))
}

// Check if a bufio.Reader is a GZIP compressed file
func IsGZIP(reader *bufio.Reader) bool {
	data, err := reader.Peek(2)
	if err != nil {
		return false
	}

	return data[0] == 31 && data[1] == 139
}
