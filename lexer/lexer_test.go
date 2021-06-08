package lexer

import (
	"testing"

	"github.com/hiroygo/go-interpreter/token"
)

func TestNextToken(t *testing.T) {
	input := `=+(){},;`
	expected := []token.Token{
		{token.ASSIGN, "="},
		{token.PLUS, "+"},
		{token.LPAREN, "("},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RBRACE, "}"},
		{token.COMMA, ","},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}

	lex := New(input)
	for _, tok := range expected {
		actual := lex.NextToken()
		if tok != actual {
			t.Fatalf("want NextToken() = %+v, got %+v", tok, actual)
		}
	}
}
