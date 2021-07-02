package ast

import (
	"bytes"

	"github.com/hiroygo/go-interpreter/token"
)

type Node interface {
	// TokenLiteral はデバッグとテストのために使用する
	TokenLiteral() string
	String() string
}

// この言語の文は 'let 文' と 'return 文' のみ
// 残りは式になる
// 'x + 10;' などは式文という
type Statement interface {
	Node
	// xxxNode はエラー検出用
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

// Program ノードは AST のルートノードになる
// Program は文のリストから構成される
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

func (p *Program) String() string {
	var b bytes.Buffer
	for _, s := range p.Statements {
		b.WriteString(s.String())
	}
	return b.String()
}

// 'let <identifier> = <expression>;'
// e.g. 'let x = 1;'
// e.g. 'let foo = add(x, y);'
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

func (l *LetStatement) String() string {
	var b bytes.Buffer

	b.WriteString(l.TokenLiteral() + " ")
	b.WriteString(l.Name.String())
	b.WriteString(" = ")
	if l.Value != nil {
		b.WriteString(l.Value.String())
	}
	b.WriteString(";")

	return b.String()
}

type Identifier struct {
	// Token = token.IDENT
	Token token.Token
	Value string
}

// Expression インタフェースを満たす
// Identifier が式になるのは 'let add = fn(x, y) { return x + y; };' で
// add(1, 2) とするような時
func (i *Identifier) expressionNode() {}

func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}

func (i *Identifier) String() string {
	return i.Value
}

type ReturnStatement struct {
	// Token = token.RETURN
	Token       token.Token
	ReturnValue Expression
}

func (r *ReturnStatement) statementNode() {}

func (r *ReturnStatement) TokenLiteral() string {
	return r.Token.Literal
}

func (r *ReturnStatement) String() string {
	var b bytes.Buffer

	b.WriteString(r.TokenLiteral() + " ")
	if r.ReturnValue != nil {
		b.WriteString(r.ReturnValue.String())
	}
	b.WriteString(";")

	return b.String()
}

type ExpressionStatement struct {
	// 式の最初の Token
	Token      token.Token
	Expression Expression
}

func (e *ExpressionStatement) statementNode() {}

func (e *ExpressionStatement) TokenLiteral() string {
	return e.Token.Literal
}

func (e *ExpressionStatement) String() string {
	if e.Expression != nil {
		return e.Expression.String()
	}
	return ""
}

type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (i *IntegerLiteral) expressionNode() {}

func (i *IntegerLiteral) TokenLiteral() string {
	return i.Token.Literal
}

func (i *IntegerLiteral) String() string {
	return i.Token.Literal
}

type PrefixExpression struct {
	// 前置トークン
	// e.g. '!'
	Token token.Token

	// '!' または '-' が入る
	Operator string

	Right Expression
}

func (p *PrefixExpression) expressionNode() {}

func (p *PrefixExpression) TokenLiteral() string {
	return p.Token.Literal
}

func (p *PrefixExpression) String() string {
	var b bytes.Buffer
	b.WriteString("(")
	b.WriteString(p.Operator)
	b.WriteString(p.Right.String())
	b.WriteString(")")
	return b.String()
}

type InfixExpression struct {
	// e.g. '+'
	Token    token.Token
	Left     Expression
	Operator string
	Right    Expression
}

func (i *InfixExpression) expressionNode() {}

func (i *InfixExpression) TokenLiteral() string {
	return i.Token.Literal
}

func (i *InfixExpression) String() string {
	var b bytes.Buffer
	b.WriteString("(")
	b.WriteString(i.Left.String())
	b.WriteString(" " + i.Operator + " ")
	b.WriteString(i.Right.String())
	b.WriteString(")")
	return b.String()
}
