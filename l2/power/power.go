package power

import (
	"time"
	"github.com/distatus/battery"
	"luna_l2/shared"
)

var MemoryPower [4]byte

func PowerController() {
	for {
		batteries, err := battery.GetAll()
		if err != nil {
			goto DONE	
		}

		for _, b := range batteries {
			percent := uint8((b.Current / b.Full) * 100)
			MemoryPower[0] = byte(percent >> 24)
			MemoryPower[1] = byte(percent >> 16)
			MemoryPower[2] = byte(percent >> 8)
			MemoryPower[3] = byte(percent & 0xFF)
			if percent < 1 {
				shared.RaiseInterrupt(0x6)
			}
		}
			
		DONE:
		time.Sleep(15 * time.Millisecond)
	}
}
