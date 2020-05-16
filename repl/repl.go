package repl

import (
	"bufio"
	"fmt"
	"monkey/lexer"
	"monkey/token"
	"os"
)

func Repl() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print(">> ")
		scanned := scanner.Scan()
		if !scanned {
			return
		}
		ch := make(chan *token.Token)
		go lexer.New(scanner.Text(), ch).Parse()
		for tok := range ch {
			fmt.Printf("%+v\n", tok)
		}
	}
}
