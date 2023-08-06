package flow

import (
	"fmt"
	"os"
	"strings"

	"github.com/danecwalker/hippo/internal/ast"
)

type CFG struct {
	filename  string
	entities  []*Entity
	relations []*Relation
}

type Entity struct {
	Alias string
	Stmts []string
}

type Relation struct {
	From string
	To   string
}

func NewCFG() *CFG {
	return &CFG{
		entities:  make([]*Entity, 0),
		relations: make([]*Relation, 0),
	}
}

func (c *CFG) AddEntity(alias string, stmts []string) *Entity {
	e := &Entity{
		Alias: alias,
		Stmts: stmts,
	}
	c.entities = append(c.entities, e)
	return e
}

func (c *CFG) AddRelation(from string, to string) {
	c.relations = append(c.relations, &Relation{
		From: from,
		To:   to,
	})
}

func (c *CFG) GetEntity() *Entity {
	return c.entities[len(c.entities)-1]
}

func (c *CFG) GenCFG(file *ast.File) {
	c.filename = file.Filename
	for _, decl := range file.Decls {
		switch decl := decl.(type) {
		case *ast.FuncDecl:
			c.genFunc(decl)
		}
	}

	c.makeDot()
}

func (c *CFG) genFunc(decl *ast.FuncDecl) {
	stmts := make([]string, 0)
	e := c.AddEntity(decl.Name.String(), stmts)
	for _, stmt := range decl.Body.Stmts {
		switch stmt := stmt.(type) {
		case *ast.ReturnStmt:
			stmts = append(stmts, c.genReturnStmt(stmt))
		case *ast.ExprStmt:
			stmts = append(stmts, c.genCallStmt(stmt))
		}
	}
	e.Stmts = stmts
}

func (c *CFG) genReturnStmt(stmt *ast.ReturnStmt) string {
	s := "return "

	for _, result := range stmt.Results {
		s += result.String()
	}

	return s
}

func (c *CFG) genCallStmt(stmt *ast.ExprStmt) string {
	s := stmt.X.String()
	c.AddRelation(c.GetEntity().Alias, stmt.X.(*ast.CallExpr).Fun.(*ast.Ident).String())
	return s
}

func (c *CFG) makeDot() {
	lines := make([]string, 0)
	lines = append(lines, fmt.Sprintf("digraph \"%s\" {", c.filename))

	for _, e := range c.entities {
		lines = append(lines, fmt.Sprintf("\t%s [label=\"<%s>\n%s\", shape=\"box\"]", e.Alias, e.Alias, strings.Join(e.Stmts, "\n")))
	}

	for _, r := range c.relations {
		lines = append(lines, fmt.Sprintf("\t%s -> %s", r.From, r.To))
	}

	lines = append(lines, "}")

	// remove .x from filename
	fname := c.filename[:len(c.filename)-2]
	f, err := os.Create(fmt.Sprintf("%s.dot", fname))
	if err != nil {
		panic(err)
	}

	for _, line := range lines {
		f.WriteString(fmt.Sprintf("%s\n", line))
	}
}
