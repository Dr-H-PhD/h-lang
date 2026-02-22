package parser

import (
	"fmt"
	"strconv"

	"github.com/Dr-H-PhD/h-lang/pkg/ast"
	"github.com/Dr-H-PhD/h-lang/pkg/lexer"
)

// Operator precedence levels
const (
	_ int = iota
	LOWEST
	ASSIGN      // =, +=, -=
	OR          // ||
	AND         // &&
	EQUALS      // ==, !=
	LESSGREATER // <, >, <=, >=
	SUM         // +, -
	PRODUCT     // *, /, %
	PREFIX      // -x, !x, &x, *x
	POSTFIX     // x++, x--
	CALL        // foo()
	INDEX       // arr[0]
	MEMBER      // obj.field
)

var precedences = map[lexer.TokenType]int{
	lexer.ASSIGN:       ASSIGN,
	lexer.PLUS_ASSIGN:  ASSIGN,
	lexer.MINUS_ASSIGN: ASSIGN,
	lexer.MUL_ASSIGN:   ASSIGN,
	lexer.DIV_ASSIGN:   ASSIGN,
	lexer.OR:           OR,
	lexer.AND:          AND,
	lexer.EQ:           EQUALS,
	lexer.NEQ:          EQUALS,
	lexer.LT:           LESSGREATER,
	lexer.GT:           LESSGREATER,
	lexer.LTE:          LESSGREATER,
	lexer.GTE:          LESSGREATER,
	lexer.PLUS:         SUM,
	lexer.MINUS:        SUM,
	lexer.ASTERISK:     PRODUCT,
	lexer.SLASH:        PRODUCT,
	lexer.PERCENT:      PRODUCT,
	lexer.INCREMENT:    POSTFIX,
	lexer.DECREMENT:    POSTFIX,
	lexer.LPAREN:       CALL,
	lexer.LBRACKET:     INDEX,
	lexer.DOT:          MEMBER,
}

// Parser parses H-lang source into AST
type Parser struct {
	l      *lexer.Lexer
	errors []string

	curToken  lexer.Token
	peekToken lexer.Token

	prefixParseFns map[lexer.TokenType]prefixParseFn
	infixParseFns  map[lexer.TokenType]infixParseFn
}

type prefixParseFn func() ast.Expression
type infixParseFn func(ast.Expression) ast.Expression

// New creates a new Parser
func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	p.prefixParseFns = make(map[lexer.TokenType]prefixParseFn)
	p.registerPrefix(lexer.IDENT, p.parseIdentifier)
	p.registerPrefix(lexer.INT, p.parseIntegerLiteral)
	p.registerPrefix(lexer.FLOAT, p.parseFloatLiteral)
	p.registerPrefix(lexer.STRING, p.parseStringLiteral)
	p.registerPrefix(lexer.CHAR, p.parseCharLiteral)
	p.registerPrefix(lexer.TRUE, p.parseBooleanLiteral)
	p.registerPrefix(lexer.FALSE, p.parseBooleanLiteral)
	p.registerPrefix(lexer.NULL, p.parseNullLiteral)
	p.registerPrefix(lexer.BANG, p.parsePrefixExpression)
	p.registerPrefix(lexer.MINUS, p.parsePrefixExpression)
	p.registerPrefix(lexer.AMPERSAND, p.parsePrefixExpression)
	p.registerPrefix(lexer.ASTERISK, p.parsePrefixExpression)
	p.registerPrefix(lexer.LPAREN, p.parseGroupedOrCast)
	p.registerPrefix(lexer.LBRACKET, p.parseArrayOrSliceLiteral)
	p.registerPrefix(lexer.ALLOC, p.parseAllocExpression)
	p.registerPrefix(lexer.LEN, p.parseLenExpression)
	p.registerPrefix(lexer.MAKE, p.parseMakeExpression)
	p.registerPrefix(lexer.MAP, p.parseMapLiteral)

	p.infixParseFns = make(map[lexer.TokenType]infixParseFn)
	p.registerInfix(lexer.PLUS, p.parseInfixExpression)
	p.registerInfix(lexer.MINUS, p.parseInfixExpression)
	p.registerInfix(lexer.ASTERISK, p.parseInfixExpression)
	p.registerInfix(lexer.SLASH, p.parseInfixExpression)
	p.registerInfix(lexer.PERCENT, p.parseInfixExpression)
	p.registerInfix(lexer.EQ, p.parseInfixExpression)
	p.registerInfix(lexer.NEQ, p.parseInfixExpression)
	p.registerInfix(lexer.LT, p.parseInfixExpression)
	p.registerInfix(lexer.GT, p.parseInfixExpression)
	p.registerInfix(lexer.LTE, p.parseInfixExpression)
	p.registerInfix(lexer.GTE, p.parseInfixExpression)
	p.registerInfix(lexer.AND, p.parseInfixExpression)
	p.registerInfix(lexer.OR, p.parseInfixExpression)
	p.registerInfix(lexer.ASSIGN, p.parseAssignExpression)
	p.registerInfix(lexer.PLUS_ASSIGN, p.parseAssignExpression)
	p.registerInfix(lexer.MINUS_ASSIGN, p.parseAssignExpression)
	p.registerInfix(lexer.MUL_ASSIGN, p.parseAssignExpression)
	p.registerInfix(lexer.DIV_ASSIGN, p.parseAssignExpression)
	p.registerInfix(lexer.INCREMENT, p.parsePostfixExpression)
	p.registerInfix(lexer.DECREMENT, p.parsePostfixExpression)
	p.registerInfix(lexer.LPAREN, p.parseCallExpression)
	p.registerInfix(lexer.LBRACKET, p.parseIndexExpression)
	p.registerInfix(lexer.DOT, p.parseMemberExpression)

	// Read two tokens to initialize curToken and peekToken
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) registerPrefix(tokenType lexer.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType lexer.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()

	// Skip comments
	for p.peekToken.Type == lexer.COMMENT {
		p.peekToken = p.l.NextToken()
	}
}

func (p *Parser) curTokenIs(t lexer.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t lexer.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t lexer.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}
	p.peekError(t)
	return false
}

func (p *Parser) peekError(t lexer.TokenType) {
	msg := fmt.Sprintf("line %d: expected %s, got %s instead",
		p.peekToken.Line, t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}

// ParseProgram parses the entire program
func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for !p.curTokenIs(lexer.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}

	return program
}

func (p *Parser) parseStatement() ast.Statement {
	// Skip comments at statement level
	for p.curTokenIs(lexer.COMMENT) {
		p.nextToken()
	}

	switch p.curToken.Type {
	case lexer.IMPORT:
		return p.parseImportStatement()
	case lexer.PUBLIC:
		return p.parsePublicStatement()
	case lexer.FUNCTION:
		return p.parseFunctionStatement(false)
	case lexer.STRUCT:
		return p.parseStructStatement(false)
	case lexer.VAR:
		return p.parseVarStatement()
	case lexer.CONST:
		return p.parseConstStatement()
	case lexer.RETURN:
		return p.parseReturnStatement()
	case lexer.IF:
		return p.parseIfStatement()
	case lexer.FOR:
		return p.parseForStatement()
	case lexer.WHILE:
		return p.parseWhileStatement()
	case lexer.FREE:
		return p.parseFreeStatement()
	case lexer.DEFER:
		return p.parseDeferStatement()
	case lexer.BREAK:
		return p.parseBreakStatement()
	case lexer.CONTINUE:
		return p.parseContinueStatement()
	case lexer.ENUM:
		return p.parseEnumStatement(false)
	case lexer.DELETE:
		return p.parseDeleteStatement()
	case lexer.IDENT:
		if p.peekTokenIs(lexer.WALRUS) {
			return p.parseInferStatement()
		}
		return p.parseExpressionStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseImportStatement() *ast.ImportStatement {
	stmt := &ast.ImportStatement{Token: p.curToken}

	if !p.expectPeek(lexer.STRING) {
		return nil
	}

	stmt.Path = p.curToken.Literal

	if p.peekTokenIs(lexer.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parsePublicStatement() ast.Statement {
	p.nextToken() // consume 'public'

	switch p.curToken.Type {
	case lexer.FUNCTION:
		return p.parseFunctionStatement(true)
	case lexer.STRUCT:
		return p.parseStructStatement(true)
	case lexer.ENUM:
		return p.parseEnumStatement(true)
	default:
		p.errors = append(p.errors, fmt.Sprintf("line %d: unexpected token after 'public': %s",
			p.curToken.Line, p.curToken.Type))
		return nil
	}
}

func (p *Parser) parseFunctionStatement(public bool) *ast.FunctionStatement {
	stmt := &ast.FunctionStatement{Token: p.curToken, Public: public}

	// Check for receiver: function (r *Type) name() { }
	if p.peekTokenIs(lexer.LPAREN) {
		p.nextToken() // consume 'function'
		p.nextToken() // consume '('

		receiver := &ast.Parameter{
			Name: &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal},
		}
		p.nextToken()
		receiver.Type = p.parseTypeAnnotation()

		if !p.expectPeek(lexer.RPAREN) {
			return nil
		}
		stmt.Receiver = receiver
		p.nextToken()
	} else {
		p.nextToken() // consume 'function'
	}

	// Function name
	if !p.curTokenIs(lexer.IDENT) {
		p.errors = append(p.errors, fmt.Sprintf("line %d: expected function name, got %s",
			p.curToken.Line, p.curToken.Type))
		return nil
	}
	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	// Parameters
	if !p.expectPeek(lexer.LPAREN) {
		return nil
	}
	stmt.Parameters = p.parseFunctionParameters()

	// Return type (optional)
	if !p.peekTokenIs(lexer.LBRACE) {
		p.nextToken()
		stmt.ReturnType = p.parseTypeAnnotation()
	}

	// Body
	if !p.expectPeek(lexer.LBRACE) {
		return nil
	}
	stmt.Body = p.parseBlockStatement()

	return stmt
}

func (p *Parser) parseFunctionParameters() []*ast.Parameter {
	params := []*ast.Parameter{}

	if p.peekTokenIs(lexer.RPAREN) {
		p.nextToken()
		return params
	}

	p.nextToken()

	param := &ast.Parameter{
		Name: &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal},
	}
	p.nextToken()
	param.Type = p.parseTypeAnnotation()
	params = append(params, param)

	for p.peekTokenIs(lexer.COMMA) {
		p.nextToken() // comma
		p.nextToken() // param name

		param := &ast.Parameter{
			Name: &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal},
		}
		p.nextToken()
		param.Type = p.parseTypeAnnotation()
		params = append(params, param)
	}

	if !p.expectPeek(lexer.RPAREN) {
		return nil
	}

	return params
}

func (p *Parser) parseTypeAnnotation() *ast.TypeAnnotation {
	typeAnn := &ast.TypeAnnotation{Token: p.curToken}

	// Check for pointer
	if p.curTokenIs(lexer.ASTERISK) {
		typeAnn.IsPtr = true
		p.nextToken()
	}

	// Check for map type: map[KeyType]ValueType
	if p.curTokenIs(lexer.MAP) {
		typeAnn.IsMap = true
		if !p.expectPeek(lexer.LBRACKET) {
			return nil
		}
		p.nextToken() // move to key type
		typeAnn.KeyType = p.parseTypeAnnotation()
		if !p.expectPeek(lexer.RBRACKET) {
			return nil
		}
		p.nextToken() // move to value type
		typeAnn.ValueType = p.parseTypeAnnotation()
		return typeAnn
	}

	// Check for array/slice
	if p.curTokenIs(lexer.LBRACKET) {
		p.nextToken()
		if p.curTokenIs(lexer.RBRACKET) {
			typeAnn.ArrayLen = -1 // slice
		} else if p.curTokenIs(lexer.INT) {
			len, _ := strconv.Atoi(p.curToken.Literal)
			typeAnn.ArrayLen = len
			p.nextToken() // consume number
		}
		p.nextToken() // consume ]
	}

	typeAnn.Name = p.curToken.Literal
	return typeAnn
}

func (p *Parser) parseStructStatement(public bool) *ast.StructStatement {
	stmt := &ast.StructStatement{Token: p.curToken, Public: public}

	if !p.expectPeek(lexer.IDENT) {
		return nil
	}
	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(lexer.LBRACE) {
		return nil
	}

	stmt.Fields = p.parseStructFields()

	return stmt
}

func (p *Parser) parseStructFields() []*ast.StructField {
	fields := []*ast.StructField{}

	for !p.peekTokenIs(lexer.RBRACE) && !p.peekTokenIs(lexer.EOF) {
		p.nextToken()

		field := &ast.StructField{}

		if p.curTokenIs(lexer.PUBLIC) {
			field.Public = true
			p.nextToken()
		}

		field.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
		p.nextToken()
		field.Type = p.parseTypeAnnotation()

		// Expect semicolon
		if p.peekTokenIs(lexer.SEMICOLON) {
			p.nextToken()
		}

		fields = append(fields, field)
	}

	p.nextToken() // consume }
	return fields
}

func (p *Parser) parseEnumStatement(public bool) *ast.EnumStatement {
	stmt := &ast.EnumStatement{Token: p.curToken, Public: public}

	if !p.expectPeek(lexer.IDENT) {
		return nil
	}
	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(lexer.LBRACE) {
		return nil
	}

	stmt.Values = p.parseEnumValues()

	return stmt
}

func (p *Parser) parseEnumValues() []*ast.EnumValue {
	values := []*ast.EnumValue{}

	for !p.peekTokenIs(lexer.RBRACE) && !p.peekTokenIs(lexer.EOF) {
		p.nextToken()

		value := &ast.EnumValue{
			Name: &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal},
		}

		// Check for explicit value: Red = 1
		if p.peekTokenIs(lexer.ASSIGN) {
			p.nextToken() // consume =
			p.nextToken() // move to value
			value.Value = p.parseExpression(LOWEST)
		}

		values = append(values, value)

		// Expect comma or closing brace
		if p.peekTokenIs(lexer.COMMA) {
			p.nextToken()
		}
	}

	p.nextToken() // consume }
	return values
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: p.curToken}
	block.Statements = []ast.Statement{}

	p.nextToken()

	for !p.curTokenIs(lexer.RBRACE) && !p.curTokenIs(lexer.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}

	return block
}

func (p *Parser) parseVarStatement() *ast.VarStatement {
	stmt := &ast.VarStatement{Token: p.curToken}

	if !p.expectPeek(lexer.IDENT) {
		return nil
	}
	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	// Type annotation (required for var)
	p.nextToken()
	stmt.Type = p.parseTypeAnnotation()

	// Optional initialization
	if p.peekTokenIs(lexer.ASSIGN) {
		p.nextToken()
		p.nextToken()
		stmt.Value = p.parseExpression(LOWEST)
	}

	if p.peekTokenIs(lexer.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseConstStatement() *ast.ConstStatement {
	stmt := &ast.ConstStatement{Token: p.curToken}

	if !p.expectPeek(lexer.IDENT) {
		return nil
	}
	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(lexer.WALRUS) {
		return nil
	}

	p.nextToken()
	stmt.Value = p.parseExpression(LOWEST)

	if p.peekTokenIs(lexer.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseInferStatement() *ast.InferStatement {
	stmt := &ast.InferStatement{Token: p.curToken}
	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(lexer.WALRUS) {
		return nil
	}

	p.nextToken()
	stmt.Value = p.parseExpression(LOWEST)

	if p.peekTokenIs(lexer.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}

	p.nextToken()

	if !p.curTokenIs(lexer.SEMICOLON) {
		stmt.Value = p.parseExpression(LOWEST)
	}

	if p.peekTokenIs(lexer.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseIfStatement() *ast.IfStatement {
	stmt := &ast.IfStatement{Token: p.curToken}

	p.nextToken()
	stmt.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(lexer.LBRACE) {
		return nil
	}
	stmt.Consequence = p.parseBlockStatement()

	if p.peekTokenIs(lexer.ELSE) {
		p.nextToken()

		if !p.expectPeek(lexer.LBRACE) {
			return nil
		}
		stmt.Alternative = p.parseBlockStatement()
	}

	return stmt
}

func (p *Parser) parseForStatement() ast.Statement {
	forToken := p.curToken
	p.nextToken()

	// Check if this is a for-range statement
	// Pattern: for ident := range ... or for ident, ident := range ...
	if p.curTokenIs(lexer.IDENT) {
		firstIdent := p.curToken

		// Check for "ident := range" or "ident, ident := range"
		if p.peekTokenIs(lexer.WALRUS) {
			// Single variable: for i := range arr
			p.nextToken() // consume :=
			if p.peekTokenIs(lexer.RANGE) {
				return p.parseForRangeStatement(forToken, &ast.Identifier{Token: firstIdent, Value: firstIdent.Literal}, nil)
			}
			// Not a range, fall back to regular for loop
			// Need to continue parsing as init statement
			return p.parseForStatementWithInit(forToken, firstIdent)
		} else if p.peekTokenIs(lexer.COMMA) {
			// Two variables: for i, v := range arr
			p.nextToken() // consume ,
			if !p.expectPeek(lexer.IDENT) {
				return nil
			}
			secondIdent := p.curToken
			if !p.expectPeek(lexer.WALRUS) {
				return nil
			}
			if !p.expectPeek(lexer.RANGE) {
				return nil
			}
			index := &ast.Identifier{Token: firstIdent, Value: firstIdent.Literal}
			value := &ast.Identifier{Token: secondIdent, Value: secondIdent.Literal}
			return p.parseForRangeStatementBody(forToken, index, value)
		}
	}

	// Check for blank identifier pattern: for _, v := range arr
	if p.curTokenIs(lexer.IDENT) && p.curToken.Literal == "_" {
		if p.peekTokenIs(lexer.COMMA) {
			p.nextToken() // consume ,
			if !p.expectPeek(lexer.IDENT) {
				return nil
			}
			secondIdent := p.curToken
			if !p.expectPeek(lexer.WALRUS) {
				return nil
			}
			if !p.expectPeek(lexer.RANGE) {
				return nil
			}
			value := &ast.Identifier{Token: secondIdent, Value: secondIdent.Literal}
			return p.parseForRangeStatementBody(forToken, nil, value)
		}
	}

	// Regular for loop
	return p.parseRegularForStatement(forToken)
}

func (p *Parser) parseForRangeStatement(forToken lexer.Token, firstIdent *ast.Identifier, secondIdent *ast.Identifier) *ast.ForRangeStatement {
	// We've already consumed "ident :=" and seen RANGE as peek
	p.nextToken() // consume range

	return p.parseForRangeStatementBody(forToken, firstIdent, secondIdent)
}

func (p *Parser) parseForRangeStatementBody(forToken lexer.Token, index *ast.Identifier, value *ast.Identifier) *ast.ForRangeStatement {
	stmt := &ast.ForRangeStatement{
		Token: forToken,
		Index: index,
		Value: value,
	}

	p.nextToken() // move to iterable
	stmt.Iterable = p.parseExpression(LOWEST)

	if !p.expectPeek(lexer.LBRACE) {
		return nil
	}
	stmt.Body = p.parseBlockStatement()

	return stmt
}

func (p *Parser) parseForStatementWithInit(forToken lexer.Token, firstIdent lexer.Token) *ast.ForStatement {
	// We've consumed "for ident :=" and it's not a range
	// Parse the rest as infer statement init
	stmt := &ast.ForStatement{Token: forToken}

	initStmt := &ast.InferStatement{Token: firstIdent}
	initStmt.Name = &ast.Identifier{Token: firstIdent, Value: firstIdent.Literal}

	p.nextToken() // move past :=
	initStmt.Value = p.parseExpression(LOWEST)
	stmt.Init = initStmt

	if p.peekTokenIs(lexer.SEMICOLON) {
		p.nextToken()
	}
	if p.curTokenIs(lexer.SEMICOLON) {
		p.nextToken()
	}

	// Condition
	if !p.curTokenIs(lexer.SEMICOLON) {
		stmt.Condition = p.parseExpression(LOWEST)
	}
	if !p.expectPeek(lexer.SEMICOLON) {
		return nil
	}
	p.nextToken()

	// Post
	if !p.curTokenIs(lexer.LBRACE) {
		stmt.Post = p.parseStatement()
	}

	if !p.curTokenIs(lexer.LBRACE) {
		if !p.expectPeek(lexer.LBRACE) {
			return nil
		}
	}
	stmt.Body = p.parseBlockStatement()

	return stmt
}

func (p *Parser) parseRegularForStatement(forToken lexer.Token) *ast.ForStatement {
	stmt := &ast.ForStatement{Token: forToken}

	// Init
	if !p.curTokenIs(lexer.SEMICOLON) {
		stmt.Init = p.parseStatement()
	}
	if p.curTokenIs(lexer.SEMICOLON) {
		p.nextToken()
	}

	// Condition
	if !p.curTokenIs(lexer.SEMICOLON) {
		stmt.Condition = p.parseExpression(LOWEST)
	}
	if !p.expectPeek(lexer.SEMICOLON) {
		return nil
	}
	p.nextToken()

	// Post
	if !p.curTokenIs(lexer.LBRACE) {
		stmt.Post = p.parseStatement()
	}

	if !p.curTokenIs(lexer.LBRACE) {
		if !p.expectPeek(lexer.LBRACE) {
			return nil
		}
	}
	stmt.Body = p.parseBlockStatement()

	return stmt
}

func (p *Parser) parseWhileStatement() *ast.WhileStatement {
	stmt := &ast.WhileStatement{Token: p.curToken}

	p.nextToken()
	stmt.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(lexer.LBRACE) {
		return nil
	}
	stmt.Body = p.parseBlockStatement()

	return stmt
}

func (p *Parser) parseFreeStatement() *ast.FreeStatement {
	stmt := &ast.FreeStatement{Token: p.curToken}

	if !p.expectPeek(lexer.LPAREN) {
		return nil
	}
	p.nextToken()

	stmt.Value = p.parseExpression(LOWEST)

	if !p.expectPeek(lexer.RPAREN) {
		return nil
	}

	if p.peekTokenIs(lexer.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseDeferStatement() *ast.DeferStatement {
	stmt := &ast.DeferStatement{Token: p.curToken}

	p.nextToken()

	// Parse the deferred statement
	stmt.Statement = p.parseStatement()

	return stmt
}

func (p *Parser) parseBreakStatement() *ast.BreakStatement {
	stmt := &ast.BreakStatement{Token: p.curToken}

	if p.peekTokenIs(lexer.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseContinueStatement() *ast.ContinueStatement {
	stmt := &ast.ContinueStatement{Token: p.curToken}

	if p.peekTokenIs(lexer.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}
	stmt.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(lexer.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.errors = append(p.errors, fmt.Sprintf("line %d: no prefix parse function for %s",
			p.curToken.Line, p.curToken.Type))
		return nil
	}
	leftExp := prefix()

	for !p.peekTokenIs(lexer.SEMICOLON) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}

		p.nextToken()
		leftExp = infix(leftExp)
	}

	return leftExp
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.curToken}

	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		p.errors = append(p.errors, fmt.Sprintf("line %d: could not parse %q as integer",
			p.curToken.Line, p.curToken.Literal))
		return nil
	}

	lit.Value = value
	return lit
}

func (p *Parser) parseFloatLiteral() ast.Expression {
	lit := &ast.FloatLiteral{Token: p.curToken}

	value, err := strconv.ParseFloat(p.curToken.Literal, 64)
	if err != nil {
		p.errors = append(p.errors, fmt.Sprintf("line %d: could not parse %q as float",
			p.curToken.Line, p.curToken.Literal))
		return nil
	}

	lit.Value = value
	return lit
}

func (p *Parser) parseStringLiteral() ast.Expression {
	return &ast.StringLiteral{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseCharLiteral() ast.Expression {
	value := byte(0)
	if len(p.curToken.Literal) > 0 {
		value = p.curToken.Literal[0]
	}
	return &ast.CharLiteral{Token: p.curToken, Value: value}
}

func (p *Parser) parseBooleanLiteral() ast.Expression {
	return &ast.BooleanLiteral{Token: p.curToken, Value: p.curTokenIs(lexer.TRUE)}
}

func (p *Parser) parseNullLiteral() ast.Expression {
	return &ast.NullLiteral{Token: p.curToken}
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}

	p.nextToken()
	expression.Right = p.parseExpression(PREFIX)

	return expression
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}

	precedence := p.curPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)

	return expression
}

func (p *Parser) parsePostfixExpression(left ast.Expression) ast.Expression {
	return &ast.PostfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}
}

func (p *Parser) parseAssignExpression(left ast.Expression) ast.Expression {
	expression := &ast.AssignExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}

	p.nextToken()
	expression.Value = p.parseExpression(LOWEST)

	return expression
}

func (p *Parser) parseGroupedOrCast() ast.Expression {
	// Could be (expr) or (type)expr
	p.nextToken()

	// Check if it's a type cast
	if p.isType() {
		typeAnn := p.parseTypeAnnotation()
		if !p.expectPeek(lexer.RPAREN) {
			return nil
		}
		p.nextToken()
		value := p.parseExpression(PREFIX)
		return &ast.CastExpression{
			Token:      p.curToken,
			TargetType: typeAnn,
			Value:      value,
		}
	}

	// It's a grouped expression
	exp := p.parseExpression(LOWEST)

	if !p.expectPeek(lexer.RPAREN) {
		return nil
	}

	return exp
}

func (p *Parser) isType() bool {
	switch p.curToken.Type {
	case lexer.TYPE_INT, lexer.TYPE_FLOAT, lexer.TYPE_STRING,
		lexer.TYPE_CHAR, lexer.TYPE_BOOL, lexer.TYPE_VOID,
		lexer.ASTERISK:
		return true
	}
	return false
}

func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	exp := &ast.CallExpression{Token: p.curToken, Function: function}
	exp.Arguments = p.parseExpressionList(lexer.RPAREN)
	return exp
}

func (p *Parser) parseExpressionList(end lexer.TokenType) []ast.Expression {
	list := []ast.Expression{}

	if p.peekTokenIs(end) {
		p.nextToken()
		return list
	}

	p.nextToken()
	list = append(list, p.parseExpression(LOWEST))

	for p.peekTokenIs(lexer.COMMA) {
		p.nextToken()
		p.nextToken()
		list = append(list, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(end) {
		return nil
	}

	return list
}

func (p *Parser) parseIndexExpression(left ast.Expression) ast.Expression {
	exp := &ast.IndexExpression{Token: p.curToken, Left: left}

	p.nextToken()
	exp.Index = p.parseExpression(LOWEST)

	if !p.expectPeek(lexer.RBRACKET) {
		return nil
	}

	return exp
}

func (p *Parser) parseMemberExpression(left ast.Expression) ast.Expression {
	exp := &ast.MemberExpression{Token: p.curToken, Object: left}

	p.nextToken()
	exp.Member = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	return exp
}

func (p *Parser) parseArrayOrSliceLiteral() ast.Expression {
	array := &ast.ArrayLiteral{Token: p.curToken}

	// Check if this is a typed array: [5]int{...} or []int{...}
	p.nextToken() // move past [

	if p.curTokenIs(lexer.RBRACKET) {
		// Slice type: []type{...}
		p.nextToken() // move past ]
		if p.curTokenIs(lexer.IDENT) || p.isType() {
			array.Type = &ast.TypeAnnotation{
				Token:    p.curToken,
				Name:     p.curToken.Literal,
				ArrayLen: -1, // slice
			}
			p.nextToken() // move past type
			if p.curTokenIs(lexer.LBRACE) {
				array.Elements = p.parseExpressionListBrace()
			}
			return array
		}
		// Just empty brackets - this is an error or empty slice
		return array
	} else if p.curTokenIs(lexer.INT) {
		// Fixed array: [5]type{...}
		length, _ := strconv.Atoi(p.curToken.Literal)
		p.nextToken() // move past number
		if !p.curTokenIs(lexer.RBRACKET) {
			return nil
		}
		p.nextToken() // move past ]
		if p.curTokenIs(lexer.IDENT) || p.isType() {
			array.Type = &ast.TypeAnnotation{
				Token:    p.curToken,
				Name:     p.curToken.Literal,
				ArrayLen: length,
			}
			p.nextToken() // move past type
			if p.curTokenIs(lexer.LBRACE) {
				array.Elements = p.parseExpressionListBrace()
			}
			return array
		}
		return nil
	}

	// Regular array literal without type: [1, 2, 3]
	// We already moved past [, so parse elements until ]
	if p.curTokenIs(lexer.RBRACKET) {
		return array
	}

	array.Elements = append(array.Elements, p.parseExpression(LOWEST))

	for p.peekTokenIs(lexer.COMMA) {
		p.nextToken() // comma
		p.nextToken() // next element
		array.Elements = append(array.Elements, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(lexer.RBRACKET) {
		return nil
	}

	return array
}

func (p *Parser) parseExpressionListBrace() []ast.Expression {
	list := []ast.Expression{}

	if !p.curTokenIs(lexer.LBRACE) {
		return list
	}

	if p.peekTokenIs(lexer.RBRACE) {
		p.nextToken()
		return list
	}

	p.nextToken()
	list = append(list, p.parseExpression(LOWEST))

	for p.peekTokenIs(lexer.COMMA) {
		p.nextToken()
		p.nextToken()
		list = append(list, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(lexer.RBRACE) {
		return nil
	}

	return list
}

func (p *Parser) parseLenExpression() ast.Expression {
	exp := &ast.CallExpression{
		Token:    p.curToken,
		Function: &ast.Identifier{Token: p.curToken, Value: "len"},
	}

	if !p.expectPeek(lexer.LPAREN) {
		return nil
	}
	exp.Arguments = p.parseExpressionList(lexer.RPAREN)

	return exp
}

func (p *Parser) parseMakeExpression() ast.Expression {
	exp := &ast.MakeExpression{Token: p.curToken}

	if !p.expectPeek(lexer.LPAREN) {
		return nil
	}
	p.nextToken()

	// Parse type
	exp.Type = p.parseTypeAnnotation()

	// Optional capacity
	if p.peekTokenIs(lexer.COMMA) {
		p.nextToken() // comma
		p.nextToken()
		exp.Length = p.parseExpression(LOWEST)

		if p.peekTokenIs(lexer.COMMA) {
			p.nextToken() // comma
			p.nextToken()
			exp.Capacity = p.parseExpression(LOWEST)
		}
	}

	if !p.expectPeek(lexer.RPAREN) {
		return nil
	}

	return exp
}

func (p *Parser) parseAllocExpression() ast.Expression {
	exp := &ast.AllocExpression{Token: p.curToken}

	if !p.expectPeek(lexer.LPAREN) {
		return nil
	}
	p.nextToken()

	exp.Type = p.parseTypeAnnotation()

	if !p.expectPeek(lexer.RPAREN) {
		return nil
	}

	return exp
}

func (p *Parser) parseMapLiteral() ast.Expression {
	ml := &ast.MapLiteral{Token: p.curToken}

	// Parse type: map[KeyType]ValueType
	ml.Type = &ast.TypeAnnotation{Token: p.curToken, IsMap: true}

	if !p.expectPeek(lexer.LBRACKET) {
		return nil
	}
	p.nextToken() // move to key type
	ml.Type.KeyType = p.parseTypeAnnotation()

	if !p.expectPeek(lexer.RBRACKET) {
		return nil
	}
	p.nextToken() // move to value type
	ml.Type.ValueType = p.parseTypeAnnotation()

	// Parse the literal body { key: value, ... }
	if !p.expectPeek(lexer.LBRACE) {
		return nil
	}

	ml.Pairs = p.parseMapPairs()

	return ml
}

func (p *Parser) parseMapPairs() []*ast.MapPair {
	pairs := []*ast.MapPair{}

	if p.peekTokenIs(lexer.RBRACE) {
		p.nextToken()
		return pairs
	}

	p.nextToken()

	pair := &ast.MapPair{}
	pair.Key = p.parseExpression(LOWEST)

	if !p.expectPeek(lexer.COLON) {
		return nil
	}
	p.nextToken()

	pair.Value = p.parseExpression(LOWEST)
	pairs = append(pairs, pair)

	for p.peekTokenIs(lexer.COMMA) {
		p.nextToken() // comma
		p.nextToken() // next key

		pair := &ast.MapPair{}
		pair.Key = p.parseExpression(LOWEST)

		if !p.expectPeek(lexer.COLON) {
			return nil
		}
		p.nextToken()

		pair.Value = p.parseExpression(LOWEST)
		pairs = append(pairs, pair)
	}

	if !p.expectPeek(lexer.RBRACE) {
		return nil
	}

	return pairs
}

func (p *Parser) parseDeleteStatement() *ast.DeleteStatement {
	stmt := &ast.DeleteStatement{Token: p.curToken}

	if !p.expectPeek(lexer.LPAREN) {
		return nil
	}
	p.nextToken()

	stmt.Map = p.parseExpression(LOWEST)

	if !p.expectPeek(lexer.COMMA) {
		return nil
	}
	p.nextToken()

	stmt.Key = p.parseExpression(LOWEST)

	if !p.expectPeek(lexer.RPAREN) {
		return nil
	}

	if p.peekTokenIs(lexer.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}
