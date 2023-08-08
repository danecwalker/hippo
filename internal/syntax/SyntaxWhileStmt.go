package syntax

import "bytes"

type WhileStmt struct {
	Cond Expression
	Body *BlockStmt
}

func NewWhileStmt(cond Expression, body *BlockStmt) *WhileStmt {
	return &WhileStmt{
		Cond: cond,
		Body: body,
	}
}

func (ws *WhileStmt) statementNode() {}
func (ws *WhileStmt) Pos() *Position {
	return ws.Cond.Pos()
}

func (ws *WhileStmt) PrettyPrint(w *bytes.Buffer, indent int) {
	addIndent(w, indent)
	w.WriteString("WhileStmt:\n")
	addIndent(w, indent+1)
	w.WriteString("Cond:\n")
	ws.Cond.PrettyPrint(w, indent+2)
	addIndent(w, indent+1)
	w.WriteString("Body:\n")
	ws.Body.PrettyPrint(w, indent+2)
}
