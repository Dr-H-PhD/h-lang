package ast

import (
	"testing"

	"github.com/Dr-H-PhD/h-lang/pkg/lexer"
)

func TestProgram_String(t *testing.T) {
	program := &Program{
		Statements: []Statement{
			&InferStatement{
				Token: lexer.Token{Type: lexer.IDENT, Literal: "x"},
				Name:  &Identifier{Token: lexer.Token{Literal: "x"}, Value: "x"},
				Value: &IntegerLiteral{Token: lexer.Token{Literal: "42"}, Value: 42},
			},
		},
	}

	expected := "x := 42;"
	if program.String() != expected {
		t.Errorf("expected %q, got %q", expected, program.String())
	}
}

func TestIdentifier_String(t *testing.T) {
	ident := &Identifier{
		Token: lexer.Token{Literal: "myVar"},
		Value: "myVar",
	}

	if ident.String() != "myVar" {
		t.Errorf("expected 'myVar', got %q", ident.String())
	}
}

func TestIntegerLiteral_String(t *testing.T) {
	lit := &IntegerLiteral{
		Token: lexer.Token{Literal: "42"},
		Value: 42,
	}

	if lit.String() != "42" {
		t.Errorf("expected '42', got %q", lit.String())
	}
}

func TestFloatLiteral_String(t *testing.T) {
	lit := &FloatLiteral{
		Token: lexer.Token{Literal: "3.14"},
		Value: 3.14,
	}

	if lit.String() != "3.14" {
		t.Errorf("expected '3.14', got %q", lit.String())
	}
}

func TestStringLiteral_String(t *testing.T) {
	lit := &StringLiteral{
		Token: lexer.Token{Literal: "hello"},
		Value: "hello",
	}

	if lit.String() != `"hello"` {
		t.Errorf("expected '\"hello\"', got %q", lit.String())
	}
}

func TestBooleanLiteral_String(t *testing.T) {
	trueLit := &BooleanLiteral{
		Token: lexer.Token{Literal: "true"},
		Value: true,
	}

	if trueLit.String() != "true" {
		t.Errorf("expected 'true', got %q", trueLit.String())
	}

	falseLit := &BooleanLiteral{
		Token: lexer.Token{Literal: "false"},
		Value: false,
	}

	if falseLit.String() != "false" {
		t.Errorf("expected 'false', got %q", falseLit.String())
	}
}

func TestNullLiteral_String(t *testing.T) {
	lit := &NullLiteral{
		Token: lexer.Token{Literal: "null"},
	}

	if lit.String() != "null" {
		t.Errorf("expected 'null', got %q", lit.String())
	}
}

func TestPrefixExpression_String(t *testing.T) {
	expr := &PrefixExpression{
		Token:    lexer.Token{Literal: "-"},
		Operator: "-",
		Right:    &IntegerLiteral{Token: lexer.Token{Literal: "5"}, Value: 5},
	}

	if expr.String() != "(-5)" {
		t.Errorf("expected '(-5)', got %q", expr.String())
	}
}

func TestInfixExpression_String(t *testing.T) {
	expr := &InfixExpression{
		Token:    lexer.Token{Literal: "+"},
		Left:     &IntegerLiteral{Token: lexer.Token{Literal: "1"}, Value: 1},
		Operator: "+",
		Right:    &IntegerLiteral{Token: lexer.Token{Literal: "2"}, Value: 2},
	}

	if expr.String() != "(1 + 2)" {
		t.Errorf("expected '(1 + 2)', got %q", expr.String())
	}
}

func TestIfStatement_String(t *testing.T) {
	stmt := &IfStatement{
		Token: lexer.Token{Literal: "if"},
		Condition: &BooleanLiteral{
			Token: lexer.Token{Literal: "true"},
			Value: true,
		},
		Consequence: &BlockStatement{
			Statements: []Statement{
				&ReturnStatement{
					Token: lexer.Token{Literal: "return"},
					Value: &IntegerLiteral{Token: lexer.Token{Literal: "1"}, Value: 1},
				},
			},
		},
	}

	result := stmt.String()
	// Check that it contains the key parts (formatting may vary)
	if result == "" {
		t.Error("if statement string should not be empty")
	}
	if !containsAll(result, "if", "true", "return", "1") {
		t.Errorf("unexpected string: %q", result)
	}
}

func containsAll(s string, substrings ...string) bool {
	for _, sub := range substrings {
		if !contains(s, sub) {
			return false
		}
	}
	return true
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && (s[0:len(substr)] == substr || contains(s[1:], substr)))
}

func TestReturnStatement_String(t *testing.T) {
	stmt := &ReturnStatement{
		Token: lexer.Token{Literal: "return"},
		Value: &IntegerLiteral{Token: lexer.Token{Literal: "42"}, Value: 42},
	}

	if stmt.String() != "return 42;" {
		t.Errorf("expected 'return 42;', got %q", stmt.String())
	}
}

func TestDeferStatement_String(t *testing.T) {
	stmt := &DeferStatement{
		Token: lexer.Token{Literal: "defer"},
		Statement: &FreeStatement{
			Token: lexer.Token{Literal: "free"},
			Value: &Identifier{Token: lexer.Token{Literal: "x"}, Value: "x"},
		},
	}

	if stmt.String() != "defer free(x);" {
		t.Errorf("expected 'defer free(x);', got %q", stmt.String())
	}
}

func TestFreeStatement_String(t *testing.T) {
	stmt := &FreeStatement{
		Token: lexer.Token{Literal: "free"},
		Value: &Identifier{Token: lexer.Token{Literal: "ptr"}, Value: "ptr"},
	}

	if stmt.String() != "free(ptr);" {
		t.Errorf("expected 'free(ptr);', got %q", stmt.String())
	}
}

func TestCallExpression_String(t *testing.T) {
	expr := &CallExpression{
		Token:    lexer.Token{Literal: "("},
		Function: &Identifier{Token: lexer.Token{Literal: "add"}, Value: "add"},
		Arguments: []Expression{
			&IntegerLiteral{Token: lexer.Token{Literal: "1"}, Value: 1},
			&IntegerLiteral{Token: lexer.Token{Literal: "2"}, Value: 2},
		},
	}

	if expr.String() != "add(1, 2)" {
		t.Errorf("expected 'add(1, 2)', got %q", expr.String())
	}
}

func TestMemberExpression_String(t *testing.T) {
	expr := &MemberExpression{
		Token:  lexer.Token{Literal: "."},
		Object: &Identifier{Token: lexer.Token{Literal: "user"}, Value: "user"},
		Member: &Identifier{Token: lexer.Token{Literal: "name"}, Value: "name"},
	}

	if expr.String() != "(user.name)" {
		t.Errorf("expected '(user.name)', got %q", expr.String())
	}
}

func TestTypeAnnotation_String(t *testing.T) {
	// Test basic type
	basic := &TypeAnnotation{Name: "int"}
	if basic.String() != "int" {
		t.Errorf("expected 'int', got %q", basic.String())
	}

	// Test pointer type
	ptr := &TypeAnnotation{Name: "int", IsPtr: true}
	result := ptr.String()
	if result != "*int" {
		t.Errorf("expected '*int', got %q", result)
	}

	// Test slice type
	slice := &TypeAnnotation{Name: "int", ArrayLen: -1}
	result = slice.String()
	if result != "[]int" {
		t.Errorf("expected '[]int', got %q", result)
	}
}
