package object

import (
	"fmt"
	"go/ast"
)

type ObjectType uint

const (
	ObjTypeError ObjectType = iota
	ObjTypeNull
	ObjTypeReturn
	ObjTypeInt
	ObjTypeBool
	ObjTypeString
	ObjTypeIdent
	ObjTypeFunc
)

type Object interface {
	Type() ObjectType
	String() string
}

type (
	ObjError  struct{ Error string }
	ObjNull   struct{}
	ObjReturn struct{ Value Object }
	ObjInt    struct{ Value int64 }
	ObjBool   struct{ Value bool }
	ObjString struct{ Value string }
	ObjFunc   struct {
		Args  []ast.Ident
		Stmts []*ast.Stmt
	}
)

func (o *ObjError) Type() ObjectType { return ObjTypeError }
func (o *ObjError) String() string   { return fmt.Sprintf("<Error: %s>", o.Error) }

func (o *ObjNull) Type() ObjectType { return ObjTypeNull }
func (o *ObjNull) String() string   { return "null" }

func (o *ObjReturn) Type() ObjectType { return ObjTypeReturn }
func (o *ObjReturn) String() string   { return fmt.Sprint(o.Value) }

func (o *ObjInt) Type() ObjectType { return ObjTypeInt }
func (o *ObjInt) String() string   { return fmt.Sprint(o.Value) }

func (o *ObjBool) Type() ObjectType { return ObjTypeBool }
func (o *ObjBool) String() string   { return fmt.Sprint(o.Value) }

func (o *ObjString) Type() ObjectType { return ObjTypeString }
func (o *ObjString) String() string   { return fmt.Sprintf("\"%s\"", o.Value) }

func (o *ObjFunc) Type() ObjectType { return ObjTypeFunc }
func (o *ObjFunc) String() string   { return "<function>" }
