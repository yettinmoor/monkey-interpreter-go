package lexer

import (
	"monkey/token"
	"testing"
)

func TestTokenize(t *testing.T) {
	input := `let x=   5;
	let y = 10;
	y = y - 4;
	y++;
	let add = fn(x, 愛) {
		return x+愛;
	};`
	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.Let, "let"},
		{token.Ident, "x"},
		{token.Assign, "="},
		{token.Int, "5"},
		{token.Semicolon, ";"},
		{token.Let, "let"},
		{token.Ident, "y"},
		{token.Assign, "="},
		{token.Int, "10"},
		{token.Semicolon, ";"},
		{token.Ident, "y"},
		{token.Assign, "="},
		{token.Ident, "y"},
		{token.Minus, "-"},
		{token.Int, "4"},
		{token.Semicolon, ";"},
		{token.Ident, "y"},
		{token.Increment, "++"},
		{token.Semicolon, ";"},
		{token.Let, "let"},
		{token.Ident, "add"},
		{token.Assign, "="},
		{token.Function, "fn"},
		{token.LParen, "("},
		{token.Ident, "x"},
		{token.Comma, ","},
		{token.Ident, "愛"},
		{token.RParen, ")"},
		{token.LBrace, "{"},
		{token.Return, "return"},
		{token.Ident, "x"},
		{token.Plus, "+"},
		{token.Ident, "愛"},
		{token.Semicolon, ";"},
		{token.RBrace, "}"},
		{token.Semicolon, ";"},
		{token.EOF, ""},
	}

	ch := make(chan *token.Token)
	l := New(input, ch)
	go l.Parse()

	for i, tt := range tests {
		tok := <-ch
		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d]: expected type %q, got type %q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d]: expected literal %q, got literal %q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestError(t *testing.T) {
	input := `let x = 3;
	let y = "hello";`
	tests := []struct {
		literal string
		row     int
		col     int
	}{
		{"let", 1, 1},
		{"x", 1, 5},
		{"=", 1, 7},
		{"3", 1, 9},
		{";", 1, 10},
		{"let", 2, 1},
		{"y", 2, 5},
		{"=", 2, 7},
		{"\"", 2, 9},
		{"hello", 2, 10},
		{"\"", 2, 15},
		{";", 2, 16},
	}

	ch := make(chan *token.Token)
	l := New(input, ch)
	go l.Parse()

	for i, tt := range tests {
		t.Logf("Test %d", i)
		tok := <-ch
		if tok.Literal != tt.literal {
			t.Errorf("Expected literal %s, got %s", tt.literal, tok.Literal)
		}
		if tok.Row != tt.row {
			t.Errorf("Expected row # %d, got %d", tt.row, tok.Row)
		}
		if tok.Col != tt.col {
			t.Errorf("Expected col # %d, got %d", tt.col, tok.Col)
		}
	}
}
