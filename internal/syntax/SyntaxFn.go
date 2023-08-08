package syntax

import (
	"bytes"
	"fmt"
)

type FuncStmt struct {
	FuncPos *Position
	Name    *Identifier
	Type    *FuncType
	Body    *BlockStmt
}

func NewFuncStmt(funcPos *Position, name *Identifier, type_ *FuncType, body *BlockStmt) *FuncStmt {
	return &FuncStmt{
		FuncPos: funcPos,
		Name:    name,
		Type:    type_,
		Body:    body,
	}
}

func (fs *FuncStmt) statementNode() {}
func (fs *FuncStmt) Pos() *Position {
	return fs.FuncPos
}

func (fs *FuncStmt) PrettyPrint(w *bytes.Buffer, indent int) {
	addIndent(w, indent)
	w.WriteString(fmt.Sprintf("FuncStmt (%s):\n", fs.FuncPos))
	addIndent(w, indent+1)
	w.WriteString("Name:\n")
	fs.Name.PrettyPrint(w, indent+2)

	addIndent(w, indent+1)
	w.WriteString("Type:\n")
	fs.Type.PrettyPrint(w, indent+2)

	addIndent(w, indent+1)
	w.WriteString("Body:\n")
	fs.Body.PrettyPrint(w, indent+2)
}
