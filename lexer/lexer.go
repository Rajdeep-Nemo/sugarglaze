package lexer

import (
	"Me/tokens"
	"unicode"
)

// Keywords map
var keywords = map[string]tokens.TokenType{
	// Keywords
	"let":      tokens.LET,
	"if":       tokens.IF,
	"else":     tokens.ELSE,
	"loop":     tokens.LOOP,
	"in":       tokens.IN,
	"break":    tokens.BREAK,
	"continue": tokens.CONTINUE,
	"return":   tokens.RETURN,
	"fn":       tokens.FN,
	"static":   tokens.STATIC,
	"const":    tokens.CONST,
	"match":    tokens.MATCH,
	"true":     tokens.TRUE,
	"false":    tokens.FALSE,
	// Type keywords
	"i8":     tokens.I8,
	"i16":    tokens.I16,
	"i32":    tokens.I32,
	"i64":    tokens.I64,
	"u8":     tokens.U8,
	"u16":    tokens.U16,
	"u32":    tokens.U32,
	"u64":    tokens.U64,
	"f32":    tokens.F32,
	"f64":    tokens.F64,
	"char":   tokens.CHAR,
	"string": tokens.STRING,
	"bool":   tokens.BOOL,
	"struct": tokens.STRUCT,
	"enum":   tokens.ENUM,
	"union":  tokens.UNION,
	// Special (Reserved literal)
	"NIL": tokens.NIL_LITERAL,
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

// Create new tokens
func (s *Scanner) createToken(token_type tokens.TokenType) tokens.Token {
	lexeme := s.source[s.start:s.current]
	return tokens.Token{
		Type:   token_type,
		Lexeme: lexeme,
		Line:   s.line,
	}
}

// Creates a error token for error reporting
func (s *Scanner) errorToken(message string) tokens.Token {
	lexeme := message
	return tokens.Token{
		Type:   tokens.ERROR_TOKEN,
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
func (s *Scanner) isCharLiteral() tokens.Token {
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
	return s.createToken(tokens.CHAR_LITERAL)
}

// Helper function to check if it is a string literal (Inside double quotes)
func (s *Scanner) isStringLiteral() tokens.Token {
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
	return s.createToken(tokens.STRING_LITERAL)
}

// Helper function to check if it is a number literal (Integer literal of float literal)
func (s *Scanner) isNumberLiteral() tokens.Token {
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
		return s.createToken(tokens.FLOAT_LITERAL)
	} else {
		return s.createToken(tokens.INT_LITERAL)
	}
}

// Function to check if it is an identifier (checks keywords as well), uses the keywords map
func (s *Scanner) isIdentifier() tokens.Token {
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
	return s.createToken(tokens.IDENTIFIER)
}

// Scan and create appropriate tokens
func (s *Scanner) ScanToken() tokens.Token {
	// Removes comments and whitespaces
	s.skipWhitespace()
	// Move the pointer to the new token
	s.start = s.current
	// Check if it is the end of file, if it is then returns EOF token
	if s.isAtEnd() {
		return s.createToken(tokens.END_OF_FILE)
	}
	// Grab the first character
	c := s.advance()
	switch c {
	// Single-character delimiters
	case '(':
		return s.createToken(tokens.OPEN_PAREN)
	case ')':
		return s.createToken(tokens.CLOSE_PAREN)
	case '{':
		return s.createToken(tokens.OPEN_BRACE)
	case '}':
		return s.createToken(tokens.CLOSE_BRACE)
	case '[':
		return s.createToken(tokens.OPEN_BRACKET)
	case ']':
		return s.createToken(tokens.CLOSE_BRACKET)
	case ',':
		return s.createToken(tokens.COMMA)
	case ':':
		return s.createToken(tokens.COLON)
	case ';':
		return s.createToken(tokens.SEMICOLON)
	case '?':
		return s.createToken(tokens.QUESTION)

	// One or two characters
	case '.':
		if s.match('.') {
			return s.createToken(tokens.DOT_DOT)
		}
		return s.createToken(tokens.DOT)
	case '+':
		if s.match('=') {
			return s.createToken(tokens.PLUS_EQUAL)
		}
		return s.createToken(tokens.PLUS)
	case '*':
		if s.match('=') {
			return s.createToken(tokens.STAR_EQUAL)
		}
		return s.createToken(tokens.STAR)
	case '/':
		if s.match('=') {
			return s.createToken(tokens.SLASH_EQUAL)
		}
		return s.createToken(tokens.SLASH)
	case '%':
		if s.match('=') {
			return s.createToken(tokens.PERCENT_EQUAL)
		}
		return s.createToken(tokens.PERCENT)

	// Three-way branches
	case '-':
		if s.match('>') {
			return s.createToken(tokens.ARROW)
		}
		if s.match('=') {
			return s.createToken(tokens.MINUS_EQUAL)
		}
		return s.createToken(tokens.MINUS)
	case '<':
		if s.match('<') {
			return s.createToken(tokens.LEFT_SHIFT)
		}
		if s.match('=') {
			return s.createToken(tokens.LESS_EQUAL)
		}
		return s.createToken(tokens.LESS)
	case '>':
		if s.match('>') {
			return s.createToken(tokens.RIGHT_SHIFT)
		}
		if s.match('=') {
			return s.createToken(tokens.GREATER_EQUAL)
		}
		return s.createToken(tokens.GREATER)

	// Comparison & Assignment
	case '=':
		if s.match('=') {
			return s.createToken(tokens.EQUAL_EQUAL)
		}
		return s.createToken(tokens.EQUAL)
	case '!':
		if s.match('=') {
			return s.createToken(tokens.BANG_EQUAL)
		}
		return s.createToken(tokens.BANG)
	case '&':
		if s.match('&') {
			return s.createToken(tokens.AND)
		}
		return s.createToken(tokens.BIT_AND)
	case '|':
		if s.match('|') {
			return s.createToken(tokens.OR)
		}
		return s.createToken(tokens.BIT_OR)

	// Single Bitwise
	case '^':
		return s.createToken(tokens.BIT_XOR)
	case '~':
		return s.createToken(tokens.BIT_NOT)

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
