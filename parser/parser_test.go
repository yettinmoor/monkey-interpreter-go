package parser

import (
	"fmt"
	"monkey/ast"
	"monkey/lexer"
	"monkey/token"
	"testing"
)

func setup(t *testing.T, input string) *ast.Program {
	ch := make(chan *token.Token)
	l := lexer.New(input, ch)
	p := New(l, ch)

	program := p.Parse()
	if program == nil {
		t.Fatalf("Parse() returned nil")
	}
	if p.Errors != nil {
		t.Errorf("parser has %d errors", len(p.Errors))
		for _, err := range p.Errors {
			t.Errorf(err.String())
		}
		t.FailNow()
	}
	return program
}

func TestParserError(t *testing.T) {
	input := `let x;
	return 3`
	ch := make(chan *token.Token)
	l := lexer.New(input, ch)
	p := New(l, ch)

	program := p.Parse()
	if program == nil {
		t.Fatalf("Parse() returned nil")
	}
	if p.Errors == nil {
		t.Fatalf("ParserError did not catch errors")
	}
	tests := []struct {
		row, col int
	}{
		{1, 5},
		{2, 8},
	}
	if len(p.Errors) != len(tests) {
		t.Errorf("Expected %d errors, caught %d", len(tests), len(p.Errors))
		for _, err := range p.Errors {
			t.Log(err.String())
		}
		t.FailNow()
	}
	for i, err := range p.Errors {
		t.Log(err.String())
		if err.row != tests[i].row {
			t.Errorf("Error #%d should be at row %d, got %d", i+1, tests[i].row, err.row)
		}
		if err.col != tests[i].col {
			t.Errorf("Error #%d should be at col %d, got %d", i+1, tests[i].col, err.col)
		}
	}
}

func TestLetStmts(t *testing.T) {
	input := `let x = 5;
	let y = 10 + 20;
	let foo = 5 * 1;`

	program := setup(t, input)

	if ls := len(program.Stmts); ls != 3 {
		t.Fatalf("Expected 3 let stmts, got %d", ls)
	}

	tests := []struct {
		expectedIdent string
	}{
		{"x"},
		{"y"},
		{"foo"},
	}

	for i, tt := range tests {
		if letStmt, ok := program.Stmts[i].(*ast.LetStmt); !ok {
			t.Errorf("Not return statement, got %T", program.Stmts[i])
		} else if name := letStmt.Name.Value; name != tt.expectedIdent {
			t.Errorf("Expected name %s, got %s", tt.expectedIdent, name)
		}
	}
}

func TestReturnStmts(t *testing.T) {
	input := `return;
	return 5;
	return 3*(5+6);`

	program := setup(t, input)

	if ls := len(program.Stmts); ls != 3 {
		t.Fatalf("Expected 3 return stmts, got %d", ls)
	}

	tests := []struct {
		value string
	}{
		{""},
		{"5"},
		{"(3*(5+6))"},
	}

	for i, tt := range tests {
		if retStmt, ok := program.Stmts[i].(*ast.ReturnStmt); !ok {
			t.Fatalf("Expected return stmt, got %T", program.Stmts[i])
		} else if tt.value != "" && retStmt.Value.String() != tt.value {
			t.Fatalf("Expected returned value %s, got %s", tt.value, retStmt.Value.String())
		}
	}
}

func TestIdentExpr(t *testing.T) {
	tests := []struct {
		input       string
		shouldParse bool
	}{
		{`foobar`, true},
		{`abc123`, true},
		{`日本語`, true},
		{`UñïĸøDə`, true},
		{`abc-123`, false},
		{`123abc`, false},
		{`x + y`, false},
		{`bad variable name`, false},
	}

	for _, tt := range tests {
		ch := make(chan *token.Token)
		l := lexer.New(fmt.Sprintf("let %s = 0;", tt.input), ch)
		p := New(l, ch)
		program := p.Parse()

		if program == nil {
			t.Fatalf("Parse() returned nil")
		}
		if !tt.shouldParse && p.Errors == nil {
			t.Fatalf("Identifier incorrectly allowed: %q", tt.input)
		}
		if tt.shouldParse && p.Errors != nil {
			for _, err := range p.Errors {
				t.Errorf(err.String())
			}
			t.FailNow()
		}
		if tt.shouldParse {
			if letStmt, ok := program.Stmts[0].(*ast.LetStmt); !ok {
				t.Fatalf("Failed to parse let stmt, got %T", program.Stmts[0])
			} else if letStmt.Name.Value != tt.input {
				t.Errorf("Token literal not %s, got %s", tt.input, letStmt.Name.Value)
			}
		}
	}
}

func TestIntLiteralExpr(t *testing.T) {
	input := `5`

	program := setup(t, input)

	if len(program.Stmts) != 1 {
		t.Fatalf("Expected 1 expr, got %d", len(program.Stmts))
	}
	if stmt, ok := program.Stmts[0].(*ast.ExprStmt); !ok {
		t.Fatalf("Not exprstmt, got %T", program.Stmts[0])
	} else if intExpr, ok := stmt.Expr.(*ast.IntLiteralExpr); !ok {
		t.Fatalf("not ident expr, got %T", stmt.Expr)
	} else if intExpr.Value != 5 {
		t.Errorf("token value not 5, got %s", intExpr.Token.Literal)
	}
}

func TestStringExpr(t *testing.T) {
	input := "let x = \"escaped \n string! \\\" quotes and stuff \\\"\";"

	program := setup(t, input)

	tests := []string{
		"escaped \n string! \\\" quotes and stuff \\\"",
	}

	if len(program.Stmts) != len(tests) {
		t.Fatalf("Expected %d expr, got %d", len(tests), len(program.Stmts))
	}
	for i, tt := range tests {
		if letStmt, ok := program.Stmts[i].(*ast.LetStmt); !ok {
			t.Errorf("Stmt %d: expected *ast.ExprStmt, got %T", i, program.Stmts[i])
		} else if strExpr, ok := letStmt.Value.(*ast.StringExpr); !ok {
			t.Errorf("Expr stmt %d is not a string expr", i)
		} else if strExpr.Value != tt {
			t.Errorf("Expected string %q, got %q", tt, strExpr.Value)
		}
	}
}

func TestPrefixExpr(t *testing.T) {
	tests := []struct {
		input, op string
		value     int64
	}{
		{"!5", "!", 5},
		{"-15", "-", 15},
	}

	for i, tt := range tests {
		program := setup(t, tt.input)
		t.Logf("Test %d", i)
		if len(program.Stmts) != 1 {
			t.Fatalf("Expected 1 expr, got %d", len(program.Stmts))
		}

		if stmt, ok := program.Stmts[0].(*ast.ExprStmt); !ok {
			t.Fatalf("Not exprstmt, got %T", program.Stmts[0])
		} else if expr, ok := stmt.Expr.(*ast.PrefixExpr); !ok {
			t.Fatalf("not prefix expr, got %T", stmt.Expr)
		} else {
			if expr.Operator != tt.op {
				t.Errorf("expression operator not %s, got %s", tt.op, expr.Operator)
			}
			testIntLit(t, expr.Right, tt.value)

		}
	}
}

func TestIncDecExpr(t *testing.T) {
	tests := []struct {
		input string
		op    token.TokenType
		ident string
	}{
		{"--x", token.Decrement, "x"},
		{"++y", token.Increment, "y"},
	}

	for _, tt := range tests {
		program := setup(t, tt.input)
		if len(program.Stmts) != 1 {
			t.Fatalf("Expected 1 stmt-expr, got %d", len(program.Stmts))
		}

		if stmt, ok := program.Stmts[0].(*ast.ExprStmt); !ok {
			t.Fatalf("Not exprstmt, got %T", program.Stmts[0])
		} else if expr, ok := stmt.Expr.(*ast.IncDecExpr); !ok {
			t.Fatalf("not incdec expr, got %T", stmt.Expr)
		} else {
			if expr.Token.Type != tt.op {
				t.Errorf("Expression operator not %s, got %s", tt.op.String(), expr.Token.Type.String())
			}
			if expr.Ident.Value != tt.ident {
				t.Errorf("Expression value not %s, got %s", tt.ident, expr.Ident.Value)
			}
		}
	}
}
func TestInfixExpr(t *testing.T) {
	tests := []struct {
		input  string
		lvalue int64
		op     string
		rvalue int64
	}{
		{"5 + 5", 5, "+", 5},
		{"5 * 5", 5, "*", 5},
		{"5 / 5", 5, "/", 5},
		{"5 == 5", 5, "==", 5},
		{"5 <= 5", 5, "<=", 5},
	}

	for i, tt := range tests {
		program := setup(t, tt.input)
		t.Logf("Test %d", i)
		if len(program.Stmts) != 1 {
			t.Fatalf("Expected 1 expr, got %d\n", len(program.Stmts))
		}
		if stmt, ok := program.Stmts[0].(*ast.ExprStmt); !ok {
			t.Fatalf("Not exprstmt, got %T", program.Stmts[0])
		} else if expr, ok := stmt.Expr.(*ast.InfixExpr); !ok {
			t.Fatalf("not infix expr, got %T", stmt.Expr)
		} else {
			if expr.Operator != tt.op {
				t.Errorf("expression operator not %s, got %s", tt.op, expr.Operator)
			}
			testIntLit(t, expr.Right, tt.lvalue)
			testIntLit(t, expr.Right, tt.rvalue)
		}
	}
}

func TestFuncExpr(t *testing.T) {
	input := `let void = fn() {  };
	let square = fn(x) { return x*x; };
	let avg = fn(a, b) { let a = 3; let b = 5; return (a + b) / 2; };`

	program := setup(t, input)

	tests := []struct {
		args   []string
		nStmts int
	}{
		{[]string{}, 0},
		{[]string{"x"}, 1},
		{[]string{"a", "b"}, 3},
	}

	for i, tt := range tests {
		if len(program.Stmts) != len(tests) {
			t.Fatalf("Expected %d expr, got %d\n", len(tests), len(program.Stmts))
		}
		if stmt, ok := program.Stmts[i].(*ast.LetStmt); !ok {
			t.Fatalf("Not exprstmt, got %T", program.Stmts[i])
		} else if funcExpr, ok := stmt.Value.(*ast.FuncExpr); !ok {
			t.Fatalf("Not func expr, got %T", stmt.Value)
		} else {
			if funcExpr.Token.Literal != "fn" {
				t.Errorf("Expression toklit not `fn`, got %q", funcExpr.Token.Literal)
			}
			if len(funcExpr.Args) != len(tt.args) {
				t.Errorf("Expected %d args, got %d", len(tt.args), len(funcExpr.Args))
			} else {
				for i, arg := range funcExpr.Args {
					if tt.args[i] != arg.String() {
						t.Errorf("Arg %d: expected %q, got %q", i, tt.args[i], arg.String())
					}
				}
			}
			if len(funcExpr.Stmts) != tt.nStmts {
				t.Errorf("Expected %d stmts, got %d", tt.nStmts, len(funcExpr.Args))
			}
		}
	}
}

func TestBlockStmt(t *testing.T) {
	input := `
	{
		let x = 3;
		let y = 4;
		let z = 5;
	}
	{
		let sq = fn(x) { return x*x; };
		let z = 5;
	}
	{
		let a = 3+3;
	}
	{}`

	program := setup(t, input)

	tests := []int{3, 2, 1, 0}

	if len(program.Stmts) != len(tests) {
		t.Fatalf("Expected %d expr, got %d\n", len(tests), len(program.Stmts))
	}

	for i, tt := range tests {
		if stmt, ok := program.Stmts[i].(*ast.BlockStmt); !ok {
			t.Fatalf("Not block stmt, got %T", program.Stmts[i])
		} else if len(stmt.Stmts) != tt {
			t.Errorf("Block has %d stmts, expected %d", len(stmt.Stmts), tt)
		}
	}
}

func TestIfExpr(t *testing.T) {
	input := `let y = if (x < 3) true else { let z = 3; z + 1 };`
	program := setup(t, input)

	if len(program.Stmts) != 1 {
		t.Fatalf("Expected 1 stmt, got %d", len(program.Stmts))
	}

	letStmt, _ := program.Stmts[0].(*ast.LetStmt)
	if ifExpr, ok := letStmt.Value.(*ast.IfExpr); !ok {
		t.Fatalf("Expected if expr, got %T", letStmt.Value)
	} else if ifExpr.Cond.String() != "(x<3)" {
		t.Errorf("Condition should be %q, got %q", "(x<3)", ifExpr.Cond.String())
	} else if _, ok := ifExpr.Then.(*ast.ExprStmt); !ok {
		t.Errorf("Then stmt should be expr-stmt, got %T", ifExpr.Then)
	} else if _, ok := ifExpr.Else.(*ast.BlockStmt); !ok {
		t.Errorf("Else stmt should be block stmt, got %T", ifExpr.Else)
	}
}

func TestFuncCallExpr(t *testing.T) {
	input := `add(1, 2, 3+4, 5*6, sub(7, 8))`
	program := setup(t, input)

	if len(program.Stmts) != 1 {
		t.Fatalf("Expected 1 stmt, got %d", len(program.Stmts))
	}

	exprStmt, _ := program.Stmts[0].(*ast.ExprStmt)
	if funcCallExpr, ok := exprStmt.Expr.(*ast.FuncCallExpr); !ok {
		t.Fatalf("Expected func call expr, got %T", exprStmt.Expr)
	} else if funcCallExpr.Func.String() != "add" {
		t.Errorf("Expected func call identifier %q, got %q", "add", funcCallExpr.Func.String())
	} else if len(funcCallExpr.Args) != 5 {
		t.Errorf("Expected %d args to func call, got %d", 5, len(funcCallExpr.Args))
	}
}

func testIntLit(t *testing.T, expr ast.Expr, value int64) {
	if intExpr, ok := expr.(*ast.IntLiteralExpr); !ok {
		t.Fatalf("Not int expr, got %T", expr)
	} else if intExpr.Value != value {
		t.Errorf("Expected value %d, got %d", value, intExpr.Value)
	}
}
