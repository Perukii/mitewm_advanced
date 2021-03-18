package main

// #cgo LDFLAGS: -lX11 -lcairo
// #include <X11/Xlib.h>
import "C"


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
	wm_display *C.Display
	background Background
)