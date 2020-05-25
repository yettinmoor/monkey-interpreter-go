package repl

import (
	"bufio"
	"fmt"
	"monkey/evaluator"
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
	"monkey/token"
	"os"
)

func Repl() {
	scanner := bufio.NewScanner(os.Stdin)
	env := object.NewEnv(nil)
	for {
		fmt.Print(">> ")
		scanned := scanner.Scan()
		if !scanned {
			return
		}
		ch := make(chan *token.Token)
		l := lexer.New(scanner.Text(), ch)
		p := parser.New(l, ch)
		prog := p.Parse()
		if p.Errors != nil {
			for _, e := range p.Errors {
				fmt.Println(e.String())
			}
		} else {
			// fmt.Printf("%s\n", prog.String())
			res := evaluator.Eval(prog, env)
			fmt.Println(res)
		}
	}
}
