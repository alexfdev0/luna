package audio

import (
	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"time"
	"luna_l2/shared"
)

var MemoryAudio [10]byte

// MEMORY layout:
	// Byte 1: Play flag
	// Byte 2, 3, 4, 5: 32-bit pointer of audio
	// Byte 6, 7, 8, 9: 32-bit size of audio
	// Byte 10: done flag
type PCMStreamer struct {	
	cursor uint32
}

var Cursor uint32
var Pointer uint32
var Size uint32

func (s *PCMStreamer) Stream(samples [][2]float64) (n int, ok bool) {	
	for i := range samples {
		if uint32(Cursor + 2) > Pointer + Size {
			return i, false
		}

		v := float64(int(int8(shared.Mapper(Cursor)))) / 128.0
		Cursor++
		samples[i][0] = v
		samples[i][1] = v
	}
	return len(samples), true
}

func (s *PCMStreamer) Err() error { return nil }
func Play() {	
	streamer := &PCMStreamer{}	
	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		MemoryAudio[9] = 1
	})))	
}

func AudioController() {
	format := beep.Format{SampleRate: 44100, NumChannels: 2, Precision: 1}
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second / 10))
	for {
		if MemoryAudio[0] == 1 {
			MemoryAudio[0] = 0
			Pointer = uint32(MemoryAudio[1]) << 24 | uint32(MemoryAudio[2]) << 16 | uint32(MemoryAudio[3]) << 8 | uint32(MemoryAudio[4])
			Size = uint32(MemoryAudio[5]) << 24 | uint32(MemoryAudio[6]) << 16 | uint32(MemoryAudio[7]) << 8 | uint32(MemoryAudio[8])
			Cursor = Pointer
			Play()	
		}
		time.Sleep(time.Duration(15) * time.Millisecond)
	}	
}
