package main

import (
	"github.com/ajcurley/mtk"
)

func main() {
	path := "/Users/acurley/projects/cfd/geometry/car.obj.gz"
	mesh, _ := mtk.NewHEMeshFromOBJFile(path)
	bounds := mesh.GetBounds()
	octree := mtk.NewOctree(bounds)

	for i := 0; i < mesh.GetNumberOfFaces(); i++ {
		vertices := mesh.GetFaceVertices(i)
		face := mtk.NewTriangle(
			mesh.GetVertex(vertices[0]).Origin,
			mesh.GetVertex(vertices[1]).Origin,
			mesh.GetVertex(vertices[2]).Origin,
		)

		octree.Insert(face)
	}

	box := mtk.NewAABB(
		mtk.NewVector3(-2, -0.5, 0.5),
		mtk.NewVector3(0.5, 0.5, 0.5),
	)

	faces := octree.Query(box)
	subset, _ := mesh.ExtractFaces(faces)
	subset.ExportOBJFile("/Users/acurley/Desktop/clipped.obj")
}
