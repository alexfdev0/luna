package audio

import (
	"luna_l2/component"
	"luna_l2/proxy"
	"runtime"
)

var AudioComponent component.Component

var CommonComponentPathPrefix string = "/usr/local/lib/l2/audio/"
var WindowsComponentPathPrefix string = "C:\\Program Files (x86)\\Luna L2\\lib\\l2\\audio\\"

func AudioController(APU string) {
	prefix := CommonComponentPathPrefix
	ext := ".so"
	if runtime.GOOS == "windows" {
		prefix = WindowsComponentPathPrefix
		ext = ".dll"
	}
	
	AudioComponent = component.InitializeComponent(prefix + APU + ext)
	go component.ReturnComponentFunction(AudioComponent, "AudioController").(func())()

	proxy.AudioWriteAudioMemory = component.ReturnComponentFunction(AudioComponent, "WriteAudioMemory").(func(uint32, byte))
	proxy.AudioReadAudioMemory = component.ReturnComponentFunction(AudioComponent, "ReadAudioMemory").(func(uint32) byte)
}
