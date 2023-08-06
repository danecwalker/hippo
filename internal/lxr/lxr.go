package lxr

import (
	"os"

	"github.com/danecwalker/hippo/internal/logger"
	"github.com/danecwalker/hippo/internal/tok"
)

type Lexer struct {
	content  []byte
	ch       byte
	Offset   int
	Line     int
	Column   int
	Filename string
}

func NewLexer(file_name string) *Lexer {
	log := logger.NewLogHandler()

	// Open file
	file, err := os.Open(file_name)
	if err != nil {
		log.Log(logger.LogFatal, "Error opening file %s: %s", file_name, err)
	}
	defer file.Close()

	// Get file size
	file_info, err := file.Stat()
	if err != nil {
		log.Log(logger.LogFatal, "Error getting file info for %s: %s", file_name, err)
	}

	// Read file
	buf := make([]byte, file_info.Size())
	_, err = file.Read(buf)
	if err != nil {
		log.Log(logger.LogFatal, "Error reading file %s: %s", file_name, err)
	}

	return &Lexer{
		content:  buf,
		Offset:   0,
		ch:       buf[0],
		Line:     1,
		Column:   1,
		Filename: file_name,
	}
}

func (l *Lexer) Next() byte {
	l.Offset++
	if l.Offset >= len(l.content) {
		l.ch = 0
		return l.ch
	}

	l.ch = l.content[l.Offset]
	l.Column++
	if l.ch == '\n' {
		l.Line++
		l.Column = 0
	}

	return l.ch
}

func (l *Lexer) Peek() byte {
	if l.Offset+1 >= len(l.content) {
		return 0
	}

	return l.content[l.Offset+1]
}

func (l *Lexer) NextToken() tok.Token {
	l.skipWhitespace()
	var token tok.Token

	switch l.ch {
	case 0:
		token = tok.NewToken(tok.EOF, l.ch, l.Line, l.Column, l.Filename)
	case '=':
		if l.Peek() == '=' {
			l.Next()
			token = tok.NewToken(tok.EQ, "==", l.Line, l.Column-1, l.Filename)
		} else {
			token = tok.NewToken(tok.ASSIGN, l.ch, l.Line, l.Column, l.Filename)
		}
	case '+':
		if l.Peek() == '+' {
			l.Next()
			token = tok.NewToken(tok.INC, "++", l.Line, l.Column-1, l.Filename)
		} else {
			token = tok.NewToken(tok.PLUS, l.ch, l.Line, l.Column, l.Filename)
		}
	case '-':
		if l.Peek() == '-' {
			l.Next()
			token = tok.NewToken(tok.DEC, "--", l.Line, l.Column-1, l.Filename)
		} else if l.Peek() == '>' {
			l.Next()
			token = tok.NewToken(tok.RARROW, "->", l.Line, l.Column-1, l.Filename)
		} else {
			token = tok.NewToken(tok.MINUS, l.ch, l.Line, l.Column, l.Filename)
		}
	case '!':
		if l.Peek() == '=' {
			l.Next()
			token = tok.NewToken(tok.NOT_EQ, "!=", l.Line, l.Column-1, l.Filename)
		} else {
			token = tok.NewToken(tok.BANG, l.ch, l.Line, l.Column, l.Filename)
		}
	case '*':
		if l.Peek() == '*' {
			l.Next()
			token = tok.NewToken(tok.EXPONENT, "**", l.Line, l.Column-1, l.Filename)
		} else {
			token = tok.NewToken(tok.ASTERISK, l.ch, l.Line, l.Column, l.Filename)
		}
	case '/':
		token = tok.NewToken(tok.SLASH, l.ch, l.Line, l.Column, l.Filename)
	case '%':
		token = tok.NewToken(tok.MODULO, l.ch, l.Line, l.Column, l.Filename)
	case '<':
		if l.Peek() == '=' {
			l.Next()
			token = tok.NewToken(tok.LTE, "<=", l.Line, l.Column-1, l.Filename)
		} else if l.Peek() == '-' {
			l.Next()
			token = tok.NewToken(tok.LARROW, "<-", l.Line, l.Column-1, l.Filename)
		} else {
			token = tok.NewToken(tok.LT, l.ch, l.Line, l.Column, l.Filename)
		}
	case '>':
		if l.Peek() == '=' {
			l.Next()
			token = tok.NewToken(tok.GTE, ">=", l.Line, l.Column-1, l.Filename)
		} else {

			token = tok.NewToken(tok.GT, l.ch, l.Line, l.Column, l.Filename)
		}
	case ',':
		token = tok.NewToken(tok.COMMA, l.ch, l.Line, l.Column, l.Filename)
	case ':':
		token = tok.NewToken(tok.COLON, l.ch, l.Line, l.Column, l.Filename)
	case ';':
		token = tok.NewToken(tok.SEMICOLON, l.ch, l.Line, l.Column, l.Filename)
	case '.':
		if l.Peek() == '.' {
			l.Next()
			token = tok.NewToken(tok.RANGE, "..", l.Line, l.Column-1, l.Filename)
		} else {
			token = tok.NewToken(tok.DOT, l.ch, l.Line, l.Column, l.Filename)
		}
	case '(':
		token = tok.NewToken(tok.LPAREN, l.ch, l.Line, l.Column, l.Filename)
	case ')':
		token = tok.NewToken(tok.RPAREN, l.ch, l.Line, l.Column, l.Filename)
	case '{':
		token = tok.NewToken(tok.LBRACE, l.ch, l.Line, l.Column, l.Filename)
	case '}':
		token = tok.NewToken(tok.RBRACE, l.ch, l.Line, l.Column, l.Filename)
	case '[':
		token = tok.NewToken(tok.LBRACKET, l.ch, l.Line, l.Column, l.Filename)
	case ']':
		token = tok.NewToken(tok.RBRACKET, l.ch, l.Line, l.Column, l.Filename)
	case '"':
		token = tok.NewToken(tok.STRING, l.readString(), l.Line, l.Column, l.Filename)
	default:
		if isLetter(l.ch) {
			line := l.Line
			column := l.Column
			ident := l.readIdentifier()
			token = tok.NewToken(tok.LookupIdent(ident), ident, line, column, l.Filename)
			return token
		}
		if isDigit(l.ch) {
			line := l.Line
			column := l.Column
			t, num := l.readNumber()
			token = tok.NewToken(t, num, line, column, l.Filename)
			return token
		}
		token = tok.NewToken(tok.ILLEGAL, l.ch, l.Line, l.Column, l.Filename)
	}

	l.Next()
	return token
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.Next()
	}
}

func (l *Lexer) readIdentifier() string {
	start := l.Offset
	for isLetter(l.ch) || isDigit(l.ch) {
		l.Next()
	}

	return string(l.content[start:l.Offset])
}

func (l *Lexer) readNumber() (tok.TokenType, string) {
	start := l.Offset
	tt := tok.INT
	for isDigit(l.ch) {
		l.Next()
	}

	return tt, string(l.content[start:l.Offset])
}

func (l *Lexer) readString() string {
	l.Next()
	start := l.Offset
	for l.ch != '"' && l.ch != 0 {
		l.Next()
	}

	return string(l.content[start:l.Offset])
}

func (l *Lexer) readComment() string {
	l.Next()
	start := l.Offset
	for l.ch != '\n' && l.ch != 0 {
		l.Next()
	}

	return string(l.content[start:l.Offset])
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}
