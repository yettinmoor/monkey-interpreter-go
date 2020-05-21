package ast

import (
	"bytes"
)

type Node interface {
	String() string
}

type Program struct {
	Stmts []Stmt
}

func (p *Program) String() string {
	var out bytes.Buffer
	for _, s := range p.Stmts {
		if s != nil {
			out.WriteString(s.String())
		}
	}
	return out.String()
}
