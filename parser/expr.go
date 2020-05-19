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
	precLogic
	precSum
	precProduct
	precPrefix
	precCall
)

var infixPrecedences = map[token.TokenType]int{
	token.Eq:    precEquals,
	token.Neq:   precEquals,
	token.Lt:    precCmp,
	token.Gt:    precCmp,
	token.Le:    precCmp,
	token.Ge:    precCmp,
	token.And:   precLogic,
	token.Or:    precLogic,
	token.Plus:  precSum,
	token.Minus: precSum,
	token.Star:  precProduct,
	token.Slash: precProduct,
}

type prefixParseFn func() ast.Expr

func (p *Parser) parseExpr(precedence int) ast.Expr {
	prefix, ok := p.prefixParseFns[p.cur.Type]
	if !ok {
		if p.cur.Type == token.Semicolon {
			return nil
		}
		p.errorf("No prefix expression found for %s", p.cur.Type.String())
		return nil
	}
	left := prefix()

	for p.cur.Type != token.EOF {
		if peekPrec, isInfix := infixPrecedences[p.peek.Type]; !isInfix || precedence >= peekPrec {
			break
		}
		p.next()
		left = p.parseInfixExpr(left)
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
	expr.Ident = p.parseIdentExpr().(*ast.IdentExpr)
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
		Args:  make([]*ast.IdentExpr, 0),
	}
	if !p.expect(token.LParen, "function expr") {
		return nil
	}
	for p.next(); p.cur.Type != token.RParen; p.next() {
		if ident := p.parseIdentExpr(); ident != nil {

		}
		ident, ok := p.parseIdentExpr().(*ast.IdentExpr)
		if !ok {
			p.errorf("Failed to parse func arguments")
			return nil
		}
		funcExpr.Args = append(funcExpr.Args, ident)
		if p.peek.Type == token.RParen {
			p.next()
			break
		}
		if !p.expect(token.Comma, "function expr") {
			return nil
		}
	}
	if !p.expect(token.LBrace, "function expr") {
		return nil
	}
	funcExpr.BlockStmt = p.parseBlockStmt()
	// p.next()
	return funcExpr
}

func (p *Parser) registerPrefixes() map[token.TokenType]prefixParseFn {
	return map[token.TokenType]prefixParseFn{
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
}
