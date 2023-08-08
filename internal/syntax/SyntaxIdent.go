package syntax

import (
	"bytes"
	"fmt"
)

type Identifier struct {
	NamePos *Position
	Name    string
	Obj     *Object
}

func NewIdentifier(namePos *Position, name string) *Identifier {
	return &Identifier{
		NamePos: namePos,
		Name:    name,
	}
}

func (i *Identifier) expressionNode() {}
func (i *Identifier) Pos() *Position {
	return i.NamePos
}

func (i *Identifier) PrettyPrint(w *bytes.Buffer, indent int) {
	addIndent(w, indent)
	w.WriteString(fmt.Sprintf("Identifier: (%s)\n", i.NamePos))
	addIndent(w, indent+1)
	w.WriteString("Name: ")
	w.WriteString(i.Name)
	w.WriteString("\n")
	addIndent(w, indent+1)
	w.WriteString("Obj: ")
	w.WriteString(fmt.Sprintf("%v\n", i.Obj))
}
