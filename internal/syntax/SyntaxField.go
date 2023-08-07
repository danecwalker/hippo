package syntax

import "bytes"

type Field struct {
	Name *Identifier
	Type *Identifier
}

func NewField(name *Identifier, type_ *Identifier) *Field {
	return &Field{
		Name: name,
		Type: type_,
	}
}

func (f *Field) PrettyPrint(w *bytes.Buffer, indent int) {
	addIndent(w, indent)
	w.WriteString("Field:\n")
	addIndent(w, indent+1)
	w.WriteString("Name:\n")
	f.Name.PrettyPrint(w, indent+2)
	addIndent(w, indent+1)
	w.WriteString("Type:\n")
	f.Type.PrettyPrint(w, indent+2)
}
