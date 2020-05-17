package parser

import (
	"fmt"
	"monkey/ast"
	"monkey/lexer"
	"monkey/token"
)

type Parser struct {
	l                 *lexer.Lexer
	ch                <-chan *token.Token
	cur, peek, buffer *token.Token
	Errors            []string
	prefixParseFns    map[token.TokenType]prefixParseFn
	infixParseFns     map[token.TokenType]bool
}

func New(l *lexer.Lexer, ch <-chan *token.Token) *Parser {
	go l.Parse()
	cur := <-ch
	peek := <-ch
	p := &Parser{
		l:    l,
		ch:   ch,
		cur:  cur,
		peek: peek,
		infixParseFns: map[token.TokenType]bool{
			token.Eq:    true,
			token.Neq:   true,
			token.Lt:    true,
			token.Gt:    true,
			token.Le:    true,
			token.Ge:    true,
			token.Plus:  true,
			token.Minus: true,
			token.Star:  true,
			token.Slash: true,
		},
	}
	p.prefixParseFns = p.registerPrefixes()
	return p
}

func (p *Parser) next() {
	p.cur = p.peek
	if p.buffer != nil {
		p.peek = p.buffer
		p.buffer = nil
	} else {
		p.peek = <-p.ch
	}
}

func (p *Parser) pushback() {
	p.buffer = p.peek
	p.peek = p.cur
	p.cur = &token.Token{}
}

func (p *Parser) Parse() *ast.Program {
	prog := &ast.Program{}
	for p.cur.Type != token.EOF {
		if stmt := p.parseStmt(); stmt != nil {
			prog.Stmts = append(prog.Stmts, stmt)
		}
		p.next()
	}
	return prog
}

func (p *Parser) parseStmt() ast.Stmt {
	switch p.cur.Type {
	case token.Let:
		return p.parseLetStmt()
	case token.Return:
		return p.parseReturnStmt()
	default:
		return p.parseExprStmt()
	}
}

func (p *Parser) parseLetStmt() *ast.LetStmt {
	stmt := &ast.LetStmt{Token: p.cur}

	if !p.expect(token.Ident) {
		return nil
	}
	stmt.Name = &ast.IdentExpr{Token: p.cur, Value: p.cur.Literal}

	if !p.expect(token.Assign) {
		return nil
	}
	p.next()
	stmt.Value = p.parseExpr(precLowest)
	if !p.expect(token.Semicolon) {
		return nil
	}

	return stmt
}

func (p *Parser) parseReturnStmt() *ast.ReturnStmt {
	stmt := &ast.ReturnStmt{Token: p.cur}

	p.next()
	stmt.Value = p.parseExpr(precLowest)
	if !p.expect(token.Semicolon) {
		return nil
	}

	return stmt
}

func (p *Parser) parseExprStmt() *ast.ExprStmt {
	stmt := &ast.ExprStmt{Token: p.cur}
	stmt.Expr = p.parseExpr(precLowest)

	if p.peek.Type == token.Semicolon {
		p.next()
	}
	return stmt
}

func (p *Parser) expect(t token.TokenType) bool {
	if p.peek.Type == t {
		p.next()
		return true
	}
	p.errorf("Expected token `%s`, got `%s`", t, p.peek.Type)
	return false
}

func (p *Parser) errorf(format string, a ...interface{}) {
	p.Errors = append(p.Errors, fmt.Sprintf(format, a...))
}
