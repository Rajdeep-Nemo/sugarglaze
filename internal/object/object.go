package object

import "fmt"

// ObjectType represents the type of a value
type ObjectType string

const (
	INTEGER_OBJ = "i32"
	FLOAT_OBJ   = "f64"
	STRING_OBJ  = "string"
	BOOL_OBJ    = "bool"
	NIL_OBJ     = "NIL"
	ERROR_OBJ   = "error"
	BUILTIN_OBJ = "builtin"
)

// Object is the interface every value in the language must implement
type Object interface {
	Type() ObjectType
	Inspect() string
}

// Integer Object
type Integer struct{ Value int64 }
func (i *Integer) Type() ObjectType { return INTEGER_OBJ }
func (i *Integer) Inspect() string  { return fmt.Sprintf("%d", i.Value) }

// Float Object
type Float struct{ Value float64 }
func (f *Float) Type() ObjectType { return FLOAT_OBJ }
func (f *Float) Inspect() string  { return fmt.Sprintf("%f", f.Value) }

// Boolean Object
type Boolean struct{ Value bool }
func (b *Boolean) Type() ObjectType { return BOOL_OBJ }
func (b *Boolean) Inspect() string  { return fmt.Sprintf("%t", b.Value) }

// String Object
type String struct{ Value string }
func (s *String) Type() ObjectType { return STRING_OBJ }
func (s *String) Inspect() string  { return s.Value }

// Error Object (for throwing type mismatches and runtime panics)
type Error struct{ Message string }
func (e *Error) Type() ObjectType { return ERROR_OBJ }
func (e *Error) Inspect() string  { return "ERROR: " + e.Message }

// Nil Object
type Nil struct{}
func (n *Nil) Type() ObjectType { return NIL_OBJ }
func (n *Nil) Inspect() string  { return "NIL" }

// BuiltinFunction represents the underlying Go function powering a Pluesi builtin
type BuiltinFunction func(args ...Object) Object

// Builtin Object wraps a Go function so it can be passed around in Pluesi
type Builtin struct {
	Fn BuiltinFunction
}
func (b *Builtin) Type() ObjectType { return BUILTIN_OBJ }
func (b *Builtin) Inspect() string  { return "built-in function" }