package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/hiroygo/go-interpreter/lexer"
	"github.com/hiroygo/go-interpreter/token"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	sc := bufio.NewScanner(in)
	for {
		fmt.Print(PROMPT)
		if !sc.Scan() {
			return
		}

		line := sc.Text()
		l := lexer.New(line)
		for t := l.NextToken(); t.Type != token.EOF; t = l.NextToken() {
			fmt.Printf("%+v\n", t)
		}
	}
}
