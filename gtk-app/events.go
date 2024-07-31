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
	buffers.mu().Lock()
	buffer := buffers.current()
	if bufferY >= 0 && bufferX >= 0 && bufferY < buffer.maxY() && bufferX < buffer.maxX() {
		if drawOrInvert {
			buffer.set(bufferX, bufferY, true)
		} else {
			buffer.set(bufferX, bufferY, !buffer.get(bufferX, bufferY))
		}
	}

	if blocked < 4 {
		blocked += 2
	}
	buffers.mu().Unlock()

	// Schedule a redraw on the main GTK thread
	glib.IdleAdd(func() {
		updateSurface()
		da.QueueDraw()
	})
}
