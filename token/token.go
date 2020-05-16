package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	Illegal   = "ILLEGAL"
	EOF       = "EOF"
	Ident     = "IDENT"
	Int       = "INT"
	Assign    = "="
	Plus      = "+"
	Minus     = "-"
	Comma     = ","
	Semicolon = ";"
	LParen    = "("
	RParen    = ")"
	LBrace    = "{"
	RBrace    = "}"
	Increment = "++"
	Decrement = "--"
	Function  = "FUNCTION"
	Let       = "LET"
	Return    = "RETURN"
)

var SymToks = map[string]TokenType{
	"=":  Assign,
	"+":  Plus,
	"-":  Minus,
	",":  Comma,
	";":  Semicolon,
	"(":  LParen,
	")":  RParen,
	"{":  LBrace,
	"}":  RBrace,
	"++": Increment,
	"--": Decrement,
}

var Keywords = map[string]TokenType{
	"let":    Let,
	"fn":     Function,
	"return": Return,
}
