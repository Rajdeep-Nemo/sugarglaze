package token

// Enum that holds all Token types
type TokenType int

const (
	// 1. Special or Control
	END_OF_FILE TokenType = iota
	ERROR_TOKEN
	// 2. Identifiers & Literals
	IDENTIFIER
	INT_LITERAL
	FLOAT_LITERAL
	STRING_LITERAL
	CHAR_LITERAL
	NIL_LITERAL // **Reserved literal**
	// 3. Keywords
	LET
	IF
	ELSE
	LOOP
	IN
	BREAK
	CONTINUE
	RETURN
	FN
	STATIC
	CONST
	MATCH
	TRUE
	FALSE
	IMPORT
	// 4. Type keywords
	I8
	I16
	I32
	I64
	U8
	U16
	U32
	U64
	F32
	F64
	CHAR
	STRING
	BOOL
	STRUCT
	ENUM
	UNION
	// 5. Operators – Arithmetic
	PLUS
	MINUS
	STAR
	SLASH
	PERCENT
	// 6. Operators – Assignment
	EQUAL
	PLUS_EQUAL
	MINUS_EQUAL
	STAR_EQUAL
	SLASH_EQUAL
	PERCENT_EQUAL
	// 7. Operators – Comparison
	EQUAL_EQUAL
	BANG_EQUAL
	LESS
	LESS_EQUAL
	GREATER
	GREATER_EQUAL
	// 8. Operators – Logical / Bitwise
	AND
	OR
	BANG
	BIT_AND
	BIT_OR
	BIT_XOR
	BIT_NOT
	LEFT_SHIFT
	RIGHT_SHIFT
	// 9. Delimiters / Punctuation
	OPEN_PAREN
	CLOSE_PAREN
	OPEN_BRACE
	CLOSE_BRACE
	OPEN_BRACKET
	CLOSE_BRACKET
	COMMA
	DOT
	QUESTION
	COLON
	DOT_DOT
	ARROW
)

// Structure that holds individual token data
type Token struct {
	Type   TokenType
	Lexeme string
	Line   uint
}
