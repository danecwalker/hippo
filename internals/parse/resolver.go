package parse

import (
	"github.com/danecwalker/hippo/internals/errors"
	"github.com/danecwalker/hippo/internals/tree"
)

type Resolver struct {
	unresolved []*tree.Ident
	scope      *tree.Scope
	*errors.ErrorHandler
}

func FindUnresolved(file *tree.File, eh *errors.ErrorHandler) []*tree.Ident {
	r := &Resolver{
		unresolved:   make([]*tree.Ident, 0),
		scope:        tree.NewScope(nil),
		ErrorHandler: eh,
	}
	// Walk the tree to find identifiers with nil OBJ
	for _, decl := range file.Decls {
		r.resolve(decl, r.scope)
	}

	return r.unresolved
}

func (r *Resolver) resolve(node tree.Node, scope *tree.Scope) {
	switch node := node.(type) {
	case *tree.StoreDecl:
		for _, spec := range node.Specs {
			r.resolve(spec, scope)
		}
	case *tree.ValueSpec:
		for _, name := range node.Names {
			scope.Insert(name.Obj)
			r.resolve(name, scope)
		}
		for _, typ := range node.Types {
			r.resolve(typ, scope)
		}
		for _, val := range node.Values {
			r.resolve(val, scope)
		}
	case *tree.Ident:
		if node.Obj == nil {
			if node.Obj = scope.Lookup(node.Name); node.Obj == nil {
				// Is ident a type?
				if node.Obj = tree.Universe.Lookup(node.Name); node.Obj == nil {
					r.ErrorWithLoc(errors.ERROR, "unresolved identifier `%s`", node.NamePos.String(), node.Name)
					r.unresolved = append(r.unresolved, node)
				}
			}
		}
	case *tree.FuncDecl:
		scope.Insert(node.Name.Obj)
		funcScope := tree.NewScope(scope)
		r.resolve(node.Type, funcScope)
		r.resolve(node.Body, funcScope)
	case *tree.FuncType:
		for _, param := range node.Params {
			scope.Insert(param.Name.Obj)
			r.resolve(param, scope)
		}
		for _, result := range node.Results {
			r.resolve(result, scope)
		}
	case *tree.BlockStmt:
		for _, stmt := range node.List {
			r.resolve(stmt, scope)
		}
	case *tree.ExprStmt:
		r.resolve(node.X, scope)
	case *tree.AssignStmt:
		for _, lhs := range node.Lhs {
			r.resolve(lhs, scope)
		}
		for _, rhs := range node.Rhs {
			r.resolve(rhs, scope)
		}
	case *tree.DeclStmt:
		for _, spec := range node.Specs {
			r.resolve(spec, scope)
		}
	}
}
