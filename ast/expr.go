package ast

import (
	"bytes"
	"monkey/token"
)

type Expr interface {
	exprNode()
	String() string
}

type IdentExpr struct {
	Token *token.Token
	Value string
}

func (i *IdentExpr) exprNode() {}
func (i *IdentExpr) String() string {
	return i.Value
}

type IntLiteralExpr struct {
	Token *token.Token
	Value int64
}

func (i *IntLiteralExpr) exprNode() {}
func (i *IntLiteralExpr) String() string {
	return i.Token.Literal
}

type BoolExpr struct {
	Token *token.Token
	Value bool
}

func (b *BoolExpr) exprNode() {}
func (b *BoolExpr) String() string {
	return b.Token.Literal
}

type PrefixExpr struct {
	Token    *token.Token
	Operator string
	Right    Expr
}

func (pe *PrefixExpr) exprNode() {}
func (pe *PrefixExpr) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")
	return out.String()
}

type InfixExpr struct {
	Token    *token.Token
	Operator string
	Left     Expr
	Right    Expr
}

func (ie *InfixExpr) exprNode() {}
func (ie *InfixExpr) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString(ie.Operator)
	out.WriteString(ie.Right.String())
	out.WriteString(")")
	return out.String()
}
