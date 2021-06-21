package parser

import (
	"fmt"

	"github.com/hiroygo/go-interpreter/ast"
	"github.com/hiroygo/go-interpreter/lexer"
	"github.com/hiroygo/go-interpreter/token"
)

type Parser struct {
	l         *lexer.Lexer
	curToken  token.Token
	peekToken token.Token
	errors    []string
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l}

	// curToken と peekToken を初期位置にセットする
	p.nextToken()
	p.nextToken()

	return p
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
		return nil
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
