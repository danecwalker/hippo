package syntax

import "bytes"

type IfStmt struct {
	If          *Position
	Cond        Expression
	Consequence *BlockStmt
	Alternative Statement
}

func NewIfStmt(if_ *Position, cond Expression, consequence *BlockStmt, alternative Statement) *IfStmt {
	return &IfStmt{
		If:          if_,
		Cond:        cond,
		Consequence: consequence,
		Alternative: alternative,
	}
}

func (is *IfStmt) statementNode() {}

func (is *IfStmt) Pos() *Position {
	return is.If
}

func (is *IfStmt) PrettyPrint(w *bytes.Buffer, indent int) {
	addIndent(w, indent)
	w.WriteString("IfStmt:\n")
	addIndent(w, indent+1)
	w.WriteString("Cond:\n")
	is.Cond.PrettyPrint(w, indent+2)
	addIndent(w, indent+1)
	w.WriteString("Consequence:\n")
	is.Consequence.PrettyPrint(w, indent+2)
	if is.Alternative != nil {
		addIndent(w, indent+1)
		w.WriteString("Alternative:\n")
		is.Alternative.PrettyPrint(w, indent+2)
	}
}
