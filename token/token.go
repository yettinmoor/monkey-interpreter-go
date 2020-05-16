package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	Assign    = "="
	Bang      = "!"
	Comma     = ","
	DQuote    = "\""
	Decrement = "--"
	EOF       = "EOF"
	String    = "STRING"
	Else      = "else"
	Eq        = "=="
	False     = "false"
	Function  = "fn"
	Ge        = ">="
	Gt        = ">"
	Ident     = "IDENT"
	If        = "if"
	Illegal   = "ILLEGAL"
	Increment = "++"
	Int       = "INT"
	LBrace    = "{"
	LParen    = "("
	Le        = "<="
	Let       = "let"
	Lt        = "<"
	Minus     = "-"
	Neq       = "!="
	Plus      = "+"
	RBrace    = "}"
	RParen    = ")"
	Return    = "return"
	SQuote    = "'"
	Semicolon = ";"
	Slash     = "/"
	Star      = "*"
	True      = "true"
)

var SymToks = map[string]TokenType{
	"!":  Bang,
	"!=": Neq,
	"'":  SQuote,
	"(":  LParen,
	")":  RParen,
	"*":  Star,
	"+":  Plus,
	"++": Increment,
	",":  Comma,
	"-":  Minus,
	"--": Decrement,
	"/":  Slash,
	";":  Semicolon,
	"<":  Lt,
	"<=": Le,
	"=":  Assign,
	"==": Eq,
	">":  Gt,
	">=": Ge,
	"\"": DQuote,
	"{":  LBrace,
	"}":  RBrace,
}

var Keywords = map[string]TokenType{
	"else":   Else,
	"false":  False,
	"fn":     Function,
	"if":     If,
	"let":    Let,
	"return": Return,
	"true":   True,
}
