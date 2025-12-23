package bios
import (
	"luna_l2/video"
	"luna_l2/shared"
	"luna_l2/audio"
	"time"	
	"os"
	"fmt"
)

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

func LoadSector(drive int, sector int, enforce bool) {	
	var file string
	switch drive {
	case 0:
		file = shared.Filename
		time.Sleep(time.Duration(12) * time.Millisecond)
	case 1:
		file = shared.SDFilename
		time.Sleep(time.Duration(2) * time.Millisecond)
	case 2:
		file = shared.OpticalFilename
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
	_, err = f.ReadAt(shared.Memory[start:start + 512], int64(start))
	if err != nil {
		fmt.Println("luna-l2: could not read from disk: ", err)
	}	
}

func WriteSector(drive int, sector int) {
	var file string
	switch drive {
	case 0:
		file = shared.Filename
		time.Sleep(time.Duration(15) * time.Millisecond)
	case 1:
		file = shared.SDFilename
		time.Sleep(time.Duration(8) * time.Millisecond)
	case 2:
		file = shared.OpticalFilename
		time.Sleep(time.Duration(200) * time.Millisecond)
	}	
	
	f, err := os.OpenFile(file, os.O_RDWR | os.O_SYNC, 0)
	if err != nil {
		fmt.Println("luna-l2: could not load/reload block device")
		return
	}
	defer f.Close()

	start := sector * 512

	_, err = f.WriteAt(shared.Memory[start:start + 512], int64(start))
	if err != nil {
		fmt.Println("luna-l2: could not write to block device")
	}
}

func IntHandler(code uint32) {
	switch code {
	case 0x01:
		char := shared.GetRegister(0x0001)
		WriteChar(string(rune(char)), uint8(shared.GetRegister(0x0002)), uint8(shared.GetRegister(0x0003)))
	case 0x02:
		// Programmable interval timer reserved
	case 0x03:
		// Unmapped
	case 0x04:
		// Unmapped
	case 0x05:
		// Keyboard reserved
	case 0x06:
		// Unmapped
	case 0x07:
		WriteLine("Illegal instruction 0x" + fmt.Sprintf("%08x", shared.GetRegister(0x0001)) + " at location 0x" + fmt.Sprintf("%08x", shared.GetRegister(0x001a)), 255, 0)
	case 0x08:
		// Unmapped
	case 0x09:
		// Unmapped
	case 0x0A:
		if shared.Bits32 == false {
			shared.SetRegister(0x0001, 0xffff)
		} else {
			shared.SetRegister(0x0001, MEMSIZE)
		}	
	case 0x0B:
		sector := shared.GetRegister(0x0001)
		drive := shared.GetRegister(0x0002)
		LoadSector(int(drive), int(sector), false)
	case 0x0C:
		video.CursorX = int(shared.GetRegister(0x0001))
		video.CursorY = int(shared.GetRegister(0x0002))
	case 0x0D:
		sector := shared.GetRegister(0x0001)
		drive := shared.GetRegister(0x0002)
		WriteSector(int(drive), int(sector))
	case 0x0E:
		shared.SetRegister(0x0001, uint32(video.CursorX))
		shared.SetRegister(0x0002, uint32(video.CursorY))
	case 0x0F:
		shared.BootDrive = int(shared.GetRegister(0x0001))
		for i, _ := range (*shared.Registers) {
			(*shared.Registers)[i].Value = uint32(0)
		}
		shared.Bits32 = false	
		(*shared.Memory) = [0x70000000]byte {}
		video.MemoryVideo = [64000]byte {}
		video.CursorX = 0
		video.CursorY = 0
	case 0x10:
		shared.SetRegister(0x0001, uint32(shared.DriveNumber))
	case 0x11:
		// Shut down machine
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
	// shared.SetRegister(0x001e, code)
	var mem_location uint32 = 0x6FFF0000 + uint32(((code - 1) * 6))	
	loc := uint32(shared.Memory[mem_location + 2]) << 24 | uint32(shared.Memory[mem_location + 3]) << 16 | uint32(shared.Memory[mem_location + 4]) << 8 | uint32(shared.Memory[mem_location + 5])
	
	switch shared.Memory[mem_location + 1] {
	case 0x00:
		IntHandler(code)
	case 0x01:
		shared.SetRegister(0x001c, next)
		shared.SetRegister(0x001a, loc)
	case 0x02:
		IntHandler(code)
		shared.SetRegister(0x001c, next)
		shared.SetRegister(0x001a, loc)
	case 0x03:
		shared.SetRegister(0x001c, next)
		shared.SetRegister(0x001a, loc)
		IntHandler(code)
	}
}
