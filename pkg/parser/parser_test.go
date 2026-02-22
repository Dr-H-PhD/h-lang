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

func TestBreakStatement(t *testing.T) {
	input := `break;`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(program.Statements))
	}

	_, ok := program.Statements[0].(*ast.BreakStatement)
	if !ok {
		t.Fatalf("expected BreakStatement, got %T", program.Statements[0])
	}
}

func TestContinueStatement(t *testing.T) {
	input := `continue;`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(program.Statements))
	}

	_, ok := program.Statements[0].(*ast.ContinueStatement)
	if !ok {
		t.Fatalf("expected ContinueStatement, got %T", program.Statements[0])
	}
}

func TestEnumStatement(t *testing.T) {
	input := `enum Color {
    Red,
    Green,
    Blue
}`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.EnumStatement)
	if !ok {
		t.Fatalf("expected EnumStatement, got %T", program.Statements[0])
	}

	if stmt.Public {
		t.Error("expected non-public enum")
	}

	if stmt.Name.Value != "Color" {
		t.Errorf("name: expected 'Color', got %q", stmt.Name.Value)
	}

	if len(stmt.Values) != 3 {
		t.Fatalf("expected 3 values, got %d", len(stmt.Values))
	}

	expectedValues := []string{"Red", "Green", "Blue"}
	for i, expected := range expectedValues {
		if stmt.Values[i].Name.Value != expected {
			t.Errorf("value %d: expected %q, got %q", i, expected, stmt.Values[i].Name.Value)
		}
		if stmt.Values[i].Value != nil {
			t.Errorf("value %d: expected no explicit value", i)
		}
	}
}

func TestEnumStatementWithExplicitValues(t *testing.T) {
	input := `enum Status {
    Pending = 0,
    Active = 1,
    Completed = 2,
    Cancelled = 100
}`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.EnumStatement)

	if len(stmt.Values) != 4 {
		t.Fatalf("expected 4 values, got %d", len(stmt.Values))
	}

	// Check explicit values
	if stmt.Values[0].Value == nil {
		t.Error("expected explicit value for Pending")
	}
	if stmt.Values[3].Name.Value != "Cancelled" {
		t.Errorf("expected 'Cancelled', got %q", stmt.Values[3].Name.Value)
	}
	if stmt.Values[3].Value == nil {
		t.Error("expected explicit value for Cancelled")
	}
}

func TestPublicEnumStatement(t *testing.T) {
	input := `public enum Priority {
    Low,
    Medium,
    High
}`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.EnumStatement)

	if !stmt.Public {
		t.Error("expected public enum")
	}

	if stmt.Name.Value != "Priority" {
		t.Errorf("name: expected 'Priority', got %q", stmt.Name.Value)
	}
}

func TestMapLiteral(t *testing.T) {
	input := `ages := map[string]int{"Alice": 30, "Bob": 25};`

	l := lexer.New(input)
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

	if stmt.Name.Value != "ages" {
		t.Errorf("name: expected 'ages', got %q", stmt.Name.Value)
	}

	ml, ok := stmt.Value.(*ast.MapLiteral)
	if !ok {
		t.Fatalf("expected MapLiteral, got %T", stmt.Value)
	}

	if !ml.Type.IsMap {
		t.Error("expected type.IsMap to be true")
	}

	if ml.Type.KeyType.Name != "string" {
		t.Errorf("key type: expected 'string', got %q", ml.Type.KeyType.Name)
	}

	if ml.Type.ValueType.Name != "int" {
		t.Errorf("value type: expected 'int', got %q", ml.Type.ValueType.Name)
	}

	if len(ml.Pairs) != 2 {
		t.Errorf("expected 2 pairs, got %d", len(ml.Pairs))
	}
}

func TestEmptyMapLiteral(t *testing.T) {
	input := `data := map[string]int{};`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.InferStatement)
	ml := stmt.Value.(*ast.MapLiteral)

	if len(ml.Pairs) != 0 {
		t.Errorf("expected 0 pairs, got %d", len(ml.Pairs))
	}
}

func TestDeleteStatement(t *testing.T) {
	input := `delete(ages, "Alice");`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.DeleteStatement)
	if !ok {
		t.Fatalf("expected DeleteStatement, got %T", program.Statements[0])
	}

	mapIdent, ok := stmt.Map.(*ast.Identifier)
	if !ok {
		t.Fatalf("expected Identifier for map, got %T", stmt.Map)
	}

	if mapIdent.Value != "ages" {
		t.Errorf("map name: expected 'ages', got %q", mapIdent.Value)
	}

	keyLit, ok := stmt.Key.(*ast.StringLiteral)
	if !ok {
		t.Fatalf("expected StringLiteral for key, got %T", stmt.Key)
	}

	if keyLit.Value != "Alice" {
		t.Errorf("key: expected 'Alice', got %q", keyLit.Value)
	}
}

func TestMapTypeAnnotation(t *testing.T) {
	input := `var scores map[string]int;`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.VarStatement)

	if !stmt.Type.IsMap {
		t.Error("expected type.IsMap to be true")
	}

	if stmt.Type.KeyType.Name != "string" {
		t.Errorf("key type: expected 'string', got %q", stmt.Type.KeyType.Name)
	}

	if stmt.Type.ValueType.Name != "int" {
		t.Errorf("value type: expected 'int', got %q", stmt.Type.ValueType.Name)
	}
}

func TestImportStatement(t *testing.T) {
	input := `import "math.hl";`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ImportStatement)
	if !ok {
		t.Fatalf("expected ImportStatement, got %T", program.Statements[0])
	}

	if stmt.Path != "math.hl" {
		t.Errorf("path: expected 'math.hl', got %q", stmt.Path)
	}
}

func TestMultipleImports(t *testing.T) {
	input := `import "math.hl";
import "utils.hl";

function main() {
    print(1);
}`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 3 {
		t.Fatalf("expected 3 statements, got %d", len(program.Statements))
	}

	imp1 := program.Statements[0].(*ast.ImportStatement)
	if imp1.Path != "math.hl" {
		t.Errorf("import 1: expected 'math.hl', got %q", imp1.Path)
	}

	imp2 := program.Statements[1].(*ast.ImportStatement)
	if imp2.Path != "utils.hl" {
		t.Errorf("import 2: expected 'utils.hl', got %q", imp2.Path)
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
