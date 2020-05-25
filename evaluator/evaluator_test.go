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
