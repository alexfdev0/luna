//go:build linux || darwin

package video

import (
	"plugin"
	"os"
	"fmt"
	"image"
)

var VideoComponent *plugin.Plugin

func InitializeComponent(Path string) {
	var err error
	VideoComponent, err = plugin.Open(Path)
	if err != nil {
		fmt.Println("luna-l2: failed to initialize video component with path '" + Path + "':", err)
		os.Exit(1)
	}

	InitializePalette, err := VideoComponent.Lookup("InitializePalette")
	if err != nil {
		fmt.Println("luna-l2: failed to send signal to video component with name 'InitializePalette':", err)
		return
	}
	
	InitializePalette.(func())()
}

func UpdateFramebuffer(Image *image.RGBA) {
	UFB, err := VideoComponent.Lookup("UpdateFramebuffer")
	if err != nil {
		fmt.Println("luna-l2: failed to send signal to video component with name 'UpdateFramebuffer':", err)
		return
	}

	UFB.(func(*image.RGBA))(Image)
}

func PrintChar(ch rune, fg byte, bg byte) {
	PC, err := VideoComponent.Lookup("PrintChar")
	if err != nil {
		fmt.Println("luna-l2: failed to send signal to video component with name 'PrintChar':", err)
		return
	}

	PC.(func(rune, byte, byte))(ch, fg, bg)
}

func SetCursor(x int, y int) {
	SC, err := VideoComponent.Lookup("SetCursor")
	if err != nil {
		fmt.Println("luna-l2: failed to send signal to video component with name 'SetCursor':", err)
		return
	}

	SC.(func(int, int))(x, y)
}

func GetCursor() (int, int) {
	GC, err := VideoComponent.Lookup("GetCursor")
	if err != nil {
		fmt.Println("luna-l2: failed to send signal to video component with name 'GetCursor':", err)
		return 0, 0
	}

	x, y := GC.(func() (int, int))()
	return x, y
}

func ClearVideoMemory() {
	CVM, err := VideoComponent.Lookup("ClearVideoMemory")
	if err != nil {
		fmt.Println("luna-l2: failed to send signal to video component with name 'ClearVideoMemory':", err)
		return
	}

	CVM.(func())()
}

func WriteVideoMemory(addr uint32, content byte) {
	WVM, err := VideoComponent.Lookup("WriteVideoMemory")
	if err != nil {
		fmt.Println("luna-l2: failed to send signal to video component with name 'WriteVideoMemory':", err)
		return
	}

	WVM.(func(uint32, byte))(addr, content)
}

func ReadVideoMemory(addr uint32) byte {
	GC, err := VideoComponent.Lookup("ReadVideoMemory")
	if err != nil {
		fmt.Println("luna-l2: failed to send signal to video component with name 'ReadVideoMemory':", err)
		return 0
	}

	content := GC.(func(uint32) byte)(addr)
	return content
}
