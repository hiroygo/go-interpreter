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

	// 予約語を表す
	FUNCTION = "FUNCTION"
	LET      = "LET"
	INT      = "INT"
)

// デバッグしやすいように string にしておく
type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

var keywords = map[string]TokenType{
	"fn":  FUNCTION,
	"let": LET,
}

func LookupIdent(s string) TokenType {
	if v, ok := keywords[s]; ok {
		return v
	}
	return IDENT
}
