package main

import (
	"log"
	"time"

	"aschoerk.de/gtk-conway/conway"
	"github.com/gotk3/gotk3/cairo"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

func main() {
	gtk.Init(nil)

	var win, canvas = createWindow()

	buffers = conway.InitBuffers(MAX_BUFFER_HISTORY, WIDTH, HEIGHT)

	doconway = false
	// Connect button signals
	canvasConnect(canvas)

	win.Connect("destroy", gtk.MainQuit)
	win.ShowAll()

	// Start the background updater
	go updateBufferInBackground(canvas)

	gtk.Main()
}

func createWindow() (*gtk.Window, *gtk.DrawingArea) {

	win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		log.Fatal("Unable to create window:", err)
	}
	win.SetTitle("GTK Conway Example")
	win.SetDefaultSize(WIDTH, HEIGHT)

	// Create a vertical box to hold the panel and canvas
	vbox, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 5)
	if err != nil {
		log.Fatal("Unable to create vbox:", err)
	}
	win.Add(vbox)

	// Create a horizontal box for the panel and canvas
	hbox, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 5)
	if err != nil {
		log.Fatal("Unable to create hbox:", err)
	}
	vbox.PackStart(hbox, true, true, 0)

	panel, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 5)
	if err != nil {
		log.Fatal("Unable to create panel:", err)
	}
	hbox.PackStart(panel, false, false, 0)
	// Create buttons and add them to the panel
	runButton, _ := gtk.ButtonNewWithLabel("start conway")
	penButton, _ := gtk.ButtonNewWithLabel("toggle boxes")
	prevInHistoryButton, _ := gtk.ButtonNewWithLabel("prev")
	nextInHistoryButton, _ := gtk.ButtonNewWithLabel("next")

	panel.PackStart(runButton, false, false, 0)
	panel.PackStart(penButton, false, false, 0)
	panel.PackStart(prevInHistoryButton, false, false, 0)
	panel.PackStart(nextInHistoryButton, false, false, 0)
	var canvas = canvasCreate()
	hbox.PackStart(canvas, true, true, 0)

	runButton.Connect("clicked", func() {
		doconway = !doconway
		if doconway {
			runButton.SetLabel("stop conway")
		} else {
			runButton.SetLabel("start conway")
		}
	})

	penButton.Connect("clicked", func() {
		drawOrInvert = !drawOrInvert
		if drawOrInvert {
			penButton.SetLabel("toggle boxes")
		} else {
			penButton.SetLabel("draw boxes")
		}
	})

	prevInHistoryButton.Connect("clicked", func() {
		if buffers.Prev() {
			glib.IdleAdd(func() {
				updateSurface()
				canvas.QueueDraw()
			})
		}
	})

	nextInHistoryButton.Connect("clicked", func() {
		if buffers.Next() {
			glib.IdleAdd(func() {
				updateSurface()
				canvas.QueueDraw()
			})
		}
	})

	return win, canvas
}

func updateBufferInBackground(drawingArea *gtk.DrawingArea) {
	for {
		time.Sleep(100 * time.Millisecond)
		if doconway {

			if blocked > 0 {
				buffers.Mu().Lock()
				blocked -= 1
				buffers.Mu().Unlock()
			} else {
				buffers.NextGeneration()
			}

			glib.IdleAdd(func() {
				updateSurface()
				drawingArea.QueueDraw()
			})
		}

	}
}

func updateSurface() {
	if surface == nil {
		return
	}

	cr := cairo.Create(surface)
	defer cr.Close()

	data := buffers.Current()

	cr.SetSourceRGB(1, 1, 1) // White background
	cr.Paint()

	cr.SetSourceRGB(0, 0, 0) // Black for drawing
	for i := uint(0); i < data.MaxX(); i++ {
		for j := uint(0); j < data.MaxY(); j++ {
			val := data.Get(i, j)
			if val {
				x := float64(i) * SIZE
				y := float64(j) * SIZE
				cr.Rectangle(x, y, SIZE, SIZE)
				cr.Fill()
			}
		}
	}
}
