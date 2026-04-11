package evaluator

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/Rajdeep-Nemo/sugarglaze/internal/object"
)

// Handles the string formatting, interpolation, and strict type checking.
// It returns the final string to print, or a runtime Error if types mismatch.
func formatOutput(args ...object.Object) (string, *object.Error) {
	if len(args) == 0 {
		return "", nil
	}

	firstArg, isString := args[0].(*object.String)

	// If the first argument isn't a string, or it doesn't contain '{',
	// just squish all arguments together without spaces
	if !isString || !strings.Contains(firstArg.Value, "{") {
		var out strings.Builder
		for _, arg := range args {
			out.WriteString(arg.Inspect())
		}
		return out.String(), nil
	}

	// Parses the format string
	format := firstArg.Value
	var out strings.Builder
	argIndex := 1 // Variables to inject start at args[1]

	i := 0
	for i < len(format) {
		if format[i] == '{' {
			// Find the matching '}'
			closeIdx := strings.IndexByte(format[i:], '}')
			if closeIdx != -1 {
				// Extract the type expected by the placeholder (e.g., "i32" from "{i32}")
				expectedType := format[i+1 : i+closeIdx]

				// Check if enough arguments are provided for the placeholders
				if argIndex >= len(args) {
					return "", &object.Error{Message: "not enough arguments for format string"}
				}

				providedArg := args[argIndex]

				// STRICT TYPE CHECKING: Compare the placeholder to the object's actual type!
				if string(providedArg.Type()) != expectedType {
					return "", &object.Error{
						Message: fmt.Sprintf("type mismatch in print: expected '{%s}', got '%s'", expectedType, providedArg.Type()),
					}
				}

				// If it passes, write the value and move to the next argument
				out.WriteString(providedArg.Inspect())
				argIndex++
				i += closeIdx + 1 // Jump past the '}'
				continue
			}
		}
		// Write normal characters
		out.WriteByte(format[i])
		i++
	}

	// Check if too many arguments are provided for the placeholders
	if argIndex < len(args) {
		return "", &object.Error{Message: "too many arguments for format string"}
	}

	return out.String(), nil
}

// Prints exact string representations
func builtinPrint(args ...object.Object) object.Object {
	output, err := formatOutput(args...)
	if err != nil {
		return err
	}
	fmt.Print(output)
	return &object.Nil{}
}

// Prints with a trailing newline
func builtinPrintln(args ...object.Object) object.Object {
	output, err := formatOutput(args...)
	if err != nil {
		return err
	}
	fmt.Println(output)
	return &object.Nil{}
}

// Trims leading and trailing whitespace from a string
func trim(s string) string {
	return strings.TrimSpace(s)
}

// builtinTrim wraps the native trim function for the evaluator
func builtinTrim(args ...object.Object) object.Object {
	if len(args) != 1 {
		return &object.Error{
			Message: fmt.Sprintf("wrong number of arguments for trim(). got=%d, want=1", len(args)),
		}
	}
	strObj, ok := args[0].(*object.String)
	if !ok {
		return &object.Error{
			Message: fmt.Sprintf("argument to trim() must be a string, got %s", args[0].Type()),
		}
	}
	result := trim(strObj.Value)
	return &object.String{Value: result}
}

// A shared buffered reader for os.Stdin.
var reader = bufio.NewReader(os.Stdin)

// Input reads a line from stdin with an optional prompt.
func Input(prompt string) string {
	if prompt != "" {
		fmt.Print(prompt)
	}
	line, err := reader.ReadString('\n')
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to read input.")
		os.Exit(1)
	}
	return strings.TrimRight(line, "\r\n")
}

// builtinInput wraps the native Input function for the evaluator
func builtinInput(args ...object.Object) object.Object {
	if len(args) > 1 {
		return &object.Error{
			Message: fmt.Sprintf("wrong number of arguments for input(). got=%d, want=0 or 1", len(args)),
		}
	}
	prompt := ""
	if len(args) == 1 {
		strObj, ok := args[0].(*object.String)
		if !ok {
			return &object.Error{
				Message: fmt.Sprintf("argument to input() must be a string, got %s", args[0].Type()),
			}
		}
		prompt = strObj.Value
	}
	result := Input(prompt)
	return &object.String{Value: result}
}

// Parses a string to a boolean value.
// Returns nil if parsing fails.
func ParseBool(s string) *bool {
	s = trim(s)
	if s == "true" {
		v := true
		return &v
	}
	if s == "false" {
		v := false
		return &v
	}
	return nil
}

// Parses a string to a boolean value.
// Exits with an error message if parsing fails.
func MustParseBool(s string) bool {
	v := ParseBool(s)
	if v == nil {
		fmt.Fprintf(os.Stderr, "Invalid value: '%s', value accepted: true / false\n", s)
		os.Exit(1)
	}
	return *v
}

// Parses a string to a single byte value.
// Returns nil if parsing fails.
func ParseChar(s string) *byte {
	s = trim(s)
	if len(s) == 0 {
		return nil
	}
	if len(s) > 1 {
		return nil
	}
	v := s[0]
	return &v
}

// Parses a string to a single byte value.
// Exits with an error message if parsing fails.
func MustParseChar(s string) byte {
	s = trim(s)
	if len(s) == 0 {
		fmt.Fprintf(os.Stderr, "Empty char literal.\n")
		os.Exit(1)
	}
	if len(s) > 1 {
		fmt.Fprintf(os.Stderr, "Multiple characters found: '%s'\n", s)
		os.Exit(1)
	}
	v := ParseChar(s)
	return *v
}

// ParseU8 parses a string to an 8-bit unsigned integer.
// Returns nil if parsing fails.
func ParseU8(s string) *uint8 {
	s = trim(s)
	if s == "" {
		return nil
	}
	v, err := strconv.ParseUint(s, 10, 8)
	if err != nil {
		return nil
	}
	val := uint8(v)
	return &val
}

// MustParseU8 parses a string to an 8-bit unsigned integer.
// Exits with an error message if parsing fails.
func MustParseU8(s string) uint8 {
	s = trim(s)
	if s == "" {
		fmt.Fprintln(os.Stderr, "Empty input.")
		os.Exit(1)
	}
	v := ParseU8(s)
	if v == nil {
		fmt.Fprintf(os.Stderr, "Invalid input: '%s' is not a valid 8-bit unsigned integer.\n", s)
		os.Exit(1)
	}
	return *v
}

// ParseU16 parses a string to a 16-bit unsigned integer.
// Returns nil if parsing fails.
func ParseU16(s string) *uint16 {
	s = trim(s)
	if s == "" {
		return nil
	}
	v, err := strconv.ParseUint(s, 10, 16)
	if err != nil {
		return nil
	}
	val := uint16(v)
	return &val
}

// MustParseU16 parses a string to a 16-bit unsigned integer.
// Exits with an error message if parsing fails.
func MustParseU16(s string) uint16 {
	s = trim(s)
	if s == "" {
		fmt.Fprintln(os.Stderr, "Empty input.")
		os.Exit(1)
	}
	v := ParseU16(s)
	if v == nil {
		fmt.Fprintf(os.Stderr, "Invalid input: '%s' is not a valid 16-bit unsigned integer.\n", s)
		os.Exit(1)
	}
	return *v
}

// ParseU32 parses a string to a 32-bit unsigned integer.
// Returns nil if parsing fails.
func ParseU32(s string) *uint32 {
	s = trim(s)
	if s == "" {
		return nil
	}
	v, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		return nil
	}
	val := uint32(v)
	return &val
}

// MustParseU32 parses a string to a 32-bit unsigned integer.
// Exits with an error message if parsing fails.
func MustParseU32(s string) uint32 {
	s = trim(s)
	if s == "" {
		fmt.Fprintln(os.Stderr, "Empty input.")
		os.Exit(1)
	}
	v := ParseU32(s)
	if v == nil {
		fmt.Fprintf(os.Stderr, "Invalid input: '%s' is not a valid 32-bit unsigned integer.\n", s)
		os.Exit(1)
	}
	return *v
}

// ParseU64 parses a string to a 64-bit unsigned integer.
// Returns nil if parsing fails.
func ParseU64(s string) *uint64 {
	s = trim(s)
	if s == "" {
		return nil
	}
	v, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return nil
	}
	val := uint64(v)
	return &val
}

// MustParseU64 parses a string to a 64-bit unsigned integer.
// Exits with an error message if parsing fails.
func MustParseU64(s string) uint64 {
	s = trim(s)
	if s == "" {
		fmt.Fprintln(os.Stderr, "Empty input.")
		os.Exit(1)
	}
	v := ParseU64(s)
	if v == nil {
		fmt.Fprintf(os.Stderr, "Invalid input: '%s' is not a valid 64-bit unsigned integer.\n", s)
		os.Exit(1)
	}
	return *v
}

// ParseI8 parses a string to an 8-bit signed integer.
// Returns nil if parsing fails.
func ParseI8(s string) *int8 {
	s = trim(s)
	if s == "" {
		return nil
	}
	v, err := strconv.ParseInt(s, 10, 8)
	if err != nil {
		return nil
	}
	val := int8(v)
	return &val
}

// MustParseI8 parses a string to an 8-bit signed integer.
// Exits with an error message if parsing fails.
func MustParseI8(s string) int8 {
	s = trim(s)
	if s == "" {
		fmt.Fprintln(os.Stderr, "Empty input.")
		os.Exit(1)
	}
	v := ParseI8(s)
	if v == nil {
		fmt.Fprintf(os.Stderr, "Invalid input: '%s' is not a valid 8-bit signed integer.\n", s)
		os.Exit(1)
	}
	return *v
}

// ParseI16 parses a string to a 16-bit signed integer.
// Returns nil if parsing fails.
func ParseI16(s string) *int16 {
	s = trim(s)
	if s == "" {
		return nil
	}
	v, err := strconv.ParseInt(s, 10, 16)
	if err != nil {
		return nil
	}
	val := int16(v)
	return &val
}

// MustParseI16 parses a string to a 16-bit signed integer.
// Exits with an error message if parsing fails.
func MustParseI16(s string) int16 {
	s = trim(s)
	if s == "" {
		fmt.Fprintln(os.Stderr, "Empty input.")
		os.Exit(1)
	}
	v := ParseI16(s)
	if v == nil {
		fmt.Fprintf(os.Stderr, "Invalid input: '%s' is not a valid 16-bit signed integer.\n", s)
		os.Exit(1)
	}
	return *v
}

// ParseI32 parses a string to a 32-bit signed integer.
// Returns nil if parsing fails.
func ParseI32(s string) *int32 {
	s = trim(s)
	if s == "" {
		return nil
	}
	v, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		return nil
	}
	val := int32(v)
	return &val
}

// MustParseI32 parses a string to a 32-bit signed integer.
// Exits with an error message if parsing fails.
func MustParseI32(s string) int32 {
	s = trim(s)
	if s == "" {
		fmt.Fprintln(os.Stderr, "Empty input.")
		os.Exit(1)
	}
	v := ParseI32(s)
	if v == nil {
		fmt.Fprintf(os.Stderr, "Invalid input: '%s' is not a valid 32-bit signed integer.\n", s)
		os.Exit(1)
	}
	return *v
}

// ParseI64 parses a string to a 64-bit signed integer.
// Returns nil if parsing fails.
func ParseI64(s string) *int64 {
	s = trim(s)
	if s == "" {
		return nil
	}
	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return nil
	}
	val := int64(v)
	return &val
}

// MustParseI64 parses a string to a 64-bit signed integer.
// Exits with an error message if parsing fails.
func MustParseI64(s string) int64 {
	s = trim(s)
	if s == "" {
		fmt.Fprintln(os.Stderr, "Empty input.")
		os.Exit(1)
	}
	v := ParseI64(s)
	if v == nil {
		fmt.Fprintf(os.Stderr, "Invalid input: '%s' is not a valid 64-bit signed integer.\n", s)
		os.Exit(1)
	}
	return *v
}

// ParseF32 parses a string to a 32-bit floating point.
// Returns nil if parsing fails.
func ParseF32(s string) *float32 {
	s = trim(s)
	if s == "" {
		return nil
	}
	v, err := strconv.ParseFloat(s, 32)
	if err != nil {
		return nil
	}
	val := float32(v)
	return &val
}

// MustParseF32 parses a string to a 32-bit floating point.
// Exits with an error message if parsing fails.
func MustParseF32(s string) float32 {
	s = trim(s)
	if s == "" {
		fmt.Fprintln(os.Stderr, "Empty input.")
		os.Exit(1)
	}
	v := ParseF32(s)
	if v == nil {
		fmt.Fprintf(os.Stderr, "Invalid input: '%s' is not a valid 32-bit floating point.\n", s)
		os.Exit(1)
	}
	return *v
}

// ParseF64 parses a string to a 64-bit floating point.
// Returns nil if parsing fails.
func ParseF64(s string) *float64 {
	s = trim(s)
	if s == "" {
		return nil
	}
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return nil
	}
	return &v
}

// MustParseF64 parses a string to a 64-bit floating point.
// Exits with an error message if parsing fails.
func MustParseF64(s string) float64 {
	s = trim(s)
	if s == "" {
		fmt.Fprintln(os.Stderr, "Empty input.")
		os.Exit(1)
	}
	v := ParseF64(s)
	if v == nil {
		fmt.Fprintf(os.Stderr, "Invalid input: '%s' is not a valid 64-bit floating point.\n", s)
		os.Exit(1)
	}
	return *v
}
