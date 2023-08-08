package syntax

import (
	"bytes"
	"fmt"
)

type ReturnStmt struct {
	ReturnPos *Position
	Result    Expression
}

func NewReturnStmt(returnPos *Position, result Expression) *ReturnStmt {
	return &ReturnStmt{
		ReturnPos: returnPos,
		Result:    result,
	}
}

func (rs *ReturnStmt) statementNode() {}
func (rs *ReturnStmt) Pos() *Position {
	return rs.ReturnPos
}

func (rs *ReturnStmt) PrettyPrint(w *bytes.Buffer, indent int) {
	addIndent(w, indent)
	w.WriteString(fmt.Sprintf("ReturnStmt (%s):\n", rs.ReturnPos))
	rs.Result.PrettyPrint(w, indent+1)
}
