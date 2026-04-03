package lexer

import (
	"unicode"

	"github.com/Rajdeep-Nemo/sugarglaze/internal/token"
)

// Keywords map
var keywords = map[string]token.TokenType{
	// Keywords
	"let":      token.LET,
	"if":       token.IF,
	"else":     token.ELSE,
	"loop":     token.LOOP,
	"in":       token.IN,
	"break":    token.BREAK,
	"continue": token.CONTINUE,
	"return":   token.RETURN,
	"fn":       token.FN,
	"static":   token.STATIC,
	"const":    token.CONST,
	"match":    token.MATCH,
	"true":     token.TRUE,
	"false":    token.FALSE,
	"import":   token.IMPORT,
	// Type keywords
	"i8":     token.I8,
	"i16":    token.I16,
	"i32":    token.I32,
	"i64":    token.I64,
	"u8":     token.U8,
	"u16":    token.U16,
	"u32":    token.U32,
	"u64":    token.U64,
	"f32":    token.F32,
	"f64":    token.F64,
	"char":   token.CHAR,
	"string": token.STRING,
	"bool":   token.BOOL,
	"struct": token.STRUCT,
	"enum":   token.ENUM,
	"union":  token.UNION,
	// Special (Reserved literal)
	"NIL": token.NIL_LITERAL,
}

// Scanner struct
type Scanner struct {
	source  string // The source string from the input file
	start   int    // Starting index
	current int    // Current index
	line    uint   // Tracks the line
}

// Initialize the scanner
func InitScanner(source string) *Scanner {
	return &Scanner{
		source:  source,
		start:   0,
		current: 0,
		line:    1,
	}
}

// Moves the pointer forward and returns the current char
func (s *Scanner) advance() byte {
	if s.current >= len(s.source) {
		return 0
	}
	c := s.source[s.current]
	s.current += 1
	return c
}

// Returns true if we are at the end of the file
func (s *Scanner) isAtEnd() bool {
	if s.current >= len(s.source) {
		return true
	}
	return false
}

// Checks the next character
func (s *Scanner) peek() byte {
	if s.isAtEnd() {
		return 0
	}
	return s.source[s.current]
}

// Checks the second next character
func (s *Scanner) peekNext() byte {
	if s.current+1 >= len(s.source) {
		return 0
	}
	return s.source[s.current+1]
}

// Skips whitespace
func (s *Scanner) skipWhitespace() {

	for {
		switch s.peek() {
		case ' ', '\t', '\r':
			s.advance()
		case '\n':
			s.line += 1
			s.advance()
		case '/':
			if s.peekNext() == '/' {
				s.advance()
				s.advance()
				for {
					if !s.isAtEnd() && s.peek() != '\n' {
						s.advance()
					} else {
						break
					}
				}
			} else {
				return
			}
		default:
			return
		}
	}
}

// Create new token
func (s *Scanner) createToken(token_type token.TokenType) token.Token {
	lexeme := s.source[s.start:s.current]
	return token.Token{
		Type:   token_type,
		Lexeme: lexeme,
		Line:   s.line,
	}
}

// Creates a error token for error reporting
func (s *Scanner) errorToken(message string) token.Token {
	lexeme := message
	return token.Token{
		Type:   token.ERROR_TOKEN,
		Lexeme: lexeme,
		Line:   s.line,
	}
}

// Helper function to evaluate conditional advances - '!=' , '=='
func (s *Scanner) match(expected byte) bool {
	if s.isAtEnd() {
		return false
	}
	if s.source[s.current] != expected {
		return false
	}
	s.current += 1
	return true
}

// Helper function to check if it is a character literal (Inside single quotes)
func (s *Scanner) isCharLiteral() token.Token {
	// Checks for empty literal
	if s.isAtEnd() {
		return s.errorToken("Unterminated character literal.")
	}
	if s.peek() == '\'' {
		return s.errorToken("Empty character literal.")
	}
	// Handle Escape Sequence
	if s.peek() == '\\' {
		// Consumes the backslash
		s.advance()
		switch s.peek() {
		case '\'', '"', '\\', 'n', '{', '}', 't', 'r', '0':
			s.advance() // Consumes the valid escape character
		default:
			return s.errorToken("Invalid escape sequence in character literal.")
		}
		// Handles regular character
	} else {
		if s.isAtEnd() {
			return s.errorToken("Unterminated character literal.")
		}
		s.advance()
	}
	// Throws error if not terminated properly
	if s.isAtEnd() {
		return s.errorToken("Unterminated character literal.")
	}
	if s.peek() != '\'' {
		return s.errorToken("Character literal containing multiple character.")
	}
	if s.isAtEnd() {
		return s.errorToken("Unterminated character literal.")
	}
	// Consume the closing single quote and return the char as a token
	s.advance()
	return s.createToken(token.CHAR_LITERAL)
}

// Helper function to check if it is a string literal (Inside double quotes)
func (s *Scanner) isStringLiteral() token.Token {
	// Loop continues until file ends or double quote found
	for s.peek() != '"' && !s.isAtEnd() {
		// Consumes the newline and increases line count
		if s.peek() == '\n' {
			s.line += 1
			s.advance()
		} else if s.peek() == '\\' { // Handles escape sequence
			s.advance()
			if s.isAtEnd() {
				return s.errorToken("Unterminated string after escape.")
			}
			switch s.peek() {
			case '\'', '"', '\\', 'n', '{', '}', 't', 'r', '0':
				s.advance()
			default:
				return s.errorToken("Invalid escape sequence.")
			}
			// Keep consuming characters
		} else {
			s.advance()
		}
	}
	// Throws error if string does not end
	if s.isAtEnd() {
		return s.errorToken("Unterminated string.")
	}
	// Consume the closing double quote and return the string as a token
	s.advance()
	return s.createToken(token.STRING_LITERAL)
}

// Helper function to check if it is a number literal (Integer literal of float literal)
func (s *Scanner) isNumberLiteral() token.Token {
	// Flag to check if float
	isFloat := false
	// Consumes digit/digits
	for unicode.IsDigit(rune(s.peek())) {
		s.advance()
	}
	// If a dot is found and the next character is also a digit it is considered a float
	if s.peek() == '.' && unicode.IsDigit(rune(s.peekNext())) {
		// Sets the flag to true
		isFloat = true
		// Consumes the dot
		s.advance()
		// Consumes remaining digit/digits
		for unicode.IsDigit(rune(s.peek())) {
			s.advance()
		}
	}
	// Create a token based on flag
	if isFloat {
		return s.createToken(token.FLOAT_LITERAL)
	} else {
		return s.createToken(token.INT_LITERAL)
	}
}

// Function to check if it is an identifier (checks keywords as well), uses the keywords map
func (s *Scanner) isIdentifier() token.Token {
	// Consume alphanumeric characters and underscores
	for !s.isAtEnd() && (unicode.IsLetter(rune(s.peek())) || unicode.IsDigit(rune(s.peek())) || s.peek() == '_') {
		s.advance()
	}
	// Slice the source string to get the recent read word
	text := s.source[s.start:s.current]
	// Look for the word in the map and returns the token type and a boolean if it exists or not
	tokenType, exists := keywords[text]
	//   ^        ^
	//   |        |
	// TokenType  bool
	if exists {
		// Creates and returns the token if it's a keyword
		return s.createToken(tokenType)
	}
	//  It's not a keyword, returns IDENTIFIER token
	return s.createToken(token.IDENTIFIER)
}

// Scan and create appropriate token
func (s *Scanner) ScanToken() token.Token {
	// Removes comments and whitespaces
	s.skipWhitespace()
	// Move the pointer to the new token
	s.start = s.current
	// Check if it is the end of file, if it is then returns EOF token
	if s.isAtEnd() {
		return s.createToken(token.END_OF_FILE)
	}
	// Grab the first character
	c := s.advance()
	switch c {
	// Single-character delimiters
	case '(':
		return s.createToken(token.OPEN_PAREN)
	case ')':
		return s.createToken(token.CLOSE_PAREN)
	case '{':
		return s.createToken(token.OPEN_BRACE)
	case '}':
		return s.createToken(token.CLOSE_BRACE)
	case '[':
		return s.createToken(token.OPEN_BRACKET)
	case ']':
		return s.createToken(token.CLOSE_BRACKET)
	case ',':
		return s.createToken(token.COMMA)
	case ':':
		return s.createToken(token.COLON)
	case '?':
		return s.createToken(token.QUESTION)

	// One or two characters
	case '.':
		if s.match('.') {
			return s.createToken(token.DOT_DOT)
		}
		return s.createToken(token.DOT)
	case '+':
		if s.match('=') {
			return s.createToken(token.PLUS_EQUAL)
		}
		return s.createToken(token.PLUS)
	case '*':
		if s.match('=') {
			return s.createToken(token.STAR_EQUAL)
		}
		return s.createToken(token.STAR)
	case '/':
		if s.match('=') {
			return s.createToken(token.SLASH_EQUAL)
		}
		return s.createToken(token.SLASH)
	case '%':
		if s.match('=') {
			return s.createToken(token.PERCENT_EQUAL)
		}
		return s.createToken(token.PERCENT)

	// Three-way branches
	case '-':
		if s.match('>') {
			return s.createToken(token.ARROW)
		}
		if s.match('=') {
			return s.createToken(token.MINUS_EQUAL)
		}
		return s.createToken(token.MINUS)
	case '<':
		if s.match('<') {
			return s.createToken(token.LEFT_SHIFT)
		}
		if s.match('=') {
			return s.createToken(token.LESS_EQUAL)
		}
		return s.createToken(token.LESS)
	case '>':
		if s.match('>') {
			return s.createToken(token.RIGHT_SHIFT)
		}
		if s.match('=') {
			return s.createToken(token.GREATER_EQUAL)
		}
		return s.createToken(token.GREATER)

	// Comparison & Assignment
	case '=':
		if s.match('=') {
			return s.createToken(token.EQUAL_EQUAL)
		}
		return s.createToken(token.EQUAL)
	case '!':
		if s.match('=') {
			return s.createToken(token.BANG_EQUAL)
		}
		return s.createToken(token.BANG)
	case '&':
		if s.match('&') {
			return s.createToken(token.AND)
		}
		return s.createToken(token.BIT_AND)
	case '|':
		if s.match('|') {
			return s.createToken(token.OR)
		}
		return s.createToken(token.BIT_OR)

	// Single Bitwise
	case '^':
		return s.createToken(token.BIT_XOR)
	case '~':
		return s.createToken(token.BIT_NOT)

	// Literals
	case '\'':
		return s.isCharLiteral()
	case '"':
		return s.isStringLiteral()

	// Default case handles Numbers, Identifiers, and Errors
	default:
		if unicode.IsDigit(rune(c)) {
			return s.isNumberLiteral()
		}
		if unicode.IsLetter(rune(c)) || c == '_' {
			return s.isIdentifier()
		}
		// If it's a completely unknown symbol
		return s.errorToken("Unexpected character.")
	}

}
