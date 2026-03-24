package parser

import (
	"lcc1/lexer"
	"lcc1/error"
	"strings"
	"fmt"	
	"strconv"
	"os"
	"math"
	// "runtime/debug"
)

var level int = 0
var L1_ALLOW_NONCONST bool = false

var Code1 string = ""
var Code2 string = ""


var IDCounter = 1

const (
	NUMBER8 int = iota // unsigned short short int
	NUMBER16		   // unsigned short int / unsigned int
	NUMBER32           // unsigned long int
	STRING             // unsigned char
	POINT              // goto labels
	NULL               // void / void*
)

type Variable_Static struct {
	Name string
	Type int
	Type2 int
	Value any
	Pointer bool
	Real string
	Scope int
	Const bool
	Extern bool
	ArgNum int
	Register bool
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

type UnpackOrder struct {
	Register string
	Label string
	Type int
	Pointer bool
}

var Variables = []Variable_Static {
	Variable_Static{Name: "_r0", Real: "r0", Register: true, Scope: 1},
	Variable_Static{Name: "_r1", Real: "r1", Register: true, Scope: 1},
	Variable_Static{Name: "_r2", Real: "r2", Register: true, Scope: 1},
	Variable_Static{Name: "_r3", Real: "r3", Register: true, Scope: 1},
	Variable_Static{Name: "_r4", Real: "r4", Register: true, Scope: 1},
	Variable_Static{Name: "_r5", Real: "r5", Register: true, Scope: 1},
	Variable_Static{Name: "_r6", Real: "r6", Register: true, Scope: 1},
	Variable_Static{Name: "_r7", Real: "r7", Register: true, Scope: 1},
	Variable_Static{Name: "_r8", Real: "r8", Register: true, Scope: 1},
	Variable_Static{Name: "_r9", Real: "r9", Register: true, Scope: 1},
	Variable_Static{Name: "_r10", Real: "r10", Register: true, Scope: 1},
	Variable_Static{Name: "_r11", Real: "r11", Register: true, Scope: 1},
	Variable_Static{Name: "_r12", Real: "r12", Register: true, Scope: 1},
	Variable_Static{Name: "_e0", Real: "e0", Register: true, Scope: 1},
	Variable_Static{Name: "_e1", Real: "e1", Register: true, Scope: 1},
	Variable_Static{Name: "_e2", Real: "e2", Register: true, Scope: 1},
	Variable_Static{Name: "_e3", Real: "e3", Register: true, Scope: 1},
	Variable_Static{Name: "_e4", Real: "e4", Register: true, Scope: 1},
	Variable_Static{Name: "_e5", Real: "e5", Register: true, Scope: 1},
	Variable_Static{Name: "_e6", Real: "e6", Register: true, Scope: 1},
	Variable_Static{Name: "_e7", Real: "e7", Register: true, Scope: 1},
	Variable_Static{Name: "_e8", Real: "e8", Register: true, Scope: 1},
	Variable_Static{Name: "_e9", Real: "e9", Register: true, Scope: 1},
	Variable_Static{Name: "_e10", Real: "e10", Register: true, Scope: 1},
	Variable_Static{Name: "_e11", Real: "e11", Register: true, Scope: 1},
	Variable_Static{Name: "_e12", Real: "e12", Register: true, Scope: 1},
	Variable_Static{Name: "_sp", Real: "sp", Register: true, Scope: 1},
	Variable_Static{Name: "_pc", Real: "pc", Register: true, Scope: 1},
	Variable_Static{Name: "_irv", Real: "irv", Register: true, Scope: 1},
	Variable_Static{Name: "_ir", Real: "ir", Register: true, Scope: 1},
	Variable_Static{Name: "_b", Real: "b", Register: true, Scope: 1},	
}

var FunctionDecls = []FunctionDecl {}
var PIE bool

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

type GlobalInt interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
    ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

func ReturnIntType(i int) string {
	switch {	
	case i <= math.MaxUint8 && i >= 0:
		return "unsigned short short int"
	case i <= math.MaxInt8 && i >= math.MinInt8:
		return "signed short short int"	
	case i <= math.MaxUint16 && i >= 0:
		return "unsigned short int"
	case i <= math.MaxInt16 && i >= math.MinInt16:
		return "signed short int"	
	case i <= math.MaxUint32 && i >= 0:
		return "unsigned long int"
	case i <= math.MaxInt32 && i >= math.MinInt32:
		return "signed long int"
	}
	return "unsigned short int"
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

// Some globals (i know its bad practice but it works so....)
var CMP_OP string = ""
var _CMP_MOP_REVERSE string = ""
var _CMP_MOP string = ""
var _BREAK_TOPLEVEL string = ""
var _CONTINUE_TOPLEVEL string = ""
func ParseExpy(tokens []lexer.Token, start int, Scope int, register string) int {
	i := start
	// CMP := false
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
			depth := 1
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
					if depth == 1 {
						ParseExpy(CurrentTokens, 0, Scope, "r7")
						Write("push r7", true)
						CurrentTokens = []lexer.Token{}
						pushed++
					} else {
						CurrentTokens = append(CurrentTokens, tokens[j])
					}
				case lexer.TokLParen:
					depth++
					CurrentTokens = append(CurrentTokens, tokens[j])
				case lexer.TokRParen:
					depth--
					if depth == 0 {
						exit = true
						break
					} else {
						CurrentTokens = append(CurrentTokens, tokens[j])
					}
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
			i = j
		}
	}
	IDENT_STRING := func(label string) string {
		if label[len(label) - 1] != '"' {
			// TODO: fix "literal not terminated and have us handle it"
			error.Error(32, "\" character", peek(-1), &tokens)
		}
		_label := fmt.Sprintf("var_%d", IDCounter)
		IDCounter++
		__label := fmt.Sprintf("var_%d", IDCounter)
		IDCounter++

		WritePre(_label + ":", false)
		WritePre(".asciz \"" + strings.ReplaceAll(label, "\"", "") + "\"", true)
		WritePre(__label + ":", false)
		WritePre(".ptrlabel " + _label, true)
		WritePre(".ptr " + _label, true)
		// Write("mov " + register + ", " + _label, true)

		return __label
	}
	_IDENT_INTENT := func(pointer bool, _type int, deref int, register bool) string {
		if peek(0).Type == lexer.TokEqual {
			// Write intent (NEVER give one free dereference)
			Write("mov r2, r1", true)
			return "write"
		} else {
			// Read intent (ALWAYS give one free dereference)
			if deref >= 0 {
				if register == false {
					switch _type {
					case NUMBER8, STRING, NULL:
						Write("lod r1, r2", true)
					case NUMBER16, NUMBER32:
						Write("lodf r1, r2", true) 	
					}
				} else {
					Write("mov r2, r1", true)
				}
			} else {
				if register == false {
					Write("mov r2, r1", true)
				} else {
					error.Error(37, "", peek(-1), &tokens)
				}
			}
		}
		return "read"
	}
	_NUMBER_PARSE := func(register string) {
		switch peek(0).Type {
		case lexer.TokIdent, lexer.TokStar, lexer.TokAmpersand:
			subslice := []lexer.Token {}
			fl_exit := false
			for {
				if fl_exit == true {
					break
				}
				switch peek(0).Type {
				case lexer.TokStar, lexer.TokAmpersand:
					subslice = append(subslice, peek(0))
					i++
				case lexer.TokIdent:
					subslice = append(subslice, peek(0))
					i++
					fl_exit = true
				}
			}
			ParseExpy(subslice, 0, Scope, register)
		case lexer.TokNumber:
			Write("mov " + register + ", " + peek(0).Value, true)
			i++
		}
	}

	deref := 0
	EQU_VT := NULL
	EQU_VAR := Variable_Static{}
	var op string = ""
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
		goto DONE
	case lexer.TokType, lexer.TokQualifier:
		bodytokens := []lexer.Token {}
		for j := i; j < len(tokens); j++ {
			bodytokens = append(bodytokens, tokens[j])
			if tokens[j].Type == lexer.TokSemi {
				i = j + 1
				break
			}
		}
		level = 0
		L1_ALLOW_NONCONST = true
		Parse(bodytokens, Scope)
		L1_ALLOW_NONCONST = true
		level = 1
		goto DONE
	case lexer.TokIf:
		// TODO: implement quick ifs

		IfScope := CreateScope(Scope)
		ElseScope := CreateScope(Scope)

		expect(lexer.TokIf)
		expect(lexer.TokLParen)

		exp_tokens := []lexer.Token {}

		depth := 1
		exit := false
		for _ = i; i < len(tokens); i++ {	
			switch tokens[i].Type {
			case lexer.TokLParen:
				depth++
				exp_tokens = append(exp_tokens, tokens[i])
			case lexer.TokRParen:
				depth--
				if depth == 0 {
					exit = true
					break
				} else {
					exp_tokens = append(exp_tokens, tokens[i])
				}
			default:
				exp_tokens = append(exp_tokens, tokens[i])
			}
			if exit == true {
				// fmt.Println(tokens[i].Value)
				break
			}
		}
		ParseExpy(exp_tokens, 0, Scope, "r11") // r12 and r5 clobbered
											   // r11 result

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
		j := i
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

		cmop := ""
		cmopr := ""
		if _CMP_MOP == "" {
			cmop = "jnz"
		} else {
			cmop = _CMP_MOP
		}

		if _CMP_MOP_REVERSE == "" {
			cmopr = "jz"
		} else {
			cmopr = _CMP_MOP_REVERSE
		}

		if peek(0).Type != lexer.TokElse {
			// Write everything
			Write(cmop + " r11, " + if_label, true)
			Write(cmopr + " r11, " + after_label, true)
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
		Write(cmop + " r11, " + if_label, true)
		Write(cmopr + " r11, " + else_label, true)	
		Write(if_label + ":", false)
		ParseExpyL1(if_tokens, 0, IfScope)
		Write("jmp " + after_label, true)
		Write(else_label + ":", false)
		ParseExpyL1(else_tokens, 0, ElseScope)
		Write(after_label + ":", false)
		goto DONE
	case lexer.TokWhile:
		expect(lexer.TokWhile)
		expect(lexer.TokLParen)
		
		subslice := []lexer.Token {}

		depth := 1
		exit := false

		top_label := "while_stmt_" + fmt.Sprintf("%d", IDCounter) + "_check"	
		middle_label := "while_stmt_" + fmt.Sprintf("%d", IDCounter) + "_body"
		bottom_label := "while_stmt_" + fmt.Sprintf("%d", IDCounter) + "_after"

		IDCounter++

		for _ = i; i < len(tokens); i++ {
			switch peek(0).Type {
			case lexer.TokLParen:
				depth++
				subslice = append(subslice, peek(0))
			case lexer.TokRParen:
				depth--
				if depth < 1 {
					exit = true
				} else {
					subslice = append(subslice, peek(0))
				}
			default:
				subslice = append(subslice, peek(0))
			}
			if exit == true {
				break
			}
		}
		expect(lexer.TokRParen)

		// Write check portion
		
		Write(top_label + ":", false)
		ParseExpy(subslice, 0, Scope, "r11")

		cmop := ""
		if _CMP_MOP == "" {
			cmop = "jnz"
		} else {
			cmop = _CMP_MOP
		}

		Write(cmop + " r11, " + middle_label, true)
		Write("jmp " + bottom_label, true)

		expect(lexer.TokLCurly)
		
		subslice2 := []lexer.Token {}
		
		depth = 1
		exit = false
		for _ = i; i < len(tokens); i++ {
			switch peek(0).Type {
			case lexer.TokLCurly:
				depth++
				subslice2 = append(subslice2, peek(0))
			case lexer.TokRCurly:
				depth--
				if depth == 0 {
					exit = true
				} else {
					subslice2 = append(subslice2, peek(0))
				}
			default:
				subslice2 = append(subslice2, peek(0))
			}
			if exit == true {
				break
			}
		}
		expect(lexer.TokRCurly)

		WScope := CreateScope(Scope)

		otln_br := _BREAK_TOPLEVEL
		_BREAK_TOPLEVEL = bottom_label
		otln_co := _CONTINUE_TOPLEVEL	
		_CONTINUE_TOPLEVEL = middle_label

		Write(middle_label + ":", false)
		ParseExpyL1(subslice2, 0, WScope)
		Write("jmp " + top_label, true)
		Write(bottom_label + ":", false)

		_BREAK_TOPLEVEL = otln_br
		_CONTINUE_TOPLEVEL = otln_co
	case lexer.TokDo:
		expect(lexer.TokDo)
		expect(lexer.TokLCurly)

		top_label := "do_stmt_" + fmt.Sprintf("%d", IDCounter) + "_top"
		middle_label := "do_stmt_" + fmt.Sprintf("%d", IDCounter) + "_check"
		bottom_label := "do_stmt_" + fmt.Sprintf("%d", IDCounter) + "_after"

		IDCounter++
		
		depth := 1
		exit := false
		subslice := []lexer.Token {}

		for _ = i; i < len(tokens); i++ {
			switch peek(0).Type {
			case lexer.TokLCurly:
				depth++
				subslice = append(subslice, peek(0))
			case lexer.TokRCurly:
				depth--
				if depth == 0 {
					exit = true
				} else {
					subslice = append(subslice, peek(0))
				}
			default:
				subslice = append(subslice, peek(0))
			}
			if exit == true {
				break
			}
		}
		expect(lexer.TokRCurly)
		expect(lexer.TokWhile)
		expect(lexer.TokLParen)
		
		subslice2 := []lexer.Token {}
		exit = false
		depth = 1
		for _ = i; i < len(tokens); i++ {
			switch peek(0).Type {
			case lexer.TokLParen:
				depth++
				subslice2 = append(subslice2, peek(0))
			case lexer.TokRParen:
				depth--
				if depth == 0 {
					exit = true
				} else {
					subslice2 = append(subslice2, peek(0))
				}
			default:
				subslice2 = append(subslice2, peek(0))
			}
			if exit == true {
				break
			}
		}

		expect(lexer.TokRParen)
		expect(lexer.TokSemi)

		otln := _BREAK_TOPLEVEL
		_BREAK_TOPLEVEL = bottom_label
		otln_co := _CONTINUE_TOPLEVEL	
		_CONTINUE_TOPLEVEL = middle_label

		Write(top_label + ":", false)
		DScope := CreateScope(Scope)
		ParseExpyL1(subslice, 0, DScope)
		Write(middle_label + ":", false)
		ParseExpy(subslice2, 0, DScope, "r11")

		cmop := ""
		if _CMP_MOP == "" {
			cmop = "jnz"
		} else {
			cmop = _CMP_MOP
		}

		Write(cmop + " r11, " + top_label, true)
		Write("jmp " + bottom_label, true)

		Write(bottom_label + ":", false)

		_BREAK_TOPLEVEL = otln
		_CONTINUE_TOPLEVEL = otln_co
	case lexer.TokFor:
		expect(lexer.TokFor)
		expect(lexer.TokLParen)

		top_label := "for_stmt_" + fmt.Sprintf("%d", IDCounter) + "_check"
		bottom_label := "for_stmt_" + fmt.Sprintf("%d", IDCounter) + "_after"

		IDCounter++
		
		subslice := []lexer.Token {}
		exit := false
		for _ = i; i < len(tokens); i++ {
			if exit == true {
				break
			}
			switch peek(0).Type {
			case lexer.TokSemi:
				subslice = append(subslice, peek(0))
				exit = true
			default:
				subslice = append(subslice, peek(0))
			}
		}
		
		subslice2 := []lexer.Token {}
		exit = false
		switch peek(0).Type {
		case lexer.TokSemi:
			subslice2 = append(subslice2, lexer.Token{Type: lexer.TokNumber, Value: "1", Line: peek(0).Line, File: peek(0).File})
			expect(lexer.TokSemi)
			goto COND_DONE	
		}	

		for _ = i; i < len(tokens); i++ {
			if exit == true {
				break
			}
			switch peek(0).Type {
			case lexer.TokSemi:
				subslice2 = append(subslice2, peek(0))
				exit = true
			default:
				subslice2 = append(subslice2, peek(0))
			}
		}

		COND_DONE:

		FScope := CreateScope(Scope)
		ParseExpyL1(subslice, 0, FScope) // Initialize variable

		Write(top_label + ":", false)
		ParseExpy(subslice2, 0, FScope, "r11")

		cmopr := ""
		if _CMP_MOP_REVERSE == "" {
			cmopr = "jz"
		} else {
			cmopr = _CMP_MOP_REVERSE
		}

		Write(cmopr + " r11, " + bottom_label, true)

		subslice3 := []lexer.Token {}
		exit = false
		depth := 1

		for _ = i; i < len(tokens); i++ {
			switch peek(0).Type {
			case lexer.TokLParen:
				depth++
				subslice3 = append(subslice3, peek(0))
			case lexer.TokRParen:
				depth--
				if depth == 0 {
					exit = true
				} else {
					subslice3 = append(subslice3, peek(0))
				}
			default:
				subslice3 = append(subslice3, peek(0))
			}
			if exit == true {
				break
			}
		}

		expect(lexer.TokRParen)
		expect(lexer.TokLCurly)

		subslice4 := []lexer.Token {}
		depth = 1
		exit = false
		
		for _ = i; i < len(tokens); i++ {
			switch peek(0).Type {
			case lexer.TokLCurly:
				depth++
				subslice4 = append(subslice4, peek(0))
			case lexer.TokRCurly:
				depth--
				if depth == 0 {
					exit = true
				} else {
					subslice4 = append(subslice4, peek(0))
				}
			default:
				subslice4 = append(subslice4, peek(0))
			}
			if exit == true {
				break
			}
		}
		subslice3 = append(subslice3, lexer.Token{Type: lexer.TokSemi, Value: ";", Line: peek(-1).Line, File: peek(-1).File})

		otln := _BREAK_TOPLEVEL
		_BREAK_TOPLEVEL = bottom_label
		otln_co := _CONTINUE_TOPLEVEL	
		_CONTINUE_TOPLEVEL = top_label

		ParseExpyL1(subslice4, 0, FScope)
		ParseExpyL1(subslice3, 0, FScope)
		Write("jmp " + top_label, true)
		Write(bottom_label + ":", false)

		_BREAK_TOPLEVEL = otln
		_CONTINUE_TOPLEVEL = otln_co
		expect(lexer.TokRCurly)
	case lexer.TokContinue:
		expect(lexer.TokContinue)
		if _CONTINUE_TOPLEVEL == "" {
			error.Error(39, "", peek(-1), &tokens)
		} else {
			Write("jmp " + _CONTINUE_TOPLEVEL, true)
		}
		expect(lexer.TokSemi)
	case lexer.TokBreak:
		expect(lexer.TokBreak)
		if _BREAK_TOPLEVEL == "" {
			error.Error(38, "", peek(-1), &tokens)
		} else {
			Write("jmp " + _BREAK_TOPLEVEL, true)
		}
		expect(lexer.TokSemi)
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
		goto DONE
	case lexer.TokIdent:
		label := expect(lexer.TokIdent)
		var variable Variable_Static

		if label[0] == '"' {
			LTL := IDENT_STRING(label)
			Write("mov r1, " + LTL, true)
			_IDENT_INTENT(true, NUMBER16, 0, false)
			Write("mov " + register + ", r2", true)
			goto CONTINUE
		}

		array := false
		switch peek(0).Type {
		case lexer.TokLParen:
			// TODO: allow comma arbitration
			IDENT_FUNC(label)	
			goto CONTINUE
		case lexer.TokColon:
			Write(label + ":", false)
			expect(lexer.TokColon)
			goto DONE	
		}	

		variable = LookupVariable(label, true, Scope, peek(-1), &tokens)

		switch peek(0).Type {
		case lexer.TokLBracket:
			if variable.ArgNum < 1 {
				error.Error(40, "", peek(0), &tokens)
			}
			array = true
			expect(lexer.TokLBracket)

			subslice := []lexer.Token {}
			depth := 1
			exit := false
			for _ = i; i < len(tokens); i++ {
				switch peek(0).Type {
				case lexer.TokLBracket:
					depth++
					subslice = append(subslice, peek(0))
				case lexer.TokRBracket:
					depth--
					if depth == 0 {
						exit = true
					} else {
						subslice = append(subslice, peek(0))
					}
				default:
					subslice = append(subslice, peek(0))
				}
				if exit == true {
					break
				}
			}
			expect(lexer.TokRBracket)

			ParseExpy(subslice, 0, Scope, "e8")
		}

		EQU_VT = variable.Type
		EQU_VAR = variable
		Write("mov r1, " + variable.Real, true)
		if array == true {
			switch variable.Type {
			case NUMBER8:
				Write("add r1, r1, e8", true)
			case NUMBER16:
				Write("mov e7, 2", true)
				Write("mul e8, e8, e7", true)
				Write("add r1, r1, e8", true)
			case NUMBER32:
				Write("mov e7, 4", true)
				Write("mul e8, e8, e7", true)
				Write("add r1, r1, e8", true)
			}
		}

		Intent := _IDENT_INTENT(variable.Pointer, variable.Type, deref, variable.Register)

		x_deref := deref
		derefs := 0
		for x_deref > 0 {
			x_deref--
			if variable.Pointer == true && ((derefs > 0 && Intent == "write") || (Intent == "read")) {
				if PIE == true {
					Write("add r2, r2, e14", true)
				}
				switch variable.Type2 {
				case NUMBER8, STRING, NULL:
					Write("lod r2, r2", true)	
				case NUMBER16, NUMBER32:
					Write("lodf r2, r2", true)
				}
			} else {
				switch variable.Type {
				case NUMBER8, STRING, NULL:
					Write("lod r2, r2", true)
				case NUMBER16, NUMBER32:
					Write("lodf r2, r2", true)
				}
			}
			derefs++
		}
		Write("mov " + register + ", r2", true)	
	case lexer.TokNumber:
		// Parse expressions
		// Load it up into r4

		_NUMBER_PARSE(register)	
	}

	switch peek(0).Type {
	case lexer.TokPlus, lexer.TokMinus, lexer.TokStar, lexer.TokSlash:
		OP_TRY:
		op = expect(peek(0).Type)
		_NUMBER_PARSE("r6")

		switch op {
		case "+":
			Write("add " + register + ", " + register + ", r6", true)
		case "-":
			Write("sub " + register + ", " + register + ", r6", true)
		case "*":
			Write("mul " + register + ", " + register + ", r6", true)
		case "/":
			Write("div " + register + ", " + register + ", r6", true)
		}
		switch peek(0).Type {
		case lexer.TokPlus, lexer.TokMinus, lexer.TokStar, lexer.TokSlash:
			goto OP_TRY
		}
	}
	
	CONTINUE:
	if peek(0).Type != lexer.TokEqual && peek(0).Type != lexer.TokEquality && peek(0).Type != lexer.TokInequality && peek(0).Type != lexer.TokGEqual && peek(0).Type != lexer.TokLEqual && peek(0).Type != lexer.TokLAngle && peek(0).Type != lexer.TokRAngle {
		// expect(lexer.TokSemi)
		return i
	}
	
	switch peek(0).Type {
	case lexer.TokEquality, lexer.TokInequality, lexer.TokGEqual, lexer.TokLEqual, lexer.TokLAngle, lexer.TokRAngle:
		cmpopreal := ""
		expect(peek(0).Type)
		switch peek(-1).Type {
		case lexer.TokEquality:
			_CMP_MOP = "jnz"
			_CMP_MOP_REVERSE = "jz"
			if CMP_OP == "" {
				cmpopreal = "cmp"
			} else {
				cmpopreal = CMP_OP
			}
		case lexer.TokInequality:
			_CMP_MOP = "jz"
			_CMP_MOP_REVERSE = "jnz"
			if CMP_OP == "" {
				cmpopreal = "cmp"
			} else {
				cmpopreal = CMP_OP
			}
		case lexer.TokGEqual:
			cmpopreal = "ilt"
			_CMP_MOP = "jz"
			_CMP_MOP_REVERSE = "jnz"
		case lexer.TokLEqual:
			cmpopreal = "igt"
			_CMP_MOP = "jz"
			_CMP_MOP_REVERSE = "jnz"
		case lexer.TokLAngle:
			cmpopreal = "ilt"
			_CMP_MOP = "jnz"
			_CMP_MOP_REVERSE = "jz"
		case lexer.TokRAngle:
			cmpopreal = "igt"
			_CMP_MOP = "jnz"
			_CMP_MOP_REVERSE = "jz"
		}

		i = ParseExpy(tokens, i, Scope, "r5")
		// Universal IF register = r11
	
		Write(cmpopreal + " r11, " + register + ", r5", true)	
	default:
		expect(lexer.TokEqual)
		i = ParseExpy(tokens, i, Scope, "r5")

		if EQU_VAR.Const == true {
			error.Error(33, "'" + EQU_VAR.Name + "' with const-qualified type", peek(-1), &tokens)
			token, stream := FuncDeclLookup(EQU_VAR.Name)
			error.Note(22, "'" + EQU_VAR.Name + "' declared here", token, stream)
		}

		if EQU_VAR.Pointer == false || (EQU_VAR.Pointer == true && deref < 1) {	
			switch EQU_VT {
			case NUMBER8, STRING, NULL:
				Write("str " + register + ", r5", true)
			case NUMBER16, NUMBER32:
				Write("strf " + register + ", r5", true)
			}
		} else {
			switch EQU_VAR.Type2 {
			case NUMBER8, STRING, NULL:
				Write("str " + register + ", r5", true)
			case NUMBER16, NUMBER32:
				Write("strf " + register + ", r5", true)
			}
		}
		

		expect(lexer.TokSemi)	
	}

	DONE:

	return i
}

func ParseNumberExpyDirect(tokens []lexer.Token, i int, Scope int) (int, int) {
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
			if toktype != lexer.TokSemi && tokens[i].Type != lexer.TokIdent {
				error.Error(1, "'" + tokens[i].Value + "'", tokens[i], &tokens)
			} else if toktype == lexer.TokSemi {
				error.Error(18, "", tokens[i - 1], &tokens)
			} else if tokens[i].Type == lexer.TokIdent {
				error.Error(35, "", tokens[i], &tokens)
			}
			i++
		}
		return value
	}	
	peek := func(lookahead int) lexer.Token {
		if i + lookahead < len(tokens) && i + lookahead >= 0 {
			return tokens[i + lookahead]
		}
		return lexer.Token{Type: lexer.TokEOF, Value: ""}
	}
	
	res := 0

	for {
		if i >= len(tokens) {
			break
		}
		
		exit := false
		exit_nodet := false

		switch peek(0).Type {
		case lexer.TokAmpersand:
			expect(lexer.TokAmpersand)
			label := expect(lexer.TokIdent)
			Variable := LookupVariable(label, true, Scope, peek(-1), &tokens)
			WritePre(".ptr " + Variable.Real, true)
			WritePre(".ptrlabel " + label, true)
			return -1, i
		}
		num := expect(lexer.TokNumber)

		OP_TRY:

		switch peek(0).Type {
		case lexer.TokPlus, lexer.TokMinus, lexer.TokStar, lexer.TokSlash:
		default:
			exit = true
		}
		if exit == true {
			if exit_nodet == false {
				n1_real, _ := strconv.ParseInt(num, 0, 64)
				res = int(n1_real)
			}
			break
		}
		
		op := peek(0).Value	
		expect(peek(0).Type)
		num2 := expect(lexer.TokNumber)

		n1_real, _ := strconv.ParseInt(num, 0, 64)
		n2_real, _ := strconv.ParseInt(num2, 0, 64)

		switch op {
		case "+":
			res = int(n1_real) + int(n2_real)
		case "-":
			res = int(n1_real) - int(n2_real)
		case "*":
			res = int(n1_real) * int(n2_real)
		case "/":
			res = int(n1_real) / int(n2_real)
		}

		exit_nodet = true 
		goto OP_TRY
	}

	return res, i 
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
			return ""
		}
		if tokens[i].Type == toktype {
			value = tokens[i].Value
		} else {
			if toktype != lexer.TokSemi {
				error.Error(1, "'" + tokens[i].Value + "'", tokens[i], &tokens)
			} else {
				error.Error(18, "", tokens[i - 1], &tokens)
			}
		}
		i++
		return value
	}	
	peek := func(lookahead int) lexer.Token {
		if i + lookahead < len(tokens) {
			return tokens[i + lookahead]
		}
		return lexer.Token{Type: lexer.TokEOF, Value: ""}
	}

	_PARSE_ATTR := func(name string) []string {
		var attrs []string
		var _RETURNS []string

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
				_RETURNS = append(_RETURNS, "norename")	
			case "noreturn":
				_RETURNS = append(_RETURNS, "noreturn")
			case "require_const":
				_RETURNS = append(_RETURNS, "require_const")
			default:
				error.Warning(11, "'" + attr + "'", tokens[i - 3], &tokens)
			}
		}
		return _RETURNS
	}

	_PARSE_TYPE := func() (int, bool) {
		ptr := false
		long := false
		short := false
		shortshort := false
		unsigned := false
		constant := false
		bits := BitPref

		for {
			if i >= len(tokens) {
				break
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
					if unsigned == true {
						error.Error(28, "'unsigned'", peek(-1), &tokens)
					}
					unsigned = true
				case "const":
					if constant == true {
						error.Error(28, "'const'", peek(-1), &tokens)
					}
					constant = true
				}
			} else {
				_type := expect(lexer.TokType)
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
				case "void":
					rtype = NULL
				}

				if peek(0).Type == lexer.TokStar {
					ptr = true
					i++
				} else {
					if rtype == NULL {
						error.Error(7, "'void'", peek(-1), &tokens)	
					}
				}

				return rtype, ptr
				break
			}
		}
		return NUMBER16, false
	}
	
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
			breakoff := false
			bits := BitPref	
			for {
				if i >= len(tokens) {
					breakoff = true
					break
				}
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
				if peek(0).Value == "__embed__" {
					pre := false
					static := false
					name := ""

					expect(lexer.TokIdent)



					arg_parse_top:
					switch peek(0).Value {
					case "pre":
						pre = true
						i++
						goto arg_parse_top
					case "static":
						static = true
						i++
						goto arg_parse_top
					}	

					name = expect(lexer.TokIdent)
					// __embed__ pre SOUND_LABEL

					expect(lexer.TokLParen)
					expect(lexer.TokLParen)
					path := expect(lexer.TokIdent)
					expect(lexer.TokRParen)
					expect(lexer.TokRParen)
					expect(lexer.TokSemi)

					if pre == false {
						Write(name + ":", false)
						Write(".embed \"" + strings.ReplaceAll(path, "\"", "") + "\"", true)
					} else {
						WritePre(name + ":", false)
						WritePre(".embed \"" + strings.ReplaceAll(path, "\"", "") + "\"", true)
					}
					if static == false {
						PreWrite(".global " + name, false)
					}
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
						if unsigned == true {
							error.Error(28, "'unsigned'", peek(-1), &tokens)
						}
						unsigned = true
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
			if breakoff == true {
				break
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

			var attrs []string
			var UnpackOrders []UnpackOrder

			allow_nonconst := L1_ALLOW_NONCONST

			switch peek(0).Type {
			case lexer.TokLParen:
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
				case lexer.TokType, lexer.TokQualifier:
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

							__arg_reg := fmt.Sprintf("e%d", register)
							__rn := fmt.Sprintf("var_%d", IDCounter)
							IDCounter++
							__rtype, __ptr := _PARSE_TYPE()
							__name := expect(lexer.TokIdent)
							
							if extern == true {
								goto ARG_DECL_DONE
							}

							UnpackOrders = append(UnpackOrders, UnpackOrder{Register: __arg_reg, Label: __rn, Type: __rtype, Pointer: __ptr})

							WritePre(__rn + ":", false)
							switch __rtype {
							case NUMBER8:
								WritePre(".byte 0x00", true)
							case STRING:	
								if __ptr == false {
									WritePre(".byte 0x00", true)
								} else {
									WritePre(".ptrlabel " + __rn, true)
									WritePre(".ptr 0x00", true)
								}
							case NULL:
								WritePre(".ptr 0x00", true)
								WritePre(".ptrlabel " + __rn, true)
							case NUMBER16:
								WritePre(".word 0x0000", true)
							case NUMBER32:
								WritePre(".dword 0x00000000", true)
							}

							ARG_DECL_DONE:
							if __ptr == false {
								Variables = append(Variables, Variable_Static{Name: __name, Type: __rtype, Value: nil, Scope: fscope, Real: __rn, Pointer: __ptr})
							} else {
								Variables = append(Variables, Variable_Static{Name: __name, Type: NUMBER16, Type2: __rtype, Value: nil, Scope: fscope, Real: __rn, Pointer: __ptr})
							}
							
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
					attrs := _PARSE_ATTR(name)
					for _, attr := range attrs {
						switch attr {
						case "noreturn":
							noreturn = true
						case "norename":
							if rns == true {
								name = "main"
							}
						default:
							error.Warning(34, "'" + attr + "' not allowed here", peek(-3), &tokens)
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

				if name == "_start" {
					PreWrite("jmp _start", false)
				}

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

					for _, UnpackOrder := range UnpackOrders {
						__rn := UnpackOrder.Label
						__arg_reg := UnpackOrder.Register
						__ptr := UnpackOrder.Pointer
						__rtype := UnpackOrder.Type

						switch __rtype {
						case NUMBER8:
							Write("mov r1, " + __rn, true)
							Write("str r1, " + __arg_reg, true)
						case STRING:
							if __ptr == false {
								Write("mov r1, " + __rn, true)
								Write("str r1, " + __arg_reg, true)
							} else {
								Write("mov r1, " + __rn, true)
								Write("strf r1, " + __arg_reg, true)
							}
						case NULL:
							Write("mov r1, " + __rn, true)
							Write("strf r1, " + __arg_reg, true)	
						case NUMBER16:
							Write("mov r1, " + __rn, true)
							Write("strf r1, " + __arg_reg, true)
						case NUMBER32:
							Write("mov r1, " + __rn, true)
							Write("strf r1, " + __arg_reg, true)
						}	
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
			case lexer.TokIdent:
				if peek(0).Value == "__attribute__" {
					attrs = _PARSE_ATTR(name)
					for _, attr := range attrs {
						switch attr {
						case "require_const":
							allow_nonconst = false	
						default:
							error.Warning(34, "'" + attr + "' not allowed here", peek(-3), &tokens)
						}
					}	
				} else {
					// Error out
					expect(lexer.TokEOF)
				}
				fallthrough
			case lexer.TokEqual:	
				expect(lexer.TokEqual)	
				switch _type {
				case "void":
					error.Error(7, "'void'", _typetoken, &tokens)
				case "int":
					_i := 0
					rn := "var_" + fmt.Sprintf("%d", IDCounter)
					IDCounter++
					

					if allow_nonconst == false {
						res := 0
						WritePre(rn + ":", false)
						res, _i = ParseNumberExpyDirect(tokens, i, Scope)

						if res == -1 {
							goto EQU_RTYPE_DONE
						}

						switch rtype {
						case NUMBER8, STRING:
							r_res := uint8(res)
							if res > math.MaxUint8 || res < 0 {
								error.Warning(36, "'" + ReturnIntType(res) + "' to 'unsigned short short int' changes value from '" + fmt.Sprintf("%d", res) + "' to '" + fmt.Sprintf("%d", r_res) + "'", peek(-1), &tokens)
							}
							WritePre(".byte " + fmt.Sprintf("0x%02x", r_res), true)
						case NUMBER16:
							r_res := uint16(res)
							if res > math.MaxUint16 || res < 0 {
								error.Warning(36, "'" + ReturnIntType(res) + "' to 'unsigned short int' changes value from '" + fmt.Sprintf("%d", res) + "' to '" + fmt.Sprintf("%d", r_res) + "'", peek(-1), &tokens)
							}
							WritePre(".word " + fmt.Sprintf("0x%04x", r_res), true)
						case NUMBER32:
							r_res := uint32(res)
							if res > math.MaxUint32 || res < 0 {
								error.Warning(36, "'" + ReturnIntType(res) + "' to 'unsigned long int' changes value from '" + fmt.Sprintf("%d", res) + "' to '" + fmt.Sprintf("%d", r_res) + "'", peek(-1), &tokens)
							}
							WritePre(".dword " + fmt.Sprintf("0x%08x", r_res), true)
						}	
						EQU_RTYPE_DONE:
					} else {
						WritePre(rn + ":", false)
						_i = ParseExpy(tokens, i, Scope, "r4")

						if ptr == false {
							switch rtype {
							case NUMBER8, STRING:
								WritePre(".byte 0x00", true)
								Write("mov r7, " + rn, true)
								Write("str r7, r4", true)
							case NUMBER16:
								WritePre(".word 0x0000", true)
								Write("mov r7, " + rn, true)
								Write("strf r7, r4", true)
							case NUMBER32:
								WritePre(".dword 0x00000000", true)
								Write("mov r7, " + rn, true)
								Write("strf r7, r4", true)
							}
						} else {
							WritePre(".ptrlabel " + rn, true)
							WritePre(".ptr 0x00", true)
							Write("mov r7, " + rn, true)
							Write("strf r7, r4", true)
						}
					}
					i = _i


					var val any
					if ptr == true {	
						Variables = append(Variables, Variable_Static{Name: name, Type: NUMBER16, Type2: rtype, Value: val, Pointer: true, Real: rn, Scope: Scope, Const: constant})
					} else {
						Variables = append(Variables, Variable_Static{Name: name, Type: rtype, Value: val, Pointer: false, Real: rn, Scope: Scope, Const: constant})
					}
				case "char":
					str, end := StringParse(tokens, i)	
					if ptr == true {
						rn := "var_" + fmt.Sprintf("%d", IDCounter)
						IDCounter++
						rn2 := "var_" + fmt.Sprintf("%d", IDCounter)
						IDCounter++
						Variables = append(Variables, Variable_Static{Name: name, Type: NUMBER16, Value: str, Pointer: true, Real: rn2, Scope: Scope, Const: constant})
						WritePre(rn + ":", false)
						WritePre(".asciz \"" + str + "\"", true)
						WritePre(rn2 + ":", false)
						WritePre(".ptrlabel " + rn, true)
						WritePre(".ptr " + rn, true)
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
			case lexer.TokSemi:
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
						rn2 := "var_" + fmt.Sprintf("%d", IDCounter)
						IDCounter++
						Variables = append(Variables, Variable_Static{Name: name, Type: NUMBER16, Type2: STRING, Value: "", Pointer: true, Real: rn2, Scope: Scope, Const: constant})
						WritePre(rn + ":", false)
						WritePre(".asciz \"\"", true)
						WritePre(rn2 + ":", false)
						WritePre(".ptrlabel " + rn, true)
						WritePre(".ptr " + rn, true)
					} else {	
						Variables = append(Variables, Variable_Static{Name: name, Type: STRING, Value: "", Pointer: false, Scope: Scope, Const: constant})
					}
				case "void":
					if ptr == true {
						if extern == false {
							rn := "var_" + fmt.Sprintf("%d", IDCounter)
							IDCounter++
							rn2 := "var_" + fmt.Sprintf("%d", IDCounter)
							IDCounter++
							WritePre(rn2 + ":", false)
							WritePre(".ptr " + rn, true)
							WritePre(rn + ":", false)
							Variables = append(Variables, Variable_Static{Name: name, Type: NUMBER16, Type2: NULL, Value: nil, Pointer: true, Real: rn, Scope: Scope, Const: constant})
						} else {
							rn := "var_" + fmt.Sprintf("%d", IDCounter)
							IDCounter++
							WritePre(rn + ":", false)
							WritePre(".ptrlabel " + rn, true)
							WritePre(".ptr " + name, true)
							Variables = append(Variables, Variable_Static{Name: name, Type: NUMBER16, Value: nil, Pointer: true, Real: rn, Scope: Scope, Const: constant})	
						}
					}
				}
			case lexer.TokLBracket:
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
			default:
				error.Error(1, "'" + peek(0).Value + "'", _typetoken, &tokens)
			}	
		}
	}
} 
