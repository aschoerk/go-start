package main

import "github.com/gotk3/gotk3/cairo"

const (
	WIDTH              = 400
	HEIGHT             = 400
	SIZE               = 5
	MAX_BUFFER_HISTORY = 100
)

var (
	buffers      BooleanBuffers
	surface      *cairo.Surface
	doconway     bool = false
	drawOrInvert bool = true
	blocked      int  = 0
)
