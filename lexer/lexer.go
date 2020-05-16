package lexer

import (
	"monkey/token"
	"strings"
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
	// log.Printf("\tEMITTING %q (type %q, pos.s %d to %d)", l.read(), t, l.start, l.pos)
	l.ch <- &token.Token{Type: t, Literal: l.consume()}
}

func (l *Lexer) Parse() {
	defer close(l.ch)
	defer l.emit(token.EOF)

	for l.r != 0 {
		l.readWhile(isWhitespace)
		l.consume()
		// log.Printf("CURRENT char %q (pos %d)", l.r, l.start)

		switch {

		case isAlpha(l.r):
			// log.Printf("\tALPHA")
			l.readWhile(isAlphaNum)
			if keywordType, ok := token.Keywords[l.read()]; ok {
				l.emit(keywordType)
			} else {
				l.emit(token.Ident)
			}

		case isNum(l.r):
			// log.Printf("\tNUM")
			l.readWhile(isNum)
			l.emit(token.Int)

		case l.r == '"':
			l.readWhile(func(r rune) bool {
				return r != '"'
			})

		default:
			var curr string
			l.readWhile(func(r rune) bool {
				next := curr + string(r)
				// log.Printf("\tSYMBOL? %q", next)
				for tok := range token.SymToks {
					if strings.HasPrefix(tok, next) {
						curr = next
						return true
					}
				}
				// return len(symToks) >= 1
				return false
			})
			if symTok, isSymTok := token.SymToks[curr]; isSymTok {
				l.emit(symTok)
			} else {
				// log.Printf("illegal")
				l.emit(token.Illegal)
			}
		}
	}
}

func (l *Lexer) readWhile(f func(r rune) bool) {
	for f(l.r) {
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
