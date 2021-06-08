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

func (l *Lexer) NextToken() token.Token {
	c := l.ch
    literal := string(c)
	l.readChar()

	switch c {
	case '=':
		return token.Token{Type: token.ASSIGN, Literal: literal}
	case ';':
		return token.Token{Type: token.SEMICOLON, Literal: literal}
	case '(':
		return token.Token{Type: token.LPAREN, Literal: literal}
	case ')':
		return token.Token{Type: token.RPAREN, Literal: literal}
	case ',':
		return token.Token{Type: token.COMMA, Literal: literal}
	case '+':
		return token.Token{Type: token.PLUS, Literal: literal}
	case '{':
		return token.Token{Type: token.LBRACE, Literal: literal}
	case '}':
		return token.Token{Type: token.RBRACE, Literal: literal}
	case 0:
		return token.Token{Type: token.EOF, Literal: ""}
	default:
		return token.Token{Type: token.ILLEGAL, Literal: literal}
	}
}
