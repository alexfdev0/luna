package parser

import (
	"lcc1/lexer"
	"lcc1/error"
	"strings"
	"fmt"
	"github.com/Knetic/govaluate"
	"strconv"
)

var level int = 0
var Code1 string = ""
var Code2 string = ""


var IDCounter = 1

const (
	NUMBER int = iota
	STRING
	POINT
	NULL
)
type Variable_Static struct {
	Name string
	Type int
	Value any
	Pointer bool
}

var Variables = []Variable_Static {}

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

func LookupVariable(Name string, Enforce bool) Variable_Static {
	for _, variable := range Variables {
		if variable.Name == Name {
			return variable
		}
	}
	if Enforce == true {
		error.Error(4, "'" + Name + "'")
		return Variable_Static{Name: "__ZERO", Type: NULL, Value: 0}
	} else {
		return Variable_Static{Name: "__ZERO", Type: NULL, Value: 0}
	}
}

func StringParse(tokens []lexer.Token, start int) (string, int) {
	// Start would be the first token
	var str string = ""
	var loc int = 0
	if strings.HasPrefix(tokens[start].Value, "\"") == false {
		error.Error(2, "\"")		
	}
	if strings.HasSuffix(tokens[start].Value, "\"") {
		tokens[start].Value = strings.Trim(tokens[start].Value, "\"")
		str = tokens[start].Value
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
		str = strings.Trim(str,  "\"")
		loc = start
	}
	
	return str, loc
}

func ParseExpression(tokens []lexer.Token, start int) (int, int) {
	if tokens[start].Type == lexer.TokNumber || LookupVariable(tokens[start].Value, false).Type == NUMBER {
		var evaltokens []string
		var end int = 0
		for i := start; i < len(tokens); i++ {
			if tokens[i].Type == lexer.TokSemi {
				end = i - 1
				break
			}
			evaltokens = append(evaltokens, tokens[i].Value)
		}
		evalstr := strings.Join(evaltokens, " ")	
		expression, err := govaluate.NewEvaluableExpression(evalstr)
		if err != nil {
			error.Error(6, "'" + evalstr + "'")
		}
		result, err := expression.Evaluate(nil)
		if err != nil {
			error.Error(6, "'" + evalstr + "'")
		}

		return int(result.(float64)), end
	} 
	return 0, start
}

func Parse(tokens []lexer.Token) {
	i := 0
	expect := func(toktype lexer.TokenType) string {
		var value string
		if i >= len(tokens) {
			error.Error(1, "'<EOF>'")
		}
		if tokens[i].Type == toktype {
			value = tokens[i].Value
			i++
		} else {
			error.Error(1, "'" + tokens[i].Value + "'")
			return ""
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
			_type := expect(lexer.TokType)
			if peek(0).Type == lexer.TokStar {
				ptr = true
				i++
			}

			name := expect(lexer.TokIdent)
			
			var rtype int	
			switch _type {
			case "int":
				rtype = NUMBER
			case "char":
				rtype = STRING
			}	

			if LookupVariable(name, false).Name != "__ZERO" {
				// print(LookupVariable(name, false).Name)
				error.Error(3, "'" + name + "'")
			}	

			if peek(0).Type == lexer.TokLParen {
				Variables = append(Variables, Variable_Static{Name: name, Type: rtype, Value: nil})
				expect(lexer.TokLParen)
				expect(lexer.TokRParen)
				expect(lexer.TokLCurly)

				var Children = []lexer.Token {}
				ending := -1
				for j := i; j < len(tokens); j++ {
					if tokens[j].Type == lexer.TokRCurly {
						ending = j
						break
					} else {	
						Children = append(Children, tokens[j])
					}
				}
				if ending == -1 {
					error.Error(2, "'}'")
				} else {
					i = ending
				}
			
				expect(lexer.TokRCurly)

				if name == "main" {
					name = "_start"
				}
				Write(name + ":", false)	
				if len(Children) > 0 {
					level = 1
					Parse(Children)
					level = 0
				}	
				i++
				if i < len(tokens) {
					print(tokens[i].Value)
					Parse(tokens)
				}
			} else if peek(0).Type == lexer.TokEqual {
				expect(lexer.TokEqual)	
				switch _type {
				case "void":
					error.Error(7, "'void'")
				case "int":	
					value, end := ParseExpression(tokens, i)		
					if ptr == true {
						Variables = append(Variables, Variable_Static{Name: name, Type: NUMBER, Value: int32(value), Pointer: true})
						Write(name + ":", false)
						Write(".word " + fmt.Sprintf("%d", int32(value)), true)
					} else {
						Variables = append(Variables, Variable_Static{Name: name, Type: NUMBER, Value: int32(value), Pointer: false})
					}
					i = end + 1
				case "char":
					str, end := StringParse(tokens, i)	
					if ptr == true {
						Variables = append(Variables, Variable_Static{Name: name, Type: STRING, Value: str, Pointer: true})
						Write(name + ":", false)
						Write(".asciz \"" + str + "\"", true)
					} else {
						if len(str) > 1 {
							error.Error(5, "'char' with an expression of type 'char*'")
						}
						Variables = append(Variables, Variable_Static{Name: name, Type: STRING, Value: str, Pointer: false})
					}
					i = end + 1
				}
				expect(lexer.TokSemi)
			} else {
				error.Error(1, "'" + peek(0).Value + "'")
			}
		case 1:	
			// Variable reassignment / function call
			var type_ lexer.TokenType = peek(0).Type
			switch type_ {
			case lexer.TokIdent:
				name := expect(lexer.TokIdent)
				if peek(0).Type == lexer.TokLParen {
					expect(lexer.TokLParen)
					var expComma bool = false
					for j := i; j < len(tokens); j++ {
						if tokens[j].Type == lexer.TokRParen {
							i = j
							break
						} else {
							if expComma == true {
								if tokens[j].Type != lexer.TokComma {
									error.Error(2, "','")
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
								variable := LookupVariable(tokens[j].Value, true)	
								if variable.Pointer == false {
									Write("push " + fmt.Sprintf("%v", variable.Value), true)
								} else {
									Write("push " + variable.Name, true)
								}
								expComma = true
							}
						}
					}

					expect(lexer.TokRParen)
					expect(lexer.TokSemi)
					Write("call " + name, true)
				} else if peek(0).Type == lexer.TokEqual {
					switch peek(1).Type {
						case lexer.TokEqual:
							expect(lexer.TokEqual)
							expect(lexer.TokEqual)
							// We'll just compare 2 variables for now...
							LookupVariable(name, true)
							name2 := expect(lexer.TokIdent)
							LookupVariable(name2, true)	
							expect(lexer.TokSemi)
							Write("mov r1, " + name, true)
							Write("lodf r1, r1", true)
							Write("mov r2, " + name2, true)
							Write("lodf r2, r2", true)
							Write("cmp r3, r1, r2", true)
						default:
							expect(lexer.TokEqual)
							variable := LookupVariable(name, true)
							switch peek(0).Type {
							case lexer.TokNumber:
								if variable.Type != NUMBER {
									error.Error(5, "'" + peek(0).Value + "'")
								}
								value, end := ParseExpression(tokens, i)
								Write("mov r1, " + name, true)
								Write("mov r2, " + fmt.Sprintf("%d", value), true)
								Write("str r1, r2", true)
								i = end
							}
							expect(lexer.TokSemi)
					}	
				} else if peek(0).Type == lexer.TokColon {
					expect(lexer.TokColon)
					Write(name + ":", false)
					Variables = append(Variables, Variable_Static{Name: name, Type: POINT, Value: NULL})
				} else {
					expect(lexer.TokSemi)
				}
			case lexer.TokReturn:
				expect(lexer.TokReturn)
				if peek(0).Type == lexer.TokIdent {
					name := expect(lexer.TokIdent)
					expect(lexer.TokSemi)
					LookupVariable(name, true)
					Write("mov t7, " + name, true)
				} else {
					expect(lexer.TokSemi)
				}
				Write("ret", true)	
			case lexer.TokSemi:
				expect(lexer.TokSemi)
			case lexer.TokIf:
				expect(lexer.TokIf)
				expect(lexer.TokLParen)
				var exptokens = []lexer.Token {}
				for j := i; j < len(tokens); j++ {
					if tokens[j].Type == lexer.TokRParen {
						i = j
						break
					}
					exptokens = append(exptokens, tokens[j])
				}
				exptokens = append(exptokens, lexer.Token{Type: lexer.TokSemi, Value: ";"})
				expect(lexer.TokRParen)
				Parse(exptokens)
			case lexer.TokGoto:
				expect(lexer.TokGoto)
				name := expect(lexer.TokIdent)	
				expect(lexer.TokSemi)
				LookupVariable(name, true)
				Write("jmp " + name, true)
			default:
				error.Error(1, "'" + tokens[i].Value + "'")
			}	
		}
	}
} 
