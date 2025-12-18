package parser

import (
	"lcc1/lexer"
	"lcc1/error"
	"strings"
	"fmt"	
	"strconv"
	"os"
	"runtime"
	"path/filepath"
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

var Variables = []Variable_Static {}
var FunctionDecls = []FunctionDecl {}

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
	
	if Enforce {
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

func ParseExpression(tokens []lexer.Token, start int, Scope int) (int, int) {
	if tokens[start].Type == lexer.TokNumber || LookupVariable(tokens[start].Value, false, Scope, tokens[start], &tokens).Type == NUMBER16 || LookupVariable(tokens[start].Value, false, Scope, tokens[start], &tokens).Type == NUMBER32 {
		var evaltokens []lexer.Token
		var end int = 0
		for i := start; i < len(tokens); i++ {
			if tokens[i].Type == lexer.TokSemi {
				end = i - 1
				break
			}
			print(tokens[i].Value + "\n")
			evaltokens = append(evaltokens, tokens[i])
		}

		for i := 0; i < len(evaltokens); i++ {
			if evaltokens[i].Type == lexer.TokNumber || evaltokens[i].Type == lexer.TokIdent || evaltokens[i].Type == lexer.TokStar || evaltokens[i].Type == lexer.TokAmpersand {
				deref := false
				// addr := false
				if evaltokens[i].Type == lexer.TokStar {
					deref = true
					i++
				} else if evaltokens[i].Type == lexer.TokAmpersand {
					// addr = true
					i++
				}
			
				print(evaltokens[i].Value + "\n")
				switch evaltokens[i].Type {
				case lexer.TokNumber:
					Write("mov r4, " + fmt.Sprintf(evaltokens[i].Value), true)
				case lexer.TokIdent:
					print("Looking up (first) " + evaltokens[i].Value + "\n")
					variable := LookupVariable(evaltokens[i].Value, true, Scope, evaltokens[i], &evaltokens)
					if deref == false {
						Write("mov r4, " + variable.Real, true)
						Write("lodf r4, r4", true)
					} else {
						Write("mov r4, " + variable.Real, true)
						Write("lodf r4, r4", true)
						Write("lodf r4, r4", true)
					}
				default:
					// Error
				}
				i++

				print(evaltokens[i + 1].Value + "\n")
				switch evaltokens[i + 1].Type {
				case lexer.TokNumber:
					Write("mov r5, " + fmt.Sprintf(evaltokens[i + 1].Value), true)
				case lexer.TokIdent:
					print("Looking up (second) " + evaltokens[i].Value + "\n")
					variable := LookupVariable(evaltokens[i + 1].Value, true, Scope, evaltokens[i + 1], &evaltokens)
					if deref == false {
						Write("mov r5, " + variable.Real, true)
						Write("lodf r5, r5", true)
					} else {
						Write("mov r4, " + variable.Real, true)
						Write("lodf r5, r5", true)
						Write("lodf r5, r5", true)
					}
				default:
					// Error
				}

				print(evaltokens[i].Value + "\n")
				switch evaltokens[i].Type {
				case lexer.TokPlus:
					Write("add r6, r4, r5", true)
				case lexer.TokMinus:
					Write("sub r6, r4, r5", true)
				case lexer.TokStar:
					Write("mul r6, r4, r5", true)
				case lexer.TokSlash:
					Write("div r6, r4, r5", true)
				}
				i++
			}
		}	
		return 0, end
	} 
	return 0, start
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
			bits := BitPref
			for {
				if peek(0).Value[0] == '#' {
					i++
					pp_dir := expect(lexer.TokIdent)	
					switch pp_dir {
					case "include":
						filename := ""
						basename := ""
						global := false
						
						if peek(0).Type != lexer.TokLAngle {
							filename = expect(lexer.TokIdent)
							filename = strings.ReplaceAll(filename, "\"", "")
							basename = filename
							baseDir := filepath.Dir(peek(0).File)
							relPath := filepath.Join(baseDir, filename)
							filename = filepath.Clean(relPath)	
						} else {
							global = true
							expect(lexer.TokLAngle)
							filename = tokens[i].Value
							i++
							filename = filename + tokens[i].Value
							i++
							filename = filename + tokens[i].Value
							i++
							basename = filename
							if runtime.GOOS != "windows" {
								filename = "/usr/local/lib/lcc/" + filename
							} else {
								filename = "C:\\luna\\lcc\\" + filename
							}
							expect(lexer.TokRAngle)
						}
						
						top:
						data, err := os.ReadFile(filename)
						if err != nil {
							if global == false {
								if runtime.GOOS != "windows" {
									filename = "/usr/local/lib/lcc/" + basename
								} else {
									filename = "C:\\luna\\lcc\\" + basename
								}
								global = true
								goto top
							} else {
								error.ErrorNoGaze(16, "'" + filename + "'", peek(-1))
							}
						}
						tokens := lexer.Lex(string(data), filename)	
						Parse(tokens, 1)
					case "pragma":
						directive := expect(lexer.TokIdent)
						switch directive {
						case "__16bit":
							BitPref = 16
						case "__32bit":
							BitPref = 32
						default:
							error.Warning(17, "'" + directive + "'", peek(-1), &tokens)	
						}
					default:
						error.Error(15, "'" + pp_dir + "'", peek(-1), &tokens)	
					}
					continue
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

				if name == "_start" {
					PreWrite("jmp _start", false)
				}

				if len(Children) > 0 {
					level = 1
					if static == false {
						WritePre(".global " + name, false)
					}
					topLevelName = name
					if name != "_start" && noreturn == false {
						Write("pop e11", true)	
					}
					if register > 0 {
						for r := register; r >= 0; r-- {
							Write("pop e" + fmt.Sprintf("%d", r), true)
						}
					}
					if name != "_start" && noreturn == false {
						Write("push e11", true)
					}
					Parse(Children, fscope)
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
					_, end := ParseExpression(tokens, i, Scope)
					var val any	

					if ptr == true {
						rn := "var_" + fmt.Sprintf("%d", IDCounter)
						IDCounter++
						Variables = append(Variables, Variable_Static{Name: name, Type: rtype, Value: val, Pointer: true, Real: rn, Scope: Scope, Const: constant})
						WritePre(rn + ":", false)
						WritePre(".ptr 0", true)

						// Move result to variable
						Write("mov r7, " + rn, true)
						Write("strf r7, r6", true)
					} else {
						rn := "var_" + fmt.Sprintf("%d", IDCounter)
						IDCounter++
						Variables = append(Variables, Variable_Static{Name: name, Type: rtype, Value: val, Pointer: false, Real: rn, Scope: Scope, Const: constant})
						WritePre(rn + ":", false)
						WritePre(".ptr 0", true)

						// Move result to variable
						Write("mov r7, " + rn, true)
						Write("strf r7, r6", true)
					}
					i = end + 1
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
						Variables = append(Variables, Variable_Static{Name: name, Type: STRING, Value: str, Pointer: false, Scope: Scope, Const: constant})
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
						WritePre(".ptr " + fmt.Sprintf("%d", 0), true)
						Variables = append(Variables, Variable_Static{Name: name, Type: rtype, Value: 0, Pointer: true, Real: rn, Scope: Scope, Const: constant})
					} else {
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
			} else {
				error.Error(1, "'" + peek(0).Value + "'", _typetoken, &tokens)
			}
		case 1:	
			// Variable reassignment / function call
			var type_ lexer.TokenType = peek(0).Type
			switch type_ {
			case lexer.TokIdent, lexer.TokStar, lexer.TokAmpersand:
				_name_token := tokens[i]

				deref := false	

				if peek(2).Type != lexer.TokLParen {
					if peek(0).Type == lexer.TokStar {
						expect(lexer.TokStar)
						deref = true
					}
				}

				_ntok := peek(0)
				name := expect(lexer.TokIdent)
				if name == "asm" || name == "__asm__" {
					if peek(0).Value == "volatile" {
						expect(lexer.TokQualifier)
					}
				}

				if peek(0).Type == lexer.TokLParen {
					expect(lexer.TokLParen)
					switch name {
						case "asm", "__asm__":
							str, end := StringParse(tokens, i)
							i = end + 1
							Write(str, true)
							expect(lexer.TokRParen)
							expect(lexer.TokSemi)
						default:
							fvar := LookupVariable(name, false, Scope, peek(-2), &tokens) 
							if fvar.Name == "__ZERO" {
								error.Error(19, "'" + name + "'; ISO C99 and later do not support implicit function declarations", peek(-2), &tokens)	
							}
							var expComma bool = false
							var pushed int = 0
							for j := i; j < len(tokens); j++ {
								if tokens[j].Type == lexer.TokRParen {
									i = j
									break
								} else {
									if expComma == true {
										if tokens[j].Type != lexer.TokComma {
											error.Error(2, "','", tokens[j], &tokens)
										} else {
											expComma = false
											continue
										}
									}
									if strings.HasPrefix(tokens[j].Value, "\"") {
										str, end := StringParse(tokens, j)
										j = end
										CreateStatic(Variable_Static{Name: "var_" + fmt.Sprintf("%d", IDCounter), Type: STRING, Value: str})
										Write("push var_" + fmt.Sprintf("%d", IDCounter), true)
										IDCounter++
										expComma = true
									} else if CheckNum(tokens[j]) == true {
										Write("push " + tokens[j].Value, true)
										expComma = true
									} else {
										variable := LookupVariable(tokens[j].Value, true, Scope, tokens[j], &tokens)	
										if variable.Pointer == false {
											Write("push " + fmt.Sprintf("%v", variable.Value), true)
										} else {
											Write("push " + variable.Real, true)
										}
										expComma = true
									}
									pushed++	
								}
							} 
							if pushed < fvar.ArgNum {
								t, s := FuncDeclLookup(name)
								error.Note(22, "'" + name + "' declared here", t, s)
								error.Error(20, "expected " + fmt.Sprintf("%d", fvar.ArgNum) + ", have " + fmt.Sprintf("%d", pushed), peek(0), &tokens)	
							} else if pushed > fvar.ArgNum {
								t, s := FuncDeclLookup(name)
								error.Note(22, "'" + name + "' declared here", t, s)
								error.Error(21, "expected " + fmt.Sprintf("%d", fvar.ArgNum) + ", have " + fmt.Sprintf("%d", pushed), peek(0), &tokens)	
							}
							expect(lexer.TokRParen)
							expect(lexer.TokSemi)
							Write("call " + name, true)
					}
				} else if peek(0).Type == lexer.TokEqual || peek(0).Type == lexer.TokExclamation || peek(0).Type == lexer.TokLAngle || peek(0).Type == lexer.TokRAngle {
					switch {
						case peek(1).Type == lexer.TokEqual || peek(0).Type == lexer.TokLAngle || peek(0).Type == lexer.TokRAngle:
							not := false
							gt := false
							lt := false
							if peek(0).Type == lexer.TokExclamation {
								not = true
								expect(lexer.TokExclamation)
							} else if peek(0).Type == lexer.TokLAngle {
								lt = true
								expect(lexer.TokLAngle)
							} else if peek(0).Type == lexer.TokRAngle {
								gt = true
								expect(lexer.TokRAngle)
							} else {
								expect(lexer.TokEqual)
							}

							if lt != true && gt != true {
								expect(lexer.TokEqual)	
							}
							// We'll just compare 2 variables for now...

							
							var_1 := LookupVariable(name, true, Scope, _name_token, &tokens)
							_var_2_token := tokens[i]

							deref_2 := false
							if peek(0).Type == lexer.TokStar {
								deref_2 = true
								expect(lexer.TokStar)
							}

							name2 := expect(lexer.TokIdent)	
							var_2 := LookupVariable(name2, true, Scope, _var_2_token, &tokens)
							expect(lexer.TokSemi)

							switch var_1.Type {
							case NUMBER8, NUMBER16, NUMBER32:
								if var_1.Pointer != var_2.Pointer {
									str1 := ""
									str2 := ""
									if var_1.Pointer == true {
										str1 = "int *"
									} else {
										str1 = "int"
									}
									if var_2.Pointer == true {
										str2 = "int *"
									} else {
										str2 = "int"
									}
									error.Warning(8, "('" + str1 + "' and '" + str2 + "')", _var_2_token, &tokens)
								}

								if deref == false {
									Write("mov r1, " + var_1.Real, true)	
									Write("lodf r1, r1", true)	
								} else {		
									Write("mov r1, " + var_1.Real, true)
									if var_1.Type == NUMBER8 {
										Write("lod r1, r1", true)
										Write("lod r1, r1", true)
									} else {
										Write("lodf r1, r1", true)
										Write("lodf r1, r1", true)
									}	
								}

								if deref_2 == false {
									Write("mov r2, " + var_2.Real, true)
									Write("lodf r2, r2", true)	
								} else {		
									Write("mov r2, " + var_2.Real, true)
									if var_1.Type == NUMBER8 {
										Write("lod r2, r2", true)
										Write("lod r2, r2", true)
									} else {
										Write("lodf r2, r2", true)
										Write("lodf r2, r2", true)
									}	
								}
								if not == true {
									Write("cmp e6, r1, r2", true)
									Write("not e6, e6", true)
								} else if lt == true {
									Write("ilt e6, r1, r2", true)
								} else if gt == true {
									Write("igt e6, r1, r2", true)
								} else {
									Write("cmp e6, r1, r2", true)
								}
							}		
						default:
							expect(lexer.TokEqual)
							variable := LookupVariable(name, true, Scope, _name_token, &tokens)	
							switch peek(0).Type {
							case lexer.TokNumber:	
								if deref == false {
									if variable.Type != NUMBER16 && variable.Type != NUMBER32 {
										error.Error(5, "'" + peek(0).Value + "'", peek(0), &tokens)
									}
									_, end := ParseExpression(tokens, i, Scope)
									Write("mov r1, " + name, true)
									Write("mov r2, r6", true)
									Write("str r1, r2", true)
									i = end
								} else {
									if variable.Pointer == false {
										error.Error(26, "", _ntok, &tokens)
									}
									_, end := ParseExpression(tokens, i, Scope)
									Write("mov r1, " + variable.Real, true)
									Write("mov r2, r6", true)

									if variable.Type == NUMBER8 {
										Write("lod r1, r1", true)
										Write("str r1, r2", true)
									} else {
										Write("lodf r1, r1", true)
										Write("strf r1, r2", true)
									}	
									i = end + 1
								}
							case lexer.TokIdent, lexer.TokAmpersand:
								addr := false
								if peek(0).Type == lexer.TokAmpersand {
									expect(lexer.TokAmpersand)
									addr = true
								}	
								name_ := expect(lexer.TokIdent)	
								variable__ := LookupVariable(name_, true, Scope, peek(-1), &tokens)
								if deref == false {
									if addr == false {
										Write("mov r1, " + variable__.Real, true)	
										Write("mov r2, " + variable.Real, true)
										Write("lodf r2, r2", true)
										Write("strf r1, r2", true)
									} else {
										Write("mov r1, " + variable__.Real, true)	
										Write("mov r2, " + variable.Real, true)
										Write("strf r1, r2", true)
									}
								} else {
									if addr == false {
										Write("mov r1, " + variable.Real, true)
										Write("lodf r1, r1", true)
										Write("mov r2, " + variable__.Real, true)
										Write("lodf r2, r2", true)
										Write("strf r1, r2", true)
									} else {
										Write("mov r1, " + variable.Real, true)
										Write("lodf r1, r1", true)
										Write("mov r2, " + variable__.Real, true)	
										Write("strf r1, r2", true)
									}
								}
							}	
							expect(lexer.TokSemi)
					}	
				} else if peek(0).Type == lexer.TokColon {
					expect(lexer.TokColon)
					Write(name + ":", false)
					Variables = append(Variables, Variable_Static{Name: name, Type: POINT, Value: NULL, Scope: 1})
				} else if peek(0).Type == lexer.TokPlus && peek(1).Type == lexer.TokPlus {
					expect(lexer.TokPlus)
					expect(lexer.TokPlus)
					_var := LookupVariable(name, true, Scope, _ntok, &tokens)
					Write("mov r4, " + _var.Real, true)
					Write("lodf r4, r5", true)
					Write("inc r5", true)
					Write("strf r4, r5", true);
					expect(lexer.TokSemi)
				} else {
					expect(lexer.TokSemi)
				}
			case lexer.TokReturn:
				expect(lexer.TokReturn)
				if peek(0).Type == lexer.TokIdent {
					_name_token := tokens[i]
					name := expect(lexer.TokIdent)
					expect(lexer.TokSemi)
					LookupVariable(name, true, Scope, _name_token, &tokens)
					Write("mov e6, " + name, true)
				} else {
					expect(lexer.TokSemi)
				}
				if topLevelName == "_start" {
					Write("pop e11", true)
					Write("ret", true)	
				}
			case lexer.TokSemi:
				expect(lexer.TokSemi)
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
			case lexer.TokGoto:
				expect(lexer.TokGoto)
				_name_token := tokens[i]
				name := expect(lexer.TokIdent)	
				expect(lexer.TokSemi)
				LookupVariable(name, true, Scope, _name_token, &tokens)
				Write("jmp " + name, true)	
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
