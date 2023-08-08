package intermediate

import "bytes"

type AssignInst struct {
	Left  Inst
	Right Inst
}

func NewAssignInst(left Inst, right Inst) *AssignInst {
	return &AssignInst{
		Left:  left,
		Right: right,
	}
}

func (inst *AssignInst) inst() {}

func (inst *AssignInst) pretty(w *bytes.Buffer, indent int) {
	inst.Left.pretty(w, indent)
	w.WriteString(" = ")
	inst.Right.pretty(w, indent)
	w.WriteRune('\n')
}
