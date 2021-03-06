package parser

import (
	"fmt"
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

func TestParsingPrefixExpressions(t *testing.T) {
	cases := []struct {
		input    string
		operator string
		v        interface{}
	}{
		{"!5;", "!", 5},
		{"-15;", "-", 15},
		{"!true;", "!", true},
		{"!false;", "!", false},
	}

	for _, c := range cases {
		c := c
		t.Run(c.input, func(t *testing.T) {
			t.Parallel()

			p := New(lexer.New(c.input))
			prg := p.ParseProgram()
			hasParserErrors(t, p)
			if len(prg.Statements) != 1 {
				t.Fatalf("want len(Program.Statements) = %v, got %v", 1, len(prg.Statements))
			}

			stmt, ok := prg.Statements[0].(*ast.ExpressionStatement)
			if !ok {
				t.Fatalf("%T.(*ast.ExpressionStatement) error", prg.Statements[0])
			}
			exp, ok := stmt.Expression.(*ast.PrefixExpression)
			if !ok {
				t.Fatalf("%T.(*ast.PrefixExpression) error", stmt)
			}
			if exp.Operator != c.operator {
				t.Fatalf("want PrefixExpression.Operator = %q, got %q", c.operator, exp.Operator)
			}
			testLiteralExpression(t, exp.Right, c.v)
		})
	}
}

func testIntegerLiteral(t *testing.T, e ast.Expression, v int64) {
	t.Helper()

	literal, ok := e.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("%T.(*ast.IntegerLiteral) error", e)
	}
	if literal.Value != v {
		t.Fatalf("want IntegerLiteral.Value = %d, got %d", v, literal.Value)
	}
	if literal.TokenLiteral() != fmt.Sprintf("%d", v) {
		t.Fatalf("want IntegerLiteral.TokenLiteral() = '%d', got %q", v, literal.TokenLiteral())
	}
}

func testIdentifier(t *testing.T, e ast.Expression, v string) {
	t.Helper()

	ident, ok := e.(*ast.Identifier)
	if !ok {
		t.Fatalf("%T.(*ast.Identifier) error", e)
	}
	if ident.Value != v {
		t.Fatalf("want Identifier.Value = %q, got %q", v, ident.Value)
	}
	if ident.TokenLiteral() != v {
		t.Fatalf("want Identifier.TokenLiteral() = %q, got %q", v, ident.TokenLiteral())
	}
}

func testBooleanLiteral(t *testing.T, e ast.Expression, v bool) {
	t.Helper()

	literal, ok := e.(*ast.Boolean)
	if !ok {
		t.Fatalf("%T.(*ast.Boolean) error", e)
	}
	if literal.Value != v {
		t.Fatalf("want Boolean.Value = %t, got %t", v, literal.Value)
	}
	if literal.TokenLiteral() != fmt.Sprintf("%t", v) {
		t.Fatalf("want IntegerLiteral.TokenLiteral() = '%t', got %q", v, literal.TokenLiteral())
	}
}

func testLiteralExpression(t *testing.T, e ast.Expression, expected interface{}) {
	t.Helper()

	switch v := expected.(type) {
	case int:
		testIntegerLiteral(t, e, int64(v))
	case int64:
		testIntegerLiteral(t, e, v)
	case string:
		testIdentifier(t, e, v)
	case bool:
		testBooleanLiteral(t, e, v)
	default:
		t.Fatalf("unexpected type, %T", v)
	}
}

func testInfixExpression(t *testing.T, e ast.Expression, left interface{},
	operator string, right interface{}) {
	t.Helper()

	opExp, ok := e.(*ast.InfixExpression)
	if !ok {
		t.Fatalf("%T.(*ast.InfixExpression) error", e)
	}
	testLiteralExpression(t, opExp.Left, left)
	if opExp.Operator != operator {
		t.Fatalf("want InfixExpression.Operator = %q, got %q", operator, opExp.Operator)
	}
	testLiteralExpression(t, opExp.Right, right)
}

func TestParsingInfixExpressions(t *testing.T) {
	cases := []struct {
		input      string
		leftValue  interface{}
		operator   string
		rightValue interface{}
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
		{"true == true", true, "==", true},
		{"true != false", true, "!=", false},
		{"false == false", false, "==", false},
	}

	for _, c := range cases {
		c := c
		t.Run(c.input, func(t *testing.T) {
			t.Parallel()

			p := New(lexer.New(c.input))
			prg := p.ParseProgram()
			hasParserErrors(t, p)
			if len(prg.Statements) != 1 {
				t.Fatalf("want len(Program.Statements) = %v, got %v", 1, len(prg.Statements))
			}

			stmt, ok := prg.Statements[0].(*ast.ExpressionStatement)
			if !ok {
				t.Fatalf("%T.(*ast.ExpressionStatement) error", prg.Statements[0])
			}
			testInfixExpression(t, stmt.Expression, c.leftValue, c.operator, c.rightValue)
		})
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{
			"-a * b",
			"((-a) * b)",
		},
		{
			"!-a",
			"(!(-a))",
		},
		{
			"a + b + c",
			"((a + b) + c)",
		},
		{
			"a + b - c",
			"((a + b) - c)",
		},
		{
			"a * b * c",
			"((a * b) * c)",
		},
		{
			"a * b / c",
			"((a * b) / c)",
		},
		{
			"a + b / c",
			"(a + (b / c))",
		},
		{
			"a + b * c + d / e - f",
			"(((a + (b * c)) + (d / e)) - f)",
		},
		{
			// len(prg.Statements) == 2
			"3 + 4; -5 * 5",
			"(3 + 4)((-5) * 5)",
		},
		{
			// len(prg.Statements) == 1
			"5 > 4 == 3 < 4",
			"((5 > 4) == (3 < 4))",
		},
		{
			"5 < 4 != 3 > 4",
			"((5 < 4) != (3 > 4))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			"true",
			"true",
		},
		{
			"false",
			"false",
		},
		{
			"3 > 5 == false",
			"((3 > 5) == false)",
		},
		{
			"3 < 5 == true",
			"((3 < 5) == true)",
		},
		{
			"1 + (2 + 3) + 4",
			"((1 + (2 + 3)) + 4)",
		},
		{
			"(5 + 5) * 2",
			"((5 + 5) * 2)",
		},
		{
			"2 / (5 + 5)",
			"(2 / (5 + 5))",
		},
		{
			"-(5 + 5)",
			"(-(5 + 5))",
		},
		{
			"!(true == true)",
			"(!(true == true))",
		},
	}

	for _, c := range cases {
		c := c
		t.Run(c.input, func(t *testing.T) {
			t.Parallel()

			p := New(lexer.New(c.input))
			prg := p.ParseProgram()
			hasParserErrors(t, p)
			actual := prg.String()
			if actual != c.expected {
				t.Fatalf("want Program.String() = %q, got %q", c.expected, actual)
			}
		})
	}
}

func TestBooleanExpression(t *testing.T) {
	cases := []struct {
		input    string
		expected bool
	}{
		{"true;", true},
		{"false;", false},
	}

	for _, c := range cases {
		c := c
		t.Run(c.input, func(t *testing.T) {
			t.Parallel()

			p := New(lexer.New(c.input))
			prg := p.ParseProgram()
			hasParserErrors(t, p)
			if len(prg.Statements) != 1 {
				t.Fatalf("want len(Program.Statements) = %v, got %v", 1, len(prg.Statements))
			}

			stmt, ok := prg.Statements[0].(*ast.ExpressionStatement)
			if !ok {
				t.Fatalf("%T.(*ast.ExpressionStatement) error", prg.Statements[0])
			}
			boo, ok := stmt.Expression.(*ast.Boolean)
			if !ok {
				t.Fatalf("%T.(*ast.Boolean) error", stmt)
			}
			if boo.Value != c.expected {
				t.Fatalf("want Boolean.Value = %t, got %t", c.expected, boo.Value)
			}
		})
	}
}
