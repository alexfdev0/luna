package main

import (
	"fmt"
	"os"
	"bytes"
	"runtime"
	"strings"
)

type binding struct {
	Name string
	Location []byte
	File string
	Global bool
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
var Globals = []string {}

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

func checkBinding(name string) binding {
	for i := range bindings {
		if bindings[i].Name == name {
			return bindings[i] 
		}
	}
	return binding{Name: "nil"}
}

func CheckGlobal(name string) bool {
	if name == "_start" {
		return true
	}
	for _, g := range Globals {
		if g == name {
			return true
		}
	}
	return false
}

var GBits32 bool = false

func Filter(data []byte, filename string) {	
	for i := 0; i < len(data); i++ {
		if bytes.HasPrefix(data[i:], []byte("LD16_")) || bytes.HasPrefix(data[i:], []byte("LD32_")) {
			var Bits32 bool = false
			if bytes.HasPrefix(data[i:], []byte("LD32_")) {
				Bits32 = true
			}

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
		
			if Bits32 == false {
				H := byte((org + len(Buffer)) >> 8)
				L := byte((org + len(Buffer)) & 0xFF)	
				bindings = append(bindings, binding{Name: name, Location: []byte{H, L}, File: filename, Global: CheckGlobal(name)})

				for i, ub := range unresolvedBindings {
					if ub.Name == name {
						if CheckGlobal(name) == true || (CheckGlobal(name) == false && ub.File == filename) {
							unresolvedBindings[i].Solved = true
							Buffer[ub.BufferLoc] = H
							Buffer[ub.BufferLoc + 1] = L
						}
					}
				}
			} else {
				HH := byte((org + len(Buffer)) >> 24)
				HL := byte((org + len(Buffer)) >> 16)
				LH := byte((org + len(Buffer)) >> 8)
				LL := byte((org + len(Buffer)) & 0xFF)
				bindings = append(bindings, binding{Name: name, Location: []byte{HH, HL, LH, LL}, File: filename, Global: CheckGlobal(name)})

				for i, ub := range unresolvedBindings {
					if ub.Name == name {
						if CheckGlobal(name) == true || (CheckGlobal(name) == false && ub.File == filename) {
							unresolvedBindings[i].Solved = true
							Buffer[ub.BufferLoc] = HH
							Buffer[ub.BufferLoc + 1] = HL
							Buffer[ub.BufferLoc + 2] = LH
							Buffer[ub.BufferLoc + 3] = LL
						}	
					}
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
			binding := checkBinding(name)	
			OK := false
			if binding.Global == false {	
				if filename == binding.File {	
					OK = true
				}
			} else {
				OK = true
			}

			if binding.Name != "nil" && OK == true {
				for _, b := range binding.Location {
					write(b)
				}
			} else {
				unresolvedBindings = append(unresolvedBindings, unresolvedBinding{Name: name, BufferLoc: len(Buffer), Solved: false, File: filename})
				if GBits32 == false {
					write(0x00)
					write(0x00)
				} else {
					write(0x00)
					write(0x00)
					write(0x00)
					write(0x00)
				}	
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
		} else if bytes.HasPrefix(data[i:], []byte("L_16BIT")) || bytes.HasPrefix(data[i:], []byte("L_32BIT")) {	
			if bytes.HasPrefix(data[i:], []byte("L_32BIT")) {
				GBits32 = true
			} else {
				GBits32 = false
			}
			i += 6
		} else if bytes.HasPrefix(data[i:], []byte("L_GLOBL_")) {
			j := i + 8
			for j < len(data) && data[j] != 0x00 {
				j++
			}
			name := string(data[i + 8:j])
			j++
			Globals = append(Globals, name)
			i = j - 1
		} else {
			write(data[i])
		}
	}	
}

var libs = make(map[string]string)
func ParseLibs() {
	file := ""
	if runtime.GOOS == "windows" {
		file = "C:\\luna\\l2ld\\libs.conf"
	} else {
		file = "/usr/local/lib/l2ld/libs.conf"
	}
	data, err := os.ReadFile(file)
	if err != nil {
		fmt.Println("l2ld: could not read libs.conf file")
	}
	data_string := string(data)
	data_words := strings.Fields(data_string)

	for i := 0; i < len(data_words); i++ {
		word := data_words[i]
		if i + 1 >= len(data_words) {
			fmt.Println("l2ld: error in libs.conf near '" + word + "'")
			break
		}
		nextword := data_words[i + 1]
		libs[word] = nextword
		i++
	}
}

func main() {
	if len(os.Args) < 2 {
		error(0, "")
	}

	var input_files []string
	var output_filename string = ""
	var auto bool = false

	for i := 1; i < len(os.Args); i++ {
		arg := os.Args[i]

		switch arg {
		case "-v":
			fmt.Println("@(#)PROGRAM:l2ld PROJECT:l2ld-2.1")
			fmt.Println("BUILD 11:17 Oct 20 2025")
			fmt.Println("configured to support archs: luna-l2")	
			os.Exit(0)
		case "-o":
			output_filename = os.Args[i + 1]
			i++
		case "-a":
			auto = true
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

	ParseLibs()
	for _, file := range input_files {
		data, err := os.ReadFile(file)
		if err != nil {
			error(1, "path=" + file)
		}	
		Filter(data, file)		
	}

	done:

	var buffer = []byte{}	
	if Entry == true {
		binding := checkBinding("_start")
		if binding.Name == "nil" {
			error(3, "\n  \"_start\", referenced from\n    <initial-undefines>")	
		}	
		buffer = append(buffer, binding.Location...)
	}
	buffer = append(buffer, Buffer...)	

	for _, ub := range unresolvedBindings {
		if ub.Solved == false {
			if libs[ub.Name] != "" && auto == true {
				file := libs[ub.Name]
				data, err := os.ReadFile(file)
				if err != nil {
					error(1, "path=" + file + " (in libs.conf)")
				}
				Filter(data, file)
				goto done
			}
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
