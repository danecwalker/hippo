package syntax

import "bytes"

type BinaryExpr struct {
	X  Expression
	Op *Token
	Y  Expression
}

func NewBinaryExpr(x Expression, op *Token, y Expression) *BinaryExpr {
	return &BinaryExpr{
		X:  x,
		Op: op,
		Y:  y,
	}
}

func (b *BinaryExpr) expressionNode() {}
func (b *BinaryExpr) PrettyPrint(w *bytes.Buffer, indent int) {
	addIndent(w, indent)
	w.WriteString("BinaryExpr:\n")
	addIndent(w, indent+1)
	w.WriteString("X:\n")
	b.X.PrettyPrint(w, indent+2)
	addIndent(w, indent+1)
	w.WriteString("Op: ")
	w.WriteString(b.Op.Literal)
	w.WriteString("\n")
	addIndent(w, indent+1)
	w.WriteString("Y:\n")
	b.Y.PrettyPrint(w, indent+2)
}
