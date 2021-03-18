package main

// #include <cairo/cairo-xlib.h>
import "C"

func setBackground(file string) {
	var rootAttributes C.XWindowAttributes
	C.XGetWindowAttributes(wm_display, C.XDefaultRootWindow(wm_display), &rootAttributes)
	background.window = C.XCreateSimpleWindow(
		wm_display, C.XDefaultRootWindow(wm_display),
		rootAttributes.x, rootAttributes.y, C.uint(rootAttributes.width), C.uint(rootAttributes.height),
		0, 0, C.XBlackPixel(wm_display, 0),
	)
	background.surface = C.cairo_xlib_surface_create(
		wm_display, background.window,
		C.XDefaultVisual(wm_display, C.XDefaultScreen(wm_display)),
		rootAttributes.width, rootAttributes.height,
	)
	background.cr = C.cairo_create(background.surface)
	background.image = C.cairo_image_surface_create_from_png(C.CString(file))

	imageWidth := C.cairo_image_surface_get_width(background.image)
	imageHeight := C.cairo_image_surface_get_height(background.image)

	if rootAttributes.width/imageWidth > rootAttributes.height/imageHeight {
		background.imageScale = C.double(rootAttributes.width) / C.double(imageWidth)
	} else {
		background.imageScale = C.double(rootAttributes.height) / C.double(imageHeight)
	}

	C.cairo_set_antialias(background.cr, C.CAIRO_ANTIALIAS_SUBPIXEL)
	C.XMapWindow(wm_display, background.window)
}

func drawBackground() {
	if background.image == nil {
		return
	}
	C.cairo_save(background.cr)
	C.cairo_scale(background.cr,
		background.imageScale, background.imageScale)
	C.cairo_set_source_surface(background.cr, background.image, 0, 0)
	C.cairo_paint(background.cr)
	C.cairo_restore(background.cr)
	C.cairo_surface_flush(background.surface)
}
