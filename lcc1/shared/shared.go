package shared

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

var Bits int = 16
