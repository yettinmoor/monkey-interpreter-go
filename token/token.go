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
	Function  = "FUNCTION"
	Let       = "LET"
	Return    = "RETURN"
)

var SingleChToks = map[rune]TokenType{
	'=': Assign,
	'+': Plus,
	'-': Minus,
	',': Comma,
	';': Semicolon,
	'(': LParen,
	')': RParen,
	'{': LBrace,
	'}': RBrace,
}

var Keywords = map[string]TokenType{
	"let":    Let,
	"fn":     Function,
	"return": Return,
}
