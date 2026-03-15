package error

import (
	"fmt"
	"os"
	"lcc1/lexer"
	"strings"
	"regexp"
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
	"unnecessary arguments passed to _start function", // 10
	"unknown attribute",
	"cannot combine with previous",
	"duplicate",
	"invalid type qualifier",
	"invalid preprocessing directive",
	"no such file or directory",
	"unknown pragma directive",
	"expected ';' after expression",
	"call to undeclared function",
	"too few arguments to function call,", // 20
	"too many arguments to function call,",
	"",
	"return type of",
	"change return type to",
	"type specifier missing, defaults to 'int'; ISO C99 and later do not support implicit int",
	"indirection requires pointer operand",
	"cannot take the address of an rvalue of type",
	"duplicate",
	"unterminated conditional directive",
	"array index", // 30
	"lvalue required as unary",
	"missing terminating",
	"cannot assign to variable",
	"attribute",
	"initializer element is not a compile-time constant",
	"implicit conversion from",
	"taking the address of a function argument is not supported",
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

func Stargaze(Tokens *[]lexer.Token, where int, errno int, kind int) {
	// Kinds:
	// 1: error
	// 2: warning
	// 3: note

	line := (*Tokens)[where].Line

	OGTVAL := (*Tokens)[where].Value

	if kind != 3 {
		(*Tokens)[where].Value = "\033[1;31m" + (*Tokens)[where].Value + "\033[0m"
	}

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
	text = strings.ReplaceAll(text, "* ", "*")
	text = strings.ReplaceAll(text, "\033[34mchar\033[0m *", "\033[34mchar\033[0m*")
	text = strings.ReplaceAll(text, "\033[34mint\033[0m *", "\033[34mint\033[0m*")
	text = strings.ReplaceAll(text, " [ ", "[")
	text = strings.ReplaceAll(text, " ] ", "] ")


	fmt.Printf("    %d | %s\n", line, text)
	if errno == 18 {
		var ansiRe = regexp.MustCompile(`\033\[[0-9;]*m`)
		visibleLen := len(ansiRe.ReplaceAllString(text, ""))
		fmt.Printf("     " + strings.Repeat(" ", len(string(line))) + "| " + strings.Repeat(" ", visibleLen) + "\033[1;32m^\033[0m\n")
		fmt.Printf("      | " + strings.Repeat(" ", visibleLen) + "\033[1;32m;\033[0m\n")
	} 

	if kind != 3 {
		(*Tokens)[where].Value = OGTVAL
	}
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

	addtl := " "
	if errno == 22 {
		addtl = ""
	} 
	fmt.Fprintln(os.Stderr, "\033[1;39m" + label + " \033[1;31merror: \033[1;39m" + errors[errno] + addtl + args + "\033[0m")
	Stargaze(tokens, find(token, tokens), errno, 1)
	Errors = Errors + 1	
	// os.Exit(1)
}

func ErrorNoGaze(errno int, args string, token lexer.Token) {
	label := "lcc:"
	if token.Line != 0 {
		label = token.File + ":" + fmt.Sprintf("%d", token.Line) + ":"
	}
	addtl := " "
	if errno == 22 {
		addtl = ""
	}
	fmt.Fprintln(os.Stderr, "\033[1;39m" + label + " \033[1;31merror: \033[1;39m" + errors[errno] + addtl + args + "\033[0m")
	Errors = Errors + 1
	// os.Exit(1)
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
	addtl := " "
	if errno == 22 {
		addtl = ""
	}
	fmt.Println("\033[1;39m" + label + " \033[1;35mwarning: \033[1;39m" + errors[errno] + addtl + args + "\033[0m")
	Stargaze(tokens, find(token, tokens), errno, 2)
	Warnings = Warnings + 1
}
func Note(errno int, args string, token lexer.Token, tokens *[]lexer.Token) {	
	label := "lcc:"
	if token.Line != 0 {
		label = token.File + ":" + fmt.Sprintf("%d", token.Line) + ":"
	}
	addtl := " "
	if errno == 22 {
		addtl = ""
	}
	fmt.Println("\033[1;39m" + label + " \033[1;36mnote: \033[1;39m" + errors[errno] + addtl + args + "\033[0m")
	Stargaze(tokens, find(token, tokens), errno, 3)
}

func Summary() {
	Warnings += lexer.Warnings
	if Errors < 1 && Warnings < 1 {
		return
	}

	str := ""	
	if Warnings > 0 {
		str = str + fmt.Sprintf("%d warning", Warnings)
	}
	if Warnings > 1 {
		str = str + "s"
	}
	if Warnings > 0 && Errors > 0 {
		str = str + " and "
	}
	if Errors > 0 {
		str = str + fmt.Sprintf("%d error", Errors)
	}	
	if Errors > 1 {
		str = str + "s"
	}	
	str = str + " generated."
	fmt.Println(str)
}
