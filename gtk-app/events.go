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
	var bufferX = int(x / SIZE)
	var bufferY = int(y / SIZE)
	buffer.mu.Lock()
	if bufferY >= 0 && bufferX >= 0 && bufferY < len(*buffer.data) && bufferX < len((*buffer.data)[0]) {
		if drawOrInvert {
			(*buffer.data)[bufferY][bufferX] = true
		} else {
			(*buffer.data)[bufferY][bufferX] = !(*buffer.data)[bufferY][bufferX]
		}
	}

	if buffer.blocked < 4 {
		buffer.blocked += 2
	}
	buffer.mu.Unlock()

	// Schedule a redraw on the main GTK thread
	glib.IdleAdd(func() {
		updateSurface()
		da.QueueDraw()
	})
}
