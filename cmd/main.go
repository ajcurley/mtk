package main

import (
	"fmt"
	"os"
	"runtime/pprof"
	"time"

	"github.com/ajcurley/mtk"
)

func main() {
	path := "/Users/acurley/projects/cfd/geometry/car.obj.gz"
    mesh, _ := mtk.NewHEMeshFromOBJFile(path)
    bounds := mesh.GetBounds()
    octree := mtk.NewOctree(bounds)

	file, _ := os.Create("cmd/ray_trace.pprof")
	defer file.Close()
	pprof.StartCPUProfile(file)
	defer pprof.StopCPUProfile()

    buildStart := time.Now()

    for i := 0; i < mesh.GetNumberOfFaces(); i++ {
        vertices := mesh.GetFaceVertices(i)
        face := mtk.NewTriangle(
            mesh.GetVertex(vertices[0]).Origin,
            mesh.GetVertex(vertices[1]).Origin,
            mesh.GetVertex(vertices[2]).Origin,
        ) 

        octree.Insert(face)
    }   

    buildTime := time.Now().Sub(buildStart).Milliseconds()
    fmt.Printf("Octree indexed in %dms\n", buildTime)

	fmt.Println(mesh.GetNumberOfFaces())
    queries := make([]mtk.IntersectsAABB, 100000)

    for i := 0; i < len(queries); i++ {
        face := mesh.GetFace(i)
        halfEdge := mesh.GetHalfEdge(face.HalfEdge)
        normal := mesh.GetFaceNormal(i).Unit()

        queries[i] = mtk.Ray{
            Origin:    mesh.GetVertex(halfEdge.Origin).Origin.Add(normal.MulScalar(0.001)),
            Direction: normal,
        }   
    }   

    queryStart := time.Now()

   	octree.QueryMany(queries)

    queryTime := time.Now().Sub(queryStart).Milliseconds()
    fmt.Printf("Octree queried in %dms\n", queryTime)
}
