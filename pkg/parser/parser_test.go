package parser

import (
	"testing"

	"github.com/Dr-H-PhD/h-lang/pkg/ast"
	"github.com/Dr-H-PhD/h-lang/pkg/lexer"
)

func TestInferStatement(t *testing.T) {
	tests := []struct {
		input         string
		expectedName  string
		expectedValue interface{}
	}{
		{"x := 5;", "x", int64(5)},
		{"y := 10;", "y", int64(10)},
		{"foobar := 838383;", "foobar", int64(838383)},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("expected 1 statement, got %d", len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.InferStatement)
		if !ok {
			t.Fatalf("expected InferStatement, got %T", program.Statements[0])
		}

		if stmt.Name.Value != tt.expectedName {
			t.Errorf("name: expected %q, got %q", tt.expectedName, stmt.Name.Value)
		}
	}
}

func TestConstStatement(t *testing.T) {
	input := `const PI := 3.14159;`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ConstStatement)
	if !ok {
		t.Fatalf("expected ConstStatement, got %T", program.Statements[0])
	}

	if stmt.Name.Value != "PI" {
		t.Errorf("name: expected 'PI', got %q", stmt.Name.Value)
	}

	lit, ok := stmt.Value.(*ast.FloatLiteral)
	if !ok {
		t.Fatalf("expected FloatLiteral, got %T", stmt.Value)
	}

	if lit.Value != 3.14159 {
		t.Errorf("value: expected 3.14159, got %f", lit.Value)
	}
}

func TestVarStatement(t *testing.T) {
	input := `var count int = 0;`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.VarStatement)
	if !ok {
		t.Fatalf("expected VarStatement, got %T", program.Statements[0])
	}

	if stmt.Name.Value != "count" {
		t.Errorf("name: expected 'count', got %q", stmt.Name.Value)
	}

	if stmt.Type.Name != "int" {
		t.Errorf("type: expected 'int', got %q", stmt.Type.Name)
	}
}

func TestReturnStatement(t *testing.T) {
	tests := []struct {
		input         string
		expectedValue interface{}
	}{
		{"return 5;", int64(5)},
		{"return x;", "x"},
		{"return;", nil},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("expected 1 statement, got %d", len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ReturnStatement)
		if !ok {
			t.Fatalf("expected ReturnStatement, got %T", program.Statements[0])
		}

		if tt.expectedValue == nil && stmt.Value != nil {
			t.Errorf("expected nil value, got %v", stmt.Value)
		}
	}
}

func TestFunctionStatement(t *testing.T) {
	input := `function add(a int, b int) int {
    return a + b;
}`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.FunctionStatement)
	if !ok {
		t.Fatalf("expected FunctionStatement, got %T", program.Statements[0])
	}

	if stmt.Public {
		t.Error("expected non-public function")
	}

	if stmt.Name.Value != "add" {
		t.Errorf("name: expected 'add', got %q", stmt.Name.Value)
	}

	if len(stmt.Parameters) != 2 {
		t.Fatalf("expected 2 parameters, got %d", len(stmt.Parameters))
	}

	if stmt.Parameters[0].Name.Value != "a" {
		t.Errorf("param 0: expected 'a', got %q", stmt.Parameters[0].Name.Value)
	}

	if stmt.ReturnType.Name != "int" {
		t.Errorf("return type: expected 'int', got %q", stmt.ReturnType.Name)
	}
}

func TestPublicFunctionStatement(t *testing.T) {
	input := `public function greet() string {
    return "hello";
}`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.FunctionStatement)

	if !stmt.Public {
		t.Error("expected public function")
	}

	if stmt.Name.Value != "greet" {
		t.Errorf("name: expected 'greet', got %q", stmt.Name.Value)
	}
}

func TestMethodStatement(t *testing.T) {
	input := `public function (u *User) greet() string {
    return "hello";
}`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.FunctionStatement)

	if stmt.Receiver == nil {
		t.Fatal("expected receiver")
	}

	if stmt.Receiver.Name.Value != "u" {
		t.Errorf("receiver name: expected 'u', got %q", stmt.Receiver.Name.Value)
	}

	if !stmt.Receiver.Type.IsPtr {
		t.Error("expected pointer receiver")
	}

	if stmt.Receiver.Type.Name != "User" {
		t.Errorf("receiver type: expected 'User', got %q", stmt.Receiver.Type.Name)
	}
}

func TestStructStatement(t *testing.T) {
	input := `public struct User {
    public name string;
    public age int;
    email string;
}`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.StructStatement)

	if !stmt.Public {
		t.Error("expected public struct")
	}

	if stmt.Name.Value != "User" {
		t.Errorf("name: expected 'User', got %q", stmt.Name.Value)
	}

	if len(stmt.Fields) != 3 {
		t.Fatalf("expected 3 fields, got %d", len(stmt.Fields))
	}

	// Check first field
	if stmt.Fields[0].Name.Value != "name" {
		t.Errorf("field 0: expected 'name', got %q", stmt.Fields[0].Name.Value)
	}
	if !stmt.Fields[0].Public {
		t.Error("field 0 should be public")
	}

	// Check private field
	if stmt.Fields[2].Name.Value != "email" {
		t.Errorf("field 2: expected 'email', got %q", stmt.Fields[2].Name.Value)
	}
	if stmt.Fields[2].Public {
		t.Error("field 2 should be private")
	}
}

func TestIfStatement(t *testing.T) {
	input := `if x > 5 { return 10; }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.IfStatement)

	if stmt.Condition == nil {
		t.Fatal("expected condition")
	}

	if stmt.Consequence == nil {
		t.Fatal("expected consequence")
	}

	if stmt.Alternative != nil {
		t.Error("expected no alternative")
	}
}

func TestIfElseStatement(t *testing.T) {
	input := `if x > 5 { return 10; } else { return 0; }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.IfStatement)

	if stmt.Alternative == nil {
		t.Fatal("expected alternative")
	}
}

func TestForStatement(t *testing.T) {
	input := `for i := 0; i < 10; i++ { print(i); }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ForStatement)

	if stmt.Init == nil {
		t.Error("expected init")
	}

	if stmt.Condition == nil {
		t.Error("expected condition")
	}

	if stmt.Post == nil {
		t.Error("expected post")
	}

	if stmt.Body == nil {
		t.Error("expected body")
	}
}

func TestWhileStatement(t *testing.T) {
	input := `while x < 10 { x++; }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.WhileStatement)

	if stmt.Condition == nil {
		t.Error("expected condition")
	}

	if stmt.Body == nil {
		t.Error("expected body")
	}
}

func TestDeferStatement(t *testing.T) {
	input := `defer free(x);`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.DeferStatement)

	if stmt.Statement == nil {
		t.Fatal("expected deferred statement")
	}

	// Should be a FreeStatement
	_, ok := stmt.Statement.(*ast.FreeStatement)
	if !ok {
		t.Errorf("expected FreeStatement, got %T", stmt.Statement)
	}
}

func TestOperatorPrecedence(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"1 + 2;", "(1 + 2);"},
		{"1 + 2 * 3;", "(1 + (2 * 3));"},
		{"1 * 2 + 3;", "((1 * 2) + 3);"},
		{"a + b * c + d;", "((a + (b * c)) + d);"},
		{"-a * b;", "((-a) * b);"},
		{"!true;", "(!true);"},
		{"a && b || c;", "((a && b) || c);"},
		{"a == b != c;", "((a == b) != c);"},
		{"a < b == c > d;", "((a < b) == (c > d));"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		actual := program.String()
		if actual != tt.expected {
			t.Errorf("input=%q: expected=%q, got=%q", tt.input, tt.expected, actual)
		}
	}
}

func TestCallExpression(t *testing.T) {
	input := `add(1, 2 * 3, 4 + 5);`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	exp := stmt.Expression.(*ast.CallExpression)

	if len(exp.Arguments) != 3 {
		t.Fatalf("expected 3 arguments, got %d", len(exp.Arguments))
	}
}

func TestMemberExpression(t *testing.T) {
	input := `user.name;`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	exp := stmt.Expression.(*ast.MemberExpression)

	if exp.Member.Value != "name" {
		t.Errorf("member: expected 'name', got %q", exp.Member.Value)
	}
}

func TestIndexExpression(t *testing.T) {
	input := `arr[0];`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	_, ok := stmt.Expression.(*ast.IndexExpression)
	if !ok {
		t.Fatalf("expected IndexExpression, got %T", stmt.Expression)
	}
}

func TestAllocExpression(t *testing.T) {
	input := `alloc(User);`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	exp := stmt.Expression.(*ast.AllocExpression)

	if exp.Type.Name != "User" {
		t.Errorf("type: expected 'User', got %q", exp.Type.Name)
	}
}

func TestCastExpression(t *testing.T) {
	input := `(int)x;`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	exp := stmt.Expression.(*ast.CastExpression)

	if exp.TargetType.Name != "int" {
		t.Errorf("target type: expected 'int', got %q", exp.TargetType.Name)
	}
}

func TestPointerTypes(t *testing.T) {
	input := `var p *int;`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.VarStatement)

	if !stmt.Type.IsPtr {
		t.Error("expected pointer type")
	}

	if stmt.Type.Name != "int" {
		t.Errorf("type: expected 'int', got %q", stmt.Type.Name)
	}
}

func TestFixedArrayLiteral(t *testing.T) {
	input := `arr := [5]int{1, 2, 3, 4, 5};`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.InferStatement)
	arr, ok := stmt.Value.(*ast.ArrayLiteral)
	if !ok {
		t.Fatalf("expected ArrayLiteral, got %T", stmt.Value)
	}

	if arr.Type == nil {
		t.Fatal("expected array type")
	}

	if arr.Type.ArrayLen != 5 {
		t.Errorf("expected ArrayLen 5, got %d", arr.Type.ArrayLen)
	}

	if arr.Type.Name != "int" {
		t.Errorf("expected type 'int', got %q", arr.Type.Name)
	}

	if len(arr.Elements) != 5 {
		t.Errorf("expected 5 elements, got %d", len(arr.Elements))
	}
}

func TestSliceLiteral(t *testing.T) {
	input := `nums := []int{10, 20, 30};`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.InferStatement)
	arr, ok := stmt.Value.(*ast.ArrayLiteral)
	if !ok {
		t.Fatalf("expected ArrayLiteral, got %T", stmt.Value)
	}

	if arr.Type == nil {
		t.Fatal("expected array type")
	}

	if arr.Type.ArrayLen != -1 {
		t.Errorf("expected slice (ArrayLen -1), got %d", arr.Type.ArrayLen)
	}

	if len(arr.Elements) != 3 {
		t.Errorf("expected 3 elements, got %d", len(arr.Elements))
	}
}

func TestMakeExpression(t *testing.T) {
	input := `buf := make([]int, 10);`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.InferStatement)
	mk, ok := stmt.Value.(*ast.MakeExpression)
	if !ok {
		t.Fatalf("expected MakeExpression, got %T", stmt.Value)
	}

	if mk.Type.ArrayLen != -1 {
		t.Errorf("expected slice type, got ArrayLen %d", mk.Type.ArrayLen)
	}

	if mk.Length == nil {
		t.Error("expected length argument")
	}
}

func TestLenExpression(t *testing.T) {
	input := `size := len(arr);`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.InferStatement)
	call, ok := stmt.Value.(*ast.CallExpression)
	if !ok {
		t.Fatalf("expected CallExpression, got %T", stmt.Value)
	}

	fn, ok := call.Function.(*ast.Identifier)
	if !ok || fn.Value != "len" {
		t.Errorf("expected 'len' function, got %v", call.Function)
	}

	if len(call.Arguments) != 1 {
		t.Errorf("expected 1 argument, got %d", len(call.Arguments))
	}
}

func TestArrayIndexing(t *testing.T) {
	input := `x := arr[0];`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.InferStatement)
	idx, ok := stmt.Value.(*ast.IndexExpression)
	if !ok {
		t.Fatalf("expected IndexExpression, got %T", stmt.Value)
	}

	if idx.Left.String() != "arr" {
		t.Errorf("expected 'arr', got %q", idx.Left.String())
	}
}

func TestFullProgram(t *testing.T) {
	input := `
public struct User {
    public name string;
}

public function (u *User) greet() string {
    return "Hello, " + u.name;
}

function main() {
    user := alloc(User);
    defer free(user);
    user.name = "Achraf";
    print(user.greet());
}
`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 3 {
		t.Fatalf("expected 3 statements, got %d", len(program.Statements))
	}

	// Struct
	_, ok := program.Statements[0].(*ast.StructStatement)
	if !ok {
		t.Errorf("statement 0: expected StructStatement, got %T", program.Statements[0])
	}

	// Method
	_, ok = program.Statements[1].(*ast.FunctionStatement)
	if !ok {
		t.Errorf("statement 1: expected FunctionStatement, got %T", program.Statements[1])
	}

	// Main
	fn, ok := program.Statements[2].(*ast.FunctionStatement)
	if !ok {
		t.Errorf("statement 2: expected FunctionStatement, got %T", program.Statements[2])
	}
	if fn.Name.Value != "main" {
		t.Errorf("expected 'main', got %q", fn.Name.Value)
	}
}

func TestForRangeStatement(t *testing.T) {
	tests := []struct {
		input     string
		hasIndex  bool
		indexName string
		hasValue  bool
		valueName string
		iterable  string
	}{
		{
			input:     "for i, v := range arr { print(v); }",
			hasIndex:  true,
			indexName: "i",
			hasValue:  true,
			valueName: "v",
			iterable:  "arr",
		},
		{
			input:     "for i := range numbers { print(i); }",
			hasIndex:  true,
			indexName: "i",
			hasValue:  false,
			iterable:  "numbers",
		},
		{
			input:     "for _, v := range items { print(v); }",
			hasIndex:  true,
			indexName: "_",
			hasValue:  true,
			valueName: "v",
			iterable:  "items",
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("input=%q: expected 1 statement, got %d", tt.input, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ForRangeStatement)
		if !ok {
			t.Fatalf("input=%q: expected ForRangeStatement, got %T", tt.input, program.Statements[0])
		}

		if tt.hasIndex {
			if stmt.Index == nil {
				t.Errorf("input=%q: expected index variable", tt.input)
			} else if stmt.Index.Value != tt.indexName {
				t.Errorf("input=%q: expected index %q, got %q", tt.input, tt.indexName, stmt.Index.Value)
			}
		} else {
			if stmt.Index != nil {
				t.Errorf("input=%q: expected no index variable, got %q", tt.input, stmt.Index.Value)
			}
		}

		if tt.hasValue {
			if stmt.Value == nil {
				t.Errorf("input=%q: expected value variable", tt.input)
			} else if stmt.Value.Value != tt.valueName {
				t.Errorf("input=%q: expected value %q, got %q", tt.input, tt.valueName, stmt.Value.Value)
			}
		} else {
			if stmt.Value != nil {
				t.Errorf("input=%q: expected no value variable, got %q", tt.input, stmt.Value.Value)
			}
		}

		if stmt.Iterable == nil {
			t.Errorf("input=%q: expected iterable", tt.input)
		}

		if stmt.Body == nil {
			t.Errorf("input=%q: expected body", tt.input)
		}
	}
}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("parser had %d errors:", len(errors))
	for _, msg := range errors {
		t.Errorf("  %s", msg)
	}
	t.FailNow()
}
