package tree

type Uni struct {
	// The universe scope contains all predeclared objects of Hippo.
	scope *Scope
}

var Universe *Uni

func init() {
	Universe = &Uni{
		scope: NewScope(nil),
	}
	Universe.initPredeclared()
}

func (u *Uni) initPredeclared() {
	u.scope.Insert(&Object{
		Kind: "i32",
		Name: "i32",
		Decl: nil,
	})

	u.scope.Insert(&Object{
		Kind: "i64",
		Name: "i64",
		Decl: nil,
	})

	u.scope.Insert(&Object{
		Kind: "f32",
		Name: "f32",
		Decl: nil,
	})

	u.scope.Insert(&Object{
		Kind: "f64",
		Name: "f64",
		Decl: nil,
	})

	u.scope.Insert(&Object{
		Kind: "string",
		Name: "str",
		Decl: nil,
	})
}

func (u *Uni) Lookup(name string) *Object {
	return u.scope.Lookup(name)
}
