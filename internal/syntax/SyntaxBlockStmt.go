package syntax

import (
	"bytes"
	"fmt"
)

type BlockStmt struct {
	Lbrace *Position
	Stmts  []Statement
	Rbrace *Position
}

func NewBlockStmt(lbrace *Position) *BlockStmt {
	return &BlockStmt{Lbrace: lbrace, Stmts: make([]Statement, 0)}
}

func (b *BlockStmt) AddStmt(stmt Statement) {
	b.Stmts = append(b.Stmts, stmt)
}

func (b *BlockStmt) Pos() *Position {
	return b.Lbrace
}

func (b *BlockStmt) End() *Position {
	return b.Rbrace
}

func (b *BlockStmt) statementNode() {}
func (b *BlockStmt) PrettyPrint(w *bytes.Buffer, indent int) {
	addIndent(w, indent)
	w.WriteString("BlockStmt:\n")
	addIndent(w, indent+1)
	w.WriteString(fmt.Sprintf("LBrace (%s)\n", b.Pos()))
	addIndent(w, indent+1)
	w.WriteString("Stmts:\n")
	for _, stmt := range b.Stmts {
		stmt.PrettyPrint(w, indent+2)
	}
	addIndent(w, indent+1)
	w.WriteString(fmt.Sprintf("RBrace (%s)\n", b.End()))
}
