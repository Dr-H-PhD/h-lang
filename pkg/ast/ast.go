package ast

import (
	"bytes"
	"strings"

	"github.com/Dr-H-PhD/h-lang/pkg/lexer"
)

// Node is the base interface for all AST nodes
type Node interface {
	TokenLiteral() string
	String() string
}

// Statement represents a statement node
type Statement interface {
	Node
	statementNode()
}

// Expression represents an expression node
type Expression interface {
	Node
	expressionNode()
}

// Program is the root node of the AST
type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}
	return ""
}

func (p *Program) String() string {
	var out bytes.Buffer
	for _, s := range p.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

// Identifier represents a variable name
type Identifier struct {
	Token lexer.Token
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String() string       { return i.Value }

// IntegerLiteral represents an integer
type IntegerLiteral struct {
	Token lexer.Token
	Value int64
}

func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *IntegerLiteral) String() string       { return il.Token.Literal }

// FloatLiteral represents a float
type FloatLiteral struct {
	Token lexer.Token
	Value float64
}

func (fl *FloatLiteral) expressionNode()      {}
func (fl *FloatLiteral) TokenLiteral() string { return fl.Token.Literal }
func (fl *FloatLiteral) String() string       { return fl.Token.Literal }

// StringLiteral represents a string
type StringLiteral struct {
	Token lexer.Token
	Value string
}

func (sl *StringLiteral) expressionNode()      {}
func (sl *StringLiteral) TokenLiteral() string { return sl.Token.Literal }
func (sl *StringLiteral) String() string       { return "\"" + sl.Value + "\"" }

// CharLiteral represents a character
type CharLiteral struct {
	Token lexer.Token
	Value byte
}

func (cl *CharLiteral) expressionNode()      {}
func (cl *CharLiteral) TokenLiteral() string { return cl.Token.Literal }
func (cl *CharLiteral) String() string       { return "'" + string(cl.Value) + "'" }

// BooleanLiteral represents true/false
type BooleanLiteral struct {
	Token lexer.Token
	Value bool
}

func (bl *BooleanLiteral) expressionNode()      {}
func (bl *BooleanLiteral) TokenLiteral() string { return bl.Token.Literal }
func (bl *BooleanLiteral) String() string       { return bl.Token.Literal }

// NullLiteral represents null
type NullLiteral struct {
	Token lexer.Token
}

func (nl *NullLiteral) expressionNode()      {}
func (nl *NullLiteral) TokenLiteral() string { return nl.Token.Literal }
func (nl *NullLiteral) String() string       { return "null" }

// TypeAnnotation represents a type
type TypeAnnotation struct {
	Token    lexer.Token
	Name     string
	IsPtr    bool   // true if *Type
	ArrayLen int    // -1 for slice, 0 for non-array, >0 for fixed array
}

func (t *TypeAnnotation) String() string {
	var out bytes.Buffer
	if t.IsPtr {
		out.WriteString("*")
	}
	if t.ArrayLen == -1 {
		out.WriteString("[]")
	} else if t.ArrayLen > 0 {
		out.WriteString("[" + string(rune(t.ArrayLen)) + "]")
	}
	out.WriteString(t.Name)
	return out.String()
}

// VarStatement: var x int = 5;
type VarStatement struct {
	Token lexer.Token
	Name  *Identifier
	Type  *TypeAnnotation // optional, can be inferred
	Value Expression
}

func (vs *VarStatement) statementNode()       {}
func (vs *VarStatement) TokenLiteral() string { return vs.Token.Literal }
func (vs *VarStatement) String() string {
	var out bytes.Buffer
	out.WriteString("var ")
	out.WriteString(vs.Name.String())
	if vs.Type != nil {
		out.WriteString(" " + vs.Type.String())
	}
	if vs.Value != nil {
		out.WriteString(" = " + vs.Value.String())
	}
	out.WriteString(";")
	return out.String()
}

// ConstStatement: const x := 5;
type ConstStatement struct {
	Token lexer.Token
	Name  *Identifier
	Type  *TypeAnnotation
	Value Expression
}

func (cs *ConstStatement) statementNode()       {}
func (cs *ConstStatement) TokenLiteral() string { return cs.Token.Literal }
func (cs *ConstStatement) String() string {
	var out bytes.Buffer
	out.WriteString("const ")
	out.WriteString(cs.Name.String())
	if cs.Type != nil {
		out.WriteString(" " + cs.Type.String())
	}
	out.WriteString(" := " + cs.Value.String() + ";")
	return out.String()
}

// InferStatement: x := 5;
type InferStatement struct {
	Token lexer.Token
	Name  *Identifier
	Value Expression
}

func (is *InferStatement) statementNode()       {}
func (is *InferStatement) TokenLiteral() string { return is.Token.Literal }
func (is *InferStatement) String() string {
	return is.Name.String() + " := " + is.Value.String() + ";"
}

// ReturnStatement: return x;
type ReturnStatement struct {
	Token lexer.Token
	Value Expression
}

func (rs *ReturnStatement) statementNode()       {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }
func (rs *ReturnStatement) String() string {
	var out bytes.Buffer
	out.WriteString("return")
	if rs.Value != nil {
		out.WriteString(" " + rs.Value.String())
	}
	out.WriteString(";")
	return out.String()
}

// ExpressionStatement wraps an expression as a statement
type ExpressionStatement struct {
	Token      lexer.Token
	Expression Expression
}

func (es *ExpressionStatement) statementNode()       {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String() + ";"
	}
	return ""
}

// BlockStatement: { ... }
type BlockStatement struct {
	Token      lexer.Token
	Statements []Statement
}

func (bs *BlockStatement) statementNode()       {}
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs *BlockStatement) String() string {
	var out bytes.Buffer
	out.WriteString("{\n")
	for _, s := range bs.Statements {
		out.WriteString("  " + s.String() + "\n")
	}
	out.WriteString("}")
	return out.String()
}

// Parameter represents a function parameter
type Parameter struct {
	Name *Identifier
	Type *TypeAnnotation
}

// FunctionStatement: function foo(x int) int { ... }
type FunctionStatement struct {
	Token      lexer.Token
	Public     bool
	Receiver   *Parameter // nil for regular functions
	Name       *Identifier
	Parameters []*Parameter
	ReturnType *TypeAnnotation
	Body       *BlockStatement
}

func (fs *FunctionStatement) statementNode()       {}
func (fs *FunctionStatement) TokenLiteral() string { return fs.Token.Literal }
func (fs *FunctionStatement) String() string {
	var out bytes.Buffer

	if fs.Public {
		out.WriteString("public ")
	}
	out.WriteString("function ")

	if fs.Receiver != nil {
		out.WriteString("(" + fs.Receiver.Name.String() + " " + fs.Receiver.Type.String() + ") ")
	}

	out.WriteString(fs.Name.String())
	out.WriteString("(")

	params := []string{}
	for _, p := range fs.Parameters {
		params = append(params, p.Name.String()+" "+p.Type.String())
	}
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(")")

	if fs.ReturnType != nil {
		out.WriteString(" " + fs.ReturnType.String())
	}

	out.WriteString(" " + fs.Body.String())
	return out.String()
}

// StructField represents a field in a struct
type StructField struct {
	Public bool
	Name   *Identifier
	Type   *TypeAnnotation
}

// StructStatement: struct Foo { ... }
type StructStatement struct {
	Token  lexer.Token
	Public bool
	Name   *Identifier
	Fields []*StructField
}

func (ss *StructStatement) statementNode()       {}
func (ss *StructStatement) TokenLiteral() string { return ss.Token.Literal }
func (ss *StructStatement) String() string {
	var out bytes.Buffer

	if ss.Public {
		out.WriteString("public ")
	}
	out.WriteString("struct ")
	out.WriteString(ss.Name.String())
	out.WriteString(" {\n")

	for _, f := range ss.Fields {
		if f.Public {
			out.WriteString("  public ")
		} else {
			out.WriteString("  ")
		}
		out.WriteString(f.Name.String() + " " + f.Type.String() + ";\n")
	}

	out.WriteString("}")
	return out.String()
}

// IfStatement: if x > 0 { ... } else { ... }
type IfStatement struct {
	Token       lexer.Token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (is *IfStatement) statementNode()       {}
func (is *IfStatement) TokenLiteral() string { return is.Token.Literal }
func (is *IfStatement) String() string {
	var out bytes.Buffer
	out.WriteString("if " + is.Condition.String() + " ")
	out.WriteString(is.Consequence.String())
	if is.Alternative != nil {
		out.WriteString(" else " + is.Alternative.String())
	}
	return out.String()
}

// ForStatement: for i := 0; i < 10; i++ { ... }
type ForStatement struct {
	Token     lexer.Token
	Init      Statement
	Condition Expression
	Post      Statement
	Body      *BlockStatement
}

func (fs *ForStatement) statementNode()       {}
func (fs *ForStatement) TokenLiteral() string { return fs.Token.Literal }
func (fs *ForStatement) String() string {
	var out bytes.Buffer
	out.WriteString("for ")
	if fs.Init != nil {
		out.WriteString(fs.Init.String() + " ")
	}
	if fs.Condition != nil {
		out.WriteString(fs.Condition.String())
	}
	out.WriteString("; ")
	if fs.Post != nil {
		out.WriteString(fs.Post.String())
	}
	out.WriteString(" " + fs.Body.String())
	return out.String()
}

// WhileStatement: while x > 0 { ... }
type WhileStatement struct {
	Token     lexer.Token
	Condition Expression
	Body      *BlockStatement
}

func (ws *WhileStatement) statementNode()       {}
func (ws *WhileStatement) TokenLiteral() string { return ws.Token.Literal }
func (ws *WhileStatement) String() string {
	return "while " + ws.Condition.String() + " " + ws.Body.String()
}

// PrefixExpression: -x, !x, &x, *x
type PrefixExpression struct {
	Token    lexer.Token
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode()      {}
func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Literal }
func (pe *PrefixExpression) String() string {
	return "(" + pe.Operator + pe.Right.String() + ")"
}

// InfixExpression: x + y
type InfixExpression struct {
	Token    lexer.Token
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) expressionNode()      {}
func (ie *InfixExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *InfixExpression) String() string {
	return "(" + ie.Left.String() + " " + ie.Operator + " " + ie.Right.String() + ")"
}

// PostfixExpression: x++, x--
type PostfixExpression struct {
	Token    lexer.Token
	Left     Expression
	Operator string
}

func (pe *PostfixExpression) expressionNode()      {}
func (pe *PostfixExpression) TokenLiteral() string { return pe.Token.Literal }
func (pe *PostfixExpression) String() string {
	return "(" + pe.Left.String() + pe.Operator + ")"
}

// CallExpression: foo(x, y)
type CallExpression struct {
	Token     lexer.Token
	Function  Expression
	Arguments []Expression
}

func (ce *CallExpression) expressionNode()      {}
func (ce *CallExpression) TokenLiteral() string { return ce.Token.Literal }
func (ce *CallExpression) String() string {
	var out bytes.Buffer
	out.WriteString(ce.Function.String())
	out.WriteString("(")

	args := []string{}
	for _, a := range ce.Arguments {
		args = append(args, a.String())
	}
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")
	return out.String()
}

// IndexExpression: arr[0]
type IndexExpression struct {
	Token lexer.Token
	Left  Expression
	Index Expression
}

func (ie *IndexExpression) expressionNode()      {}
func (ie *IndexExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IndexExpression) String() string {
	return "(" + ie.Left.String() + "[" + ie.Index.String() + "])"
}

// MemberExpression: user.name
type MemberExpression struct {
	Token  lexer.Token
	Object Expression
	Member *Identifier
}

func (me *MemberExpression) expressionNode()      {}
func (me *MemberExpression) TokenLiteral() string { return me.Token.Literal }
func (me *MemberExpression) String() string {
	return "(" + me.Object.String() + "." + me.Member.String() + ")"
}

// AssignExpression: x = 5
type AssignExpression struct {
	Token    lexer.Token
	Left     Expression
	Operator string // =, +=, -=, *=, /=
	Value    Expression
}

func (ae *AssignExpression) expressionNode()      {}
func (ae *AssignExpression) TokenLiteral() string { return ae.Token.Literal }
func (ae *AssignExpression) String() string {
	return "(" + ae.Left.String() + " " + ae.Operator + " " + ae.Value.String() + ")"
}

// CastExpression: (int)x
type CastExpression struct {
	Token      lexer.Token
	TargetType *TypeAnnotation
	Value      Expression
}

func (ce *CastExpression) expressionNode()      {}
func (ce *CastExpression) TokenLiteral() string { return ce.Token.Literal }
func (ce *CastExpression) String() string {
	return "((" + ce.TargetType.String() + ")" + ce.Value.String() + ")"
}

// AllocExpression: alloc(User)
type AllocExpression struct {
	Token lexer.Token
	Type  *TypeAnnotation
}

func (ae *AllocExpression) expressionNode()      {}
func (ae *AllocExpression) TokenLiteral() string { return ae.Token.Literal }
func (ae *AllocExpression) String() string {
	return "alloc(" + ae.Type.String() + ")"
}

// FreeStatement: free(ptr);
type FreeStatement struct {
	Token lexer.Token
	Value Expression
}

func (fs *FreeStatement) statementNode()       {}
func (fs *FreeStatement) TokenLiteral() string { return fs.Token.Literal }
func (fs *FreeStatement) String() string {
	return "free(" + fs.Value.String() + ");"
}

// DeferStatement: defer stmt;
type DeferStatement struct {
	Token     lexer.Token
	Statement Statement
}

func (ds *DeferStatement) statementNode()       {}
func (ds *DeferStatement) TokenLiteral() string { return ds.Token.Literal }
func (ds *DeferStatement) String() string {
	return "defer " + ds.Statement.String()
}

// ArrayLiteral: [1, 2, 3] or [5]int{1, 2, 3, 4, 5}
type ArrayLiteral struct {
	Token    lexer.Token
	Type     *TypeAnnotation
	Elements []Expression
}

func (al *ArrayLiteral) expressionNode()      {}
func (al *ArrayLiteral) TokenLiteral() string { return al.Token.Literal }
func (al *ArrayLiteral) String() string {
	var out bytes.Buffer
	if al.Type != nil {
		out.WriteString(al.Type.String())
	}
	out.WriteString("{")
	elems := []string{}
	for _, e := range al.Elements {
		elems = append(elems, e.String())
	}
	out.WriteString(strings.Join(elems, ", "))
	out.WriteString("}")
	return out.String()
}
