package tok

import "fmt"

type Bounds struct {
	Start Position
	End   Position
}

func NewBounds(start Position, end Position) *Bounds {
	return &Bounds{
		Start: start,
		End:   end,
	}
}

type TokenType int

const (
	// Special tokens
	ILLEGAL TokenType = iota
	EOF
	COMMENT

	// Identifiers and literals
	IDENT
	INT
	FLOAT
	STRING

	// Operators
	ASSIGN
	PLUS
	MINUS
	BANG
	ASTERISK
	SLASH
	MODULO
	EXPONENT

	INC
	DEC
	RANGE
	EQ
	NOT_EQ
	LT
	GT
	LTE
	GTE

	// Delimiters
	COMMA
	SEMICOLON
	COLON
	DOT
	LARROW
	RARROW

	LPAREN
	RPAREN
	LBRACE
	RBRACE
	LBRACKET
	RBRACKET

	// Keywords
	VAR
	CONST
	FUNCTION
	TRUE
	FALSE
	IF
	ELSE
	FOR
	RETURN
)

type Token struct {
	Type    TokenType
	Literal string
	Bounds  *Bounds
}

func NewToken[T string | rune | byte](tokType TokenType, ch T, line int, column int, filename string) Token {
	return Token{
		Type:    tokType,
		Literal: string(ch),
		Bounds: NewBounds(
			NewPosition(line, column, filename),
			NewPosition(line, column+len(string(ch)), filename),
		),
	}
}

var keywords = map[string]TokenType{
	"var":   VAR,
	"const": CONST,
	"fn":    FUNCTION,
	"true":  TRUE,
	"false": FALSE,
	"if":    IF,
	"else":  ELSE,
	"for":   FOR,
	"ret":   RETURN,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}

func (t Token) String() string {
	// return fmt.Sprintf("Token(%s, %s) %s:%d:%d -> %s:%d:%d", stringer[t.Type], t.Literal, t.Bounds.Start.Filename, t.Bounds.Start.Line, t.Bounds.Start.Column, t.Bounds.End.Filename, t.Bounds.End.Line, t.Bounds.End.Column)
	return fmt.Sprintf("Token(%s, %s)", stringer[t.Type], t.Literal)
}

var stringer = map[TokenType]string{
	// Special tokens
	ILLEGAL: "ILLEGAL",
	EOF:     "EOF",
	COMMENT: "COMMENT",

	// Identifiers and literals
	IDENT:  "IDENT",
	INT:    "INT",
	FLOAT:  "FLOAT",
	STRING: "STRING",

	// Operators
	ASSIGN:   "ASSIGN",
	PLUS:     "PLUS",
	MINUS:    "MINUS",
	BANG:     "BANG",
	ASTERISK: "ASTERISK",
	SLASH:    "SLASH",
	MODULO:   "MODULO",
	EXPONENT: "EXPONENT",

	INC:    "INC",
	DEC:    "DEC",
	RANGE:  "RANGE",
	EQ:     "EQ",
	NOT_EQ: "NOT_EQ",
	LT:     "LT",
	GT:     "GT",
	LTE:    "LTE",
	GTE:    "GTE",

	// Delimiters
	COMMA:     "COMMA",
	SEMICOLON: "SEMICOLON",
	COLON:     "COLON",
	DOT:       "DOT",
	LARROW:    "LARROW",
	RARROW:    "RARROW",

	LPAREN:   "LPAREN",
	RPAREN:   "RPAREN",
	LBRACE:   "LBRACE",
	RBRACE:   "RBRACE",
	LBRACKET: "LBRACKET",
	RBRACKET: "RBRACKET",

	// Keywords
	VAR:      "VAR",
	CONST:    "CONST",
	FUNCTION: "FUNCTION",
	TRUE:     "TRUE",
	FALSE:    "FALSE",
	IF:       "IF",
	ELSE:     "ELSE",
	FOR:      "FOR",
	RETURN:   "RETURN",
}

func (t TokenType) String() string {
	return stringer[t]
}
