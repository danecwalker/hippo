package syntax

import "fmt"

type TokenType int

const (
	TokenIllegal TokenType = iota // Illegal token
	TokenEOF                      // End of file
	TokenComment                  // Comment

	// Identifiers and basic type literals
	TokenIdent // main
	TokenInt   // 12345

	// Operators and delimiters
	TokenAssign // =
	TokenPlus   // +
	TokenMinus  // -

	TokenSemicolon // ;
	TokenColon     // :
	TokenComma     // ,
	TokenLParen    // (
	TokenRParen    // )

	// Keywords
	TokenFunc  // fn
	TokenVar   // var
	TokenConst // const
)

var tokenNames = map[TokenType]string{
	TokenIllegal: "ILLEGAL",
	TokenEOF:     "EOF",
	TokenComment: "COMMENT",

	TokenIdent: "IDENT",
	TokenInt:   "INT",

	TokenAssign: "=",
	TokenPlus:   "+",
	TokenMinus:  "-",

	TokenSemicolon: "SEMICOLON",
	TokenColon:     "COLON",
	TokenComma:     "COMMA",
	TokenLParen:    "(",
	TokenRParen:    ")",

	TokenFunc:  "fn",
	TokenVar:   "var",
	TokenConst: "const",
}

func (tt TokenType) String() string {
	if name, ok := tokenNames[tt]; ok {
		return name
	}

	return "UNKNOWN"
}

var keywords = map[string]TokenType{
	"fn":    TokenFunc,
	"var":   TokenVar,
	"const": TokenConst,
}

type Token struct {
	Type     TokenType
	Literal  string
	Position *Position
}

func NewToken(tt TokenType, lit string, pos *Position) *Token {
	return &Token{
		Type:     tt,
		Literal:  lit,
		Position: pos,
	}
}

func (t *Token) String() string {
	return fmt.Sprintf("Token(%s, %s)", t.Type, t.Literal)
}

func (t *Token) IsEOF() bool {
	return t.Type == TokenEOF
}

func (t *Token) IsIllegal() bool {
	return t.Type == TokenIllegal
}

func (t *Token) IsComment() bool {
	return t.Type == TokenComment
}

func LookupIdent(lit string) TokenType {
	if tok, ok := keywords[lit]; ok {
		return tok
	}

	return TokenIdent
}
