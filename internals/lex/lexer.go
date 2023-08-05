package lex

import (
	"github.com/danecwalker/hippo/internals/symbol"
	"github.com/danecwalker/hippo/internals/utils"
)

type Lexer struct {
	input []byte

	*utils.Location
	Offset int
	ch     byte
}

func (l *Lexer) next() {
	l.Offset += 1
	if l.Offset >= len(l.input) {
		l.ch = 0
		return
	}

	l.ch = l.input[l.Offset]

	if l.ch == '\n' {
		l.Column = 0
		l.Line += 1
	} else {
		l.Column += 1
	}
}

func NewLexer(file_name string, input []byte) *Lexer {
	return &Lexer{
		input: input,
		Location: &utils.Location{
			Line:     1,
			Column:   1,
			Filename: file_name,
		},
		Offset: 0,
		ch:     input[0],
	}
}

func (l *Lexer) NextSymbol() symbol.Symbol {
	var sym symbol.Symbol
	l.skipWS()
	switch l.ch {
	case 0:
		sym = symbol.NewSymbol(symbol.EOF, l.ch, *l.Location)
		return sym
	case '=':
		sym = symbol.NewSymbol(symbol.ASSIGN, l.ch, *l.Location)
	case ':':
		sym = symbol.NewSymbol(symbol.COLON, l.ch, *l.Location)
	case ',':
		sym = symbol.NewSymbol(symbol.COMMA, l.ch, *l.Location)
	case '{':
		sym = symbol.NewSymbol(symbol.LBRACE, l.ch, *l.Location)
	case '}':
		sym = symbol.NewSymbol(symbol.RBRACE, l.ch, *l.Location)
	case '(':
		sym = symbol.NewSymbol(symbol.LPAREN, l.ch, *l.Location)
	case ')':
		sym = symbol.NewSymbol(symbol.RPAREN, l.ch, *l.Location)
	case '"':
		loc := *l.Location
		str := l.readString()
		sym = symbol.NewSymbol(symbol.STRING, str, loc)
		return sym
	default:
		if isLetter(l.ch) {
			ident, loc := l.readIdentifier()
			sym = symbol.NewSymbol(symbol.LookupIdent(ident), ident, loc)
			return sym
		} else if isNumber(l.ch) {
			t, num, loc := l.readNumber()
			sym = symbol.NewSymbol(t, num, loc)
			return sym
		} else {
			sym = symbol.NewSymbol(symbol.ILLEGAL, l.ch, *l.Location)
		}
	}
	l.next()
	return sym
}

func (l *Lexer) skipWS() {
	for isWhitespace(l.ch) {
		l.next()
	}
}

func (l *Lexer) readString() string {
	offset := l.Offset
	for {
		l.next()
		if l.ch == '"' || l.ch == 0 {
			break
		}
	}
	l.next()
	return string(l.input[offset:l.Offset])
}

func (l *Lexer) readIdentifier() (string, utils.Location) {
	offset := l.Offset
	loc := *l.Location
	for isLetter(l.ch) || isNumber(l.ch) {
		l.next()
	}
	return string(l.input[offset:l.Offset]), loc
}

func (l *Lexer) readNumber() (symbol.SymbolType, string, utils.Location) {
	offset := l.Offset
	loc := *l.Location
	t := symbol.INT
	for isNumber(l.ch) || l.ch == '.' {
		if l.ch == '.' {
			t = symbol.FLOAT
		}
		l.next()
	}
	return t, string(l.input[offset:l.Offset]), loc
}

func isWhitespace(ch byte) bool {
	return ch == ' ' || ch == '\t' || ch == '\n'
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isNumber(ch byte) bool {
	return '0' <= ch && ch <= '9'
}
