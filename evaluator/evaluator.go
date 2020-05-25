package evaluator

import (
	"fmt"
	"monkey/ast"
	"monkey/object"
)

type propagate object.Object

var (
	nullObj  = &object.ObjNull{}
	trueObj  = &object.ObjBool{Value: true}
	falseObj = &object.ObjBool{Value: false}
)

func Eval(n ast.Node) object.Object {
	if n == nil {
		return nullObj
	}
	switch n := n.(type) {
	case *ast.Program:
		return evalProgram(n)

	case *ast.ReturnStmt:
		return propagate(Eval(n.Value))
	case *ast.BlockStmt:
		return evalBlockStmt(n.Stmts)
	case *ast.ExprStmt:
		return Eval(n.Expr)

	case *ast.PrefixExpr:
		right := Eval(n.Right)
		if isError(right) {
			return right
		}
		return evalPrefixExpr(n.Operator, right)

	case *ast.InfixExpr:
		left := Eval(n.Left)
		if isError(left) {
			return left
		}
		right := Eval(n.Right)
		if isError(right) {
			return right
		}
		return evalInfixExpr(n.Operator, left, right)

	case *ast.IfExpr:
		cond := Eval(n.Cond)
		if isError(cond) {
			return cond
		}
		then := Eval(n.Then)
		if isError(then) {
			return then
		}
		else_ := Eval(n.Else)
		if isError(else_) {
			return else_
		}

		if isTruthy(cond) {
			return then
		}
		return else_

	case *ast.IntLiteralExpr:
		return &object.ObjInt{Value: n.Value}
	case *ast.BoolExpr:
		if n.Value {
			return trueObj
		}
		return falseObj
	case *ast.StringExpr:
		return &object.ObjString{Value: n.Value}

	default:
		panic("Invalid")
	}
}

func evalProgram(n *ast.Program) object.Object {
	var ret object.Object
	for _, stmt := range n.Stmts {
		ret = Eval(stmt)
		switch ret := ret.(type) {
		case *object.ObjError:
			return ret
		case *object.ObjReturn:
			return ret.Value
		}
	}
	return ret
}

func evalBlockStmt(stmts []*ast.Stmt) object.Object {
	for _, stmt := range stmts {
		if ret, ok := Eval(*stmt).(propagate); ok {
			return ret
		}
	}
	return nullObj
}

func evalPrefixExpr(op string, right object.Object) object.Object {
	switch op {
	case "-":
		if right.Type() != object.ObjTypeInt {
			return errorf("Bad int prefix %s", op)
		}
		return &object.ObjInt{Value: -right.(*object.ObjInt).Value}
	case "!":
		return getBool(!isTruthy(right))
	default:
		return errorf("Bad prefix %s", op)
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
		case "==":
			return getBool(leftVal == rightVal)
		case "!=":
			return getBool(leftVal != rightVal)
		case "<":
			return getBool(leftVal < rightVal)
		case "<=":
			return getBool(leftVal <= rightVal)
		case ">":
			return getBool(leftVal > rightVal)
		case ">=":
			return getBool(leftVal >= rightVal)
		default:
			return errorf("Bad int operator %q", op)
		}
	default:
		return errorf("Bad expression: %s %s %s", left, op, right)
	}
}

func isTruthy(o object.Object) bool {
	switch o := o.(type) {
	case *object.ObjInt:
		return o.Value != 0
	case *object.ObjBool:
		return o.Value
	case *object.ObjString:
		return len(o.Value) > 0
	}
	return false
}

func getBool(b bool) *object.ObjBool {
	if b {
		return trueObj
	}
	return falseObj
}

func isError(o object.Object) bool {
	return o.Type() == object.ObjTypeError
}

func errorf(msg string, a ...interface{}) *object.ObjError {
	return &object.ObjError{Error: fmt.Sprintf(msg, a...)}
}
