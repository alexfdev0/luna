package error

import (
	"fmt"
	"os"
	"lcc1/lexer"
	"strings"
)

var errors = []string {
	"no input files",
	"unexpected token",
	"expected",
	"redefinition of",
	"use of undeclared identifier",
	"incompatible type conversion",
	"could not evaluate mathematical expression",
	"variable has incomplete type",
	"comparison between pointer and non-pointer",
	"too many arguments passed to function (max 6)",
	"unnecessary arguments passed to _start function",
	"unknown attribute",
	"cannot combine with previous",
	"duplicate",
	"invalid type qualifier",
	"invalid preprocessing directive",
	"no such file or directory",
	"unknown pragma directive",
	"expected ';' after expression",
	"call to undeclared function",
	"too few arguments to function call,",
	"too many arguments to function call,",
	"",
}

var Warnings int = 0
var Errors int = 0
var Upgrade bool

func Clamp(num int, mini int, maxi int) int {
	if num < mini {
		return mini
	}
	if num > maxi {
		return maxi
	}
	return num
}

func Stargaze(Tokens *[]lexer.Token, where int) {
	line := (*Tokens)[where].Line	

	start := where
	for start > 0 && (*Tokens)[start - 1].Line == line {
		start--
	}
	end := where
	for end < len(*Tokens) - 1 && (*Tokens)[end + 1].Line == line {
		end++
	}

	words := make([]string, 0, end - start + 1)
	for j := start; j <= end; j++ {	
		if (*Tokens)[j].Type != lexer.TokType && (*Tokens)[j].Type != lexer.TokQualifier {
			words = append(words, (*Tokens)[j].Value)
		} else {
			words = append(words, "\033[34m" + (*Tokens)[j].Value + "\033[0m")
		}	
	}

	text := strings.Join(words, " ")	
	text = strings.ReplaceAll(text, " ( ", "(")
	text = strings.ReplaceAll(text, "( ", "(")
	text = strings.ReplaceAll(text, " )", ")")
	text = strings.ReplaceAll(text, " ;", ";")
	text = strings.ReplaceAll(text, " ,", ",")
	text = strings.ReplaceAll(text, "# ", "#")
	text = strings.ReplaceAll(text, "\033[34mchar\033[0m *", "\033[34mchar\033[0m*")
	text = strings.ReplaceAll(text, "\033[34mint\033[0m *", "\033[34mint\033[0m*")

	fmt.Printf("    %d | %s\n", line, text)
}

func find(token lexer.Token, tokens *[]lexer.Token) int {
	for i, t := range (*tokens) {
		if t == token {
			return i
		}
	}
	return 0
}

func Error(errno int, args string, token lexer.Token, tokens *[]lexer.Token) {
	label := "lcc:"
	if token.Line != 0 {
		label = token.File + ":" + fmt.Sprintf("%d", token.Line) + ":"
	}
	fmt.Fprintln(os.Stderr, "\033[1;39m" + label + " \033[1;31merror: \033[1;39m" + errors[errno] + " " + args + "\033[0m")
	Stargaze(tokens, find(token, tokens))
	Errors = Errors + 1
	os.Exit(1)
}

func ErrorNoGaze(errno int, args string, token lexer.Token) {
	label := "lcc:"
	if token.Line != 0 {
		label = token.File + ":" + fmt.Sprintf("%d", token.Line) + ":"
	}
	fmt.Fprintln(os.Stderr, "\033[1;39m" + label + " \033[1;31merror: \033[1;39m" + errors[errno] + " " + args + "\033[0m")
	Errors = Errors + 1
	os.Exit(1)
}

func Warning(errno int, args string, token lexer.Token, tokens *[]lexer.Token) {
	if Upgrade == true {
		Error(errno, args, token, tokens)
		return
	}
	label := "lcc:"
	if token.Line != 0 {
		label = token.File + ":" + fmt.Sprintf("%d", token.Line) + ":"
	}
	fmt.Println("\033[1;39m" + label + " \033[1;35mwarning: \033[1;39m" + errors[errno] + " " + args + "\033[0m")
	Stargaze(tokens, find(token, tokens))
	Warnings = Warnings + 1
}
func Note(errno int, args string, token lexer.Token, tokens *[]lexer.Token) {	
	label := "lcc:"
	if token.Line != 0 {
		label = token.File + ":" + fmt.Sprintf("%d", token.Line) + ":"
	}
	fmt.Println("\033[1;39m" + label + " \033[1;36mnote: \033[1;39m" + errors[errno] + " " + args + "\033[0m")
	Stargaze(tokens, find(token, tokens))
}
