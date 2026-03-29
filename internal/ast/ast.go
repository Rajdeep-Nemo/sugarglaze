package ast

import (
	"bytes"
	"pluesi/internal/token"
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
