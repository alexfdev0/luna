package main

import (	
	"lcc1/lexer"
	"lcc1/parser"
	"lcc1/error"	
	"os"
	"fmt"
	"strings"
	"runtime"
	"os/exec"
	"path/filepath"
)

func execute(command string) bool {
	shell := "sh"
	flag := "-c"
	if runtime.GOOS == "windows" {
		shell = "cmd"
		flag = "/C"
	}

	cmd := exec.Command(shell, flag, command)
	output, err := cmd.CombinedOutput()
	fmt.Printf(string(output))

	if err != nil {
		return false	
	}
	return true
}

func splitFile(path string) (name string, ext string) {
	ext = filepath.Ext(path)
	name = filepath.Base(path)	
	if ext != "" {
		name = name[:len(name)-len(ext)]
	}
	return
}

func main() {
	if len(os.Args) < 2 {
		error.ErrorNoGaze(0, "", lexer.Token{Line: 0})
		os.Exit(1)
	}

	var input_files = []string {}
	var output_file string = ""
	var noassemble bool = false
	var nolink bool = false

	for i := 1; i < len(os.Args); i++ {
		arg := os.Args[i]
		switch arg {
		case "-o":
			output_file = os.Args[i + 1]
			i++
		case "-v":
			fmt.Println("Luna Compiler Collection version 3.0")
			fmt.Println("Target: luna-l2")
			os.Exit(0)
		case "-S":
			noassemble = true
		case "-c":
			nolink = true
		case "-Werror":
			error.Upgrade = true
		default:
			input_files = append(input_files, arg)
		}
	}

	if output_file == "" {
		if nolink == true {
			output_file = "a.o"
		} else {
			output_file = "a.bin"
		}
	}

	assembly_files := []string {}
	for _, file := range input_files {
		data, err := os.ReadFile(file)
		if err != nil {
			os.Exit(1)
		}
		code := lexer.Preprocessor(string(data), file, false)
		tokens := lexer.Lex(code, file)
		parser.Parse(tokens, 1)
		name, _ := splitFile(file)
		assembly_files = append(assembly_files, name + ".s")
		os.WriteFile(name + ".s", []byte(".text\n" + parser.Code1 + "\n" + parser.Code2), 0644)
	}	

	if error.Errors > 0 {
		os.Exit(1)
	}

	if noassemble == true {
		os.Exit(0)
	}

	if nolink == true {
		execute("lcc -c " + strings.Join(assembly_files, " "))
	} else {
		execute("lcc " + strings.Join(assembly_files, " ") + " -o " + output_file)
	}

	for _, file := range assembly_files {
		if runtime.GOOS != "windows" {
			execute("rm -f " + file)
		} else {
			execute("del /f " + file)
		}
	}
}
