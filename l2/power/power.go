package power

import (
	"time"
	"github.com/distatus/battery"
	"luna_l2/shared"
)

var MemoryPower [1]byte

func PowerController() {
	for {
		batteries, err := battery.GetAll()
		if err != nil {
			percent := 100
			MemoryPower[0] = byte(percent & 0xFF)
			goto DONE	
		}
		if len(batteries) < 1 {
			percent := 100
			MemoryPower[0] = byte(percent & 0xFF)	
			goto DONE
		}

		for _, b := range batteries {
			percent := uint8((b.Current / b.Full) * 100)
			MemoryPower[0] = byte(percent & 0xFF)
			if percent < 1 {
				shared.RaiseInterrupt(0x6)
			}
		}
			
		DONE:
		time.Sleep(15 * time.Millisecond)
	}
}
