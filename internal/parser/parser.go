package parser

import (
	"fmt"
	"pluesi/internal/token"
)

// Parser struct to hold the tokens and current position
type Parser struct {
	tokens []token.Token // List of tokens to parse
	pos    int           // Current position in the token list
	errors []string      // List of errors encountered during parsing
}

// New creates a new parser instance with the given tokens
func New(tokens []token.Token) *Parser {
	return &Parser{tokens: tokens, pos: 0}
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
