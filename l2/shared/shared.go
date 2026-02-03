package shared

import (
	"math/rand"
	"luna_l2/video"
)
/* 
shared.go:
main functions that are exported so other independent libraries can use them

*/

const (
	MEMSIZE uint32 = 0x70000000
	MEMCAP uint32 = 0x6FFFFFFF
)

type Register struct {
	Address uint32
	Name    string
	Value   uint32
}

var Registers *[]Register
var Memory *[0x70000000]byte
var MemoryVideo *[64000]byte
var MemoryAudio *[10]byte
var MemoryMouse *[8]byte
var MemoryKeyboard *[1]byte
var MemoryNetwork *[22]byte
var MemoryRTC *[6]byte
var MemoryPIT *[8]byte

func Mapper(address uint32) byte {
	if Bits32 == true {
		switch {
		case address >= 0x00000000 && address <= MEMCAP:
			return (*Memory)[address]
		case address >= 0x70000000 && address <= 0x7000F9FF:
			return (*MemoryVideo)[address - MEMSIZE]
		case address >= 0x7000FA00 && address <= 0x7000FA09:
			return (*MemoryAudio)[address - 0x7000FA00]
		case address >= 0x7000FA0A && address <= 0x7000FA11:
			return (*MemoryMouse)[address - 0x7000FA0A]
		case address >= 0x7000FA12 && address <= 0x7000FA12:
			return (*MemoryKeyboard)[address - 0x7000FA12]
		case address >= 0x7000FA13 && address <= 0x7000FA1A:
			return (*MemoryPIT)[address - 0x7000FA13]
		case address >= 0x7001A644 && address <= 0x7001A659:
			return (*MemoryNetwork)[address - 0x7001A644]
		case address >= 0x7001B65E && address <= 0x7001B663:
			return (*MemoryRTC)[address - 0x7001B65E]
		}
	} else {
		switch {
		case address >= 0x0000 && address <= 0xEFFF:
			return (*Memory)[address]
		case address >= 0xF000 && address <= 0xF009:
			return (*MemoryAudio)[address - 0xF000]
		case address >= 0xFA0A && address <= 0xFA11:
			return (*MemoryMouse)[address - 0xFA0A]
		case address == 0xFA12:
			return (*MemoryKeyboard)[address - 0xFA12]
		case address >= 0xFA13 && address <= 0xFA1A:
			return (*MemoryPIT)[address - 0xFA13]
		case address >= 0xFA1B && address <= 0xFA30:
			return (*MemoryNetwork)[address - 0xFA30]
		case address >= 0xFA31 && address <= 0xFA36:
			return (*MemoryRTC)[address - 0xFA31]
		case address >= 0xFA37 && address <= 0xFC36:
			// IDT
			return (*Memory)[0x6FFF0000 + (address - 0xFA37)]
		case address >= 0xFE00 && address <= 0xFFFF:
			if GetRegister(0x001F) <= 124 {
				return (*MemoryVideo)[video.Clamp((GetRegister(0x001f) * 0x200) + (address - 0xFE00), 0, 63999)]
			}
		}
	}
	return byte(rand.Intn(0xFF - 0x00) + 0x00)
}

func MapperWrite(address uint32, content byte) {
	if Bits32 == true {
		switch {
		case address >= 0x00000000 && address <= MEMCAP:
			(*Memory)[address] = content
		case address >= 0x70000000 && address <= 0x7000F9FF:
			(*MemoryVideo)[address - MEMSIZE] = content
		case address >= 0x7000FA00 && address <= 0x7000FA09:
			(*MemoryAudio)[address - 0x7000FA00] = content
		case address >= 0x7000FA0A && address <= 0x7000FA11:
			(*MemoryMouse)[address - 0x7000FA0A] = content
		case address >= 0x7000FA12 && address <= 0x7000FA12:
			(*MemoryKeyboard)[address - 0x7000FA12] = content
		case address >= 0x7000FA13 && address <= 0x7000FA1A:
			(*MemoryPIT)[address - 0x7000FA13] = content
		case address >= 0x7001A644 && address <= 0x7001A659:
			(*MemoryNetwork)[address - 0x7001A644] = content
		case address >= 0x7001B65E && address <= 0x7001B663:
			(*MemoryRTC)[address - 0x7001B65E] = content
		}
	} else {
		switch {
		case address >= 0x0000 && address <= 0xEFFF:
			(*Memory)[address] = content
		case address >= 0xF000 && address <= 0xF009:
			(*MemoryAudio)[address - 0xF000] = content
		case address >= 0xFA0A && address <= 0xFA11:
			(*MemoryMouse)[address - 0xFA0A] = content
		case address == 0xFA12:
			(*MemoryKeyboard)[address - 0xFA12] = content
		case address >= 0xFA13 && address <= 0xFA1A:
			(*MemoryPIT)[address - 0xFA13] = content
		case address >= 0xFA1B && address <= 0xFA30:
			(*MemoryNetwork)[address - 0xFA30] = content
		case address >= 0xFA31 && address <= 0xFA36:
			(*MemoryRTC)[address - 0xFA31] = content
		case address >= 0xFA37 && address <= 0xFC36:
			// IDT
			(*Memory)[0x6FFF0000 + (address - 0xFA37)] = content
		case address >= 0xFE00 && address <= 0xFFFF:
			if GetRegister(0x001F) <= 124 {
				(*MemoryVideo)[(GetRegister(0x001f) * 0x200) + (address - 0xFE00)] = content
			}
		}
	}	
}

func SetRegister(address uint32, value uint32) {
	for i := range (*Registers) {
		if (*Registers)[i].Address == address {
			if Bits32 == false {
				(*Registers)[i].Value = uint32(uint16(value))
			} else {
				(*Registers)[i].Value = value
			}
		}
	}
}

func GetRegister(address uint32) uint32 {
	for _, register := range (*Registers) {
		if register.Address == address {
			return register.Value
		}
	}
	return 0x0000
}

func RaiseInterrupt(code uint32)  {
	code--
	if code >= 32 {
		return
	}
	SetRegister(0x001e, GetRegister(0x001e) | (1 << code))
}

var Bits32 bool = false
var Filename string = ""
var SDFilename string = ""
var OpticalFilename string = ""
var DriveNumber int = 0
var BootDrive int = 0
var IntRaiseCode uint32 = 0
