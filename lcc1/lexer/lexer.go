package lexer

import (
	"text/scanner"
	"strconv"
	"strings"
	"fmt"
	"unicode"
	"os"
	"path/filepath"
)

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
    s.Init(strings.NewReader(code))
    s.Mode = scanner.ScanIdents | scanner.ScanInts | scanner.ScanChars | scanner.ScanStrings | scanner.SkipComments
    for tok := s.Scan(); tok != scanner.EOF; tok = s.Scan() {	
		content := s.TokenText()
		if content == "int" || content == "void" || content == "char" {
			tokens = append(tokens, Token{Type: TokType, Value: content, Line: s.Pos().Line, File: filename})
		} else if content == "volatile" || content == "unsigned" || content == "long" || content == "short" || content == "static" || content == "const" || content == "extern" {
			tokens = append(tokens, Token{Type: TokQualifier, Value: content, Line: s.Pos().Line, File: filename})
		} else if content == "return" {
			tokens = append(tokens, Token{Type: TokReturn, Value: content, Line: s.Pos().Line, File: filename})
		} else if content == "if" {
			tokens = append(tokens, Token{Type: TokIf, Value: content, Line: s.Pos().Line, File: filename})
		} else if content == "else" {
			tokens = append(tokens, Token{Type: TokElse, Value: content, Line: s.Pos().Line, File: filename})
		} else if num, err := strconv.ParseInt(content, 0, 64); err == nil {
			tokens = append(tokens, Token{Type: TokNumber, Value: fmt.Sprintf("%d", num), Line: s.Pos().Line, File: filename})
		} else if content == "(" {
			tokens = append(tokens, Token{Type: TokLParen, Value: content, Line: s.Pos().Line, File: filename})
		} else if content == ")" {
			tokens = append(tokens, Token{Type: TokRParen, Value: content, Line: s.Pos().Line, File: filename})
		} else if content == "{" {
			tokens = append(tokens, Token{Type: TokLCurly, Value: content, Line: s.Pos().Line, File: filename})
		} else if content == "}" {
			tokens = append(tokens, Token{Type: TokRCurly, Value: content, Line: s.Pos().Line, File: filename})
		} else if content == ";" {
			tokens = append(tokens, Token{Type: TokSemi, Value: content, Line: s.Pos().Line, File: filename})
		} else if content == "+" {
			tokens = append(tokens, Token{Type: TokPlus, Value: content, Line: s.Pos().Line, File: filename})
		} else if content == "-" {
			tokens = append(tokens, Token{Type: TokMinus, Value: content, Line: s.Pos().Line, File: filename})
		} else if content == "*" {
			tokens = append(tokens, Token{Type: TokStar, Value: content, Line: s.Pos().Line, File: filename})
		} else if content == "/" {
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
				tokens = append(tokens, Token{Type: TokSlash, Value: content, Line: s.Pos().Line, File: filename})
			}	
		} else if content == "=" {
			tokens = append(tokens, Token{Type: TokEqual, Value: content, Line: s.Pos().Line, File: filename})
		} else if content == "," {
			tokens = append(tokens, Token{Type: TokComma, Value: content, Line: s.Pos().Line, File: filename})
		} else if content == ":" {
			tokens = append(tokens, Token{Type: TokColon, Value: content, Line: s.Pos().Line, File: filename})
		} else if content == "goto" {
			tokens = append(tokens, Token{Type: TokGoto, Value: content, Line: s.Pos().Line, File: filename})
		} else if content == "for" {
			tokens = append(tokens, Token{Type: TokFor, Value: content, Line: s.Pos().Line, File: filename})
		} else if content == "while" {
			tokens = append(tokens, Token{Type: TokWhile, Value: content, Line: s.Pos().Line, File: filename})
		} else if content == "do" {
			tokens = append(tokens, Token{Type: TokDo, Value: content, Line: s.Pos().Line, File: filename})
		} else if content == "<" {
			tokens = append(tokens, Token{Type: TokLAngle, Value: content, Line: s.Pos().Line, File: filename})
		} else if content == ">" {
			tokens = append(tokens, Token{Type: TokRAngle, Value: content, Line: s.Pos().Line, File: filename})
		} else if content == "&" {
			tokens = append(tokens, Token{Type: TokAmpersand, Value: content, Line: s.Pos().Line, File: filename})
		} else if content == "!" {
			tokens = append(tokens, Token{Type: TokExclamation, Value: content, Line: s.Pos().Line, File: filename})
		} else if content == "//" {

		} else {
			tokens = append(tokens, Token{Type: TokIdent, Value: content, Line: s.Pos().Line, File: filename})
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
		case "#ifdef":
			alias := tokens[i + 1].Value
			i += 2
			if _, ok := defines[alias]; ok {
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
			fmt.Println("\033[1;39m" + filename + ":" + fmt.Sprintf("%d", tokens[i - 1].Line) + "\033[0m \033[1;31merror:\033[0m \033[1;39m" + tokens[i].Value + "\033[0m")
			os.Exit(1)
		case "#include":
			base := filepath.Dir(filename)
			i++
			again = true
			raw := strings.ReplaceAll(tokens[i].Value, "\"", "")
			path := raw

			if !filepath.IsAbs(raw) {
				path = filepath.Join(base, raw)
			}

			contents, err := os.ReadFile(path)
			if err != nil {
				fmt.Println("\033[1;39m" + filename + ":" + fmt.Sprintf("%d", tokens[i - 1].Line) + ":\033[0m \033[1;31merror:\033[0m \033[1;39mno such file or directory '" + path + "'\033[0m")
				os.Exit(1)
			}
			i++
			ntokens := tokenize(string(contents)) 
			for j := 0; j < len(ntokens); j++ {
				out = append(out, ntokens[j])
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

	fmt.Println(out_text)
	return out_text 
}

