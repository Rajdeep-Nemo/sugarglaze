package evaluator

import (
	"fmt"
	"testing"

	"github.com/Rajdeep-Nemo/sugarglaze/internal/lexer"
	"github.com/Rajdeep-Nemo/sugarglaze/internal/object"
	"github.com/Rajdeep-Nemo/sugarglaze/internal/parser"
	"github.com/Rajdeep-Nemo/sugarglaze/internal/token"
)

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
		objType  object.ObjectType
	}{
		{"5", 5, object.I32_OBJ},
		{"10", 10, object.I32_OBJ},
		{"-5", -5, object.I32_OBJ},
		{"-10", -10, object.I32_OBJ},
		{"5 + 5 + 5 + 5 - 10", 10, object.I32_OBJ},
		{"2 * 2 * 2 * 2 * 2", 32, object.I32_OBJ},
		{"-50 + 100 + -50", 0, object.I32_OBJ},
		{"5 * 2 + 10", 20, object.I32_OBJ},
		{"5 + 2 * 10", 25, object.I32_OBJ},
		{"20 + 2 * -10", 0, object.I32_OBJ},
		{"50 / 2 * 2 + 10", 60, object.I32_OBJ},
		{"2 * (5 + 10)", 30, object.I32_OBJ},
		{"3 * 3 * 3 + 10", 37, object.I32_OBJ},
		{"3 * (3 * 3) + 10", 37, object.I32_OBJ},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50, object.I32_OBJ},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected, tt.objType)
	}
}

func TestEvalFloatExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
		objType  object.ObjectType
	}{
		{"5.5", 5.5, object.F32_OBJ},
		{"-5.5", -5.5, object.F32_OBJ},
		{"5.0 + 5.5", 10.5, object.F32_OBJ},
		{"10.5 - 5.0", 5.5, object.F32_OBJ},
		{"2.0 * 2.5", 5.0, object.F32_OBJ},
		{"10.0 / 2.0", 5.0, object.F32_OBJ},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testFloatObject(t, evaluated, tt.expected, tt.objType)
	}
}

func TestLetStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
		objType  object.ObjectType
	}{
		{"let a = 5 \n a", 5, object.I32_OBJ},
		{"let a = 5 * 5 \n a", 25, object.I32_OBJ},
		{"let a = 5 \n let b = a \n b", 5, object.I32_OBJ},
		{"let a = 5 \n let b = a \n let c = a + b + 5 \n c", 15, object.I32_OBJ},
		{"let a: i8 = 10 \n a", 10, object.I8_OBJ},
		{"let a: i16 = 200 \n a", 200, object.I16_OBJ},
		{"let a: i64 = 50 \n a", 50, object.I64_OBJ},
		{"let a: u8 = 255 \n a", 255, object.U8_OBJ},
		{"let a: f64 = 5 \n a", 5, object.F64_OBJ},
	}

	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected, tt.objType)
	}
}

func TestConstStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
		objType  object.ObjectType
	}{
		// FIX: Added the mandatory type hints for const to pass the parser check
		{"const a: i32 = 5 \n a", 5, object.I32_OBJ},
		{"const a: i8 = 10 \n a", 10, object.I8_OBJ},
	}

	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected, tt.objType)
	}
}

func TestAssignStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
		objType  object.ObjectType
	}{
		{"let a = 5 \n a = 10 \n a", 10, object.I32_OBJ},
		{"let a = 5 \n a += 10 \n a", 15, object.I32_OBJ},
		{"let a = 10 \n a -= 5 \n a", 5, object.I32_OBJ},
		{"let a = 5 \n a *= 2 \n a", 10, object.I32_OBJ},
		{"let a = 10 \n a /= 2 \n a", 5, object.I32_OBJ},
		{"let a: i8 = 5 \n a = 10 \n a", 10, object.I8_OBJ},
		{"let a: i8 = 5 \n a += 10 \n a", 15, object.I8_OBJ},
	}

	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected, tt.objType)
	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input           string
		expectedMessage string
	}{
		{
			"let a: i8 = 300",
			"value 300 out of bounds for i8",
		},
		{
			"let a: u8 = -5",
			"value -5 out of bounds for u8",
		},
		{
			"let a: i32 \n a",
			"cannot access uninitialized variable 'a'",
		},
		{
			"const a: i32 = 5 \n a = 10",
			"cannot reassign to const variable 'a'",
		},
		{
			"let a: i32 = 5 \n a = 5.5",
			"type mismatch: cannot assign 'f32' to variable 'a' (expected 'i32')",
		},
		{
			"let a: f32 = 5.5 \n let b: i32 = a",
			// FIX: Now uses the professional variable-context error string!
			"type mismatch: cannot assign 'f32' to variable 'b' (expected 'i32')",
		},
		{
			"let a: i8 = 100 \n a += 100",
			"value 200 out of bounds for i8",
		},
		{
			"let a: u8 = 5 \n a -= 10",
			"value -5 out of bounds for u8",
		},
		{
			"-true",
			"unknown operator: -bool",
		},
		{
			"let a: u32 = 10 \n let b = -a",
			"cannot negate unsigned integer type u32",
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		if evaluated == nil {
			t.Errorf("expected error object, got nil for input: %q (Check parser for syntax errors)", tt.input)
			continue
		}

		errObj, ok := evaluated.(*object.Error)
		if !ok {
			t.Errorf("no error object returned. got=%T(%v) for input: %q",
				evaluated, evaluated, tt.input)
			continue
		}

		if errObj.Message != tt.expectedMessage {
			t.Errorf("wrong error message. expected=%q, got=%q",
				tt.expectedMessage, errObj.Message)
		}
	}
}

func testEval(input string) object.Object {
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

	if program == nil || len(program.Statements) == 0 {
		fmt.Printf("Warning: Parser returned empty program for input: %q\n", input)
		return nil
	}

	env := object.NewEnvironment()
	return Eval(program, env)
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64, expectedType object.ObjectType) bool {
	if obj == nil {
		t.Errorf("object is nil. Expected %d (%s)", expected, expectedType)
		return false
	}

	if obj.Type() != expectedType {
		t.Errorf("object is not %s. got=%T (%s)", expectedType, obj, obj.Type())
		return false
	}

	var val int64
	switch v := obj.(type) {
	case *object.Int8:
		val = int64(v.Value)
	case *object.Int16:
		val = int64(v.Value)
	case *object.Int32:
		val = int64(v.Value)
	case *object.Int64:
		val = v.Value
	case *object.Uint8:
		val = int64(v.Value)
	case *object.Uint16:
		val = int64(v.Value)
	case *object.Uint32:
		val = int64(v.Value)
	case *object.Uint64:
		val = int64(v.Value)
	case *object.Float64:
		val = int64(v.Value)
	default:
		t.Errorf("object is not an integer type. got=%T", obj)
		return false
	}

	if val != expected {
		t.Errorf("object has wrong value. got=%d, want=%d", val, expected)
		return false
	}

	return true
}

func testFloatObject(t *testing.T, obj object.Object, expected float64, expectedType object.ObjectType) bool {
	if obj == nil {
		t.Errorf("object is nil. Expected %f (%s)", expected, expectedType)
		return false
	}

	if obj.Type() != expectedType {
		t.Errorf("object is not %s. got=%T (%s)", expectedType, obj, obj.Type())
		return false
	}

	var val float64
	switch v := obj.(type) {
	case *object.Float32:
		val = float64(v.Value)
	case *object.Float64:
		val = v.Value
	default:
		t.Errorf("object is not a float type. got=%T", obj)
		return false
	}

	if val != expected {
		t.Errorf("object has wrong value. got=%f, want=%f", val, expected)
		return false
	}

	return true
}
