package types

type Register struct {
	Address uint32
	Name    string
	Value   uint32
}

var Bits32 bool = false
var Filename string = ""
var SDFilename string = ""
var OpticalFilename string = ""
var DriveNumber int = 0
var BootDrive int = 0
