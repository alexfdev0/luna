package main

import (	
	"image"
	"os"	
	"time"
	"fmt"
	"strconv"
	"bufio"
	"runtime"	

	"luna_l2/bios"		
	"luna_l2/video"
	"luna_l2/shared"
	"luna_l2/audio"
	"luna_l2/network"
	"luna_l2/rtc"
	"luna_l2/keyboard"
	"luna_l2/pit"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"	
)

// Basic elements of CPU
var Registers = []shared.Register {
	{0x0000, "R0", 0},
	{0x0001, "R1", 0},
	{0x0002, "R2", 0},
	{0x0003, "R3", 0},
	{0x0004, "R4", 0},
	{0x0005, "R5", 0},
	{0x0006, "R6", 0},
	{0x0007, "R7", 0},
	{0x0008, "R8", 0},
	{0x0009, "R9", 0},
	{0x000a, "R10", 0},
	{0x000b, "R11", 0},
	{0x000c, "R12", 0},
	{0x000d, "E0", 0},
	{0x000e, "E1", 0},
	{0x000f, "E2", 0},
	{0x0010, "E3", 0},
	{0x0011, "E4", 0},
	{0x0012, "E5", 0},
	{0x0013, "E6", 0},
	{0x0014, "E7", 0},
	{0x0015, "E8", 0},
	{0x0016, "E9", 0},
	{0x0017, "E10", 0},
	{0x0018, "E11", 0},
	{0x001b, "E12", 0},
	{0x0019, "SP", 0},
	{0x001a, "PC", 0},
	{0x001c, "IRV", 0},
	{0x001e, "IR", 0},
	{0x001f, "B", 0},
}

var Memory [0x70000000]byte
const (
	MEMSIZE uint32 = 0x70000000
	MEMCAP uint32 = 0x6FFFFFFF
)

// Register controls
func setRegister(address uint32, value uint32) {
	for i := range Registers {
		if Registers[i].Address == address {
			if shared.Bits32 == false {
				Registers[i].Value = uint32(uint16(value))
			} else {
				Registers[i].Value = value
			}
		}
	}
}

func getRegister(address uint32) uint32 {
	for _, register := range Registers {
		if register.Address == address {
			return register.Value
		}
	}
	return 0x0000
}

func getRegisterName[T uint32 | byte](address T) string {
	addr := uint32(address)
	for _, register := range Registers {
		if register.Address == addr {
			return register.Name
		}
	}
	return ""
}

// Memory map (32 bit mode):
	// 0x00000000 - 0x6FFEFFFF: general purpose RAM (1.7499 GiB)
	// 0x6FFF0000 - 0x6FFFFFFF: reserved region for IDT (65 KB)
	// 0x70000000 - 0x7000F9FF: video RAM (64 KB)
	// 0x7000FA00 - 0x7000FA09: audio RAM (9 B)
	// 0x7000FA0A - 0x7001A643: empty region
	// 0x7001A644 - 0x7001B65D: network RAM (4.1 KB)
	// 0x7001B65E - 0x7001B663: clock RAM (6 B)

// Memory map (16 bit mode):
	// 0x0000 - 0xEFFF - general purpose RAM (61.4KB)
	// 0xF000 - 0xFFFF - MMIO device (accessible via bank)
// Bank map:
	// 0 - 15: VRAM 
	// 16: Audio RAM + Mouse RAM + Keyboard RAM + PIT RAM
	// 17: Network RAM + RTC RAM
// Size formula
	// end = start + size - 1

// Meta-code
var LogOn bool = false
var Debug bool = false
var ClockSpeed int64 = 33000000
var BIOS_REBOOT bool = false
var BIOS_SHUTDOWN bool = false
func Log(text string) {
	if LogOn == true {
		fmt.Println("\033[33m" + fmt.Sprintf("0x%08x: ", getRegister(0x001a)) + text + "\033[0m")	
	}	
}

// CPU code
var accumulated int64
func stall(cycles int64) { 
	cycleTime := int64(int(time.Second)) / ClockSpeed
	accumulated += cycles
	if accumulated >= 66000 {
		time.Sleep(time.Duration(cycleTime * accumulated))
		accumulated = 0
	}	
}

func execute() {
	for {
		ProgramCounter := getRegister(0x001a)
		op := shared.Mapper(ProgramCounter)

		// Handle interrupts	
		var IntHandled bool
		for i := 0; i < 32; i++ {
			if (getRegister(0x001e) & (1 << i)) != 0 {
				code := i + 1
				IntHandled = true
				setRegister(0x001e, getRegister(0x001e) &^ (1 << i))
				bios.IntWrapper(uint32(code), ProgramCounter)
				
				if code == 0x0f {
					Log("system reboot")
					BIOS_REBOOT = true
					return
				} else if code == 0x11 {
					Log("system shutdown")
					BIOS_SHUTDOWN = true
					return
				}
				break
			}
		}
		if IntHandled == true {
			continue
		}

		switch op {
		case 0x00:
			Log("null")
			now := ProgramCounter
			bios.IntWrapper(0x7, ProgramCounter + 1)
			if getRegister(0x001a) == now {
				return
			}
		case 0x01:
			// MOV
			mode := shared.Mapper(ProgramCounter + 1)
			dst := shared.Mapper(ProgramCounter + 2)

			if mode == 0x01 {
				var imm uint32 = 0
				var next uint32 = 0
				if shared.Bits32 == false {
					imm = uint32(uint16(Memory[ProgramCounter + 3]) << 8 | uint16(Memory[ProgramCounter + 4]))
					next = ProgramCounter + 5
				} else {
					imm = uint32(Memory[ProgramCounter + 3]) << 24 | uint32(Memory[ProgramCounter + 4])	<< 16 | uint32(Memory[ProgramCounter + 5]) << 8 | uint32(Memory[ProgramCounter + 6])
					next = ProgramCounter + 7
				}
				setRegister(uint32(dst), imm)
				setRegister(0x001a, next)
				Log("mov " + getRegisterName(uint32(dst)) + ", " + fmt.Sprintf("0x%08x", imm))
			} else if mode == 0x02 {
				frm := uint32(Memory[ProgramCounter+3])
				setRegister(uint32(dst), uint32(getRegister(frm)))
				setRegister(0x001a, ProgramCounter+4)
				Log("mov " + getRegisterName(uint32(dst)) + ", " + getRegisterName(frm))
			}	
			stall(4)
		case 0x02:
			// HLT
			Log("hlt")
			now := ProgramCounter
			for {
				if getRegister(0x001a) != now || getRegister(0x001e) != 0 {
					break	
				}	
				time.Sleep(time.Duration(15) * time.Millisecond)
			}
			setRegister(0x001a, ProgramCounter + 1)
		case 0x03:
			// JMP	
			mode := Memory[ProgramCounter + 1]

			if mode == 0x01 {
				var loc uint32 = 0
				if shared.Bits32 == false {
					loc = uint32(uint16(Memory[ProgramCounter + 2]) << 8 | uint16(Memory[ProgramCounter + 3]))
				} else {
					loc = uint32(Memory[ProgramCounter + 2]) << 24 | uint32(Memory[ProgramCounter + 3])	<< 16 | uint32(Memory[ProgramCounter + 4]) << 8 | uint32(Memory[ProgramCounter + 5])
				}
				Log("jmp " + fmt.Sprintf("0x%08x", loc))
				setRegister(0x001a, loc)	
			} else if mode == 0x02 {
				frm := uint32(Memory[ProgramCounter+2])
				loc := getRegister(frm)
				Log("jmp " + getRegisterName(frm))
				setRegister(0x001a, loc)	
			}
			stall(8)
		case 0x04:
			// INT
			var code uint32 = 0
			var next uint32 = 0
	
			if shared.Bits32 == false {
				code = uint32(uint16(Memory[ProgramCounter + 1]) << 8 | uint16(Memory[ProgramCounter + 2]))
				next = ProgramCounter + 3
			} else {
				code = uint32(Memory[ProgramCounter + 1]) << 24 | uint32(Memory[ProgramCounter + 2])	<< 16 | uint32(Memory[ProgramCounter + 3]) << 8 | uint32(Memory[ProgramCounter + 4])
				next = ProgramCounter + 5
			}	
	
			shared.RaiseInterrupt(code)
			setRegister(0x001a, next)	

			Log("int " + fmt.Sprintf("0x%08x", code))	
			stall(34)
		case 0x05:
			// JNZ
			// jnz <mode (01 or 02)> <check register> <loc (register or raw addr)>
			mode := Memory[ProgramCounter+1]
			checkRegister := Memory[ProgramCounter+2]
			var loc uint32 = 0
			var not uint32 = 0

			if mode == 0x01 {	
				if shared.Bits32 == false {
					loc = uint32(uint16(Memory[ProgramCounter + 3]) << 8 | uint16(Memory[ProgramCounter + 4]))
					not = ProgramCounter + 5
				} else {
					loc = uint32(Memory[ProgramCounter + 3]) << 24 | uint32(Memory[ProgramCounter + 4])	<< 16 | uint32(Memory[ProgramCounter + 5]) << 8 | uint32(Memory[ProgramCounter + 6])
					not = ProgramCounter + 7
				}	
				Log("jnz " + getRegisterName(uint32(checkRegister)) + ", " + fmt.Sprintf("0x%08x", loc))
			} else if mode == 0x02 {
				frm := uint32(Memory[ProgramCounter+3])
				loc = getRegister(frm)
				not = ProgramCounter + 4
				Log("jnz " + getRegisterName(uint32(checkRegister)) + ", " + getRegisterName(frm))
			}

			if getRegister(uint32(checkRegister)) != 0 {
				setRegister(0x001a, loc)
			} else {
				setRegister(0x001a, not)
			}
			stall(8)
		case 0x06:
			// NOP
			setRegister(0x001a, ProgramCounter+1)
			Log("nop")
			stall(1)
		case 0x07:
			// CMP
			// Syntax: CMP <to> <r1> <r2>
			to := Memory[ProgramCounter+1]
			first := Memory[ProgramCounter+2]
			second := Memory[ProgramCounter+3]
			Log("cmp " + getRegisterName(uint32(to)) + ", " + getRegisterName(first) + ", " + getRegisterName(second))

			if getRegister(uint32(first)) == getRegister(uint32(second)) {
				setRegister(uint32(to), uint32(1))
			} else {
				setRegister(uint32(to), uint32(0))
			}
			setRegister(0x001a, ProgramCounter+4)
			stall(4)
		case 0x08:
			// JZ
			// jz <mode (01 or 02)> <check register> <loc (register or raw addr)>
			mode := Memory[ProgramCounter+1]
			checkRegister := Memory[ProgramCounter+2]
			var loc uint32 = 0
			var not uint32 = 0

			if mode == 0x01 {	
				if shared.Bits32 == false {
					loc = uint32(uint16(shared.Mapper(ProgramCounter + 3)) << 8 | uint16(shared.Mapper(ProgramCounter + 4)))
					not = ProgramCounter + 5
				} else {
					loc = uint32(shared.Mapper(ProgramCounter + 3)) << 24 | uint32(shared.Mapper(ProgramCounter + 4)) << 16 | uint32(shared.Mapper(ProgramCounter + 5)) << 8 | uint32(shared.Mapper(ProgramCounter + 6))
					not = ProgramCounter + 7
				}	
				Log("jz " + getRegisterName(checkRegister) + ", " + fmt.Sprintf("0x%08x", loc))
			} else if mode == 0x02 {
				frm := uint32(Memory[ProgramCounter+3])
				loc = getRegister(frm)
				not = ProgramCounter + 4
				Log("jz " + getRegisterName(checkRegister) + ", " + getRegisterName(frm))
			}

			if getRegister(uint32(checkRegister)) == 0 {
				setRegister(0x001a, loc)
			} else {
				setRegister(0x001a, not)
			}
			stall(8)
		case 0x09:
			// INC
			// inc <register>
			register := uint32(Memory[ProgramCounter+1])
			setRegister(register, getRegister(register)+1)
			setRegister(0x001a, ProgramCounter+2)
			Log("inc " + getRegisterName(register))
			stall(1)
		case 0x0a:
			// DEC
			// dec <register>
			register := uint32(Memory[ProgramCounter+1])
			setRegister(register, getRegister(register)-1)
			setRegister(0x001a, ProgramCounter+2)
			Log("dec " + getRegisterName(register))
			stall(1)
		case 0x0b:
			// PUSH
			// push <mode> <immediate or register>
			mode := Memory[ProgramCounter + 1]
			var value uint32	
			if mode == 0x1 {	
				var next uint32 = 0
				if shared.Bits32 == false {
					value = uint32(uint16(shared.Mapper(ProgramCounter + 2)) << 8 | uint16(shared.Mapper(ProgramCounter + 3)))
					next = ProgramCounter + 4
				} else {
					value = uint32(shared.Mapper(ProgramCounter + 2)) << 24 | uint32(shared.Mapper(ProgramCounter + 3)) << 16 | uint32(shared.Mapper(ProgramCounter + 4)) << 8 | uint32(shared.Mapper(ProgramCounter + 5))
					next = ProgramCounter + 6
				}	
				setRegister(0x001a, next)
				Log("push " + fmt.Sprintf("0x%08x", value))
			} else if mode == 0x2 {
				value = getRegister(uint32(shared.Mapper(ProgramCounter + 2)))
				setRegister(0x001a, ProgramCounter + 3)
				Log("push " + getRegisterName(uint32(shared.Mapper(ProgramCounter + 2))))
			}	
			sp := getRegister(0x0019)
			if shared.Bits32 == false {
				sp = video.Clamp(sp - 2, 0, MEMCAP)
				shared.MapperWrite(sp, byte(value & 0xFF))
				shared.MapperWrite(sp + 1, byte(value >> 8))
			} else {
				sp = video.Clamp(sp - 4, 0, MEMCAP)
				shared.MapperWrite(sp, byte(value & 0xFF))
				shared.MapperWrite(sp + 1, byte(value >> 8))
				shared.MapperWrite(sp + 2, byte(value >> 16))
				shared.MapperWrite(sp + 3, byte(value >> 24))
			}	
			setRegister(0x0019, uint32(sp))	
			stall(2)
		case 0x0c:
			// POP
			// pop <register>	
			register := shared.Mapper(ProgramCounter + 1)
			sp := getRegister(0x0019)
			var value uint32
			if shared.Bits32 == false {
				value = uint32(uint16(shared.Mapper(sp)) | uint16(shared.Mapper(sp + 1)) << 8) 
			} else {	
				value = uint32(shared.Mapper(sp)) | uint32(shared.Mapper(sp + 1)) << 8 | uint32(shared.Mapper(sp + 2)) << 16 | uint32(shared.Mapper(sp + 3)) << 24
			}
			Log("value: " + fmt.Sprintf("0x%08x", value))
			setRegister(uint32(register), uint32(value))
			if shared.Bits32 == false {
				sp = video.Clamp(sp + 2, 0, MEMCAP)
			} else {
				sp = video.Clamp(sp + 4, 0, MEMCAP)
			}
			setRegister(0x0019, uint32(sp))
			setRegister(0x001a, ProgramCounter + 2)
			Log("pop " + getRegisterName(register))
			stall(2)
		case 0x0d:
			// ADD
			// add <register> <register> <register>
			toregister := Memory[ProgramCounter+1]
			regone := Memory[ProgramCounter+2]
			regtwo := Memory[ProgramCounter+3]
			setRegister(uint32(toregister), getRegister(uint32(regone))+getRegister(uint32(regtwo)))
			setRegister(0x001a, ProgramCounter+4)
			Log("add " + getRegisterName(toregister) + ", " + getRegisterName(regone) + ", " + getRegisterName(regtwo))
			stall(7)
		case 0x0e:
			// SUB
			// SUB <register> <register> <register>
			toregister := Memory[ProgramCounter+1]
			regone := Memory[ProgramCounter+2]
			regtwo := Memory[ProgramCounter+3]
			setRegister(uint32(toregister), getRegister(uint32(regone))-getRegister(uint32(regtwo)))
			setRegister(0x001a, ProgramCounter+4)
			Log("sub " + getRegisterName(toregister) + ", " + getRegisterName(regone) + ", " + getRegisterName(regtwo))
			stall(7)
		case 0x0f:
			// MUL
			// mul <register> <register> <register>
			toregister := Memory[ProgramCounter+1]
			regone := Memory[ProgramCounter+2]
			regtwo := Memory[ProgramCounter+3]
			setRegister(uint32(toregister), getRegister(uint32(regone))*getRegister(uint32(regtwo)))
			setRegister(0x001a, ProgramCounter+4)
			Log("mul " + getRegisterName(toregister) + ", " + getRegisterName(regone) + ", " + getRegisterName(regtwo))
			stall(70)
		case 0x10:
			// DIV
			// div <register> <register> <register>
			toregister := Memory[ProgramCounter+1]
			regone := Memory[ProgramCounter+2]
			regtwo := Memory[ProgramCounter+3]
			setRegister(uint32(toregister), getRegister(uint32(regone))/getRegister(uint32(regtwo)))
			setRegister(0x001a, ProgramCounter+4)
			Log("div " + getRegisterName(toregister) + ", " + getRegisterName(regone) + ", " + getRegisterName(regtwo))
			stall(140)
		case 0x11:
			// IGT
			// igt <register> <register> <register>
			toregister := Memory[ProgramCounter+1]
			regone := Memory[ProgramCounter+2]
			regtwo := Memory[ProgramCounter+3]
			if getRegister(uint32(regone)) > getRegister(uint32(regtwo)) {
				setRegister(uint32(toregister), uint32(1))
			} else {
				setRegister(uint32(toregister), uint32(0))
			}
			setRegister(0x001a, ProgramCounter + 4)
			Log("igt " + getRegisterName(toregister) + ", " + getRegisterName(regone) + ", " + getRegisterName(regtwo))
			stall(4)
		case 0x12:
			// ILT
			// ilt <register> <register> <register>
			toregister := Memory[ProgramCounter+1]
			regone := Memory[ProgramCounter+2]
			regtwo := Memory[ProgramCounter+3]
			if getRegister(uint32(regone)) < getRegister(uint32(regtwo)) {
				setRegister(uint32(toregister), uint32(1))
			} else {
				setRegister(uint32(toregister), uint32(0))
			}
			setRegister(0x001a, ProgramCounter + 4)
			Log("ilt " + getRegisterName(toregister) + ", " + getRegisterName(regone) + ", " + getRegisterName(regtwo))
			stall(4)
		case 0x13:
			// AND
			// and <register> <register> <register>
			toregister := Memory[ProgramCounter+1]
			regone := Memory[ProgramCounter+2]
			regtwo := Memory[ProgramCounter+3]
			setRegister(uint32(toregister), getRegister(uint32(regone)) & getRegister(uint32(regtwo)))	
			setRegister(0x001a, ProgramCounter + 4)
			Log("and " + getRegisterName(toregister) + ", " + getRegisterName(regone) + ", " + getRegisterName(regtwo))
			stall(1)
		case 0x14:
			// OR
			// or <register> <register> <register>
			toregister := Memory[ProgramCounter+1]
			regone := Memory[ProgramCounter+2]
			regtwo := Memory[ProgramCounter+3]
			setRegister(uint32(toregister), getRegister(uint32(regone)) | getRegister(uint32(regtwo)))	
			setRegister(0x001a, ProgramCounter + 4)
			Log("or " + getRegisterName(toregister) + ", " + getRegisterName(regone) + ", " + getRegisterName(regtwo))
			stall(1)
		case 0x15:
			// NOT
			// not <register> <register>
			toregister := Memory[ProgramCounter+1]
			regone := Memory[ProgramCounter+2]
			setRegister(uint32(uint32(toregister)), ^getRegister(uint32(regone)))	
			setRegister(0x001a, ProgramCounter + 3)
			Log("not " + getRegisterName(toregister) + ", " + getRegisterName(regone))
			stall(1)
		case 0x16:
			// XOR
			// xor <register> <register> <register>
			toregister := Memory[ProgramCounter+1]
			regone := Memory[ProgramCounter+2]
			regtwo := Memory[ProgramCounter+3]
			setRegister(uint32(toregister), getRegister(uint32(regone)) ^ getRegister(uint32(regtwo)))	
			setRegister(0x001a, ProgramCounter + 4)
			Log("xor " + getRegisterName(toregister) + ", " + getRegisterName(regone) + ", " + getRegisterName(regtwo))
			stall(6)
		case 0x17:
			// LOD
			// lod <addr (register)> <destination register>	
			addr := getRegister(uint32(Memory[ProgramCounter+1]))
			toregister := uint32(Memory[ProgramCounter+2])
			setRegister(toregister, uint32(shared.Mapper(addr)))
			setRegister(0x001a, ProgramCounter + 3)
			Log("lod " + getRegisterName(uint32(Memory[ProgramCounter + 1])) + ", " + getRegisterName(toregister) + " (" + fmt.Sprintf("0x%02x", shared.Mapper(addr)) + ")")
			stall(100)
		case 0x18:
			// STRF
			// strf <addr (register)> <value (register)>	
			addr := getRegister(uint32(Memory[ProgramCounter+1]))
			value := uint32(Memory[ProgramCounter+2])
			if shared.Bits32 == false {
				shared.MapperWrite(addr, byte(getRegister(value) >> 8))
				shared.MapperWrite(addr + 1, byte(getRegister(value) & 0xFF))
			} else {
				shared.MapperWrite(addr, byte(getRegister(value) >> 24))
				shared.MapperWrite(addr + 1, byte(getRegister(value) >> 16))
				shared.MapperWrite(addr + 2, byte(getRegister(value) >> 8))
				shared.MapperWrite(addr + 3, byte(getRegister(value) & 0xFF))
			}	
			setRegister(0x001a, ProgramCounter + 3)
			Log("str " + getRegisterName(uint32(Memory[ProgramCounter + 1])) + ", " + getRegisterName(value))
			stall(100)
		case 0x19:
			// LODF
			// lodf <addr (register)> <destination register>
			addr := getRegister(uint32(Memory[ProgramCounter+1]))
			toregister := uint32(Memory[ProgramCounter+2])
			if shared.Bits32 == false {
				setRegister(toregister, uint32(uint16(shared.Mapper(addr)) << 8 | uint16(shared.Mapper(addr + 1))))
			} else {
				setRegister(toregister, uint32(shared.Mapper(addr)) << 24 | uint32(shared.Mapper(addr + 1)) << 16 | uint32(shared.Mapper(addr + 2)) << 8 | uint32(shared.Mapper(addr + 3)))
			}
			setRegister(0x001a, ProgramCounter + 3)
			Log("lodf " + getRegisterName(uint32(Memory[ProgramCounter + 1])) + ", " + getRegisterName(toregister))
			stall(100)
		case 0x1a:
			// SET
			// set <00 or 01>
			mode := uint32(Memory[ProgramCounter + 1])
			if mode == 0 {
				shared.Bits32 = false
				Log("16 bit mode")
			} else if mode == 1 {
				shared.Bits32 = true
				Log("32 bit mode")
			}
			setRegister(0x001a, ProgramCounter + 2)
			stall(1)
		case 0x1b:
			// STR
			// str <addr> <register>
			addr := getRegister(uint32(Memory[ProgramCounter + 1]))
			value := uint32(Memory[ProgramCounter + 2])
			shared.MapperWrite(addr, byte(getRegister(value)))
			setRegister(0x001a, ProgramCounter + 3)
			stall(100)
		case 0x1c:
			// SHL
			// shl <dest> <value> <by>
			dest := uint32(Memory[ProgramCounter + 1])
			value := getRegister(uint32(Memory[ProgramCounter + 2]))
			by := getRegister(uint32(Memory[ProgramCounter + 3]))
			setRegister(dest, uint32(value) << uint32(by))
			stall(95)
		case 0x1d:
			// SHR
			// shr <dest> <value> <by>
			dest := uint32(Memory[ProgramCounter + 1])
			value := getRegister(uint32(Memory[ProgramCounter + 2]))
			by := getRegister(uint32(Memory[ProgramCounter + 3]))
			setRegister(dest, uint32(value) >> uint32(by))
			stall(95)
		default:
			setRegister(0x0001, uint32(op))
			Log("\033[31mIllegal instruction 0x" + fmt.Sprintf("%08x", uint32(op)) + "\033[33m")
			now := ProgramCounter
			bios.IntWrapper(0x7, ProgramCounter + 1)	
			if Debug == true {
				setRegister(0x001a, ProgramCounter + 1)
			} else {
				if getRegister(0x001a) == now {
					return
				}
			}
		}

		if Debug == true {
			bufio.NewReader(os.Stdin).ReadBytes('\n')
		}
	}
}

// Frontend code
var Ready bool
var img = image.NewRGBA(image.Rect(0, 0, 320, 200))
var Vertices = []float32 {
	-1, -1, 0, 1,
     1, -1, 1, 1,
     1,  1, 1, 0,

    -1, -1, 0, 1,
     1,  1, 1, 0,
    -1,  1, 0, 0,	
}

func UpdateFramebuffer() {
	i := 0
	for y := 0; y < 200; y++ {
		for x := 0; x < 320; x++ {
			img.Set(x, y, video.Palette[video.MemoryVideo[i]])
			i++
		}
	}
}

func ToggleGrab(window *glfw.Window, Grab bool) {
	if Grab == true {
		window.SetTitle("Luna L2 - Press Ctrl+Alt+G to release grab")
		window.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)
	} else {
		window.SetTitle("Luna L2")
		window.SetInputMode(glfw.CursorMode, glfw.CursorNormal)
	}
}

var FS bool
func ToggleFullscreen(window *glfw.Window) {
	if FS == false {
		window.SetFramebufferSizeCallback(func(w *glfw.Window, width, height int) {
   			gl.Viewport(0, 0, int32(width), int32(height))
		})
		window.SetMonitor(glfw.GetPrimaryMonitor(), 0, 0, 640, 400, 60)
		FS = true
	} else {
		window.SetMonitor(nil, 960, 540, 640, 400, 0)
		FS = false
	}
}

var Grab bool
func InitializeWindow() {
	wd, _ := os.Getwd()
	video.InitializePalette()
	err := glfw.Init()
	if err != nil {
		fmt.Println("luna-l2: could not initialize window: ", err)
		os.Exit(1)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 2)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(640, 400, "Luna L2", nil, nil)
	if err != nil {
		fmt.Println("luna-l2: could not initialize window: ", err)
		os.Exit(1)
	}
	window.MakeContextCurrent()

	err = gl.Init();
	if err != nil {
		fmt.Println("luna-l2: could not initialize window: ", err)
		os.Exit(1)
	}	

	gl.Viewport(0, 0, 640, 400)
	gl.ClearColor(0, 0, 0, 1)

	program := video.CreateProgram()
	gl.UseProgram(program)

	loc := gl.GetUniformLocation(program, gl.Str("tex\x00"))	
	gl.Uniform1i(loc, 0)

	var vao, vbo uint32

	gl.GenVertexArrays(1, &vao)
	gl.GenBuffers(1, &vbo)

	gl.BindVertexArray(vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)

	gl.BufferData(
		gl.ARRAY_BUFFER,
		len(Vertices) * 4,
		gl.Ptr(Vertices),
		gl.STATIC_DRAW,
	)

	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 2, gl.FLOAT, false, 4*4, gl.PtrOffset(0))

	gl.EnableVertexAttribArray(1)
	gl.VertexAttribPointer(1, 2, gl.FLOAT, false, 4*4, gl.PtrOffset(2*4))

	window.SetKeyCallback(func(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		if action == glfw.Press || action == glfw.Repeat {
			shift := (mods & glfw.ModShift) != 0
			alt := (mods & glfw.ModAlt) != 0
			ctrl := (mods & glfw.ModControl) != 0

			if ctrl && alt && key == glfw.KeyG {
				if Grab == true {
					ToggleGrab(window, false)
					Grab = false
					return
				}	
			}
			if ctrl && alt && key == glfw.KeyF {
				ToggleFullscreen(window)
				return
			}

			var char string
			switch key {
			case glfw.KeySpace:
				char = string(byte(0x20))
			case glfw.KeyEnter:
				char = string(byte(0x0A))
			case glfw.KeyBackspace:
				char = string(byte(0xC3))	
			default:
				char = glfw.GetKeyName(key, scancode)	
			}

			if shift {
				char = keyboard.Upper(char)
			} else {
				char = keyboard.Lower(char)
			}

			if len(char) > 0 {
				keyboard.MemoryKeyboard[0] = byte(char[0])
				shared.RaiseInterrupt(bios.KeyInterruptCode)
				setRegister(0x001b, uint32(char[0]))
				bios.IntHandler(bios.KeyInterruptCode)
			}
		}
	})

	window.SetMouseButtonCallback(func(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
		if button == glfw.MouseButtonLeft && action == glfw.Press {
			if Grab == false {
				ToggleGrab(window, true)
				Grab = true
				return
			}
		}
	})

	window.SetCursorPosCallback(func(w *glfw.Window, xpos float64, ypos float64) {
		if Grab == false {
			return
		}

		if xpos > 320 {
			xpos = 320
		} else if xpos < 0 {
			xpos = 0
		}

		if ypos > 320 {
			ypos = 320
		} else if ypos < 0 {
			ypos = 0
		}
	
		ixh := int(xpos) >> 8
		ixl := int(xpos) & 0xFF

		iyh := int(ypos) >> 8
		iyl := int(ypos) & 0xFF

		keyboard.MemoryMouse[2] = byte(ixh)
		keyboard.MemoryMouse[3] = byte(ixl)
		keyboard.MemoryMouse[6] = byte(iyh)
		keyboard.MemoryMouse[7] = byte(iyl)
		
		shared.RaiseInterrupt(0x12)
	})

	var texture uint32
	gl.GenTextures(1, &texture)
	gl.BindTexture(gl.TEXTURE_2D, texture)

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)

	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA8,
		320, 200,
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		nil,
	)	

	os.Chdir(wd)
	next := time.Now()	
	for !window.ShouldClose() {
		Ready = true
    	UpdateFramebuffer()

		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D, texture)
		gl.TexSubImage2D(
			gl.TEXTURE_2D,
			0,
			0, 0,
			320, 200,
			gl.RGBA,
			gl.UNSIGNED_BYTE,
			gl.Ptr(img.Pix),
		)

		gl.Clear(gl.COLOR_BUFFER_BIT)

		gl.UseProgram(program)
		gl.BindVertexArray(vao)

		gl.DrawArrays(gl.TRIANGLES, 0, 6)

		next = next.Add(time.Second / 70)
		sleep := time.Until(next)
		if sleep > 0 {
			time.Sleep(sleep)
		} else {
			next = time.Now()
		}

		window.SwapBuffers()
		glfw.PollEvents()
	}
}

var RequireDevicePresent bool = true
func main() {
	runtime.LockOSThread()

	shared.Registers = &Registers
	shared.Memory = &Memory
	shared.MemoryVideo = &video.MemoryVideo
	shared.MemoryAudio = &audio.MemoryAudio
	shared.MemoryMouse = &keyboard.MemoryMouse
	shared.MemoryKeyboard = &keyboard.MemoryKeyboard
	shared.MemoryNetwork = &network.MemoryNetwork
	shared.MemoryRTC = &rtc.MemoryRTC
	shared.MemoryPIT = &pit.MemoryPIT

	go func() {
		if Ready == false {	
			for {
				if Ready == true {
					break
				} else {
					time.Sleep(500)
				}
			}
		}	

		if bios.CheckArgs() == false {
			return
		}	
	
		for i := 1; i < len(os.Args); i++ {
			arg := os.Args[i]
			switch arg {
			case "--speed":
				if i + 1 >= len(os.Args) { fmt.Println("Not enough arguments to --speed"); i++; continue }
				speed, err := strconv.ParseInt(os.Args[i + 1], 0, 64)
				if err != nil {
					fmt.Println("Invalid clock speed")
					i++
					continue
				}
				ClockSpeed = int64(speed)
				i++
			case "--log":
				LogOn = true
			case "--debug":
				Debug = true
				LogOn = true
			case "-sd":
				shared.SDFilename = os.Args[i + 1]
				i++
			case "-dvd":
				shared.OpticalFilename = os.Args[i + 1]
				i++	
			case "-boot":
				next := os.Args[i + 1]
				switch next {
				case "hdd":
					shared.BootDrive = 0
				case "sd":
					shared.BootDrive = 1
				case "dvd":
					shared.BootDrive = 2
				default:
					fmt.Println("luna-l2: invalid boot drive")
				}
				i++
			default:
				shared.Filename = arg
			}
		}

		boot:
		bios.Splash()

		switch shared.BootDrive {
		case 0:
			bios.WriteLine("Booting from hard disk...", 255, 0)
			bios.LoadSector(0, 0, RequireDevicePresent)
			shared.DriveNumber = 0
		case 1:
			bios.WriteLine("Booting from SD...", 255, 0)
			bios.LoadSector(1, 0, RequireDevicePresent)
			shared.DriveNumber = 1
		case 2:
			bios.WriteLine("Booting from DVD...", 255, 0)
			bios.LoadSector(2, 0, RequireDevicePresent)
			shared.DriveNumber = 2
		default:
			bios.WriteLine("No bootable device", 255, 0)
			return
		}
		RequireDevicePresent = false

		// Initialize components
		go network.NetController()
		go audio.AudioController()
		go rtc.RTCController()
		// go pit.PITController()
		copy(Memory[0x6FFF0000:], []byte{ // IDT setup
			0x01, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x02, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x03, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x04, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x05, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x06, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x07, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x08, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x09, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x0A, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x0B, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x0C, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x0D, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x0E, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x0F, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x10, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x11, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x12, 0x00, 0x00, 0x00, 0x00, 0x00,
		})

		// Execute
		execute()

		if BIOS_REBOOT == true {
			BIOS_REBOOT = false
			goto boot
		} else if BIOS_SHUTDOWN == true {
			os.Exit(0)
		}
	}()	
	InitializeWindow()
}
