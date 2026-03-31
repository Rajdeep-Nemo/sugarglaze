package evaluator

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"glaze/internal/lexer"
	"glaze/internal/object"
	"glaze/internal/parser"
	"glaze/internal/token"
)

// Helper function to process code and return the final evaluated object
func testEval(t *testing.T, input string) object.Object {
	s := lexer.InitScanner(input)
	var tokens []token.Token
	for {
		tok := s.ScanToken()
		tokens = append(tokens, tok)
		if tok.Type == token.END_OF_FILE {
			break
		}
	}

	p := parser.New(tokens)
	program := p.ParseProgram()

	if len(p.Errors()) != 0 {
		t.Fatalf("Parser encountered errors in input: %q\n%v", input, p.Errors())
	}

	env := object.NewEnvironment()
	return Eval(program, env)
}

// Test variable and typing
func TestVariablesAndStrictTyping(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{} // Can be an int, string, or an Error message substring
	}{
		// --- Success Cases ---
		{"let a = 5; a;", int64(5)},
		{"let a: i32 = 10; a;", int64(10)},
		{"let a = 5; let b = a; b;", int64(5)},
		{"let a = 5; a = 10; a;", int64(10)},
		{"let a = 10; a += 5; a;", int64(15)},
		{"const PI = 3; PI;", int64(3)},

		// --- Uninitialized Defaults ---
		{"let empty: string; empty;", nil}, // Should be NIL object

		// --- Error Cases (Should Fail Fast!) ---
		{"let a: string = 5;", "type mismatch"},
		{"let a = 5; a = \"hello\";", "type mismatch"},
		{"const a = 5; a = 10;", "cannot reassign to const variable"},
		{"let empty: i32; empty += 5;", "type mismatch"}, // Math on NIL should crash
		{"unknownVar = 10;", "cannot assign to undefined variable"},
	}

	for _, tt := range tests {
		evaluated := testEval(t, tt.input)

		switch expected := tt.expected.(type) {
		case int64:
			testIntegerObject(t, evaluated, expected)
		case string:
			errObj, ok := evaluated.(*object.Error)
			if !ok {
				t.Errorf("Expected an error for input %q, got %T (%+v)", tt.input, evaluated, evaluated)
				continue
			}
			if !strings.Contains(errObj.Message, expected) {
				t.Errorf("Expected error to contain %q, got %q", expected, errObj.Message)
			}
		case nil:
			if evaluated.Type() != object.NIL_OBJ {
				t.Errorf("Expected NIL for input %q, got %T", tt.input, evaluated)
			}
		}
	}
}

// Small helper to verify integers
func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	result, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("object is not Integer. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%d, want=%d", result.Value, expected)
		return false
	}
	return true
}

// Test the built-in print and println functions, including format string handling and error cases
func TestBuiltinPrintFormatting(t *testing.T) {
	tests := []struct {
		input          string
		expectedOutput string
		expectedError  string
	}{
		// --- Success Cases ---
		{`print("Hello");`, "Hello", ""},
		{`println("Hello " + "World");`, "Hello World\n", ""},
		{`let age = 25; print("Age: {i32}", age);`, "Age: 25", ""},
		{`let n = "Pluesi"; let v = 1; println("Lang: {string}, Ver: {i32}", n, v);`, "Lang: Pluesi, Ver: 1\n", ""},

		// --- Error Cases (Format Engine) ---
		{`print("Age: {i32}", "twenty");`, "", "type mismatch in print"},
		{`print("Score: {f64}", 100);`, "", "type mismatch in print"}, // strict float vs int check!
		{`print("Name: {string} {string}", "Pluesi");`, "", "not enough arguments"},
		{`print("Name: {string}", "Pluesi", "Extra");`, "", "too many arguments"},
	}

	for _, tt := range tests {
		// Hijack os.Stdout to capture printed output
		oldStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		// Run the evaluator
		evaluated := testEval(t, tt.input)

		// Restore os.Stdout
		w.Close()
		os.Stdout = oldStdout
		var buf bytes.Buffer
		io.Copy(&buf, r)
		actualOutput := buf.String()

		// 1. Check Output
		if actualOutput != tt.expectedOutput {
			t.Errorf("Input: %q\nExpected output: %q\nGot: %q", tt.input, tt.expectedOutput, actualOutput)
		}

		// 2. Check Errors
		if tt.expectedError != "" {
			errObj, ok := evaluated.(*object.Error)
			if !ok {
				t.Errorf("Input: %q\nExpected error containing %q, but got no error! Result: %s", tt.input, tt.expectedError, evaluated.Inspect())
				continue
			}
			if !strings.Contains(errObj.Message, tt.expectedError) {
				t.Errorf("Input: %q\nExpected error to contain %q, got %q", tt.input, tt.expectedError, errObj.Message)
			}
		}
	}
}

// Test if-else statements and block scoping, including variable mutation and shadowing
func TestIfElseAndBlockScoping(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		// --- Basic Truthiness ---
		{"if (true) { 10 }", int64(10)},
		{"if (false) { 10 }", nil},
		{"if (1) { 10 }", int64(10)}, // 1 is truthy
		{"if (false) { 10 } else { 20 }", int64(20)},

		// --- Math Conditions ---
		{"if (5 < 10) { 10 } else { 20 }", int64(10)},
		{"if (5 > 10) { 10 } else { 20 }", int64(20)},
		{"if (10 == 10) { 10 } else { 20 }", int64(10)},

		// --- SCOPING: Mutating Outer Variables ---
		{`
			let x = 10;
			if (true) { x = 20; }
			x;
		`, int64(20)}, // The outer x should be updated

		// --- SCOPING: Shadowing (Local variables don't leak) ---
		{`
			let x = 10;
			if (true) { let x = 50; }
			x;
		`, int64(10)}, // Outer x is untouched because 'let' created a local block variable

		// --- Nested Blocks ---
		{`
			if (true) {
				if (true) {
					100
				}
			}
		`, int64(100)},
	}

	for _, tt := range tests {
		evaluated := testEval(t, tt.input)
		integer, ok := tt.expected.(int64)
		if ok {
			testIntegerObject(t, evaluated, integer)
		} else {
			if evaluated != nil && evaluated.Type() != object.NIL_OBJ {
				t.Errorf("Expected NIL/nil, got %T (%+v)", evaluated, evaluated)
			}
		}
	}
}
