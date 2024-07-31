package main

import (
	"log"

	"github.com/gotk3/gotk3/cairo"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

func canvasCreate() *gtk.DrawingArea {

	canvas, err := gtk.DrawingAreaNew()
	if err != nil {
		log.Fatal("Unable to create drawing area:", err)
	}

	canvas.SetSizeRequest(WIDTH, HEIGHT)
	canvas.AddEvents(int(gdk.BUTTON_PRESS_MASK | gdk.POINTER_MOTION_MASK))
	return canvas

}

func canvasConnect(canvas *gtk.DrawingArea) {
	canvas.Connect("configure-event", canvasConfigure)

	canvas.Connect("button-press-event", func(da *gtk.DrawingArea, event *gdk.Event) bool {
		buttonEvent := gdk.EventButtonNewFromEvent(event)
		handleMousePress(da, buttonEvent)
		return false
	})

	canvas.Connect("motion-notify-event", func(da *gtk.DrawingArea, event *gdk.Event) bool {
		motionEvent := gdk.EventMotionNewFromEvent(event)
		handleMouseMotion(da, motionEvent)
		return false
	})

	canvas.Connect("draw", func(da *gtk.DrawingArea, cr *cairo.Context) {
		if surface != nil {
			cr.SetSourceSurface(surface, 0, 0)
			cr.Paint()
		}
	})
}

func canvasConfigure(canvas *gtk.DrawingArea, event *gdk.Event) {
	if surface != nil {
		surface.Close()
	}
	win, err := canvas.GetWindow()
	if err != nil {
		log.Fatal("unable get window", err)
	}
	width := canvas.GetAllocatedWidth()
	height := canvas.GetAllocatedHeight()
	surface, err = win.CreateSimilarSurface(cairo.CONTENT_COLOR, width, height)
	if err != nil {
		log.Fatal("unable get surface", err)
	}

	// Update the buffer size
	buffers.mu().Lock()
	buffers.current().changeSizeNotDestructing(uint(width/SIZE), uint(height/SIZE))
	buffers.mu().Unlock()

	// Redraw the surface
	updateSurface()
}
