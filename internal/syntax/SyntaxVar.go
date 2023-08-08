package syntax

import (
	"bytes"
	"fmt"
)

type VarStatement struct {
	Position *Position
	Names    []*Identifier
	Type     *Identifier
	Values   []Expression
}

func NewVarStatement(position *Position, names []*Identifier, type_ *Identifier, values []Expression) *VarStatement {
	return &VarStatement{
		Position: position,
		Names:    names,
		Type:     type_,
		Values:   values,
	}
}

func (vs *VarStatement) statementNode() {}
func (vs *VarStatement) Pos() *Position {
	return vs.Position
}

func (vs *VarStatement) PrettyPrint(w *bytes.Buffer, indent int) {
	addIndent(w, indent)
	w.WriteString(fmt.Sprintf("VarStatement (%s):\n", vs.Position))
	addIndent(w, indent+1)
	w.WriteString("Names:\n")
	for _, name := range vs.Names {
		name.PrettyPrint(w, indent+2)
	}
	addIndent(w, indent+1)
	w.WriteString("Type:\n")
	if vs.Type != nil {
		vs.Type.PrettyPrint(w, indent+2)
	} else {
		addIndent(w, indent+2)
		w.WriteString("nil\n")
	}

	addIndent(w, indent+1)
	w.WriteString("Values:\n")
	for _, value := range vs.Values {
		value.PrettyPrint(w, indent+2)
	}
	w.WriteString("\n")
}
