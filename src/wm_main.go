package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

// #cgo LDFLAGS: -lX11 -lcairo
// #include <X11/Xlib.h>
// #include <X11/cursorfont.h>
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

const (
	CLIENT_WIDGETS = 2 // 構成ウインドウの数
	CLIENT_BOX     = 0 // Box
	CLIENT_APP     = 1 // Application
	CLIENT_MASK    = 2 // Mask

	RESIZE_ANGLE_TOP    = 0
	RESIZE_ANGLE_BOTTOM = 1
	RESIZE_ANGLE_START  = 2
	RESIZE_ANGLE_END    = 3
	EXIT_PRESSED        = 4

	CONFIG_BOX_BORDER = 12
	// クライアントの枠線の幅のうち、影を描画しない外側の部分。
	CONFIG_EMPTY_BOX_BORDER = 4
	// クライアントのTITLEBARの右端とAPPの右端のX座標の差 = EXITの幅
	CONFIG_TITLEBAR_WIDTH_MARGIN = 25
	// TITLEBARの高さ
	CONFIG_TITLEBAR_HEIGHT = 25
	// クライアントの影の粗さ。必ず0より大きい値に！
	// 1.5~2.5 ぐらいがいい感じ。粗すぎると変な見た目になるし、逆に
	CONFIG_SHADOW_ROUGHNESS = 1.5

	WIDTH_DIFF  = 2 * CONFIG_BOX_BORDER
	HEIGHT_DIFF = CONFIG_TITLEBAR_HEIGHT + 2*CONFIG_BOX_BORDER
)

var (
	display    *C.Display
	rootWindow C.Window
)

func main() {
	var config MiteWMConfig
	if len(os.Args) > 1 {
		log.Println(os.Args[1])
		data, err := ioutil.ReadFile(os.Args[1])
		if err != nil {
			log.Println("Failed to load Configuration file.")
		} else {
			err = json.Unmarshal(data, &config)
			if err != nil {
				log.Fatalln("Failed to parse Configuration file.")
			}
		}
	}
	display = C.XOpenDisplay(nil)
	defer C.XCloseDisplay(display)
	if display == nil {
		log.Fatalln("Failed to Open Display.")
	}

	clientTable := make(ClientTable)

	rootWindow = C.XDefaultRootWindow(display)

	C.XSelectInput(
		display, rootWindow,
		C.ButtonPressMask|C.ButtonReleaseMask|C.PointerMotionMask|C.SubstructureNotifyMask,
	)

	if config.BackgroundImageFile != "" {
		setBackground(config.BackgroundImageFile)
	}

	var info GlobalInfo
	info.window = C.None

	var lastUngrabbedApp C.Window = C.None

	C.XDefineCursor(
		display, rootWindow,
		C.XCreateFontCursor(display, C.XC_left_ptr),
	)
	for true {
		var event C.XEvent
		if C.XPending(display) < 0 || C.XNextEvent(display, &event) < 0 {
			break
		}
		handle_event(event, &clientTable, &lastUngrabbedApp, &info)
	}
}
