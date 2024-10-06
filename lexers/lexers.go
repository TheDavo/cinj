package lexers

type Lexer interface {
	Lex()
	NewLexer() *Lexer
	MatchKeyword(Token) TokenType
}

type TokenType string

type Token struct {
	Type          TokenType
	Literal       string
	Line          int
	Column        int
	Depth         int
	StartPosition int
	EndPosition   int
}

// Use to make the root node of the TokenTree struct
const ROOT = "ROOT"

type TokenTree struct {
	Node     Token
	Children []Token
}

type Line struct {
	StartPosition int
	EndPosition   int
}
