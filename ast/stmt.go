package ast

import (
	"bytes"
	"monkey/token"
)

type Stmt interface {
	stmtNode()
	String() string
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
