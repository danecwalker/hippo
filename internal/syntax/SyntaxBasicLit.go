package syntax

import (
	"bytes"
	"fmt"
)

type BasicLit struct {
	ValuePos *Position
	Kind     string
	Value    string
}

func NewBasicLit(position *Position, kind string, value string) *BasicLit {
	return &BasicLit{
		ValuePos: position,
		Kind:     kind,
		Value:    value,
	}
}

func (bl *BasicLit) expressionNode() {}
func (bl *BasicLit) PrettyPrint(w *bytes.Buffer, indent int) {
	addIndent(w, indent)
	w.WriteString(fmt.Sprintf("BasicLit: (%s)\n", bl.ValuePos))
	addIndent(w, indent+1)
	w.WriteString("Kind: ")
	w.WriteString(bl.Kind)
	w.WriteString("\n")
	addIndent(w, indent+1)
	w.WriteString("Value: ")
	w.WriteString(bl.Value)
	w.WriteString("\n")
}
