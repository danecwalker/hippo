package intermediate

import (
	"fmt"
	"os"

	"github.com/danecwalker/hippo/internal/syntax"
)

func (ir *IR) NewZeroInst(t *syntax.Object) Inst {
	if t.Kind != syntax.ObjKindType {
		fmt.Fprintln(os.Stderr, NewError(t.Decl.Pos(), "Cannot create zero value for non-type: "+t.Name))
		return nil
	}

	switch t.Type {
	case "i32":
		return ir.NewIntZeroInst()
	default:
		fmt.Fprintln(os.Stderr, NewError(t.Decl.Pos(), "Cannot create zero value for type: "+t.Name))
		return nil
	}
}

func (ir *IR) NewIntZeroInst() Inst {
	return NewBasicLitInst("0", ir.GetObject(nil, "i32"))
}
