package main

// #include <X11/Xlib.h>
// #include <cairo/cairo-xlib.h>
import "C"

type GlobalInfo struct {
	button, eventProperty uint
	xRoot, yRoot          int
	window                C.Window
	attributes            C.XWindowAttributes
}

type MiteWMConfig struct {
	BackgroundImageFile string `json:"background_image_file"`
}

type Client struct {
	window                              [CLIENT_WIDGETS]C.Window
	surface                             [1]*C.cairo_surface_t
	cr                                  [1]*C.cairo_t
	title                               string
	localBorderWidth, localBorderHeight int
}

type Background struct {
	window     C.Window
	surface    *C.cairo_surface_t
	image      *C.cairo_surface_t
	cr         *C.cairo_t
	imageScale C.double
}
