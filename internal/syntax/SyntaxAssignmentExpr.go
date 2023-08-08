package syntax

import (
	"bytes"
	"fmt"
)

type AssignmentExpr struct {
	Lhs Expression
	Op  *Token
	Rhs Expression
}

func NewAssignmentExpr(lhs Expression, op *Token, rhs Expression) *AssignmentExpr {
	return &AssignmentExpr{
		Lhs: lhs,
		Op:  op,
		Rhs: rhs,
	}
}

func (e *AssignmentExpr) expressionNode() {}
func (e *AssignmentExpr) Pos() *Position {
	return e.Lhs.Pos()
}

func (e *AssignmentExpr) PrettyPrint(w *bytes.Buffer, indent int) {
	addIndent(w, indent)
	w.WriteString("AssignmentExpr:\n")
	addIndent(w, indent+1)
	w.WriteString("Lhs:\n")
	e.Lhs.PrettyPrint(w, indent+2)
	addIndent(w, indent+1)
	w.WriteString(fmt.Sprintf("Op: %s\n", e.Op.Literal))
	addIndent(w, indent+1)
	w.WriteString("Rhs:\n")
	e.Rhs.PrettyPrint(w, indent+2)
}
