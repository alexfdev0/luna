package lexer

import (
	"text/scanner"
	"strconv"
	"strings"
	"fmt"
	"unicode"
	"os"
	"path/filepath"
	"runtime"
)

var Bits int = 16
var Warnings = 0

type TokenType int

const (
	TokType TokenType = iota
	TokReturn
	TokIf
	TokElse
	TokIdent
	TokNumber
	TokLParen
	TokRParen
	TokLCurly
	TokRCurly	
	TokSemi
	TokPlus
	TokMinus
	TokStar
	TokSlash
	TokEqual
	TokComma
	TokEOF
	TokColon
	TokGoto
	TokQualifier
	TokFor
	TokWhile
	TokDo
	TokLAngle
	TokRAngle
	TokAmpersand
	TokExclamation
	TokLBracket
	TokRBracket
	TokTypedef
	TokEquality
	TokInequality
	TokGEqual
	TokLEqual
	TokBreak
	TokContinue
	TokStruct
)

type Token struct {
	Type TokenType
	Value string
	Line int
	File string
}

func contains(set string, c byte) bool {
    for i := 0; i < len(set); i++ {
        if set[i] == c {
            return true
        }
    }
    return false
}

func Lex(code string, filename string) []Token {
	var tokens = []Token {}
	var s scanner.Scanner

	Add := func(Type TokenType, Value string) {
		tokens = append(tokens, Token{
			Type: Type,
			Value: Value,
			Line: s.Pos().Line,
			File: filename,
		})
	}

	
    s.Init(strings.NewReader(code))
    s.Mode = scanner.ScanIdents | scanner.ScanInts | scanner.ScanChars | scanner.ScanStrings | scanner.SkipComments
    for tok := s.Scan(); tok != scanner.EOF; tok = s.Scan() {	
		content := s.TokenText()

		switch content {	
		case "int", "void", "char":
			Add(TokType, content)
		case "volatile", "unsigned", "short", "long", "static", "const", "extern":
			Add(TokQualifier, content)
		case "return":
			Add(TokReturn, content)
		case "if":
			Add(TokIf, content)
		case "else":
			Add(TokElse, content)
		case "break":
			Add(TokBreak, content)
		case "continue":
			Add(TokContinue, content)
		case "(":
			Add(TokLParen, content)
		case ")":
			Add(TokRParen, content)
		case "{":
			Add(TokLCurly, content)
		case "}":
			Add(TokRCurly, content)
		case ";":
			Add(TokSemi, content)
		case "+":
			Add(TokPlus, content)
		case "-":
			Add(TokMinus, content)
		case "*":
			Add(TokStar, content)
		case "/":
			next := s.Peek()
			if next == '/' {
				for {
					r := s.Next()
					if r == '\n' || r == scanner.EOF {
						break
					}
				}
				continue
			} else if next == '*' {
				s.Next()
				for {
					r := s.Next()
					if r == scanner.EOF {
						break
					}
					if r == '*' && s.Peek() == '/' {
						s.Next()
						break
					}
				}
				continue
			} else {
				Add(TokSlash, content)	
			}
		case "=":
			if s.Peek() == '=' {
				s.Next()
				Add(TokEquality, "==")
			} else {
				Add(TokEqual, content)
			}
		case ",":
			Add(TokComma, content)
		case ":":
			Add(TokColon, content)
		case "goto":
			Add(TokGoto, content)
		case "for":
			Add(TokFor, content)
		case "while":
			Add(TokWhile, content)
		case "do":
			Add(TokDo, content)
		case "<":
			if s.Peek() == '=' {
				s.Next()
				Add(TokLEqual, "<=")
			} else {
				Add(TokLAngle, content)
			}
		case ">":
			if s.Peek() == '=' {
				s.Next()
				Add(TokGEqual, ">=")
			} else {
				Add(TokRAngle, content)
			}
		case "&":
			Add(TokAmpersand, content)
		case "!":
			if s.Peek() == '=' {
				s.Next()
				Add(TokInequality, "!=")
			} else {
				Add(TokExclamation, content)
			}
		case "//":
		case "[":
			Add(TokLBracket, content)
		case "]":
			Add(TokRBracket, content)
		case "typedef":
		case "struct":
		default:
			num, err := strconv.ParseInt(content, 0, 64)
			if err == nil {
				Add(TokNumber, fmt.Sprintf("%d", num))
			} else {
				Add(TokIdent, content)
			}
		}	
	}

	return tokens
}

type SmallToken struct {
	Value string
	Line int
}

func Preprocessor(text string, filename string, just_split bool) string {
	var out []SmallToken
	defines := make(map[string][]SmallToken)
	defines["__LCC__"] = []SmallToken{SmallToken{Value: "1", Line: 0}}
	again := false	

	tokenize := func(text string) []SmallToken {
		currentLine := 1
		inString := false
		var tokens []SmallToken
		var buf []rune
		for i, r := range text {
			switch {
			case r == '"':
				buf = append(buf, r)
				if inString { 
					tokens = append(tokens, SmallToken{Value: string(buf), Line: currentLine})
					buf = buf[:0]
				}
				inString = !inString	
			case r == '\n':
				if len(buf) > 0 {
					tokens = append(tokens, SmallToken{Value: string(buf), Line: currentLine})
					buf = buf[:0]
				}
				tokens = append(tokens, SmallToken{Value: "\n", Line: currentLine})
				currentLine++
			case unicode.IsSpace(r) && inString == false:
				if len(buf) > 0 {
					tokens = append(tokens, SmallToken{Value: string(buf), Line: currentLine})
					buf = buf[:0]
				}
			default:
				buf = append(buf, r)
			}
			if i == len(text) - 1 && len(buf) > 0 {
				tokens = append(tokens, SmallToken{Value: string(buf), Line: currentLine})
			}
		}
		return tokens
	}

	tokens := tokenize(text)	

	for i := 0; i < len(tokens); i++ {
		switch tokens[i].Value {
		case "#define":
			alias := tokens[i + 1].Value
			actual := []SmallToken {}
			i += 2
			
			for j := i; j < len(tokens); j++ {
				i++
				if tokens[j].Value == "\n" {
					break
				}	
				actual = append(actual, tokens[j]) 
			}

			defines[alias] = actual
		case "#ifdef", "#ifndef":
			alias := tokens[i + 1].Value
			i += 2
			if _, ok := defines[alias]; (ok && tokens[i - 2].Value == "#ifdef") || (!ok && tokens[i - 2].Value == "#ifndef") {
				for j := i; j < len(tokens); j++ {
					if tokens[j].Value != "#else" && tokens[j].Value != "#endif" {
						again = true
						out = append(out, tokens[j])
						i++
					} else {
						if tokens[j].Value == "#endif" {
							i++
						} else {
							for k := j; k < len(tokens); k++ {
								i++
								if tokens[k].Value == "#endif" {
									break
								}
							}
						}
						break
					}
				}
					
			} else {
				for j := i; j < len(tokens); j++ {
					i++
					if tokens[j].Value == "#endif" {
						break
					} else if tokens[j].Value == "#else" {
						for k := j + 1; k < len(tokens); k++ {
							i++
							if tokens[k].Value != "#endif" {
								again = true
								out = append(out, tokens[k])
							} else {
								break
							}
						}
					}
				}
			}
		case "#error":
			i++
			fmt.Println("\033[1;39m" + filename + ":" + fmt.Sprintf("%d", tokens[i - 1].Line) + ":\033[0m \033[1;31merror:\033[0m \033[1;39m" + tokens[i].Value + "\033[0m")
			os.Exit(1)
		case "#warning":
			i++
			Warnings++
			fmt.Println("\033[1;39m" + filename + ":" + fmt.Sprintf("%d", tokens[i - 1].Line) + ":\033[0m \033[1;35mwarning:\033[0m \033[1;39m" + tokens[i].Value + "\033[0m")	
		case "#include":
			base := filepath.Dir(filename)
			i++
			again = true
			raw := strings.ReplaceAll(tokens[i].Value, "\"", "")
			path := raw

			if !filepath.IsAbs(raw) {
				path = filepath.Join(base, raw)
			}

			ogpath := path
			times := 0

			ff_top:
			contents, err := os.ReadFile(path)
			if err != nil {
				times++
				if times >= 2 {
					fmt.Println("\033[1;39m" + filename + ":" + fmt.Sprintf("%d", tokens[i - 1].Line) + ":\033[0m \033[1;31mfatal error:\033[0m \033[1;39mno such file or directory '" + ogpath + "'\033[0m")
					os.Exit(1)
				} else {
					switch runtime.GOOS {
					case "windows":
						path = "C:\\Program Files (x86)\\Luna L2\\lib\\lcc\\" + raw
					default:
						path = "/usr/local/lib/lcc/" + raw
					} 
					goto ff_top
				}
			}
			i++

			ntokens := tokenize(string(contents)) 
			for j := 0; j < len(ntokens); j++ {
				out = append(out, ntokens[j])
			}
		case "#pragma":
			i++
			switch tokens[i].Value {
			case "bits":
				i++
				switch tokens[i].Value {
				case "16":
					Bits = 16
				case "32":
					Bits = 32
				default:
					fmt.Println("\033[1;39m" + filename + ":" + fmt.Sprintf("%d", tokens[i - 1].Line) + ":\033[0m \033[1;35mwarning:\033[0m \033[1;39minvalid number for '#pragma bits'\033[0m")
					Warnings++	
				}
			}
		default:
			token := tokens[i]
			if tokens[i].Value[0] == '#' {
				fmt.Println("\033[1;39m" + filename + ":" + fmt.Sprintf("%d", tokens[i - 1].Line) + ":\033[0m \033[1;31merror:\033[0m \033[1;39minvalid preprocessor directive '" + tokens[i].Value + "'\033[0m")
				os.Exit(1)
			}

			if repl, ok := defines[token.Value]; ok {
				for j := 0; j < len(repl); j++ {
					out = append(out, repl[j])
				}
				i++
			} else {
				out = append(out, token)
			}	
		}	
	}	

	out_text := ""
	for i := 0; i < len(out); i++ {
		out_text = out_text + out[i].Value + " "
	}

	if again == true {
		return Preprocessor(out_text, filename, false)
	}

	return out_text 
}

