package test

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Dr-H-PhD/h-lang/pkg/codegen"
	"github.com/Dr-H-PhD/h-lang/pkg/lexer"
	"github.com/Dr-H-PhD/h-lang/pkg/parser"
)

// TestPipeline_Lexer validates that source code can be tokenized
func TestPipeline_Lexer(t *testing.T) {
	source := `
function main() {
    x := 42;
    print(x);
}
`
	l := lexer.New(source)

	// Should be able to tokenize without errors
	tokenCount := 0
	for {
		tok := l.NextToken()
		tokenCount++
		if tok.Type == lexer.EOF {
			break
		}
		if tok.Type == lexer.ILLEGAL {
			t.Fatalf("illegal token: %v", tok)
		}
	}

	if tokenCount < 10 {
		t.Errorf("expected at least 10 tokens, got %d", tokenCount)
	}
}

// TestPipeline_Parser validates that source code can be parsed
func TestPipeline_Parser(t *testing.T) {
	source := `
function main() {
    x := 42;
    print(x);
}
`
	l := lexer.New(source)
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		t.Fatalf("parser errors: %v", p.Errors())
	}

	if len(program.Statements) != 1 {
		t.Errorf("expected 1 statement, got %d", len(program.Statements))
	}
}

// TestPipeline_CodeGen validates that AST can be converted to C
func TestPipeline_CodeGen(t *testing.T) {
	source := `
function main() {
    x := 42;
}
`
	l := lexer.New(source)
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		t.Fatalf("parser errors: %v", p.Errors())
	}

	g := codegen.New()
	cCode := g.Generate(program)

	if !strings.Contains(cCode, "int x = 42") {
		t.Errorf("expected C code to contain variable declaration")
	}
}

// TestCompilation tests that H-lang code compiles to valid C
func TestCompilation_HelloWorld(t *testing.T) {
	source := `
function main() {
    x := 42;
}
`
	err := compileOnly(t, source)
	if err != nil {
		t.Fatalf("compilation failed: %v", err)
	}
}

func TestCompilation_Arithmetic(t *testing.T) {
	source := `
function main() {
    x := 10 + 20 * 2;
}
`
	err := compileOnly(t, source)
	if err != nil {
		t.Fatalf("failed: %v", err)
	}
}

func TestCompilation_Variables(t *testing.T) {
	source := `
function main() {
    x := 42;
    const PI := 3.14;
    var count int = 0;
}
`
	err := compileOnly(t, source)
	if err != nil {
		t.Fatalf("failed: %v", err)
	}
}

func TestCompilation_IfStatement(t *testing.T) {
	source := `
function main() {
    x := 5;
    if x > 0 {
        x = 1;
    } else {
        x = 0;
    }
}
`
	err := compileOnly(t, source)
	if err != nil {
		t.Fatalf("failed: %v", err)
	}
}

func TestCompilation_ForLoop(t *testing.T) {
	source := `
function main() {
    sum := 0;
    for i := 1; i <= 5; i++ {
        sum = sum + i;
    }
}
`
	err := compileOnly(t, source)
	if err != nil {
		t.Fatalf("failed: %v", err)
	}
}

func TestCompilation_WhileLoop(t *testing.T) {
	source := `
function main() {
    x := 0;
    while x < 3 {
        x++;
    }
}
`
	err := compileOnly(t, source)
	if err != nil {
		t.Fatalf("failed: %v", err)
	}
}

func TestCompilation_FunctionCall(t *testing.T) {
	source := `
function add(a int, b int) int {
    return a + b;
}

function main() {
    result := add(10, 20);
}
`
	err := compileOnly(t, source)
	if err != nil {
		t.Fatalf("failed: %v", err)
	}
}

func TestCompilation_Recursion(t *testing.T) {
	source := `
function factorial(n int) int {
    if n <= 1 {
        return 1;
    }
    return n * factorial(n - 1);
}

function main() {
    result := factorial(5);
}
`
	err := compileOnly(t, source)
	if err != nil {
		t.Fatalf("failed: %v", err)
	}
}

func TestCompilation_Pointers(t *testing.T) {
	source := `
function main() {
    x := 42;
    ptr := &x;
    *ptr = 100;
}
`
	err := compileOnly(t, source)
	if err != nil {
		t.Fatalf("failed: %v", err)
	}
}

func TestCompilation_Cast(t *testing.T) {
	source := `
function main() {
    f := 3.7;
    i := (int)f;
}
`
	err := compileOnly(t, source)
	if err != nil {
		t.Fatalf("failed: %v", err)
	}
}

func TestCompilation_CompoundAssignment(t *testing.T) {
	source := `
function main() {
    x := 10;
    x += 5;
    x -= 2;
    x *= 3;
    x /= 2;
}
`
	err := compileOnly(t, source)
	if err != nil {
		t.Fatalf("failed: %v", err)
	}
}

func TestCompilation_BooleanLogic(t *testing.T) {
	source := `
function main() {
    x := 5;
    y := 10;
    if x < y && y > 0 {
        x = 1;
    }
}
`
	err := compileOnly(t, source)
	if err != nil {
		t.Fatalf("failed: %v", err)
	}
}

func TestCompilation_ForRange(t *testing.T) {
	source := `
function main() {
    arr := [5]int{1, 2, 3, 4, 5};
    sum := 0;
    for i, v := range arr {
        sum = sum + v;
    }
}
`
	err := compileOnly(t, source)
	if err != nil {
		t.Fatalf("failed: %v", err)
	}
}

func TestCompilation_ForRangeIndexOnly(t *testing.T) {
	source := `
function main() {
    arr := [3]int{10, 20, 30};
    for i := range arr {
        print(i);
    }
}
`
	err := compileOnly(t, source)
	if err != nil {
		t.Fatalf("failed: %v", err)
	}
}

func TestCompilation_Break(t *testing.T) {
	source := `
function main() {
    for i := 0; i < 10; i++ {
        if i == 5 {
            break;
        }
    }
}
`
	err := compileOnly(t, source)
	if err != nil {
		t.Fatalf("failed: %v", err)
	}
}

func TestCompilation_Continue(t *testing.T) {
	source := `
function main() {
    sum := 0;
    for i := 0; i < 10; i++ {
        if i == 5 {
            continue;
        }
        sum = sum + i;
    }
}
`
	err := compileOnly(t, source)
	if err != nil {
		t.Fatalf("failed: %v", err)
	}
}

func TestCompilation_Enum(t *testing.T) {
	source := `
enum Color {
    Red,
    Green,
    Blue
}

function main() {
    c := Color_Red;
    if c == Color_Green {
        print(1);
    }
}
`
	err := compileOnly(t, source)
	if err != nil {
		t.Fatalf("failed: %v", err)
	}
}

func TestCompilation_EnumWithValues(t *testing.T) {
	source := `
enum Status {
    Pending = 0,
    Active = 1,
    Completed = 2,
    Cancelled = 100
}

function main() {
    s := Status_Active;
    if s == 1 {
        print(s);
    }
}
`
	err := compileOnly(t, source)
	if err != nil {
		t.Fatalf("failed: %v", err)
	}
}

func TestCompilation_WhileBreak(t *testing.T) {
	source := `
function main() {
    x := 0;
    while true {
        x++;
        if x >= 10 {
            break;
        }
    }
}
`
	err := compileOnly(t, source)
	if err != nil {
		t.Fatalf("failed: %v", err)
	}
}

func TestCompilation_MapLiteral(t *testing.T) {
	source := `
function main() {
    ages := map[string]int{"Alice": 30, "Bob": 25};
    x := ages["Alice"];
    print(x);
}
`
	err := compileOnly(t, source)
	if err != nil {
		t.Fatalf("failed: %v", err)
	}
}

func TestCompilation_MapEmpty(t *testing.T) {
	source := `
function main() {
    data := map[string]int{};
    data["key"] = 42;
    print(data["key"]);
}
`
	err := compileOnly(t, source)
	if err != nil {
		t.Fatalf("failed: %v", err)
	}
}

func TestCompilation_MapDelete(t *testing.T) {
	source := `
function main() {
    ages := map[string]int{"Alice": 30};
    delete(ages, "Alice");
    print(len(ages));
}
`
	err := compileOnly(t, source)
	if err != nil {
		t.Fatalf("failed: %v", err)
	}
}

func TestCompilation_MapFree(t *testing.T) {
	source := `
function main() {
    ages := map[string]int{"Alice": 30, "Bob": 25};
    print(len(ages));
    free(ages);
}
`
	err := compileOnly(t, source)
	if err != nil {
		t.Fatalf("failed: %v", err)
	}
}

// compileOnly compiles H-lang code to C and verifies C compilation succeeds
func compileOnly(t *testing.T, source string) error {
	t.Helper()

	// Skip if no C compiler available
	compiler := findCompiler()
	if compiler == "" {
		t.Skip("no C compiler found (gcc, clang)")
	}

	// Parse
	l := lexer.New(source)
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		return &compileError{errors: p.Errors()}
	}

	// Generate C code
	g := codegen.New()
	cCode := g.Generate(program)

	// Create temp directory
	tmpDir, err := os.MkdirTemp("", "hlang-test-*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpDir)

	// Write C file
	cFile := filepath.Join(tmpDir, "test.c")
	if err := os.WriteFile(cFile, []byte(cCode), 0644); err != nil {
		return err
	}

	// Compile (but don't run)
	binary := filepath.Join(tmpDir, "test")
	cmd := exec.Command(compiler, "-o", binary, cFile)
	if output, err := cmd.CombinedOutput(); err != nil {
		return &compileError{
			errors: []string{
				"C compilation failed: " + err.Error(),
				"Output: " + string(output),
				"Generated C:\n" + cCode,
			},
		}
	}

	return nil
}

// Helper function to compile and run H-lang code
func compileAndRun(t *testing.T, source string) (string, error) {
	t.Helper()

	// Skip if no C compiler available
	compiler := findCompiler()
	if compiler == "" {
		t.Skip("no C compiler found (gcc, clang)")
	}

	// Parse
	l := lexer.New(source)
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		return "", &compileError{errors: p.Errors()}
	}

	// Generate C code
	g := codegen.New()
	cCode := g.Generate(program)

	// Create temp directory
	tmpDir, err := os.MkdirTemp("", "hlang-test-*")
	if err != nil {
		return "", err
	}
	defer os.RemoveAll(tmpDir)

	// Write C file
	cFile := filepath.Join(tmpDir, "test.c")
	if err := os.WriteFile(cFile, []byte(cCode), 0644); err != nil {
		return "", err
	}

	// Compile
	binary := filepath.Join(tmpDir, "test")
	cmd := exec.Command(compiler, "-o", binary, cFile)
	if output, err := cmd.CombinedOutput(); err != nil {
		return "", &compileError{
			errors: []string{
				"C compilation failed: " + err.Error(),
				"Output: " + string(output),
				"Generated C:\n" + cCode,
			},
		}
	}

	// Run
	runCmd := exec.Command(binary)
	output, err := runCmd.CombinedOutput()
	if err != nil {
		return string(output), err
	}

	return string(output), nil
}

func findCompiler() string {
	compilers := []string{"gcc", "clang", "cc"}
	for _, c := range compilers {
		if _, err := exec.LookPath(c); err == nil {
			return c
		}
	}
	return ""
}

type compileError struct {
	errors []string
}

func (e *compileError) Error() string {
	return strings.Join(e.errors, "\n")
}
