package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gotk3/gotk3/cairo"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

const (
	WIDTH  = 400
	HEIGHT = 400
	SIZE   = 5
)

type Buffer struct {
	data    [][]bool
	blocked int
	mu      sync.Mutex
}

func (b *Buffer) Update(newData [][]bool) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.data = newData
}

func (b *Buffer) Get() [][]bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.data
}

func (b *Buffer) ToggleCell(x, y int) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if x >= 0 && x < len(b.data[0]) && y >= 0 && y < len(b.data) {
		b.data[y][x] = !b.data[y][x]
	}
}

func (b *Buffer) NextGeneration() {
	b.mu.Lock()
	defer b.mu.Unlock()

	newData := make([][]bool, len(b.data))
	for i := range newData {
		newData[i] = make([]bool, len(b.data[i]))
	}

	for y := range b.data {
		for x := range b.data[y] {
			neighbors := b.countNeighbors(x, y)
			if b.data[y][x] {
				newData[y][x] = neighbors == 2 || neighbors == 3
			} else {
				newData[y][x] = neighbors == 3
			}
		}
	}

	b.data = newData
}

func (b *Buffer) countNeighbors(x, y int) int {
	count := 0
	for dy := -1; dy <= 1; dy++ {
		for dx := -1; dx <= 1; dx++ {
			if dx == 0 && dy == 0 {
				continue
			}
			nx, ny := x+dx, y+dy
			if nx >= 0 && nx < len(b.data[0]) && ny >= 0 && ny < len(b.data) && b.data[ny][nx] {
				count++
			}
		}
	}
	return count
}

var (
	buffer  *Buffer
	surface *cairo.Surface
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
	button1, _ := gtk.ButtonNewWithLabel("Button 1")
	button2, _ := gtk.ButtonNewWithLabel("Button 2")
	button3, _ := gtk.ButtonNewWithLabel("Button 3")

	panel.PackStart(button1, false, false, 0)
	panel.PackStart(button2, false, false, 0)
	panel.PackStart(button3, false, false, 0)

	canvas, err := gtk.DrawingAreaNew()
	if err != nil {
		log.Fatal("Unable to create drawing area:", err)
	}

	canvas.SetSizeRequest(WIDTH, HEIGHT)
	canvas.AddEvents(int(gdk.BUTTON_PRESS_MASK | gdk.BUTTON_RELEASE_MASK | gdk.POINTER_MOTION_MASK))

	buffer = &Buffer{data: make([][]bool, HEIGHT)}
	for i := range buffer.data {
		buffer.data[i] = make([]bool, WIDTH)
	}

	hbox.PackStart(canvas, true, true, 0)

	// Connect button signals
	button1.Connect("clicked", func() {
		fmt.Println("Button 1 clicked")
	})
	button2.Connect("clicked", func() {
		fmt.Println("Button 2 clicked")
	})
	button3.Connect("clicked", func() {
		fmt.Println("Button 3 clicked")
	})

	canvas.Connect("configure-event", func(da *gtk.DrawingArea, event *gdk.Event) {
		if surface != nil {
			surface.Close()
		}
		win, err := da.GetWindow()
		if err != nil {
			log.Fatal("unable get window", err)
		}
		width := da.GetAllocatedWidth()
		height := da.GetAllocatedHeight()
		surface, err = win.CreateSimilarSurface(cairo.CONTENT_COLOR, width, height)
		if err != nil {
			log.Fatal("unable get surface", err)
		}

		// Update the buffer size
		buffer.mu.Lock()
		var orgData = buffer.data

		buffer.data = make([][]bool, height/SIZE)
		for i := range buffer.data {
			buffer.data[i] = make([]bool, width/SIZE)
			if i < len(orgData) {
				for j := range orgData[i] {
					if j < len(buffer.data[i]) {
						buffer.data[i][j] = orgData[i][j]
					}
				}
			}
		}
		buffer.mu.Unlock()

		// Redraw the surface
		updateSurface()
	})

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

	win.Connect("destroy", gtk.MainQuit)
	win.ShowAll()

	// Start the background updater
	go updateBufferInBackground(canvas)

	gtk.Main()
}

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
	if bufferY >= 0 && bufferX >= 0 && bufferY < len(buffer.data) && bufferX < len(buffer.data[0]) {
		buffer.data[bufferY][bufferX] = true
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

func updateBufferInBackground(drawingArea *gtk.DrawingArea) {
	for {
		time.Sleep(100 * time.Millisecond)
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
