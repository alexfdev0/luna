package shared

import (
	"math/rand"
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
var MemoryNetwork *[4122]byte
var MemoryRTC *[6]byte
var MemoryPIT *[8]byte

func Mapper(address uint32) byte {
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
	case address >= 0x7001A644 && address <= 0x7001B65D:
		return (*MemoryNetwork)[address - 0x7001A644]
	case address >= 0x7001B65E && address <= 0x7001B663:
		return (*MemoryRTC)[address - 0x7001B65E]
	}
	return byte(rand.Intn(0xFF - 0x00) + 0x00)
}

func MapperWrite(address uint32, content byte) {
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
	case address >= 0x7001A644 && address <= 0x7001B65C:
		(*MemoryNetwork)[address - 0x7001A644] = content
	case address >= 0x7001B65E && address <= 0x7001B663:
		(*MemoryRTC)[address - 0x7001B65E] = content
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
