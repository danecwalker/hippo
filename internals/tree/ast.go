package tree

import (
	"github.com/danecwalker/hippo/internals/symbol"
	"github.com/danecwalker/hippo/internals/utils"
)

type Node interface {
	Pos() utils.Location
}

type Decl interface {
	Node
	declNode()
}

type Expr interface {
	Node
	exprNode()
}

type File struct {
	Decls []Decl
}

type FuncDecl struct {
	Name *Ident
	Type *FuncType
	Body *BlockStmt
}

func (c *FuncDecl) declNode() {}
func (c *FuncDecl) Pos() utils.Location {
	return c.Type.FuncPos
}

type FuncType struct {
	FuncPos utils.Location
	Params  []*Field
	Results []*Field
}

func (c *FuncType) exprNode() {}
func (c *FuncType) Pos() utils.Location {
	return c.FuncPos
}

type BlockStmt struct {
	Lbrace utils.Location
	List   []Stmt
	Rbrace utils.Location
}

func (c *BlockStmt) exprNode() {}
func (c *BlockStmt) Pos() utils.Location {
	return c.Lbrace
}

type Field struct {
	Name *Ident
	Type *Ident
}

func (c *Field) exprNode() {}
func (c *Field) Pos() utils.Location {
	return c.Name.NamePos
}

type Stmt interface {
	Node
	stmtNode()
}

type AssignStmt struct {
	Lhs    []Expr
	TokPos utils.Location
	Tok    symbol.Symbol
	Rhs    []Expr
}

func (c *AssignStmt) exprNode() {}
func (c *AssignStmt) Pos() utils.Location {
	return c.TokPos
}

type StoreDecl struct {
	Sym    symbol.Symbol
	SymPos utils.Location

	Specs []Spec
}

func (c *StoreDecl) declNode() {}
func (c *StoreDecl) Pos() utils.Location {
	return c.SymPos
}

type DeclStmt struct {
	Sym    symbol.Symbol
	SymPos utils.Location

	Specs []Spec
}

func (c *DeclStmt) stmtNode() {}
func (c *DeclStmt) Pos() utils.Location {
	return c.SymPos
}

type Spec interface {
	Node
	specNode()
}

type ValueSpec struct {
	Names  []*Ident
	Types  []*Ident
	Values []Expr
}

func (c *ValueSpec) specNode() {}
func (c *ValueSpec) Pos() utils.Location {
	return c.Names[0].NamePos
}

type Ident struct {
	NamePos utils.Location
	Name    string
	Obj     *Object
}

func (c *Ident) exprNode() {}
func (c *Ident) Pos() utils.Location {
	return c.NamePos
}

type BasicLit struct {
	ValuePos utils.Location
	Kind     string
	Value    string
}

func (c *BasicLit) exprNode() {}
func (c *BasicLit) Pos() utils.Location {
	return c.ValuePos
}

type Object struct {
	Kind string
	Name string
	Decl Node
}

type ExprStmt struct {
	X Expr
}

func (c *ExprStmt) stmtNode() {}
func (c *ExprStmt) Pos() utils.Location {
	return c.X.Pos()
}
