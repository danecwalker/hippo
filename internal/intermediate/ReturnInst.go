package intermediate

import "bytes"

type ReturnInst struct {
	Expr Inst
}

func NewReturnInst(expr Inst) *ReturnInst {
	return &ReturnInst{
		Expr: expr,
	}
}

func (inst *ReturnInst) inst() {}

func (inst *ReturnInst) pretty(w *bytes.Buffer, indent int) {
	w.WriteString("return")
	if inst.Expr != nil {
		w.WriteRune(' ')
		inst.Expr.pretty(w, indent)
	}
	w.WriteRune('\n')
}
