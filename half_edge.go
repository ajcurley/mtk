package mtk

import (
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"slices"
	"strings"
)

var (
	ErrNonManifoldMesh = errors.New("non-manifold mesh")
	ErrOpenMesh        = errors.New("mesh must be closed")
)

type HEMesh struct {
	vertices  []HEVertex
	faces     []HEFace
	halfEdges []HEHalfEdge
	patches   []HEPatch
}

// Construct a half edge mesh from a PolygonSoup
func NewHEMeshFromPolygonSoup(soup *PolygonSoup) (*HEMesh, error) {
	nVertices := soup.NumberOfVertices()
	nFaces := soup.NumberOfFaces()
	nPatches := soup.NumberOfPatches()

	mesh := HEMesh{
		vertices:  make([]HEVertex, 0, nVertices),
		faces:     make([]HEFace, 0, nFaces),
		halfEdges: make([]HEHalfEdge, 0, 3*nFaces),
		patches:   make([]HEPatch, 0),
	}

	// Index the patches. Each face will be assigned to the patch when the
	// faces are indexed.
	for pi := 0; pi < nPatches; pi++ {
		patch := HEPatch{
			Name: soup.Patch(pi),
		}
		mesh.patches = append(mesh.patches, patch)
	}

	// Index the vertices without reference to their originating half edge
	// which will be indexed later during the face insertion.
	indexHalfEdges := make(map[[2]int][]int)

	for vi := 0; vi < nVertices; vi++ {
		vertex := HEVertex{
			Origin: soup.Vertex(vi),
		}
		mesh.vertices = append(mesh.vertices, vertex)
	}

	// Index the faces and all edges of each face as half edges. Each half
	// edge will not have a reference to its twin until later.
	for fi := 0; fi < nFaces; fi++ {
		faceVertices := soup.Face(fi)
		nFaceVertices := len(faceVertices)
		nHalfEdges := len(mesh.halfEdges)

		face := HEFace{
			HalfEdge: nHalfEdges,
			Patch:    soup.FacePatch(fi),
		}
		mesh.faces = append(mesh.faces, face)

		for hi := 0; hi < nFaceVertices; hi++ {
			halfEdge := HEHalfEdge{
				Origin: faceVertices[hi],
				Face:   fi,
				Prev:   nHalfEdges + (hi+nFaceVertices-1)%nFaceVertices,
				Next:   nHalfEdges + (hi+nFaceVertices+1)%nFaceVertices,
				Twin:   -1,
			}

			mesh.halfEdges = append(mesh.halfEdges, halfEdge)

			// Update the originating vertex's half edge
			vertex := mesh.vertices[halfEdge.Origin]
			vertex.HalfEdge = nHalfEdges + hi
			mesh.vertices[halfEdge.Origin] = vertex

			// Index the half edges by the vertices defining the edge. The vertices
			// are sorted since the face orientation does not matter when assigning
			// twin half edges.
			x := faceVertices[hi]
			y := faceVertices[(hi+1)%nFaceVertices]
			k := [2]int{min(x, y), max(x, y)}

			if shared, ok := indexHalfEdges[k]; ok {
				if len(shared) == 2 {
					return nil, ErrNonManifoldMesh
				}

				indexHalfEdges[k] = append(shared, hi+nHalfEdges)
			} else {
				indexHalfEdges[k] = []int{hi + nHalfEdges}
			}
		}
	}

	for _, shared := range indexHalfEdges {
		if len(shared) == 2 {
			for i, index := range shared {
				halfEdge := mesh.halfEdges[index]
				halfEdge.Twin = shared[(i+1)%2]
				mesh.halfEdges[index] = halfEdge
			}
		}
	}

	return &mesh, nil
}

// Construct a half edge mesh from an OBJ file reader
func NewHEMeshFromOBJ(reader io.Reader) (*HEMesh, error) {
	objReader := NewOBJReader()
	soup, err := objReader.Read(reader)

	if err != nil {
		return nil, err
	}

	return NewHEMeshFromPolygonSoup(soup)
}

// Construct a half edge mesh from an OBJ file
func NewHEMeshFromOBJFile(path string) (*HEMesh, error) {
	objReader := NewOBJReader()
	soup, err := objReader.ReadFile(path)

	if err != nil {
		return nil, err
	}

	return NewHEMeshFromPolygonSoup(soup)
}

// Compute the axis-aligned bounding box
func (m *HEMesh) Bounds() AABB {
	minBound := Vector3{1, 1, 1}.MulScalar(math.Inf(1))
	maxBound := Vector3{1, 1, 1}.MulScalar(math.Inf(-1))

	for _, vertex := range m.vertices {
		for i := 0; i < 3; i++ {
			if vertex.Origin[i] < minBound[i] {
				minBound[i] = vertex.Origin[i]
			}

			if vertex.Origin[i] > maxBound[i] {
				maxBound[i] = vertex.Origin[i]
			}
		}
	}

	center := maxBound.Add(minBound).MulScalar(0.5)
	halfSize := maxBound.Sub(minBound).MulScalar(0.5)

	return NewAABB(center, halfSize)
}

// Check if the half edge mesh is closed (no open boundaries)
func (m *HEMesh) IsClosed() bool {
	for _, halfEdge := range m.halfEdges {
		if halfEdge.IsBoundary() {
			return false
		}
	}

	return true
}

// Check if the half edge mesh has consistently oriented faces
func (m *HEMesh) IsConsistent() bool {
	for _, halfEdge := range m.halfEdges {
		if !halfEdge.IsBoundary() {
			twin := m.HalfEdge(halfEdge.Twin)

			if twin.Origin == halfEdge.Origin {
				return false
			}
		}
	}

	return true
}

// Get the number of vertices
func (m *HEMesh) NumberOfVertices() int {
	return len(m.vertices)
}

// Get the vertex by ID
func (m *HEMesh) Vertex(id int) HEVertex {
	return m.vertices[id]
}

// Get the neighboring vertices for a vertex by ID
func (m *HEMesh) VertexNeighbors(id int) []int {
	neighbors := make([]int, 0)
	vertex := m.Vertex(id)
	current := vertex.HalfEdge

	for {
		halfEdge := m.HalfEdge(current)

		if halfEdge.IsBoundary() {
			panic("vertex neighbors requires a closed mesh")
		}

		if halfEdge.Origin == id {
			neighbor := m.HalfEdge(halfEdge.Next).Origin
			neighbors = append(neighbors, neighbor)
		} else {
			neighbors = append(neighbors, halfEdge.Origin)
		}

		twin := m.HalfEdge(halfEdge.Twin)

		if twin.Origin != id {
			current = twin.Next
		} else {
			current = twin.Prev
		}

		if current == vertex.HalfEdge {
			break
		}
	}

	return neighbors
}

// Get the faces using the vertex by ID
func (m *HEMesh) VertexFaces(id int) []int {
	faces := make([]int, 0)
	vertex := m.Vertex(id)
	current := vertex.HalfEdge

	for {
		halfEdge := m.HalfEdge(current)

		if halfEdge.IsBoundary() {
			panic("vertex faces requires a closed mesh")
		}

		faces = append(faces, halfEdge.Face)
		twin := m.HalfEdge(halfEdge.Twin)

		if twin.Origin != id {
			current = twin.Next
		} else {
			current = twin.Prev
		}

		if current == vertex.HalfEdge {
			break
		}
	}

	return faces
}

// Get the Gaussian curvature at the vertex by ID. This assumes the mesh
// is composed of strictly triangular elements and is oriented.
func (m *HEMesh) VertexCurvature(id int) (float64, error) {
	var angle, area float64
	angle = 2 * math.Pi

	vertex := m.Vertex(id)
	current := vertex.HalfEdge

	for {
		halfEdge := m.HalfEdge(current)
		nextHalfEdge := m.HalfEdge(halfEdge.Next)
		prevHalfEdge := m.HalfEdge(halfEdge.Prev)

		p := m.Vertex(prevHalfEdge.Origin).Origin
		q := m.Vertex(halfEdge.Origin).Origin
		r := m.Vertex(nextHalfEdge.Origin).Origin

		u := p.Sub(q)
		v := r.Sub(q)

		angle -= u.AngleTo(v)
		area += u.Cross(v).Mag() * 0.5

		if halfEdge.IsBoundary() {
			return math.NaN(), ErrOpenMesh
		}

		twin := m.HalfEdge(halfEdge.Twin)
		current = twin.Next

		if current == vertex.HalfEdge {
			break
		}
	}

	return 3 * angle / area, nil
}

// Get if the vertex is on a boundary. This assumes the mesh is consistently
// oriented. For inconsistent meshes, this may yield incorrect results.
func (m *HEMesh) IsVertexOnBoundary(id int) bool {
	vertex := m.Vertex(id)
	current := vertex.HalfEdge

	for {
		halfEdge := m.HalfEdge(current)

		if halfEdge.IsBoundary() {
			return true
		}

		current = m.HalfEdge(halfEdge.Twin).Next

		if current == vertex.HalfEdge {
			break
		}
	}

	return false
}

// Get the number of faces
func (m *HEMesh) NumberOfFaces() int {
	return len(m.faces)
}

// Get the face by ID
func (m *HEMesh) Face(id int) HEFace {
	return m.faces[id]
}

// Get the vertices defining the face by ID
func (m *HEMesh) FaceVertices(id int) []int {
	faceHalfEdges := m.FaceHalfEdges(id)
	vertices := make([]int, len(faceHalfEdges))

	for i, faceHalfEdge := range faceHalfEdges {
		vertices[i] = m.HalfEdge(faceHalfEdge).Origin
	}

	return vertices
}

// Get the neighboring faces for a face by ID
func (m *HEMesh) FaceNeighbors(id int) []int {
	faceHalfEdges := m.FaceHalfEdges(id)
	neighbors := make([]int, 0, len(faceHalfEdges))

	for _, faceHalfEdge := range faceHalfEdges {
		halfEdge := m.HalfEdge(faceHalfEdge)

		if !halfEdge.IsBoundary() {
			twin := m.HalfEdge(halfEdge.Twin)
			neighbors = append(neighbors, twin.Face)
		}
	}

	return neighbors
}

// Get the half edges of the face by ID
func (m *HEMesh) FaceHalfEdges(id int) []int {
	halfEdges := make([]int, 0)

	face := m.Face(id)
	current := face.HalfEdge

	for {
		halfEdges = append(halfEdges, current)
		current = m.HalfEdge(current).Next

		if current == face.HalfEdge {
			break
		}
	}

	return halfEdges
}

// Get the unit normal vector of the face by ID
func (m *HEMesh) FaceNormal(id int) Vector3 {
	normal := Vector3{0, 0, 0}
	vertices := m.FaceVertices(id)

	for i := 0; i < len(vertices); i++ {
		p := m.Vertex(vertices[i]).Origin
		q := m.Vertex(vertices[(i+1)%len(vertices)]).Origin
		normal = normal.Add(p.Cross(q))
	}

	return normal.Unit()
}

// Get the number of half edges
func (m *HEMesh) NumberOfHalfEdges() int {
	return len(m.halfEdges)
}

// Get the half edge by ID
func (m *HEMesh) HalfEdge(id int) HEHalfEdge {
	return m.halfEdges[id]
}

// Get the number of patches
func (m *HEMesh) NumberOfPatches() int {
	return len(m.patches)
}

// Get the patch by ID
func (m *HEMesh) Patch(id int) HEPatch {
	return m.patches[id]
}

// Get the list of patch names
func (m *HEMesh) PatchNames() []string {
	names := make([]string, m.NumberOfPatches())

	for i, patch := range m.patches {
		names[i] = patch.Name
	}

	return names
}

// Get the faces assigned to the patch by ID
func (m *HEMesh) PatchFaces(id int) []int {
	faces := make([]int, 0)

	for i, face := range m.faces {
		if face.Patch == id {
			faces = append(faces, i)
		}
	}

	return faces
}

// Get the distinct components (connected faces). Each component is
// defined by the indices of the faces.
func (m *HEMesh) Components() [][]int {
	components := make([][]int, 0)
	visited := make([]bool, m.NumberOfFaces())

	for next, isVisited := range visited {
		if !isVisited {
			component := make([]int, 0)
			queue := []int{next}

			for len(queue) > 0 {
				current := queue[0]
				queue = queue[1:]

				if !visited[current] {
					visited[current] = true
					component = append(component, current)

					for _, neighbor := range m.FaceNeighbors(current) {
						if !visited[neighbor] {
							queue = append(queue, neighbor)
						}
					}
				}
			}

			components = append(components, component)
		}
	}

	return components
}

// Get the shared vertices between two faces.
func (m *HEMesh) SharedVertices(i, j int) []Vector3 {
	index := make(map[int]struct{})
	vertices := make([]Vector3, 0)

	for _, vertex := range m.FaceVertices(i) {
		index[vertex] = struct{}{}
	}

	for _, vertex := range m.FaceVertices(j) {
		if _, ok := index[vertex]; ok {
			origin := m.Vertex(vertex).Origin
			vertices = append(vertices, origin)
		}
	}

	return vertices
}

// Get the half edge pairs with adjacent faces exceeding the angle
// threshold (in radians) between face normals.
func (m *HEMesh) FeatureEdges(threshold float64) [][2]int {
	visited := make(map[int]struct{})
	featureEdges := make([][2]int, 0)

	for i, halfEdge := range m.halfEdges {
		if _, ok := visited[i]; !ok {
			visited[i] = struct{}{}

			if !halfEdge.IsBoundary() {
				visited[halfEdge.Twin] = struct{}{}
				twin := m.halfEdges[halfEdge.Twin]

				u := m.FaceNormal(halfEdge.Face)
				v := m.FaceNormal(twin.Face)

				if u.AngleTo(v) >= threshold {
					featureEdges = append(featureEdges, [2]int{i, halfEdge.Twin})
				}
			}
		}
	}

	return featureEdges
}

// Naively copy another half edge mesh into the current. This does not
// merge any duplicate vertices or faces.
func (m *HEMesh) Merge(other *HEMesh) {
	indexPatches := make(map[string]int)

	for _, patch := range m.patches {
		indexPatches[patch.Name] = len(indexPatches)
	}

	for _, patch := range other.patches {
		if _, ok := indexPatches[patch.Name]; !ok {
			indexPatches[patch.Name] = len(indexPatches)
		}
	}

	offsetVertices := m.NumberOfVertices()
	offsetFaces := m.NumberOfFaces()
	offsetHalfEdges := m.NumberOfHalfEdges()

	m.vertices = append(m.vertices, other.vertices...)
	m.faces = append(m.faces, other.faces...)
	m.halfEdges = append(m.halfEdges, other.halfEdges...)
	m.patches = make([]HEPatch, len(indexPatches))

	for patchName, i := range indexPatches {
		m.patches[i] = HEPatch{Name: patchName}
	}

	for i, vertex := range m.vertices[offsetVertices:] {
		vertex.HalfEdge += offsetHalfEdges
		m.vertices[i+offsetVertices] = vertex
	}

	for i, face := range m.faces[offsetFaces:] {
		face.HalfEdge += offsetHalfEdges

		if face.Patch >= 0 {
			patchName := other.patches[face.Patch].Name
			face.Patch = indexPatches[patchName]
		}

		m.faces[i+offsetFaces] = face
	}

	for i, halfEdge := range m.halfEdges[offsetHalfEdges:] {
		halfEdge.Origin += offsetVertices
		halfEdge.Face += offsetFaces
		halfEdge.Prev += offsetHalfEdges
		halfEdge.Next += offsetHalfEdges

		if !halfEdge.IsBoundary() {
			halfEdge.Twin += offsetHalfEdges
		}

		m.halfEdges[i+offsetHalfEdges] = halfEdge
	}
}

// Orient the mesh such that the faces of each distinct component share
// the same normal vector orientaton. This does not guarantee that all
// components will have the same orientation.
func (m *HEMesh) Orient() {
	oriented := make([]bool, m.NumberOfFaces())

	for _, component := range m.Components() {
		queue := []int{component[0]}

		for len(queue) > 0 {
			current := queue[0]
			queue = queue[1:]

			if !oriented[current] {
				oriented[current] = true

				for _, neighbor := range m.FaceNeighbors(current) {
					if !oriented[neighbor] {
						queue = append(queue, neighbor)

						if !m.isFaceConsistent(current, neighbor) {
							m.flipFace(neighbor)
						}
					}
				}
			}
		}
	}
}

// Check if the two faces share a consistent orientation if they share
// an edge. If no edge is shared, this returns false.
func (m *HEMesh) isFaceConsistent(i, j int) bool {
	index := make(map[int]bool)

	for _, faceHalfEdge := range m.FaceHalfEdges(i) {
		index[faceHalfEdge] = true
	}

	for _, faceHalfEdge := range m.FaceHalfEdges(j) {
		halfEdge := m.HalfEdge(faceHalfEdge)

		if _, ok := index[halfEdge.Twin]; ok {
			twin := m.HalfEdge(halfEdge.Twin)
			return halfEdge.Origin != twin.Origin
		}
	}

	return false
}

// Flip the orientation of the face by reversing the half edges
func (m *HEMesh) flipFace(id int) {
	faceHalfEdges := m.FaceHalfEdges(id)
	halfEdges := make([]HEHalfEdge, len(faceHalfEdges))

	for i, faceHalfEdge := range faceHalfEdges {
		halfEdge := m.HalfEdge(faceHalfEdge)
		prev := halfEdge.Prev
		next := halfEdge.Next
		origin := m.HalfEdge(prev).Origin

		halfEdge.Next = prev
		halfEdge.Prev = next
		halfEdge.Origin = origin

		halfEdges[i] = halfEdge
	}

	for i, faceHalfEdge := range faceHalfEdges {
		m.halfEdges[faceHalfEdge] = halfEdges[i]
	}
}

// Extract a subset of the mesh by face IDs
func (m *HEMesh) ExtractFaces(ids []int) (*HEMesh, error) {
	soup := NewPolygonSoup()
	indexVertices := make(map[int]int)
	indexPatches := make(map[int]int)

	for _, originalFace := range ids {
		face := m.Face(originalFace)
		faceVertices := m.FaceVertices(originalFace)

		for i, originalVertex := range faceVertices {
			if _, ok := indexVertices[originalVertex]; !ok {
				vertex := m.Vertex(originalVertex)
				soup.InsertVertex(vertex.Origin)
				indexVertices[originalVertex] = len(indexVertices)
			}

			faceVertices[i] = indexVertices[originalVertex]
		}

		if face.Patch >= 0 {
			if _, ok := indexPatches[face.Patch]; !ok {
				patch := m.Patch(face.Patch)
				soup.InsertPatch(patch.Name)
				indexPatches[face.Patch] = len(indexPatches)
			}

			soup.InsertFaceWithPatch(faceVertices, indexPatches[face.Patch])
		} else {
			soup.InsertFace(faceVertices)
		}
	}

	return NewHEMeshFromPolygonSoup(soup)
}

// Extract a subset of the mesh by patch IDs
func (m *HEMesh) ExtractPatches(ids []int) (*HEMesh, error) {
	indexPatches := make(map[int]bool)
	faces := make([]int, 0)

	for _, i := range ids {
		indexPatches[i] = true
	}

	for i, face := range m.faces {
		if _, ok := indexPatches[face.Patch]; ok {
			faces = append(faces, i)
		}
	}

	return m.ExtractFaces(faces)
}

// Extract a subset of the mesh by patch names
func (m *HEMesh) ExtractPatchNames(names []string) (*HEMesh, error) {
	indexPatches := make(map[string]bool)
	patches := make([]int, 0)

	for _, name := range names {
		indexPatches[name] = true
	}

	for i, patch := range m.patches {
		if _, ok := indexPatches[patch.Name]; ok {
			patches = append(patches, i)
		}
	}

	return m.ExtractPatches(patches)
}

// Zip open edges by merging vertices within the geometric tolerance and
// are on edges without a twin.
func (m *HEMesh) ZipEdges() error {
	bounds := m.Bounds().Buffer(GeometricTolerance)
	octree := NewOctree(bounds)

	if !m.IsConsistent() {
		m.Orient()
	}

	vertices := make([]HEVertex, 0)
	indexLookup := make(map[int]int)
	vertexLookup := make(map[int]int)

	for i, vertex := range m.vertices {
		if m.IsVertexOnBoundary(i) {
			query := NewSphere(vertex.Origin, GeometricTolerance)
			duplicates := octree.Query(query)

			if len(duplicates) > 0 {
				index := slices.Min(duplicates)
				vertexLookup[i] = indexLookup[index]
			} else {
				indexLookup[octree.NumberOfItems()] = i
				vertexLookup[i] = len(vertices)
				vertices = append(vertices, vertex)
				octree.Insert(vertex.Origin)
			}
		} else {
			vertexLookup[i] = len(vertices)
			vertices = append(vertices, vertex)
		}
	}

	// Update the vertices
	m.vertices = vertices

	// Update the half edges to reference the condensed vertices
	for i, halfEdge := range m.halfEdges {
		halfEdge.Origin = vertexLookup[halfEdge.Origin]
		m.halfEdges[i] = halfEdge
	}

	// Update the half edge twins by indexing the sorted pair of
	// vertices defining the edge. If there are more than two edges
	// indexed together, the mesh is non-manifold.
	sharedEdges := make(map[[2]int][]int)

	for i, halfEdge := range m.halfEdges {
		ki := halfEdge.Origin
		kn := m.halfEdges[halfEdge.Next].Origin
		key := [2]int{min(ki, kn), max(ki, kn)}

		if shared, ok := sharedEdges[key]; ok {
			if len(shared) == 2 {
				vi := m.vertices[ki].Origin
				vn := m.vertices[kn].Origin
				loc := vi.Add(vn).MulScalar(0.5)
				return fmt.Errorf("%v: near %v", ErrNonManifoldMesh, loc)
			}

			sharedEdges[key] = append(shared, i)
		} else {
			sharedEdges[key] = []int{i}
		}
	}

	for _, shared := range sharedEdges {
		if len(shared) == 2 {
			for i := range shared {
				halfEdge := m.halfEdges[shared[i]]
				halfEdge.Twin = shared[(i+1)%2]
				m.halfEdges[shared[i]] = halfEdge
			}
		}
	}

	return nil
}

// Export the mesh to OBJ
func (m *HEMesh) ExportOBJ(w io.Writer) error {
	vertices := make([]Vector3, m.NumberOfVertices())
	faces := make([][]int, m.NumberOfFaces())
	faceGroups := make([]int, m.NumberOfFaces())
	groups := make([]string, m.NumberOfPatches())

	for i, vertex := range m.vertices {
		vertices[i] = vertex.Origin
	}

	for i, face := range m.faces {
		faces[i] = m.FaceVertices(i)
		faceGroups[i] = face.Patch
	}

	for i, patch := range m.patches {
		groups[i] = patch.Name
	}

	objWriter := NewOBJWriter()
	objWriter.SetVertices(vertices)
	objWriter.SetFaces(faces)
	objWriter.SetFaceGroups(faceGroups)
	objWriter.SetGroups(groups)

	return objWriter.Write(w)
}

// Export the mesh to an OBJ file
func (m *HEMesh) ExportOBJFile(path string) error {
	var writer io.Writer

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	writer = file

	if strings.HasSuffix(strings.ToLower(path), ".gz") {
		gzipFile := gzip.NewWriter(file)
		defer gzipFile.Close()
		writer = gzipFile
	}

	return m.ExportOBJ(writer)
}

// Half edge mesh vertex
type HEVertex struct {
	Origin   Vector3
	HalfEdge int
}

// Half edge mesh face
type HEFace struct {
	HalfEdge int
	Patch    int
}

// Half edge mesh half edge
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

// Half edge mesh patch
type HEPatch struct {
	Name string
}
