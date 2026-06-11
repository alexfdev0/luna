package shared

import (
	"math/rand"
	"cmp"
)
/* 
shared.go:
main functions that are exported so other independent libraries can use them

*/

var (
	MEMSIZE uint32 = 0x70000000
	MEMCAP uint32 = 0x6FFFFFFF
)

type Register struct {
	Address uint32
	Name    string
	Value   uint32
}

var Registers *[]Register
var Memory *[]byte
var MemoryAudio *[10]byte
var MemoryMouse *[8]byte
var MemoryKeyboard *[1]byte
var MemoryRTC *[6]byte
var MemoryPIT *[8]byte
var MemoryPower *[4]byte
var WriteVideoMemory = func(addr uint32, content byte) {}
var ReadVideoMemory = func(addr uint32) byte { return 0x00 }

func Mapper(address uint32) byte {
	if Bits32 == true {
		switch {
		case address >= 0x00000000 && address <= MEMCAP:
			return (*Memory)[address]
		case address >= 0x70000000 && address <= 0x7000F9FF:
			return ReadVideoMemory(address - 0x70000000)
		case address >= 0x7000FA00 && address <= 0x7000FA09:
			return (*MemoryAudio)[address - 0x7000FA00]
		case address >= 0x7000FA0A && address <= 0x7000FA11:
			return (*MemoryMouse)[address - 0x7000FA0A]
		case address >= 0x7000FA12 && address <= 0x7000FA12:
			return (*MemoryKeyboard)[address - 0x7000FA12]
		case address >= 0x7000FA13 && address <= 0x7000FA1A:
			return (*MemoryPIT)[address - 0x7000FA13]	
		case address >= 0x7001B65E && address <= 0x7001B663:
			return (*MemoryRTC)[address - 0x7001B65E]
		case address >= 0x7001B664 && address <= 0x7001B667:
			return (*MemoryPower)[address - 0x7001B664]
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
		case address >= 0xFD41 && address <= 0xFD46:
			return (*MemoryRTC)[address - 0xFD41]
		case address >= 0xFA37 && address <= 0xFC36:
			// IDT
			return (*Memory)[0x6FFF0000 + (address - 0xFA37)]
		case address >= 0xFC37 && address <= 0xFC3A:
			return (*MemoryPower)[address - 0xFC37]
		case address >= 0xFE00 && address <= 0xFFFF:
			if GetRegister(0x001F) <= 124 {
				return ReadVideoMemory(Clamp((GetRegister(0x0020) * 0x200) + (address - 0xFE00), 0, 63999))
			}
		}
	}
	return byte(rand.Intn(0xFF))
}

func MapperWrite(address uint32, content byte) {
	if Bits32 == true {
		switch {
		case address >= 0x00000000 && address <= MEMCAP:
			(*Memory)[address] = content
		case address >= 0x70000000 && address <= 0x7000F9FF:
			WriteVideoMemory(address - 0x70000000, content)
		case address >= 0x7000FA00 && address <= 0x7000FA09:
			(*MemoryAudio)[address - 0x7000FA00] = content
		case address >= 0x7000FA0A && address <= 0x7000FA11:
			(*MemoryMouse)[address - 0x7000FA0A] = content
		case address >= 0x7000FA12 && address <= 0x7000FA12:
			(*MemoryKeyboard)[address - 0x7000FA12] = content
		case address >= 0x7000FA13 && address <= 0x7000FA1A:
			(*MemoryPIT)[address - 0x7000FA13] = content
		case address >= 0x7001B65E && address <= 0x7001B663:
			(*MemoryRTC)[address - 0x7001B65E] = content
		case address >= 0x7001B664 && address <= 0x7001B667:
			(*MemoryPower)[address - 0x7001B664] = content
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
		case address >= 0xFD41 && address <= 0xFD46:
			(*MemoryRTC)[address - 0xFD41] = content
		case address >= 0xFA37 && address <= 0xFC36:
			// IDT
			(*Memory)[0x6FFF0000 + (address - 0xFA37)] = content
		case address >= 0xFC37 && address <= 0xFC3A:
			(*MemoryPower)[address - 0xFC37] = content
		case address >= 0xFE00 && address <= 0xFFFF:
			if GetRegister(0x001F) <= 124 {
				WriteVideoMemory((GetRegister(0x0020) * 0x200) + (address - 0xFE00), content)
			}
		}
	}	
}

func SetRegister(address uint32, value uint32) {
	if address < uint32(len((*Registers))) {
		if Bits32 == false && address != 0x001f {
			(*Registers)[address].Value = uint32(uint16(value))
		} else {
			(*Registers)[address].Value = value
		}
	}
}

func GetRegister(address uint32) uint32 {
	if address < uint32(len((*Registers))) {
		return (*Registers)[address].Value
	}	
	return 0x0000
}

func RaiseInterrupt(code uint32)  {
	code--
	if code >= 32 {
		return
	}
	SetRegister(0x001f, GetRegister(0x001f) | (1 << code))
}

func Clamp[T cmp.Ordered](x T, min T, max T) T {
    if x < min {
        return min
    }
    if x > max {
        return max
    }
    return x
}

var Bits32 bool = false
var Filename string = ""
var SDFilename string = ""
var OpticalFilename string = ""
var DriveNumber int = 0
var BootDrive int = 0
var IntRaiseCode uint32 = 0
var LogOn bool = false
var Debug bool = false
