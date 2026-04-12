package parser

import (
	"testing"

	"github.com/Rajdeep-Nemo/sugarglaze/internal/ast"
	"github.com/Rajdeep-Nemo/sugarglaze/internal/token"
)

// ============================================================
// Helpers
// ============================================================

// checkParserErrors fails the test if the parser collected any errors.
func checkParserErrors(t *testing.T, p *Parser) {
	t.Helper()
	for _, err := range p.Errors() {
		t.Errorf("parser error: %s", err)
	}
}

// checkNoParserErrors is an alias for clarity in negative tests.
func expectParserErrors(t *testing.T, p *Parser, count int) {
	t.Helper()
	if len(p.Errors()) != count {
		t.Errorf("expected %d parser error(s), got %d: %v", count, len(p.Errors()), p.Errors())
	}
}

// makeProgram builds a Parser from a token slice and returns the parsed Program.
func makeProgram(tokens []token.Token) (*Parser, *ast.Program) {
	p := New(tokens)
	program := p.ParseProgram()
	return p, program
}

// checkStatementCount fails the test if the program doesn't have exactly n statements.
func checkStatementCount(t *testing.T, program *ast.Program, n int) {
	t.Helper()
	if len(program.Statements) != n {
		t.Fatalf("expected %d statement(s), got %d", n, len(program.Statements))
	}
}

// eof is a convenience token for END_OF_FILE.
var eof = token.Token{Type: token.END_OF_FILE, Lexeme: ""}

// ============================================================
// Let Statement Tests
// ============================================================

// let x = 5  (inferred type, initialized)
func TestLetInferredType(t *testing.T) {
	tokens := []token.Token{
		{Type: token.LET, Lexeme: "let"},
		{Type: token.IDENTIFIER, Lexeme: "x"},
		{Type: token.EQUAL, Lexeme: "="},
		{Type: token.INT_LITERAL, Lexeme: "5"},
		eof,
	}
	p, program := makeProgram(tokens)
	checkParserErrors(t, p)
	checkStatementCount(t, program, 1)

	stmt, ok := program.Statements[0].(*ast.LetStatement)
	if !ok {
		t.Fatalf("expected *ast.LetStatement, got %T", program.Statements[0])
	}
	if stmt.Name.Value != "x" {
		t.Errorf("expected name 'x', got %q", stmt.Name.Value)
	}
	if stmt.TypeHint != nil {
		t.Errorf("expected no type hint, got %q", stmt.TypeHint.Name)
	}
	if stmt.Value == nil {
		t.Errorf("expected a value, got nil")
	}
}

// let x: i32 = 5  (explicit type, initialized)
func TestLetExplicitType(t *testing.T) {
	tokens := []token.Token{
		{Type: token.LET, Lexeme: "let"},
		{Type: token.IDENTIFIER, Lexeme: "x"},
		{Type: token.COLON, Lexeme: ":"},
		{Type: token.I32, Lexeme: "i32"},
		{Type: token.EQUAL, Lexeme: "="},
		{Type: token.INT_LITERAL, Lexeme: "5"},
		eof,
	}
	p, program := makeProgram(tokens)
	checkParserErrors(t, p)
	checkStatementCount(t, program, 1)

	stmt, ok := program.Statements[0].(*ast.LetStatement)
	if !ok {
		t.Fatalf("expected *ast.LetStatement, got %T", program.Statements[0])
	}
	if stmt.TypeHint == nil {
		t.Fatal("expected a type hint, got nil")
	}
	if stmt.TypeHint.Name != "i32" {
		t.Errorf("expected type 'i32', got %q", stmt.TypeHint.Name)
	}
	if stmt.Value == nil {
		t.Errorf("expected a value, got nil")
	}
}

// let x: i32  (uninitialized, type provided — valid)
func TestLetUninitializedWithType(t *testing.T) {
	tokens := []token.Token{
		{Type: token.LET, Lexeme: "let"},
		{Type: token.IDENTIFIER, Lexeme: "x"},
		{Type: token.COLON, Lexeme: ":"},
		{Type: token.I32, Lexeme: "i32"},
		eof,
	}
	p, program := makeProgram(tokens)
	checkParserErrors(t, p)
	checkStatementCount(t, program, 1)

	stmt, ok := program.Statements[0].(*ast.LetStatement)
	if !ok {
		t.Fatalf("expected *ast.LetStatement, got %T", program.Statements[0])
	}
	if stmt.TypeHint == nil {
		t.Fatal("expected a type hint, got nil")
	}
	if stmt.Value != nil {
		t.Errorf("expected no value, got %v", stmt.Value)
	}
}

// let x  (uninitialized, no type — invalid, must error)
func TestLetUninitializedNoType(t *testing.T) {
	tokens := []token.Token{
		{Type: token.LET, Lexeme: "let"},
		{Type: token.IDENTIFIER, Lexeme: "x"},
		eof,
	}
	p, _ := makeProgram(tokens)
	expectParserErrors(t, p, 1)
}

// let  (missing identifier — invalid)
func TestLetMissingIdentifier(t *testing.T) {
	tokens := []token.Token{
		{Type: token.LET, Lexeme: "let"},
		{Type: token.EQUAL, Lexeme: "="},
		eof,
	}
	p, _ := makeProgram(tokens)
	expectParserErrors(t, p, 2)
}

// let x: <unknown type> = 5  (invalid type keyword)
func TestLetInvalidType(t *testing.T) {
	tokens := []token.Token{
		{Type: token.LET, Lexeme: "let"},
		{Type: token.IDENTIFIER, Lexeme: "x"},
		{Type: token.COLON, Lexeme: ":"},
		{Type: token.IDENTIFIER, Lexeme: "i65"}, // not a valid type keyword
		{Type: token.EQUAL, Lexeme: "="},
		{Type: token.INT_LITERAL, Lexeme: "5"},
		eof,
	}
	p, _ := makeProgram(tokens)
	expectParserErrors(t, p, 1)
}

// Test all supported type keywords are accepted in let statements
func TestLetAllTypeKeywords(t *testing.T) {
	typeTests := []struct {
		tokenType token.TokenType
		lexeme    string
		expected  string
	}{
		{token.I8, "i8", "i8"},
		{token.I16, "i16", "i16"},
		{token.I32, "i32", "i32"},
		{token.I64, "i64", "i64"},
		{token.U8, "u8", "u8"},
		{token.U16, "u16", "u16"},
		{token.U32, "u32", "u32"},
		{token.U64, "u64", "u64"},
		{token.F32, "f32", "f32"},
		{token.F64, "f64", "f64"},
		{token.CHAR, "char", "char"},
		{token.STRING, "string", "string"},
		{token.BOOL, "bool", "bool"},
	}

	for _, tt := range typeTests {
		t.Run(tt.lexeme, func(t *testing.T) {
			tokens := []token.Token{
				{Type: token.LET, Lexeme: "let"},
				{Type: token.IDENTIFIER, Lexeme: "x"},
				{Type: token.COLON, Lexeme: ":"},
				{Type: tt.tokenType, Lexeme: tt.lexeme},
				eof,
			}
			p, program := makeProgram(tokens)
			checkParserErrors(t, p)
			checkStatementCount(t, program, 1)

			stmt := program.Statements[0].(*ast.LetStatement)
			if stmt.TypeHint.Name != tt.expected {
				t.Errorf("expected type %q, got %q", tt.expected, stmt.TypeHint.Name)
			}
		})
	}
}

// ============================================================
// Const Statement Tests
// ============================================================

// const MAX = 100  (missing type — invalid, must error)
func TestConstMissingType(t *testing.T) {
	tokens := []token.Token{
		{Type: token.CONST, Lexeme: "const"},
		{Type: token.IDENTIFIER, Lexeme: "MAX"},
		{Type: token.EQUAL, Lexeme: "="},
		{Type: token.INT_LITERAL, Lexeme: "100"},
		eof,
	}

	p, _ := makeProgram(tokens)
	expectParserErrors(t, p, 2) // We expect this to fail now!
}

// const MAX: i32 = 100  (explicit type, initialized)
func TestConstExplicitType(t *testing.T) {
	tokens := []token.Token{
		{Type: token.CONST, Lexeme: "const"},
		{Type: token.IDENTIFIER, Lexeme: "MAX"},
		{Type: token.COLON, Lexeme: ":"},
		{Type: token.I32, Lexeme: "i32"},
		{Type: token.EQUAL, Lexeme: "="},
		{Type: token.INT_LITERAL, Lexeme: "100"},
		eof,
	}
	p, program := makeProgram(tokens)
	checkParserErrors(t, p)
	checkStatementCount(t, program, 1)

	stmt, ok := program.Statements[0].(*ast.ConstStatement)
	if !ok {
		t.Fatalf("expected *ast.ConstStatement, got %T", program.Statements[0])
	}
	if stmt.TypeHint == nil {
		t.Fatal("expected a type hint, got nil")
	}
	if stmt.TypeHint.Name != "i32" {
		t.Errorf("expected type 'i32', got %q", stmt.TypeHint.Name)
	}
}

// const MAX: i32  (no value — invalid, must error)
func TestConstUninitializedWithType(t *testing.T) {
	tokens := []token.Token{
		{Type: token.CONST, Lexeme: "const"},
		{Type: token.IDENTIFIER, Lexeme: "MAX"},
		{Type: token.COLON, Lexeme: ":"},
		{Type: token.I32, Lexeme: "i32"},
		eof,
	}
	p, _ := makeProgram(tokens)
	expectParserErrors(t, p, 1)
}

// const MAX  (no type, no value — invalid, must error)
func TestConstUninitializedNoType(t *testing.T) {
	tokens := []token.Token{
		{Type: token.CONST, Lexeme: "const"},
		{Type: token.IDENTIFIER, Lexeme: "MAX"},
		eof,
	}
	p, _ := makeProgram(tokens)
	expectParserErrors(t, p, 1)
}

// const  (missing identifier — invalid)
func TestConstMissingIdentifier(t *testing.T) {
	tokens := []token.Token{
		{Type: token.CONST, Lexeme: "const"},
		{Type: token.EQUAL, Lexeme: "="},
		eof,
	}
	p, _ := makeProgram(tokens)
	expectParserErrors(t, p, 2)
}

// ============================================================
// Assign Statement Tests
// ============================================================

// x = 10  (simple assignment)
func TestAssignSimple(t *testing.T) {
	tokens := []token.Token{
		{Type: token.IDENTIFIER, Lexeme: "x"},
		{Type: token.EQUAL, Lexeme: "="},
		{Type: token.INT_LITERAL, Lexeme: "10"},
		eof,
	}
	p, program := makeProgram(tokens)
	checkParserErrors(t, p)
	checkStatementCount(t, program, 1)

	stmt, ok := program.Statements[0].(*ast.AssignStatement)
	if !ok {
		t.Fatalf("expected *ast.AssignStatement, got %T", program.Statements[0])
	}
	if stmt.Target.Value != "x" {
		t.Errorf("expected target 'x', got %q", stmt.Target.Value)
	}
	if stmt.Operator != ast.Assign {
		t.Errorf("expected operator '=', got %q", stmt.Operator)
	}
}

// Test all compound assignment operators
func TestAssignCompoundOperators(t *testing.T) {
	opTests := []struct {
		tokenType token.TokenType
		lexeme    string
		expected  ast.AssignOperator
	}{
		{token.PLUS_EQUAL, "+=", ast.PlusAssign},
		{token.MINUS_EQUAL, "-=", ast.MinusAssign},
		{token.STAR_EQUAL, "*=", ast.StarAssign},
		{token.SLASH_EQUAL, "/=", ast.SlashAssign},
		{token.PERCENT_EQUAL, "%=", ast.PercentAssign},
	}

	for _, tt := range opTests {
		t.Run(tt.lexeme, func(t *testing.T) {
			tokens := []token.Token{
				{Type: token.IDENTIFIER, Lexeme: "x"},
				{Type: tt.tokenType, Lexeme: tt.lexeme},
				{Type: token.INT_LITERAL, Lexeme: "5"},
				eof,
			}
			p, program := makeProgram(tokens)
			checkParserErrors(t, p)
			checkStatementCount(t, program, 1)

			stmt, ok := program.Statements[0].(*ast.AssignStatement)
			if !ok {
				t.Fatalf("expected *ast.AssignStatement, got %T", program.Statements[0])
			}
			if stmt.Operator != tt.expected {
				t.Errorf("expected operator %q, got %q", tt.expected, stmt.Operator)
			}
		})
	}
}

// ============================================================
// String() / Pretty Print Tests
// ============================================================

func TestLetString(t *testing.T) {
	stmt := &ast.LetStatement{
		Token: token.Token{Type: token.LET, Lexeme: "let"},
		Name:  &ast.Identifier{Token: token.Token{Type: token.IDENTIFIER, Lexeme: "x"}, Value: "x"},
		TypeHint: &ast.TypeAnnotation{
			Token: token.Token{Type: token.I32, Lexeme: "i32"},
			Name:  "i32",
		},
		Value: nil,
	}
	expected := "let x: i32"
	if stmt.String() != expected {
		t.Errorf("expected %q, got %q", expected, stmt.String())
	}
}

func TestConstString(t *testing.T) {
	stmt := &ast.ConstStatement{
		Token:    token.Token{Type: token.CONST, Lexeme: "const"},
		Name:     &ast.Identifier{Token: token.Token{Type: token.IDENTIFIER, Lexeme: "MAX"}, Value: "MAX"},
		TypeHint: nil,
		Value:    &ast.Identifier{Token: token.Token{Type: token.IDENTIFIER, Lexeme: "100"}, Value: "100"},
	}
	expected := "const MAX = 100"
	if stmt.String() != expected {
		t.Errorf("expected %q, got %q", expected, stmt.String())
	}
}

func TestAssignString(t *testing.T) {
	stmt := &ast.AssignStatement{
		Token:    token.Token{Type: token.PLUS_EQUAL, Lexeme: "+="},
		Target:   &ast.Identifier{Token: token.Token{Type: token.IDENTIFIER, Lexeme: "x"}, Value: "x"},
		Operator: ast.PlusAssign,
		Value:    &ast.Identifier{Token: token.Token{Type: token.IDENTIFIER, Lexeme: "5"}, Value: "5"},
	}
	expected := "x += 5"
	if stmt.String() != expected {
		t.Errorf("expected %q, got %q", expected, stmt.String())
	}
}

// ============================================================
// Multi-statement Tests
// ============================================================

// Parse multiple statements in sequence
// Parse multiple statements in sequence
func TestMultipleStatements(t *testing.T) {
	tokens := []token.Token{
		// let x: i32 = 5
		{Type: token.LET, Lexeme: "let"},
		{Type: token.IDENTIFIER, Lexeme: "x"},
		{Type: token.COLON, Lexeme: ":"},
		{Type: token.I32, Lexeme: "i32"},
		{Type: token.EQUAL, Lexeme: "="},
		{Type: token.INT_LITERAL, Lexeme: "5"},

		// const MAX: i32 = 100  <-- FIXED: Added type annotation here
		{Type: token.CONST, Lexeme: "const"},
		{Type: token.IDENTIFIER, Lexeme: "MAX"},
		{Type: token.COLON, Lexeme: ":"}, // ADDED
		{Type: token.I32, Lexeme: "i32"}, // ADDED
		{Type: token.EQUAL, Lexeme: "="},
		{Type: token.INT_LITERAL, Lexeme: "100"},

		// x += 1
		{Type: token.IDENTIFIER, Lexeme: "x"},
		{Type: token.PLUS_EQUAL, Lexeme: "+="},
		{Type: token.INT_LITERAL, Lexeme: "1"},
		eof,
	}
	p, program := makeProgram(tokens)
	checkParserErrors(t, p)
	checkStatementCount(t, program, 3)

	// first statement is a LetStatement
	if _, ok := program.Statements[0].(*ast.LetStatement); !ok {
		t.Errorf("expected statement 0 to be *ast.LetStatement, got %T", program.Statements[0])
	}
	// second is a ConstStatement
	if _, ok := program.Statements[1].(*ast.ConstStatement); !ok {
		t.Errorf("expected statement 1 to be *ast.ConstStatement, got %T", program.Statements[1])
	}
	// third is an AssignStatement
	if _, ok := program.Statements[2].(*ast.AssignStatement); !ok {
		t.Errorf("expected statement 2 to be *ast.AssignStatement, got %T", program.Statements[2])
	}
}

// Test import statements
func TestImportStatement(t *testing.T) {
	tests := []struct {
		name            string
		inputTokens     []token.Token
		expectedModules []string
	}{
		{
			name: "Single Import",
			// import "io";
			inputTokens: []token.Token{
				{Type: token.IMPORT, Lexeme: "import"},
				{Type: token.STRING_LITERAL, Lexeme: "\"io\""}, // Lexer includes quotes
				eof,
			},
			expectedModules: []string{"io"}, // Note: parseStringLiteral strips the quotes for the AST Value
		},
		{
			name: "Grouped Imports",
			// import ("io", "math");
			inputTokens: []token.Token{
				{Type: token.IMPORT, Lexeme: "import"},
				{Type: token.OPEN_PAREN, Lexeme: "("},
				{Type: token.STRING_LITERAL, Lexeme: "\"io\""},
				{Type: token.COMMA, Lexeme: ","},
				{Type: token.STRING_LITERAL, Lexeme: "\"math\""},
				{Type: token.CLOSE_PAREN, Lexeme: ")"},
				eof,
			},
			expectedModules: []string{"io", "math"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, program := makeProgram(tt.inputTokens)

			// We shouldn't have any parsing errors
			checkParserErrors(t, p)

			// There should be exactly 1 statement parsed
			checkStatementCount(t, program, 1)

			// 1. Check if the statement is an ImportStatement
			stmt, ok := program.Statements[0].(*ast.ImportStatement)
			if !ok {
				t.Fatalf("expected *ast.ImportStatement, got %T", program.Statements[0])
			}

			// 2. Check that the correct number of modules were parsed
			if len(stmt.Modules) != len(tt.expectedModules) {
				t.Fatalf("expected %d modules, got %d", len(tt.expectedModules), len(stmt.Modules))
			}

			// 3. Verify the names of the imported modules
			for i, expected := range tt.expectedModules {
				actual := stmt.Modules[i].Value
				if actual != expected {
					t.Errorf("expected module %q, got %q", expected, actual)
				}
			}
		})
	}
}

// Empty program — should parse fine with zero statements
func TestEmptyProgram(t *testing.T) {
	tokens := []token.Token{eof}
	p, program := makeProgram(tokens)
	checkParserErrors(t, p)
	checkStatementCount(t, program, 0)
}

// Test comprehensive call expressions in both let and expression statements, with various argument types and nesting levels
func TestComprehensiveCallExpressions(t *testing.T) {
	tests := []struct {
		name         string
		inputTokens  []token.Token
		isLet        bool     // Is this a Let statement (true) or an Expression statement (false)?
		expectedFunc string   // The name of the function being called
		expectedArgs []string // The string representation of each argument
	}{
		{
			name: "Standalone with no arguments",
			// print();
			inputTokens: []token.Token{
				{Type: token.IDENTIFIER, Lexeme: "print"},
				{Type: token.OPEN_PAREN, Lexeme: "("},
				{Type: token.CLOSE_PAREN, Lexeme: ")"},
				eof,
			},
			isLet:        false,
			expectedFunc: "print",
			expectedArgs: []string{},
		},
		{
			name: "Standalone with format string and variable",
			// print("Val: {i32}", x);
			inputTokens: []token.Token{
				{Type: token.IDENTIFIER, Lexeme: "print"},
				{Type: token.OPEN_PAREN, Lexeme: "("},
				{Type: token.STRING_LITERAL, Lexeme: "\"Val: {i32}\""},
				{Type: token.COMMA, Lexeme: ","},
				{Type: token.IDENTIFIER, Lexeme: "x"},
				{Type: token.CLOSE_PAREN, Lexeme: ")"},
				eof,
			},
			isLet:        false,
			expectedFunc: "print",
			expectedArgs: []string{"\"Val: {i32}\"", "x"},
		},
		{
			name: "Call inside a let declaration",
			// let result = add(5, 10);
			inputTokens: []token.Token{
				{Type: token.LET, Lexeme: "let"},
				{Type: token.IDENTIFIER, Lexeme: "result"},
				{Type: token.EQUAL, Lexeme: "="},
				{Type: token.IDENTIFIER, Lexeme: "add"},
				{Type: token.OPEN_PAREN, Lexeme: "("},
				{Type: token.INT_LITERAL, Lexeme: "5"},
				{Type: token.COMMA, Lexeme: ","},
				{Type: token.INT_LITERAL, Lexeme: "10"},
				{Type: token.CLOSE_PAREN, Lexeme: ")"},
				eof,
			},
			isLet:        true,
			expectedFunc: "add",
			expectedArgs: []string{"5", "10"},
		},
		{
			name: "Nested function call",
			// print(get_value());
			inputTokens: []token.Token{
				{Type: token.IDENTIFIER, Lexeme: "print"},
				{Type: token.OPEN_PAREN, Lexeme: "("},
				{Type: token.IDENTIFIER, Lexeme: "get_value"},
				{Type: token.OPEN_PAREN, Lexeme: "("},
				{Type: token.CLOSE_PAREN, Lexeme: ")"},
				{Type: token.CLOSE_PAREN, Lexeme: ")"},
				eof,
			},
			isLet:        false,
			expectedFunc: "print",
			expectedArgs: []string{"get_value()"}, // The inner call expression stringifies to this
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, program := makeProgram(tt.inputTokens)
			checkParserErrors(t, p)
			checkStatementCount(t, program, 1)

			var callExp *ast.CallExpression

			// Extract the CallExpression based on the statement type
			if tt.isLet {
				stmt, ok := program.Statements[0].(*ast.LetStatement)
				if !ok {
					t.Fatalf("expected *ast.LetStatement, got %T", program.Statements[0])
				}
				callExp, ok = stmt.Value.(*ast.CallExpression)
				if !ok {
					t.Fatalf("expected stmt.Value to be *ast.CallExpression, got %T", stmt.Value)
				}
			} else {
				stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
				if !ok {
					t.Fatalf("expected *ast.ExpressionStatement, got %T", program.Statements[0])
				}
				callExp, ok = stmt.Expression.(*ast.CallExpression)
				if !ok {
					t.Fatalf("expected stmt.Expression to be *ast.CallExpression, got %T", stmt.Expression)
				}
			}

			// 1. Check the function name
			if callExp.Function.String() != tt.expectedFunc {
				t.Errorf("expected function name %q, got %q", tt.expectedFunc, callExp.Function.String())
			}

			// 2. Check the arguments length
			if len(callExp.Arguments) != len(tt.expectedArgs) {
				t.Fatalf("expected %d arguments, got %d", len(tt.expectedArgs), len(callExp.Arguments))
			}

			// 3. Check each argument's string representation
			for i, arg := range callExp.Arguments {
				if arg.String() != tt.expectedArgs[i] {
					t.Errorf("argument %d: expected %q, got %q", i, tt.expectedArgs[i], arg.String())
				}
			}
		})
	}
}

// ============================================================
// Control Flow Tests (If / Else / Blocks)
// ============================================================

func TestIfElseExpression(t *testing.T) {
	// if (x < y) { x } else { y }
	tokens := []token.Token{
		{Type: token.IF, Lexeme: "if"},
		{Type: token.OPEN_PAREN, Lexeme: "("},
		{Type: token.IDENTIFIER, Lexeme: "x"},
		{Type: token.LESS, Lexeme: "<"},
		{Type: token.IDENTIFIER, Lexeme: "y"},
		{Type: token.CLOSE_PAREN, Lexeme: ")"},
		{Type: token.OPEN_BRACE, Lexeme: "{"},
		{Type: token.IDENTIFIER, Lexeme: "x"},
		{Type: token.CLOSE_BRACE, Lexeme: "}"},
		{Type: token.ELSE, Lexeme: "else"},
		{Type: token.OPEN_BRACE, Lexeme: "{"},
		{Type: token.IDENTIFIER, Lexeme: "y"},
		{Type: token.CLOSE_BRACE, Lexeme: "}"},
		eof,
	}

	p, program := makeProgram(tokens)
	checkParserErrors(t, p)
	checkStatementCount(t, program, 1)

	// The if/else should be wrapped in an ExpressionStatement
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("expected *ast.ExpressionStatement, got %T", program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("expected stmt.Expression to be *ast.IfExpression, got %T", stmt.Expression)
	}

	// Check condition parsing
	if exp.Condition == nil {
		t.Fatalf("expected condition, got nil")
	}

	// Check Consequence Block { x }
	if len(exp.Consequence.Statements) != 1 {
		t.Errorf("consequence is not 1 statement. got=%d\n", len(exp.Consequence.Statements))
	}

	// Check Alternative Block { y }
	if exp.Alternative == nil {
		t.Fatalf("expected alternative block, got nil")
	}
	if len(exp.Alternative.Statements) != 1 {
		t.Errorf("alternative is not 1 statement. got=%d\n", len(exp.Alternative.Statements))
	}
}
