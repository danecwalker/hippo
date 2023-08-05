package tree

import (
	"fmt"
)

// Print ast in form of
// File
// .	Decls
// .	.	Decl
// .	.	.	VarDecl
// .	.	.	.	Ident
// .	.	.	.	.	foo
// .	.	.	.	Ident
// .	.	.	.	.	bar

var lineCount int = 0

func Print(file *File) {
	lineCount = 0
	printFile(file, 0)
}

func printPrefix(depth int) {
	lineCount++
	prefix := fmt.Sprintf("%4d  ", lineCount)
	for i := 0; i < depth; i++ {
		prefix += ".  "
	}
	fmt.Print(prefix)
}

func printFile(file *File, depth int) {
	printPrefix(depth)
	fmt.Println("*tree.File {")
	printDecls(file.Decls, depth+1)

	printPrefix(depth)
	fmt.Println("}")
}

func printDecls(decls []Decl, depth int) {
	printPrefix(depth)
	fmt.Printf("Decls: []tree.Decl (len = %d) [\n", len(decls))
	for i, decl := range decls {
		printDecl(i, decl, depth+1)
	}
	printPrefix(depth)
	fmt.Println("]")
}

func printDecl(i int, decl Decl, depth int) {
	switch decl := decl.(type) {
	case *StoreDecl:
		printStoreDecl(i, decl, depth)
	case *FuncDecl:
		printFuncDecl(i, decl, depth)
	}
}

func printFuncDecl(i int, decl *FuncDecl, depth int) {
	printPrefix(depth)
	fmt.Printf("%d: *tree.FuncDecl {\n", i)
	printPrefix(depth + 1)
	fmt.Print("Name: ")
	printIdent(-1, decl.Name, depth+1)
	printFuncType(decl.Type, depth+1)
	printFuncBody(decl.Body, depth+1)

	printPrefix(depth)
	fmt.Println("}")
}

func printFuncBody(body *BlockStmt, depth int) {
	printPrefix(depth)
	fmt.Println("Body: *tree.BlockStmt {")
	printPrefix(depth + 1)
	fmt.Printf("LBrace: %s\n", body.Lbrace)
	printPrefix(depth + 1)
	fmt.Println("Stmts: []tree.Stmt {")
	for i, stmt := range body.List {
		printStmt(i, stmt, depth+2)
	}
	printPrefix(depth + 1)
	fmt.Println("}")
	printPrefix(depth + 1)
	fmt.Printf("RBrace: %s\n", body.Rbrace)
	printPrefix(depth)
	fmt.Println("}")
}

func printStmt(i int, stmt Stmt, depth int) {
	switch stmt := stmt.(type) {
	case *DeclStmt:
		printDeclStmt(i, stmt, depth)
	case *ExprStmt:
		printExprStmt(i, stmt, depth)
	}
}

func printDeclStmt(i int, decl *DeclStmt, depth int) {
	printPrefix(depth)
	fmt.Printf("%d: *tree.DeclStmt {\n", i)
	printPrefix(depth + 1)
	fmt.Println("SymPos:", decl.Sym.Loc)
	printPrefix(depth + 1)
	fmt.Println("Sym:", decl.Sym.Type)
	printSpecs(decl.Specs, depth+1)
	printPrefix(depth)
	fmt.Println("}")
}

func printExprStmt(i int, stmt *ExprStmt, depth int) {
	printPrefix(depth)
	fmt.Printf("%d: *tree.ExprStmt {\n", i)
	printPrefix(depth + 1)
	fmt.Print("X: ")
	printExpr(i, stmt.X, depth+1)
	printPrefix(depth)
	fmt.Println("}")
}

func printFuncType(ft *FuncType, depth int) {
	printPrefix(depth)
	fmt.Println("Type: *tree.FuncType {")
	printPrefix(depth + 1)
	fmt.Print("Params: ")
	printFieldList(ft.Params, depth+2)
	printPrefix(depth + 1)
	fmt.Print("Results: ")
	printFieldList(ft.Results, depth+2)
	printPrefix(depth)
	fmt.Println("}")
}

func printFieldList(fl []*Field, depth int) {
	fmt.Printf("[]*tree.Field (len = %d) [", len(fl))
	if len(fl) == 0 {
		fmt.Println("]")
		return
	}
	fmt.Println()

	for i, field := range fl {
		printField(i, field, depth)
	}
	printPrefix(depth - 1)
	fmt.Println("]")
}

func printField(i int, field *Field, depth int) {
	printPrefix(depth)
	fmt.Printf("%d: *tree.Field {\n", i)
	printPrefix(depth + 1)
	if field.Name != nil {
		fmt.Print("Name: ")
		printIdent(-1, field.Name, depth+1)
	} else {
		fmt.Println("Name: nil")
	}
	printPrefix(depth + 1)
	fmt.Print("Type: ")
	printIdent(-1, field.Type, depth+1)
	printPrefix(depth)
	fmt.Println("}")
}

func printStoreDecl(i int, decl *StoreDecl, depth int) {
	printPrefix(depth)
	fmt.Printf("%d: *tree.StoreDecl {\n", i)
	printPrefix(depth + 1)
	fmt.Println("SymPos:", decl.Sym.Loc)
	printPrefix(depth + 1)
	fmt.Println("Sym:", decl.Sym.Type)
	printSpecs(decl.Specs, depth+1)
	printPrefix(depth)
	fmt.Println("}")
}

func printSpecs(specs []Spec, depth int) {
	printPrefix(depth)
	fmt.Printf("Specs: []tree.Spec (len = %d) [\n", len(specs))
	for i, spec := range specs {
		printSpec(i, spec, depth+1)
	}
	printPrefix(depth)
	fmt.Println("]")
}

func printSpec(i int, spec Spec, depth int) {
	switch spec := spec.(type) {
	case *ValueSpec:
		printValueSpec(i, spec, depth)
	}
}

func printValueSpec(i int, spec *ValueSpec, depth int) {
	printPrefix(depth)
	fmt.Printf("%d: *tree.ValueSpec {\n", i)
	printPrefix(depth + 1)
	fmt.Printf("Names: []*tree.Ident (len = %d) [\n", len(spec.Names))
	for j, name := range spec.Names {
		printExpr(j, name, depth+2)
	}
	printPrefix(depth + 1)
	fmt.Println("]")

	printPrefix(depth + 1)
	fmt.Printf("Types: []*tree.Ident (len = %d) [\n", len(spec.Types))
	for j, typ := range spec.Types {
		printExpr(j, typ, depth+2)
	}
	printPrefix(depth + 1)
	fmt.Println("]")
	printPrefix(depth + 1)
	fmt.Printf("Values: []tree.Expr (len = %d) [\n", len(spec.Values))
	for j, value := range spec.Values {
		printExpr(j, value, depth+2)
	}
	printPrefix(depth + 1)
	fmt.Println("]")
	printPrefix(depth)
	fmt.Println("}")
}

func printExpr(i int, expr Expr, depth int) {
	switch expr := expr.(type) {
	case *Ident:
		printIdent(i, expr, depth)
	case *BasicLit:
		printBasicLit(i, expr, depth)
	}
}

func printIdent(i int, ident *Ident, depth int) {
	if i >= 0 {
		printPrefix(depth)
		fmt.Printf("%d: *tree.Ident {\n", i)
		printPrefix(depth + 1)
		fmt.Println("NamePos:", ident.NamePos)
		printPrefix(depth + 1)
		fmt.Println("Name:", ident.Name)
		printPrefix(depth + 1)
		if ident.Obj != nil {
			fmt.Printf("Obj: *tree.Object {\n")
			printObj(ident.Obj, depth+2)
			printPrefix(depth + 1)
			fmt.Println("}")
		} else {
			fmt.Println("Obj: nil")
		}
		printPrefix(depth)
		fmt.Println("}")
	} else {
		fmt.Print("*tree.Ident {\n")
		printPrefix(depth + 1)
		fmt.Println("NamePos:", ident.NamePos)
		printPrefix(depth + 1)
		fmt.Println("Name:", ident.Name)
		printPrefix(depth + 1)
		if ident.Obj != nil {
			fmt.Printf("Obj: *tree.Object {\n")
			printObj(ident.Obj, depth+2)
			printPrefix(depth + 1)
			fmt.Println("}")
		} else {
			fmt.Println("Obj: nil")
		}
		printPrefix(depth)
		fmt.Println("}")
	}
}

func printBasicLit(i int, lit *BasicLit, depth int) {
	printPrefix(depth)
	fmt.Printf("%d: *tree.BasicLit {\n", i)
	printPrefix(depth + 1)
	fmt.Println("ValuePos:", lit.ValuePos)
	printPrefix(depth + 1)
	fmt.Println("Kind:", lit.Kind)
	printPrefix(depth + 1)
	fmt.Println("Value:", lit.Value)
	printPrefix(depth)
	fmt.Println("}")
}

func printObj(obj *Object, depth int) {
	printPrefix(depth)
	fmt.Println("Kind:", obj.Kind)
	printPrefix(depth)
	fmt.Println("Name:", obj.Name)
	printPrefix(depth)
	switch obj.Decl.(type) {
	case *StoreDecl:
		fmt.Println("Decl: *tree.StoreDecl { ... }")
	case *ValueSpec:
		fmt.Println("Decl: *tree.ValueSpec { ... }")
	case *FuncDecl:
		fmt.Println("Decl: *tree.FuncDecl { ... }")
	case *Ident:
		fmt.Println("Decl: *tree.Ident { ... }")
	case *BasicLit:
		fmt.Println("Decl: *tree.BasicLit { ... }")
	case *Field:
		fmt.Println("Decl: *tree.Field { ... }")
	default:
		fmt.Println("Decl:", obj.Decl)
	}
}
