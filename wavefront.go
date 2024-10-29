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
	PrefixVertex = "v"
	PrefixFace   = "f"
	PrefixGroup  = "g"
)

var (
	ErrInvalidVertex = errors.New("invalid vertex")
	ErrInvalidFace   = errors.New("invalid face")
)

type OBJReader struct {
	vertices    []Vector3
	faces       []int
	faceOffsets []int
	faceGroups  []int
	groups      []string
}

func NewOBJReader() *OBJReader {
	return &OBJReader{
		vertices:    make([]Vector3, 0),
		faces:       make([]int, 0),
		faceOffsets: make([]int, 0),
		faceGroups:  make([]int, 0),
		groups:      make([]string, 0),
	}
}

// Read an OBJ file from an io.Reader interface
func (r *OBJReader) Read(reader io.Reader) error {
	count := 1
	buffer := bufio.NewReader(reader)

	if IsGZIP(buffer) {
		gzipFile, err := gzip.NewReader(buffer)
		if err != nil {
			return err
		}
		defer gzipFile.Close()

		buffer = bufio.NewReader(gzipFile)
	}

	for {
		data, err := buffer.ReadBytes('\n')

		if errors.Is(err, io.EOF) {
			return nil
		}

		data = bytes.TrimSpace(data)
		prefix := r.parsePrefix(data)

		switch string(prefix) {
		case PrefixVertex:
			err = r.parseVertex(data)
		case PrefixFace:
			err = r.parseFace(data)
		case PrefixGroup:
			r.parseGroup(data)
		}

		if err != nil {
			return fmt.Errorf("line %d: %v", count, err)
		}

		count++
	}
}

// Read an OBJ file from path
func (r *OBJReader) ReadFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
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
	fields := bytes.Fields(data[len(PrefixVertex):])

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
	r.vertices = append(r.vertices, vertex)

	return nil
}

// Parse a face from a line
func (r *OBJReader) parseFace(data []byte) error {
	fields := bytes.Fields(data[len(PrefixFace):])

	if len(fields) <= 2 {
		return ErrInvalidFace
	}

	faceOffset := len(r.faces)
	faceGroup := len(r.groups) - 1

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

		r.faces = append(r.faces, value-1)
	}

	r.faceOffsets = append(r.faceOffsets, faceOffset)
	r.faceGroups = append(r.faceGroups, faceGroup)

	return nil
}

// Parse a group from a line
func (r *OBJReader) parseGroup(data []byte) {
	group := bytes.TrimSpace(data[len(PrefixGroup):])
	r.groups = append(r.groups, string(group))
}

// Check if a bufio.Reader is a GZIP compressed file
func IsGZIP(reader *bufio.Reader) bool {
	data, err := reader.Peek(2)
	if err != nil {
		return false
	}

	return data[0] == 31 && data[1] == 139
}
