package lexer

import (
	"unicode"
)

// Lexer tokenizes H-lang source code
type Lexer struct {
	input   string
	pos     int  // current position in input
	readPos int  // next reading position
	ch      byte // current character
	line    int
	column  int
}

// New creates a new Lexer
func New(input string) *Lexer {
	l := &Lexer{
		input:  input,
		line:   1,
		column: 0,
	}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPos >= len(l.input) {
		l.ch = 0 // EOF
	} else {
		l.ch = l.input[l.readPos]
	}
	l.pos = l.readPos
	l.readPos++
	l.column++

	if l.ch == '\n' {
		l.line++
		l.column = 0
	}
}

func (l *Lexer) peekChar() byte {
	if l.readPos >= len(l.input) {
		return 0
	}
	return l.input[l.readPos]
}

// NextToken returns the next token
func (l *Lexer) NextToken() Token {
	var tok Token

	l.skipWhitespace()

	tok.Line = l.line
	tok.Column = l.column

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			l.readChar()
			tok = Token{Type: EQ, Literal: "==", Line: tok.Line, Column: tok.Column}
		} else if l.peekChar() == '>' {
			l.readChar()
			tok = Token{Type: ARROW, Literal: "=>", Line: tok.Line, Column: tok.Column}
		} else {
			tok = l.newToken(ASSIGN, l.ch)
		}
	case '+':
		if l.peekChar() == '+' {
			l.readChar()
			tok = Token{Type: INCREMENT, Literal: "++", Line: tok.Line, Column: tok.Column}
		} else if l.peekChar() == '=' {
			l.readChar()
			tok = Token{Type: PLUS_ASSIGN, Literal: "+=", Line: tok.Line, Column: tok.Column}
		} else {
			tok = l.newToken(PLUS, l.ch)
		}
	case '-':
		if l.peekChar() == '-' {
			l.readChar()
			tok = Token{Type: DECREMENT, Literal: "--", Line: tok.Line, Column: tok.Column}
		} else if l.peekChar() == '=' {
			l.readChar()
			tok = Token{Type: MINUS_ASSIGN, Literal: "-=", Line: tok.Line, Column: tok.Column}
		} else {
			tok = l.newToken(MINUS, l.ch)
		}
	case '*':
		if l.peekChar() == '=' {
			l.readChar()
			tok = Token{Type: MUL_ASSIGN, Literal: "*=", Line: tok.Line, Column: tok.Column}
		} else {
			tok = l.newToken(ASTERISK, l.ch)
		}
	case '/':
		if l.peekChar() == '/' {
			// Single-line comment
			tok.Type = COMMENT
			tok.Literal = l.readLineComment()
			tok.Line = tok.Line
			tok.Column = tok.Column
			return tok
		} else if l.peekChar() == '*' {
			// Multi-line comment
			tok.Type = COMMENT
			tok.Literal = l.readBlockComment()
			tok.Line = tok.Line
			tok.Column = tok.Column
			return tok
		} else if l.peekChar() == '=' {
			l.readChar()
			tok = Token{Type: DIV_ASSIGN, Literal: "/=", Line: tok.Line, Column: tok.Column}
		} else {
			tok = l.newToken(SLASH, l.ch)
		}
	case '#':
		// Shell-style comment
		tok.Type = COMMENT
		tok.Literal = l.readLineComment()
		return tok
	case '%':
		tok = l.newToken(PERCENT, l.ch)
	case '!':
		if l.peekChar() == '=' {
			l.readChar()
			tok = Token{Type: NEQ, Literal: "!=", Line: tok.Line, Column: tok.Column}
		} else {
			tok = l.newToken(BANG, l.ch)
		}
	case '&':
		if l.peekChar() == '&' {
			l.readChar()
			tok = Token{Type: AND, Literal: "&&", Line: tok.Line, Column: tok.Column}
		} else {
			tok = l.newToken(AMPERSAND, l.ch)
		}
	case '|':
		if l.peekChar() == '|' {
			l.readChar()
			tok = Token{Type: OR, Literal: "||", Line: tok.Line, Column: tok.Column}
		} else {
			tok = l.newToken(ILLEGAL, l.ch)
		}
	case '<':
		if l.peekChar() == '=' {
			l.readChar()
			tok = Token{Type: LTE, Literal: "<=", Line: tok.Line, Column: tok.Column}
		} else {
			tok = l.newToken(LT, l.ch)
		}
	case '>':
		if l.peekChar() == '=' {
			l.readChar()
			tok = Token{Type: GTE, Literal: ">=", Line: tok.Line, Column: tok.Column}
		} else {
			tok = l.newToken(GT, l.ch)
		}
	case ':':
		if l.peekChar() == '=' {
			l.readChar()
			tok = Token{Type: WALRUS, Literal: ":=", Line: tok.Line, Column: tok.Column}
		} else {
			tok = l.newToken(COLON, l.ch)
		}
	case ',':
		tok = l.newToken(COMMA, l.ch)
	case ';':
		tok = l.newToken(SEMICOLON, l.ch)
	case '.':
		tok = l.newToken(DOT, l.ch)
	case '(':
		tok = l.newToken(LPAREN, l.ch)
	case ')':
		tok = l.newToken(RPAREN, l.ch)
	case '{':
		tok = l.newToken(LBRACE, l.ch)
	case '}':
		tok = l.newToken(RBRACE, l.ch)
	case '[':
		tok = l.newToken(LBRACKET, l.ch)
	case ']':
		tok = l.newToken(RBRACKET, l.ch)
	case '"':
		tok.Type = STRING
		tok.Literal = l.readString()
		return tok
	case '\'':
		tok.Type = CHAR
		tok.Literal = l.readChar2()
		return tok
	case 0:
		tok.Literal = ""
		tok.Type = EOF
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = LookupIdent(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			tok.Literal, tok.Type = l.readNumber()
			return tok
		} else {
			tok = l.newToken(ILLEGAL, l.ch)
		}
	}

	l.readChar()
	return tok
}

func (l *Lexer) newToken(tokenType TokenType, ch byte) Token {
	return Token{
		Type:    tokenType,
		Literal: string(ch),
		Line:    l.line,
		Column:  l.column,
	}
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func (l *Lexer) readIdentifier() string {
	pos := l.pos
	for isLetter(l.ch) || isDigit(l.ch) || l.ch == '_' {
		l.readChar()
	}
	return l.input[pos:l.pos]
}

func (l *Lexer) readNumber() (string, TokenType) {
	pos := l.pos
	tokenType := INT

	for isDigit(l.ch) {
		l.readChar()
	}

	// Check for float
	if l.ch == '.' && isDigit(l.peekChar()) {
		tokenType = FLOAT
		l.readChar() // consume '.'
		for isDigit(l.ch) {
			l.readChar()
		}
	}

	return l.input[pos:l.pos], tokenType
}

func (l *Lexer) readString() string {
	l.readChar() // skip opening "
	pos := l.pos

	for l.ch != '"' && l.ch != 0 {
		if l.ch == '\\' {
			l.readChar() // skip escape char
		}
		l.readChar()
	}

	str := l.input[pos:l.pos]
	l.readChar() // skip closing "
	return str
}

func (l *Lexer) readChar2() string {
	l.readChar() // skip opening '
	pos := l.pos

	if l.ch == '\\' {
		l.readChar() // escape sequence
	}
	l.readChar()

	str := l.input[pos:l.pos]
	l.readChar() // skip closing '
	return str
}

func (l *Lexer) readLineComment() string {
	// Skip // or #
	if l.ch == '/' {
		l.readChar()
	}
	l.readChar()

	pos := l.pos
	for l.ch != '\n' && l.ch != 0 {
		l.readChar()
	}
	return l.input[pos:l.pos]
}

func (l *Lexer) readBlockComment() string {
	// Skip /*
	l.readChar()
	l.readChar()

	pos := l.pos
	for {
		if l.ch == 0 {
			break
		}
		if l.ch == '*' && l.peekChar() == '/' {
			end := l.pos
			l.readChar() // skip *
			l.readChar() // skip /
			return l.input[pos:end]
		}
		l.readChar()
	}
	return l.input[pos:l.pos]
}

func isLetter(ch byte) bool {
	return unicode.IsLetter(rune(ch)) || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}
