package lexer

import "testing"

func TestNextToken_SingleCharacters(t *testing.T) {
	input := `=+-*/%!<>,;:.(){}[]&`

	tests := []struct {
		expectedType    TokenType
		expectedLiteral string
	}{
		{ASSIGN, "="},
		{PLUS, "+"},
		{MINUS, "-"},
		{ASTERISK, "*"},
		{SLASH, "/"},
		{PERCENT, "%"},
		{BANG, "!"},
		{LT, "<"},
		{GT, ">"},
		{COMMA, ","},
		{SEMICOLON, ";"},
		{COLON, ":"},
		{DOT, "."},
		{LPAREN, "("},
		{RPAREN, ")"},
		{LBRACE, "{"},
		{RBRACE, "}"},
		{LBRACKET, "["},
		{RBRACKET, "]"},
		{AMPERSAND, "&"},
		{EOF, ""},
	}

	l := New(input)
	for i, tt := range tests {
		tok := l.NextToken()
		if tok.Type != tt.expectedType {
			t.Errorf("tests[%d] - type wrong. expected=%v, got=%v",
				i, tt.expectedType, tok.Type)
		}
		if tok.Literal != tt.expectedLiteral {
			t.Errorf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestNextToken_TwoCharOperators(t *testing.T) {
	input := `:= == != <= >= && || ++ -- += -= *= /= =>`

	tests := []struct {
		expectedType    TokenType
		expectedLiteral string
	}{
		{WALRUS, ":="},
		{EQ, "=="},
		{NEQ, "!="},
		{LTE, "<="},
		{GTE, ">="},
		{AND, "&&"},
		{OR, "||"},
		{INCREMENT, "++"},
		{DECREMENT, "--"},
		{PLUS_ASSIGN, "+="},
		{MINUS_ASSIGN, "-="},
		{MUL_ASSIGN, "*="},
		{DIV_ASSIGN, "/="},
		{ARROW, "=>"},
		{EOF, ""},
	}

	l := New(input)
	for i, tt := range tests {
		tok := l.NextToken()
		if tok.Type != tt.expectedType {
			t.Errorf("tests[%d] - type wrong. expected=%v, got=%v",
				i, tt.expectedType, tok.Type)
		}
		if tok.Literal != tt.expectedLiteral {
			t.Errorf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestNextToken_Keywords(t *testing.T) {
	input := `function struct enum if else for while return const var public null true false alloc free defer len make range break continue int float string char bool void`

	tests := []struct {
		expectedType    TokenType
		expectedLiteral string
	}{
		{FUNCTION, "function"},
		{STRUCT, "struct"},
		{ENUM, "enum"},
		{IF, "if"},
		{ELSE, "else"},
		{FOR, "for"},
		{WHILE, "while"},
		{RETURN, "return"},
		{CONST, "const"},
		{VAR, "var"},
		{PUBLIC, "public"},
		{NULL, "null"},
		{TRUE, "true"},
		{FALSE, "false"},
		{ALLOC, "alloc"},
		{FREE, "free"},
		{DEFER, "defer"},
		{LEN, "len"},
		{MAKE, "make"},
		{RANGE, "range"},
		{BREAK, "break"},
		{CONTINUE, "continue"},
		{TYPE_INT, "int"},
		{TYPE_FLOAT, "float"},
		{TYPE_STRING, "string"},
		{TYPE_CHAR, "char"},
		{TYPE_BOOL, "bool"},
		{TYPE_VOID, "void"},
		{EOF, ""},
	}

	l := New(input)
	for i, tt := range tests {
		tok := l.NextToken()
		if tok.Type != tt.expectedType {
			t.Errorf("tests[%d] - type wrong. expected=%v, got=%v",
				i, tt.expectedType, tok.Type)
		}
		if tok.Literal != tt.expectedLiteral {
			t.Errorf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestNextToken_Identifiers(t *testing.T) {
	input := `foo bar_baz myVar _private camelCase`

	tests := []struct {
		expectedType    TokenType
		expectedLiteral string
	}{
		{IDENT, "foo"},
		{IDENT, "bar_baz"},
		{IDENT, "myVar"},
		{IDENT, "_private"},
		{IDENT, "camelCase"},
		{EOF, ""},
	}

	l := New(input)
	for i, tt := range tests {
		tok := l.NextToken()
		if tok.Type != tt.expectedType {
			t.Errorf("tests[%d] - type wrong. expected=%v, got=%v",
				i, tt.expectedType, tok.Type)
		}
		if tok.Literal != tt.expectedLiteral {
			t.Errorf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestNextToken_Numbers(t *testing.T) {
	input := `42 0 123456 3.14 0.5 100.001`

	tests := []struct {
		expectedType    TokenType
		expectedLiteral string
	}{
		{INT, "42"},
		{INT, "0"},
		{INT, "123456"},
		{FLOAT, "3.14"},
		{FLOAT, "0.5"},
		{FLOAT, "100.001"},
		{EOF, ""},
	}

	l := New(input)
	for i, tt := range tests {
		tok := l.NextToken()
		if tok.Type != tt.expectedType {
			t.Errorf("tests[%d] - type wrong. expected=%v, got=%v",
				i, tt.expectedType, tok.Type)
		}
		if tok.Literal != tt.expectedLiteral {
			t.Errorf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestNextToken_Strings(t *testing.T) {
	input := `"hello" "world" "with spaces" ""`

	tests := []struct {
		expectedType    TokenType
		expectedLiteral string
	}{
		{STRING, "hello"},
		{STRING, "world"},
		{STRING, "with spaces"},
		{STRING, ""},
		{EOF, ""},
	}

	l := New(input)
	for i, tt := range tests {
		tok := l.NextToken()
		if tok.Type != tt.expectedType {
			t.Errorf("tests[%d] - type wrong. expected=%v, got=%v",
				i, tt.expectedType, tok.Type)
		}
		if tok.Literal != tt.expectedLiteral {
			t.Errorf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestNextToken_Characters(t *testing.T) {
	input := `'a' 'Z' '0'`

	tests := []struct {
		expectedType    TokenType
		expectedLiteral string
	}{
		{CHAR, "a"},
		{CHAR, "Z"},
		{CHAR, "0"},
		{EOF, ""},
	}

	l := New(input)
	for i, tt := range tests {
		tok := l.NextToken()
		if tok.Type != tt.expectedType {
			t.Errorf("tests[%d] - type wrong. expected=%v, got=%v",
				i, tt.expectedType, tok.Type)
		}
		if tok.Literal != tt.expectedLiteral {
			t.Errorf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestNextToken_Comments(t *testing.T) {
	input := `
	// C-style single line
	x := 42;
	# Shell-style comment
	y := 10;
	/* Multi
	   line */
	z := 5;
	`

	l := New(input)

	// Skip comment, find x
	tok := l.NextToken()
	for tok.Type == COMMENT {
		tok = l.NextToken()
	}
	if tok.Type != IDENT || tok.Literal != "x" {
		t.Errorf("expected IDENT 'x', got %v '%s'", tok.Type, tok.Literal)
	}

	// Skip to y
	for tok.Type != IDENT || tok.Literal != "y" {
		tok = l.NextToken()
		if tok.Type == EOF {
			t.Fatal("unexpected EOF looking for 'y'")
		}
	}
	if tok.Literal != "y" {
		t.Errorf("expected 'y', got '%s'", tok.Literal)
	}

	// Skip to z
	for tok.Type != IDENT || tok.Literal != "z" {
		tok = l.NextToken()
		if tok.Type == EOF {
			t.Fatal("unexpected EOF looking for 'z'")
		}
	}
	if tok.Literal != "z" {
		t.Errorf("expected 'z', got '%s'", tok.Literal)
	}
}

func TestNextToken_LineTracking(t *testing.T) {
	input := `x
y
z`

	l := New(input)

	tok := l.NextToken() // x on line 1
	if tok.Line != 1 {
		t.Errorf("expected line 1, got %d", tok.Line)
	}

	tok = l.NextToken() // y on line 2
	if tok.Line != 2 {
		t.Errorf("expected line 2, got %d", tok.Line)
	}

	tok = l.NextToken() // z on line 3
	if tok.Line != 3 {
		t.Errorf("expected line 3, got %d", tok.Line)
	}
}

func TestNextToken_FullProgram(t *testing.T) {
	input := `function main() {
    x := 42;
    print(x);
}`

	expected := []TokenType{
		FUNCTION, IDENT, LPAREN, RPAREN, LBRACE,
		IDENT, WALRUS, INT, SEMICOLON,
		IDENT, LPAREN, IDENT, RPAREN, SEMICOLON,
		RBRACE, EOF,
	}

	l := New(input)
	for i, expectedType := range expected {
		tok := l.NextToken()
		if tok.Type != expectedType {
			t.Errorf("token[%d] - expected %v, got %v (%q)",
				i, expectedType, tok.Type, tok.Literal)
		}
	}
}

func TestLookupIdent(t *testing.T) {
	tests := []struct {
		ident    string
		expected TokenType
	}{
		{"function", FUNCTION},
		{"if", IF},
		{"defer", DEFER},
		{"myVar", IDENT},
		{"notAKeyword", IDENT},
	}

	for _, tt := range tests {
		result := LookupIdent(tt.ident)
		if result != tt.expected {
			t.Errorf("LookupIdent(%q) = %v, want %v",
				tt.ident, result, tt.expected)
		}
	}
}
