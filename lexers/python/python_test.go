package python

import (
	"os"
	"testing"

	lex "github.com/TheDavo/cinj/lexers"
)

func TestNextToken(t *testing.T) {
	input := `class Teehee:
	def __init__():
		pass

	def test1():
		pass
		`

	tests := []struct {
		expectedType    lex.TokenType
		expectedLiteral string
	}{
		{CLASS, "class"},
		{IDENT, "Teehee"},
		{COLON, ":"},
		{NEWLINE, "\\n"},
		{FUNCTION, "def"},
		{IDENT, "__init__"},
		{LPAREN, "("},
		{RPAREN, ")"},
		{COLON, ":"},
		{NEWLINE, "\\n"},
		{IDENT, "pass"},
		{NEWLINE, "\\n"},
		{NEWLINE, "\\n"},
		{FUNCTION, "def"},
		{IDENT, "test1"},
		{LPAREN, "("},
		{RPAREN, ")"},
		{COLON, ":"},
		{NEWLINE, "\\n"},
		{IDENT, "pass"},
		{NEWLINE, "\\n"},
		{EOF, ""},
	}

	l := NewLexer(input, 2)
	l.Lex()
	// t.Log("Size of tokens slice:", len(l.tokens))
	// t.Log("Size of tests slice:", len(tests))
	for i, tt := range tests {

		tok := l.tokens[i]
		// t.Log(tok.Literal, tok.Type, tok.Column, tok.Depth, tok.Line)
		if lex.TokenType(tok.Literal) != lex.TokenType(tt.expectedLiteral) {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q on line %d on column %d",
				i, tt.expectedLiteral, tok.Literal, tok.Line, tok.Column)
		}

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q on line %d",
				i, tt.expectedType, tok.Type, tok.Line)
		}

	}
}

func TestFindToken(t *testing.T) {
	input := `class Teehee:
	def __init__():
		var_value = 5

	def test1():
		pass
`
	l := NewLexer(input, 2)
	l.Lex()

	lineTest1 := 5

	// t.Log("Size of l.tokens,", len(l.tokens))

	// for _, token := range l.tokens {
	// 	t.Log(token.Literal, token.Type, token.Line, token.Column, token.Depth)
	// }
	tok, _, err := l.findToken(IDENT, "test1")
	if err != nil {
		t.Fatalf("test FindToken for token literal %s error'd\nwith error: %s",
			"test1", err.Error())
	}

	if lineTest1 != tok.Line {
		t.Fatalf("finding function 'test1' line, expected %d, got %d",
			lineTest1, tok.Line)
	}
}

func TestFindBlockRange(t *testing.T) {
	input := `class Teehee:
	def __init__():
		var_value = 5

	def test1():
		hello = "world"
		another_val = 5
		`

	l := NewLexer(input, 2)
	l.Lex()

	tests := []struct {
		tType     lex.TokenType
		tLiteral  string
		startStop []int
	}{
		{IDENT, "Teehee", []int{1, 8}},
		{IDENT, "__init__", []int{2, 5}},
		{IDENT, "hello", []int{6, 7}},
		{IDENT, "test1", []int{5, 8}},
	}
	// for _, token := range l.tokens {
	// 	t.Log(token.Literal, token.Type, token.Line, token.Column, token.Depth)
	// }

	for i, test := range tests {
		start, end, err := l.findBlockRange(test.tType, test.tLiteral)
		if err != nil {
			t.Fatal(err.Error())
		}

		expectedStart := test.startStop[0]
		expectedEnd := test.startStop[1]

		if start != expectedStart || end != expectedEnd {
			for _, b := range l.blockStack {
				t.Log(b.Block, b.Line, b.Depth)
			}
			t.Fatalf(
				"tests[%d] - finding block range for token with literal %s, expected start, end of %d, %d, got start, end of %d, %d",
				i,
				test.tLiteral,
				expectedStart,
				expectedEnd,
				start,
				end,
			)
		}
	}
}

func TestGetClass(t *testing.T) {
	input := `class Teehee:
  def __init__():
    var_value = 5

  def test1():
    hello = "world"
    another_val = 5

class Test2:
  def __init__():
    gotta_have_one = True
    `

	l := NewLexer(input, 2)
	l.Lex()

	class, err := l.GetClass("Teehee")
	expected := `class Teehee:
  def __init__():
    var_value = 5

  def test1():
    hello = "world"
    another_val = 5

`
	if err != nil {
		t.Fatal(err.Error())
	}

	if class != expected {
		t.Fatalf("Expected \n%s\nGot \n%s", expected, class)
	}
}

func TestFindBlockRangeFromFile(t *testing.T) {
	fileName := "./pySample.py"
	input, err := os.ReadFile(fileName)
	if err != nil {
		t.Fatal(err.Error())
	}

	l := NewLexer(string(input), 2)

	l.Lex()

	tests := []struct {
		tType     lex.TokenType
		tLiteral  string
		startStop []int
	}{
		{IDENT, "Test", []int{6, 19}},
		{IDENT, "print_content", []int{10, 13}},
		{IDENT, "MyABC", []int{19, 29}},
	}
	// for _, token := range l.tokens {
	// 	t.Log(token.Literal, token.Type, token.Line, token.Column, token.Depth)
	// }

	for i, test := range tests {
		start, end, err := l.findBlockRange(test.tType, test.tLiteral)
		if err != nil {
			t.Fatal(err.Error())
		}

		expectedStart := test.startStop[0]
		expectedEnd := test.startStop[1]

		if start != expectedStart || end != expectedEnd {
			for _, b := range l.blockStack {
				t.Log(b.Block, b.Line, b.Depth)
			}
			t.Fatalf(
				"tests[%d] - finding block range for token with literal %s, expected start, end of %d, %d, got start, end of %d, %d",
				i,
				test.tLiteral,
				expectedStart,
				expectedEnd,
				start,
				end,
			)
		}
	}
}

func TestGetLine(t *testing.T) {
	input := `class Teehee:
  def __init__():
    var_value = 5

  def test1():
    hello = "world"
    another_val = 5

class Test2:
  def __init__():
    gotta_have_one = True
    `

	l := NewLexer(input, 2)
	l.Lex()

	tests := []struct {
		line     int
		expected string
	}{
		{0, "class Teehee:\n"},
		{1, "  def __init__():\n"},
	}

	for i, test := range tests {
		gotLine, err := l.getLine(test.line)
		if err != nil {
			t.Fatalf(err.Error())
		}
		if gotLine != test.expected {
			t.Fatalf("tests[%d]: expected %s, got %s", i, test.expected, gotLine)
		}
	}
}
