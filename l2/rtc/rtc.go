package rtc

import (
	"time"
)

var MemoryRTC [6]byte

func RTCController() {
	for {
		now := time.Now().UTC()
		MemoryRTC[0x0000] = byte(now.Second())
		MemoryRTC[0x0001] = byte(now.Minute())
		MemoryRTC[0x0002] = byte(now.Hour())
		MemoryRTC[0x0003] = byte(now.Day())
		MemoryRTC[0x0004] = byte(now.Month())
		MemoryRTC[0x0005] = byte(now.Year() - 2000)
		time.Sleep(time.Duration(1000) * time.Millisecond) 
	}
}
