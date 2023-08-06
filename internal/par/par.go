package par

import (
	"github.com/danecwalker/hippo/internal/ast"
	"github.com/danecwalker/hippo/internal/logger"
	"github.com/danecwalker/hippo/internal/lxr"
	"github.com/danecwalker/hippo/internal/tok"
)

type Parser struct {
	lexer *lxr.Lexer

	// Current token
	cur_tok tok.Token

	// Next token
	peek_tok tok.Token

	// Error handler
	log *logger.LogHandler

	// Expression parsing functions
	prefixParseFns map[tok.TokenType]prefixParseFn
	infixParseFns  map[tok.TokenType]infixParseFn
}

type (
	prefixParseFn func() ast.Expr
	infixParseFn  func(ast.Expr) ast.Expr
)

func (p *Parser) registerPrefix(token_type tok.TokenType, fn prefixParseFn) {
	p.prefixParseFns[token_type] = fn
}

func (p *Parser) registerInfix(token_type tok.TokenType, fn infixParseFn) {
	p.infixParseFns[token_type] = fn
}

func NewParser(file_name string) *Parser {
	log := logger.NewLogHandler()
	return &Parser{
		lexer:          lxr.NewLexer(file_name),
		log:            log,
		prefixParseFns: make(map[tok.TokenType]prefixParseFn),
		infixParseFns:  make(map[tok.TokenType]infixParseFn),
	}
}

func (p *Parser) Parse() *ast.File {
	p.log.InfoLog("Parsing file %s", p.lexer.Filename)

	// Register prefix parsing functions
	p.registerPrefix(tok.IDENT, p.parseIdent)
	p.registerPrefix(tok.INT, p.parseIntegerLiteral)
	p.registerPrefix(tok.STRING, p.parseStringLiteral)

	// Register infix parsing functions
	p.registerInfix(tok.PLUS, p.parseBinaryExpr)
	p.registerInfix(tok.MINUS, p.parseBinaryExpr)
	p.registerInfix(tok.ASTERISK, p.parseBinaryExpr)
	p.registerInfix(tok.SLASH, p.parseBinaryExpr)
	p.registerInfix(tok.LPAREN, p.parseCallExpr)
	p.registerInfix(tok.RANGE, p.parseRangeExpr)

	p.nextToken()
	p.nextToken()

	f := p.parseFile()
	f.Filename = p.lexer.Filename

	p.log.MaybeFatal()
	return f
}

func (p *Parser) nextToken() {
	p.cur_tok = p.peek_tok
	p.peek_tok = p.lexer.NextToken()
}

func (p *Parser) expectPeek(t tok.TokenType) bool {
	if p.peek_tok.Type == t {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

func (p *Parser) acceptPeek(t tok.TokenType) bool {
	if p.peek_tok.Type == t {
		p.nextToken()
		return true
	} else {
		return false
	}
}

func (p *Parser) curTokenIs(t tok.TokenType) bool {
	return p.cur_tok.Type == t
}

func (p *Parser) peekTokenIs(t tok.TokenType) bool {
	return p.peek_tok.Type == t
}

func (p *Parser) peekError(t tok.TokenType) {
	p.log.LogAt(logger.LogError, &p.cur_tok.Bounds.Start, "Expected next token to be %s, got %s instead", t, p.peek_tok.Type)
}

func (p *Parser) parseFile() *ast.File {
	file := &ast.File{}

	for !p.curTokenIs(tok.EOF) {
		decl := p.parseDecl()
		if decl != nil {
			file.Decls = append(file.Decls, decl)
		}
	}

	return file
}

func (p *Parser) parseDecl() ast.Decl {
	switch p.cur_tok.Type {
	case
		tok.CONST,
		tok.VAR:
		return p.parseInitDecl()
	case tok.FUNCTION:
		return p.parseFuncDecl()
	default:
		p.log.LogAt(logger.LogFatal, &p.cur_tok.Bounds.Start, "Expected declaration, got %s instead", p.cur_tok.Type)
		return nil
	}
}

func (p *Parser) parseInitDecl() *ast.InitDecl {
	decl := &ast.InitDecl{}
	if !p.expectPeek(tok.IDENT) {
		return nil
	}

	decl.Names = append(decl.Names, &ast.Ident{Name: p.cur_tok.Literal, NamePos: p.cur_tok.Bounds.Start})

	for p.acceptPeek(tok.COMMA) {
		if !p.expectPeek(tok.IDENT) {
			return nil
		}
		decl.Names = append(decl.Names, &ast.Ident{Name: p.cur_tok.Literal, NamePos: p.cur_tok.Bounds.Start})
	}

	if !p.expectPeek(tok.IDENT) {
		return nil
	}

	decl.Type = &ast.Ident{Name: p.cur_tok.Literal, NamePos: p.cur_tok.Bounds.Start}

	if p.acceptPeek(tok.ASSIGN) {
		p.nextToken()

		decl.Values = p.parseExprList()
	}

	p.nextToken()
	return decl
}

func (p *Parser) parseFuncDecl() *ast.FuncDecl {
	decl := &ast.FuncDecl{FuncPos: p.cur_tok.Bounds.Start}

	if !p.expectPeek(tok.IDENT) {
		return nil
	}

	decl.Name = &ast.Ident{Name: p.cur_tok.Literal, NamePos: p.cur_tok.Bounds.Start}

	decl.Type = &ast.FuncType{
		Params:  []*ast.Field{},
		Results: []*ast.Field{},
	}

	if p.acceptPeek(tok.COLON) {
		for p.expectPeek(tok.IDENT) {
			field := &ast.Field{}
			field.Names = append(field.Names, &ast.Ident{Name: p.cur_tok.Literal, NamePos: p.cur_tok.Bounds.Start})

			for p.acceptPeek(tok.COMMA) {
				if !p.expectPeek(tok.IDENT) {
					return nil
				}
				field.Names = append(field.Names, &ast.Ident{Name: p.cur_tok.Literal, NamePos: p.cur_tok.Bounds.Start})
			}

			if !p.expectPeek(tok.IDENT) {
				return nil
			}

			field.Type = &ast.Ident{Name: p.cur_tok.Literal, NamePos: p.cur_tok.Bounds.Start}

			decl.Type.Params = append(decl.Type.Params, field)

			if !p.acceptPeek(tok.COMMA) {
				break
			}
		}
	}

	if p.acceptPeek(tok.RARROW) {
		for p.expectPeek(tok.IDENT) {
			field := &ast.Field{}
			field.Type = &ast.Ident{Name: p.cur_tok.Literal, NamePos: p.cur_tok.Bounds.Start}

			decl.Type.Results = append(decl.Type.Results, field)

			if !p.acceptPeek(tok.COMMA) {
				break
			}
		}
	}

	if !p.expectPeek(tok.LBRACE) {
		return nil
	}

	decl.Body = p.parseBlock()

	if !p.expectPeek(tok.RBRACE) {
		return nil
	}

	p.nextToken()

	return decl
}

func (p *Parser) parseBlock() *ast.BlockStmt {
	block := &ast.BlockStmt{}

	for !p.peekTokenIs(tok.RBRACE) && !p.peekTokenIs(tok.EOF) {
		p.nextToken()
		stmt := p.parseStmt()

		if stmt != nil {
			block.Stmts = append(block.Stmts, stmt)
		}
	}

	return block
}

func (p *Parser) parseStmt() ast.Stmt {
	switch p.cur_tok.Type {
	case tok.FOR:
		return p.parseRangeStmt()
	case tok.RETURN:
		return p.parseReturnStmt()
	case tok.VAR, tok.CONST:
		return p.parseDeclStmt()
	case tok.IDENT:
		if p.peekTokenIs(tok.LPAREN) {
			return p.parseCallStmt()
		} else {
			p.log.LogAt(logger.LogError, &p.cur_tok.Bounds.Start, "Expected assignment or call statement, got %s instead", p.peek_tok.Type)
			return nil
		}
	default:
		return nil
	}
}

func (p *Parser) parseRangeStmt() ast.Stmt {
	stmt := &ast.RangeStmt{}

	if !p.expectPeek(tok.IDENT) {
		return nil
	}

	stmt.Names = append(stmt.Names, &ast.Ident{Name: p.cur_tok.Literal, NamePos: p.cur_tok.Bounds.Start})

	for p.acceptPeek(tok.COMMA) {
		if !p.expectPeek(tok.IDENT) {
			return nil
		}

		stmt.Names = append(stmt.Names, &ast.Ident{Name: p.cur_tok.Literal, NamePos: p.cur_tok.Bounds.Start})
	}

	if !p.expectPeek(tok.LARROW) {
		return nil
	}

	p.nextToken()

	stmt.Iter = p.parseExpr(LOWEST)

	if !p.expectPeek(tok.LBRACE) {
		return nil
	}

	stmt.Body = p.parseBlock()

	if !p.expectPeek(tok.RBRACE) {
		return nil
	}

	return stmt
}

func (p *Parser) parseDeclStmt() ast.Stmt {
	decl := &ast.DeclStmt{}

	if !p.expectPeek(tok.IDENT) {
		return nil
	}
	decl.Names = append(decl.Names, &ast.Ident{Name: p.cur_tok.Literal, NamePos: p.cur_tok.Bounds.Start})

	for p.acceptPeek(tok.COMMA) {
		if !p.expectPeek(tok.IDENT) {
			return nil
		}
		decl.Names = append(decl.Names, &ast.Ident{Name: p.cur_tok.Literal, NamePos: p.cur_tok.Bounds.Start})
	}

	if p.acceptPeek(tok.IDENT) {
		decl.Type = &ast.Ident{Name: p.cur_tok.Literal, NamePos: p.cur_tok.Bounds.Start}
	}

	if p.acceptPeek(tok.ASSIGN) {
		p.nextToken()

		decl.Values = p.parseExprList()
	}
	return decl
}

func (p *Parser) parseCallStmt() *ast.ExprStmt {
	stmt := &ast.ExprStmt{}

	ident := &ast.Ident{Name: p.cur_tok.Literal, NamePos: p.cur_tok.Bounds.Start}

	if !p.expectPeek(tok.LPAREN) {
		return nil
	}

	stmt.X = p.parseCallExpr(ident)

	return stmt
}

func (p *Parser) parseReturnStmt() *ast.ReturnStmt {
	stmt := &ast.ReturnStmt{ReturnPos: p.cur_tok.Bounds.Start, Results: make([]ast.Expr, 0)}

	p.nextToken()

	stmt.Results = p.parseExprList()

	return stmt
}

func (p *Parser) parseExprList() []ast.Expr {
	list := []ast.Expr{}

	list = append(list, p.parseExpr(LOWEST))

	for p.acceptPeek(tok.COMMA) {
		p.nextToken()
		list = append(list, p.parseExpr(LOWEST))
	}

	return list
}

type P int

const (
	_ P = iota
	LOWEST
	ASSIGN
	OR
	AND
	EQUALS
	LESSGREATER
	SUM
	PRODUCT
	PREFIX
	CALL
)

var precedences = map[tok.TokenType]P{
	tok.EQ:       EQUALS,
	tok.NOT_EQ:   EQUALS,
	tok.LT:       LESSGREATER,
	tok.GT:       LESSGREATER,
	tok.PLUS:     SUM,
	tok.MINUS:    SUM,
	tok.SLASH:    PRODUCT,
	tok.ASTERISK: PRODUCT,
	tok.LPAREN:   CALL,
	tok.RANGE:    CALL,
}

func (p *Parser) peekPrecedence() P {
	if p, ok := precedences[p.peek_tok.Type]; ok {
		return p
	}

	return LOWEST
}

func (p *Parser) curPrecedence() P {
	if p, ok := precedences[p.cur_tok.Type]; ok {
		return p
	}

	return LOWEST
}

func (p *Parser) parseExpr(precedence P) ast.Expr {
	prefix := p.prefixParseFns[p.cur_tok.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.cur_tok.Type)
		return nil
	}

	leftExp := prefix()

	for !p.peekTokenIs(tok.EOF) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peek_tok.Type]
		if infix == nil {
			return leftExp
		}

		p.nextToken()

		leftExp = infix(leftExp)
	}

	return leftExp
}

func (p *Parser) noPrefixParseFnError(t tok.TokenType) {
	p.log.LogAt(logger.LogError, &p.cur_tok.Bounds.Start, "No prefix parse function for %s found", t)
}

func (p *Parser) parseIdent() ast.Expr {
	return &ast.Ident{Name: p.cur_tok.Literal, NamePos: p.cur_tok.Bounds.Start}
}

func (p *Parser) parseIntegerLiteral() ast.Expr {
	lit := &ast.BasicLit{ValuePos: p.cur_tok.Bounds.Start, Kind: "i32"}
	lit.Value = p.cur_tok.Literal

	return lit
}

func (p *Parser) parseStringLiteral() ast.Expr {
	lit := &ast.BasicLit{ValuePos: p.cur_tok.Bounds.Start, Kind: "string"}
	lit.Value = p.cur_tok.Literal

	return lit
}

func (p *Parser) parseBoolean() ast.Expr {
	lit := &ast.BasicLit{ValuePos: p.cur_tok.Bounds.Start, Kind: "bool"}
	lit.Value = p.cur_tok.Literal

	return lit
}

func (p *Parser) parseBinaryExpr(left ast.Expr) ast.Expr {
	expr := &ast.BinaryExpr{X: left, Op: p.cur_tok.Literal}

	precedence := p.curPrecedence()
	p.nextToken()
	expr.Y = p.parseExpr(precedence)

	return expr
}

func (p *Parser) parseCallExpr(left ast.Expr) ast.Expr {
	call := &ast.CallExpr{Fun: left, Lparen: p.cur_tok.Bounds.Start, Args: make([]ast.Expr, 0)}

	p.nextToken()

	call.Args = p.parseExprList()

	if !p.expectPeek(tok.RPAREN) {
		return nil
	}

	call.Rparen = p.cur_tok.Bounds.Start

	return call
}

func (p *Parser) parseRangeExpr(left ast.Expr) ast.Expr {
	expr := &ast.RangeExpr{Low: left}

	precedence := p.curPrecedence()
	p.nextToken()
	expr.High = p.parseExpr(precedence)

	return expr
}
