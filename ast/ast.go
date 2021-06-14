package ast

import "github.com/hiroygo/go-interpreter/token"

type Node interface {
	// TokenLiteral はデバッグとテストのために使用する
	TokenLiteral() string
}

type Statement interface {
	Node
	// xxxNode メソッドはエラー検出用
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

// Program ノードは AST のルートノードになる
type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

// `let <identifier> = <expression>;`
// e.g. `let x = 1;`
// e.g. `let foo = add(x, y);`
type LetStatement struct {
	// Token = token.LET
	Token token.Token
	Name  *Identifier
	Value Expression
}

// Statement インタフェースを満たす
func (l *LetStatement) statementNode() {}

func (l *LetStatement) TokenLiteral() string {
	return l.Token.Literal
}

type Identifier struct {
	// Token = token.IDENT
	Token token.Token
	Value string
}

// Expression インタフェースを満たす
// Identifier が式になるのは `let add = fn(x, y) { return x + y; };` で
// add(1, 2) とするような時
func (i *Identifier) expressionNode() {}

func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}
