package ast

import (
	"bytes"
	"strings"

	"github.com/Rajdeep-Nemo/sugarglaze/internal/token"
)

// Node is the base interface for all nodes in the tree
type Node interface {
	TokenLiteral() string
	String() string
}

// Interface for statements
type Statement interface {
	Node
	statementNode()
}

// Interface for expressions
type Expression interface {
	Node
	expressionNode()
}

// Root node
type Program struct {
	Statements []Statement
}

// Helper function to get the literal value of the first token in the program
func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}
	return ""
}

// Helper function to convert the program to a string for debugging
func (p *Program) String() string {
	var out bytes.Buffer
	for _, s := range p.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

// Identifier node
type Identifier struct {
	Token token.Token // The identifier token
	Value string      // The name of the identifier
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Lexeme }
func (i *Identifier) String() string       { return i.Value }

// Type annotation node
type TypeAnnotation struct {
	Token token.Token // The type annotation token e.g. I32, BOOL
	Name  string      // The name of the type e.g. "i32", "bool"
}

// String method for TypeAnnotation
func (t *TypeAnnotation) String() string {
	return t.Name
}

// Let statement node
type LetStatement struct {
	Token    token.Token     // The 'let' token
	Name     *Identifier     // The variable name
	TypeHint *TypeAnnotation // Optional, nil if type is inferred
	Value    Expression      // Optional, nil if uninitialized
}

func (ls *LetStatement) statementNode()       {}
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Lexeme }

// Return a string representation of the let statement
// Format: "let <name>" or "let <name>: <type>" or "let <name>: <type> = <value>"
func (ls *LetStatement) String() string {
	var out bytes.Buffer
	out.WriteString(ls.Token.Lexeme + " ")
	out.WriteString(ls.Name.String())
	if ls.TypeHint != nil {
		out.WriteString(": " + ls.TypeHint.String())
	}
	if ls.Value != nil {
		out.WriteString(" = " + ls.Value.String())
	}
	return out.String()
}

// Const statement node
type ConstStatement struct {
	Token    token.Token     // The 'const' token
	Name     *Identifier     // The variable name
	TypeHint *TypeAnnotation // Optional, nil if type is inferred
	Value    Expression      // Required, const must always be initialized
}

func (ls *ConstStatement) statementNode()       {}
func (ls *ConstStatement) TokenLiteral() string { return ls.Token.Lexeme }

// Return a string representation of the const statement
// Format: "const <name> = <value>" or "const <name>: <type> = <value>"
func (ls *ConstStatement) String() string {
	var out bytes.Buffer
	out.WriteString(ls.Token.Lexeme + " ")
	out.WriteString(ls.Name.String())
	if ls.TypeHint != nil {
		out.WriteString(": " + ls.TypeHint.String())
	}
	if ls.Value != nil {
		out.WriteString(" = " + ls.Value.String())
	}
	return out.String()
}

// ImportStatement represents an import directive: import "io" or import ("io", "math")
type ImportStatement struct {
	Token   token.Token      // The 'import' token
	Modules []*StringLiteral // List of module names to import
}

func (is *ImportStatement) statementNode()       {}
func (is *ImportStatement) TokenLiteral() string { return is.Token.Lexeme }

func (is *ImportStatement) String() string {
	var out bytes.Buffer
	out.WriteString(is.TokenLiteral() + " ")

	if len(is.Modules) > 1 {
		out.WriteString("(")
	}

	var mods []string
	for _, m := range is.Modules {
		mods = append(mods, m.String())
	}
	out.WriteString(strings.Join(mods, ", "))

	if len(is.Modules) > 1 {
		out.WriteString(")")
	}
	out.WriteString(";")
	return out.String()
}

// Holds the operator for assignment statements like "=", "+=", etc.
type AssignOperator string

const (
	Assign        AssignOperator = "="
	PlusAssign    AssignOperator = "+="
	MinusAssign   AssignOperator = "-="
	StarAssign    AssignOperator = "*="
	SlashAssign   AssignOperator = "/="
	PercentAssign AssignOperator = "%="
)

// Assignment statement node
type AssignStatement struct {
	Token    token.Token    // The operator token e.g. "=", "+=", etc.
	Target   *Identifier    // The variable being assigned to
	Operator AssignOperator // The operator e.g. "=", "+=", etc.
	Value    Expression     // The expression being assigned (right hand side)
}

func (as *AssignStatement) statementNode()       {}
func (as *AssignStatement) TokenLiteral() string { return as.Token.Lexeme }

// Return a string representation of the assignment statement
// Format: "<target> <operator> <value>" e.g. "x += 5"
func (as *AssignStatement) String() string {
	var out bytes.Buffer
	out.WriteString(as.Target.String() + " " + string(as.Operator) + " " + as.Value.String())
	return out.String()
}

// ExpressionStatement represents an expression standing alone on a line (e.g., a function call)
type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

func (es *ExpressionStatement) statementNode()       {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Lexeme }

func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String() + ";"
	}
	return ""
}

// Literal nodes for different literal types
type IntegerLiteral struct {
	Token token.Token
	Value uint64
}

func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Lexeme }
func (il *IntegerLiteral) String() string       { return il.Token.Lexeme }

type FloatLiteral struct {
	Token token.Token
	Value float64
}

func (fl *FloatLiteral) expressionNode()      {}
func (fl *FloatLiteral) TokenLiteral() string { return fl.Token.Lexeme }
func (fl *FloatLiteral) String() string       { return fl.Token.Lexeme }

type StringLiteral struct {
	Token token.Token
	Value string
}

func (sl *StringLiteral) expressionNode()      {}
func (sl *StringLiteral) TokenLiteral() string { return sl.Token.Lexeme }
func (sl *StringLiteral) String() string       { return sl.Token.Lexeme }

type BoolLiteral struct {
	Token token.Token
	Value bool
}

func (bl *BoolLiteral) expressionNode()      {}
func (bl *BoolLiteral) TokenLiteral() string { return bl.Token.Lexeme }
func (bl *BoolLiteral) String() string       { return bl.Token.Lexeme }

type CharLiteral struct {
	Token token.Token
	Value rune
}

func (cl *CharLiteral) expressionNode()      {}
func (cl *CharLiteral) TokenLiteral() string { return cl.Token.Lexeme }
func (cl *CharLiteral) String() string       { return cl.Token.Lexeme }

type NilLiteral struct {
	Token token.Token
}

func (nl *NilLiteral) expressionNode()      {}
func (nl *NilLiteral) TokenLiteral() string { return nl.Token.Lexeme }
func (nl *NilLiteral) String() string       { return nl.Token.Lexeme }

// Prefix operator nodes for unary expressions
type PrefixExpression struct {
	Token    token.Token // The prefix token, e.g. '!' or '-'
	Operator string      // The operator as a string, e.g. "!" or "-"
	Right    Expression  // The expression to the right of the operator
}

func (pe *PrefixExpression) String() string {
	return "(" + pe.Operator + pe.Right.String() + ")"
}
func (pe *PrefixExpression) expressionNode()      {}
func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Lexeme }

// Infix operator nodes for binary expressions
type InfixExpression struct {
	Token    token.Token // The operator token, e.g. '+'
	Left     Expression  // The expression to the left of the operator
	Operator string      // The operator as a string, e.g. "+"
	Right    Expression  // The expression to the right of the operator
}

func (ie *InfixExpression) String() string {
	return "(" + ie.Left.String() + " " + ie.Operator + " " + ie.Right.String() + ")"
}
func (ie *InfixExpression) expressionNode()      {}
func (ie *InfixExpression) TokenLiteral() string { return ie.Token.Lexeme }

type CallExpression struct {
	Token     token.Token  // The '(' token
	Function  Expression   // The function being called (could be an identifier or a more complex expression)
	Arguments []Expression // The arguments passed to the function
}

func (ce *CallExpression) expressionNode()      {}
func (ce *CallExpression) TokenLiteral() string { return ce.Token.Lexeme }
func (ce *CallExpression) String() string {
	var out bytes.Buffer
	var args []string
	for _, a := range ce.Arguments {
		args = append(args, a.String())
	}
	out.WriteString(ce.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")
	return out.String()
}

// BlockStatement represents a sequence of statements enclosed in { }
type BlockStatement struct {
	Token      token.Token // the { token
	Statements []Statement
}

func (bs *BlockStatement) statementNode()       {}
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Lexeme }
func (bs *BlockStatement) String() string {
	var out bytes.Buffer
	for _, s := range bs.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

// IfExpression represents an if/else conditional
type IfExpression struct {
	Token       token.Token // The 'if' token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement // This can be nil if there is no 'else'
}

func (ie *IfExpression) expressionNode()      {}
func (ie *IfExpression) TokenLiteral() string { return ie.Token.Lexeme }
func (ie *IfExpression) String() string {
	var out bytes.Buffer
	out.WriteString("if")
	out.WriteString(ie.Condition.String())
	out.WriteString(" ")
	out.WriteString(ie.Consequence.String())

	if ie.Alternative != nil {
		out.WriteString("else ")
		out.WriteString(ie.Alternative.String())
	}
	return out.String()
}
