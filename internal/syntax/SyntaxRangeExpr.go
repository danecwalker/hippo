package syntax

import "bytes"

type RangeExpr struct {
	Low  Expression
	High Expression
	Kind string
}

func NewRangeExpr(low Expression, high Expression, type_ string) *RangeExpr {
	return &RangeExpr{
		Low:  low,
		High: high,
		Kind: type_,
	}
}

func (re *RangeExpr) expressionNode() {}
func (re *RangeExpr) Pos() *Position {
	return re.Low.Pos()
}

func (re *RangeExpr) PrettyPrint(w *bytes.Buffer, indent int) {
	addIndent(w, indent)
	w.WriteString("RangeExpr:\n")
	addIndent(w, indent+1)
	w.WriteString("Low:\n")
	re.Low.PrettyPrint(w, indent+2)
	addIndent(w, indent+1)
	w.WriteString("High:\n")
	re.High.PrettyPrint(w, indent+2)
}
