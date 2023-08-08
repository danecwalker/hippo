package intermediate

import (
	"bytes"

	"github.com/danecwalker/hippo/internal/syntax"
)

type AllocInst struct {
	Name  string
	Value Inst
	Type  *syntax.Object
}

func NewAllocInst(name string, value Inst, type_ *syntax.Object) *AllocInst {
	return &AllocInst{
		Name:  name,
		Value: value,
		Type:  type_,
	}
}
func (inst *AllocInst) inst() {}
func (inst *AllocInst) pretty(w *bytes.Buffer, indent int) {
	w.WriteString(inst.Name)
	w.WriteString(" = ")
	inst.Value.pretty(w, indent)
	w.WriteRune(' ')
	w.WriteRune(':')
	w.WriteRune(' ')
	w.WriteString(inst.Type.Name)
	w.WriteRune('\n')
}
