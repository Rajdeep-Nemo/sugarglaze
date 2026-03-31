package evaluator

import (
	"fmt"
	"pluesi/internal/object"
)

// builtins is the central registry of all native Pluesi functions written in Go.
// Because this is in the "evaluator" package, evaluator.go can read it automatically!
var builtins = map[string]*object.Builtin{
	
	// print() - Prints any number of arguments to the console
	"print": {
		Fn: func(args ...object.Object) object.Object {
			for _, arg := range args {
				fmt.Print(arg.Inspect(), " ")
			}
			fmt.Println()
			return &object.Nil{} // Pluesi functions must always return an Object!
		},
	},

	// len() - Returns the length of a string
	"len": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return &object.Error{
					Message: fmt.Sprintf("wrong number of arguments for len(). got=%d, want=1", len(args)),
				}
			}

			switch arg := args[0].(type) {
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}
			default:
				return &object.Error{
					Message: fmt.Sprintf("argument to 'len' not supported, got %s", args[0].Type()),
				}
			}
		},
	},
}