package parser

import (
	"fmt"
	"strconv"

	"github.com/Rajdeep-Nemo/sugarglaze/internal/ast"
	"github.com/Rajdeep-Nemo/sugarglaze/internal/token"
)

// Parser struct to hold the tokens and current position
type Parser struct {
	tokens []token.Token // List of tokens to parse
	pos    int           // Current position in the token list
	errors []string      // List of errors encountered during parsing
}

// New creates a new parser instance with the given tokens
func New(tokens []token.Token) *Parser {
	p := &Parser{tokens: tokens, pos: 0}
	p.registerParseFns()
	return p
}

// To retrieve errors encountered during parsing
func (p *Parser) Errors() []string {
	return p.errors
}

// Helper function to get the current token
func (p *Parser) currentToken() token.Token {
	if p.pos < len(p.tokens) {
		return p.tokens[p.pos]
	}
	return token.Token{Type: token.END_OF_FILE, Lexeme: ""}
}

// Helper function to peek at the next token without advancing the position
func (p *Parser) peekToken() token.Token {
	if p.pos+1 < len(p.tokens) {
		return p.tokens[p.pos+1]
	}
	return token.Token{Type: token.END_OF_FILE, Lexeme: ""}
}

// Helper function to check if the current token is of the expected type
func (p *Parser) check(tt token.TokenType) bool {
	return p.currentToken().Type == tt
}

// Helper function to advance to the next token
func (p *Parser) advance() token.Token {
	t := p.currentToken()
	p.pos += 1
	return t
}

// Helper function to add an error message to the parser's error list
func (p *Parser) errorf(format string, args ...any) {
	p.errors = append(p.errors, fmt.Sprintf(format, args...))
}

// Helper function to check if the current token is of the expected type
func (p *Parser) expect(tt token.TokenType) (token.Token, bool) {
	t := p.currentToken()
	if t.Type != tt {
		p.errorf("expected %v but got %q at line %d", tt, t.Lexeme, t.Line)
		return token.Token{}, false
	}
	p.advance()
	return t, true
}

// Loops through the tokens and parse statements until the end of the file is reached
func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	for !p.check(token.END_OF_FILE) {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
	}
	return program
}

// Checks the current token and decides which statement parsing function to call based on the token type
func (p *Parser) parseStatement() ast.Statement {
	switch p.currentToken().Type {
	case token.LET:
		return p.parseLetStatement()
	case token.CONST:
		return p.parseConstStatement()
	case token.IMPORT:
		return p.parseImportStatement()
	case token.IDENTIFIER:
		if _, ok := assignOperators[p.peekToken().Type]; ok {
			return p.parseAssignmentStatement()
		}
		fallthrough
	default:
		return p.parseExpressionStatement()
	}
}

// Type tokens for type annotations
var typeTokens = map[token.TokenType]string{
	token.I8:     "i8",
	token.I16:    "i16",
	token.I32:    "i32",
	token.I64:    "i64",
	token.U8:     "u8",
	token.U16:    "u16",
	token.U32:    "u32",
	token.U64:    "u64",
	token.F32:    "f32",
	token.F64:    "f64",
	token.CHAR:   "char",
	token.STRING: "string",
	token.BOOL:   "bool",
}

// Parses a type annotation after a colon in a let/const statement, e.g. `: i32`
func (p *Parser) parseTypeAnnotation() *ast.TypeAnnotation {
	t := p.currentToken()
	name, ok := typeTokens[t.Type]
	if !ok {
		p.errorf("expected a type but got %q at line %d", t.Lexeme, t.Line)
		return nil
	}
	p.advance()
	return &ast.TypeAnnotation{Token: t, Name: name}
}

// Parses a let statement, which can optionally include a type annotation.
func (p *Parser) parseLetStatement() ast.Statement {
	stmt := &ast.LetStatement{Token: p.currentToken()}

	p.advance() // Move past the 'let' keyword

	// 1. Get the identifier name
	identToken, ok := p.expect(token.IDENTIFIER)
	if !ok {
		return nil
	}
	stmt.Name = &ast.Identifier{Token: identToken, Value: identToken.Lexeme}

	// 2. OPTIONAL Type Hint
	if p.currentToken().Type == token.COLON {
		p.advance() // Move past the ':'

		stmt.TypeHint = p.parseTypeAnnotation()
		if stmt.TypeHint == nil {
			return nil
		}
	}

	// 3. Handle Initialization vs Uninitialized
	if p.currentToken().Type == token.EQUAL {
		p.advance() // Move past the '='

		// Parse the right side of the equals sign
		stmt.Value = p.parseExpression(LOWEST)
	} else {
		// If there is NO equals sign, there MUST be a type hint (e.g., `let a;` is invalid)
		if stmt.TypeHint == nil {
			p.errorf("expected '=' after identifier without type hint at line %d", p.currentToken().Line)
			return nil
		}
	}

	return stmt
}

// Parses a const statement, which MUST include a type annotation and an initializer.
func (p *Parser) parseConstStatement() ast.Statement {
	stmt := &ast.ConstStatement{Token: p.currentToken()}

	p.advance() // Move past the 'const' keyword

	// 1. Get the identifier name
	identToken, ok := p.expect(token.IDENTIFIER)
	if !ok {
		return nil
	}
	stmt.Name = &ast.Identifier{Token: identToken, Value: identToken.Lexeme}

	// 2. MANDATORY Type Hint
	if _, ok := p.expect(token.COLON); !ok {
		return nil
	}

	stmt.TypeHint = p.parseTypeAnnotation()
	if stmt.TypeHint == nil {
		return nil
	}

	// 3. Expect the Equals sign
	if _, ok := p.expect(token.EQUAL); !ok {
		return nil
	}

	// 4. Parse the right side of the equals sign
	stmt.Value = p.parseExpression(LOWEST)

	return stmt
}

// Parses an identifier (variable) in an expression
func (p *Parser) parseIdentifier() ast.Expression {
	ident := &ast.Identifier{Token: p.currentToken(), Value: p.currentToken().Lexeme}
	p.advance()
	return ident
}

// Map of assignment operators for easy lookup
var assignOperators = map[token.TokenType]ast.AssignOperator{
	token.EQUAL:         ast.Assign,
	token.PLUS_EQUAL:    ast.PlusAssign,
	token.MINUS_EQUAL:   ast.MinusAssign,
	token.STAR_EQUAL:    ast.StarAssign,
	token.SLASH_EQUAL:   ast.SlashAssign,
	token.PERCENT_EQUAL: ast.PercentAssign,
}

// Parses an assignment statement, which must be in the form of `<identifier> <assignment operators> <expression>`.
func (p *Parser) parseAssignmentStatement() *ast.AssignStatement {
	nameToken, ok := p.expect(token.IDENTIFIER)
	if !ok {
		return nil
	}
	name := &ast.Identifier{Token: nameToken, Value: nameToken.Lexeme}
	opToken := p.currentToken()
	assignOp, ok := assignOperators[opToken.Type]
	if !ok {
		p.errorf("expected an assignment operator but got %q at line %d", p.currentToken().Lexeme, p.currentToken().Line)
		p.advance()
		return nil
	}
	p.advance()
	value := p.parseExpression(LOWEST)
	return &ast.AssignStatement{Token: opToken, Target: name, Operator: assignOp, Value: value}
}

// Parses an import statement: import "io"; OR import ("io", "math");
func (p *Parser) parseImportStatement() *ast.ImportStatement {
	stmt := &ast.ImportStatement{Token: p.currentToken()}
	p.advance() // Consume the 'import' keyword

	// Handle grouped imports: import ("io", "math")
	if p.check(token.OPEN_PAREN) {
		p.advance() // Consume '('

		for !p.check(token.CLOSE_PAREN) && !p.check(token.END_OF_FILE) {
			if !p.check(token.STRING_LITERAL) {
				p.errorf("expected string literal in import, got %q at line %d", p.currentToken().Lexeme, p.currentToken().Line)
				return nil
			}

			// Parse the string literal and cast it
			strNode := p.parseStringLiteral().(*ast.StringLiteral)
			stmt.Modules = append(stmt.Modules, strNode)

			// If the next token is a comma, consume it so we can read the next string
			if p.check(token.COMMA) {
				p.advance()
			}
		}

		if !p.check(token.CLOSE_PAREN) {
			p.errorf("expected ')' after grouped imports at line %d", p.currentToken().Line)
			return nil
		}
		p.advance() // Consume ')'
	} else {
		// Handle single import: import "io"
		if !p.check(token.STRING_LITERAL) {
			p.errorf("expected string literal after import at line %d", p.currentToken().Line)
			return nil
		}

		strNode := p.parseStringLiteral().(*ast.StringLiteral)
		stmt.Modules = append(stmt.Modules, strNode)
	}
	return stmt
}

// Parses an expression using Pratt parsing technique.
func (p *Parser) parseExpression(precedence Precedence) ast.Expression {
	// Fetch the prefix parse function based on the current token type
	prefix := prefixParseFns[p.currentToken().Type]
	if prefix == nil {
		p.errorf("no prefix parse function for %q at line %d", p.currentToken().Lexeme, p.currentToken().Line)
		p.advance()
		return nil
	}
	// Call the prefix parse function to get the left-hand side expression
	leftExp := prefix()
	// Loop to parse infix expressions as long as the next token has a higher precedence
	for precedence < p.currentPrecedence() {
		infix := infixParseFns[p.currentToken().Type]
		if infix == nil {
			return leftExp
		}

		leftExp = infix(leftExp)
	}
	return leftExp
}

// Precedence orders for operators, used in expression parsing to determine which operations bind more tightly than others
type Precedence int

const (
	LOWEST     Precedence = iota
	OR                    // ||
	AND                   // &&
	EQUALS                // == !=
	COMPARISON            // < > <= >=
	SUM                   // + -
	PRODUCT               // * / %
	PREFIX                // -x !x
	CALL                  // foo()
	INDEX                 // arr[]
)

// Precedence lookup table operators
var precedences = map[token.TokenType]Precedence{
	token.OR:            OR,
	token.AND:           AND,
	token.EQUAL_EQUAL:   EQUALS,
	token.BANG_EQUAL:    EQUALS,
	token.LESS:          COMPARISON,
	token.LESS_EQUAL:    COMPARISON,
	token.GREATER:       COMPARISON,
	token.GREATER_EQUAL: COMPARISON,
	token.PLUS:          SUM,
	token.MINUS:         SUM,
	token.STAR:          PRODUCT,
	token.SLASH:         PRODUCT,
	token.PERCENT:       PRODUCT,
	token.OPEN_PAREN:    CALL,
}

func (p *Parser) currentPrecedence() Precedence {
	if pr, ok := precedences[p.currentToken().Type]; ok {
		return pr
	}
	return LOWEST
}

// Function mapping to determine how to parse different expressions based on the current token
type prefixParseFn func() ast.Expression
type infixParseFn func(ast.Expression) ast.Expression

var prefixParseFns map[token.TokenType]prefixParseFn
var infixParseFns map[token.TokenType]infixParseFn

func (p *Parser) registerParseFns() {
	// Prefix functions
	prefixParseFns = map[token.TokenType]prefixParseFn{
		token.IDENTIFIER:     p.parseIdentifier,
		token.INT_LITERAL:    p.parseIntegerLiteral,
		token.FLOAT_LITERAL:  p.parseFloatLiteral,
		token.STRING_LITERAL: p.parseStringLiteral,
		token.CHAR_LITERAL:   p.parseCharLiteral,
		token.TRUE:           p.parseBoolLiteral,
		token.FALSE:          p.parseBoolLiteral,
		token.NIL_LITERAL:    p.parseNilLiteral,
		token.BANG:           p.parsePrefixExpression,
		token.MINUS:          p.parsePrefixExpression,
		token.OPEN_PAREN:     p.parseGroupedExpression,
		token.IF:             p.parseIfExpression,
	}
	// Infix functions
	infixParseFns = map[token.TokenType]infixParseFn{
		token.PLUS:          p.parseInfixExpression,
		token.MINUS:         p.parseInfixExpression,
		token.STAR:          p.parseInfixExpression,
		token.SLASH:         p.parseInfixExpression,
		token.PERCENT:       p.parseInfixExpression,
		token.EQUAL_EQUAL:   p.parseInfixExpression,
		token.BANG_EQUAL:    p.parseInfixExpression,
		token.LESS:          p.parseInfixExpression,
		token.LESS_EQUAL:    p.parseInfixExpression,
		token.GREATER:       p.parseInfixExpression,
		token.GREATER_EQUAL: p.parseInfixExpression,
		token.AND:           p.parseInfixExpression,
		token.OR:            p.parseInfixExpression,
		token.OPEN_PAREN:    p.parseCallExpression,
	}
}

// Parses an integer literal: e.g., 5, 100
func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.currentToken()}

	value, err := strconv.ParseUint(p.currentToken().Lexeme, 10, 64)
	if err != nil {
		p.errorf("could not parse %q as integer at line %d", p.currentToken().Lexeme, p.currentToken().Line)
		return nil
	}

	lit.Value = value
	p.advance()
	return lit
}

// Parses a float literal: e.g., 3.14
func (p *Parser) parseFloatLiteral() ast.Expression {
	lit := &ast.FloatLiteral{Token: p.currentToken()}

	value, err := strconv.ParseFloat(p.currentToken().Lexeme, 64)
	if err != nil {
		p.errorf("could not parse %q as float at line %d", p.currentToken().Lexeme, p.currentToken().Line)
		return nil
	}

	lit.Value = value
	p.advance()
	return lit
}

// Parses a string literal: e.g., "hello world"
func (p *Parser) parseStringLiteral() ast.Expression {
	lexeme := p.currentToken().Lexeme
	value := lexeme
	if len(lexeme) >= 2 && lexeme[0] == '"' && lexeme[len(lexeme)-1] == '"' {
		value = lexeme[1 : len(lexeme)-1] // Strip the " quotes
	}

	lit := &ast.StringLiteral{Token: p.currentToken(), Value: value}
	p.advance()
	return lit
}

// Parses boolean literals: true, false
func (p *Parser) parseBoolLiteral() ast.Expression {
	lit := &ast.BoolLiteral{
		Token: p.currentToken(),
		Value: p.currentToken().Type == token.TRUE,
	}
	p.advance()
	return lit
}

// Parses a char literal: e.g., 'a'
func (p *Parser) parseCharLiteral() ast.Expression {
	lexeme := p.currentToken().Lexeme
	var val rune
	if len(lexeme) >= 3 {
		val = rune(lexeme[1])
	}

	lit := &ast.CharLiteral{Token: p.currentToken(), Value: val}
	p.advance()
	return lit
}

// Parses the NIL literal
func (p *Parser) parseNilLiteral() ast.Expression {
	lit := &ast.NilLiteral{Token: p.currentToken()}
	p.advance()
	return lit
}

// Parses prefix expressions like -5 or !true
func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    p.currentToken(),
		Operator: p.currentToken().Lexeme,
	}

	p.advance() // Consume the operator (e.g., '-')
	expression.Right = p.parseExpression(PREFIX)

	return expression
}

// Parses infix expressions like 5 + 5 or x == y
func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    p.currentToken(),
		Operator: p.currentToken().Lexeme,
		Left:     left,
	}

	precedence := p.currentPrecedence()
	p.advance() // Consume the operator (e.g., '+')
	expression.Right = p.parseExpression(precedence)

	return expression
}

// Parses expressions grouped by parentheses
func (p *Parser) parseGroupedExpression() ast.Expression {
	p.advance() // Consume the '('
	exp := p.parseExpression(LOWEST)

	if !p.check(token.CLOSE_PAREN) {
		p.errorf("expected ')' but got %q at line %d", p.currentToken().Lexeme, p.currentToken().Line)
		return nil
	}

	p.advance() // Consume the ')'
	return exp
}

// Parses a function call expression
func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	exp := &ast.CallExpression{Token: p.currentToken(), Function: function}
	exp.Arguments = p.parseCallArguments()
	return exp
}

// Helper function to parse a comma-separated list of arguments
func (p *Parser) parseCallArguments() []ast.Expression {
	args := []ast.Expression{}

	p.advance() // Move past the '('

	// Check if the very next token is ')', meaning no arguments: print()
	if p.check(token.CLOSE_PAREN) {
		p.advance() // Consume the ')'
		return args
	}

	// Parse the first argument
	args = append(args, p.parseExpression(LOWEST))

	// While the current token is a comma, keep parsing arguments
	for p.check(token.COMMA) {
		p.advance() // Consume the comma
		args = append(args, p.parseExpression(LOWEST))
	}

	// Enforce the closing parenthesis
	if !p.check(token.CLOSE_PAREN) {
		p.errorf("expected ')' but got %q at line %d", p.currentToken().Lexeme, p.currentToken().Line)
		return nil
	}

	p.advance() // Consume the ')'
	return args
}

// Parses an expression standing as its own statement
func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.currentToken()}

	// Parse the expression (like print("hello")) using the lowest precedence
	stmt.Expression = p.parseExpression(LOWEST)
	return stmt
}

// Parses a block of statements enclosed in { }
func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: p.currentToken()}
	block.Statements = []ast.Statement{}

	p.advance() // Consume the opening '{'

	// Loop until we hit a '}' or run out of tokens
	for p.currentToken().Type != token.CLOSE_BRACE && p.currentToken().Type != token.END_OF_FILE {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
	}

	// Consume the '}'
	p.expect(token.CLOSE_BRACE)

	return block
}

// Parses an if/else expression
func (p *Parser) parseIfExpression() ast.Expression {
	expression := &ast.IfExpression{Token: p.currentToken()}

	p.advance() // Move past the 'if' keyword

	// Use expect for '('
	if _, ok := p.expect(token.OPEN_PAREN); !ok {
		return nil
	}

	expression.Condition = p.parseExpression(LOWEST)

	// Use expect for ')'
	if _, ok := p.expect(token.CLOSE_PAREN); !ok {
		return nil
	}

	// FIX: Use check() here so we DON'T double-advance!
	// parseBlockStatement needs to see the '{' to consume it properly.
	if !p.check(token.OPEN_BRACE) {
		p.errorf("expected '{' but got %q at line %d", p.currentToken().Lexeme, p.currentToken().Line)
		return nil
	}
	expression.Consequence = p.parseBlockStatement()

	// Check if there is an 'else' block
	if p.currentToken().Type == token.ELSE {
		p.advance() // Move past 'else'

		// FIX: Use check() here as well!
		if !p.check(token.OPEN_BRACE) {
			p.errorf("expected '{' after else but got %q at line %d", p.currentToken().Lexeme, p.currentToken().Line)
			return nil
		}
		expression.Alternative = p.parseBlockStatement()
	}

	return expression
}
