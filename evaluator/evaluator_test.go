package evaluator

import (
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
	"monkey/token"
	"testing"
)

func TestEvalIntExpr(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
		{"-1", -1},
		{"2 - 1", 1},
		{"10 * 5 + 15", 65},
	}
	for _, tt := range tests {
		ch := make(chan *token.Token)
		l := lexer.New(tt.input, ch)
		program := parser.New(l, ch).Parse()
		output := Eval(program, object.NewEnv(nil))
		if output, ok := output.(*object.ObjInt); !ok {
			t.Errorf("Expected int64 value, got %T", output)
		} else if output.Value != tt.expected {
			t.Errorf("Expected value %v, got %v", tt.expected, output.Value)
		}
	}
}

func TestFuncCall(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"let add = fn(x, y) { return x+y; }; add(3, 4)", "7"},
		{"fn(x, y){ return x-y; }(2, 1)", "1"},
	}
	for _, tt := range tests {
		ch := make(chan *token.Token)
		l := lexer.New(tt.input, ch)
		program := parser.New(l, ch).Parse()
		output := Eval(program, object.NewEnv(nil)).String()
		if output != tt.expected {
			t.Errorf("Expected output %q, got %q", tt.expected, output)
		}
	}
}
