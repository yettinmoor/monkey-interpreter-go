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
	case token.LBrace:
		return p.parseBlockStmt()
	default:
		return p.parseExprStmt()
	}
}

func (p *Parser) parseLetStmt() *ast.LetStmt {
	stmt := &ast.LetStmt{Token: p.cur}

	if !p.expect(token.Ident, "let stmt") {
		return nil
	}
	stmt.Name = &ast.IdentExpr{Token: p.cur, Value: p.cur.Literal}

	if !p.expect(token.Assign, "let stmt") {
		return nil
	}
	p.next()
	stmt.Value = p.parseExpr(precLowest)
	if !p.expect(token.Semicolon, "let stmt") {
		return nil
	}

	return stmt
}

func (p *Parser) parseReturnStmt() *ast.ReturnStmt {
	stmt := &ast.ReturnStmt{Token: p.cur}

	if p.accept(token.Semicolon) {
		return stmt
	}

	p.next()
	stmt.Value = p.parseExpr(precLowest)
	if !p.expect(token.Semicolon, "return stmt") {
		return nil
	}

	return stmt
}

func (p *Parser) parseExprStmt() *ast.ExprStmt {
	stmt := &ast.ExprStmt{Token: p.cur}
	stmt.Expr = p.parseExpr(precLowest)
	return stmt
}

func (p *Parser) parseBlockStmt() *ast.BlockStmt {
	block := &ast.BlockStmt{Token: p.cur, Stmts: make([]*ast.Stmt, 0)}
	for p.next(); p.cur.Type != token.RBrace; p.next() {
		stmt := p.parseStmt()
		if stmt == nil {
			return nil
		}
		block.Stmts = append(block.Stmts, &stmt)
	}
	return block
}
