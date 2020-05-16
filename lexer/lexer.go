package lexer

import (
	"fmt"
	"monkey/token"
	"unicode/utf8"
)

type Lexer struct {
	input string
	start int
	pos   int
	r     rune
	ch    chan<- *token.Token
}

func New(input string, ch chan<- *token.Token) *Lexer {
	l := &Lexer{input: input, ch: ch}
	l.step()
	return l
}

func (l *Lexer) step() {
	if l.pos >= len(l.input) {
		l.r = 0
		l.pos++
	} else {
		r, w := utf8.DecodeRuneInString(l.input[l.pos:])
		l.r = r
		l.pos += w
	}
}

func (l *Lexer) peek() rune {
	if l.pos >= len(l.input) {
		return 0
	}
	ch, _ := utf8.DecodeRuneInString(l.input[l.pos:])
	return ch
}

func (l *Lexer) read() string {
	return l.input[l.start : l.pos-1]
}

func (l *Lexer) consume() string {
	s := l.read()
	l.start = l.pos - 1
	return s
}

func (l *Lexer) emit(t token.TokenType) {
	fmt.Printf("emitting %q (type %q, pos.s %d to %d)\n", l.read(), t, l.start, l.pos)
	l.ch <- &token.Token{Type: t, Literal: l.consume()}
}

func (l *Lexer) Parse() {
	defer close(l.ch)
	defer l.emit(token.EOF)

	for l.r != 0 {
		l.readWhile(isWhitespace)
		fmt.Printf("consuming %q\n", l.read())
		l.consume()
		fmt.Printf("curr input is %q\n", l.input[l.start:])
		fmt.Printf("curr char %q (pos %d) is ", l.r, l.start)

		switch {

		case isAlpha(l.r):
			fmt.Printf("alpha\n")
			l.readWhile(isAlphaNum)
			if keywordType, ok := token.Keywords[l.read()]; ok {
				l.emit(keywordType)
			} else {
				l.emit(token.Ident)
			}

		case isNum(l.r):
			fmt.Printf("num\n")
			l.readWhile(isNum)
			l.emit(token.Int)

		default:
			if singleTokType, isSingleCh := token.SingleChToks[l.r]; isSingleCh {
				l.step()
				l.emit(singleTokType)
			} else {
				fmt.Printf("illegal\n")
				l.emit(token.Illegal)
			}
		}
	}
}

func (l *Lexer) readWhile(f func(r rune) bool) {
	for f(l.r) {
		fmt.Printf("\treading %q at pos %d\n", l.r, l.pos-1)
		l.step()
	}
}

func isAlpha(r rune) bool {
	return 'a' <= r && r <= 'z' || 'A' <= r && r <= 'Z' || r == '_'
}

func isNum(r rune) bool {
	return '0' <= r && r <= '9'
}

func isAlphaNum(r rune) bool {
	return isAlpha(r) || isNum(r)
}

func isWhitespace(r rune) bool {
	return r == ' ' || r == '\t' || r == '\n' || r == '\r'
}
