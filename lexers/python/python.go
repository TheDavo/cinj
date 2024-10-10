package python

import (
	"errors"
	"fmt"
	"strings"

	lex "github.com/TheDavo/cinj/lexers"
)

const (
	LPAREN     = "LPAREN"
	RPAREN     = "RPAREN"
	CLASS      = "CLASS"
	TAB        = "TAB"
	EOF        = "EOF"
	IDENT      = "IDENT"
	FUNCTION   = "FUNCTION"
	COLON      = "COLON"
	NEWLINE    = "NEWLINE"
	STARTBLOCK = "STARTBLOCK"
	ENDBLOCK   = "ENDBLOCK"
	IGNORE     = "IGNORE"
	DECORATOR  = "@"
	IMPORT     = "IMPORT"
	FROM       = "FROM"
)

var keywords = map[string]lex.TokenType{
	"class":  CLASS,
	"def":    FUNCTION,
	"@":      DECORATOR,
	"import": IMPORT,
	"from":   FROM,
}

func keywordFromTokenType(tt lex.TokenType) string {
	for k, v := range keywords {
		if v == tt {
			return k
		}
	}
	return ""
}

type PyBlock struct {
	Block lex.TokenType // The STARTBLOCK and ENDBLOCK tokens
	Line  int
	Depth int
}

type PythonLexer struct {
	input          string
	position       int
	readPosition   int
	line           int
	column         int
	ch             byte
	depth          int
	lastNewLinePos int
	indentSize     int
	blockStack     []PyBlock
	tokens         []lex.Token
	lines          []lex.Line
	tokenTree      []lex.TokenTree
}

func NewLexer(input string, indentSize int) *PythonLexer {
	stackInit := []PyBlock{{STARTBLOCK, 1, 1}}
	return &PythonLexer{
		position:     0,
		readPosition: 0,
		line:         1,
		input:        input,
		column:       1,
		blockStack:   stackInit,
		depth:        1,
		indentSize:   indentSize,
		tokenTree: []lex.TokenTree{
			{
				Node: lex.Token{
					Type:    lex.ROOT,
					Literal: lex.ROOT,
				},
				Children: []lex.Token{},
			},
		},
	}
}

func (pl *PythonLexer) Lex() {
	pl.readChar()
	for !pl.isAtEnd() {
		pl.nextToken()
	}
}

func (pl *PythonLexer) nextToken() lex.Token {
	var tok lex.Token

	currDepth := pl.depth

	// Start of a new line
	if pl.column == 1 {
		pl.lines = append(pl.lines, lex.Line{
			StartPosition: pl.position,
			EndPosition:   0,
		})
		pl.skipIndentation()
		if pl.depth > currDepth {
			pl.blockStack = append(pl.blockStack,
				PyBlock{STARTBLOCK, pl.line, pl.depth})
		} else if pl.depth < currDepth || pl.isAtEnd() {
			pl.blockStack = append(pl.blockStack,
				PyBlock{ENDBLOCK, pl.line, currDepth})
		}
		// skipIndentation does not do anything since a new line is starting
		// from the beginning, so check for that, and reset depth
		if pl.column == 1 {
			pl.depth = 1
			pl.blockStack = append(pl.blockStack,
				PyBlock{ENDBLOCK, pl.line, 2})
		}
	}
	tok.Line = pl.line
	tok.Column = pl.column
	tok.StartPosition = pl.position
	tok.Depth = pl.depth

	// Always skip whitespace, even after indentation
	pl.skipWhitespace()

	switch pl.ch {
	case '\n':
		pl.line++
		pl.column = 1
		pl.lastNewLinePos = pl.readPosition
		pl.lines[len(pl.lines)-1].EndPosition = pl.lastNewLinePos
		tok.Type = NEWLINE
		tok.Literal = "\\n"
	case ':':
		tok.Type = COLON
		tok.Literal = ":"
	case '(':
		tok.Type = LPAREN
		tok.Literal = "("
	case ')':
		tok.Type = RPAREN
		tok.Literal = ")"
	case 0:
		tok.Type = EOF
		tok.Literal = ""
		for pl.depth > 1 {
			pl.depth--
			pl.blockStack = append(pl.blockStack,
				PyBlock{ENDBLOCK, pl.line, pl.depth})
		}
		tok.Depth = 1
		tok.EndPosition = pl.readPosition
		pl.lines[len(pl.lines)-1].EndPosition = pl.lastNewLinePos
		pl.tokens = append(pl.tokens, tok)
		return tok
	default:
		if isLetter(pl.ch) {
			tok.Literal = pl.getIdentifier()
			tok.Type = pl.MatchKeyword(tok)
			tok.Depth = pl.depth
			// tok.Column = l.column
			tok.EndPosition = pl.readPosition
			pl.tokens = append(pl.tokens, tok)
			return tok
		} else {
			tok.Literal = IGNORE
			tok.Type = IGNORE
		}

	}
	pl.readChar()
	if pl.isAtEnd() {
		for pl.depth > 1 {
			pl.depth--
			pl.blockStack = append(pl.blockStack,
				PyBlock{ENDBLOCK, pl.line, pl.depth})
		}
		tok.Type = EOF
		tok.Literal = ""
		tok.Depth = pl.depth
		tok.Column = 1
		tok.EndPosition = pl.readPosition
		pl.tokens = append(pl.tokens, tok)
		return tok
	}
	tok.EndPosition = pl.readPosition
	pl.tokens = append(pl.tokens, tok)
	return tok
}

func (pl *PythonLexer) skipIndentation() {
	spaceCount := 1
	currDepth := pl.depth
	pl.depth = 1

	if pl.ch == ' ' {
		for pl.ch == ' ' {
			if spaceCount%pl.indentSize == 0 {
				pl.depth++
				pl.column++
				spaceCount = 0
			}
			spaceCount++
			pl.readChar()
		}
	} else if pl.ch == '\t' {
		for pl.ch == '\t' {
			pl.depth++
			pl.column++
			pl.readChar()
		}
	} else if pl.ch == '\n' {
		pl.depth = currDepth
	}
}

func (pl PythonLexer) MatchKeyword(t lex.Token) lex.TokenType {
	if tType, ok := keywords[t.Literal]; ok {
		return tType
	}

	return IDENT
}

func (pl PythonLexer) isAtEnd() bool {
	return pl.readPosition >= len(pl.input)
}

func (pl *PythonLexer) readChar() {
	if pl.isAtEnd() {
		pl.ch = 0
	} else {
		pl.ch = pl.input[pl.readPosition]
	}

	pl.position = pl.readPosition
	pl.readPosition += 1

	// columns are typically 1 indexed, so add that
	pl.column = pl.position - pl.lastNewLinePos + 1
}

func (pl *PythonLexer) peekChar() byte {
	if pl.isAtEnd() {
		return 0
	} else {
		return pl.input[pl.readPosition]
	}
}

func (l *PythonLexer) getIdentifier() string {
	start := l.position
	for isLetter(l.ch) || isDigit(l.ch) {
		l.readChar()
	}

	return l.input[start:l.position]
}

func (pl *PythonLexer) skipWhitespace() {
	for pl.ch == ' ' || pl.ch == '\r' {
		pl.readChar()
	}
}

func (pl *PythonLexer) skipWhitespaceWithNewline() {
	for pl.ch == ' ' || pl.ch == '\r' || pl.ch == '\n' {
		pl.readChar()
	}
}

// findToken finds the first instance of the token with identifer `ident`
// and of a particular TokenType `tt`
func (pl PythonLexer) findToken(tt lex.TokenType, ident string) (lex.Token,
	int, error,
) {
	var emptyTok lex.Token
	for i, token := range pl.tokens {
		if token.Type == tt && token.Literal == ident {
			return token, i, nil
		}
	}
	return emptyTok, 0, fmt.Errorf("Could not find %s", ident)
}

func (pl PythonLexer) findTokenParent(t lex.Token) (lex.Token, error) {
	if t.Depth == 1 {
		return lex.Token{
			Type:    lex.ROOT,
			Literal: lex.ROOT,
		}, nil
	}

	_, idx, err := pl.findToken(t.Type, t.Literal)
	if err != nil {
		return lex.Token{}, err
	}

	for i := idx; i > 0; i-- {
		if pl.tokens[i].Depth == t.Depth-1 {
			return pl.tokens[i], nil
		}
	}

	return lex.Token{}, errors.New("Cannot find parent token")
}

func (pl PythonLexer) findTokenParentFromIndex(idx int) (lex.Token, error) {
	if pl.tokens[idx].Depth == 1 {
		return lex.Token{
			Type:    lex.ROOT,
			Literal: lex.ROOT,
		}, nil
	}

	for i := idx; i > 0; i-- {
		// fmt.Println(pl.tokens[i])
		if pl.tokens[i].Depth == pl.tokens[idx].Depth-1 && pl.tokens[i].Type != NEWLINE {
			return pl.tokens[i], nil
		}
	}

	return lex.Token{}, errors.New("Cannot find parent token")
}

// findTokens returns a slice of integers with the line location for a token
// with identifer `ident` and of TokenType `tt`
func (pl PythonLexer) findTokens(tt lex.TokenType, ident string) ([]lex.Token,
	[]int, error,
) {
	foundTokens := []lex.Token{}
	tokenLocs := []int{}
	for i, token := range pl.tokens {
		if token.Type == tt && token.Literal == ident {
			foundTokens = append(foundTokens, token)
			tokenLocs = append(tokenLocs, i)
		}
	}

	if len(foundTokens) == 0 {
		return foundTokens, []int{}, errors.New("No tokens found")
	}

	return foundTokens, tokenLocs, nil
}

func (pl PythonLexer) findKeywordThenIdentifierLine(tt lex.TokenType, ident string) (int, error) {
	_, locs, err := pl.findTokens(tt, keywordFromTokenType(tt))
	if err != nil {
		return 0, err
	}

	for _, i := range locs {
		if i != len(pl.tokens) {
			if pl.tokens[i+1].Literal == ident {
				return pl.tokens[i].Line, nil
			}
		}
	}

	return 0, nil
}

// findKeywordThenIdentifierLines is used to find specific two-token pairs
// such as class ClassName or def function_name
func (pl PythonLexer) findKeywordThenIdentifierLines(
	tt lex.TokenType,
	ident string,
) ([]int, error) {
	_, locs, err := pl.findTokens(tt, keywordFromTokenType(tt))
	lines := []int{}
	if err != nil {
		return lines, err
	}

	for _, i := range locs {
		if i != len(pl.tokens) {
			if pl.tokens[i+1].Literal == ident {
				lines = append(lines, pl.tokens[i].Line)
			}
		}
	}

	return lines, nil
}

// findBlockRange returns the starting line and ending line of a block that
// contains the token being searched for.
// This function is used to find a block in Python such as a class block
// or function block
func (pl PythonLexer) findBlockRange(tt lex.TokenType, ident string) (int,
	int, error,
) {
	tok, idx, err := pl.findToken(tt, ident)
	if err != nil {
		return 0, 0, err
	}

	remainingTokens := pl.tokens[idx+1:]

	for _, token := range remainingTokens {
		if token.Type != NEWLINE && token.Depth <= tok.Depth &&
			tok.Line != token.Line {
			return tok.Line, token.Line, nil
		}
	}

	return 0, 0, errors.New("Cannot find block range for the token")
}

func (pl PythonLexer) findBlockRangePos(tt lex.TokenType,
	ident string,
) (int, int, error) {
	searchedToken, idx, err := pl.findToken(tt, ident)
	if err != nil {
		return 0, 0, err
	}

	remainingTokens := pl.tokens[idx+1:]

	for i, token := range remainingTokens {
		if token.Type != NEWLINE && token.Depth <= searchedToken.Depth &&
			searchedToken.Line != token.Line {
			// Bounds check
			if idx+i < len(pl.tokens)-2 {
				return searchedToken.StartPosition,
					remainingTokens[i].StartPosition, nil
			}
			if token.EndPosition >= len(pl.input) {
				return searchedToken.StartPosition, token.StartPosition, nil
			}
			return searchedToken.StartPosition, token.EndPosition, nil
		}
	}

	return 0, 0, errors.New("Cannot find block range for the token")
}

func (pl PythonLexer) findBlockRangePosFromToken(t lex.Token, idx int) (int, int, error) {
	if idx == len(pl.tokens)-1 {
		return 0, 0, errors.New("index parameter at length of tokens slice")
	}
	remainingTokens := pl.tokens[idx+1:]

	for i, token := range remainingTokens {
		if token.Type != NEWLINE && token.Depth <= t.Depth &&
			t.Line != token.Line {
			// Bounds check
			if idx+i < len(pl.tokens)-2 {
				return t.StartPosition,
					remainingTokens[i].StartPosition, nil
			}
			if token.EndPosition >= len(pl.input) {
				return t.StartPosition, token.StartPosition, nil
			}
			return t.StartPosition, token.EndPosition, nil
		}
	}

	return 0, 0, errors.New("Cannot find block range for the token")
}

// GetClass returns a string corresponding to a class block in the text input
// of the lexer
func (pl *PythonLexer) GetClass(className string) (string, error) {
	_, foundTokenLocs, err := pl.findTokens(IDENT, className)
	if err != nil {
		return "", err
	}

	for _, idx := range foundTokenLocs {
		if pl.tokens[idx-1].Type == CLASS {
			decorators := pl.findDecoratorsAboveToken(pl.tokens[idx-1])
			_, end, err := pl.findBlockRangePos(IDENT, className)
			if err != nil {
				return "", err
			}
			return fmt.Sprintf(
				"%s%s",
				decorators,
				pl.input[pl.tokens[idx-1].StartPosition:end],
			), nil
		}
	}

	return "", errors.New("Could not find class")
}

// GetFunction returns a string corresponding to a function
// block in the text input of the lexer
func (pl *PythonLexer) GetFunction(functionName string,
	className string,
) (string, error) {
	_, foundTokenLocs, err := pl.findTokens(IDENT, functionName)
	if err != nil {
		return "", err
	}

	keyClassLines, err := pl.findKeywordThenIdentifierLines(CLASS, className)
	if err != nil {
		return "", err
	}

	for _, idx := range foundTokenLocs {
		// checks that the preceding token is a function
		// to validate the current token is not a function call
		// but a definition
		if pl.tokens[idx-1].Type == FUNCTION {
			_, end, err := pl.findBlockRangePosFromToken(pl.tokens[idx-1], idx-1)
			if err != nil {
				return "", err
			}

			decorators := pl.findDecoratorsAboveToken(pl.tokens[idx-1])

			if className == "" {
				return fmt.Sprintf(
					"%s%s",
					decorators,
					pl.input[pl.tokens[idx-2].StartPosition:end],
				), nil
				// return pl.input[pl.tokens[idx-2].StartPosition:end], nil
			}

			// Add flavor text showing which class the function is in
			for _, classLine := range keyClassLines {
				if classLine < pl.tokens[idx].Line {
					parentLine, err := pl.getLine(
						classLine - 1,
					) // offset since the line parameter is 1-indexed
					if err != nil {
						return "", err
					}
					function := fmt.Sprintf("%s#----%s%s",
						// idx-2 here to grab the correct spacing
						// idx-1 does not work as the tokenizer skips indentation
						// and the start position of idx-1 is after
						// the indentation
						parentLine, decorators,
						pl.input[pl.tokens[idx-2].StartPosition:end])
					return function, nil

				}
			}
		}
	}
	return "", errors.New("Could not find class")
}

// findDecoratorsAboveToken is a helper function to find decorators that
// are used in functions and classes in Python
// This function is needed as the currently there is no tree style setup
// that links the decorators to a function or token
func (pl PythonLexer) findDecoratorsAboveToken(t lex.Token) string {
	tLine := t.Line
	decoratorOffset := 1
	line, err := pl.getLine(tLine - decoratorOffset - 1)
	decoratorStr := ""

	if err != nil || t.Line == 1 {
		return ""
	}

	for len(line) >= t.Column && (line != NEWLINE && line != "") {
		if line[t.Column-1] == '@' {
			decoratorOffset++
		} else {
			decoratorOffset--
			break
		}
		line, err = pl.getLine(tLine - decoratorOffset - 1)
		if err != nil {
			return ""
		}
	}

	for i := tLine - decoratorOffset - 1; i < tLine-1; i++ {
		line, err = pl.getLine(i)
		if err != nil {
			return ""
		}
		decoratorStr += line

	}
	decoratorStr = strings.TrimSuffix(decoratorStr, "\r\n")
	return decoratorStr
}

// getLine returns a string corresponding to the line of the text input of
// the lexer
func (pl PythonLexer) getLine(line int) (string, error) {
	if line < 0 || line >= len(pl.lines) {
		return "", errors.New("line value out of bounds")
	}
	start := pl.lines[line].StartPosition
	end := pl.lines[line].EndPosition

	return pl.input[start:end], nil
}

func isLetter(ch byte) bool {
	return ch == '_' || (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}

func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}
