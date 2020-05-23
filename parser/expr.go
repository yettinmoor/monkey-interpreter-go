package parser

import (
	"monkey/ast"
	"monkey/token"
	"strconv"
)

const (
	_ int = iota
	precLowest
	precEquals
	precCmp
	precOr
	precAnd
	precSum
	precProduct
	precPrefix
	precCall
)

var infixPrecedences = map[token.TokenType]int{
	token.Eq:     precEquals,
	token.Neq:    precEquals,
	token.Lt:     precCmp,
	token.Gt:     precCmp,
	token.Le:     precCmp,
	token.Ge:     precCmp,
	token.Or:     precOr,
	token.And:    precAnd,
	token.Plus:   precSum,
	token.Minus:  precSum,
	token.Star:   precProduct,
	token.Slash:  precProduct,
	token.Modulo: precProduct,
	token.LParen: precCall,
}

func (p *Parser) parseExpr(prec int) ast.Expr {
	prefix, ok := p.prefixParseFns[p.cur.Type]
	if !ok {
		if p.cur.Type != token.Semicolon {
			p.errorf("No prefix expression found for %s", p.cur.Type.String())
		}
		return nil
	}
	left := prefix()
	for {
		if peek, _ := infixPrecedences[p.peek.Type]; prec >= peek {
			break
		}
		p.next()
		switch p.cur.Type {
		case token.LParen:
			left = p.parseFuncCallExpr(left)
		default:
			left = p.parseInfixExpr(left)
		}
	}
	return left
}

func (p *Parser) parseIdentExpr() ast.Expr {
	return &ast.IdentExpr{Token: p.cur, Value: p.cur.Literal}
}

func (p *Parser) parseIntLiteralExpr() ast.Expr {
	value, err := strconv.ParseInt(p.cur.Literal, 0, 64)
	if err != nil {
		p.errorf("Failed to parse %q as int", p.cur.Literal)
		return nil
	}
	return &ast.IntLiteralExpr{Token: p.cur, Value: value}
}

func (p *Parser) parseStringExpr() ast.Expr {
	expr := &ast.StringExpr{Token: p.cur}
	p.next()
	expr.Value = p.cur.Literal
	if !p.expect(token.DQuote, "string expr") {
		return nil
	}
	return expr
}

func (p *Parser) parseBoolExpr() ast.Expr {
	return &ast.BoolExpr{Token: p.cur, Value: p.cur.Literal == "true"}
}

func (p *Parser) parsePrefixExpr() ast.Expr {
	expr := &ast.PrefixExpr{Token: p.cur, Operator: p.cur.Literal}
	p.next()
	expr.Right = p.parseExpr(precPrefix)
	return expr
}

func (p *Parser) parseIncDecExpr() ast.Expr {
	expr := &ast.IncDecExpr{Token: p.cur, Operator: p.cur.Literal}
	if !p.expect(token.Ident, "inc-dec stmt") {
		return nil
	}
	expr.Ident = *p.parseIdentExpr().(*ast.IdentExpr)
	return expr
}

func (p *Parser) parseInfixExpr(left ast.Expr) ast.Expr {
	expr := &ast.InfixExpr{
		Token:    p.cur,
		Operator: p.cur.Literal,
		Left:     left,
	}
	prec := infixPrecedences[p.cur.Type]
	p.next()
	expr.Right = p.parseExpr(prec)
	return expr
}

func (p *Parser) parseGroupedExpr() ast.Expr {
	p.next()
	expr := p.parseExpr(precLowest)
	if !p.expect(token.RParen, "grouped expr") {
		return nil
	}
	return expr
}

func (p *Parser) parseFuncExpr() ast.Expr {
	funcExpr := &ast.FuncExpr{
		Token: p.cur,
		Args:  make([]ast.IdentExpr, 0),
	}
	if !p.expect(token.LParen, "function expr") {
		return nil
	}
	for p.next(); p.cur.Type != token.RParen; p.next() {
		switch p.cur.Type {
		case token.Ident:
			ident, _ := p.parseIdentExpr().(*ast.IdentExpr)
			funcExpr.Args = append(funcExpr.Args, *ident)
		case token.Comma:
			break
		default:
			p.errorf(
				"While parsing function arguments: Expected expression or comma, got %q %q",
				p.cur.Type.String(),
				p.cur.Literal,
			)
			return nil
		}
	}
	if !p.expect(token.LBrace, "function expr") {
		return nil
	}
	funcExpr.BlockStmt = p.parseBlockStmt()
	return funcExpr
}

func (p *Parser) parseIfExpr() ast.Expr {
	ifExpr := &ast.IfExpr{Token: p.cur}
	// if !p.expect(token.LParen, "if expr") {
	// 	return nil
	// }
	// ifExpr.Cond = p.parseGroupedExpr()
	p.next()
	ifExpr.Cond = p.parseExpr(precLowest)
	p.next()
	ifExpr.Then = p.parseStmt()
	if p.accept(token.Else) {
		p.next()
		ifExpr.Else = p.parseStmt()
	}
	return ifExpr
}

func (p *Parser) parseFuncCallExpr(f ast.Expr) ast.Expr {
	callExpr := &ast.FuncCallExpr{
		Token: p.cur,
		Func:  f,
	}
	if p.accept(token.RParen) {
		return callExpr
	}
	p.next()
	callExpr.Args = append(callExpr.Args, p.parseExpr(precLowest))
	for p.accept(token.Comma) {
		p.next()
		callExpr.Args = append(callExpr.Args, p.parseExpr(precLowest))
	}
	if !p.expect(token.RParen, "func call expr") {
		return nil
	}
	return callExpr
}
