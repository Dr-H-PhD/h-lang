# H-lang

A statically-typed, compiled programming language that blends C and Go syntax with PHP pragmatism. H-lang transpiles to C, leveraging existing compilers for optimization and portability.

## Features

- **Static types with inference** — Type safety without verbosity
- **Manual memory management** — Full control with `alloc`/`free`
- **C/Go syntax blend** — Familiar to systems programmers
- **Transpiles to C** — Portable and optimized
- **Three comment styles** — `//`, `/* */`, and `#`

## Quick Start

```bash
# Build the compiler
go build -o hlc ./cmd/hlc

# Compile and run
./hlc -run examples/hello.hl
```

## Syntax Example

```h
# hello.hl

public struct User {
    public name string;
    public age int;
}

public function (u *User) greet() string {
    return "Hello, " + u.name;
}

function main() {
    x := 42;
    const PI := 3.14159;

    user := alloc(User);
    user.name = "Achraf";
    user.age = 30;

    if user != null {
        print(user.greet());
    }

    for i := 0; i < 10; i++ {
        print(i);
    }

    free(user);
}
```

## Language Specification

| Aspect | Choice |
|--------|--------|
| File extension | `.hl` |
| Entry point | `function main()` |
| Semicolons | Required `;` |
| Comments | `//`, `/* */`, `#` |
| Visibility | `public` keyword |
| Functions | `function` keyword |
| Types | Static with strong inference |
| Mutability | All mutable, `const` for immutable |
| Pointers | C-style `&` / `*` |
| Null | Allowed |
| Casting | C-style `(int)x` |
| Memory | Manual (`alloc`/`free`) |
| Target | Transpile to C |

## CLI Usage

```
hlc [options] <file.hl>

Options:
  -o <file>     Output file name
  -emit-c       Emit C code instead of compiling
  -run          Compile and run immediately
  -version      Print version
  -help         Print help
```

## Project Structure

```
h-lang/
├── cmd/hlc/          # Compiler CLI
├── pkg/
│   ├── lexer/        # Tokenizer
│   ├── parser/       # AST builder
│   ├── ast/          # AST types
│   └── codegen/      # C code generator
├── examples/         # Example programs
├── tutorials/        # HTML tutorial series
└── go.mod
```

## Tutorials

See the [tutorials](tutorials/index.html) directory for a complete guide to building H-lang from scratch.

## Author

Achraf SOLTANI | [achrafsoltani.com](https://achrafsoltani.com)

## License

MIT
