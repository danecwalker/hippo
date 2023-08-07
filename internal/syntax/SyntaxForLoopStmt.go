package syntax

import "bytes"

type ForLoopStmt struct {
	Init Statement
	Cond Expression
	Post Statement
	Body *BlockStmt
}

func NewForLoopStmt(init Statement, cond Expression, post Statement, body *BlockStmt) *ForLoopStmt {
	return &ForLoopStmt{Init: init, Cond: cond, Post: post, Body: body}
}

func (f *ForLoopStmt) statementNode() {}
func (fs *ForLoopStmt) PrettyPrint(w *bytes.Buffer, indent int) {
	addIndent(w, indent)
	w.WriteString("ForLoopStmt:\n")
	addIndent(w, indent+1)
	w.WriteString("Init:\n")
	fs.Init.PrettyPrint(w, indent+2)
	addIndent(w, indent+1)
	w.WriteString("Cond:\n")
	fs.Cond.PrettyPrint(w, indent+2)
	addIndent(w, indent+1)
	w.WriteString("Post:\n")
	fs.Post.PrettyPrint(w, indent+2)
	addIndent(w, indent+1)
	w.WriteString("Body:\n")
	fs.Body.PrettyPrint(w, indent+2)
}
