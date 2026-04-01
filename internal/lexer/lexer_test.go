// This test case is AI generated - Gemini 3.1 Pro
package lexer

import (
	"testing"

	"github.com/Rajdeep-Nemo/sugarglaze/internal/token"
)

// expectedToken is a small helper struct just for our tests
type expectedToken struct {
	expectedType   token.TokenType
	expectedLexeme string
	expectedLine   int
}

func TestScanner(t *testing.T) {
	// 1. Define our Table-Driven Tests
	tests := []struct {
		name     string          // Name of the test
		input    string          // The raw source code to lex
		expected []expectedToken // The exact token we expect to get back
	}{
		{
			name:  "Single Character Punctuation",
			input: "(){}[],.?",
			expected: []expectedToken{
				{token.OPEN_PAREN, "(", 1},
				{token.CLOSE_PAREN, ")", 1},
				{token.OPEN_BRACE, "{", 1},
				{token.CLOSE_BRACE, "}", 1},
				{token.OPEN_BRACKET, "[", 1},
				{token.CLOSE_BRACKET, "]", 1},
				{token.COMMA, ",", 1},
				{token.DOT, ".", 1},
				{token.QUESTION, "?", 1},
				{token.END_OF_FILE, "", 1},
			},
		},
		{
			name:  "Two-Character Operators",
			input: "!= == <= >= -> ..",
			expected: []expectedToken{
				{token.BANG_EQUAL, "!=", 1},
				{token.EQUAL_EQUAL, "==", 1},
				{token.LESS_EQUAL, "<=", 1},
				{token.GREATER_EQUAL, ">=", 1},
				{token.ARROW, "->", 1},
				{token.DOT_DOT, "..", 1},
				{token.END_OF_FILE, "", 1},
			},
		},
		{
			name:  "Mathematical Operators",
			input: "+ - * / %",
			expected: []expectedToken{
				{token.PLUS, "+", 1},
				{token.MINUS, "-", 1},
				{token.STAR, "*", 1},
				{token.SLASH, "/", 1},
				{token.PERCENT, "%", 1},
				{token.END_OF_FILE, "", 1},
			},
		},
		{
			name:  "Compound Assignment Operators",
			input: "+= -= *= /= %=",
			expected: []expectedToken{
				{token.PLUS_EQUAL, "+=", 1},
				{token.MINUS_EQUAL, "-=", 1},
				{token.STAR_EQUAL, "*=", 1},
				{token.SLASH_EQUAL, "/=", 1},
				{token.PERCENT_EQUAL, "%=", 1},
				{token.END_OF_FILE, "", 1},
			},
		},
		{
			name:  "Bitwise and Logical Operators",
			input: "& && | || ^ ~ << >> !",
			expected: []expectedToken{
				{token.BIT_AND, "&", 1},
				{token.AND, "&&", 1},
				{token.BIT_OR, "|", 1},
				{token.OR, "||", 1},
				{token.BIT_XOR, "^", 1},
				{token.BIT_NOT, "~", 1},
				{token.LEFT_SHIFT, "<<", 1},
				{token.RIGHT_SHIFT, ">>", 1},
				{token.BANG, "!", 1},
				{token.END_OF_FILE, "", 1},
			},
		},
		{
			name:  "Control Flow Keywords",
			input: "if else loop in break continue match return",
			expected: []expectedToken{
				{token.IF, "if", 1},
				{token.ELSE, "else", 1},
				{token.LOOP, "loop", 1},
				{token.IN, "in", 1},
				{token.BREAK, "break", 1},
				{token.CONTINUE, "continue", 1},
				{token.MATCH, "match", 1},
				{token.RETURN, "return", 1},
				{token.END_OF_FILE, "", 1},
			},
		},
		{
			name:  "Data Type Keywords",
			input: "i8 u32 f64 char string bool struct enum union",
			expected: []expectedToken{
				{token.I8, "i8", 1},
				{token.U32, "u32", 1},
				{token.F64, "f64", 1},
				{token.CHAR, "char", 1},
				{token.STRING, "string", 1},
				{token.BOOL, "bool", 1},
				{token.STRUCT, "struct", 1},
				{token.ENUM, "enum", 1},
				{token.UNION, "union", 1},
				{token.END_OF_FILE, "", 1},
			},
		},
		{
			name:  "Boolean and Nil Literals",
			input: "true false NIL",
			expected: []expectedToken{
				{token.TRUE, "true", 1},
				{token.FALSE, "false", 1},
				{token.NIL_LITERAL, "NIL", 1},
				{token.END_OF_FILE, "", 1},
			},
		},
		{
			name:  "Character Literals",
			input: "'a' '\\n'",
			expected: []expectedToken{
				{token.CHAR_LITERAL, "'a'", 1},
				{token.CHAR_LITERAL, "'\\n'", 1},
				{token.END_OF_FILE, "", 1},
			},
		},
		{
			name:  "Complex Identifiers",
			input: "_privateVar camelCase123 another_var",
			expected: []expectedToken{
				{token.IDENTIFIER, "_privateVar", 1},
				{token.IDENTIFIER, "camelCase123", 1},
				{token.IDENTIFIER, "another_var", 1},
				{token.END_OF_FILE, "", 1},
			},
		},
		{
			name:  "String Literal and Whitespace",
			input: "   \n \"hello world\"  \n",
			expected: []expectedToken{
				{token.STRING_LITERAL, "\"hello world\"", 2},
				{token.END_OF_FILE, "", 3},
			},
		},
		{
			name:  "Unterminated String Error",
			input: "\"this string never ends",
			expected: []expectedToken{
				{token.ERROR_TOKEN, "Unterminated string.", 1}, // Added the period based on your earlier code, adjust if yours lacks it
				{token.END_OF_FILE, "", 1},
			},
		},
		{
			name: "Full Snippet with Comments",
			input: `
                // This is a comment
                fn calculate() i32 {
                    return 100
                }
            `,
			expected: []expectedToken{
				{token.FN, "fn", 3},
				{token.IDENTIFIER, "calculate", 3},
				{token.OPEN_PAREN, "(", 3},
				{token.CLOSE_PAREN, ")", 3},
				{token.I32, "i32", 3},
				{token.OPEN_BRACE, "{", 3},
				{token.RETURN, "return", 4},
				{token.INT_LITERAL, "100", 4},
				{token.CLOSE_BRACE, "}", 5},
				{token.END_OF_FILE, "", 6},
			},
		},
		{
			name: "EVIL: Alien Characters",
			// Someone drops a random @ and $ in the code
			input: "let @x = $5",
			expected: []expectedToken{
				{token.LET, "let", 1},
				{token.ERROR_TOKEN, "Unexpected character.", 1}, // The @
				{token.IDENTIFIER, "x", 1},
				{token.EQUAL, "=", 1},
				{token.ERROR_TOKEN, "Unexpected character.", 1}, // The $
				{token.INT_LITERAL, "5", 1},                     // It should recover and keep reading!
				{token.END_OF_FILE, "", 1},
			},
		},
		{
			name: "EVIL: EOF Triggered Mid-Escape Sequence",
			// A string that ends literally right after the escape backslash
			input: "\"This string ends with an escape \\",
			expected: []expectedToken{
				{token.ERROR_TOKEN, "Unterminated string after escape.", 1},
				{token.END_OF_FILE, "", 1},
			},
		},
		{
			name: "EVIL: EOF Triggered Mid-Comment",
			// A comment with no newline at the end of the file
			input: "let x = 10 // This file ends right he...",
			expected: []expectedToken{
				{token.LET, "let", 1},
				{token.IDENTIFIER, "x", 1},
				{token.EQUAL, "=", 1},
				{token.INT_LITERAL, "10", 1},
				{token.END_OF_FILE, "", 1}, // It should gracefully hit EOF, not freeze
			},
		},
		{
			name: "EVIL: Whitespace Chaos",
			// Mixing tabs, spaces, Windows carriage returns (\r), and Linux newlines (\n)
			input: " \t\r\n  \n\n  let\r\n\tx",
			expected: []expectedToken{
				{token.LET, "let", 4}, // Should correctly track lines despite the mess
				{token.IDENTIFIER, "x", 5},
				{token.END_OF_FILE, "", 5},
			},
		},
		{
			name: "EVIL: Operator Soup (Maximal Munch Test)",
			// Lexers use "Maximal Munch" (grab the longest match possible).
			// ===>= should become `==`, `=`, `>=`
			input: "===>=",
			expected: []expectedToken{
				{token.EQUAL_EQUAL, "==", 1},
				{token.EQUAL, "=", 1},
				{token.GREATER_EQUAL, ">=", 1},
				{token.END_OF_FILE, "", 1},
			},
		},
		{
			name: "EVIL: Numbers glued to Identifiers",
			// Standard C-like languages separate these into two token
			input: "99bottles",
			expected: []expectedToken{
				{token.INT_LITERAL, "99", 1},
				{token.IDENTIFIER, "bottles", 1},
				{token.END_OF_FILE, "", 1},
			},
		},
		{
			name:  "EVIL: Completely Empty File",
			input: "    \t   \n  ",
			expected: []expectedToken{
				{token.END_OF_FILE, "", 2}, // Should just skip whitespace and exit safely
			},
		},
	}

	// 2. Loop through all the test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Initialize scanner with the test input
			scanner := InitScanner(tt.input)

			// Check every expected token
			for i, exp := range tt.expected {
				tok := scanner.ScanToken()

				// Check Token Type
				if tok.Type != exp.expectedType {
					t.Errorf("[%s] Token %d: Expected Type %v, got %v (Lexeme: '%s')",
						tt.name, i, exp.expectedType, tok.Type, tok.Lexeme)
				}

				// Check Lexeme
				if tok.Lexeme != exp.expectedLexeme {
					t.Errorf("[%s] Token %d: Expected Lexeme '%s', got '%s'",
						tt.name, i, exp.expectedLexeme, tok.Lexeme)
				}

				// Check Line Number
				if tok.Line != uint(exp.expectedLine) {
					t.Errorf("[%s] Token %d: Expected Line %d, got %d (Lexeme: '%s')",
						tt.name, i, exp.expectedLine, tok.Line, tok.Lexeme)
				}

				// Break early if we hit EOF so the loop stops
				if tok.Type == token.END_OF_FILE {
					break
				}
			}
		})
	}
}
