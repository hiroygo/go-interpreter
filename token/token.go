package token

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// 識別子(変数名)を表す
	IDENT = "IDENT"

	// 記号を表す
	ASSIGN    = "="
	PLUS      = "+"
	COMMA     = ","
	SEMICOLON = ";"
	LPAREN    = "("
	RPAREN    = ")"
	LBRACE    = "{"
	RBRACE    = "}"
	MINUS     = "-"
	BANG      = "!"
	ASTERISK  = "*"
	SLASH     = "/"
	LT        = "<"
	GT        = ">"
	EQ        = "=="
	NOT_EQ    = "!="

	// 予約語を表す
	// "return", "let" など
	INT      = "INT"
	FUNCTION = "FUNCTION"
	LET      = "LET"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	IF       = "IF"
	ELSE     = "ELSE"
	RETURN   = "RETURN"
)

// デバッグしやすいように string にしておく
type TokenType string

// 変数トークンのときは Token{Type: IDENT, Literal: "foo"} のようになる
type Token struct {
	Type    TokenType
	Literal string
}

var keywords = map[string]TokenType{
	"fn":     FUNCTION,
	"let":    LET,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
}

func LookupIdent(s string) TokenType {
	if v, ok := keywords[s]; ok {
		return v
	}
	return IDENT
}
