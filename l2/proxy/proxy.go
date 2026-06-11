package proxy

var VideoPrintChar func(rune, byte, byte)
var VideoSetCursor func(int, int)
var VideoGetCursor func() (int, int)
var VideoClearVideoMemory func()
var VideoWriteVideoMemory func(uint32, byte)
var VideoReadVideoMemory func(uint32) byte
var AudioWriteAudioMemory func(uint32, byte)
var AudioReadAudioMemory func(uint32) byte
