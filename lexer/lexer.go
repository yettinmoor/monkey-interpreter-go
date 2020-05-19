package lexer

import (
	"monkey/token"
	"strings"
	"unicode/utf8"
)

type Lexer struct {
	input      string
	start, pos int
	r          rune
	ch         chan<- *token.Token
	row        int
	lastRowPos int
}

func New(input string, ch chan<- *token.Token) *Lexer {
	l := &Lexer{
		input: input,
		ch:    ch,
		row:   1,
	}
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
	if l.pos > len(l.input) {
		return l.input[l.start:]
	}
	s := l.input[l.start:l.pos]
	_, w := utf8.DecodeLastRuneInString(s)
	return s[:len(s)-w]
}

func (l *Lexer) consume() string {
	s := l.read()
	w := utf8.RuneLen(l.r)
	l.start = l.pos - w
	return s
}

func (l *Lexer) emit(t token.TokenType) {
	col := 1 + l.start - l.lastRowPos
	l.ch <- &token.Token{
		Type:    t,
		Literal: l.consume(),
		Row:     l.row,
		Col:     col,
	}
}

func (l *Lexer) Parse() {
	defer close(l.ch)
	defer l.emit(token.EOF)

	for l.r != 0 {
		l.readWhile(func(r rune) bool {
			if r == '\n' {
				l.row++
				l.lastRowPos = l.pos + 1
			}
			return IsWhitespace(r)
		})
		l.consume()

		switch {
		case IsValidIdentifierHead(l.r):
			l.readWhile(IsValidIdentifierRune)
			if keywordType, ok := token.Keywords[l.read()]; ok {
				l.emit(keywordType)
			} else {
				l.emit(token.Ident)
			}

		case IsNum(l.r):
			l.readWhile(IsNum)
			l.emit(token.Int)

		case l.r == '"':
			l.step()
			l.emit(token.DQuote)
			escape := false
			l.readWhile(func(r rune) bool {
				ret := (escape || r != '"') && r != 0
				escape = !escape && l.r == '\\'
				return ret
			})
			l.emit(token.String)
			if l.r == '"' {
				l.step()
				l.emit(token.DQuote)
			}

		default:
			var curr string
			l.readWhile(func(r rune) bool {
				next := curr + string(r)
				for tok := range token.SymToks {
					if strings.HasPrefix(tok, next) {
						curr = next
						return true
					}
				}
				return false
			})
			if symTok, isSymTok := token.SymToks[curr]; isSymTok {
				l.emit(symTok)
			} else {
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
