package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/Dr-H-PhD/h-lang/pkg/codegen"
	"github.com/Dr-H-PhD/h-lang/pkg/lexer"
	"github.com/Dr-H-PhD/h-lang/pkg/parser"
	"github.com/Dr-H-PhD/h-lang/pkg/version"
)

func main() {
	// Flags
	outputFlag := flag.String("o", "", "Output file name")
	emitC := flag.Bool("emit-c", false, "Emit C code instead of compiling")
	runFlag := flag.Bool("run", false, "Compile and run immediately")
	versionFlag := flag.Bool("version", false, "Print version")
	helpFlag := flag.Bool("help", false, "Print help")

	flag.Parse()

	if *versionFlag {
		fmt.Println(version.Full())
		os.Exit(0)
	}

	if *helpFlag || flag.NArg() == 0 {
		printUsage()
		os.Exit(0)
	}

	inputFile := flag.Arg(0)

	// Validate input file
	if !strings.HasSuffix(inputFile, ".hl") {
		fmt.Fprintf(os.Stderr, "Error: input file must have .hl extension\n")
		os.Exit(1)
	}

	// Read input file
	source, err := os.ReadFile(inputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	// Compile
	cCode, errors := compile(string(source))
	if len(errors) > 0 {
		fmt.Fprintf(os.Stderr, "Compilation errors:\n")
		for _, e := range errors {
			fmt.Fprintf(os.Stderr, "  %s\n", e)
		}
		os.Exit(1)
	}

	// Determine output names
	baseName := strings.TrimSuffix(filepath.Base(inputFile), ".hl")
	cFileName := baseName + ".c"
	outputName := baseName
	if *outputFlag != "" {
		outputName = *outputFlag
		if *emitC {
			cFileName = outputName
		}
	}

	if *emitC {
		// Just emit C code
		err := os.WriteFile(cFileName, []byte(cCode), 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error writing C file: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Generated: %s\n", cFileName)
		return
	}

	// Write temporary C file
	tmpDir, err := os.MkdirTemp("", "hlc-*")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating temp directory: %v\n", err)
		os.Exit(1)
	}
	defer os.RemoveAll(tmpDir)

	tmpCFile := filepath.Join(tmpDir, cFileName)
	err = os.WriteFile(tmpCFile, []byte(cCode), 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing temp C file: %v\n", err)
		os.Exit(1)
	}

	// Compile with gcc/clang
	compiler := findCompiler()
	if compiler == "" {
		fmt.Fprintf(os.Stderr, "Error: no C compiler found (tried gcc, clang)\n")
		os.Exit(1)
	}

	cmd := exec.Command(compiler, "-o", outputName, tmpCFile)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error compiling C code: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Compiled: %s\n", outputName)

	// Run if requested
	if *runFlag {
		fmt.Println("---")
		runCmd := exec.Command("./" + outputName)
		runCmd.Stdout = os.Stdout
		runCmd.Stderr = os.Stderr
		runCmd.Stdin = os.Stdin
		runCmd.Run()
	}
}

func compile(source string) (string, []string) {
	// Lexer
	l := lexer.New(source)

	// Parser
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		return "", p.Errors()
	}

	// Code generation
	g := codegen.New()
	cCode := g.Generate(program)

	return cCode, nil
}

func findCompiler() string {
	compilers := []string{"gcc", "clang", "cc"}
	for _, c := range compilers {
		path, err := exec.LookPath(c)
		if err == nil {
			return path
		}
	}
	return ""
}

func printUsage() {
	fmt.Println("H-lang Compiler (hlc)")
	fmt.Println()
	fmt.Println("Usage: hlc [options] <file.hl>")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  -o <file>     Output file name")
	fmt.Println("  -emit-c       Emit C code instead of compiling")
	fmt.Println("  -run          Compile and run immediately")
	fmt.Println("  -version      Print version")
	fmt.Println("  -help         Print this help")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  hlc hello.hl              Compile hello.hl to ./hello")
	fmt.Println("  hlc -o myapp hello.hl     Compile to ./myapp")
	fmt.Println("  hlc -emit-c hello.hl      Generate hello.c")
	fmt.Println("  hlc -run hello.hl         Compile and run")
}
