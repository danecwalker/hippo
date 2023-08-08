package syntax

import "fmt"

type ObjectKind int

const (
	ObjKindInvalid ObjectKind = iota
	ObjKindFunc
	ObjKindVar
	ObjKindType
	ObjKindConst
	ObjKindField
)

func (ok ObjectKind) String() string {
	switch ok {
	case ObjKindFunc:
		return "func"
	case ObjKindVar:
		return "var"
	case ObjKindType:
		return "type"
	case ObjKindConst:
		return "const"
	case ObjKindField:
		return "field"
	default:
		return "invalid"
	}
}

type Object struct {
	Kind ObjectKind
	Name string
	Decl Node
	Type string
}

func NewObject(kind ObjectKind, name string, decl Node) *Object {
	return &Object{
		Kind: kind,
		Name: name,
		Decl: decl,
	}
}

func (o *Object) String() string {
	n := o.Name
	k := o.Kind
	var v Node
	switch decl := o.Decl.(type) {
	case *VarStatement:
		for i, name := range decl.Names {
			if name.Name == n {
				v = decl.Values[i]
				return fmt.Sprintf("%s %s = %s", k, n, v)
			}
		}
		return fmt.Sprintf("%s %s", k, n)
	default:
		return n
	}
}
