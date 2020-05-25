package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"monkey/evaluator"
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
	"monkey/repl"
	"monkey/token"
	"os"
)

func main() {
	if len(os.Args) == 1 {
		repl.Repl()
		return
	}

	input, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	ch := make(chan *token.Token)
	p := parser.New(lexer.New(string(input)[:len(input)-1], ch), ch)
	prog := p.Parse()

	if len(p.Errors) > 0 {
		for _, err := range p.Errors {
			println(err.String())
		}
		return
	}

	ret := evaluator.Eval(prog, object.NewEnv(nil))
	fmt.Println(ret.String())
}
