package parser

import (
	"testing"

	"github.com/hiroygo/go-interpreter/ast"
	"github.com/hiroygo/go-interpreter/lexer"
)

func TestLetStatements(t *testing.T) {
	// 入力は Token ではなく、文字列として与える
	// 文字列の方がテストが読みやすく、理解しやすいため
	input := `
let x = 5;
let y = 10;
let foobar = 838383;
`
	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}

	p := New(lexer.New(input))
	prg := p.ParseProgram()
	hasParserErrors(t, p)
	if len(prg.Statements) != len(tests) {
		t.Fatalf("want len(Program.Statements) = %v, got %v", len(tests), len(prg.Statements))
	}

	for i, tt := range tests {
		testLetStatement(t, prg.Statements[i], tt.expectedIdentifier)
	}
}

func testLetStatement(t *testing.T, s ast.Statement, name string) {
	t.Helper()

	if s.TokenLiteral() != "let" {
		t.Fatalf("want Statement.TokenLiteral() = 'let', got %q", s.TokenLiteral())
	}
	let, ok := s.(*ast.LetStatement)
	if !ok {
		t.Fatalf("%T.(*ast.LetStatement) error", s)
	}
	if let.Name.Value != name {
		t.Fatalf("want LetStatement.Name.Value = %q, got %q", name, let.Name.Value)
	}
	if let.Name.TokenLiteral() != name {
		t.Fatalf("want LetStatement.Name.TokenLiteral() = %q, got %q", name, let.Name.Value)
	}
}

func TestReturnStatements(t *testing.T) {
	input := `
return 5;
return 10;
return 993322;
`
	p := New(lexer.New(input))
	prg := p.ParseProgram()
	hasParserErrors(t, p)
	if len(prg.Statements) != 3 {
		t.Fatalf("want len(Program.Statements) = %v, got %v", 3, len(prg.Statements))
	}

	for _, s := range prg.Statements {
		returnStmt, ok := s.(*ast.ReturnStatement)
		if !ok {
			t.Fatalf("%T.(*ast.ReturnStatement) error", s)
		}
		if returnStmt.TokenLiteral() != "return" {
			t.Fatalf("want ReturnStatement.TokenLiteral() = %q, got %q", "return", returnStmt.TokenLiteral())
		}
	}
}

func TestIdentifierExpression(t *testing.T) {
	input := "foobar;"

	p := New(lexer.New(input))
	prg := p.ParseProgram()
	hasParserErrors(t, p)
	if len(prg.Statements) != 1 {
		t.Fatalf("want len(Program.Statements) = %v, got %v", 1, len(prg.Statements))
	}

	stmt, ok := prg.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("%T.(*ast.ExpressionStatement) error", prg.Statements[0])
	}
	ident, ok := stmt.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("%T.(*ast.Identifier) error", stmt)
	}
	if ident.Value != "foobar" {
		t.Fatalf("want Identifier.Value = %q, got %q", "foobar", ident.Value)
	}
	if ident.TokenLiteral() != "foobar" {
		t.Fatalf("want Identifier.TokenLiteral() = %q, got %q", "foobar", ident.TokenLiteral())
	}
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := "5;"

	p := New(lexer.New(input))
	prg := p.ParseProgram()
	hasParserErrors(t, p)
	if len(prg.Statements) != 1 {
		t.Fatalf("want len(Program.Statements) = %v, got %v", 1, len(prg.Statements))
	}

	stmt, ok := prg.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("%T.(*ast.ExpressionStatement) error", prg.Statements[0])
	}
	literal, ok := stmt.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("%T.(*ast.IntegerLiteral) error", stmt)
	}
	if literal.Value != 5 {
		t.Fatalf("want IntegerLiteral.Value = %d, got %d", 5, literal.Value)
	}
	if literal.TokenLiteral() != "5" {
		t.Fatalf("want IntegerLiteral.TokenLiteral() = %q, got %q", "5", literal.TokenLiteral())
	}
}

func hasParserErrors(t *testing.T, p *Parser) {
	t.Helper()
	errs := p.Errors()
	if len(errs) == 0 {
		return
	}

	t.Errorf("Parser has %d errors", len(errs))
	for i, s := range errs {
		t.Errorf("Parser error %d: %s", i, s)
	}
	t.FailNow()
}
