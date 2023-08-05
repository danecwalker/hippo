package parse

import (
	"strings"

	"github.com/danecwalker/hippo/internals/errors"
	"github.com/danecwalker/hippo/internals/lex"
	"github.com/danecwalker/hippo/internals/symbol"
	"github.com/danecwalker/hippo/internals/tree"
)

type prefixParseFn func() tree.Expr
type infixParseFn func(tree.Expr) tree.Expr

type Parser struct {
	*errors.ErrorHandler
	*lex.Lexer
	cSym symbol.Symbol
	pSym symbol.Symbol

	prefixParseFns map[symbol.SymbolType]prefixParseFn
	infixParseFns  map[symbol.SymbolType]infixParseFn
}

func (p *Parser) registerPrefix(t symbol.SymbolType, fn prefixParseFn) {
	p.prefixParseFns[t] = fn
}

func (p *Parser) registerInfix(t symbol.SymbolType, fn infixParseFn) {
	p.infixParseFns[t] = fn
}

func (p *Parser) next() {
	p.cSym = p.pSym
	p.pSym = p.NextSymbol()
}

func (p *Parser) expect(t symbol.SymbolType) bool {
	if p.pSym.Type != t {
		p.ErrorWithLoc(errors.ERROR, "expected %s, got `%s` (type %s)", p.pSym.Loc.String(), t, p.pSym.Lit, p.pSym.Type)
		p.ShouldExit()
		return false
	}
	p.next()
	return true
}

func (p *Parser) accept(t symbol.SymbolType) bool {
	if p.pSym.Type != t {
		return false
	}
	p.next()
	return true
}

func ParseFile(l *lex.Lexer, eh *errors.ErrorHandler) *tree.File {
	p := &Parser{
		ErrorHandler: eh,
		Lexer:        l,
	}
	p.prefixParseFns = make(map[symbol.SymbolType]prefixParseFn)
	p.registerPrefix(symbol.IDENT, p.parseIdent)
	p.registerPrefix(symbol.INT, p.parseIntegerLiteral)
	p.registerPrefix(symbol.FLOAT, p.parseFloatLiteral)
	p.registerPrefix(symbol.STRING, p.parseStringLiteral)

	p.infixParseFns = make(map[symbol.SymbolType]infixParseFn)
	p.registerInfix(symbol.ASSIGN, p.parseAssign)

	p.next()
	p.next()
	return p.parseFile()
}

func (p *Parser) parseFile() *tree.File {
	return &tree.File{
		Decls: p.parseDecls(),
	}
}

func (p *Parser) parseDecls() []tree.Decl {
	var decls []tree.Decl
	for p.cSym.Type != symbol.EOF {
		decls = append(decls, p.parseDecl())
	}
	return decls
}

func (p *Parser) parseDecl() tree.Decl {
	switch p.cSym.Type {
	case symbol.CONST, symbol.VAR:
		return p.parseStoreDecl()
	case symbol.FN:
		return p.parseFnDecl()
	default:
		p.ErrorWithLoc(errors.ERROR, "unexpected declaration", p.cSym.Loc.String())
		p.ShouldExit()
		return nil
	}
}

func (p *Parser) parseStoreDecl() *tree.StoreDecl {
	decl := &tree.StoreDecl{
		Sym:    p.cSym,
		SymPos: p.cSym.Loc,
		Specs:  make([]tree.Spec, 0),
	}

	// Parse new ValueSpec
	spec := &tree.ValueSpec{
		Names:  make([]*tree.Ident, 0),
		Types:  make([]*tree.Ident, 0),
		Values: make([]tree.Expr, 0),
	}

	if !p.expect(symbol.IDENT) {
		return nil
	}

	name := &tree.Ident{
		NamePos: p.cSym.Loc,
		Name:    p.cSym.Lit,
		Obj: &tree.Object{
			Kind: "VAR",
			Name: p.cSym.Lit,
			Decl: spec,
		},
	}

	if !p.expect(symbol.COLON) {
		return nil
	}

	if !p.expect(symbol.IDENT) {
		return nil
	}

	typ := &tree.Ident{
		NamePos: p.cSym.Loc,
		Name:    p.cSym.Lit,
		Obj:     nil,
	}

	spec.Names = append(spec.Names, name)
	spec.Types = append(spec.Types, typ)

	for p.accept(symbol.COMMA) {
		if !p.expect(symbol.IDENT) {
			return nil
		}

		name := &tree.Ident{
			NamePos: p.cSym.Loc,
			Name:    p.cSym.Lit,
			Obj: &tree.Object{
				Kind: "VAR",
				Name: p.cSym.Lit,
				Decl: spec,
			},
		}

		if !p.expect(symbol.COLON) {
			return nil
		}

		if !p.expect(symbol.IDENT) {
			return nil
		}

		typ := &tree.Ident{
			NamePos: p.cSym.Loc,
			Name:    p.cSym.Lit,
			Obj:     nil,
		}

		spec.Names = append(spec.Names, name)
		spec.Types = append(spec.Types, typ)
	}

	// Accept assignment operator
	if !p.expect(symbol.ASSIGN) {
		return nil
	}

	p.next()

	spec.Values = append(spec.Values, p.parseExprList()...)

	if len(spec.Types) > len(spec.Names) {
		p.Error(errors.ERROR, strings.Join([]string{
			"variable assignment mismatch",
			"the number of types must be less than or equal to the number of variables",
		}, "\n\t"))
		p.ShouldExit()
		return nil
	}

	if len(spec.Values) != len(spec.Names) {
		p.Error(errors.ERROR, strings.Join([]string{
			"variable assignment mismatch",
			"expression count does not match variable count",
		}, "\n\t"))
		p.ShouldExit()
		return nil
	}

	// Add spec to list of specs
	decl.Specs = append(decl.Specs, spec)

	p.next()
	return decl
}

func (p *Parser) parseFnDecl() *tree.FuncDecl {
	fn := &tree.FuncDecl{
		Type: &tree.FuncType{
			FuncPos: p.cSym.Loc,
		},
	}

	if !p.expect(symbol.IDENT) {
		return nil
	}

	fn.Name = &tree.Ident{
		NamePos: p.cSym.Loc,
		Name:    p.cSym.Lit,
		Obj: &tree.Object{
			Kind: "FUNC",
			Name: p.cSym.Lit,
			Decl: fn,
		},
	}

	if !p.expect(symbol.LPAREN) {
		return nil
	}

	for p.accept(symbol.IDENT) {
		field := &tree.Field{}
		field.Name = &tree.Ident{
			NamePos: p.cSym.Loc,
			Name:    p.cSym.Lit,
			Obj: &tree.Object{
				Kind: "VAR",
				Name: p.cSym.Lit,
				Decl: field,
			},
		}

		if !p.expect(symbol.COLON) {
			return nil
		}

		if !p.expect(symbol.IDENT) {
			return nil
		}

		field.Type = &tree.Ident{
			NamePos: p.cSym.Loc,
			Name:    p.cSym.Lit,
			Obj:     nil,
		}

		fn.Type.Params = append(fn.Type.Params, field)

		p.accept(symbol.COMMA)
	}

	if !p.expect(symbol.RPAREN) {
		return nil
	}

	for p.accept(symbol.IDENT) {
		field := &tree.Field{}
		field.Name = nil

		field.Type = &tree.Ident{
			NamePos: p.cSym.Loc,
			Name:    p.cSym.Lit,
			Obj:     nil,
		}

		fn.Type.Results = append(fn.Type.Results, field)

		p.accept(symbol.COMMA)
	}

	if !p.expect(symbol.LBRACE) {
		return nil
	}

	fn.Body = p.parseBlockStmt()

	fn.Body.Rbrace = p.cSym.Loc
	p.next()

	return fn
}

func (p *Parser) parseBlockStmt() *tree.BlockStmt {
	block := &tree.BlockStmt{
		Lbrace: p.cSym.Loc,
		List:   []tree.Stmt{},
	}

	p.next()
	for {
		if p.cSym.Type == symbol.EOF || p.cSym.Type == symbol.RBRACE {
			break
		}
		block.List = append(block.List, p.parseStmt())
	}

	return block
}

func (p *Parser) parseStmt() tree.Stmt {
	switch p.cSym.Type {
	case symbol.VAR, symbol.CONST:
		return p.parseDeclStmt()
	default:
		return p.parseSimpleStmt()
	}
}

func (p *Parser) parseDeclStmt() tree.Stmt {
	decl := p.parseStoreDecl()
	return &tree.DeclStmt{
		Sym:    decl.Sym,
		SymPos: decl.SymPos,
		Specs:  decl.Specs,
	}
}

func (p *Parser) parseSimpleStmt() tree.Stmt {
	return &tree.ExprStmt{
		X: p.parseExpr(LOWEST),
	}
}

// ----------------------  Expression Parsing  ----------------------

func (p *Parser) parseAssign(lhs tree.Expr) tree.Expr {
	assign := &tree.AssignStmt{
		Lhs: []tree.Expr{lhs},
	}

	assign.TokPos = p.cSym.Loc
	assign.Tok = p.cSym

	assign.Rhs = append(assign.Rhs, p.parseExpr(LOWEST))

	return assign
}

func (p *Parser) parseIdent() tree.Expr {
	ident := &tree.Ident{
		NamePos: p.cSym.Loc,
		Name:    p.cSym.Lit,
		Obj:     nil,
	}
	return ident
}

func (p *Parser) precedence() Precedence {
	if p, ok := precedences[p.pSym.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) parseExpr(pre Precedence) tree.Expr {
	prefix := p.prefixParseFns[p.cSym.Type]
	if prefix == nil {
		p.ErrorWithLoc(errors.ERROR, "no prefix parse function for %s", p.cSym.Loc.String(), p.cSym.Type)
		return nil
	}
	left := prefix()

	for pre < p.precedence() {
		infix := p.infixParseFns[p.pSym.Type]
		if infix == nil {
			return left
		}
		p.next()
		left = infix(left)
	}
	return left
}

func (p *Parser) parseExprList() []tree.Expr {
	var exprs []tree.Expr
	exprs = append(exprs, p.parseExpr(LOWEST))
	for p.accept(symbol.COMMA) {
		p.next()
		exprs = append(exprs, p.parseExpr(LOWEST))
	}

	return exprs
}

func (p *Parser) parseIntegerLiteral() tree.Expr {
	return &tree.BasicLit{
		ValuePos: p.cSym.Loc,
		Kind:     "i32",
		Value:    p.cSym.Lit,
	}
}

func (p *Parser) parseFloatLiteral() tree.Expr {
	return &tree.BasicLit{
		ValuePos: p.cSym.Loc,
		Kind:     "f32",
		Value:    p.cSym.Lit,
	}
}

func (p *Parser) parseStringLiteral() tree.Expr {
	return &tree.BasicLit{
		ValuePos: p.cSym.Loc,
		Kind:     "string",
		Value:    p.cSym.Lit,
	}
}
