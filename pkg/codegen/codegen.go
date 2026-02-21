package codegen

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/Dr-H-PhD/h-lang/pkg/ast"
)

// Generator generates C code from H-lang AST
type Generator struct {
	output     bytes.Buffer
	indent     int
	structs    map[string]*ast.StructStatement
	functions  map[string]*ast.FunctionStatement
}

// New creates a new code generator
func New() *Generator {
	return &Generator{
		structs:   make(map[string]*ast.StructStatement),
		functions: make(map[string]*ast.FunctionStatement),
	}
}

// Generate produces C code from the AST
func (g *Generator) Generate(program *ast.Program) string {
	// First pass: collect struct and function declarations
	for _, stmt := range program.Statements {
		switch s := stmt.(type) {
		case *ast.StructStatement:
			g.structs[s.Name.Value] = s
		case *ast.FunctionStatement:
			g.functions[s.Name.Value] = s
		}
	}

	// Generate header
	g.writeLine("#include <stdio.h>")
	g.writeLine("#include <stdlib.h>")
	g.writeLine("#include <string.h>")
	g.writeLine("#include <stdbool.h>")
	g.writeLine("")

	// Generate type definitions for strings
	g.writeLine("typedef char* h_string;")
	g.writeLine("")

	// Generate struct forward declarations
	for name := range g.structs {
		g.writeLine(fmt.Sprintf("typedef struct %s %s;", name, name))
	}
	if len(g.structs) > 0 {
		g.writeLine("")
	}

	// Generate struct definitions
	for _, stmt := range program.Statements {
		if s, ok := stmt.(*ast.StructStatement); ok {
			g.generateStruct(s)
		}
	}

	// Generate function forward declarations
	for _, stmt := range program.Statements {
		if s, ok := stmt.(*ast.FunctionStatement); ok {
			g.generateFunctionDeclaration(s)
		}
	}
	if len(g.functions) > 0 {
		g.writeLine("")
	}

	// Generate function implementations
	for _, stmt := range program.Statements {
		if s, ok := stmt.(*ast.FunctionStatement); ok {
			g.generateFunction(s)
		}
	}

	return g.output.String()
}

func (g *Generator) write(s string) {
	g.output.WriteString(s)
}

func (g *Generator) writeLine(s string) {
	g.output.WriteString(strings.Repeat("    ", g.indent))
	g.output.WriteString(s)
	g.output.WriteString("\n")
}

func (g *Generator) generateStruct(s *ast.StructStatement) {
	g.writeLine(fmt.Sprintf("struct %s {", s.Name.Value))
	g.indent++

	for _, field := range s.Fields {
		cType := g.typeToC(field.Type)
		g.writeLine(fmt.Sprintf("%s %s;", cType, field.Name.Value))
	}

	g.indent--
	g.writeLine("};")
	g.writeLine("")
}

func (g *Generator) generateFunctionDeclaration(f *ast.FunctionStatement) {
	returnType := "void"
	if f.ReturnType != nil {
		returnType = g.typeToC(f.ReturnType)
	}

	funcName := f.Name.Value
	if f.Receiver != nil {
		// Method: StructName_methodName
		typeName := f.Receiver.Type.Name
		if f.Receiver.Type.IsPtr {
			typeName = strings.TrimPrefix(typeName, "*")
		}
		funcName = fmt.Sprintf("%s_%s", typeName, f.Name.Value)
	}

	params := g.generateParams(f)
	g.writeLine(fmt.Sprintf("%s %s(%s);", returnType, funcName, params))
}

func (g *Generator) generateFunction(f *ast.FunctionStatement) {
	returnType := "void"
	if f.ReturnType != nil {
		returnType = g.typeToC(f.ReturnType)
	}

	funcName := f.Name.Value
	if f.Receiver != nil {
		typeName := f.Receiver.Type.Name
		if f.Receiver.Type.IsPtr {
			typeName = strings.TrimPrefix(typeName, "*")
		}
		funcName = fmt.Sprintf("%s_%s", typeName, f.Name.Value)
	}

	params := g.generateParams(f)
	g.writeLine(fmt.Sprintf("%s %s(%s) {", returnType, funcName, params))
	g.indent++

	g.generateBlock(f.Body)

	g.indent--
	g.writeLine("}")
	g.writeLine("")
}

func (g *Generator) generateParams(f *ast.FunctionStatement) string {
	var params []string

	// Add receiver as first parameter for methods
	if f.Receiver != nil {
		cType := g.typeToC(f.Receiver.Type)
		params = append(params, fmt.Sprintf("%s %s", cType, f.Receiver.Name.Value))
	}

	for _, p := range f.Parameters {
		cType := g.typeToC(p.Type)
		params = append(params, fmt.Sprintf("%s %s", cType, p.Name.Value))
	}

	if len(params) == 0 {
		return "void"
	}
	return strings.Join(params, ", ")
}

func (g *Generator) generateBlock(block *ast.BlockStatement) {
	for _, stmt := range block.Statements {
		g.generateStatement(stmt)
	}
}

func (g *Generator) generateStatement(stmt ast.Statement) {
	switch s := stmt.(type) {
	case *ast.VarStatement:
		g.generateVarStatement(s)
	case *ast.ConstStatement:
		g.generateConstStatement(s)
	case *ast.InferStatement:
		g.generateInferStatement(s)
	case *ast.ReturnStatement:
		g.generateReturnStatement(s)
	case *ast.IfStatement:
		g.generateIfStatement(s)
	case *ast.ForStatement:
		g.generateForStatement(s)
	case *ast.WhileStatement:
		g.generateWhileStatement(s)
	case *ast.FreeStatement:
		g.generateFreeStatement(s)
	case *ast.ExpressionStatement:
		g.writeLine(g.generateExpression(s.Expression) + ";")
	}
}

func (g *Generator) generateVarStatement(s *ast.VarStatement) {
	cType := g.typeToC(s.Type)
	if s.Value != nil {
		g.writeLine(fmt.Sprintf("%s %s = %s;", cType, s.Name.Value, g.generateExpression(s.Value)))
	} else {
		g.writeLine(fmt.Sprintf("%s %s;", cType, s.Name.Value))
	}
}

func (g *Generator) generateConstStatement(s *ast.ConstStatement) {
	// Infer type from value
	cType := g.inferType(s.Value)
	g.writeLine(fmt.Sprintf("const %s %s = %s;", cType, s.Name.Value, g.generateExpression(s.Value)))
}

func (g *Generator) generateInferStatement(s *ast.InferStatement) {
	cType := g.inferType(s.Value)
	g.writeLine(fmt.Sprintf("%s %s = %s;", cType, s.Name.Value, g.generateExpression(s.Value)))
}

func (g *Generator) generateReturnStatement(s *ast.ReturnStatement) {
	if s.Value != nil {
		g.writeLine(fmt.Sprintf("return %s;", g.generateExpression(s.Value)))
	} else {
		g.writeLine("return;")
	}
}

func (g *Generator) generateIfStatement(s *ast.IfStatement) {
	g.writeLine(fmt.Sprintf("if (%s) {", g.generateExpression(s.Condition)))
	g.indent++
	g.generateBlock(s.Consequence)
	g.indent--

	if s.Alternative != nil {
		g.writeLine("} else {")
		g.indent++
		g.generateBlock(s.Alternative)
		g.indent--
	}
	g.writeLine("}")
}

func (g *Generator) generateForStatement(s *ast.ForStatement) {
	init := ""
	if s.Init != nil {
		init = g.generateStatementInline(s.Init)
	}

	cond := ""
	if s.Condition != nil {
		cond = g.generateExpression(s.Condition)
	}

	post := ""
	if s.Post != nil {
		post = g.generateStatementInline(s.Post)
	}

	g.writeLine(fmt.Sprintf("for (%s; %s; %s) {", init, cond, post))
	g.indent++
	g.generateBlock(s.Body)
	g.indent--
	g.writeLine("}")
}

func (g *Generator) generateWhileStatement(s *ast.WhileStatement) {
	g.writeLine(fmt.Sprintf("while (%s) {", g.generateExpression(s.Condition)))
	g.indent++
	g.generateBlock(s.Body)
	g.indent--
	g.writeLine("}")
}

func (g *Generator) generateFreeStatement(s *ast.FreeStatement) {
	g.writeLine(fmt.Sprintf("free(%s);", g.generateExpression(s.Value)))
}

func (g *Generator) generateStatementInline(stmt ast.Statement) string {
	switch s := stmt.(type) {
	case *ast.InferStatement:
		cType := g.inferType(s.Value)
		return fmt.Sprintf("%s %s = %s", cType, s.Name.Value, g.generateExpression(s.Value))
	case *ast.ExpressionStatement:
		return g.generateExpression(s.Expression)
	}
	return ""
}

func (g *Generator) generateExpression(expr ast.Expression) string {
	switch e := expr.(type) {
	case *ast.Identifier:
		return e.Value
	case *ast.IntegerLiteral:
		return fmt.Sprintf("%d", e.Value)
	case *ast.FloatLiteral:
		return fmt.Sprintf("%f", e.Value)
	case *ast.StringLiteral:
		return fmt.Sprintf("\"%s\"", e.Value)
	case *ast.CharLiteral:
		return fmt.Sprintf("'%c'", e.Value)
	case *ast.BooleanLiteral:
		if e.Value {
			return "true"
		}
		return "false"
	case *ast.NullLiteral:
		return "NULL"
	case *ast.PrefixExpression:
		return fmt.Sprintf("(%s%s)", e.Operator, g.generateExpression(e.Right))
	case *ast.InfixExpression:
		left := g.generateExpression(e.Left)
		right := g.generateExpression(e.Right)
		// Handle string concatenation
		if e.Operator == "+" {
			if g.isStringExpr(e.Left) || g.isStringExpr(e.Right) {
				return fmt.Sprintf("h_string_concat(%s, %s)", left, right)
			}
		}
		return fmt.Sprintf("(%s %s %s)", left, e.Operator, right)
	case *ast.PostfixExpression:
		return fmt.Sprintf("(%s%s)", g.generateExpression(e.Left), e.Operator)
	case *ast.AssignExpression:
		return fmt.Sprintf("(%s %s %s)", g.generateExpression(e.Left), e.Operator, g.generateExpression(e.Value))
	case *ast.CallExpression:
		return g.generateCallExpression(e)
	case *ast.IndexExpression:
		return fmt.Sprintf("%s[%s]", g.generateExpression(e.Left), g.generateExpression(e.Index))
	case *ast.MemberExpression:
		obj := g.generateExpression(e.Object)
		// Use -> for pointers
		if g.isPointerExpr(e.Object) {
			return fmt.Sprintf("%s->%s", obj, e.Member.Value)
		}
		return fmt.Sprintf("%s.%s", obj, e.Member.Value)
	case *ast.CastExpression:
		return fmt.Sprintf("((%s)%s)", g.typeToC(e.TargetType), g.generateExpression(e.Value))
	case *ast.AllocExpression:
		return fmt.Sprintf("(%s*)malloc(sizeof(%s))", e.Type.Name, e.Type.Name)
	case *ast.ArrayLiteral:
		var elements []string
		for _, el := range e.Elements {
			elements = append(elements, g.generateExpression(el))
		}
		return fmt.Sprintf("{%s}", strings.Join(elements, ", "))
	}
	return ""
}

func (g *Generator) generateCallExpression(e *ast.CallExpression) string {
	funcName := g.generateExpression(e.Function)

	// Handle built-in functions
	if funcName == "print" {
		if len(e.Arguments) > 0 {
			arg := e.Arguments[0]
			argStr := g.generateExpression(arg)
			switch arg.(type) {
			case *ast.StringLiteral:
				return fmt.Sprintf("printf(\"%%s\\n\", %s)", argStr)
			case *ast.IntegerLiteral:
				return fmt.Sprintf("printf(\"%%d\\n\", %s)", argStr)
			case *ast.FloatLiteral:
				return fmt.Sprintf("printf(\"%%f\\n\", %s)", argStr)
			default:
				// Default to string
				return fmt.Sprintf("printf(\"%%s\\n\", %s)", argStr)
			}
		}
		return "printf(\"\\n\")"
	}

	// Check if it's a method call (obj.method())
	if member, ok := e.Function.(*ast.MemberExpression); ok {
		// Convert to StructName_method(obj, args)
		obj := g.generateExpression(member.Object)
		method := member.Member.Value

		var args []string
		args = append(args, obj)
		for _, arg := range e.Arguments {
			args = append(args, g.generateExpression(arg))
		}

		// Try to determine struct type
		typeName := g.getExprType(member.Object)
		return fmt.Sprintf("%s_%s(%s)", typeName, method, strings.Join(args, ", "))
	}

	var args []string
	for _, arg := range e.Arguments {
		args = append(args, g.generateExpression(arg))
	}
	return fmt.Sprintf("%s(%s)", funcName, strings.Join(args, ", "))
}

func (g *Generator) typeToC(t *ast.TypeAnnotation) string {
	if t == nil {
		return "void"
	}

	var cType string
	switch t.Name {
	case "int":
		cType = "int"
	case "float":
		cType = "double"
	case "string":
		cType = "h_string"
	case "char":
		cType = "char"
	case "bool":
		cType = "bool"
	case "void":
		cType = "void"
	default:
		// User-defined type (struct)
		cType = t.Name
	}

	if t.ArrayLen == -1 {
		cType = cType + "*" // Slice becomes pointer
	} else if t.ArrayLen > 0 {
		cType = fmt.Sprintf("%s[%d]", cType, t.ArrayLen)
	}

	if t.IsPtr {
		cType = cType + "*"
	}

	return cType
}

func (g *Generator) inferType(expr ast.Expression) string {
	switch e := expr.(type) {
	case *ast.IntegerLiteral:
		return "int"
	case *ast.FloatLiteral:
		return "double"
	case *ast.StringLiteral:
		return "h_string"
	case *ast.CharLiteral:
		return "char"
	case *ast.BooleanLiteral:
		return "bool"
	case *ast.NullLiteral:
		return "void*"
	case *ast.AllocExpression:
		return e.Type.Name + "*"
	case *ast.PrefixExpression:
		if e.Operator == "&" {
			return g.inferType(e.Right) + "*"
		}
		if e.Operator == "*" {
			inner := g.inferType(e.Right)
			return strings.TrimSuffix(inner, "*")
		}
		return g.inferType(e.Right)
	case *ast.InfixExpression:
		return g.inferType(e.Left)
	case *ast.CallExpression:
		// Look up function return type
		if ident, ok := e.Function.(*ast.Identifier); ok {
			if fn, exists := g.functions[ident.Value]; exists && fn.ReturnType != nil {
				return g.typeToC(fn.ReturnType)
			}
		}
		return "int"
	}
	return "int"
}

func (g *Generator) isStringExpr(expr ast.Expression) bool {
	_, ok := expr.(*ast.StringLiteral)
	return ok
}

func (g *Generator) isPointerExpr(expr ast.Expression) bool {
	switch e := expr.(type) {
	case *ast.Identifier:
		// TODO: look up in symbol table
		return false
	case *ast.AllocExpression:
		return true
	case *ast.PrefixExpression:
		return e.Operator == "&"
	}
	return false
}

func (g *Generator) getExprType(expr ast.Expression) string {
	switch e := expr.(type) {
	case *ast.Identifier:
		// TODO: proper symbol table lookup
		return "Unknown"
	case *ast.AllocExpression:
		return e.Type.Name
	}
	return "Unknown"
}
