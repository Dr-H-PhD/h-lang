package lexer

import "fmt"

// TokenType represents the type of a token
type TokenType int

const (
	// Special tokens
	ILLEGAL TokenType = iota
	EOF
	COMMENT

	// Identifiers and literals
	IDENT  // variable names, function names
	INT    // 123
	FLOAT  // 3.14
	STRING // "hello"
	CHAR   // 'a'

	// Operators
	ASSIGN    // =
	PLUS      // +
	MINUS     // -
	ASTERISK  // *
	SLASH     // /
	PERCENT   // %
	BANG      // !
	AMPERSAND // &

	LT  // <
	GT  // >
	LTE // <=
	GTE // >=
	EQ  // ==
	NEQ // !=
	AND // &&
	OR  // ||

	INCREMENT // ++
	DECREMENT // --

	PLUS_ASSIGN  // +=
	MINUS_ASSIGN // -=
	MUL_ASSIGN   // *=
	DIV_ASSIGN   // /=

	WALRUS // :=

	// Delimiters
	COMMA     // ,
	SEMICOLON // ;
	COLON     // :
	DOT       // .
	ARROW     // =>

	LPAREN   // (
	RPAREN   // )
	LBRACE   // {
	RBRACE   // }
	LBRACKET // [
	RBRACKET // ]

	// Keywords
	FUNCTION
	STRUCT
	ENUM
	IMPORT
	IF
	ELSE
	FOR
	WHILE
	RETURN
	CONST
	VAR
	PUBLIC
	NULL
	TRUE
	FALSE
	ALLOC
	FREE
	DEFER
	LEN
	MAKE
	RANGE
	BREAK
	CONTINUE
	MAP
	DELETE

	// Types
	TYPE_INT
	TYPE_FLOAT
	TYPE_STRING
	TYPE_CHAR
	TYPE_BOOL
	TYPE_VOID
)

var tokenNames = map[TokenType]string{
	ILLEGAL:      "ILLEGAL",
	EOF:          "EOF",
	COMMENT:      "COMMENT",
	IDENT:        "IDENT",
	INT:          "INT",
	FLOAT:        "FLOAT",
	STRING:       "STRING",
	CHAR:         "CHAR",
	ASSIGN:       "=",
	PLUS:         "+",
	MINUS:        "-",
	ASTERISK:     "*",
	SLASH:        "/",
	PERCENT:      "%",
	BANG:         "!",
	AMPERSAND:    "&",
	LT:           "<",
	GT:           ">",
	LTE:          "<=",
	GTE:          ">=",
	EQ:           "==",
	NEQ:          "!=",
	AND:          "&&",
	OR:           "||",
	INCREMENT:    "++",
	DECREMENT:    "--",
	PLUS_ASSIGN:  "+=",
	MINUS_ASSIGN: "-=",
	MUL_ASSIGN:   "*=",
	DIV_ASSIGN:   "/=",
	WALRUS:       ":=",
	COMMA:        ",",
	SEMICOLON:    ";",
	COLON:        ":",
	DOT:          ".",
	ARROW:        "=>",
	LPAREN:       "(",
	RPAREN:       ")",
	LBRACE:       "{",
	RBRACE:       "}",
	LBRACKET:     "[",
	RBRACKET:     "]",
	FUNCTION:     "function",
	STRUCT:       "struct",
	ENUM:         "enum",
	IMPORT:       "import",
	IF:           "if",
	ELSE:         "else",
	FOR:          "for",
	WHILE:        "while",
	RETURN:       "return",
	CONST:        "const",
	VAR:          "var",
	PUBLIC:       "public",
	NULL:         "null",
	TRUE:         "true",
	FALSE:        "false",
	ALLOC:        "alloc",
	FREE:         "free",
	DEFER:        "defer",
	LEN:          "len",
	MAKE:         "make",
	RANGE:        "range",
	BREAK:        "break",
	CONTINUE:     "continue",
	MAP:          "map",
	DELETE:       "delete",
	TYPE_INT:     "int",
	TYPE_FLOAT:   "float",
	TYPE_STRING:  "string",
	TYPE_CHAR:    "char",
	TYPE_BOOL:    "bool",
	TYPE_VOID:    "void",
}

func (t TokenType) String() string {
	if name, ok := tokenNames[t]; ok {
		return name
	}
	return "UNKNOWN"
}

var keywords = map[string]TokenType{
	"function": FUNCTION,
	"struct":   STRUCT,
	"enum":     ENUM,
	"import":   IMPORT,
	"if":       IF,
	"else":     ELSE,
	"for":      FOR,
	"while":    WHILE,
	"return":   RETURN,
	"const":    CONST,
	"var":      VAR,
	"public":   PUBLIC,
	"null":     NULL,
	"true":     TRUE,
	"false":    FALSE,
	"alloc":    ALLOC,
	"free":     FREE,
	"defer":    DEFER,
	"len":      LEN,
	"make":     MAKE,
	"range":    RANGE,
	"break":    BREAK,
	"continue": CONTINUE,
	"map":      MAP,
	"delete":   DELETE,
	"int":      TYPE_INT,
	"float":    TYPE_FLOAT,
	"string":   TYPE_STRING,
	"char":     TYPE_CHAR,
	"bool":     TYPE_BOOL,
	"void":     TYPE_VOID,
}

// LookupIdent checks if an identifier is a keyword
func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}

// Token represents a lexical token
type Token struct {
	Type    TokenType
	Literal string
	Line    int
	Column  int
}

// Position returns formatted position
func (t Token) Position() string {
	return fmt.Sprintf("%d:%d", t.Line, t.Column)
}
