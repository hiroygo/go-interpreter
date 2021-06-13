package lexer

import "github.com/hiroygo/go-interpreter/token"

type Lexer struct {
	input        string
	position     int  // 入力における現在の位置(現在の文字を指し示す)
	readPosition int  // これから読み込む位置(現在の文字の次)
	ch           byte // 現在検査中の文字
}

func New(s string) *Lexer {
	l := &Lexer{input: s}
	// NextToken の実行前に呼び出す必要がある
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++
}

func (l *Lexer) eatWhiteSpace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func (l *Lexer) readIdentifier() string {
	head := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[head:l.position]
}

func (l *Lexer) readNumber() string {
	head := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[head:l.position]
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

func (l *Lexer) NextToken() token.Token {
	l.eatWhiteSpace()

	c := l.ch
	t := token.Token{}
	switch c {
	case '=':
		p := l.peekChar()
		if p == '=' {
			t = token.Token{Type: token.EQ, Literal: string(c) + string(p)}
			l.readChar()
		} else {
			t = newToken(token.ASSIGN, c)
		}
	case '+':
		t = newToken(token.PLUS, c)
	case '-':
		t = newToken(token.MINUS, c)
	case '!':
		p := l.peekChar()
		if p == '=' {
			t = token.Token{Type: token.NOT_EQ, Literal: string(c) + string(p)}
			l.readChar()
		} else {
			t = newToken(token.BANG, c)
		}
	case '/':
		t = newToken(token.SLASH, c)
	case '*':
		t = newToken(token.ASTERISK, c)
	case '<':
		t = newToken(token.LT, c)
	case '>':
		t = newToken(token.GT, c)
	case ';':
		t = newToken(token.SEMICOLON, c)
	case ',':
		t = newToken(token.COMMA, c)
	case '{':
		t = newToken(token.LBRACE, c)
	case '}':
		t = newToken(token.RBRACE, c)
	case '(':
		t = newToken(token.LPAREN, c)
	case ')':
		t = newToken(token.RPAREN, c)
	case 0:
		t = token.Token{Type: token.EOF, Literal: ""}
	default:
        // 言語のキーワードか変数名か判定する
		if isLetter(c) {
			ident := l.readIdentifier()
			tt := token.LookupIdent(ident)
			// ここで return するのは readIdentifier() で
			// 次の読み取るべき位置に移動済だから
			return token.Token{Type: tt, Literal: ident}
		}
		if isDigit(c) {
			strNum := l.readNumber()
			return token.Token{Type: token.INT, Literal: strNum}
		}
		t = newToken(token.ILLEGAL, c)
	}

	l.readChar()
	return t
}

func newToken(t token.TokenType, c byte) token.Token {
	return token.Token{Type: t, Literal: string(c)}
}

func isLetter(c byte) bool {
	return 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z' || c == '_'
}

func isDigit(c byte) bool {
	return '0' <= c && c <= '9'
}
