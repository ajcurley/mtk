# Mesh Tool Kit (mtk)
Mesh tool kit for polygonal mesh processing. Visit [pkg.go.dev](https://pkg.go.dev/github.com/ajcurley/mtk) for a more detailed documentation.

## Quickstart
`mtk` is primarily intended for reading/writing polygonal surface meshes as well as performing minor modifications (like zipping open edges and orienting faces). Additionally, the library supports spatial indexing for fast retrieval of intersecting entities.

While a minimal polygon soup mesh is implemented, it is recommended to use the half edge mesh data structure `HEMesh`. The half edge mesh data structure only supports manifold surface meshes and will return an error upon import if a non-manifold mesh is provided. Below is an example of some common use cases of a half edge mesh.

```go
package main

import (
  "fmt"

  "github.com/ajcurley/mtk"
)

func main() {
  path := "/some/path/to/model.obj" // also supports .obj.gz
  mesh, err := mtk.NewHEMeshFromOBJFile(path)

  if err != nil {
    panic(err)
  }

  // Print out a basic summary of the contents
  fmt.Printf("Mesh summary:\n")
  fmt.Printf("Number of vertices:   %d", mesh.GetNumberOfVertices())
  fmt.Printf("Number of faces:      %d", mesh.GetNumberOfFaces())
  fmt.Printf("Number of half edges: %d", mesh.GetNumberOfHalfEdges())
  fmt.Printf("Number of patches:    %d", mesh.GetNumberOfPatches())

  // Check if the mesh has any open edges
  if !mesh.IsClosed() {
    fmt.Println("Open edges found! Fixing them now.")

    // Zip any open edges. This may result in a non-manifold mesh if three faces
    // with open edges collapse into a single edge.
    if err := mesh.ZipEdges(); err != nil {
      panic(err)
    }
  }

  // Check if the mesh has any inconsistently oriented faces
  if !mesh.IsConsistent() {
    fmt.Println("Inconsistent faces found! Fixing them now.")

    // Orient the mesh such that the faces of each connected component are the
    // same. Note: for meshes with multiple independent components, the orientation
    // of each component may be different.
    mesh.Orient()
  }

  // Extract a subset of the mesh. In this case, we subset using patch names; however,
  // we could subset by a list of face IDs.
  patchNames := []string{"patch1", "patch2"}
  submesh := mesh.ExtractPatchNames(patchNames)

  // Write the submesh to an OBJ file (also supports .obj.gz)
  submesh.ExportOBJFile("/path/to/some/output.obj")
}
```

## Spatial Indexing
`mtk` supports spatial indexing using a linear octree data structure. The `Octree` type implements three main methods: `Insert`, `Query`, and `QueryMany` among other helpful methods. `QueryMany` uses the available number of CPU by default.

One or more geometries types may be inserted into a single octree. When performing a query, only the types implementing the appropriate interface for the query geometry can be returned.

**Example: Index a triangle mesh from OBJ in an octree**
```go
func main() {
  // Import the mesh. This supports both `.obj` and `.obj.gz` extensions.
  mesh, _ := NewHEMeshFromOBJFile("./some_mesh.obj.gz")

  // Create the bounded octree
  bounds := mesh.GetBounds()
  octree := NewOctree(bounds)

  // Insert each face into the octree
  for i := 0; i < mesh.GetNumberOfFaces(); i++ {
    vertices := mesh.GetFaceVertices(i)

    // Check that the face is a triangle since the HEMesh supports polygon elements
    // but collision detection is only implemented for triangles.
    if len(vertices) != 3 {
      triangle := NewTriangle(vertices[0].Origin, vertices[1].Origin, vertices[2].Origin)
      octree.Insert(triangle)
    }
  }

  // Query for all faces intersecting the AABB. In this case, the AABB is centered
  // at the origin (0, 0, 0) and ranges from (-0.5, -0.5, -0.5) to (0.5, 0.5, 0.5)
  center := NewVector3(0, 0, 0)
  halfSize := NewVector3(0.5, 0.5, 0.5)
  query := NewAABB(center, halfSize)

  // Get the list of item IDs intersecting the AABB
  results := octree.Query(query) 
}
```

**Example: Find all points within a radius**
```go
func findPointsInside(octree *mtk.Octree, loc Vector3) []Vector3 {
  // Given an octree indexing many Vector3 items, search for all points within
  // 1e-3 distance of the a query location.
  query := NewSphere(loc, 1e-3)
  items := make([]Vector3)

  for i, index := range octree.Query(query) {
    // Cast the item to the appropriate type
    if item, ok := octree.GetItem(index).(Vector3) {
      items = append(items, item)
    }
  }

  return items
}
```

## Collision Detection
The following interfaces are used for collision detection between geometric types. The table below shows which interfaces are currently implemented by each geometric type.
- `IntersectsVector3(Vector3) bool`
- `IntersectsRay(Ray) bool`
- `IntersectsSphere(Sphere) bool`
- `IntersectsAABB(AABB) bool`
- `IntersectsTriangle(Triangle) bool`

|          | Vector3          | Ray              | Sphere           | AABB             | Triangle         |
|----------|:----------------:|:----------------:|:----------------:|:----------------:|:----------------:|
| Vector3  |                  |                  |:heavy_check_mark:|:heavy_check_mark:|                  |
| Ray      |                  |                  |                  |:heavy_check_mark:|:heavy_check_mark:|
| Sphere   |:heavy_check_mark:|                  |                  |:heavy_check_mark:|                  |
| AABB     |:heavy_check_mark:|:heavy_check_mark:|:heavy_check_mark:|:heavy_check_mark:|:heavy_check_mark:|
| Triangle |                  |:heavy_check_mark:|                  |:heavy_check_mark:|                  |

## Issues
Help improve this library by [reporting issues](https://github.com/ajcurley/mtk/issues)

## License
`mtk` is licensed under the [MIT License](./LICENSE)
