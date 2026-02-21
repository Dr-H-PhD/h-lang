# H-lang

A statically-typed compiled programming language that transpiles to C.

H-lang blends C and Go syntax with pragmatic design principles. It features manual memory management, type inference, and compiles to portable C code.

## Quick Start

```bash
# Clone and build
git clone https://github.com/Dr-H-PhD/h-lang.git
cd h-lang
go build -o hlc ./cmd/hlc

# Compile and run
./hlc examples/hello.hl
./hello

# Or compile and run in one step
./hlc -run examples/hello.hl
```

## Example

```hlang
# fibonacci.hl

function fibonacci(n int) int {
    if n <= 1 {
        return n;
    }
    return fibonacci(n - 1) + fibonacci(n - 2);
}

function main() {
    # Arrays
    nums := [5]int{10, 20, 30, 40, 50};

    # For-range loop
    for i, v := range nums {
        print(v);
    }

    # Heap allocation with defer cleanup
    data := alloc(int);
    defer free(data);

    *data = fibonacci(10);
    print(*data);  # 55
}
```

## Features

| Feature | Syntax | Description |
|---------|--------|-------------|
| Type inference | `x := 42;` | Infer type from value |
| Constants | `const PI := 3.14;` | Immutable values |
| Explicit types | `var x int = 0;` | Explicit type declaration |
| Pointers | `ptr := &x; *ptr = 10;` | C-style pointers |
| Structs | `struct User { name string; }` | User-defined types |
| Methods | `function (u *User) greet() string` | Methods on structs |
| Arrays | `[5]int{1, 2, 3, 4, 5}` | Fixed-size arrays |
| Slices | `[]int{1, 2, 3}` | Dynamic arrays |
| For loops | `for i := 0; i < 10; i++` | C-style for loops |
| For-range | `for i, v := range arr` | Iterate collections |
| While loops | `while x > 0 { }` | Condition-based loops |
| If/else | `if x > 0 { } else { }` | Conditionals |
| Defer | `defer free(ptr);` | LIFO cleanup at function exit |
| Alloc/Free | `alloc(Type)` / `free(ptr)` | Manual memory management |
| Casting | `(int)x` | C-style type casting |
| Comments | `//`, `/* */`, `#` | Three comment styles |

## Language Specification

| Aspect | Choice |
|--------|--------|
| File extension | `.hl` |
| Entry point | `function main()` |
| Semicolons | Required |
| Visibility | `public` keyword |
| Memory | Manual (`alloc`/`free`) |
| Null | Allowed |
| Target | Transpiles to C |

## Compiler Architecture

```
┌─────────┐     ┌─────────┐     ┌────────┐     ┌─────────┐     ┌───────┐
│ Source  │────▶│  Lexer  │────▶│ Parser │────▶│ Codegen │────▶│  GCC  │
│  (.hl)  │     │         │     │  (AST) │     │   (C)   │     │       │
└─────────┘     └─────────┘     └────────┘     └─────────┘     └───────┘
```

## Project Structure

```
h-lang/
├── cmd/hlc/           # Compiler CLI
├── pkg/
│   ├── ast/           # Abstract Syntax Tree
│   ├── codegen/       # C code generator
│   ├── lexer/         # Tokenizer
│   ├── parser/        # Pratt parser
│   └── version/       # Version info
├── examples/          # Example programs
├── tutorials/         # HTML tutorials
├── scripts/           # Build scripts
└── test/              # Integration tests
```

## CLI Usage

```bash
# Compile to binary
./hlc program.hl

# Compile and run
./hlc -run program.hl

# Output to specific file
./hlc -o myprogram program.hl

# Emit C code only
./hlc -emit-c program.hl

# Show version
./hlc --version

# Show help
./hlc --help
```

## Development

```bash
# Run all tests
go test ./...

# Build and test with version increment
./scripts/build-and-test.sh

# Format code
go fmt ./...

# Vet code
go vet ./...
```

## Tutorials

The `tutorials/` directory contains a complete guide to building the compiler:

| Part | Topic | Description |
|------|-------|-------------|
| 0 | [Introduction](tutorials/00-introduction.html) | Language spec and architecture |
| 1 | [Project Setup](tutorials/01-project-setup.html) | Initialize the Go project |
| 2 | [Lexer](tutorials/02-lexer.html) | Tokenize source code |
| 3 | [AST](tutorials/03-ast.html) | Define syntax tree nodes |
| 4 | [Parser](tutorials/04-parser.html) | Build the AST (Pratt parser) |
| 5 | [Code Generation](tutorials/05-codegen.html) | Generate C code |
| 6 | [CLI](tutorials/06-cli.html) | Build the compiler command |
| 7 | [Testing](tutorials/07-testing.html) | Unit and integration tests |
| 8 | [Defer](tutorials/08-defer.html) | Deferred execution |
| 9 | [Arrays & Slices](tutorials/09-arrays.html) | Collection types |
| 10 | [For-Range Loops](tutorials/10-forrange.html) | Iterating collections |

## Examples

| Example | Description |
|---------|-------------|
| `hello.hl` | Hello World |
| `fibonacci.hl` | Recursion (fibonacci, factorial, power, gcd) |
| `control.hl` | Control flow (if, for, while) |
| `operators.hl` | Arithmetic, comparison, logical operators |
| `pointers.hl` | Pointer operations |
| `structs.hl` | Structs and methods |
| `arrays.hl` | Arrays and slices |
| `defer.hl` | Defer statement |
| `forrange.hl` | For-range loops |
| `demo.hl` | Full feature demonstration |

## Requirements

- Go 1.21+
- GCC or Clang (for compiling generated C)

## Author

Achraf SOLTANI | [achrafsoltani.com](https://achrafsoltani.com)

## License

MIT
