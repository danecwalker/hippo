package lexer

import (
	"fmt"

	"github.com/danecwalker/hippo/internal/syntax"
)

type Error struct {
	Pos syntax.Position
	Msg string
}

func NewError(pos *syntax.Position, msg string) *Error {
	return &Error{
		Pos: *pos,
		Msg: msg,
	}
}

func NewEOFError(pos *syntax.Position) *Error {
	return NewError(pos, "unexpected end of file")
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s:%d:%d:\n\t\033[31merror:\033[0m %s", e.Pos.Filename, e.Pos.Line, e.Pos.Column, e.Msg)
}
