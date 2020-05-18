package parser

import (
	"monkey/ast"
	"monkey/token"
)

func (p *Parser) parseStmt() ast.Stmt {
	switch p.cur.Type {
	case token.Let:
		return p.parseLetStmt()
	case token.Return:
		return p.parseReturnStmt()
	default:
		expr := p.parseExprStmt()
		return expr
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
