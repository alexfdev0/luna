package pit

import (
	"time"
	"luna_l2/shared"
)

// Memory layout:
	// Bytes 1, 2, 3, 4: programmed countdown value
	// Bytes 5, 6, 7, 8: actual countdown value

var MemoryPIT [8]byte


func PITController() {
	copy(MemoryPIT[:], []byte{ 0x00, 0x00, 0x03, 0xE8, 0x00, 0x00, 0x03, 0xE8, })
	for {
		current :=
			uint32(MemoryPIT[4]) << 24 |
			uint32(MemoryPIT[5]) << 16 |
			uint32(MemoryPIT[6]) << 8 |
			uint32(MemoryPIT[7])

		if current == 0 {
			MemoryPIT[4] = MemoryPIT[0]
			MemoryPIT[5] = MemoryPIT[1]
			MemoryPIT[6] = MemoryPIT[2]
			MemoryPIT[7] = MemoryPIT[3]
			shared.RaiseInterrupt(0x2)
		} else {
			current--
			MemoryPIT[4] = byte(current >> 24)
			MemoryPIT[5] = byte(current >> 16)
			MemoryPIT[6] = byte(current >> 8)
			MemoryPIT[7] = byte(current)
		}

		time.Sleep(time.Millisecond)
	}
}

