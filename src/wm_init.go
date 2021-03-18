package main 

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

// #include <X11/Xlib.h>
// #include <X11/cursorfont.h>
import "C"

func wm_setup(clientTable *ClientTable, config *MiteWMConfig, info *GlobalInfo){

	C.XSelectInput(
		wm_display, C.XDefaultRootWindow(wm_display),
		C.ButtonPressMask|C.ButtonReleaseMask|C.PointerMotionMask|C.SubstructureNotifyMask,
	)
	if config.BackgroundImageFile != "" {
		setBackground(config.BackgroundImageFile)
	}
	info.window = C.None
	C.XDefineCursor(
		wm_display, C.XDefaultRootWindow(wm_display),
		C.XCreateFontCursor(wm_display, C.XC_left_ptr),
	)
}

func wm_read_config_file(config *MiteWMConfig){
	// JSONファイルの読み込み
	if len(os.Args) > 1 {
		log.Println(os.Args[1])
		data, err := ioutil.ReadFile(os.Args[1])
		if err != nil {
			log.Println("Failed to load Configuration file.")
		} else {
			err = json.Unmarshal(data, config)
			if err != nil {
				log.Fatalln("Failed to parse Configuration file.")
			}
		}
	}
}