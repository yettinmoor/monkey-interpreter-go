package parser

import (
	"fmt"
	"monkey/ast"
	"monkey/lexer"
	"monkey/token"
)

type Parser struct {
	l              *lexer.Lexer
	ch             <-chan *token.Token
	cur, peek      *token.Token
	Errors         []ParserError
	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]bool
}

type ParserError struct {
	row, col int
	msg      string
}

func (pe *ParserError) String() string {
	return fmt.Sprintf("Row %d, col %d: %s", pe.row, pe.col, pe.msg)
}

func New(l *lexer.Lexer, ch <-chan *token.Token) *Parser {
	go l.Parse()
	p := &Parser{
		l:    l,
		ch:   ch,
		peek: <-ch,
	}
	p.prefixParseFns = map[token.TokenType]prefixParseFn{
		token.Ident:     p.parseIdentExpr,
		token.Int:       p.parseIntLiteralExpr,
		token.DQuote:    p.parseStringExpr,
		token.Bang:      p.parsePrefixExpr,
		token.Minus:     p.parsePrefixExpr,
		token.Increment: p.parseIncDecExpr,
		token.Decrement: p.parseIncDecExpr,
		token.True:      p.parseBoolExpr,
		token.False:     p.parseBoolExpr,
		token.LParen:    p.parseGroupedExpr,
		token.Function:  p.parseFuncExpr,
	}
	p.next()
	return p
}

func (p *Parser) next() {
	if p.peek.Type != token.EOF {
		p.cur, p.peek = p.peek, <-p.ch
	}
}

func (p *Parser) Parse() *ast.Program {
	prog := &ast.Program{}
	for p.peek.Type != token.EOF {
		if stmt := p.parseStmt(); stmt != nil {
			prog.Stmts = append(prog.Stmts, stmt)
		}
		p.next()
	}
	return prog
}

func (p *Parser) expect(t token.TokenType, caller string) bool {
	if p.peek.Type == t {
		p.next()
		return true
	}
	p.errorf("While parsing %s: Expected token `%s`, got `%s`", caller, t.String(), p.peek.Type.String())
	p.next()
	return false
}

func (p *Parser) errorf(format string, a ...interface{}) {
	pe := ParserError{row: p.cur.Row, col: p.cur.Col, msg: fmt.Sprintf(format, a...)}
	p.Errors = append(p.Errors, pe)
}
