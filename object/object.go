package object

type ObjectType uint

const (
	ObjTypeNull = iota
	ObjTypeInt
	ObjTypeBool
	ObjTypeString
)

type Object interface {
	Type() ObjectType
}

type (
	ObjNull   struct{}
	ObjInt    struct{ Value int64 }
	ObjBool   struct{ Value bool }
	ObjString struct{ Value string }
)

func (o *ObjNull) Type() ObjectType   { return ObjTypeNull }
func (o *ObjInt) Type() ObjectType    { return ObjTypeInt }
func (o *ObjBool) Type() ObjectType   { return ObjTypeBool }
func (o *ObjString) Type() ObjectType { return ObjTypeString }
