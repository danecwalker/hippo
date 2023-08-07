package parse

import (
	"fmt"
	"os"

	"github.com/danecwalker/hippo/internal/lexer"
	"github.com/danecwalker/hippo/internal/syntax"
)

type (
	prefixParseFn func() syntax.Expression
	infixParseFn  func(syntax.Expression) syntax.Expression
	stmtParseFn   func() syntax.Statement
)

type Prec int

const (
	_ Prec = iota
	LOWEST
	EQUALS      // ==
	LESSGREATER // > or <
	SUM         // +
	PRODUCT     // *
	PREFIX      // -X or !X
	CALL        // myFunction(X)
	INDEX       // array[index]
)

type Parser struct {
	cur_token  *syntax.Token
	peek_token *syntax.Token

	lex *lexer.Lexer

	errors []*Error

	precedences map[syntax.TokenType]Prec

	prefixParseFns map[syntax.TokenType]prefixParseFn
	infixParseFns  map[syntax.TokenType]infixParseFn

	stmtParseFns map[syntax.TokenType]stmtParseFn
}

func (p *Parser) registerPrefix(token_type syntax.TokenType, fn prefixParseFn) {
	p.prefixParseFns[token_type] = fn
}

func (p *Parser) registerInfix(token_type syntax.TokenType, fn infixParseFn) {
	p.infixParseFns[token_type] = fn
}

func (p *Parser) registerStmt(token_type syntax.TokenType, fn stmtParseFn) {
	p.stmtParseFns[token_type] = fn
}

func NewParser(filename string) *Parser {
	lex := lexer.NewLexer(filename)
	p := &Parser{
		lex: lex,
	}

	p.prefixParseFns = make(map[syntax.TokenType]prefixParseFn)
	p.registerPrefix(syntax.TokenIdent, p.parseIdentifier)
	p.registerPrefix(syntax.TokenInt, p.parseIntegerLiteral)

	p.infixParseFns = make(map[syntax.TokenType]infixParseFn)
	p.registerInfix(syntax.TokenPlus, p.parseInfixExpression)
	p.registerInfix(syntax.TokenMinus, p.parseInfixExpression)
	p.registerInfix(syntax.TokenStar, p.parseInfixExpression)
	p.registerInfix(syntax.TokenSlash, p.parseInfixExpression)
	p.registerInfix(syntax.TokenLParen, p.parseCallExpression)
	p.registerInfix(syntax.TokenRange, p.parseRangeExpression)

	p.stmtParseFns = make(map[syntax.TokenType]stmtParseFn)
	p.registerStmt(syntax.TokenVar, p.parseVarStatement)
	p.registerStmt(syntax.TokenConst, p.parseVarStatement)
	p.registerStmt(syntax.TokenFunc, p.parseFunctionStatement)
	p.registerStmt(syntax.TokenReturn, p.parseReturnStatement)
	p.registerStmt(syntax.TokenFor, p.parseForStatement)
	p.registerStmt(syntax.TokenIdent, p.parseExpressionStatement)

	p.precedences = make(map[syntax.TokenType]Prec)
	p.precedences[syntax.TokenAssign] = EQUALS
	p.precedences[syntax.TokenPlus] = SUM
	p.precedences[syntax.TokenMinus] = SUM
	p.precedences[syntax.TokenSlash] = PRODUCT
	p.precedences[syntax.TokenStar] = PRODUCT
	p.precedences[syntax.TokenLParen] = CALL
	p.precedences[syntax.TokenRange] = INDEX

	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) Errors() []*Error {
	return p.errors
}

func (p *Parser) nextToken() {
	p.cur_token = p.peek_token
	p.peek_token = p.lex.NextToken()
}

func (p *Parser) peekTokenIs(token_type syntax.TokenType) bool {
	return p.peek_token.Type == token_type
}

func (p *Parser) curTokenIs(token_type syntax.TokenType) bool {
	return p.cur_token.Type == token_type
}

func (p *Parser) expectPeek(token_type syntax.TokenType) bool {
	if p.peek_token.Type == token_type {
		p.nextToken()
		return true
	} else {
		NewPeekError(p.cur_token.Position, token_type.String(), p.peek_token.Type.String())
		return false
	}
}

func (p *Parser) acceptPeek(token_type syntax.TokenType) bool {
	if p.peek_token.Type == token_type {
		p.nextToken()
		return true
	} else {
		return false
	}
}

func (p *Parser) parseStatement() syntax.Statement {
	if fn := p.stmtParseFns[p.cur_token.Type]; fn != nil {
		return fn()
	} else {
		p.errors = append(p.errors, NewUnexpectedTokenError(p.cur_token.Position, p.cur_token))
	}

	return nil
}

func (p *Parser) peekPrecedence() Prec {
	if precedence, ok := p.precedences[p.peek_token.Type]; ok {
		return precedence
	}

	return LOWEST
}

func (p *Parser) curPrecedence() Prec {
	if precedence, ok := p.precedences[p.cur_token.Type]; ok {
		return precedence
	}

	return LOWEST
}

func (p *Parser) parseExpression(precedence Prec) syntax.Expression {
	if fn := p.prefixParseFns[p.cur_token.Type]; fn != nil {
		left_exp := fn()

		for p.peek_token.Type != syntax.TokenEOF && precedence < p.peekPrecedence() {
			if fn := p.infixParseFns[p.peek_token.Type]; fn != nil {
				p.nextToken()
				left_exp = fn(left_exp)
			} else {
				p.errors = append(p.errors, NewUnexpectedTokenError(p.cur_token.Position, p.cur_token))
				return left_exp
			}
		}

		return left_exp
	} else {
		p.errors = append(p.errors, NewUnexpectedTokenError(p.cur_token.Position, p.cur_token))
	}

	return nil
}

func (p *Parser) ParseProgram() *syntax.Program {
	program := syntax.NewProgram()

	for p.cur_token.Type != syntax.TokenEOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.AddStatement(stmt)
		}
		p.nextToken()
	}

	program.
		PrettyPrint()

	if len(p.errors) > 0 {
		for _, err := range p.errors {
			fmt.Fprintln(os.Stderr, err)
		}
		os.Exit(1)
	}

	return program
}

func (p *Parser) parseVarStatement() syntax.Statement {
	position := p.cur_token.Position
	kind := p.cur_token.Literal

	if !p.expectPeek(syntax.TokenIdent) {
		return nil
	}

	names := p.parseIdentifierList()

	var type_ *syntax.Identifier
	if p.acceptPeek(syntax.TokenIdent) {
		type_ = p.parseIdentifier().(*syntax.Identifier)
	}

	var values []syntax.Expression
	if p.acceptPeek(syntax.TokenAssign) {
		p.nextToken()
		values = p.parseExpressionList()
	}

	return syntax.NewNVarStatement(position, kind, names, type_, values)
}

func (p *Parser) parseFunctionStatement() syntax.Statement {
	position := p.cur_token.Position

	if !p.expectPeek(syntax.TokenIdent) {
		return nil
	}

	name, ok := p.parseIdentifier().(*syntax.Identifier)
	if !ok {
		fmt.Fprintln(os.Stderr, NewError(p.cur_token.Position, "function requires a name"))
		os.Exit(1)
	}

	var params_ []*syntax.NField
	if p.acceptPeek(syntax.TokenColon) {
		p.nextToken()
		params_ = p.parseNFields()
	}

	var return_type []*syntax.Identifier
	if p.acceptPeek(syntax.TokenArrow) {
		p.nextToken()
		return_type = p.parseIdentifierList()
	}

	if !p.expectPeek(syntax.TokenLBrace) {
		return nil
	}

	body := p.parseBlockStatement()

	return syntax.NewFuncStmt(position, name, syntax.NewFuncType(params_, return_type), body)
}

func (p *Parser) parseNFields() []*syntax.NField {
	fields := make([]*syntax.NField, 0)

	names := p.parseIdentifierList()

	if !p.expectPeek(syntax.TokenIdent) {
		return nil
	}

	type_ := p.parseIdentifier().(*syntax.Identifier)

	fields = append(fields, syntax.NewNField(names, type_))

	for p.acceptPeek(syntax.TokenComma) {
		p.nextToken()

		names := p.parseIdentifierList()

		if !p.expectPeek(syntax.TokenIdent) {
			return nil
		}

		type_ := p.parseIdentifier().(*syntax.Identifier)

		fields = append(fields, syntax.NewNField(names, type_))
	}

	return fields
}

func (p *Parser) parseExpressionList() []syntax.Expression {
	exprs := make([]syntax.Expression, 0)

	exprs = append(exprs, p.parseExpression(LOWEST))

	for p.acceptPeek(syntax.TokenComma) {
		p.nextToken()
		exprs = append(exprs, p.parseExpression(LOWEST))
	}

	return exprs
}

func (p *Parser) parseIdentifierList() []*syntax.Identifier {
	identifiers := make([]*syntax.Identifier, 0)

	identifiers = append(identifiers, p.parseIdentifier().(*syntax.Identifier))

	for p.acceptPeek(syntax.TokenComma) {
		p.nextToken()
		identifiers = append(identifiers, p.parseIdentifier().(*syntax.Identifier))
	}

	return identifiers
}

func (p *Parser) parseFields() []*syntax.Field {
	var fields []*syntax.Field = make([]*syntax.Field, 0)

	name := p.parseIdentifier().(*syntax.Identifier)

	if !p.expectPeek(syntax.TokenIdent) {
		return nil
	}

	type_ := p.parseIdentifier().(*syntax.Identifier)

	fields = append(fields, syntax.NewField(name, type_))
	return fields
}

func (p *Parser) parseBlockStatement() *syntax.BlockStmt {
	position := p.cur_token.Position

	block := syntax.NewBlockStmt(position)

	p.nextToken()

	for !p.curTokenIs(syntax.TokenRBrace) && !p.curTokenIs(syntax.TokenEOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			block.AddStmt(stmt)
		}
		p.nextToken()
	}

	block.Rbrace = p.cur_token.Position

	return block
}

func (p *Parser) parseReturnStatement() syntax.Statement {
	position := p.cur_token.Position

	p.nextToken()
	results := p.parseExpression(LOWEST)

	return syntax.NewReturnStmt(position, results)
}

func (p *Parser) parseForStatement() syntax.Statement {
	position := p.cur_token.Position
	if p.peekTokenIs(syntax.TokenIdent) {
		return p.parseForRangeStatement(position)
	} else if p.peekTokenIs(syntax.TokenLBrace) {
		return p.parseWhileStatement(position, syntax.NewIdentifier(position, "true"))
	} else {
		return p.parseForLoopStatement(position)
	}
}

func (p *Parser) parseForRangeStatement(position *syntax.Position) syntax.Statement {
	p.nextToken()
	expr_ := p.parseExpression(LOWEST)
	name, ok := expr_.(*syntax.Identifier)
	if !ok {
		return p.parseWhileStatement(position, expr_)
	}

	if !p.expectPeek(syntax.TokenInfer) {
		return nil
	}

	p.nextToken()

	expr := p.parseExpression(LOWEST)

	if !p.expectPeek(syntax.TokenLBrace) {
		return nil
	}

	body := p.parseBlockStatement()

	if !p.expectPeek(syntax.TokenRBrace) {
		return nil
	}

	return syntax.NewForRangeStmt(position, name, expr, body)
}

func (p *Parser) parseWhileStatement(position *syntax.Position, expr syntax.Expression) syntax.Statement {
	if !p.expectPeek(syntax.TokenLBrace) {
		return nil
	}

	body := p.parseBlockStatement()

	if !p.expectPeek(syntax.TokenRBrace) {
		return nil
	}

	return syntax.NewWhileStmt(expr, body)
}

func (p *Parser) parseForLoopStatement(position *syntax.Position) syntax.Statement {
	p.nextToken()

	init := p.parseStatement()

	if !p.expectPeek(syntax.TokenSemicolon) {
		return nil
	}

	p.nextToken()
	cond := p.parseExpression(LOWEST)

	if !p.expectPeek(syntax.TokenSemicolon) {
		return nil
	}

	p.nextToken()
	inc := p.parseStatement()

	if !p.expectPeek(syntax.TokenLBrace) {
		return nil
	}

	body := p.parseBlockStatement()

	if !p.expectPeek(syntax.TokenRBrace) {
		return nil
	}

	return syntax.NewForLoopStmt(init, cond, inc, body)
}

func (p *Parser) parseIdentifier() syntax.Expression {
	return syntax.NewIdentifier(p.cur_token.Position, p.cur_token.Literal)
}

func (p *Parser) parseIntegerLiteral() syntax.Expression {
	position := p.cur_token.Position
	literal := p.cur_token.Literal
	return syntax.NewBasicLit(position, "i32", literal)
}

func (p *Parser) parseInfixExpression(expr syntax.Expression) syntax.Expression {
	operator := p.cur_token
	precedence := p.curPrecedence()
	p.nextToken()
	right := p.parseExpression(precedence)
	return syntax.NewBinaryExpr(expr, operator, right)
}

func (p *Parser) parseCallExpression(expr syntax.Expression) syntax.Expression {
	position := p.cur_token.Position

	name, ok := expr.(*syntax.Identifier)
	if !ok {
		fmt.Fprintln(os.Stderr, NewError(p.cur_token.Position, "call requires a name"))
		os.Exit(1)
	}

	p.nextToken()
	args := p.parseExpression(LOWEST)

	if !p.expectPeek(syntax.TokenRParen) {
		return nil
	}

	return syntax.NewCallExpr(name, position, args, p.cur_token.Position)
}

func (p *Parser) parseExpressionStatement() syntax.Statement {
	expr := p.parseExpression(LOWEST)
	return syntax.NewExpressionStmt(expr)
}

func (p *Parser) parseRangeExpression(expr syntax.Expression) syntax.Expression {
	p.nextToken()
	right := p.parseExpression(LOWEST)
	return syntax.NewRangeExpr(expr, right)
}
