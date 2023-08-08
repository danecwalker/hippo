package intermediate

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/danecwalker/hippo/internal/syntax"
)

type IR struct {
	Blocks []*Block
	errors []*Error
}

type Inst interface {
	inst()
	pretty(*bytes.Buffer, int)
}

type Block struct {
	Objects map[string]*syntax.Object
	Insts   []Inst
	Parent  *Block
}

func NewIR() *IR {
	return &IR{
		Blocks: []*Block{},
	}
}

func (ir *IR) NewBlock(parent *Block) *Block {
	block := &Block{
		Objects: map[string]*syntax.Object{},
		Parent:  parent,
	}
	ir.Blocks = append(ir.Blocks, block)
	return block
}

func (ir *IR) GetBlock() *Block {
	return ir.Blocks[len(ir.Blocks)-1]
}

func (ir *IR) GetObject(pos *syntax.Position, name string) *syntax.Object {
	b := ir.GetBlock()
	for b != nil {
		if obj, ok := b.Objects[name]; ok {
			return obj
		}
		b = b.Parent
	}
	NewUndefinedObjectError(pos, name)
	return nil
}

func (ir *IR) SetObject(name string, obj *syntax.Object) {
	ir.GetBlock().Objects[name] = obj
}

func (ir *IR) AddInstruction(inst Inst) {
	ir.GetBlock().Insts = append(ir.GetBlock().Insts, inst)
}

func (ir *IR) Generate(prog *syntax.Program) {
	ir.NewBlock(nil)
	ir.addBuiltins()

	for _, stmt := range prog.Statements {
		switch stmt := stmt.(type) {
		case *syntax.VarStatement:
			ir.generateVar(stmt)
		case *syntax.FuncStmt:
			ir.generateFunc(stmt)
		default:
			ir.errors = append(ir.errors, NewDisallowedTopLevelStatementError(stmt.Pos()))
		}
	}
	ir.Pretty()

	if len(ir.errors) > 0 {
		for _, err := range ir.errors {
			fmt.Fprintln(os.Stderr, err)
		}
		os.Exit(1)
	}

}

func (ir *IR) addBuiltins() {
	ir.SetObject("i32", &syntax.Object{
		Name: "i32",
		Kind: syntax.ObjKindType,
		Type: "i32",
	})
}

func (ir *IR) generateVar(stmt *syntax.VarStatement) {
	for i, name := range stmt.Names {
		ir.SetObject(name.Name, name.Obj)
		n := name.Name

		var t *syntax.Object
		if stmt.Type != nil {
			if stmt.Type.Obj != nil {
				t = stmt.Type.Obj
			} else {
				t = ir.GetObject(stmt.Type.NamePos, stmt.Type.Name)
				if t == nil {
					fmt.Fprintln(os.Stderr, NewError(stmt.Type.NamePos, "Undefined type: "+stmt.Type.Name))
				}
			}
		}

		var v Inst
		if stmt.Values[i] != nil {
			v = ir.generateExpr(stmt.Values[i])
		} else {
			v = ir.NewZeroInst(t)
		}

		t = InferType(v)

		ir.AddInstruction(NewAllocInst(n, v, t))
	}
}

func InferType(inst Inst) *syntax.Object {
	switch inst := inst.(type) {
	case *BasicLitInst:
		return inst.Type
	default:
		return nil
	}
}

func (ir *IR) generateExpr(expr syntax.Expression) Inst {
	switch expr := expr.(type) {
	case *syntax.Identifier:
		return ir.generateIdent(expr)
	case *syntax.BasicLit:
		return ir.generateBasicLit(expr)
	case *syntax.BinaryExpr:
		return ir.generateBinaryExpr(expr)
	case *syntax.AssignmentExpr:
		return ir.generateAssignmentExpr(expr)
	default:
		ir.errors = append(ir.errors, NewUnexpectedExpr(expr.Pos()))
		return nil
	}
}

func (ir *IR) generateAssignmentExpr(expr *syntax.AssignmentExpr) Inst {
	l := ir.generateExpr(expr.Lhs)
	r := ir.generateExpr(expr.Rhs)

	switch expr.Op.Type {
	case syntax.TokenAssign:
		return NewAssignInst(l, r)
	default:
		ir.errors = append(ir.errors, NewUnexpectedExpr(expr.Pos()))
		return nil
	}
}

func (ir *IR) generateIdent(ident *syntax.Identifier) Inst {
	obj := ir.GetObject(ident.Pos(), ident.Name)
	if obj == nil {
		ir.errors = append(ir.errors, NewUndefinedObjectError(ident.Pos(), ident.Name))
		return nil
	}
	return NewIdentInst(obj.Name)
}

func (ir *IR) generateBasicLit(lit *syntax.BasicLit) Inst {
	return NewBasicLitInst(lit.Value, ir.GetObject(lit.Pos(), lit.Kind))
}

func (ir *IR) generateBinaryExpr(expr *syntax.BinaryExpr) Inst {
	l := ir.generateExpr(expr.X)
	r := ir.generateExpr(expr.Y)

	switch expr.Op.Type {
	case syntax.TokenPlus:
		return NewAddInst(l, r)
	case syntax.TokenMinus:
		return NewSubInst(l, r)
	case syntax.TokenStar:
		return NewMulInst(l, r)
	case syntax.TokenSlash:
		return NewDivInst(l, r)
	default:
		ir.errors = append(ir.errors, NewUnexpectedExpr(expr.Pos()))
		return nil
	}
}

func (ir *IR) generateFunc(stmt *syntax.FuncStmt) {
	ir.NewBlock(ir.GetBlock())
	for _, param := range stmt.Type.Params {
		for _, name := range param.Names {
			ir.SetObject(name.Name, name.Obj)
		}
	}

	ir.generateBlockStmt(stmt.Body)
}

func (ir *IR) generateStmt(stmt syntax.Statement) {
	switch stmt := stmt.(type) {
	case *syntax.VarStatement:
		ir.generateVar(stmt)
	case *syntax.FuncStmt:
		ir.generateFunc(stmt)
	case *syntax.ReturnStmt:
		ir.generateReturn(stmt)
	case *syntax.ForRangeStmt:
		ir.generateForRange(stmt)
	case *syntax.BlockStmt:
		ir.generateBlockStmt(stmt)
	case *syntax.ExpressionStmt:
		ir.AddInstruction(ir.generateExpr(stmt.X))
	default:
		ir.errors = append(ir.errors, NewUnexpectedStmt(stmt.Pos()))
	}
}

func (ir *IR) generateBlockStmt(stmt *syntax.BlockStmt) {
	for _, stmt := range stmt.Stmts {
		ir.generateStmt(stmt)
	}
}

func (ir *IR) generateForRange(stmt *syntax.ForRangeStmt) {
	ir.generateVar(stmt.Key.Obj.Decl.(*syntax.VarStatement))
	ir.NewBlock(ir.GetBlock())

	ir.generateStmt(stmt.Body)

	name := syntax.NewIdentifier(nil, stmt.Key.Obj.Name)
	ir.AddInstruction(ir.generateAssignmentExpr(syntax.NewAssignmentExpr(
		name,
		syntax.NewToken(syntax.TokenAssign, "=", nil),
		syntax.NewBinaryExpr(
			name,
			syntax.NewToken(syntax.TokenPlus, "+", nil),
			syntax.NewBasicLit(nil, "i32", "1"),
		),
	)))
}

func (ir *IR) generateReturn(stmt *syntax.ReturnStmt) {
	var r Inst
	if stmt.Result != nil {
		r = ir.generateExpr(stmt.Result)
	}

	ir.AddInstruction(NewReturnInst(r))
}

func (ir *IR) Pretty() {
	var w bytes.Buffer

	for i, b := range ir.Blocks {
		b.pretty(&w, i, 0)
	}

	fmt.Println(w.String())
}

func addIndent(w *bytes.Buffer, indent int) {
	w.WriteString(strings.Repeat("  ", indent))
}

func (b *Block) pretty(w *bytes.Buffer, count int, indent int) {
	w.WriteString(fmt.Sprintf("b%d:\n", count))
	addIndent(w, indent+1)
	w.WriteString("Objects:\n")
	for _, o := range b.Objects {
		addIndent(w, indent+2)
		w.WriteString(fmt.Sprintf("<%s> %s\n", o.Name, o.Kind))
	}
	addIndent(w, indent+1)
	w.WriteString("Instructions:\n")
	for _, inst := range b.Insts {
		addIndent(w, indent+2)
		inst.pretty(w, indent)
	}
}
