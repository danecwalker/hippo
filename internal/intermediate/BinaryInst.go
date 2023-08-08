package intermediate

import "bytes"

type BinaryInst struct {
	Op    string
	Left  Inst
	Right Inst
}

func NewBinaryInst(op string, left Inst, right Inst) *BinaryInst {
	return &BinaryInst{
		Op:    op,
		Left:  left,
		Right: right,
	}
}

func (inst *BinaryInst) inst() {}
func (inst *BinaryInst) pretty(w *bytes.Buffer, indent int) {
	inst.Left.pretty(w, indent)
	w.WriteRune(' ')
	w.WriteString(inst.Op)
	w.WriteRune(' ')
	inst.Right.pretty(w, indent)
}

func NewAddInst(left Inst, right Inst) Inst {
	return NewBinaryInst("+", left, right)
}

func NewSubInst(left Inst, right Inst) Inst {
	return NewBinaryInst("-", left, right)
}

func NewMulInst(left Inst, right Inst) Inst {
	return NewBinaryInst("*", left, right)
}

func NewDivInst(left Inst, right Inst) Inst {
	return NewBinaryInst("/", left, right)
}
