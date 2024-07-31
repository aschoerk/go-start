package main

import (
	"log"
	"time"

	"github.com/gotk3/gotk3/cairo"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

const (
	WIDTH  = 400
	HEIGHT = 400
	SIZE   = 5
)

var (
	buffer   *Buffer
	surface  *cairo.Surface
	doconway bool
)

func main() {
	gtk.Init(nil)

	win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		log.Fatal("Unable to create window:", err)
	}
	win.SetTitle("GTK Double Buffering Example")
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
	button1, _ := gtk.ButtonNewWithLabel("start conway")
	button2, _ := gtk.ButtonNewWithLabel("Button 2")
	button3, _ := gtk.ButtonNewWithLabel("Button 3")

	panel.PackStart(button1, false, false, 0)
	panel.PackStart(button2, false, false, 0)
	panel.PackStart(button3, false, false, 0)
	var canvas = canvasCreate()

	buffer = &Buffer{data: make([][]bool, HEIGHT)}
	for i := range buffer.data {
		buffer.data[i] = make([]bool, WIDTH)
	}

	hbox.PackStart(canvas, true, true, 0)

	doconway = false
	// Connect button signals
	button1.Connect("clicked", func() {
		doconway = !doconway
		if doconway {
			button1.SetLabel("stop conway")
		} else {
			button1.SetLabel("start conway")
		}
	})

	canvasConnect(canvas)

	win.Connect("destroy", gtk.MainQuit)
	win.ShowAll()

	// Start the background updater
	go updateBufferInBackground(canvas)

	gtk.Main()
}

func updateBufferInBackground(drawingArea *gtk.DrawingArea) {
	for {
		time.Sleep(100 * time.Millisecond)
		if doconway {

			if buffer.blocked > 0 {
				buffer.mu.Lock()
				buffer.blocked--
				buffer.mu.Unlock()
			} else {
				buffer.NextGeneration()

				// Schedule a redraw on the main GTK thread

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

	data := buffer.Get()

	cr.SetSourceRGB(1, 1, 1) // White background
	cr.Paint()

	cr.SetSourceRGB(0, 0, 0) // Black for drawing
	for i, row := range data {
		for j, val := range row {
			if val {
				x := float64(j) * SIZE
				y := float64(i) * SIZE
				cr.Rectangle(x, y, SIZE, SIZE)
				cr.Fill()
			}
		}
	}
}
