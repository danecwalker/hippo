package ast

import (
	"strings"

	"github.com/danecwalker/hippo/internal/tok"
)

type Node interface{}

type Decl interface {
	Node
	declNode()
}

type Stmt interface {
	Node
	stmtNode()
}

type Expr interface {
	Node
	exprNode()
	String() string
}

type File struct {
	Filename string
	Decls    []Decl
}

type FuncDecl struct {
	FuncPos tok.Position
	Name    *Ident
	Type    *FuncType
	Body    *BlockStmt
}

func (fd *FuncDecl) declNode() {}

type FuncType struct {
	Params  []*Field
	Results []*Field
}

type BlockStmt struct {
	Stmts []Stmt
}

func (bs *BlockStmt) stmtNode() {}

type Field struct {
	Names []*Ident
	Type  *Ident
}

type Ident struct {
	Name    string
	NamePos tok.Position
}

func (i *Ident) exprNode() {}
func (i *Ident) String() string {
	return i.Name
}

type ReturnStmt struct {
	ReturnPos tok.Position
	Results   []Expr
}

func (rs *ReturnStmt) stmtNode() {}

type BinaryExpr struct {
	X  Expr
	Op string
	Y  Expr
}

func (be *BinaryExpr) exprNode() {}
func (be *BinaryExpr) String() string {
	return be.X.String() + " " + be.Op + " " + be.Y.String()
}

type BasicLit struct {
	ValuePos tok.Position
	Value    string
	Kind     string
}

func (bl *BasicLit) exprNode() {}
func (bl *BasicLit) String() string {
	if bl.Kind == "string" {
		return "\\\"" + bl.Value + "\\\""
	}
	return bl.Value
}

type ParenExpr struct {
	Lparen tok.Position
	X      Expr
	Rparen tok.Position
}

func (pe *ParenExpr) exprNode() {}
func (pe *ParenExpr) String() string {
	return "(" + pe.X.String() + ")"
}

type UnaryExpr struct {
	OpPos tok.Position
	Op    string
	X     Expr
}

func (ue *UnaryExpr) exprNode() {}
func (ue *UnaryExpr) String() string {
	return ue.Op + ue.X.String()
}

type CallExpr struct {
	Fun    Expr
	Lparen tok.Position
	Args   []Expr
	Rparen tok.Position
}

func (ce *CallExpr) exprNode() {}
func (ce *CallExpr) String() string {
	return ce.Fun.String() + "(" + ce.argsString() + ")"
}
func (ce *CallExpr) argsString() string {
	var args []string
	for _, arg := range ce.Args {
		args = append(args, arg.String())
	}
	return strings.Join(args, ", ")
}

type ExprStmt struct {
	X Expr
}

func (cs *ExprStmt) stmtNode() {}

type InitDecl struct {
	Names  []*Ident
	Type   *Ident
	Values []Expr
}

func (id *InitDecl) declNode() {}

type AssignStmt struct {
	Lhs []Expr
	Tok string
	Rhs []Expr
}

func (as *AssignStmt) stmtNode() {}

type DeclStmt struct {
	Names  []*Ident
	Type   *Ident
	Values []Expr
}

func (ds *DeclStmt) stmtNode() {}

type RangeStmt struct {
	Names []*Ident
	Iter  Expr
	Body  *BlockStmt
}

func (rs *RangeStmt) stmtNode() {}

type RangeExpr struct {
	Low  Expr
	High Expr
}

func (re *RangeExpr) exprNode() {}
func (re *RangeExpr) String() string {
	return re.Low.String() + ":" + re.High.String()
}
