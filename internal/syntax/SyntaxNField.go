package syntax

import "bytes"

type NField struct {
	Names []*Identifier
	Type  *Identifier
}

func NewNField(names []*Identifier, type_ *Identifier) *NField {
	return &NField{
		Names: names,
		Type:  type_,
	}
}

func (f *NField) PrettyPrint(w *bytes.Buffer, indent int) {
	addIndent(w, indent)
	w.WriteString("Field:\n")
	addIndent(w, indent+1)
	w.WriteString("Names:\n")
	for _, name := range f.Names {
		name.PrettyPrint(w, indent+2)
	}
	addIndent(w, indent+1)
	w.WriteString("Type:\n")
	f.Type.PrettyPrint(w, indent+2)
}
