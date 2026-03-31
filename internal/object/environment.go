package object

import "fmt"

// Record holds the actual value and whether it is a constant
type Record struct {
	Value   Object
	IsConst bool
}

// Environment is a map of strings to Records, with a pointer to an outer scope
type Environment struct {
	store map[string]Record
	outer *Environment
}

// NewEnvironment creates a fresh, global environment
func NewEnvironment() *Environment {
	s := make(map[string]Record)
	return &Environment{store: s, outer: nil}
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
func (e *Environment) Define(name string, val Object, isConst bool) Object {
	e.store[name] = Record{Value: val, IsConst: isConst}
	return val
}

// Assign updates an EXISTING variable (used for `x = 10` or `x += 1`)
func (e *Environment) Assign(name string, val Object) Object {
	record, ok := e.store[name]

	// 1. If it's not in the local scope, try to assign it in the outer scope
	if !ok {
		if e.outer != nil {
			return e.outer.Assign(name, val)
		}
		// If we hit the top without finding it, the user is assigning an undeclared variable
		return &Error{Message: fmt.Sprintf("cannot assign to undefined variable '%s'", name)}
	}

	// 2. Enforce Const Immutability!
	if record.IsConst {
		return &Error{Message: fmt.Sprintf("cannot reassign to const variable '%s'", name)}
	}

	// 3. Update the value
	e.store[name] = Record{Value: val, IsConst: false}
	return val
}