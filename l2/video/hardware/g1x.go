/*
Component: Luna G1X
Resolutions:
320x200

Description:
The Luna G1X is an improvement over the previous
Luna G1, being able to handle basic transparency
by passing color 0xE3
*/
package main

import (
	"image"
	"image/color"
	"luna_l2/font"
	"luna_l2/shared"
)

var CursorX int = 0
var CursorY int = 0
var MemoryVideo = make([]byte, 64000)
var Palette = make([]color.NRGBA, 256)
var img = image.NewRGBA(image.Rect(0, 0, 320, 200))

func scrollUp() {
    const (
        screenWidth  = 320
        screenHeight = 200
        charHeight   = 8
    )

    lineSize := screenWidth * charHeight
    visibleLines := screenHeight / charHeight
 
    copy(MemoryVideo[0:], MemoryVideo[lineSize:])
 
    bottomStart := (visibleLines - 1) * lineSize
    for i := bottomStart; i < len(MemoryVideo); i++ {
        MemoryVideo[i] = 0
    }
 
    CursorY = visibleLines - 1
}

func PushChar(x, y int, ch rune, fg byte, bg byte) {
    idx := int(ch)
    glyph := font.Font[0x00]

    if idx >= 0 && idx < len(font.Font) {
        glyph = font.Font[idx]
    }

    for row := 0; row < 8; row++ {
        line := glyph[row]
        
		for col := 0; col < 8; col++ {
			mask := byte(1 << col)
			var color byte
			if line & mask != 0 {
				color = fg
			} else {
				color = bg
			}
			px := (y + row) * 320 + (x+col)
			MemoryVideo[shared.Clamp(int(px), 0, 63999)] = color
		}

    }
}

func PrintChar(ch rune, fg byte, bg byte) {
	if ch == 0x0a {
		CursorY++
		CursorX = 0
	} else if ch == 0x0d {
		CursorX = 0
		return
	} else if ch == 0x00 {
		return
	}

	if CursorY >= 200 / 8 {
		scrollUp()	
	}

	if ch == 0x0a {
		return
	}

	x := CursorX * 8
	y := CursorY * 8		

	PushChar(x, y, ch, fg, bg)	

	CursorX++
	if CursorX >= 320/8 {
		CursorY++
		CursorX = 0
	}
}

func InitializePalette() {
	for i := 0; i < 256; i++ {
		r := (i >> 5) & 0x07
        g := (i >> 2) & 0x07 
        b := i & 0x03        
 
        R := uint8(r * 255 / 7)
        G := uint8(g * 255 / 7)
        B := uint8(b * 255 / 3)

        Palette[i] = color.NRGBA{R, G, B, 255}
    }
}

func SetCursor(x int, y int) {
	CursorX = x
	CursorY = y
}

func GetCursor() (int, int) {
	return CursorX, CursorY
}

func ClearVideoMemory() {
	for i := 0; i < len(MemoryVideo); i++ {
		MemoryVideo[i] = 0x00
	}
}

func WriteVideoMemory(addr uint32, content byte) {
	MemoryVideo[addr] = content
}

func ReadVideoMemory(addr uint32) byte {
	return MemoryVideo[addr]
}

func ReturnFramebuffer() *image.RGBA {
	i := 0
	for y := 0; y < 200; y++ {
		for x := 0; x < 320; x++ {
			if MemoryVideo[i] != 0xE3 {
				img.Set(x, y, Palette[MemoryVideo[i]])
			}
			i++
		}
	}

	return img
}
