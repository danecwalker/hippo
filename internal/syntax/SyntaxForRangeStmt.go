package syntax

import (
	"bytes"
)

type ForRangeStmt struct {
	ForPos *Position
	Key    *Identifier
	X      Expression
	Body   *BlockStmt
}

func NewForRangeStmt(forPos *Position, key *Identifier, x Expression, body *BlockStmt) *ForRangeStmt {
	return &ForRangeStmt{
		ForPos: forPos,
		Key:    key,
		X:      x,
		Body:   body,
	}
}

func (fs *ForRangeStmt) statementNode() {}

func (fs *ForRangeStmt) Pos() *Position {
	return fs.ForPos
}

func (fs *ForRangeStmt) End() *Position {
	return fs.Body.End()
}

func (fs *ForRangeStmt) PrettyPrint(w *bytes.Buffer, indent int) {
	addIndent(w, indent)
	w.WriteString("ForRangeStmt:\n")
	addIndent(w, indent+1)
	w.WriteString("Key:\n")
	fs.Key.PrettyPrint(w, indent+2)
	addIndent(w, indent+1)
	w.WriteString("X:\n")
	fs.X.PrettyPrint(w, indent+2)
	addIndent(w, indent+1)
	w.WriteString("Body:\n")
	fs.Body.PrettyPrint(w, indent+2)
}

func (fs *ForRangeStmt) String() string {
	var buf bytes.Buffer
	fs.PrettyPrint(&buf, 0)
	return buf.String()
}
