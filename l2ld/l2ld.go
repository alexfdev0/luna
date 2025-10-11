package main

import (
	"fmt"
	"os"
	"bytes"	
)

type binding struct {
	Name string
	Location []byte
}

type unresolvedBinding struct {
	Name string
	BufferLoc int
	Solved bool
	File string
}

var Buffer []byte

var section string = "text"

var bindings = []binding {}
var unresolvedBindings = []unresolvedBinding {}

var FillSize int = 0
var Org int = 0
var Entry bool = true

var errors = []string {
	"no object files specified",
	"file cannot be open()ed, errno=2",
	"multiple definitions of",
	"Undefined symbol for architecture luna-l2:",
	"File size exceeds padding directive:",
}
func error(errno int, args string) {
	fmt.Fprintln(os.Stderr, "l2ld: " + errors[errno] + " " + args)
	os.Exit(1)
}

func write(content byte) {
	Buffer = append(Buffer, content)	
}

func checkBinding(name string) ([]byte, bool) {
	for i := range bindings {
		if bindings[i].Name == name {
			return bindings[i].Location, true
		}
	}
	return nil, false
}

func Filter(data []byte, filename string) {	
	for i := 0; i < len(data); i++ {
		if bytes.HasPrefix(data[i:], []byte("LD16_")) || bytes.HasPrefix(data[i:], []byte("LD32_")) {
			// Add 32 bit defs later
			j := i + 5
			for j < len(data) && data[j] != 0x00 {
				j++
			}
			name := string(data[i + 5:j])
			j++
			
			org := 2
			if Org != 0 {
				org = Org
			}
			H := byte((org + len(Buffer)) >> 8)
			L := byte((org + len(Buffer)) & 0xFF)
			bindings = append(bindings, binding{Name: name, Location: []byte{H, L}})

			for i, ub := range unresolvedBindings {
				if ub.Name == name {
					unresolvedBindings[i].Solved = true
					Buffer[ub.BufferLoc] = H
					Buffer[ub.BufferLoc + 1] = L
				}
			}
			i = j - 1
		} else if bytes.HasPrefix(data[i:], []byte("LR_")) {
			j := i + 3
			for j < len(data) && data[j] != 0x00 {
				j++
			}
			name := string(data[i + 3:j])
			j++
			location, found := checkBinding(name)
			if found == true {
				for _, b := range location {
					write(b)
				}
			} else {
				unresolvedBindings = append(unresolvedBindings, unresolvedBinding{Name: name, BufferLoc: len(Buffer), Solved: false, File: filename})
				write(0x00)
				write(0x00)
			}
			i = j - 1
		} else if bytes.HasPrefix(data[i:], []byte("LP_")) {
			i += 3
			H := data[i]
			L := data[i + 1]	
			FillSize = int(H) << 8 | int(L) 	
			i++
		} else if bytes.HasPrefix(data[i:], []byte("LO_")) {
			i += 3
			H := data[i]
			L := data[i + 1]	
			Org = int(H) << 8 | int(L) 	
			i++
		} else if bytes.HasPrefix(data[i:], []byte("L_NOENTRY")) {
			i += 8
			Entry = false	
		} else {
			write(data[i])
		}
	}	
}

func main() {
	if len(os.Args) < 2 {
		error(0, "")
	}

	var input_files []string
	var output_filename string = ""

	for i := 1; i < len(os.Args); i++ {
		arg := os.Args[i]

		switch arg {
		case "-v":
			fmt.Println("@(#)PROGRAM:l2ld PROJECT:l2ld-2.0")
			fmt.Println("BUILD 16:02 Oct 6 2025")
			fmt.Println("configured to support archs: luna-l2")	
			os.Exit(0)
		case "-o":
			output_filename = os.Args[i + 1]
			i++
		default:
			input_files = append(input_files, arg)
		}
	}

	if len(input_files) < 1 {
		error(0, "")
	}
	if output_filename == "" {
		output_filename = "a.bin"
	}

	for _, file := range input_files {
		data, err := os.ReadFile(file)
		if err != nil {
			error(1, "path=" + file)
		}
		Filter(data, file)		
	}

	var buffer = []byte{}	
	if Entry == true {
		startloc, found := checkBinding("_start")
		if found == false {
			error(3, "\n  \"_start\", referenced from\n    <initial-undefines>")	
		}	
		buffer = append(buffer, startloc...)
	}
	buffer = append(buffer, Buffer...)	

	for _, ub := range unresolvedBindings {
		if ub.Solved == false {
			error(3, "\n  \"" + ub.Name + "\", referenced from\n    " + ub.File)
		}
	}
	if FillSize > 0 {	
		if len(buffer) > FillSize {
			error(4, "\n  directive: " + fmt.Sprintf("%d", FillSize) + ", actual: " + fmt.Sprintf("%d", len(buffer)))
		} else if len(buffer) < FillSize {
			for len(buffer) < FillSize {
				buffer = append(buffer, 0x00)
			}	
		}
	}
	os.WriteFile(output_filename, []byte(buffer), 0644)
}
