package main

import (
	//"encoding/json"
	//"io/ioutil"
	"log"
	//"os"
)

// #cgo LDFLAGS: -lX11 -lcairo
// #include <X11/Xlib.h>
import "C"

func main() {

	wm_display = C.XOpenDisplay(nil)
	defer C.XCloseDisplay(wm_display)
	if wm_display == nil {
		log.Fatalln("Failed to Open Display.")
	}

	var clientTable ClientTable = make(ClientTable)
	var config MiteWMConfig
	var info GlobalInfo
	var lastUngrabbedApp C.Window = C.None

	wm_read_config_file(&config)
	wm_setup(&clientTable, &config, &info)

	for true {
		var event C.XEvent
		if C.XPending(wm_display) < 0 || C.XNextEvent(wm_display, &event) < 0 {
			break
		}
		handle_event(event, &clientTable, &lastUngrabbedApp, &info)
	}
}
