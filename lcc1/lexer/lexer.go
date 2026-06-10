package lexer

import (
	"strconv"
	"strings"
	"fmt"
	"unicode"
	"os"
	"path/filepath"
	"runtime"
	"lcc1/shared"
	"lcc1/error"
)

type SmallToken struct {
	Value string
	Line int
	Filename string
}

func Lex(code []SmallToken, filename string) []shared.Token {
	var tokens = []shared.Token {}

	Add := func(Type shared.TokenType, Value string, ST SmallToken) {
		tokens = append(tokens, shared.Token{
			Type: Type,
			Value: Value,
			Line: ST.Line,
			File: ST.Filename,
		})
	}
   	 
	for i := 0; i < len(code); i++ {
		SToken := code[i]
		content := SToken.Value

		switch content {	
		case "int", "void", "char":
			Add(shared.TokType, content, SToken)
		case "volatile", "unsigned", "short", "long", "static", "const", "extern":
			Add(shared.TokQualifier, content, SToken)
		case "return":
			Add(shared.TokReturn, content, SToken)
		case "if":
			Add(shared.TokIf, content, SToken)
		case "else":
			Add(shared.TokElse, content, SToken)
		case "break":
			Add(shared.TokBreak, content, SToken)
		case "continue":
			Add(shared.TokContinue, content, SToken)
		case "(":
			Add(shared.TokLParen, content, SToken)
		case ")":
			Add(shared.TokRParen, content, SToken)
		case "{":
			Add(shared.TokLCurly, content, SToken)
		case "}":
			Add(shared.TokRCurly, content, SToken)
		case ";":
			Add(shared.TokSemi, content, SToken)
		case "+":
			Add(shared.TokPlus, content, SToken)
		case "-":
			Add(shared.TokMinus, content, SToken)
		case "*":
			Add(shared.TokStar, content, SToken)
		case "/":	
			Add(shared.TokSlash, content, SToken)	
		case "=":
			Add(shared.TokEqual, content, SToken)
		case "==":
			Add(shared.TokEquality, content, SToken)
		case ",":
			Add(shared.TokComma, content, SToken)
		case ":":
			Add(shared.TokColon, content, SToken)
		case "goto":
			Add(shared.TokGoto, content, SToken)
		case "for":
			Add(shared.TokFor, content, SToken)
		case "while":
			Add(shared.TokWhile, content, SToken)
		case "do":
			Add(shared.TokDo, content, SToken)
		case "<":
			Add(shared.TokLAngle, content, SToken)
		case "<=":
			Add(shared.TokLEqual, content, SToken)
		case ">":
			Add(shared.TokRAngle, content, SToken)
		case ">=":
			Add(shared.TokGEqual, content, SToken)
		case "&":
			Add(shared.TokAmpersand, content, SToken)
		case "!":
			Add(shared.TokExclamation, content, SToken)
		case "!=":
			Add(shared.TokInequality, content, SToken)
		case "//":
		case "\n":
		case "[":
			Add(shared.TokLBracket, content, SToken)
		case "]":
			Add(shared.TokRBracket, content, SToken)
		case "typedef":
			Add(shared.TokTypedef, content, SToken)
		case "struct":
			Add(shared.TokStruct, content, SToken)
		default:
			num, err := strconv.ParseInt(content, 0, 64)
			if err == nil {
				Add(shared.TokNumber, fmt.Sprintf("%d", num), SToken)
			} else {
				Add(shared.TokIdent, content, SToken)
			}
		}	
	}

	return tokens
}

func Tokenize (text string, filename string) []SmallToken {
	var tokens []SmallToken
	runes := []rune(text)
	i := 0
	Line := 1

	peek := func(offset int) rune {
		j := i + offset
		if j < len(runes) {
			return runes[j]
		}
		return 0
	}

	emit := func(Value string) {
		tokens = append(tokens, SmallToken{Value: Value, Line: Line, Filename: filename})
	}

	for i < len(runes) {
		r := runes[i]

		if r == '/' && peek(1) == '/' {
			i += 2
			for i < len(runes) && runes[i] != '\n' {
				i++
			}
			continue
		}

		if r == '/' && peek(1) == '*' {
			i += 2
			for i < len(runes) {
				if runes[i] == '*' && peek(1) == '/' {
					i += 2
					break
				}
				if runes[i] == '\n' {
					Line++
				}
				i++
			}
			continue
		}

		if r == '\n' {
			emit("\n")
			Line++
			i++
			continue
		}

		if unicode.IsSpace(r) {
			i++
			continue
		}

		if r == '"' || (r == '<' && tokens[len(tokens) - 1].Value == "#include")  {
			var termonrangle bool
			var buf []rune

			if r == '<' && tokens[len(tokens) - 1].Value == "#include" {
				termonrangle = true
			}

			buf = append(buf, r)
			i++
			for i < len(runes) {
				c := runes[i]
				buf = append(buf, c)
				if c == '\\' && peek(1) != 0 {
					i++
					buf = append(buf, runes[i])
				} else if c == '"' || (c == '>' && termonrangle == true) {
					break
				} else if c == '\n' {
					Line++
				}
				i++
			}
			i++

			emit(string(buf))
			continue
		}

		// ==
		// >=
		// <=
		// !=

		if r == '=' && peek(1) == '=' {
			emit("==")
			i += 2
			continue
		}

		if r == '>' && peek(1) == '=' {
			emit(">=")
			i += 2
			continue
		}

		if r == '<' && peek(1) == '=' {
			emit("<=")
			i += 2
			continue
		}

		if r == '!' && peek(1) == '=' {
			emit("!=")
			i += 2
			continue
		}

		if strings.ContainsRune("+-*/%&|^~<>=!?:;.,()[]{}@", r) {
			emit(string(r))
			i++
			continue
		}

		var buf []rune
		for i < len(runes) {
			c := runes[i]
			if unicode.IsSpace(c) || strings.ContainsRune("+-*/%&|^~<>=!?:;.,()[]{}@\"", c) {
				break
			}
			buf = append(buf, c)
			i++
		}
		if len(buf) > 0 {
			emit(string(buf))
		}
	}

	return tokens
}

type DefineEntry struct {
	Name string
	ReplacementList []SmallToken
}

var Defines = []DefineEntry {
	DefineEntry { 
		Name: "__LCC__", 
		ReplacementList: []SmallToken { 
			SmallToken { 
				Value: "1", 
				Filename: "__DEFAULT__",
			}, 
		},
	},
}

func CheckDefined(Name string) (DefineEntry, bool) {
	for _, Entry := range Defines {
		if Entry.Name == Name {
			return Entry, true
		}
	}
	return DefineEntry{Name: "__NOTFOUND__"}, false
}

func Preprocessor(text string, filename string) []SmallToken {
	var out []SmallToken

	

	again := false

	var pragma_once_files []string
	var included_files []string

	CheckPragmaOnce := func(filename string) bool {
		for _, file := range pragma_once_files {
			if file == filename {
				return true
			}
		}
		return false
	}

	tokens := Tokenize(text, filename)
PREPROCESSOR_TOP:
	for i := 0; i < len(tokens); i++ {
		switch tokens[i].Value {
		case "#define":
			var Define DefineEntry
			i++

			Define.Name = tokens[i].Value
			i++

			for j := i; j < len(tokens); j++ {
				t := tokens[j]
				if t.Value == "\n" {
					i = j
					break
				}
				Define.ReplacementList = append(Define.ReplacementList, t)
			}
			Defines = append(Defines, Define)
			again = true
		case "#ifdef", "#ifndef":
			var FakeStream []shared.Token
			for j := i; j < len(tokens); j++ {
				t := tokens[j]
				if t.Value == "\n" {
					break
				}
				FakeStream = append(FakeStream, shared.Token{
					Type: shared.TokIdent,
					Value: t.Value,
					Line: t.Line,
					File: t.Filename,
				})
			}
			origin := FakeStream[0]

			dir := tokens[i].Value
			i++
			look_for := tokens[i].Value
			_, Found := CheckDefined(look_for)
			i++
			if (Found == true && dir == "#ifdef") || (Found == false && dir == "#ifndef") {
				end := 0
				for j := i; j < len(tokens); j++ {
					t := tokens[j]
					if t.Value == "#endif" {
						i = j
						end = j
						break
					}
					out = append(out, t)
				}
				if end == 0 {
					error.Error(29, "", origin, &FakeStream)
				}
				again = true
			} else {
				// TODO: add #else
				end := 0
				for j := i; j < len(tokens); j++ {
					t := tokens[j]
					if t.Value == "#endif" {
						i = j
						end = j
						break
					}
				}
				if end == 0 {
					error.Error(29, "", origin, &FakeStream)
				}
			}
		case "#error":
			var FakeStream []shared.Token
			for j := i; j < len(tokens); j++ {
				t := tokens[j]
				if t.Value == "\n" {
					break
				}
				FakeStream = append(FakeStream, shared.Token{
					Type: shared.TokIdent,
					Value: t.Value,
					Line: t.Line,
					File: t.Filename,
				})
			}
			origin := FakeStream[0]
			i++
			error.Error(22, tokens[i].Value, origin, &FakeStream)
		case "#warning":
			var FakeStream []shared.Token
			for j := i; j < len(tokens); j++ {
				t := tokens[j]
				if t.Value == "\n" {
					break
				}
				FakeStream = append(FakeStream, shared.Token{
					Type: shared.TokIdent,
					Value: t.Value,
					Line: t.Line,
					File: t.Filename,
				})
			}
			origin := FakeStream[0]
			i++
			error.Warning(22, tokens[i].Value, origin, &FakeStream)	
		case "#include":
			again = true
			i++

			raw := ""
			times := 0
			root_by_default := false
			base := filepath.Dir(filename)
			if tokens[i].Value[0] == '"' {
				raw = strings.ReplaceAll(tokens[i].Value, "\"", "")
			} else if tokens[i].Value[0] == '<' {
				raw = strings.ReplaceAll(tokens[i].Value, "<", "")
				raw = strings.ReplaceAll(raw, ">", "")
				root_by_default = true
				times = 1
			}

			path := raw
			if !filepath.IsAbs(raw) {
				path = filepath.Join(base, raw)
			}
			ogpath := path

			if root_by_default == true {
				switch runtime.GOOS {
				case "windows":
					path = "C:\\Program Files (x86)\\Luna L2\\lib\\lcc\\" + raw
				default:
					path = "/usr/local/lib/lcc/" + raw
				}
				ogpath = path
			}

			ff_top:
			if CheckPragmaOnce(path) == true {
				continue
			}

			contents, err := os.ReadFile(path)
			if err != nil {
				times++
				if times >= 2 {
					var FakeStream []shared.Token
					for j := i - 1; j < len(tokens); j++ {
						t := tokens[j]
						if t.Value == "\n" {
							break
						}
						FakeStream = append(FakeStream, shared.Token{
							Type: shared.TokIdent,
							Value: t.Value,
							Line: t.Line,
							File: t.Filename,
						})
					}
					origin := FakeStream[1]
					error.Error(16, "\"" + ogpath + "\"", origin, &FakeStream)
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

			ntokens := Tokenize(string(contents), path) 
			for j := 0; j < len(ntokens); j++ {
				out = append(out, ntokens[j])
			}

			included_files = append(included_files, path)
		case "#pragma":
			i++
			switch tokens[i].Value {
			case "bits":
				i++
				switch tokens[i].Value {
				case "16":
					shared.Bits = 16
				case "32":
					shared.Bits = 32
				default:
					var FakeStream []shared.Token
					for j := i - 2; j < len(tokens); j++ {
						t := tokens[j]
						if t.Value == "\n" {
							break
						}
						FakeStream = append(FakeStream, shared.Token{
							Type: shared.TokIdent,
							Value: t.Value,
							Line: t.Line,
							File: t.Filename,
						})
					}
					origin := FakeStream[2]
					error.Warning(43, "", origin, &FakeStream)
				}
			case "once":
				file := tokens[i - 1].Filename
				pragma_once_files = append(pragma_once_files, file)
			}

			for j := i; j < len(tokens); j++ {
				if tokens[j].Value == "\n" {
					break
				}
				i++
			}
		default:
			token := tokens[i]
			if tokens[i].Value[0] == '#' {
				var FakeStream []shared.Token
				for j := i; j < len(tokens); j++ {
					t := tokens[j]
					if t.Value == "\n" {
						break
					}
					FakeStream = append(FakeStream, shared.Token{
						Type: shared.TokIdent,
						Value: t.Value,
						Line: t.Line,
						File: t.Filename,
					})
				}
				origin := FakeStream[0]
				error.Error(42, "\"" + token.Value + "\"", origin, &FakeStream)
			} else {
				Entry, Found := CheckDefined(token.Value)
				if Found == true {
					for _, Token := range Entry.ReplacementList {
						Token.Line = token.Line
						Token.Filename = token.Filename
						out = append(out, Token)
					}
				} else {
					out = append(out, token)
				}
			}
		}	
	}	

	if again == true {
		tokens = out
		out = []SmallToken {}
		again = false
		goto PREPROCESSOR_TOP
	}

	return out 
}

