package intermediate

import (
	"bytes"

	"github.com/danecwalker/hippo/internal/syntax"
)

type BasicLitInst struct {
	Value string
	Type  *syntax.Object
}

func NewBasicLitInst(value string, type_ *syntax.Object) *BasicLitInst {
	return &BasicLitInst{
		Value: value,
		Type:  type_,
	}
}

func (inst *BasicLitInst) inst() {}
func (inst *BasicLitInst) pretty(w *bytes.Buffer, indent int) {
	w.WriteString(inst.Value)
}
