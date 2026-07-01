package main

import (	
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
	"luna_l2/rtc"
	"luna_l2/keyboard"
	"luna_l2/pit"
	"luna_l2/power"	
)

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
	{0x0019, "E12", 0},
	{0x001a, "E13", 0},
	{0x001b, "E14", 0},
	{0x001c, "SP", 0},
	{0x001d, "PC", 0},
	{0x001e, "IRV", 0},
	{0x001f, "IR", 0},
	{0x0020, "B", 0},
		
}

var Memory []byte 
const (
	MEMSIZE uint32 = 0x70000000
	MEMCAP uint32 = 0x6FFFFFFF
)

func setRegister(address uint32, value uint32) {
	if address < uint32(len(Registers)) {
		if shared.Bits32 == false && address != 0x001f {
			Registers[address].Value = uint32(uint16(value))
		} else {
			Registers[address].Value = value
		}
	}	
}

func getRegister(address uint32) uint32 {
	if address < uint32(len(Registers)) {
		return Registers[address].Value
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

var ClockSpeed int64 = 33000000
var BIOS_REBOOT bool = false
var BIOS_SHUTDOWN bool = false

func Log(text string) {
	if shared.LogOn == true {
		fmt.Println("\033[33m" + fmt.Sprintf("0x%08x: ", getRegister(0x001d)) + text + "\033[0m")
	}	
}

var accumulated int64
func stall(cycles int64) { 
	cycleTime := int64(int(time.Second)) / ClockSpeed
	accumulated += cycles
	if accumulated >= 66000 {
		time.Sleep(time.Duration(cycleTime * accumulated))
		accumulated = 0
	}	
}

var ins int64 = 0

func execute() {
	II_HALT := func(ProgramCounter uint32, next uint32) bool {
		now := ProgramCounter
		bios.IntWrapper(0x7, next)
		if shared.Debug == true {
			setRegister(0x001d, next)
		} else {
			if getRegister(0x001d) == now {
				return true
			} else {
				fmt.Println(now, getRegister(0x001d))
			}
		}
		return false
	}
	for {
		ProgramCounter := getRegister(0x001d)
		op := shared.Mapper(ProgramCounter)	

		// Handle interrupts	
		var IntHandled bool
		for i := 0; i < 32; i++ {
			if (getRegister(0x001f) & (1 << i)) != 0 {
				code := i + 1
				IntHandled = true
				setRegister(0x001f, getRegister(0x001f) &^ (1 << i))
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
		case 0x01:
			// MOV
			mode := shared.Mapper(ProgramCounter + 1)
			dst := shared.Mapper(ProgramCounter + 2)

			if mode == 0x01 {
				var imm uint32 = 0
				var next uint32 = 0
				if shared.Bits32 == false {
					imm = uint32(uint16(shared.Mapper(ProgramCounter + 3)) << 8 | uint16(shared.Mapper(ProgramCounter + 4)))
					next = ProgramCounter + 5
				} else {
					imm = uint32(shared.Mapper(ProgramCounter + 3)) << 24 | uint32(shared.Mapper(ProgramCounter + 4)) << 16 | uint32(shared.Mapper(ProgramCounter + 5)) << 8 | uint32(shared.Mapper(ProgramCounter + 6))
					next = ProgramCounter + 7
				}
				setRegister(uint32(dst), imm)
				setRegister(0x001d, next)
				Log("mov " + getRegisterName(uint32(dst)) + ", " + fmt.Sprintf("0x%08x", imm))
			} else if mode == 0x02 {
				frm := uint32(shared.Mapper(ProgramCounter + 3))
				setRegister(uint32(dst), uint32(getRegister(frm)))
				setRegister(0x001d, ProgramCounter + 4)
				Log("mov " + getRegisterName(uint32(dst)) + ", " + getRegisterName(frm))
			}	
			stall(4)
		case 0x02:
			// HLT
			Log("hlt")
			now := ProgramCounter
			for {
				if getRegister(0x001d) != now || getRegister(0x001f) != 0 {
					break	
				}	
				time.Sleep(time.Duration(15) * time.Millisecond)
			}
			setRegister(0x001d, ProgramCounter + 1)
		case 0x03:
			// JMP	
			mode := shared.Mapper(ProgramCounter + 1)

			if mode == 0x01 {
				var loc uint32 = 0
				if shared.Bits32 == false {
					loc = uint32(uint16(shared.Mapper(ProgramCounter + 2)) << 8 | uint16(shared.Mapper(ProgramCounter + 3)))
				} else {
					loc = uint32(shared.Mapper(ProgramCounter + 2)) << 24 | uint32(shared.Mapper(ProgramCounter + 3)) << 16 | uint32(shared.Mapper(ProgramCounter + 4)) << 8 | uint32(shared.Mapper(ProgramCounter + 5))
				}
				Log("jmp " + fmt.Sprintf("0x%08x", loc))
				setRegister(0x001d, loc)	
			} else if mode == 0x02 {
				frm := uint32(shared.Mapper(ProgramCounter + 2))
				loc := getRegister(frm)
				Log("jmp " + getRegisterName(frm))
				setRegister(0x001d, loc)	
			}
			stall(8)
		case 0x04:
			// INT
			var code uint32 = 0
			var next uint32 = 0
	
			if shared.Bits32 == false {
				code = uint32(uint16(shared.Mapper(ProgramCounter + 1)) << 8 | uint16(shared.Mapper(ProgramCounter + 2)))
				next = ProgramCounter + 3
			} else {
				code = uint32(shared.Mapper(ProgramCounter + 1)) << 24 | uint32(shared.Mapper(ProgramCounter + 2)) << 16 | uint32(shared.Mapper(ProgramCounter + 3)) << 8 | uint32(shared.Mapper(ProgramCounter + 4))
				next = ProgramCounter + 5
			}	
	
			shared.RaiseInterrupt(code)
			setRegister(0x001d, next)	

			Log("int " + fmt.Sprintf("0x%08x", code))	
			stall(34)
		case 0x05:
			// JNZ
			// jnz <mode (01 or 02)> <check register> <loc (register or raw addr)>
			mode := shared.Mapper(ProgramCounter + 1)
			checkRegister := shared.Mapper(ProgramCounter + 2)
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
				Log("jnz " + getRegisterName(uint32(checkRegister)) + ", " + fmt.Sprintf("0x%08x", loc))
			} else if mode == 0x02 {
				frm := uint32(shared.Mapper(ProgramCounter + 3))
				loc = getRegister(frm)
				not = ProgramCounter + 4
				Log("jnz " + getRegisterName(uint32(checkRegister)) + ", " + getRegisterName(frm))
			}

			if getRegister(uint32(checkRegister)) != 0 {
				setRegister(0x001d, loc)
			} else {
				setRegister(0x001d, not)
			}
			stall(8)
		case 0x06:
			// NOP
			setRegister(0x001d, ProgramCounter + 1)
			Log("nop")
			stall(1)
		case 0x07:
			// CMP
			// Syntax: CMP <to> <r1> <r2>
			to := shared.Mapper(ProgramCounter + 1)
			first := shared.Mapper(ProgramCounter + 2)
			second := shared.Mapper(ProgramCounter + 3)
			Log("cmp " + getRegisterName(uint32(to)) + ", " + getRegisterName(first) + ", " + getRegisterName(second))

			if getRegister(uint32(first)) == getRegister(uint32(second)) {
				setRegister(uint32(to), uint32(1))
			} else {
				setRegister(uint32(to), uint32(0))
			}
			setRegister(0x001d, ProgramCounter + 4)
			stall(4)
		case 0x08:
			// JZ
			// jz <mode (01 or 02)> <check register> <loc (register or raw addr)>
			mode := shared.Mapper(ProgramCounter + 1)
			checkRegister := shared.Mapper(ProgramCounter + 2)
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
				frm := uint32(shared.Mapper(ProgramCounter + 3))
				loc = getRegister(frm)
				not = ProgramCounter + 4
				Log("jz " + getRegisterName(checkRegister) + ", " + getRegisterName(frm))
			}

			if getRegister(uint32(checkRegister)) == 0 {
				setRegister(0x001d, loc)
			} else {
				setRegister(0x001d, not)
			}
			stall(8)
		case 0x09:
			// INC
			// inc <register>
			register := uint32(shared.Mapper(ProgramCounter + 1))
			setRegister(register, getRegister(register) + 1)
			setRegister(0x001d, ProgramCounter+2)
			Log("inc " + getRegisterName(register))
			stall(1)
		case 0x0a:
			// DEC
			// dec <register>
			register := uint32(shared.Mapper(ProgramCounter + 1))
			setRegister(register, getRegister(register) - 1)
			setRegister(0x001d, ProgramCounter + 2)
			Log("dec " + getRegisterName(register))
			stall(1)
		case 0x0b:
			// PUSH
			// push <mode> <immediate or register>
			mode := shared.Mapper(ProgramCounter + 1)
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
				setRegister(0x001d, next)
				Log("push " + fmt.Sprintf("0x%08x", value))
			} else if mode == 0x2 {
				value = getRegister(uint32(shared.Mapper(ProgramCounter + 2)))
				setRegister(0x001d, ProgramCounter + 3)
				Log("push " + getRegisterName(uint32(shared.Mapper(ProgramCounter + 2))))
			}	
			sp := getRegister(0x001c)
			if shared.Bits32 == false {
				sp = shared.Clamp(sp - 2, 0, MEMCAP)
				shared.MapperWrite(sp, byte(value & 0xFF))
				shared.MapperWrite(sp + 1, byte(value >> 8))
			} else {
				sp = shared.Clamp(sp - 4, 0, MEMCAP)
				shared.MapperWrite(sp, byte(value & 0xFF))
				shared.MapperWrite(sp + 1, byte(value >> 8))
				shared.MapperWrite(sp + 2, byte(value >> 16))
				shared.MapperWrite(sp + 3, byte(value >> 24))
			}	
			setRegister(0x001c, uint32(sp))	
			stall(2)
		case 0x0c:
			// POP
			// pop <register>	
			register := shared.Mapper(ProgramCounter + 1)
			sp := getRegister(0x001c)
			var value uint32
			if shared.Bits32 == false {
				value = uint32(uint16(shared.Mapper(sp)) | uint16(shared.Mapper(sp + 1)) << 8) 
			} else {	
				value = uint32(shared.Mapper(sp)) | uint32(shared.Mapper(sp + 1)) << 8 | uint32(shared.Mapper(sp + 2)) << 16 | uint32(shared.Mapper(sp + 3)) << 24
			}
			Log("value: " + fmt.Sprintf("0x%08x", value))
			setRegister(uint32(register), uint32(value))
			if shared.Bits32 == false {
				sp = shared.Clamp(sp + 2, 0, MEMCAP)
			} else {
				sp = shared.Clamp(sp + 4, 0, MEMCAP)
			}
			setRegister(0x001c, uint32(sp))
			setRegister(0x001d, ProgramCounter + 2)
			Log("pop " + getRegisterName(register))
			stall(2)
		case 0x0d:
			// ADD
			// add <register> <register> <register>
			toregister := shared.Mapper(ProgramCounter + 1)
			regone := shared.Mapper(ProgramCounter + 2)
			regtwo := shared.Mapper(ProgramCounter + 3)
			setRegister(uint32(toregister), getRegister(uint32(regone)) + getRegister(uint32(regtwo)))
			setRegister(0x001d, ProgramCounter + 4)
			Log("add " + getRegisterName(toregister) + ", " + getRegisterName(regone) + ", " + getRegisterName(regtwo))
			stall(7)
		case 0x0e:
			// SUB
			// SUB <register> <register> <register>
			toregister := shared.Mapper(ProgramCounter + 1)
			regone := shared.Mapper(ProgramCounter + 2)
			regtwo := shared.Mapper(ProgramCounter + 3)
			setRegister(uint32(toregister), getRegister(uint32(regone)) - getRegister(uint32(regtwo)))
			setRegister(0x001d, ProgramCounter + 4)
			Log("sub " + getRegisterName(toregister) + ", " + getRegisterName(regone) + ", " + getRegisterName(regtwo))
			stall(7)
		case 0x0f:
			// MUL
			// mul <register> <register> <register>
			toregister := shared.Mapper(ProgramCounter + 1)
			regone := shared.Mapper(ProgramCounter + 2)
			regtwo := shared.Mapper(ProgramCounter + 3)
			setRegister(uint32(toregister), getRegister(uint32(regone)) * getRegister(uint32(regtwo)))
			setRegister(0x001d, ProgramCounter + 4)
			Log("mul " + getRegisterName(toregister) + ", " + getRegisterName(regone) + ", " + getRegisterName(regtwo))
			stall(70)
		case 0x10:
			// DIV
			// div <register> <register> <register>
			toregister := shared.Mapper(ProgramCounter + 1)
			regone := shared.Mapper(ProgramCounter + 2)
			regtwo := shared.Mapper(ProgramCounter + 3)
			if (getRegister(uint32(regtwo)) == 0) {
				setRegister(0x0001, uint32(op))
				setRegister(0x0002, getRegister(0x001d))
				stall(10)
				if II_HALT(ProgramCounter, ProgramCounter + 4) == true {
					return
				}
			} else {
				setRegister(uint32(toregister), getRegister(uint32(regone)) / getRegister(uint32(regtwo)))
				setRegister(0x001d, ProgramCounter + 4)
				stall(140)
			}
			Log("div " + getRegisterName(toregister) + ", " + getRegisterName(regone) + ", " + getRegisterName(regtwo))
		case 0x11:
			// IGT
			// igt <register> <register> <register>
			toregister := shared.Mapper(ProgramCounter + 1)
			regone := shared.Mapper(ProgramCounter + 2)
			regtwo := shared.Mapper(ProgramCounter + 3)
			if getRegister(uint32(regone)) > getRegister(uint32(regtwo)) {
				setRegister(uint32(toregister), uint32(1))
			} else {
				setRegister(uint32(toregister), uint32(0))
			}
			setRegister(0x001d, ProgramCounter + 4)
			Log("igt " + getRegisterName(toregister) + ", " + getRegisterName(regone) + ", " + getRegisterName(regtwo))
			stall(4)
		case 0x12:
			// ILT
			// ilt <register> <register> <register>
			toregister := shared.Mapper(ProgramCounter + 1)
			regone := shared.Mapper(ProgramCounter + 2)
			regtwo := shared.Mapper(ProgramCounter + 3)
			if getRegister(uint32(regone)) < getRegister(uint32(regtwo)) {
				setRegister(uint32(toregister), uint32(1))
			} else {
				setRegister(uint32(toregister), uint32(0))
			}
			setRegister(0x001d, ProgramCounter + 4)
			Log("ilt " + getRegisterName(toregister) + ", " + getRegisterName(regone) + ", " + getRegisterName(regtwo))
			stall(4)
		case 0x13:
			// AND
			// and <register> <register> <register>
			toregister := shared.Mapper(ProgramCounter + 1)
			regone := shared.Mapper(ProgramCounter + 2)
			regtwo := shared.Mapper(ProgramCounter + 3)
			setRegister(uint32(toregister), getRegister(uint32(regone)) & getRegister(uint32(regtwo)))	
			setRegister(0x001d, ProgramCounter + 4)
			Log("and " + getRegisterName(toregister) + ", " + getRegisterName(regone) + ", " + getRegisterName(regtwo))
			stall(1)
		case 0x14:
			// OR
			// or <register> <register> <register>
			toregister := shared.Mapper(ProgramCounter + 1)
			regone := shared.Mapper(ProgramCounter + 2)
			regtwo := shared.Mapper(ProgramCounter + 3)
			setRegister(uint32(toregister), getRegister(uint32(regone)) | getRegister(uint32(regtwo)))	
			setRegister(0x001d, ProgramCounter + 4)
			Log("or " + getRegisterName(toregister) + ", " + getRegisterName(regone) + ", " + getRegisterName(regtwo))
			stall(1)
		case 0x15:
			// NOT
			// not <register> <register>
			toregister := shared.Mapper(ProgramCounter + 1)
			regone := shared.Mapper(ProgramCounter + 2)
			setRegister(uint32(uint32(toregister)), ^getRegister(uint32(regone)))	
			setRegister(0x001d, ProgramCounter + 3)
			Log("not " + getRegisterName(toregister) + ", " + getRegisterName(regone))
			stall(1)
		case 0x16:
			// XOR
			// xor <register> <register> <register>
			toregister := shared.Mapper(ProgramCounter + 1)
			regone := shared.Mapper(ProgramCounter + 2)
			regtwo := shared.Mapper(ProgramCounter + 3)
			setRegister(uint32(toregister), getRegister(uint32(regone)) ^ getRegister(uint32(regtwo)))	
			setRegister(0x001d, ProgramCounter + 4)
			Log("xor " + getRegisterName(toregister) + ", " + getRegisterName(regone) + ", " + getRegisterName(regtwo))
			stall(6)
		case 0x17:
			// LOD
			// lod <addr (register)> <destination register>	
			addr := getRegister(uint32(shared.Mapper(ProgramCounter + 1)))
			toregister := uint32(shared.Mapper(ProgramCounter+2))
			setRegister(toregister, uint32(shared.Mapper(addr)))
			setRegister(0x001d, ProgramCounter + 3)
			Log("lod " + getRegisterName(uint32(shared.Mapper(ProgramCounter + 1))) + ", " + getRegisterName(toregister) + " (" + fmt.Sprintf("0x%02x", shared.Mapper(addr)) + ")")
			stall(100)
		case 0x18:
			// STR16
			// str16 <addr (register)> <value (register)>	
			addr := getRegister(uint32(shared.Mapper(ProgramCounter + 1)))
			value := uint32(shared.Mapper(ProgramCounter + 2))
			shared.MapperWrite(addr, byte(getRegister(value) >> 8))
			shared.MapperWrite(addr + 1, byte(getRegister(value) & 0xFF))	
			setRegister(0x001d, ProgramCounter + 3)
			Log("str16 " + getRegisterName(uint32(shared.Mapper(ProgramCounter + 1))) + ", " + getRegisterName(value))
			stall(100)
		case 0x19:
			// LOD16
			// lod16 <addr (register)> <destination register>
			addr := getRegister(uint32(shared.Mapper(ProgramCounter + 1)))
			toregister := uint32(shared.Mapper(ProgramCounter + 2))
			setRegister(toregister, uint32(uint16(shared.Mapper(addr)) << 8 | uint16(shared.Mapper(addr + 1))))
			setRegister(0x001d, ProgramCounter + 3)
			Log("value: " + fmt.Sprintf("0x%08x", getRegister(toregister)))
			Log("lod16 " + getRegisterName(uint32(shared.Mapper(ProgramCounter + 1))) + ", " + getRegisterName(toregister))
			stall(100)
		case 0x1a:
			// SET
			// set <00 or 01>
			mode := uint32(shared.Mapper(ProgramCounter + 1))
			if mode == 0 {
				shared.Bits32 = false
				Log("16 bit mode")
			} else if mode == 1 {
				shared.Bits32 = true
				Log("32 bit mode")
			}
			setRegister(0x001d, ProgramCounter + 2)
			stall(1)
		case 0x1b:
			// STR
			// str <addr> <register>
			addr := getRegister(uint32(shared.Mapper(ProgramCounter + 1)))
			value := uint32(shared.Mapper(ProgramCounter + 2))
			shared.MapperWrite(addr, byte(getRegister(value)))
			setRegister(0x001d, ProgramCounter + 3)
			stall(100)
		case 0x1c:
			// SHL
			// shl <dest> <value> <by>
			dest := uint32(shared.Mapper(ProgramCounter + 1))
			value := getRegister(uint32(shared.Mapper(ProgramCounter + 2)))
			by := getRegister(uint32(shared.Mapper(ProgramCounter + 3)))
			setRegister(dest, uint32(value) << uint32(by))
			setRegister(0x001d, ProgramCounter + 4)
			stall(95)
		case 0x1d:
			// SHR
			// shr <dest> <value> <by>
			dest := uint32(shared.Mapper(ProgramCounter + 1))
			value := getRegister(uint32(shared.Mapper(ProgramCounter + 2)))
			by := getRegister(uint32(shared.Mapper(ProgramCounter + 3)))
			setRegister(dest, uint32(value) >> uint32(by))
			setRegister(0x001d, ProgramCounter + 4)
			stall(95)
		case 0x1e:
			// LOD32
			// lod32 <addr (register)> <destination register>
			addr := getRegister(uint32(shared.Mapper(ProgramCounter + 1)))
			toregister := uint32(shared.Mapper(ProgramCounter + 2))	
			setRegister(toregister, uint32(shared.Mapper(addr)) << 24 | uint32(shared.Mapper(addr + 1)) << 16 | uint32(shared.Mapper(addr + 2)) << 8 | uint32(shared.Mapper(addr + 3)))
			setRegister(0x001d, ProgramCounter + 3)
			Log("value: " + fmt.Sprintf("0x%08x", getRegister(toregister)))
			Log("lod32 " + getRegisterName(uint32(shared.Mapper(ProgramCounter + 1))) + ", " + getRegisterName(toregister))
			stall(100)
		case 0x1f:
			// STR32
			// str32 <addr (register)> <value (register)>	
			addr := getRegister(uint32(shared.Mapper(ProgramCounter + 1)))
			value := uint32(shared.Mapper(ProgramCounter + 2))
			shared.MapperWrite(addr, byte(getRegister(value) >> 24))
			shared.MapperWrite(addr + 1, byte(getRegister(value) >> 16))
			shared.MapperWrite(addr + 2, byte(getRegister(value) >> 8))
			shared.MapperWrite(addr + 3, byte(getRegister(value) & 0xFF))
			setRegister(0x001d, ProgramCounter + 3)
			Log("str32 " + getRegisterName(uint32(shared.Mapper(ProgramCounter + 1))) + ", " + getRegisterName(value))
			stall(100)
		case 0x20:
			// MOD
			// mod <dest> <reg1> <reg2>
			toregister := shared.Mapper(ProgramCounter + 1)
			regone := shared.Mapper(ProgramCounter + 2)
			regtwo := shared.Mapper(ProgramCounter + 3)
			setRegister(uint32(toregister), getRegister(uint32(regone)) % getRegister(uint32(regtwo)))
			setRegister(0x001d, ProgramCounter + 4)
			Log("mod " + getRegisterName(toregister) + ", " + getRegisterName(regone) + ", " + getRegisterName(regtwo))
			stall(70)
		case 0x21:
			shared.LogOn = true
			shared.Debug = true
			Log("\033[31m----- BREAKPOINT -----\033[33m")
			setRegister(0x001d, ProgramCounter + 1)
		default:
			setRegister(0x0001, uint32(op))
			setRegister(0x0002, getRegister(0x001d))
			Log("\033[31mIllegal instruction 0x" + fmt.Sprintf("%08x", uint32(op)) + "\033[33m")
			if II_HALT(ProgramCounter, ProgramCounter + 1) == true {
				return
			}
		}

		if shared.Debug == true {
			for i := 0; i < len(Registers); i++ {
				Register := Registers[i]
				fmt.Println(Register.Name + ": " + fmt.Sprintf("0x%8x", Register.Value))
			}
			bufio.NewReader(os.Stdin).ReadBytes('\n')	
		}	
	}
}

var RequireDevicePresent bool = true
func main() {
	runtime.LockOSThread()
 
	shared.Registers = &Registers
	shared.Memory = &Memory
	shared.MemoryMouse = &keyboard.MemoryMouse
	shared.MemoryKeyboard = &keyboard.MemoryKeyboard
	shared.MemoryRTC = &rtc.MemoryRTC
	shared.MemoryPIT = &pit.MemoryPIT
	shared.MemoryPower = &power.MemoryPower

	var MemorySetting uint32 = 0x70000000
	var GPU string = "g1x"
	var APU string = "s1"

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
		case "-ram":
			if i + 1 >= len(os.Args) { fmt.Println("Not enough arguments to -ram"); i++; continue }
			memset, err := strconv.ParseInt(os.Args[i + 1], 0, 32)
			if err != nil && memset > 0x70000000 {
				fmt.Println("Invalid RAM amount")
				i++
				continue
			}
			MemorySetting = uint32(memset)
			i++
		case "--log":
			shared.LogOn = true
		case "--debug":
			shared.Debug = true
			shared.LogOn = true
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
		case "-gpu":
			if i + 1 >= len(os.Args) { fmt.Println("Not enough args to -gpu"); i++; continue; }
			i++
			GPU = os.Args[i]
		case "-apu":
			if i + 1 >= len(os.Args) { fmt.Println("Not enough args to -apu"); i++; continue; }
			i++
			APU = os.Args[i]
		default:
			shared.Filename = arg
		}
	}

	Memory = make([]byte, MemorySetting)
	shared.MEMCAP = MemorySetting - 1

	// Initialize components
	audio.AudioController(APU)
	go rtc.RTCController()
	go pit.PITController()
	go power.PowerController()	

	go func() {
		if video.Ready == false {	
			for {
				if video.Ready == true {
					break
				} else {
					time.Sleep(500)
				}
			}
		}	

		boot:
		bios.Splash()
		boot2:

		RetryBoot := func() {
			bios.WriteLine("Could not read the boot disk\n", 255, 0)
			shared.BootDrive++
		}

		switch shared.BootDrive {
		case 0:
			bios.WriteLine("Booting from hard disk...", 255, 0)

			if shared.Filename == "" || bios.LoadSector(0, 0, 0) == false {
				RetryBoot()
				goto boot2
			}
			shared.DriveNumber = 0
		case 1:
			bios.WriteLine("Booting from SD...", 255, 0)

			if shared.SDFilename == "" || bios.LoadSector(1, 0, 0) == false {
				RetryBoot()
				goto boot2
			}
			shared.DriveNumber = 1
		case 2:
			bios.WriteLine("Booting from DVD...", 255, 0)

			if shared.OpticalFilename == "" || bios.LoadSector(2, 0, 0) == false {
				RetryBoot()
				goto boot2
			}
			shared.DriveNumber = 2
		default:
			bios.WriteLine("No bootable device", 255, 0)
			return
		}	
		// IDT increments of +6 for every interrupt, starting at 0	

		// Execute
		execute()

		if BIOS_REBOOT == true {
			BIOS_REBOOT = false
			goto boot
		} else if BIOS_SHUTDOWN == true {
			os.Exit(0)
		}
	}()	
	video.InitializeWindow(GPU)
}
