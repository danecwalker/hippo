package parse

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

func (e *Error) Error() string {
	return fmt.Sprintf("%s:%d:%d:\n\t\033[31merror:\033[0m %s", e.Pos.Filename, e.Pos.Line, e.Pos.Column, e.Msg)
}

func NewEOFError(pos *syntax.Position) *Error {
	return NewError(pos, "unexpected end of file")
}

func NewPeekError(pos *syntax.Position, expected, actual string) *Error {
	msg := "expected %s, got %s"
	return NewError(pos, fmt.Sprintf(msg, expected, actual))
}

func NewUnexpectedTokenError(pos *syntax.Position, expected *syntax.Token) *Error {
	msg := "unexpected token %s"
	return NewError(pos, fmt.Sprintf(msg, expected))
}

func NewRedeclaredError(pos *syntax.Position, name string) *Error {
	msg := "redeclared %s"
	return NewError(pos, fmt.Sprintf(msg, name))
}

func NewInvalidRangeError(pos *syntax.Position) *Error {
	msg := "invalid range"
	return NewError(pos, fmt.Sprintf(msg))
}
