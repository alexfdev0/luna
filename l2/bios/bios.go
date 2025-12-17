package bios
import (
	"luna_l2/video"
	"luna_l2/types"
	"luna_l2/audio"
	"time"	
	"os"
	"fmt"
)

var TypeOut bool = false
var KeyTrap bool = false
var Registers *[]types.Register
var Memory *[0x70000000]byte
var KeyInterruptCode uint32 = 0x5
const (
	MEMSIZE uint32 = 0x70000000
	MEMCAP uint32 = 0x6FFFFFFF
)

// Interrupt modes:
	// 0: BIOS handle
	// 1: Software handle
	// 2: Dual handle (BIOS first)
	// 3: Dual handle (Software first)

func WriteChar(char string, fg uint8, bg uint8) {
	video.PrintChar(rune(char[0]), byte(fg), byte(bg))
}

func WriteString(str string, fg uint8, bg uint8) {
	for _, r := range str {
		WriteChar(string(r), fg, bg)
	}
}

func WriteLine(str string, fg uint8, bg uint8) {
	WriteString(str + "\n", fg, bg)
}

func setRegister(address uint32, value uint32) {
	for i := range (*Registers) {
		if (*Registers)[i].Address == address {
			if types.Bits32 == false {
				(*Registers)[i].Value = uint32(uint16(value))
			} else {
				(*Registers)[i].Value = value
			}
		}
	}
}

func getRegister(address uint32) uint32 {
	for _, register := range (*Registers) {
		if register.Address == address {
			return register.Value
		}
	}
	return 0x0000
}

func LoadSector(drive int, sector int, enforce bool) {	
	var file string
	switch drive {
	case 0:
		file = types.Filename
		time.Sleep(time.Duration(12) * time.Millisecond)
	case 1:
		file = types.SDFilename
		time.Sleep(time.Duration(2) * time.Millisecond)
	case 2:
		file = types.OpticalFilename
		time.Sleep(time.Duration(110) * time.Millisecond)
	}

	f, err := os.OpenFile(file, os.O_RDONLY, 0)
	if err != nil {
		if enforce == false {
			fmt.Println("luna-l2: could not load/reload block device")
			return
		} else {
			fmt.Println("luna-l2: could not open '" + file + "'", err)
			os.Exit(1)
		}
	}
	defer f.Close()

	start := sector * 512
	_, err = f.ReadAt((*Memory)[start:start + 512], int64(start))
	if err != nil {
		fmt.Println("luna-l2: could not read from disk: ", err)
	}	
}

func WriteSector(drive int, sector int) {
	var file string
	switch drive {
	case 0:
		file = types.Filename
		time.Sleep(time.Duration(15) * time.Millisecond)
	case 1:
		file = types.SDFilename
		time.Sleep(time.Duration(8) * time.Millisecond)
	case 2:
		file = types.OpticalFilename
		time.Sleep(time.Duration(200) * time.Millisecond)
	}
	
	f, err := os.OpenFile(file, os.O_RDWR | os.O_SYNC, 0)
	if err != nil {
		fmt.Println("luna-l2: could not load/reload block device")
		return
	}
	defer f.Close()

	start := sector * 512

	_, err = f.WriteAt((*Memory)[start:start + 512], int64(start))
	if err != nil {
		fmt.Println("luna-l2: could not write to block device")
	}
}

func IntHandler(code uint32) {
	if code == 0x01 {
		// BIOS print to screen
		// start address in R1
		// Foreground in R2
		// Background in R3
		char := getRegister(0x0001)
		WriteChar(string(rune(char)), uint8(getRegister(0x0002)), uint8(getRegister(0x0003)))
	} else if code == 0x02 {
		// BIOS sleep
		// seconds in R1
		timeToSleep := getRegister(0x0001)
		time.Sleep(time.Duration(timeToSleep) * time.Millisecond)
	} else if code == 0x03 {
		// BIOS write to VRAM
		// address in R1, word in R2
		address := getRegister(0x0001)
		word := getRegister(0x0002)
		if types.Bits32 == false {
			video.MemoryVideo[video.Clamp(address, 0, 63999)] = byte(uint16(word) >> 8)
			video.MemoryVideo[video.Clamp(address + 1, 0, 63999)] = byte(uint16(word) & 0xFF)
		} else {
			video.MemoryVideo[video.Clamp(address, 0, 63999)] = byte(uint32(word) >> 24)
			video.MemoryVideo[video.Clamp(address + 1, 0, 63999)] = byte(uint32(word) >> 16)
			video.MemoryVideo[video.Clamp(address + 2, 0, 63999)] = byte(uint32(word) >> 8)
			video.MemoryVideo[video.Clamp(address + 3, 0, 63999)] = byte(uint32(word) & 0xFF)
		}
	} else if code == 0x4 {
		// BIOS configure input mode
		// Mode 1: no type output
		// Mode 2: type output
		// In R1
		if getRegister(0x0001) == 1 {
			TypeOut = true
		} else {
			TypeOut = false
		}
	} else if code == 0x5 {
		// BIOS key event	
		if TypeOut == true {
			WriteChar(string(rune(getRegister(0x001b))), uint8(255), uint8(0))	
		}
		if KeyTrap == true {
			KeyTrap = false
			setRegister(0x0001, getRegister(0x001b))
		}
	} else if code == 0x6 {
		// BIOS wait for key
		// Return in R1 via interrupt 5
		KeyTrap = true
		for {
			if KeyTrap == true {
				time.Sleep(time.Duration(15) * time.Millisecond)
			} else {
				break
			}
		}
	} else if code == 0x7 {
		WriteLine("Illegal instruction 0x" + fmt.Sprintf("%08x", getRegister(0x0001)) + " at location 0x" + fmt.Sprintf("%08x", getRegister(0x001a)), 255, 0)
		return
	} else if code == 0x8 {
		// BIOS write to ARAM
		// address in R1, word in R2
		address := getRegister(0x0001)
		word := getRegister(0x0002)
		if types.Bits32 == false {
			audio.MemoryAudio[video.Clamp(address, 0, MEMCAP)] = byte(uint16(word) >> 8)
			audio.MemoryAudio[video.Clamp(address + 1, 0, MEMCAP)] = byte(uint16(word) & 0xFF)
		} else {
			audio.MemoryAudio[video.Clamp(address, 0, MEMCAP)] = byte(uint32(word) >> 24)
			audio.MemoryAudio[video.Clamp(address + 1, 0, MEMCAP)] = byte(uint32(word) >> 16)
			audio.MemoryAudio[video.Clamp(address + 2, 0, MEMCAP)] = byte(uint32(word) >> 8)
			audio.MemoryAudio[video.Clamp(address + 3, 0, MEMCAP)] = byte(uint32(word) & 0xFF)
		}	
	} else if code == 0x9 {
		audio.Play()	
	} else if code == 0xa {
		if types.Bits32 == false {
			setRegister(0x0001, 0xffff)
		} else {
			setRegister(0x0001, MEMSIZE)
		}
	} else if code == 0xb {
		sector := getRegister(0x0001)
		drive := getRegister(0x0002)
		LoadSector(int(drive), int(sector), false)
	} else if code == 0xc {
		video.CursorX = int(getRegister(0x0001))
		video.CursorY = int(getRegister(0x0002))
	} else if code == 0xd {	
		sector := getRegister(0x0001)
		drive := getRegister(0x0002)
		WriteSector(int(drive), int(sector))
	} else if code == 0xe {
		setRegister(0x0001, uint32(video.CursorX))
		setRegister(0x0002, uint32(video.CursorY))
	} else if code == 0xf {
		types.BootDrive = int(getRegister(0x0001))
		for i, _ := range (*Registers) {
			(*Registers)[i].Value = uint32(0)
		}
		types.Bits32 = false	
		*Memory = [0x70000000]byte {}
		video.MemoryVideo = [64000]byte {}
		video.CursorX = 0
		video.CursorY = 0
	} else if code == 0x10 {
		setRegister(0x0001, uint32(types.DriveNumber))
	}
}

func Splash() {
	WriteLine("Luna L2", 255, 0)
	WriteLine("BIOS: Integrated BIOS", 255, 0)	
	WriteLine("Copyright (c) 2025 Alexander Flax\n", 255, 0)
}

func CheckArgs() bool {
	if len(os.Args) < 2 {
		Splash()
		WriteLine("No bootable device", 255, 0)
		return false
	}
	return true
}

func IntWrapper(code uint32, next uint32) {
	var mem_location uint32 = 0x6FFF0000 + uint32(((code - 1) * 6))	
	if Memory[mem_location + 1] != 0 { 
		sp := getRegister(0x0019)
		if types.Bits32 == false {
			sp = video.Clamp(sp - 2, 0, MEMCAP)
			Memory[sp] = byte(next & 0xFF)
			Memory[sp + 1] = byte(next >> 8)
		} else {
			sp = video.Clamp(sp - 4, 0, MEMCAP)
			Memory[sp] = byte(next & 0xFF)
			Memory[sp + 1] = byte(next >> 8)
			Memory[sp + 2] = byte(next >> 16)
			Memory[sp + 3] = byte(next >> 24)
		}
		setRegister(0x0019, sp)
	}

	loc := uint32(Memory[mem_location + 2]) << 24 | uint32(Memory[mem_location + 3]) << 16 | uint32(Memory[mem_location + 4]) << 8 | uint32(Memory[mem_location + 5])	
	switch Memory[mem_location + 1] {
	case 0x00:
		IntHandler(code)
	case 1:	
		setRegister(0x001a, loc)
	case 2:
		IntHandler(code)
		setRegister(0x001a, loc)
	case 3:	
		setRegister(0x001a, loc)
		IntHandler(code)	
	}
}
