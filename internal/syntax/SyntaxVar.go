package syntax

import (
	"bytes"
	"fmt"
)

type VarStatement struct {
	Position *Position
	Name     *Identifier
	Type     *Identifier
	Value    Expression
	Kind     string
}

func NewVarStatement(position *Position, kind string, name *Identifier, type_ *Identifier, value Expression) *VarStatement {
	return &VarStatement{
		Position: position,
		Name:     name,
		Type:     type_,
		Value:    value,
		Kind:     kind,
	}
}

func (vs *VarStatement) statementNode() {}
func (vs *VarStatement) PrettyPrint(w *bytes.Buffer, indent int) {
	addIndent(w, indent)
	w.WriteString(fmt.Sprintf("VarStatement (%s):\n", vs.Position))
	addIndent(w, indent+1)
	w.WriteString("Names:\n")
	vs.Name.PrettyPrint(w, indent+2)
	addIndent(w, indent+1)
	w.WriteString("Type:\n")
	if vs.Type != nil {
		vs.Type.PrettyPrint(w, indent+2)
	} else {
		addIndent(w, indent+2)
		w.WriteString("nil\n")
	}

	addIndent(w, indent+1)
	w.WriteString("Value:")
	if vs.Value != nil {
		w.WriteString("\n")
		vs.Value.PrettyPrint(w, indent+2)
	} else {
		addIndent(w, indent+2)
		w.WriteString(" nil\n")
	}
	addIndent(w, indent+1)
	w.WriteString("Kind: ")
	w.WriteString(vs.Kind)
	w.WriteString("\n")

}
