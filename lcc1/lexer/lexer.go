package lexer

import (
	"text/scanner"
	"strconv"
	"strings"
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
	TokLBrack
	TokRBrack	
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
)

type Token struct {
	Type TokenType
	Value string
	Line int
}

func contains(set string, c byte) bool {
    for i := 0; i < len(set); i++ {
        if set[i] == c {
            return true
        }
    }
    return false
}

func Lex(code string) []Token {
	var tokens = []Token {}
	var s scanner.Scanner
    s.Init(strings.NewReader(code))
    s.Mode = scanner.ScanIdents | scanner.ScanInts | scanner.ScanChars | scanner.ScanStrings | scanner.SkipComments
    for tok := s.Scan(); tok != scanner.EOF; tok = s.Scan() {	
		content := s.TokenText()
		if content == "int" || content == "void" || content == "char" {
			tokens = append(tokens, Token{Type: TokType, Value: content, Line: s.Pos().Line})
		} else if content == "volatile" || content == "unsigned" || content == "long" || content == "short" || content == "static" || content == "const" || content == "extern" {
			tokens = append(tokens, Token{Type: TokQualifier, Value: content, Line: s.Pos().Line})
		} else if content == "return" {
			tokens = append(tokens, Token{Type: TokReturn, Value: content, Line: s.Pos().Line})
		} else if content == "if" {
			tokens = append(tokens, Token{Type: TokIf, Value: content, Line: s.Pos().Line})
		} else if content == "else" {
			tokens = append(tokens, Token{Type: TokElse, Value: content, Line: s.Pos().Line})
		} else if _, err := strconv.ParseInt(content, 0, 64); err == nil {
			tokens = append(tokens, Token{Type: TokNumber, Value: content, Line: s.Pos().Line})
		} else if content == "(" {
			tokens = append(tokens, Token{Type: TokLParen, Value: content, Line: s.Pos().Line})
		} else if content == ")" {
			tokens = append(tokens, Token{Type: TokRParen, Value: content, Line: s.Pos().Line})
		} else if content == "{" {
			tokens = append(tokens, Token{Type: TokLCurly, Value: content, Line: s.Pos().Line})
		} else if content == "}" {
			tokens = append(tokens, Token{Type: TokRCurly, Value: content, Line: s.Pos().Line})
		} else if content == ";" {
			tokens = append(tokens, Token{Type: TokSemi, Value: content, Line: s.Pos().Line})
		} else if content == "+" {
			tokens = append(tokens, Token{Type: TokPlus, Value: content, Line: s.Pos().Line})
		} else if content == "-" {
			tokens = append(tokens, Token{Type: TokMinus, Value: content, Line: s.Pos().Line})
		} else if content == "*" {
			tokens = append(tokens, Token{Type: TokStar, Value: content, Line: s.Pos().Line})
		} else if content == "/" {
			tokens = append(tokens, Token{Type: TokSlash, Value: content, Line: s.Pos().Line})
		} else if content == "=" {
			tokens = append(tokens, Token{Type: TokEqual, Value: content, Line: s.Pos().Line})
		} else if content == "," {
			tokens = append(tokens, Token{Type: TokComma, Value: content, Line: s.Pos().Line})
		} else if content == ":" {
			tokens = append(tokens, Token{Type: TokColon, Value: content, Line: s.Pos().Line})
		} else if content == "goto" {
			tokens = append(tokens, Token{Type: TokGoto, Value: content, Line: s.Pos().Line})
		} else if content == "for" {
			tokens = append(tokens, Token{Type: TokFor, Value: content, Line: s.Pos().Line})
		} else if content == "while" {
			tokens = append(tokens, Token{Type: TokWhile, Value: content, Line: s.Pos().Line})
		} else if content == "do" {
			tokens = append(tokens, Token{Type: TokDo, Value: content, Line: s.Pos().Line})
		} else if content == "<" {
			tokens = append(tokens, Token{Type: TokLAngle, Value: content, Line: s.Pos().Line})
		} else if content == ">" {
			tokens = append(tokens, Token{Type: TokRAngle, Value: content, Line: s.Pos().Line})
		} else {
			tokens = append(tokens, Token{Type: TokIdent, Value: content, Line: s.Pos().Line})
		} 
	}
	return tokens
}
