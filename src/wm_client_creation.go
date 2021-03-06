package main

// #include <X11/Xutil.h>
// #include <cairo/cairo-xlib.h>
import "C"

func (client *Client) getSizeHints(hints *C.XSizeHints) {
	var supplied C.long
	C.XGetWMNormalHints(wm_display, client.window[CLIENT_APP], hints, &supplied)
}

func (client *Client) getClassHint(hint *C.XClassHint) {
	C.XGetClassHint(wm_display, client.window[CLIENT_APP], hint)
}

func newClient(
	clientTable *ClientTable,
	targetWindow C.Window,
	lastUngrabbedApp *C.Window,
) {
	if targetWindow == C.None {
		return
	}
	if _, exist := (*clientTable)[targetWindow]; exist {
		return
	}

	var targetAttributes C.XWindowAttributes
	C.XGetWindowAttributes(wm_display, targetWindow, &targetAttributes)

	var client Client
	client.window[CLIENT_APP] = targetWindow

	C.XGrabButton(
		wm_display,
		C.AnyButton,
		C.AnyModifier,
		client.window[CLIENT_APP],
		C.False,
		C.ButtonPressMask,
		C.GrabModeAsync,
		C.GrabModeAsync,
		C.None,
		C.None,
	)

	if *lastUngrabbedApp != C.None {
		C.XGrabButton(
			wm_display,
			C.AnyButton,
			C.AnyModifier,
			*lastUngrabbedApp,
			C.False,
			C.ButtonPressMask,
			C.GrabModeAsync,
			C.GrabModeAsync,
			C.None,
			C.None,
		)
	}

	C.XSetInputFocus(wm_display, client.window[CLIENT_APP], C.RevertToNone, C.CurrentTime)

	client.localBorderWidth = WIDTH_DIFF + int(targetAttributes.border_width)*2
	client.localBorderHeight = HEIGHT_DIFF + int(targetAttributes.border_width)*2

	boxWidth := client.localBorderWidth + int(targetAttributes.width)
	boxHeight := client.localBorderHeight + int(targetAttributes.height)

	var boxVisualInfo C.XVisualInfo
	C.XMatchVisualInfo(wm_display, C.XDefaultScreen(wm_display), 32, C.TrueColor, &boxVisualInfo)

	var boxAttributes C.XSetWindowAttributes
	boxAttributes.colormap = C.XCreateColormap(wm_display, C.XDefaultRootWindow(wm_display), boxVisualInfo.visual, C.AllocNone)
	boxAttributes.override_redirect = 1

	// BOXを作成
	client.window[CLIENT_BOX] = C.XCreateWindow(
		wm_display,
		C.XDefaultRootWindow(wm_display),
		targetAttributes.x-CONFIG_BOX_BORDER,
		targetAttributes.y-CONFIG_BOX_BORDER,
		C.uint(boxWidth),
		C.uint(boxHeight),
		0,
		boxVisualInfo.depth,
		C.InputOutput,
		boxVisualInfo.visual,
		C.CWColormap|C.CWBorderPixel|C.CWBackPixel|C.CWOverrideRedirect,
		&boxAttributes,
	)

	// ウインドウの設定
	C.XReparentWindow(
		wm_display,
		client.window[CLIENT_APP],
		client.window[CLIENT_BOX],
		CONFIG_BOX_BORDER,
		CONFIG_TITLEBAR_HEIGHT+CONFIG_BOX_BORDER,
	)
	C.XMapRaised(wm_display, client.window[CLIENT_BOX])
	C.XSelectInput(wm_display, client.window[CLIENT_BOX], C.SubstructureNotifyMask)

	// ウインドウのcairoコンテキストを作成
	client.surface[CLIENT_BOX] = C.cairo_xlib_surface_create(
		wm_display,
		client.window[CLIENT_BOX],
		boxVisualInfo.visual,
		C.int(boxWidth),
		C.int(boxHeight),
	)
	client.cr[CLIENT_BOX] = C.cairo_create(client.surface[CLIENT_BOX])

	var hint C.XClassHint
	client.getClassHint(&hint)
	client.title = C.GoString(hint.res_class)

	C.cairo_set_operator(client.cr[CLIENT_BOX], C.CAIRO_OPERATOR_SOURCE)
	C.cairo_set_antialias(client.cr[CLIENT_BOX], C.CAIRO_ANTIALIAS_SUBPIXEL)
	C.cairo_set_line_width(client.cr[CLIENT_BOX], CONFIG_SHADOW_ROUGHNESS)

	// コンテキストを更新
	client.drawClient()
	(*clientTable)[client.window[CLIENT_BOX]] = &client
}
