package codegen

import (
	"strings"
	"testing"

	"github.com/Dr-H-PhD/h-lang/pkg/lexer"
	"github.com/Dr-H-PhD/h-lang/pkg/parser"
)

func TestGenerate_HelloWorld(t *testing.T) {
	input := `function main() {
    print("Hello, H-lang!");
}`

	code := compile(t, input)

	assertContains(t, code, "#include <stdio.h>")
	assertContains(t, code, "int main(void)")
	assertContains(t, code, "printf")
	assertContains(t, code, "Hello, H-lang!")
}

func TestGenerate_Variables(t *testing.T) {
	input := `function main() {
    x := 42;
    const PI := 3.14;
    var count int = 0;
}`

	code := compile(t, input)

	assertContains(t, code, "int x = 42;")
	assertContains(t, code, "const double PI = 3.14")
	assertContains(t, code, "int count = 0;")
}

func TestGenerate_Struct(t *testing.T) {
	input := `public struct User {
    public name string;
    public age int;
}`

	code := compile(t, input)

	assertContains(t, code, "typedef struct User User;")
	assertContains(t, code, "struct User {")
	assertContains(t, code, "h_string name;")
	assertContains(t, code, "int age;")
}

func TestGenerate_Function(t *testing.T) {
	input := `function add(a int, b int) int {
    return a + b;
}`

	code := compile(t, input)

	assertContains(t, code, "int add(int a, int b)")
	assertContains(t, code, "return (a + b);")
}

func TestGenerate_Method(t *testing.T) {
	input := `public struct User {
    public name string;
}

public function (u *User) greet() string {
    return "hello";
}`

	code := compile(t, input)

	// Method should be translated to User_greet with receiver as first param
	assertContains(t, code, "h_string User_greet(User* u)")
}

func TestGenerate_IfStatement(t *testing.T) {
	input := `function main() {
    x := 5;
    if x > 0 {
        print("positive");
    }
}`

	code := compile(t, input)

	assertContains(t, code, "if ((x > 0))")
}

func TestGenerate_IfElseStatement(t *testing.T) {
	input := `function main() {
    x := 5;
    if x > 0 {
        print("positive");
    } else {
        print("non-positive");
    }
}`

	code := compile(t, input)

	assertContains(t, code, "if ((x > 0))")
	assertContains(t, code, "} else {")
}

func TestGenerate_ForLoop(t *testing.T) {
	input := `function main() {
    for i := 0; i < 10; i++ {
        print(i);
    }
}`

	code := compile(t, input)

	assertContains(t, code, "for (int i = 0; (i < 10); (i++))")
}

func TestGenerate_WhileLoop(t *testing.T) {
	input := `function main() {
    x := 0;
    while x < 10 {
        x++;
    }
}`

	code := compile(t, input)

	assertContains(t, code, "while ((x < 10))")
}

func TestGenerate_Alloc(t *testing.T) {
	input := `public struct User {
    public name string;
}

function main() {
    user := alloc(User);
}`

	code := compile(t, input)

	assertContains(t, code, "(User*)malloc(sizeof(User))")
}

func TestGenerate_Free(t *testing.T) {
	input := `public struct User {
    public name string;
}

function main() {
    user := alloc(User);
    free(user);
}`

	code := compile(t, input)

	assertContains(t, code, "free(user);")
}

func TestGenerate_Defer(t *testing.T) {
	input := `public struct Data {
    public value int;
}

function process() int {
    x := alloc(Data);
    defer free(x);
    return 1;
}`

	code := compile(t, input)

	// Defer should emit free before return
	// The return value should be saved first
	assertContains(t, code, "__ret_val")
	assertContains(t, code, "free(x);")
	assertContains(t, code, "return __ret_val;")
}

func TestGenerate_DeferLIFO(t *testing.T) {
	input := `function test() {
    defer print("first");
    defer print("second");
    defer print("third");
}`

	code := compile(t, input)

	// Find positions of print calls
	firstPos := strings.Index(code, `printf("%s\n", "first")`)
	secondPos := strings.Index(code, `printf("%s\n", "second")`)
	thirdPos := strings.Index(code, `printf("%s\n", "third")`)

	// Should be in LIFO order (third, second, first)
	if thirdPos >= secondPos || secondPos >= firstPos {
		t.Errorf("defer should emit in LIFO order: third=%d, second=%d, first=%d",
			thirdPos, secondPos, firstPos)
	}
}

func TestGenerate_Cast(t *testing.T) {
	input := `function main() {
    x := 3.7;
    y := (int)x;
}`

	code := compile(t, input)

	assertContains(t, code, "((int)x)")
}

func TestGenerate_Pointer(t *testing.T) {
	input := `function main() {
    x := 42;
    ptr := &x;
    y := *ptr;
}`

	code := compile(t, input)

	assertContains(t, code, "(&x)")
	assertContains(t, code, "(*ptr)")
}

func TestGenerate_CompoundAssignment(t *testing.T) {
	input := `function main() {
    x := 10;
    x += 5;
    x -= 2;
    x *= 3;
    x /= 2;
}`

	code := compile(t, input)

	assertContains(t, code, "(x += 5)")
	assertContains(t, code, "(x -= 2)")
	assertContains(t, code, "(x *= 3)")
	assertContains(t, code, "(x /= 2)")
}

func TestGenerate_BooleanOperators(t *testing.T) {
	input := `function main() {
    x := true && false;
    y := true || false;
    z := !true;
}`

	code := compile(t, input)

	assertContains(t, code, "(true && false)")
	assertContains(t, code, "(true || false)")
	assertContains(t, code, "(!true)")
}

func TestGenerate_ComparisonOperators(t *testing.T) {
	input := `function main() {
    a := 1 == 2;
    b := 1 != 2;
    c := 1 < 2;
    d := 1 <= 2;
    e := 1 > 2;
    f := 1 >= 2;
}`

	code := compile(t, input)

	assertContains(t, code, "(1 == 2)")
	assertContains(t, code, "(1 != 2)")
	assertContains(t, code, "(1 < 2)")
	assertContains(t, code, "(1 <= 2)")
	assertContains(t, code, "(1 > 2)")
	assertContains(t, code, "(1 >= 2)")
}

func TestGenerate_NullCheck(t *testing.T) {
	input := `public struct User {
    public name string;
}

function main() {
    user := alloc(User);
    if user != null {
        print("not null");
    }
}`

	code := compile(t, input)

	assertContains(t, code, "(user != NULL)")
}

func TestGenerate_IncDec(t *testing.T) {
	input := `function main() {
    x := 0;
    x++;
    x--;
}`

	code := compile(t, input)

	assertContains(t, code, "(x++)")
	assertContains(t, code, "(x--)")
}

func TestGenerate_Headers(t *testing.T) {
	input := `function main() {}`

	code := compile(t, input)

	assertContains(t, code, "#include <stdio.h>")
	assertContains(t, code, "#include <stdlib.h>")
	assertContains(t, code, "#include <string.h>")
	assertContains(t, code, "#include <stdbool.h>")
	assertContains(t, code, "typedef char* h_string;")
}

func TestGenerate_FunctionForwardDeclaration(t *testing.T) {
	input := `function main() {
    x := add(1, 2);
}

function add(a int, b int) int {
    return a + b;
}`

	code := compile(t, input)

	// Should have forward declaration before main
	declPos := strings.Index(code, "int add(int a, int b);")
	mainPos := strings.Index(code, "int main(void) {")

	if declPos == -1 || mainPos == -1 {
		t.Fatal("missing forward declaration or main")
	}

	if declPos > mainPos {
		t.Error("forward declaration should come before main")
	}
}

func TestGenerate_FixedArray(t *testing.T) {
	input := `function main() {
    arr := [5]int{1, 2, 3, 4, 5};
}`

	code := compile(t, input)

	assertContains(t, code, "int arr[5] = {1, 2, 3, 4, 5}")
}

func TestGenerate_SliceLiteral(t *testing.T) {
	input := `function main() {
    nums := []int{10, 20, 30};
}`

	code := compile(t, input)

	assertContains(t, code, "int nums[] = {10, 20, 30}")
}

func TestGenerate_ArrayIndexing(t *testing.T) {
	input := `function main() {
    arr := [3]int{1, 2, 3};
    x := arr[0];
    arr[1] = 100;
}`

	code := compile(t, input)

	assertContains(t, code, "arr[0]")
	assertContains(t, code, "arr[1]")
}

func TestGenerate_LenFunction(t *testing.T) {
	input := `function main() {
    arr := [5]int{1, 2, 3, 4, 5};
    size := len(arr);
}`

	code := compile(t, input)

	assertContains(t, code, "sizeof(arr)/sizeof(arr[0])")
}

func TestGenerate_MakeSlice(t *testing.T) {
	input := `function main() {
    buf := make([]int, 10);
}`

	code := compile(t, input)

	assertContains(t, code, "int* buf")
	assertContains(t, code, "calloc(10, sizeof(int))")
}

func TestGenerate_ForRangeLoop(t *testing.T) {
	input := `function main() {
    arr := [5]int{1, 2, 3, 4, 5};
    for i, v := range arr {
        print(v);
    }
}`

	code := compile(t, input)

	assertContains(t, code, "for (int i = 0;")
	assertContains(t, code, "sizeof(arr)/sizeof(arr[0])")
	assertContains(t, code, "int v = arr[i];")
}

func TestGenerate_ForRangeIndexOnly(t *testing.T) {
	input := `function main() {
    arr := [3]int{10, 20, 30};
    for i := range arr {
        print(i);
    }
}`

	code := compile(t, input)

	assertContains(t, code, "for (int i = 0;")
	assertContains(t, code, "sizeof(arr)/sizeof(arr[0])")
}

func TestGenerate_BreakStatement(t *testing.T) {
	input := `function main() {
    for i := 0; i < 10; i++ {
        if i == 5 {
            break;
        }
    }
}`

	code := compile(t, input)

	assertContains(t, code, "break;")
}

func TestGenerate_ContinueStatement(t *testing.T) {
	input := `function main() {
    for i := 0; i < 10; i++ {
        if i == 5 {
            continue;
        }
    }
}`

	code := compile(t, input)

	assertContains(t, code, "continue;")
}

func TestGenerate_Enum(t *testing.T) {
	input := `enum Color {
    Red,
    Green,
    Blue
}`

	code := compile(t, input)

	assertContains(t, code, "typedef enum {")
	assertContains(t, code, "Color_Red")
	assertContains(t, code, "Color_Green")
	assertContains(t, code, "Color_Blue")
	assertContains(t, code, "} Color;")
}

func TestGenerate_EnumWithValues(t *testing.T) {
	input := `enum Status {
    Pending = 0,
    Active = 1,
    Cancelled = 100
}`

	code := compile(t, input)

	assertContains(t, code, "Status_Pending = 0")
	assertContains(t, code, "Status_Active = 1")
	assertContains(t, code, "Status_Cancelled = 100")
}

func TestGenerate_EnumUsage(t *testing.T) {
	input := `enum Color {
    Red,
    Green,
    Blue
}

function main() {
    c := Color_Red;
    if c == Color_Red {
        print(1);
    }
}`

	code := compile(t, input)

	assertContains(t, code, "int c = Color_Red;")
	assertContains(t, code, "(c == Color_Red)")
}

func TestGenerate_MapLiteral(t *testing.T) {
	input := `function main() {
    ages := map[string]int{"Alice": 30, "Bob": 25};
    print(ages["Alice"]);
}`

	code := compile(t, input)

	// Check map helpers are generated
	assertContains(t, code, "h_map* ages = h_map_new();")
	assertContains(t, code, "h_map_set(ages,")
	assertContains(t, code, "h_map_get(ages,")
}

func TestGenerate_MapAssignment(t *testing.T) {
	input := `function main() {
    ages := map[string]int{};
    ages["Charlie"] = 35;
}`

	code := compile(t, input)

	assertContains(t, code, "h_map* ages = h_map_new();")
	assertContains(t, code, "h_map_set(ages,")
}

func TestGenerate_MapDelete(t *testing.T) {
	input := `function main() {
    ages := map[string]int{"Alice": 30};
    delete(ages, "Alice");
}`

	code := compile(t, input)

	assertContains(t, code, "h_map_delete(ages,")
}

func TestGenerate_MapLen(t *testing.T) {
	input := `function main() {
    ages := map[string]int{"Alice": 30, "Bob": 25};
    count := len(ages);
}`

	code := compile(t, input)

	assertContains(t, code, "h_map_len(ages)")
}

func TestGenerate_MapFree(t *testing.T) {
	input := `function main() {
    ages := map[string]int{};
    free(ages);
}`

	code := compile(t, input)

	assertContains(t, code, "h_map_free(ages);")
}

func TestGenerate_MapHelpers(t *testing.T) {
	input := `function main() {
    m := map[string]int{};
}`

	code := compile(t, input)

	// Verify map helper functions are generated
	assertContains(t, code, "typedef struct h_map_entry")
	assertContains(t, code, "typedef struct {")
	assertContains(t, code, "h_map* h_map_new()")
	assertContains(t, code, "void h_map_set(h_map* m,")
	assertContains(t, code, "void* h_map_get(h_map* m,")
	assertContains(t, code, "void h_map_delete(h_map* m,")
	assertContains(t, code, "int h_map_len(h_map* m)")
	assertContains(t, code, "void h_map_free(h_map* m)")
}

// Helper functions

func compile(t *testing.T, input string) string {
	t.Helper()

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		t.Fatalf("parser errors: %v", p.Errors())
	}

	g := New()
	return g.Generate(program)
}

func assertContains(t *testing.T, code, substr string) {
	t.Helper()
	if !strings.Contains(code, substr) {
		t.Errorf("expected code to contain %q\n\nGenerated code:\n%s", substr, code)
	}
}
