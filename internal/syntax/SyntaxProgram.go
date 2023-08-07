package syntax

import (
	"bytes"
	"fmt"
)

type Program struct {
	Statements []Statement
}

func NewProgram() *Program {
	return &Program{
		Statements: make([]Statement, 0),
	}
}

func (p *Program) AddStatement(stmt Statement) {
	p.Statements = append(p.Statements, stmt)
}

func (p *Program) PrettyPrint() {
	var w bytes.Buffer
	w.WriteString("File:\n")
	for _, stmt := range p.Statements {
		stmt.PrettyPrint(&w, 1)
	}

	fmt.Println(w.String())
}

func addIndent(w *bytes.Buffer, indent int) {
	for i := 0; i < indent; i++ {
		w.WriteString("  ")
	}
}
