package tree

type Scope struct {
	Parent  *Scope
	Objects map[string]*Object
}

func NewScope(parent *Scope) *Scope {
	return &Scope{
		Parent:  parent,
		Objects: make(map[string]*Object),
	}
}

func (s *Scope) Lookup(name string) *Object {
	if obj, ok := s.Objects[name]; ok {
		return obj
	}
	if s.Parent != nil {
		return s.Parent.Lookup(name)
	}
	return nil
}

func (s *Scope) Insert(obj *Object) {
	s.Objects[obj.Name] = obj
}

func (s *Scope) Remove(name string) {
	delete(s.Objects, name)
}
