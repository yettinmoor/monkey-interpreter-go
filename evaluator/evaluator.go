package evaluator

import (
	"monkey/ast"
	"monkey/object"
)

type ReturnValue object.Object

var (
	null   = &object.ObjNull{}
	true_  = &object.ObjBool{Value: true}
	false_ = &object.ObjBool{Value: false}
)

func Eval(n ast.Node) object.Object {
	switch n := n.(type) {
	case *ast.Program:
		return evalProgram(n)

	case *ast.ExprStmt:
		return Eval(n.Expr)
	case *ast.PrefixExpr:
		return evalPrefixExpr(n.Operator, Eval(n.Right))
	case *ast.InfixExpr:
		return evalInfixExpr(n.Operator, Eval(n.Left), Eval(n.Right))

	case *ast.IntLiteralExpr:
		return &object.ObjInt{Value: n.Value}
	default:
		panic("Invalid")
	}
}

func evalProgram(n *ast.Program) object.Object {
	var ret object.Object
	for _, stmt := range n.Stmts {
		ret = Eval(stmt)
	}
	return ret
}

func evalPrefixExpr(op string, right object.Object) object.Object {
	switch op {
	case "-":
		if right.Type() != object.ObjTypeInt {
			return null
		}
		return &object.ObjInt{Value: -right.(*object.ObjInt).Value}
	default:
		return null
	}
}

func evalInfixExpr(op string, left, right object.Object) object.Object {
	switch {
	case left.Type() == object.ObjTypeInt && right.Type() == object.ObjTypeInt:
		leftVal := left.(*object.ObjInt).Value
		rightVal := right.(*object.ObjInt).Value
		switch op {
		case "+":
			return &object.ObjInt{Value: leftVal + rightVal}
		case "-":
			return &object.ObjInt{Value: leftVal - rightVal}
		case "*":
			return &object.ObjInt{Value: leftVal * rightVal}
		case "/":
			return &object.ObjInt{Value: leftVal / rightVal}
		case "%":
			return &object.ObjInt{Value: leftVal % rightVal}
		default:
			return null
		}
	// case left.Type() == object.ObjTypeBool && right.Type() == object.ObjTypeBool:
	// 	return evalBoolInfixExpr(n.Operator, left, right)
	default:
		return null
	}
}
