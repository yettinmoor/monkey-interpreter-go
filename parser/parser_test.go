package parser

import (
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
		{2, 11},
		{3, 9},
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
	let y = 10;
	let foo = 5;`

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
	input := `
	return 10;
	return 5;
	return 5+6;`

	program := setup(t, input)

	if ls := len(program.Stmts); ls != 3 {
		t.Fatalf("Expected 3 return stmts, got %d", ls)
	}

	tests := []struct {
		value string
	}{
		{"10"},
		{"5"},
		{"(5+6)"},
	}

	for i, tt := range tests {
		if retStmt, ok := program.Stmts[i].(*ast.ReturnStmt); !ok {
			t.Fatalf("Expected return stmt, got %T", program.Stmts[i])
		} else if retStmt.Value.String() != tt.value {
			t.Fatalf("Expected returned value %s, got %s", tt.value, retStmt.Value.String())
		}
	}
}

func TestIdentExpr(t *testing.T) {
	input := `foobar;`

	program := setup(t, input)

	if len(program.Stmts) != 1 {
		t.Fatalf("Expected 1 expr, got %d", len(program.Stmts))
	}
	if stmt, ok := program.Stmts[0].(*ast.ExprStmt); !ok {
		t.Fatalf("Not exprstmt, got %T", program.Stmts[0])
	} else if ident, ok := stmt.Expr.(*ast.IdentExpr); !ok {
		t.Fatalf("not ident expr, got %T", stmt.Expr)
	} else if ident.Token.Literal != "foobar" {
		t.Errorf("Token literal not foobar, got %s", ident.Token.Literal)
	}
}

func TestIntLiteralExpr(t *testing.T) {
	input := `5;`

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

func TestPrefixExpr(t *testing.T) {
	tests := []struct {
		input, op string
		value     int64
	}{
		{"!5;", "!", 5},
		{"-15;", "-", 15},
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

func testIntLit(t *testing.T, expr ast.Expr, value int64) {
	if intExpr, ok := expr.(*ast.IntLiteralExpr); !ok {
		t.Fatalf("Not int expr, got %T", expr)
	} else if intExpr.Value != value {
		t.Errorf("Expected value %d, got %d", value, intExpr.Value)
	}
}
