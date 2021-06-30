package parser

import (
	"fmt"
	"strconv"

	"github.com/hiroygo/go-interpreter/ast"
	"github.com/hiroygo/go-interpreter/lexer"
	"github.com/hiroygo/go-interpreter/token"
)

const (
	_ = iota
	LOWEST
	EQUALS      // ==
	LESSGREATER // > or <
	SUM         // +
	PRODUCT     // *
	PREFIX      // -x or !x
	CALL        // myFunction(x)
)

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

type Parser struct {
	l      *lexer.Lexer
	errors []string

	curToken  token.Token
	peekToken token.Token

	// TokenType と構文解析関数を対応させる
	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l}

	// curToken と peekToken を初期位置にセットする
	p.nextToken()
	p.nextToken()

	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)

	return p
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{
		Token: p.curToken, Value: p.curToken.Literal,
	}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	literal := &ast.IntegerLiteral{Token: p.curToken}

	v, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		s := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, s)
		return nil
	}
	literal.Value = v

	return literal
}

func (p *Parser) registerPrefix(t token.TokenType, f prefixParseFn) {
	p.prefixParseFns[t] = f
}

func (p *Parser) registerInfix(t token.TokenType, f infixParseFn) {
	p.infixParseFns[t] = f
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	// lexer を前進させる
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
	prg := &ast.Program{}
	for p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			prg.Statements = append(prg.Statements, stmt)
		}
		p.nextToken()
	}
	return prg
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	// e.g. 'let x = 10;'
	let := &ast.LetStatement{Token: p.curToken}

	// 次の Token が INDENT のときは
	// curToken に INDENT を読み込ませる
	// 'x' を読み込み
	if !p.expectPeek(token.IDENT) {
		return nil
	}
	let.Name = &ast.Identifier{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}

	// '='
	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	// TODO:
	// セミコロンに遭遇するまで式を読み飛ばしている
	// セミコロンが出現しないときは?
	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return let
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) peekError(t token.TokenType) {
	s := fmt.Sprintf("expected next token to be %q, got %q instead", t, p.peekToken.Type)
	p.errors = append(p.errors, s)
}

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}
	p.peekError(t)
	return false
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	// e.g. 'return 10;'
	r := &ast.ReturnStatement{Token: p.curToken}

	p.nextToken()

	// TODO:
	// セミコロンに遭遇するまで式を読み飛ばしている
	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return r
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	// e.g. 'foobar;'
	es := &ast.ExpressionStatement{Token: p.curToken}
	es.Expression = p.parseExpression(LOWEST)
	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return es
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		return nil
	}
	leftExp := prefix()
	return leftExp
}
