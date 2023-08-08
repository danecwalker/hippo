package lexer

import (
	"io/ioutil"
	"os"

	"github.com/danecwalker/hippo/internal/syntax"
)

type Lexer struct {
	Offset int

	Line     int
	Column   int
	Filename string

	input []byte
}

func NewLexer(filename string) *Lexer {
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	input, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}

	return &Lexer{
		Offset:   0,
		Line:     1,
		Column:   1,
		Filename: filename,
		input:    input,
	}
}

func (l *Lexer) Pos() *syntax.Position {
	return syntax.NewPosition(l.Offset, l.Line, l.Column, l.Filename)
}

func (l *Lexer) Next() (byte, error) {
	if l.Offset >= len(l.input) {
		return 0, NewEOFError(l.Pos())
	}

	b := l.input[l.Offset]
	l.Offset++
	l.Column++

	if b == '\n' {
		l.Line++
		l.Column = 1
	}

	return b, nil
}

func (l *Lexer) Peek() (byte, error) {
	return l.PeekN(1)
}

func (l *Lexer) PeekN(n int) (byte, error) {
	if l.Offset+n >= len(l.input) {
		return 0, NewEOFError(l.Pos())
	}

	return l.input[l.Offset+n], nil
}

func (l *Lexer) NextToken() *syntax.Token {
	l.skip() // skip whitespace

	b, err := l.PeekN(0)
	if err != nil {
		return syntax.NewToken(syntax.TokenEOF, "", l.Pos())
	}

	switch b {
	case '-':
		p := l.Pos()
		n, err := l.Peek()
		if err != nil {
			return syntax.NewToken(syntax.TokenEOF, "", l.Pos())
		}
		if n == '>' {
			l.Next()
			l.Next()
			return syntax.NewToken(syntax.TokenArrow, "->", p)
		} else {
			l.Next()
			return syntax.NewToken(syntax.TokenMinus, "-", p)
		}
	case '=':
		p := l.Pos()
		l.Next()
		return syntax.NewToken(syntax.TokenAssign, "=", p)
	case '+':
		p := l.Pos()
		l.Next()
		return syntax.NewToken(syntax.TokenPlus, "+", p)
	case '*':
		p := l.Pos()
		l.Next()
		return syntax.NewToken(syntax.TokenStar, "*", p)
	case '/':
		p := l.Pos()
		l.Next()
		return syntax.NewToken(syntax.TokenSlash, "/", p)
	case '<':
		p := l.Pos()
		n, err := l.Peek()
		if err != nil {
			return syntax.NewToken(syntax.TokenEOF, "", l.Pos())
		}
		if n == '-' {
			l.Next()
			l.Next()
			return syntax.NewToken(syntax.TokenInfer, "<-", p)
		} else {
			l.Next()
			return syntax.NewToken(syntax.TokenLt, "<", l.Pos())
		}
	case '>':
		p := l.Pos()
		l.Next()
		return syntax.NewToken(syntax.TokenGt, ">", p)
	case ',':
		p := l.Pos()
		l.Next()
		return syntax.NewToken(syntax.TokenComma, ",", p)
	case '.':
		p := l.Pos()
		n, err := l.Peek()
		if err != nil {
			return syntax.NewToken(syntax.TokenEOF, "", l.Pos())
		}
		if n == '.' {
			l.Next()
			l.Next()
			return syntax.NewToken(syntax.TokenRange, "..", p)
		} else {
			l.Next()
			return syntax.NewToken(syntax.TokenDot, ".", p)
		}
	case ';':
		p := l.Pos()
		l.Next()
		return syntax.NewToken(syntax.TokenSemicolon, ";", p)
	case ':':
		p := l.Pos()
		l.Next()
		return syntax.NewToken(syntax.TokenColon, ":", p)
	case '(':
		p := l.Pos()
		l.Next()
		return syntax.NewToken(syntax.TokenLParen, "(", p)
	case ')':
		p := l.Pos()
		l.Next()
		return syntax.NewToken(syntax.TokenRParen, ")", p)
	case '{':
		p := l.Pos()
		l.Next()
		return syntax.NewToken(syntax.TokenLBrace, "{", p)
	case '}':
		p := l.Pos()
		l.Next()
		return syntax.NewToken(syntax.TokenRBrace, "}", p)
	default:
		if isLetter(b) {
			return l.scanIdentifier()
		} else if isDigit(b) {
			return l.scanNumber()
		}
		l.Next()
		return syntax.NewToken(syntax.TokenIllegal, string(b), l.Pos())
	}
}

func (l *Lexer) skip() {
	for {
		b, err := l.PeekN(0)
		if err != nil {
			return
		}

		if b == ' ' || b == '\t' || b == '\n' || b == '\r' {
			l.Next()
		} else {
			return
		}
	}
}

func isDigit(b byte) bool {
	return '0' <= b && b <= '9'
}

func (l *Lexer) scanNumber() *syntax.Token {
	pos := l.Pos()

	b, err := l.PeekN(0)
	if err != nil {
		return syntax.NewToken(syntax.TokenEOF, "", l.Pos())
	}

	literal := string(b)
	l.Next()

	for {
		b, err := l.PeekN(0)
		if err != nil {
			break
		}

		if isDigit(b) {
			literal += string(b)
			l.Next()
		} else {
			break
		}
	}

	return syntax.NewToken(syntax.TokenInt, literal, pos)
}

func isLetter(b byte) bool {
	return 'a' <= b && b <= 'z' || 'A' <= b && b <= 'Z' || b == '_'
}

func (l *Lexer) scanIdentifier() *syntax.Token {
	pos := l.Pos()

	b, err := l.PeekN(0)
	if err != nil {
		return syntax.NewToken(syntax.TokenEOF, "", l.Pos())
	}

	literal := string(b)
	l.Next()

	for {
		b, err := l.PeekN(0)
		if err != nil {
			break
		}

		if isLetter(b) || isDigit(b) {
			literal += string(b)
			l.Next()
		} else {
			break
		}
	}

	return syntax.NewToken(syntax.LookupIdent(literal), literal, pos)
}
