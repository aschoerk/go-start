package main

import (
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

func handleMouseMotion(da *gtk.DrawingArea, event *gdk.EventMotion) {
	// Use 'da' as needed
	x, y := event.MotionVal()
	state := event.State()

	if state&gdk.BUTTON1_MASK != 0 {
		handleMouse(da, x, y)
	}

}

func handleMousePress(da *gtk.DrawingArea, event *gdk.EventButton) {
	x := event.X()
	y := event.Y()
	handleMouse(da, x, y)
}

func handleMouse(da *gtk.DrawingArea, x float64, y float64) {
	var bufferX = uint(x / SIZE)
	var bufferY = uint(y / SIZE)
	buffers.Mu().Lock()
	buffer := buffers.Current()
	if bufferY >= 0 && bufferX >= 0 && bufferY < buffer.MaxY() && bufferX < buffer.MaxX() {
		if drawOrInvert {
			buffer.Set(bufferX, bufferY, true)
		} else {
			buffer.Set(bufferX, bufferY, !buffer.Get(bufferX, bufferY))
		}
	}

	if blocked < 4 {
		blocked += 2
	}
	buffers.Mu().Unlock()

	// Schedule a redraw on the main GTK thread
	glib.IdleAdd(func() {
		updateSurface()
		da.QueueDraw()
	})
}
