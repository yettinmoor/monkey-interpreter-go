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

func Eval(n ast.Node, env object.Env) object.Object {
	if n == nil {
		return nullObj
	}

	switch n := n.(type) {
	case *ast.Program:
		return evalProgram(n, env)

	case *ast.LetStmt:
		val := Eval(n.Value, env)
		if isError(val) {
			return val
		}
		env.Set(n.Name.Value, val)
		return nullObj

	case *ast.ReturnStmt:
		return propagate(Eval(n.Value, env))
	case *ast.BlockStmt:
		return evalBlockStmt(n.Stmts, env)
	case *ast.ExprStmt:
		return Eval(n.Expr, env)

	case *ast.IdentExpr:
		if val, ok := env.Get(n.Value); ok {
			return val
		}
		return errorf("identifier not found: %s", n.Value)

	case *ast.PrefixExpr:
		right := Eval(n.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpr(n.Operator, right)

	case *ast.InfixExpr:
		left := Eval(n.Left, env)
		if isError(left) {
			return left
		}
		right := Eval(n.Right, env)
		if isError(right) {
			return right
		}
		return evalInfixExpr(n.Operator, left, right)

	case *ast.IfExpr:
		cond := Eval(n.Cond, env)
		if isError(cond) {
			return cond
		}
		then := Eval(n.Then, env)
		if isError(then) {
			return then
		}
		else_ := Eval(n.Else, env)
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

func evalProgram(n *ast.Program, env object.Env) object.Object {
	var ret object.Object
	for _, stmt := range n.Stmts {
		ret = Eval(stmt, env)
		switch ret := ret.(type) {
		case *object.ObjError:
			return ret
		case *object.ObjReturn:
			return ret.Value
		}
	}
	return ret
}

func evalBlockStmt(stmts []*ast.Stmt, env object.Env) object.Object {
	env = object.NewEnv(&env)
	for _, stmt := range stmts {
		if ret, ok := Eval(*stmt, env).(propagate); ok {
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
