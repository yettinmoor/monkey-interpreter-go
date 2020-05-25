package token

type TokenType uint8

func (t *TokenType) String() string {
	return allTokens[*t]
}

type Token struct {
	Type     TokenType
	Literal  string
	Row, Col int
}

type tokenGroup map[string]TokenType

const (
	_ TokenType = iota
	And
	Assign
	Bang
	Comma
	DQuote
	Decrement
	EOF
	Else
	Eq
	False
	Function
	Ge
	Gt
	Hash
	Ident
	If
	Illegal
	Increment
	Int
	LBrace
	LParen
	Le
	Let
	Lt
	Minus
	Modulo
	Neq
	Or
	Plus
	RBrace
	RParen
	Return
	SQuote
	Semicolon
	Slash
	Star
	String
	True
)

var allTokens = func() map[TokenType]string {
	allTokens := make(map[TokenType]string, len(SymToks)+len(Keywords)+len(special))
	for _, tokGroup := range []tokenGroup{SymToks, Keywords, special} {
		for name, tokID := range tokGroup {
			allTokens[tokID] = name
		}
	}
	return allTokens
}()

var SymToks = tokenGroup{
	"!":  Bang,
	"!=": Neq,
	"#":  Hash,
	"%":  Modulo,
	"&&": And,
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
	"||": Or,
	"}":  RBrace,
}

var Keywords = tokenGroup{
	"else":   Else,
	"false":  False,
	"fn":     Function,
	"if":     If,
	"let":    Let,
	"return": Return,
	"true":   True,
}

var special = tokenGroup{
	"EOF":     EOF,
	"INT":     Int,
	"IDENT":   Ident,
	"ILLEGAL": Illegal,
}
