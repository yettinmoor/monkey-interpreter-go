package ast

import (
	"monkey/token"
	"testing"
)

func TestString(t *testing.T) {
	program := &Program{
		Stmts: []Stmt{
			&LetStmt{
				Token: &token.Token{Type: token.Let, Literal: "let"},
				Name: &IdentExpr{
					Token: &token.Token{Type: token.Ident, Literal: "myVar"},
					Value: "myVar",
				},
				Value: &IdentExpr{
					Token: &token.Token{Type: token.Ident, Literal: "otherVar"},
					Value: "otherVar",
				},
			},
		},
	}

	if program.String() != "let myVar = otherVar;" {
		t.Errorf("program.String() returned wrong value: %q", program.String())
	}
}
