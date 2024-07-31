package main

import (
	"log"
	"time"

	"github.com/gotk3/gotk3/cairo"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

func main() {
	gtk.Init(nil)

	var win, canvas = createWindow()

	tmp := make([][]bool, HEIGHT)
	buffer = &Buffer{data: &tmp}
	for i := range *buffer.data {
		(*buffer.data)[i] = make([]bool, WIDTH)
	}

	bufferHistory = make([]*[][]bool, MAX_BUFFER_HISTORY)

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
		handleHistory(canvas, -1)
	})

	nextInHistoryButton.Connect("clicked", func() {
		handleHistory(canvas, +1)
	})

	return win, canvas
}

func handleHistory(canvas *gtk.DrawingArea, dir int) {
	buffer.mu.Lock()
	defer buffer.mu.Unlock()
	var saveDoConway = doconway
	doconway = false
	var tmp = (actBufferHistoryIndex + dir) % MAX_BUFFER_HISTORY
	if tmp < 0 {
		tmp = tmp + MAX_BUFFER_HISTORY
	}
	if bufferHistory[tmp] != nil {
		bufferHistory[actBufferHistoryIndex] = buffer.data
		buffer.data = bufferHistory[tmp]
		actBufferHistoryIndex = tmp
	}
	doconway = saveDoConway
	glib.IdleAdd(func() {
		updateSurface()
		canvas.QueueDraw()
	})
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
	for i, row := range *data {
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
