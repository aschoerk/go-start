package main

import (
	"time"

	"github.com/g3n/engine/app"
	"github.com/g3n/engine/camera"
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/geometry"
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/light"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/renderer"
	"github.com/g3n/engine/window"
)

func main() {
	// Create application and scene
	a := app.App()
	scene := core.NewNode()

	// Create camera
	cam := camera.New(1)
	cam.SetPosition(5, 5, 15)
	cam.LookAt(&math32.Vector3{0, 0, 0}, &math32.Vector3{0, 1, 0})
	scene.Add(cam)

	// Create and add a light source
	light1 := light.NewDirectional(&math32.Color{1, 1, 1}, 1.0)
	light1.SetPosition(10, 10, 10)
	scene.Add(light1)

	// Create and add the second light source orthogonal to the first one
	light2 := light.NewDirectional(&math32.Color{1, 1, 1}, 1.0)
	light2.SetPosition(-10, 10, 10)
	scene.Add(light2)

	// Create grid
	grid := createGrid(10, 10, 10)
	scene.Add(grid)

	// Add coordinate axes
	addCoordinateAxes(scene)

	// Set up orbit control for rotation
	camera.NewOrbitControl(cam)

	// Handle window resize
	onResize := func(evname string, ev interface{}) {
		// Get framebuffer size and update viewport accordingly
		width, height := a.GetSize()
		a.Gls().Viewport(0, 0, int32(width), int32(height))
		cam.SetAspect(float32(width) / float32(height))
	}
	a.Subscribe(window.OnWindowSize, onResize)
	onResize("", nil) // Call once to set initial size

	// Run the application
	a.Run(func(renderer *renderer.Renderer, deltaTime time.Duration) {
		a.Gls().Clear(gls.DEPTH_BUFFER_BIT | gls.STENCIL_BUFFER_BIT | gls.COLOR_BUFFER_BIT)
		renderer.Render(scene, cam)
	})
}

func createGrid(width, height, depth int) *core.Node {
	grid := core.NewNode()

	halfW, halfH, halfD := float32(width)/2, float32(height)/2, float32(depth)/2

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			for z := 0; z < depth; z++ {
				if (x+y+z)%2 == 0 {
					// Create a cube for each grid position
					cube := geometry.NewCube(0.5)
					mat := material.NewStandard(math32.NewColor("white"))
					mesh := graphic.NewMesh(cube, mat)
					mesh.SetPosition(float32(x)-halfW+0.5, float32(y)-halfH+0.5, float32(z)-halfD+0.5)

					// Set initial boolean value (you can modify this logic)
					if (x+y+z)%2 == 0 {
						mesh.SetMaterial(material.NewStandard(math32.NewColor("blue")))
					}

					grid.Add(mesh)
				}
			}
		}
	}

	return grid
}

func addCoordinateAxes(scene *core.Node) {
	axisLength := float32(5.0)

	// X-axis (red)
	xAxis := createLine(math32.NewVector3(0, 0, 0), math32.NewVector3(axisLength, 0, 0), math32.NewColor("yellow"))
	scene.Add(xAxis)

	// Y-axis (green)
	yAxis := createLine(math32.NewVector3(0, 0, 0), math32.NewVector3(0, axisLength, 0), math32.NewColor("red"))
	scene.Add(yAxis)

	// Z-axis (blue)
	zAxis := createLine(math32.NewVector3(0, 0, 0), math32.NewVector3(0, 0, axisLength), math32.NewColor("green"))
	scene.Add(zAxis)

}

func createLine(start, end *math32.Vector3, color *math32.Color) *graphic.Lines {
	geom := geometry.NewGeometry()
	vertices := math32.NewArrayF32(0, 6)
	vertices.AppendVector3(start)
	vertices.AppendVector3(end)
	geom.AddVBO(gls.NewVBO(vertices).AddAttrib(gls.VertexPosition))

	mat := material.NewStandard(color)
	lines := graphic.NewLines(geom, mat)
	return lines
}
