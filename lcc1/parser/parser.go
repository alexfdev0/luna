package parser

import (
	"lcc1/lexer"
	"lcc1/error"
	"strings"
	"fmt"	
	"strconv"
	"os"
	// "runtime/debug"
)

var level int = 0

var Code1 string = ""
var Code2 string = ""


var IDCounter = 1

const (
	NUMBER8 int = iota
	NUMBER16
	NUMBER32
	STRING
	POINT
	NULL
	ARRAY
)

type Variable_Static struct {
	Name string
	Type int
	Value any
	Pointer bool
	Real string
	Scope int
	Const bool
	Extern bool
	ArgNum int
}

type FunctionDecl struct {
	Name string
	Token lexer.Token
	Set []lexer.Token
}

type Scope struct {
	ID int
	Parent int
}

type Define struct {
	From string
	To string
}

var Variables = []Variable_Static {}
var FunctionDecls = []FunctionDecl {}
var Defines = []Define {}

var Scopes = []Scope {
	Scope{ID: 1, Parent: -1},
}

func Write(text string, spaced bool) {
	if spaced == false {
		Code2 = Code2 + text + "\n"
	} else {
		Code2 = Code2 + "    " + text + "\n"
	}
}

func WritePre(text string, spaced bool) {
	if spaced == false {
		Code1 = Code1 + text + "\n"
	} else {
		Code1 = Code1 + "    " + text + "\n"
	}
}

func PreWrite(text string, spaced bool) {
	if spaced == false {
		Code1 = text + "\n" + Code1
	} else {
		Code1 = "    " + text + "\n" + Code1
	}
}

func CheckNum(token lexer.Token) bool {
	if _, err := strconv.ParseInt(token.Value, 0, 64); err == nil {
		return true
	} else {
		return false
	}
}

func CreateStatic(variable Variable_Static) {
	WritePre(variable.Name + ":\n    .asciz \"" + variable.Value.(string) + "\"", false)	
}

func LookupParent(Scope int) int {
	if Scope == 1 {
		return -1
	}
	for _, s := range Scopes {
		if s.ID == Scope {
			return s.Parent
		}
	}
	return 1
}

func CreateScope(Parent int) int {
	Scopes = append(Scopes, Scope{ID: IDCounter, Parent: Parent})
	IDCounter++
	return IDCounter - 1
}

func ParseExpyL1(tokens []lexer.Token, i int, Scope int) int {
	for {
		if i >= len(tokens) {
			break
		}
		i = ParseExpy(tokens, i, Scope, "r4")
	}
	return i
}

func ParseExpy(tokens []lexer.Token, start int, Scope int, register string) int {
	i := start
	expect := func(toktype lexer.TokenType) string {
		var value string
		if i >= len(tokens) {
			if toktype != lexer.TokSemi {
				error.Error(1, "'<EOF>'", tokens[i - 1], &tokens)
			} else {
				error.Error(18, "", tokens[i - 1], &tokens)
			}
			return ""
		}
		if tokens[i].Type == toktype {
			value = tokens[i].Value
			i++
		} else {
			if toktype != lexer.TokSemi {
				error.Error(1, "'" + tokens[i].Value + "'", tokens[i], &tokens)
			} else {
				error.Error(18, "", tokens[i - 1], &tokens)
			}
		}
		return value
	}	
	peek := func(lookahead int) lexer.Token {
		if i + lookahead < len(tokens) && i + lookahead >= 0 {
			return tokens[i + lookahead]
		}
		return lexer.Token{Type: lexer.TokEOF, Value: ""}
	}


	IDENT_FUNC := func(label string) {
		expect(lexer.TokLParen)

		switch label {
		case "asm", "__asm__":
			if peek(0).Type == lexer.TokQualifier && peek(0).Value == "volatile" {
				expect(lexer.TokQualifier)
			}
			asmval := expect(lexer.TokIdent)
			expect(lexer.TokRParen)
			// expect(lexer.TokSemi)
			// TODO: make quotes check here
			Write(strings.ReplaceAll(asmval, "\"", ""), true)
		case "sizeof":
			val := 0
			_label := expect(lexer.TokIdent)
			expect(lexer.TokRParen)
			Variable := LookupVariable(_label, true, Scope, peek(-2), &tokens)
			switch Variable.Type {
			case NUMBER8:
				val = 1
			case STRING:
				if Variable.Pointer == true {
					switch lexer.Bits {
					case 16:
						val = 2
					case 32:
						val = 4
					}
				} else {
					val = 1
				}
			case NUMBER16:
				val = 2
			case NUMBER32:
				val = 4
			}
			Write("mov " + register + ", " + fmt.Sprintf("%d", val), true)
		default:
			Function_Variable := LookupVariable(label, false, Scope, peek(-2), &tokens)
			if Function_Variable.Name == "__ZERO" {
				error.Error(19, "'" + label + "'; ISO C99 and later do not support implicit function declarations", peek(-2), &tokens)	
			}
			// Parse arguments
			pushed := 0
			j := i
			exit := false
			var CurrentTokens []lexer.Token
			for j = i; j < len(tokens); j++ {
				if exit == true {
					break
				}
				switch tokens[j].Type {
				case lexer.TokComma:
					ParseExpy(CurrentTokens, 0, Scope, "r7")
					Write("push r7", true)
					CurrentTokens = []lexer.Token{}
					pushed++
				case lexer.TokRParen:
						// fmt.Println("breaking")
						exit = true
						break
				default:
					CurrentTokens = append(CurrentTokens, tokens[j])	
				}
			}

			// Push last args 
			if len(CurrentTokens) > 0 {
				ParseExpy(CurrentTokens, 0, Scope, "r7")
				Write("push r7", true)
				pushed++
			}
			Write("call " + label, true)
			Write("mov " + register + ", e6", true)

			if pushed < Function_Variable.ArgNum {
				t, s := FuncDeclLookup(label)	
				error.Error(20, "expected " + fmt.Sprintf("%d", Function_Variable.ArgNum) + ", have " + fmt.Sprintf("%d", pushed), peek(0), &tokens)
				error.Note(22, "'" + label + "' declared here", t, s)
			} else if pushed > Function_Variable.ArgNum {
				t, s := FuncDeclLookup(label)	
				error.Error(21, "expected " + fmt.Sprintf("%d", Function_Variable.ArgNum) + ", have " + fmt.Sprintf("%d", pushed), peek(0), &tokens)
				error.Note(22, "'" + label + "' declared here", t, s)
			}
			i = j + 1	
		}
	}
	IDENT_STRING := func(label string) {
		if label[len(label) - 1] != '"' {
			// TODO: fix "literal not terminated and have us handle it"
			error.Error(32, "\" character", peek(-1), &tokens)
		}
		_label := fmt.Sprintf("var_%d", IDCounter)
		WritePre(_label + ":", false)
		WritePre(".asciz \"" + strings.ReplaceAll(label, "\"", "") + "\"", true)
		Write("mov " + register + ", " + _label, true)
		IDCounter++
	}

	deref := 0
	EQU_VT := NULL
	EQU_VAR := Variable_Static{}
	EXPY_TOP:
	switch peek(0).Type {
	default:
		expect(lexer.TokEOF)
	case lexer.TokSemi:
		expect(lexer.TokSemi)
	case lexer.TokStar:
		expect(lexer.TokStar)
		deref++
		goto EXPY_TOP
	case lexer.TokAmpersand:
		expect(lexer.TokAmpersand)
		deref--
		if deref <= -2 {
			error.Error(27, "'int'", peek(-1), &tokens)
		}
		goto EXPY_TOP
	case lexer.TokGoto:
		// TODO: add goto var checks
		expect(lexer.TokGoto)
		label := expect(lexer.TokIdent)
		expect(lexer.TokSemi)
		Write("jmp " + label, true)
	case lexer.TokIf:
		// TODO: implement quick ifs 
		expect(lexer.TokIf)
		expect(lexer.TokLParen)

		first_tokens := []lexer.Token {}
		second_tokens := []lexer.Token {}

		j := i
		exit := false
		for j = i; j < len(tokens); j++ {
			if exit == true {
				break
			}
			switch tokens[j].Type {
			case lexer.TokEquality, lexer.TokInequality:
				exit = true	
			default:
				first_tokens = append(first_tokens, tokens[j])
			}	
		}
		i = j - 1
		
		op := ""
		op_reverse := ""
		IfScope := CreateScope(Scope)
		ElseScope := CreateScope(Scope)
		ParseExpy(first_tokens, 0, Scope, "r9")

		switch peek(0).Type {
		case lexer.TokEquality:
			op = "jnz"
			op_reverse = "jz"
			expect(lexer.TokEquality)
		case lexer.TokInequality:
			op = "jz"
			op_reverse = "jnz"
			expect(lexer.TokInequality)
		default:
			// Error out
			expect(lexer.TokEOF)
		}

		j = i
		depth := 1
		exit = false
		for j = i; j < len(tokens); j++ {
			if exit == true {
				break
			}
			switch tokens[j].Type {
			case lexer.TokLParen:
				depth++
			case lexer.TokRParen:
				depth--
				if depth == 0 {
					exit = true
				}
			default:
				second_tokens = append(second_tokens, tokens[j])
			}	
		}
		i = j - 1

		ParseExpy(second_tokens, 0, Scope, "r10")

		expect(lexer.TokRParen)

		if_label := fmt.Sprintf("if_stmt_%d", IDCounter)
		IDCounter++
		else_label := fmt.Sprintf("else_stmt_%d", IDCounter)
		IDCounter++
		after_label := fmt.Sprintf("after_stmt_%d", IDCounter)
		IDCounter++	

		if_tokens := []lexer.Token {}
		else_tokens := []lexer.Token {}
		
		expect(lexer.TokLCurly)	
		j = i
		depth = 1
		exit = false
		for j = i; j < len(tokens); j++ {
			if exit == true {
				break
			}
			switch tokens[j].Type {
			case lexer.TokLCurly:
				depth++
				if_tokens = append(if_tokens, tokens[j])
			case lexer.TokRCurly:
				depth--
				if depth == 0 {
					exit = true
				} else {
					if_tokens = append(if_tokens, tokens[j])
				}
			default:
				if_tokens = append(if_tokens, tokens[j])
			}
		}
		i = j - 1
		expect(lexer.TokRCurly)
		if peek(0).Type != lexer.TokElse {
			// Write everything
			Write("cmp r11, r9, r10", true)
			Write(op + " r11, " + if_label, true)
			Write(op_reverse + " r11, " + after_label, true)
			Write(if_label + ":", false)
			ParseExpyL1(if_tokens, 0, IfScope)
			Write("jmp " + after_label, true)
			Write(after_label + ":", false)
			goto DONE
		}
		
		expect(lexer.TokElse)

		expect(lexer.TokLCurly)	
		j = i
		depth = 1
		exit = false
		for j = i; j < len(tokens); j++ {
			if exit == true {
				break
			}
			switch tokens[j].Type {
			case lexer.TokLCurly:
				depth++
				else_tokens = append(else_tokens, tokens[j])
			case lexer.TokRCurly:
				depth--
				if depth == 0 {
					exit = true
				} else {
					else_tokens = append(else_tokens, tokens[j])
				}
			default:
				else_tokens = append(else_tokens, tokens[j])
			}
		}
		i = j - 1
		expect(lexer.TokRCurly)

		// Write everything
		Write("cmp r11, r9, r10", true)
		Write(op + " r11, " + if_label, true)
		Write(op_reverse + " r11, " + else_label, true)	
		Write(if_label + ":", false)
		ParseExpyL1(if_tokens, 0, IfScope)
		Write("jmp " + after_label, true)
		Write(else_label + ":", false)
		ParseExpyL1(else_tokens, 0, ElseScope)
		Write(after_label + ":", false)
	case lexer.TokReturn:
		expect(lexer.TokReturn)

		switch peek(0).Type {
		case lexer.TokSemi:
			
		default:
			i = ParseExpy(tokens, i, Scope, "e6")
		}	
		expect(lexer.TokSemi)
		Write("pop e11", true)
		Write("ret", true)
	case lexer.TokIdent:
		label := expect(lexer.TokIdent)
		var variable Variable_Static

		if label[0] == '"' {
			IDENT_STRING(label)
			goto CONTINUE
		}

		if peek(0).Type == lexer.TokLParen {
			IDENT_FUNC(label)
			goto CONTINUE
		}

		if peek(0).Type == lexer.TokColon {
			Write(label + ":", false)
			expect(lexer.TokColon)
			goto DONE
		}

		variable = LookupVariable(label, true, Scope, peek(-1), &tokens)
		EQU_VT = variable.Type
		EQU_VAR = variable
		Write("mov r1, " + variable.Real, true)
		switch variable.Type {
		case NUMBER8:
			Write("lod r1, r2", true)
		case STRING:
			if variable.Pointer == false {
				Write("lod r1, r2", true)
			} else {
				Write("mov r2, r1", true)
			}
		case NUMBER16, NUMBER32:
			if variable.Pointer == false {
				Write("lodf r1, r2", true)
			} else {
				Write("mov r2, r1", true)
			}
		case NULL:
			if variable.Pointer == true {
				Write("mov r2, r1", true)
			}
		}

		for deref > 0 {
			deref--
			switch variable.Type {
			case NUMBER8, STRING, NULL:
				Write("lod r2, r2", true)
			case NUMBER16, NUMBER32:
				Write("lodf r2, r2", true)
			}
		}
		Write("mov " + register + ", r2", true)			
	case lexer.TokNumber:
		// Parse expressions
		// Load it up into r4
		exit := false
		num1 := expect(lexer.TokNumber)
		var num2 string = "" 
		var op string = ""
		OP_TRY:

		switch peek(0).Type {
		case lexer.TokPlus, lexer.TokMinus, lexer.TokStar, lexer.TokSlash:
			exit = false
		default:
			exit = true
		}
		if exit == true {
			goto CHECK_FOR_PTR
		}

		op = expect(peek(0).Type)
		num2 = expect(lexer.TokNumber)
		switch op {
		case "+":
			Write("mov r5, " + num1, true)
			Write("mov r6, " + num2, true)
			Write("add " + register + ", r5, r6", true)
		case "-":
			Write("mov r5, " + num1, true)
			Write("mov r6, " + num2, true)
			Write("sub " + register + ", r5, r6", true)
		case "*":
			Write("mov r5, " + num1, true)
			Write("mov r6, " + num2, true)
			Write("mul " + register + ", r5, r6", true)
		case "/":
			Write("mov r5, " + num1, true)
			Write("mov r6, " + num2, true)
			Write("div " + register + ", r5, r6", true)
		}
		switch peek(0).Type {
		case lexer.TokPlus, lexer.TokMinus, lexer.TokStar, lexer.TokSlash:
			num1 = "r5"
			goto OP_TRY
		}
		goto CONTINUE

		CHECK_FOR_PTR:
		if deref < 0 {
			error.Error(31, "'&' operand", peek(-1), &tokens)
		}
		derefed := false
		Write("mov r1, " + num1, true)
		for _ = deref; deref > 0; deref-- {
			derefed = true
			Write("lod r1, r2", true)
		}
		if derefed == false {
			Write("mov " + register + ", r1", true)
		} else {
			Write("mov " + register + ", r2", true)
		}
	}
	CONTINUE:
	
	if peek(0).Type != lexer.TokEqual {
		// fmt.Println("PEEK 0:", peek(0).Value)
		// fmt.Println("PEEK -1:", peek(-1).Value)	
		// expect(lexer.TokSemi)
		return i
	}

	expect(lexer.TokEqual)

	i = ParseExpy(tokens, i, Scope, "r5")

	if EQU_VAR.Const == true {
		error.Error(33, "'" + EQU_VAR.Name + "' with const-qualified type", peek(-1), &tokens)
		token, stream := FuncDeclLookup(EQU_VAR.Name)
		error.Note(22, "'" + EQU_VAR.Name + "' declared here", token, stream)
	}

	switch EQU_VT {
	case NUMBER8, STRING:
		Write("str " + register + ", r5", true)
	case NUMBER16, NUMBER32:
		Write("strf " + register + ", r5", true)
	}

	expect(lexer.TokSemi)

	DONE:

	return i
}

func LookupVariable(Name string, Enforce bool, Scope int, Token lexer.Token, Tokens *[]lexer.Token) Variable_Static {
	for {
		for _, variable := range Variables {
			if variable.Name == Name && variable.Scope == Scope {	
				return variable
			}
		}
		parent := LookupParent(Scope)
		if parent == -1 {	
			break
		}	
		Scope = parent
	}
	
	if Enforce == true {
		error.Error(4, "'" + Name + "'", Token, Tokens)
	}
	return Variable_Static{Name: "__ZERO", Type: NULL, Value: 0}
}

func StringParse(tokens []lexer.Token, start int) (string, int) {
	// Start would be the first token
	var str string = ""
	var loc int = 0
	
	if strings.HasPrefix(tokens[start].Value, "\"") == false {
		error.Error(2, "\"", tokens[start], &tokens)		
	}
	if strings.HasSuffix(tokens[start].Value, "\"") {
		word := strings.Trim(tokens[start].Value, "\"")
		str = word
		loc = start
	} else {
		var strtokens = []string { tokens[start].Value }
		for k := start + 1; k < len(tokens); k++ {
			strtokens = append(strtokens, tokens[k].Value)
			if strings.HasSuffix(tokens[k].Value, "\"") {
				start = k
				break
			}
		}
		str = strings.Join(strtokens, " ")
		str = strings.TrimSuffix(str,  "\"")
		loc = start
	}
	
	return str, loc
}

func FuncDeclLookup(Name string) (lexer.Token, *[]lexer.Token) {
	for _, d := range FunctionDecls {
		if d.Name == Name {
			return d.Token, &d.Set
		}
	}
	fmt.Println("Compiler fault: func not found on lookup")
	os.Exit(1)
	return lexer.Token{Type: lexer.TokEOF}, &[]lexer.Token {}
}

var topLevelName string
var BitPref int = 16
func Parse(tokens []lexer.Token, Scope int) {
	i := 0
	expect := func(toktype lexer.TokenType) string {
		var value string
		if i >= len(tokens) {
			if toktype != lexer.TokSemi {
				error.Error(1, "'<EOF>'", tokens[i - 1], &tokens)
			} else {
				error.Error(18, "", tokens[i - 1], &tokens)
			}
		}
		if tokens[i].Type == toktype {
			value = tokens[i].Value
			i++
		} else {
			if toktype != lexer.TokSemi {
				error.Error(1, "'" + tokens[i].Value + "'", tokens[i], &tokens)
			} else {
				error.Error(18, "", tokens[i - 1], &tokens)
			}
		}
		return value
	}	
	peek := func(lookahead int) lexer.Token {
		if i + lookahead < len(tokens) {
			return tokens[i + lookahead]
		}
		return lexer.Token{Type: lexer.TokEOF, Value: ""}
	}

	PreWrite("jmp _init", false)
	Write("_init:", false)
	
	for {
		if i >= len(tokens) {
			break
		}	

		switch level {
		case 0:	
			var ptr bool = false
			_typetoken := tokens[i]

			long := false
			short := false
			shortshort := false
			unsigned := false
			constant := false
			extern := false
			static := false
			signed := false
			bits := BitPref	
			for {	
				if peek(0).Value == "asm" || peek(0).Value == "__asm__" {
					expect(lexer.TokIdent)
					if peek(0).Value == "volatile" {
						expect(lexer.TokQualifier)
					}
					expect(lexer.TokLParen)
					str, end := StringParse(tokens, i)
					i = end + 1
					WritePre(str, false)
					expect(lexer.TokRParen)
					expect(lexer.TokSemi)
					continue
				}
				if peek(0).Type == lexer.TokTypedef {
					expect(lexer.TokTypedef)
					// TODO: add typedef
				}
				if peek(0).Type == lexer.TokQualifier {	
					qual := expect(lexer.TokQualifier)	
					switch qual {
					case "short":
						if long == true {
							error.Error(12, "'long' declaration specifier", peek(-1), &tokens)
						}
						if short == true && shortshort == true {
							error.Warning(13, "'short' declaration specifier", peek(-1), &tokens)
						} else if short == true && shortshort == false {
							shortshort = true
							bits = 8
						} else {
							short = true
							bits = 16
						}
					case "long":
						if short == true {
							error.Error(12, "'short' declaration specifier", peek(-1), &tokens)
						}
						if long == true {
							error.Warning(13, "'long' declaration specifier", peek(-1), &tokens)
						}
						long = true
						bits = 32
					case "unsigned":
						if signed == true {
							error.Error(12, "'signed' declaration specifier", peek(-1), &tokens)
						}
						if unsigned == true {
							error.Error(28, "'unsigned'", peek(-1), &tokens)
						}
						unsigned = true
					case "signed":
						if unsigned == true {
							error.Error(12, "'unsigned' declaration specifier", peek(-1), &tokens)	
						}
						if signed == true {
							error.Error(28, "'signed'", peek(-1), &tokens)
						}
						signed = true
					case "const":
						if constant == true {
							error.Error(28, "'const'", peek(-1), &tokens)
						}
						constant = true
					case "extern":
						if extern == true {
							error.Error(28, "'extern'", peek(-1), &tokens)
						}
						extern = true
					case "static":
						if static == true {
							error.Error(28, "'static'", peek(-1), &tokens)
						}
						static = true
					}
				} else {
					break
				}
			}

			if peek(0).Type == lexer.TokIdent {
				error.Error(25, "", peek(0), &tokens)
			}

			_typetok := tokens[i]
			_type := expect(lexer.TokType)
			if peek(0).Type == lexer.TokStar {
				ptr = true
				i++
			}

			name := expect(lexer.TokIdent)
		
			var rtype int	
			switch _type {
			case "int":
				switch bits {
				case 8:
					rtype = NUMBER8
				case 16:
					rtype = NUMBER16
				case 32:
					rtype = NUMBER32
				default:
					rtype = NUMBER16
				}
			case "char":
				if long == true || short == true || unsigned == true {
					error.Error(14, "for type 'char'", peek(-2), &tokens)
				}
				rtype = STRING
			}	

			_variable := LookupVariable(name, false, Scope, tokens[i - 1], &tokens) 
			if _variable.Name != "__ZERO" && _variable.Scope == Scope {	
				error.Error(3, "'" + name + "'", tokens[i - 1], &tokens)
			}

			FunctionDecls = append(FunctionDecls, FunctionDecl{Name: name, Token: peek(-1), Set: tokens})

			if peek(0).Type == lexer.TokLParen {
				rns := false
				if name == "main" {
					rns = true
					name = "_start"
				}	
				expect(lexer.TokLParen)
				fscope := CreateScope(Scope)	
		
				register := 0
				nargs := 0
				switch peek(0).Type {
				case lexer.TokType:
					if name == "_start" {
						error.Warning(10, "", peek(0), &tokens)
					}
					register = 0
					expComma := false
					for j := i; j < len(tokens); j++ {
						if peek(0).Type == lexer.TokRParen {
							expect(lexer.TokRParen)
							break
						}
						if expComma == false {
							if register >= 6 {
								error.Error(9, "", peek(0), &tokens)	
							}
							ptr := false
							expect(lexer.TokType)
							if peek(0).Type == lexer.TokStar {
								ptr = true
								expect(lexer.TokStar)
							}	
							__name := expect(lexer.TokIdent)	
							Variables = append(Variables, Variable_Static{Name: __name, Type: rtype, Value: nil, Scope: fscope, Real: fmt.Sprintf("e%d", register), Pointer: ptr})
							register++
							nargs++
							expComma = true
						} else {
							expect(lexer.TokComma)
							expComma = false
						}	
					}	
				case lexer.TokRParen:
					expect(lexer.TokRParen)
				}

				Variables = append(Variables, Variable_Static{Name: name, Type: rtype, Value: nil, Scope: Scope, Real: "e6", Extern: extern, ArgNum: nargs})

				noreturn := false
				if peek(0).Value == "__attribute__" {
					var attrs []string
					expect(lexer.TokIdent)
					expect(lexer.TokLParen)
					expect(lexer.TokLParen)

					expComma := false
					for {
						if expComma == false {
							attr := expect(lexer.TokIdent)
							attrs = append(attrs, attr)
							expComma = true
							if peek(0).Type == lexer.TokRParen {
								break
							}
						} else {
							expect(lexer.TokComma)
							expComma = false
						}
					}	
					
					expect(lexer.TokRParen)
					expect(lexer.TokRParen)

					for _, attr := range attrs {
						switch attr {
						case "norename":
							if rns == true {
								name = "main"
							}
						case "noreturn":
							noreturn = true
						default:
							error.Warning(11, "'" + attr + "'", tokens[i - 3], &tokens)
						}
					}
				}

				if peek(0).Type == lexer.TokSemi {
					expect(lexer.TokSemi)
					continue
				}

				if name == "_start" && _type != "void" {
					error.Warning(23, "'_start' is not 'void'", _typetok, &tokens)
					error.Note(24, "'void'", _typetok, &tokens)
				}
				expect(lexer.TokLCurly)	

				var Children = []lexer.Token {}
				ending := -1

				depth := 1
				for j := i; j < len(tokens); j++ {
					if tokens[j].Type == lexer.TokRCurly {
						depth--
						if depth == 0 {
							ending = j
							break
						} else {
							Children = append(Children, tokens[j])
						}	
					} else if tokens[j].Type == lexer.TokLCurly {
						depth++
						Children = append(Children, tokens[j])
					} else {
						Children = append(Children, tokens[j])
					}
				}
				if ending == -1 {
					error.Error(2, "'}'", tokens[i], &tokens)
				} else {
					i = ending
				}
			
				expect(lexer.TokRCurly)
	
				Write(name + ":", false)

				/*
				if name == "_start" {
					PreWrite("jmp _start", false)
				}
				*/

				if len(Children) > 0 {
					level = 1
					if static == false {
						PreWrite(".global " + name, false)
					}
					topLevelName = name
					if name != "_start" && noreturn == false {
						Write("pop e11", true)	
					}
					if register > 0 {
						for r := nargs - 1; nargs > 0; nargs-- {
							Write("pop e" + fmt.Sprintf("%d", r), true)
							r--
						}
					}
					if name != "_start" && noreturn == false {
						Write("push e11", true)
					}
	
					ParseExpyL1(Children, 0, fscope)

					if name != "_start" && noreturn == false {
						Write("pop e11", true)
						Write("ret", true)
					}
					IDCounter++
					topLevelName = ""
					level = 0
				}	
			} else if peek(0).Type == lexer.TokEqual {
				expect(lexer.TokEqual)	
				switch _type {
				case "void":
					error.Error(7, "'void'", _typetoken, &tokens)
				case "int":
					end := ParseExpy(tokens, i, Scope, "r4")	
					var val any

					rn := "var_" + fmt.Sprintf("%d", IDCounter)
					IDCounter++
					WritePre(rn + ":", false)
					switch rtype {
					case NUMBER8, STRING:
						WritePre(".byte 0x00", true)
					case NUMBER16:
						WritePre(".word 0x0000", true)
					case NUMBER32:
						WritePre(".dword 0x00000000", true)
					}

					if ptr == true {	
						Variables = append(Variables, Variable_Static{Name: name, Type: rtype, Value: val, Pointer: true, Real: rn, Scope: Scope, Const: constant})
						// Move result to variable
						Write("mov r7, " + rn, true)
						Write("strf r7, r4", true)
					} else {
						Variables = append(Variables, Variable_Static{Name: name, Type: rtype, Value: val, Pointer: false, Real: rn, Scope: Scope, Const: constant})
						// Move result to variable
						Write("mov r7, " + rn, true)
						Write("strf r7, r4", true)
					}
					i = end
				case "char":
					str, end := StringParse(tokens, i)	
					if ptr == true {
						rn := "var_" + fmt.Sprintf("%d", IDCounter)
						IDCounter++
						Variables = append(Variables, Variable_Static{Name: name, Type: STRING, Value: str, Pointer: true, Real: rn, Scope: Scope, Const: constant})
						WritePre(rn + ":", false)
						WritePre(".asciz \"" + str + "\"", true)
					} else {	
						if len(str) > 1 {
							error.Error(5, "'char' with an expression of type 'char*'", tokens[i], &tokens)
						}
						rn := "var_" + fmt.Sprintf("%d", IDCounter)
						IDCounter++
						WritePre(rn + ":", false)
						WritePre(".byte " + fmt.Sprintf("0x%02x", str[0]), true)
						Variables = append(Variables, Variable_Static{Name: name, Type: STRING, Value: str, Pointer: false, Scope: Scope, Const: constant, Real: rn})
					}
					i = end + 1
				}
				expect(lexer.TokSemi)
			} else if peek(0).Type == lexer.TokSemi {
				expect(lexer.TokSemi)

				switch _type {
				case "int":
					if ptr == true {
						rn := "var_" + fmt.Sprintf("%d", IDCounter)
						WritePre(rn + ":", false)
						IDCounter++
						switch rtype {
						case NUMBER8, STRING:
							WritePre(".byte 0x00", true)
						case NUMBER16:
							WritePre(".word 0x0000", true)
						case NUMBER32:
							WritePre(".dword 0x00000000", true)
						}

						Variables = append(Variables, Variable_Static{Name: name, Type: rtype, Value: 0, Pointer: true, Real: rn, Scope: Scope, Const: constant})
					} else {
						rn := "var_" + fmt.Sprintf("%d", IDCounter)
						WritePre(rn + ":", false)
						IDCounter++
						switch rtype {
						case NUMBER8, STRING:
							WritePre(".byte 0x00", true)
						case NUMBER16:
							WritePre(".word 0x0000", true)
						case NUMBER32:
							WritePre(".dword 0x00000000", true)
						}

						Variables = append(Variables, Variable_Static{Name: name, Type: rtype, Value: 0, Pointer: false, Scope: Scope, Const: constant})	
					}
				case "char":
					if ptr == true {
						rn := "var_" + fmt.Sprintf("%d", IDCounter)
						IDCounter++
						Variables = append(Variables, Variable_Static{Name: name, Type: STRING, Value: "", Pointer: true, Real: rn, Scope: Scope, Const: constant})
						WritePre(rn + ":", false)
						WritePre(".asciz \"\"", true)
					} else {	
						Variables = append(Variables, Variable_Static{Name: name, Type: STRING, Value: "", Pointer: false, Scope: Scope, Const: constant})
					}
				case "void":
					if ptr == true {
						Variables = append(Variables, Variable_Static{Name: name, Type: NULL, Value: nil, Pointer: true, Real: name, Scope: Scope, Const: constant})
					}
				}
			} else if peek(0).Type == lexer.TokLBracket {
				expect(lexer.TokLBracket)
				// Length next
				length := expect(lexer.TokNumber)
				length_real, _ := strconv.ParseInt(length, 0, 64)
				expect(lexer.TokRBracket);

				rn := "var_" + fmt.Sprintf("%d", IDCounter)
				IDCounter++
				Variables = append(Variables, Variable_Static{Name: name, Type: rtype, Value: nil, Pointer: false, Real: rn, Scope: Scope, Const: constant, ArgNum: int(length_real) })
				Write(rn + ":", false)

				if rtype == NUMBER8 || rtype == STRING {
					length_real = length_real
				} else if rtype == NUMBER16 {
					length_real = length_real * 2
				} else if rtype == NUMBER32 {
					length_real = length_real * 4
				}

				Write(".pad " + fmt.Sprintf("%d", length_real), true)

				expect(lexer.TokSemi)
			} else {
				error.Error(1, "'" + peek(0).Value + "'", _typetoken, &tokens)
			}
		case 1:	
			// Variable reassignment / function call
			_FUNC_CALL := func(name string, expect_semi bool) bool {	
				return false
			}

			var type_ lexer.TokenType = peek(0).Type
			switch type_ {
			case lexer.TokIdent, lexer.TokStar, lexer.TokAmpersand:
				// _name_token := tokens[i]

				// deref := false	

				if peek(2).Type != lexer.TokLParen {
					if peek(0).Type == lexer.TokStar {
						expect(lexer.TokStar)
						// deref = true
					}
				}

				_ntok := peek(0)
				name := expect(lexer.TokIdent)
				if name == "asm" || name == "__asm__" {
					if peek(0).Value == "volatile" {
						expect(lexer.TokQualifier)
					}
				}

				ASSIGNMENT_TOP:
				if peek(0).Type == lexer.TokLParen {
					if _FUNC_CALL(name, true) == true {
						goto ASSIGNMENT_TOP
					}
				} else if peek(0).Type == lexer.TokEqual || peek(0).Type == lexer.TokExclamation || peek(0).Type == lexer.TokLAngle || peek(0).Type == lexer.TokRAngle || peek(0).Type == lexer.TokLBracket {
					switch {
						case peek(1).Type == lexer.TokEqual || peek(0).Type == lexer.TokLAngle || peek(0).Type == lexer.TokRAngle:
										
						default:
							
					}	
				} else if peek(0).Type == lexer.TokColon {
					// TODO: make POINT vars in ParseExpy
					Variables = append(Variables, Variable_Static{Name: name, Type: POINT, Value: NULL, Scope: 1})
				} else if peek(0).Type == lexer.TokPlus && peek(1).Type == lexer.TokPlus {
					// TODO: implement ++ system
					expect(lexer.TokPlus)
					expect(lexer.TokPlus)
					_var := LookupVariable(name, true, Scope, _ntok, &tokens)
					Write("mov r4, " + _var.Real, true)
					Write("lodf r4, r5", true)
					Write("inc r5", true)
					Write("strf r4, r5", true);
					expect(lexer.TokSemi)
				} else {
					
				}
			case lexer.TokIf:
				expect(lexer.TokIf)
				expect(lexer.TokLParen)
				var exptokens = []lexer.Token {}
				var bodytokens = []lexer.Token {}
				var elsetokens = []lexer.Token {}

				depth := 1
				for j := i; j < len(tokens); j++ {
					if tokens[j].Type == lexer.TokRParen {
						depth--
						if depth == 0 {
							i = j
							break
						}
					} else if tokens[j].Type == lexer.TokLParen {
						depth++
					}
					exptokens = append(exptokens, tokens[j])
				}
				exptokens = append(exptokens, lexer.Token{Type: lexer.TokSemi, Value: ";"})
				Parse(exptokens, CreateScope(Scope))
				expect(lexer.TokRParen)	
				if peek(0).Type == lexer.TokLCurly {
					expect(lexer.TokLCurly)

					depth := 1
					for j := i; j < len(tokens); j++ {	
						if tokens[j].Type == lexer.TokRCurly {
							depth--
							if depth == 0 {
								i = j	
								break
							}
						} else if tokens[j].Type == lexer.TokLCurly {
							depth++
						}
						bodytokens = append(bodytokens, tokens[j])
					}	
					expect(lexer.TokRCurly)	
				} else {
					for j := i; j < len(tokens); j++ {
						bodytokens = append(bodytokens, tokens[j])
						if tokens[j].Type == lexer.TokSemi {
							i = j
							break
						}	
					}
					expect(lexer.TokSemi)
				}
				afterName := topLevelName + "_after_" + fmt.Sprintf("%d", IDCounter)
				IDCounter++
				elseName := "else_stmt_" + fmt.Sprintf("%d", IDCounter)
				IDCounter++	
				if peek(0).Type == lexer.TokElse {
					Write("jz e6, " + elseName, true)
				} else {
					Write("jz e6, " + afterName, true)
				}
				Write("if_stmt_" + fmt.Sprintf("%v", IDCounter) + ":", false)
				otln := topLevelName
				topLevelName = "if_stmt_" + fmt.Sprintf("%d", IDCounter)
				IDCounter++
				Parse(bodytokens, CreateScope(Scope))	
				topLevelName = otln
				IDCounter++
				Write("jmp " + afterName, true)
				if peek(0).Type == lexer.TokElse {
					expect(lexer.TokElse)
					Write(elseName + ":", false)
					if peek(0).Type == lexer.TokLCurly {
						expect(lexer.TokLCurly)
						depth := 1
						for j := i; j < len(tokens); j++ {	
							if tokens[j].Type == lexer.TokRCurly {
								depth--
								if depth == 0 {
									i = j	
									break
								}
							} else if tokens[j].Type == lexer.TokLCurly {
								depth++
							}
							elsetokens = append(elsetokens, tokens[j])
						}	
						expect(lexer.TokRCurly)
					} else {
						for j := i; j < len(tokens); j++ {
							elsetokens = append(elsetokens, tokens[j])
							if tokens[j].Type == lexer.TokSemi {
								i = j
								break
							}	
						}	
					}	
					oltn := topLevelName
					topLevelName = elseName	
					IDCounter++
					Parse(elsetokens, CreateScope(Scope))
					topLevelName = oltn
				}
				Write(afterName + ":", false)	
			case lexer.TokType, lexer.TokQualifier:
				bodytokens := []lexer.Token {}
				for j := i; j < len(tokens); j++ {
					bodytokens = append(bodytokens, tokens[j])
					if tokens[j].Type == lexer.TokSemi {
						i = j
						break
					}
				}
				level = 0
				Parse(bodytokens, Scope)
				level = 1
			case lexer.TokFor:
				expect(lexer.TokFor)
				expect(lexer.TokLParen)

				init_tokens := []lexer.Token {}
				for j := i; j < len(tokens); j++ {
					init_tokens = append(init_tokens, tokens[j])
					if tokens[j].Type == lexer.TokSemi {
						j = i
						break
					}
				}
			default:
				error.Error(1, "'" + tokens[i].Value + "'", tokens[i], &tokens)
			}	
		}
	}
} 
