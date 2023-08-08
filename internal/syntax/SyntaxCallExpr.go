package syntax

import "bytes"

type CallExpr struct {
	Func   *Identifier
	Lparen *Position
	Args   Expression
	Rparen *Position
}

func NewCallExpr(func_ *Identifier, lparen *Position, args Expression, rparen *Position) *CallExpr {
	return &CallExpr{
		Func:   func_,
		Lparen: lparen,
		Args:   args,
		Rparen: rparen,
	}
}

func (c *CallExpr) expressionNode() {}
func (c *CallExpr) Pos() *Position {
	return c.Func.Pos()
}

func (c *CallExpr) PrettyPrint(w *bytes.Buffer, indent int) {
	addIndent(w, indent)
	w.WriteString("CallExpr:\n")
	addIndent(w, indent+1)
	w.WriteString("Func:\n")
	c.Func.PrettyPrint(w, indent+2)
	addIndent(w, indent+1)
	w.WriteString("Args:\n")
	c.Args.PrettyPrint(w, indent+2)
}
