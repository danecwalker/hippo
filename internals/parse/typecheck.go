package parse

import (
	"fmt"

	"github.com/danecwalker/hippo/internals/errors"
	"github.com/danecwalker/hippo/internals/tree"
)

type Checker struct {
	typeMap map[string]*tree.Object
	*errors.ErrorHandler
}

func Typecheck(file *tree.File, eh *errors.ErrorHandler) map[string]*tree.Object {
	t := &Checker{
		typeMap:      make(map[string]*tree.Object),
		ErrorHandler: eh,
	}
	// Walk the tree to find identifiers with nil OBJ
	for _, decl := range file.Decls {
		t.check(decl)
	}

	return t.typeMap
}

func (t *Checker) check(node tree.Node) {
	switch node := node.(type) {
	case *tree.StoreDecl:
		for _, spec := range node.Specs {
			t.check(spec)
		}
	case *tree.ValueSpec:
		t.checkValueSpec(node)
	}
}

func (t *Checker) getKind(node tree.Node, args ...string) string {
	switch node := node.(type) {
	case *tree.Ident:
		if node.Obj.Decl == nil {
			return node.Obj.Kind
		}
		return t.getKind(node.Obj.Decl, node.Name)
	case *tree.ValueSpec:
		if len(args) < 1 {
			t.ErrorWithLoc(errors.ERROR, "invalid number of arguments for getKind", node.Names[0].NamePos.String())
			t.ShouldExit()
		}
		for i, name := range node.Names {
			if name.Name == args[0] {
				return t.getKind(node.Types[i])
			}
		}

		t.ErrorWithLoc(errors.ERROR, "could not find identifier `%s`", node.Names[0].NamePos.String(), args[0])
		t.ShouldExit()
		return ""
	case *tree.BasicLit:
		return node.Kind
	default:
		return ""
	}
}

var NumPrecedence = map[string]int{
	"i32": 1,
	"i64": 2,
	"u32": 3,
	"u64": 4,
	"f32": 5,
	"f64": 6,
}

func (t *Checker) checkValueSpec(node *tree.ValueSpec) {
	for i, name := range node.Names {
		typ := t.getKind(name)
		vtyp := t.getKind(node.Values[i])
		if typ != vtyp {
			// type coercion
			fmt.Println(NumPrecedence[typ])
			if NumPrecedence[typ] > NumPrecedence[vtyp] && NumPrecedence[typ] > 0 && NumPrecedence[vtyp] > 0 {
				if blit, ok := node.Values[i].(*tree.BasicLit); ok {
					blit.Kind = typ
					return
				}
			}

			t.ErrorWithLoc(errors.ERROR, "cannot assign value of type `%s` to variable of type `%s`", node.Values[i].Pos().String(), vtyp, typ)
		}
	}
}
