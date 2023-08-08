package parse

import "github.com/danecwalker/hippo/internal/syntax"

type Scope struct {
	Parent  *Scope
	Objects map[string]*syntax.Object
}

func NewScope(parent *Scope) *Scope {
	return &Scope{
		Parent:  parent,
		Objects: make(map[string]*syntax.Object),
	}
}

func (s *Scope) Insert(obj *syntax.Object) {
	s.Objects[obj.Name] = obj
}

func (s *Scope) Lookup(name string) *syntax.Object {
	if obj, ok := s.Objects[name]; ok {
		return obj
	}

	if s.Parent != nil {
		return s.Parent.Lookup(name)
	}

	return nil
}
