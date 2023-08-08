package syntax

import "bytes"

type ExpressionStmt struct {
	X Expression
}

func NewExpressionStmt(x Expression) *ExpressionStmt {
	return &ExpressionStmt{
		X: x,
	}
}

func (es *ExpressionStmt) statementNode() {}
func (es *ExpressionStmt) Pos() *Position {
	return es.X.Pos()
}

func (es *ExpressionStmt) PrettyPrint(w *bytes.Buffer, indent int) {
	addIndent(w, indent)
	w.WriteString("ExpressionStmt:\n")
	es.X.PrettyPrint(w, indent+1)
}
