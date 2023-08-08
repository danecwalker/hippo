package intermediate

import "bytes"

type IdentInst struct {
	Name string
}

func NewIdentInst(name string) *IdentInst {
	return &IdentInst{
		Name: name,
	}
}

func (inst *IdentInst) inst() {}
func (inst *IdentInst) pretty(w *bytes.Buffer, indent int) {
	w.WriteString(inst.Name)
}
