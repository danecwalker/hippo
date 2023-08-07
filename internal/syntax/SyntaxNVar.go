package syntax

import (
	"bytes"
	"fmt"
)

type NVarStatement struct {
	Position *Position
	Names    []*Identifier
	Type     *Identifier
	Values   []Expression
	Kind     string
}

func NewNVarStatement(position *Position, kind string, names []*Identifier, type_ *Identifier, values []Expression) *NVarStatement {
	return &NVarStatement{
		Position: position,
		Names:    names,
		Type:     type_,
		Values:   values,
		Kind:     kind,
	}
}

func (vs *NVarStatement) statementNode() {}
func (vs *NVarStatement) PrettyPrint(w *bytes.Buffer, indent int) {
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
	addIndent(w, indent+1)
	w.WriteString("Kind: ")
	w.WriteString(vs.Kind)
	w.WriteString("\n")

}
