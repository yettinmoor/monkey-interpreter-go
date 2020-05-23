package ast

import (
	"bytes"
	"monkey/token"
	"strings"
)

type Stmt interface {
	Node
	stmtNode()
}

type LetStmt struct {
	Token *token.Token
	Name  *IdentExpr
	Value Expr
}

func (ls *LetStmt) stmtNode() {}
func (ls *LetStmt) String() string {
	var out bytes.Buffer
	out.WriteString(ls.Token.Literal + " ")
	out.WriteString(ls.Name.String())
	out.WriteString(" = ")
	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}
	out.WriteString(";")
	return out.String()
}

type ReturnStmt struct {
	Token *token.Token
	Value Expr
}

func (rs *ReturnStmt) stmtNode() {}
func (rs *ReturnStmt) String() string {
	var out bytes.Buffer
	out.WriteString(rs.Token.Literal + " ")
	if rs.Value != nil {
		out.WriteString(rs.Value.String())
	}
	out.WriteString(";")
	return out.String()
}

type ExprStmt struct {
	Token *token.Token
	Expr  Expr
}

func (es *ExprStmt) stmtNode() {}
func (es *ExprStmt) String() string {
	if es.Expr != nil {
		return es.Expr.String()
	}
	return ""
}

type BlockStmt struct {
	Token *token.Token
	Stmts []*Stmt
}

func (bs *BlockStmt) stmtNode() {}
func (bs *BlockStmt) String() string {
	var out bytes.Buffer
	out.WriteString("{")
	stmts := make([]string, len(bs.Stmts))
	for i, stmt := range bs.Stmts {
		stmts[i] = (*stmt).String()
	}
	out.WriteString(strings.Join(stmts, " "))
	out.WriteString("}")
	return out.String()
}
