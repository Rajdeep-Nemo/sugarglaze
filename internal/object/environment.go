package object

import "fmt"

// Record holds the actual value and whether it is a constant, as well as its type (for type checking)
type Record struct {
	Value         Object
	IsConst       bool
	Type          string
	IsInitialized bool
}

// Environment is a map of strings to Records, with a pointer to an outer scope
type Environment struct {
	store map[string]Record
	outer *Environment
}

// NewEnvironment creates a fresh, global environment
func NewEnvironment() *Environment {
	return &Environment{store: make(map[string]Record), outer: nil}
}

// NewEnclosedEnvironment creates a new local scope (e.g., for a function or if-block)
func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}

// Looks up a variable by name. It checks the local scope first, then walks up the outer scopes.
func (e *Environment) Get(name string) (Record, bool) {
	record, ok := e.store[name]
	if !ok && e.outer != nil {
		return e.outer.Get(name)
	}
	return record, ok
}

// Creates a brand new variable in the CURRENT scope (used for `let` and `const`)
func (e *Environment) Define(name string, val Object, isConst bool, varType string, isInit bool) Object {
	e.store[name] = Record{Value: val, IsConst: isConst, Type: varType, IsInitialized: isInit}
	return val
}

// Assign updates an EXISTING variable (used for `x = 10` or `x += 1`)
func (e *Environment) Assign(name string, val Object) Object {
	record, ok := e.store[name]

	if !ok {
		if e.outer != nil {
			return e.outer.Assign(name, val)
		}
		return &Error{Message: fmt.Sprintf("cannot assign to undefined variable '%s'", name)}
	}

	if record.IsConst {
		return &Error{Message: fmt.Sprintf("cannot reassign to const variable '%s'", name)}
	}

	if record.Type != "" && record.Type != string(val.Type()) {
		return &Error{Message: fmt.Sprintf("type mismatch: cannot assign '%s' to variable '%s' (expected '%s')",
			val.Type(), name, record.Type)}
	}

	e.store[name] = Record{Value: val, IsConst: false, Type: record.Type, IsInitialized: true}
	return val
}
