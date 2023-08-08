package intermediate

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
	if e.Pos.Filename != "" {
		return fmt.Sprintf("%s:%d:%d:\n\t\033[31merror:\033[0m %s", e.Pos.Filename, e.Pos.Line, e.Pos.Column, e.Msg)
	}

	return fmt.Sprintf("\033[31merror:\033[0m %s", e.Msg)
}

func NewUndefinedObjectError(pos *syntax.Position, name string) *Error {
	msg := "undefined object %s"
	return NewError(pos, fmt.Sprintf(msg, name))
}

func NewDisallowedTopLevelStatementError(pos *syntax.Position) *Error {
	return NewError(pos, "only variable and function declarations are allowed at the top level")
}

func NewUnexpectedStmt(pos *syntax.Position) *Error {
	return NewError(pos, "unexpected statement")
}

func NewUnexpectedExpr(pos *syntax.Position) *Error {
	return NewError(pos, "unexpected expression")
}

func NewUnexpectedBasicLit(pos *syntax.Position) *Error {
	return NewError(pos, "unexpected basic literal")
}
