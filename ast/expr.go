package ast

import (
	"bytes"
	"fmt"
	"monkey/token"
	"strings"
)

type Expr interface {
	Node
	exprNode()
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

type StringExpr struct {
	Token *token.Token
	Value string
}

func (s *StringExpr) exprNode() {}
func (s *StringExpr) String() string {
	return fmt.Sprintf("\"%s\"", s.Value)
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

type IncDecExpr struct {
	Token    *token.Token
	Operator string
	Ident    IdentExpr
}

func (ide *IncDecExpr) exprNode() {}
func (ide *IncDecExpr) String() string {
	return ide.Operator + ide.Ident.String()
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

type FuncExpr struct {
	Token *token.Token
	Args  []*IdentExpr
	*BlockStmt
}

func (fe *FuncExpr) exprNode() {}
func (fe *FuncExpr) String() string {
	var out bytes.Buffer
	args := make([]string, len(fe.Args))
	for i, arg := range fe.Args {
		args[i] = arg.String()
	}
	out.WriteString("fn(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(") ")
	out.WriteString(fe.BlockStmt.String())
	return out.String()
}

type IfExpr struct {
	Token *token.Token
	Cond  Expr
	Then  Stmt
	Else  Stmt
}

func (ie *IfExpr) exprNode() {}
func (ie *IfExpr) String() string {
	var out bytes.Buffer
	out.WriteString("if ")
	out.WriteString(ie.Cond.String())
	out.WriteString(" ")
	out.WriteString(ie.Then.String())
	if ie.Else != nil {
		out.WriteString(" else ")
		out.WriteString(ie.Else.String())
	}
	return out.String()
}

type FuncCallExpr struct {
	Token *token.Token
	Func  Expr
	Args  []Expr
}

func (fce *FuncCallExpr) exprNode() {}
func (fce *FuncCallExpr) String() string {
	var out bytes.Buffer
	out.WriteString(fce.Func.String())
	out.WriteString("(")
	args := make([]string, 0, len(fce.Args))
	for _, arg := range fce.Args {
		args = append(args, arg.String())
	}
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")
	return out.String()
}
