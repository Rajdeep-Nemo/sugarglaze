// internal/object/object.go
package object

import "fmt"

type ObjectType string

const (
	I8_OBJ      = "i8"
	I16_OBJ     = "i16"
	I32_OBJ     = "i32"
	I64_OBJ     = "i64"
	U8_OBJ      = "u8"
	U16_OBJ     = "u16"
	U32_OBJ     = "u32"
	U64_OBJ     = "u64"
	F32_OBJ     = "f32"
	F64_OBJ     = "f64"
	STRING_OBJ  = "string"
	BOOL_OBJ    = "bool"
	NIL_OBJ     = "NIL"
	ERROR_OBJ   = "error"
	BUILTIN_OBJ = "builtin"
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Int8 struct{ Value int8 }

func (i *Int8) Type() ObjectType { return I8_OBJ }
func (i *Int8) Inspect() string  { return fmt.Sprintf("%d", i.Value) }

type Int16 struct{ Value int16 }

func (i *Int16) Type() ObjectType { return I16_OBJ }
func (i *Int16) Inspect() string  { return fmt.Sprintf("%d", i.Value) }

type Int32 struct{ Value int32 }

func (i *Int32) Type() ObjectType { return I32_OBJ }
func (i *Int32) Inspect() string  { return fmt.Sprintf("%d", i.Value) }

type Int64 struct{ Value int64 }

func (i *Int64) Type() ObjectType { return I64_OBJ }
func (i *Int64) Inspect() string  { return fmt.Sprintf("%d", i.Value) }

type Uint8 struct{ Value uint8 }

func (i *Uint8) Type() ObjectType { return U8_OBJ }
func (i *Uint8) Inspect() string  { return fmt.Sprintf("%d", i.Value) }

type Uint16 struct{ Value uint16 }

func (i *Uint16) Type() ObjectType { return U16_OBJ }
func (i *Uint16) Inspect() string  { return fmt.Sprintf("%d", i.Value) }

type Uint32 struct{ Value uint32 }

func (i *Uint32) Type() ObjectType { return U32_OBJ }
func (i *Uint32) Inspect() string  { return fmt.Sprintf("%d", i.Value) }

type Uint64 struct{ Value uint64 }

func (i *Uint64) Type() ObjectType { return U64_OBJ }
func (i *Uint64) Inspect() string  { return fmt.Sprintf("%d", i.Value) }

type Float32 struct{ Value float32 }

func (f *Float32) Type() ObjectType { return F32_OBJ }
func (f *Float32) Inspect() string  { return fmt.Sprintf("%f", f.Value) }

type Float64 struct{ Value float64 }

func (f *Float64) Type() ObjectType { return F64_OBJ }
func (f *Float64) Inspect() string  { return fmt.Sprintf("%f", f.Value) }

type Boolean struct{ Value bool }

func (b *Boolean) Type() ObjectType { return BOOL_OBJ }
func (b *Boolean) Inspect() string  { return fmt.Sprintf("%t", b.Value) }

type String struct{ Value string }

func (s *String) Type() ObjectType { return STRING_OBJ }
func (s *String) Inspect() string  { return s.Value }

type Error struct{ Message string }

func (e *Error) Type() ObjectType { return ERROR_OBJ }
func (e *Error) Inspect() string  { return "ERROR: " + e.Message }

type Nil struct{}

func (n *Nil) Type() ObjectType { return NIL_OBJ }
func (n *Nil) Inspect() string  { return "NIL" }

type BuiltinFunction func(args ...Object) Object

type Builtin struct {
	Fn BuiltinFunction
}

func (b *Builtin) Type() ObjectType { return BUILTIN_OBJ }
func (b *Builtin) Inspect() string  { return "built-in function" }
