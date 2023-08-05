package symbol

import (
	"fmt"

	"github.com/danecwalker/hippo/internals/utils"
)

type SymbolType int

const (
	ILLEGAL SymbolType = iota
	EOF

	// Identifiers + literals
	IDENT  // main
	INT    // 12345
	FLOAT  // 123.45
	STRING // "abc"

	// Operators
	ASSIGN // =
	PLUS   // +

	// Delimiters
	COLON  // :
	COMMA  // ,
	LBRACE // {
	RBRACE // }
	LPAREN // (
	RPAREN // )

	// Keywords
	CONST // const
	VAR   // variable
	FN    // function
	RET   // return
)

func (s SymbolType) String() string {
	switch s {
	case ILLEGAL:
		return "illegal"
	case EOF:
		return "end of file"
	case IDENT:
		return "identifier"
	case STRING:
		return "string"
	case INT:
		return "integer"
	case FLOAT:
		return "float"
	case ASSIGN:
		return "="
	case PLUS:
		return "+"
	case COLON:
		return ":"
	case COMMA:
		return ","
	case LBRACE:
		return "{"
	case RBRACE:
		return "}"
	case LPAREN:
		return "("
	case RPAREN:
		return ")"
	case CONST:
		return "const"
	case VAR:
		return "variable"
	case FN:
		return "function"
	case RET:
		return "return"
	default:
		panic(fmt.Sprintf("undefined string representations for symbol type %d", s))
	}
}

var keywords = map[string]SymbolType{
	"const": CONST,
	"var":   VAR,
	"fn":    FN,
}

func LookupIdent(ident string) SymbolType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}

type Symbol struct {
	Type SymbolType
	Lit  string
	Loc  utils.Location
}

func NewSymbol[T string | byte | rune](typ SymbolType, lit T, loc utils.Location) Symbol {
	return Symbol{typ, string(lit), loc}
}

func (s Symbol) String() string {
	return fmt.Sprintf("%s: %s %s", s.Loc, s.Type, s.Lit)
}
