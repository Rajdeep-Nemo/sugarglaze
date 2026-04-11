package evaluator

import (
	"fmt"

	"github.com/Rajdeep-Nemo/sugarglaze/internal/object"
)

// builtins is the central registry of all native functions.
// The string keys to the separate logic functions defined in other files.
var builtins = map[string]*object.Builtin{
	// Output functions
	"print":   {Fn: builtinPrint},
	"println": {Fn: builtinPrintln},
	// Input functions
	"trim":  {Fn: builtinTrim},
	"input": {Fn: builtinInput},
	// Extra functions
	"len": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return &object.Error{
					Message: fmt.Sprintf("wrong number of arguments for len(). got=%d, want=1", len(args)),
				}
			}

			switch arg := args[0].(type) {
			case *object.String:
				return &object.Int64{Value: int64(len(arg.Value))}
			default:
				return &object.Error{
					Message: fmt.Sprintf("argument to 'len' not supported, got %s", args[0].Type()),
				}
			}
		},
	},
}
