package evaluator

import (
	"fmt"
	"pluesi/internal/ast"
	"pluesi/internal/object"
)

// Optimization: We only ever need two boolean objects in memory!
var (
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

// Eval is the core recursive function that walks the AST
func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {

	// STATEMENTS
	case *ast.Program:
		return evalProgram(node, env)

	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)

	// LITERALS
	case *ast.IntegerLiteral:
		return &object.Integer{Value: int64(node.Value)}

	case *ast.FloatLiteral:
		return &object.Float{Value: node.Value}

	case *ast.StringLiteral:
		return &object.String{Value: node.Value}

	case *ast.CharLiteral:
		// We store chars as single-character strings for simplicity
		return &object.String{Value: string(node.Value)}

	case *ast.BoolLiteral:
		return nativeBoolToBooleanObject(node.Value)

	case *ast.NilLiteral:
		return &object.Nil{}

	// EXPRESSIONS
	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right)

	case *ast.InfixExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalInfixExpression(node.Operator, left, right)

	// VARIABLES & IDENTIFIERS
	case *ast.LetStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		// Save it to the environment as mutable (isConst = false)
		env.Define(node.Name.Value, val, false)
		return &object.Nil{}

	case *ast.ConstStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		// Save it to the environment as IMMUTABLE (isConst = true)
		env.Define(node.Name.Value, val, true)
		return &object.Nil{}

	case *ast.AssignStatement:
		// 1. Evaluate the right side of the equals sign (e.g., x + y)
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}

		// 2. If it is a standard "=" assignment
		if node.Operator == ast.Assign {
			result := env.Assign(node.Target.Value, val)
			// env.Assign returns an Error object if the variable doesn't exist or is a const!
			if isError(result) {
				return result
			}
			return &object.Nil{}
		}

		// 3. Handle compound assignments like +=, -=, *=, /=
		currentRecord, exists := env.Get(node.Target.Value)
		if !exists {
			return newError("cannot assign to undefined variable '%s'", node.Target.Value)
		}

		// Map the AssignOperator to a standard math operator
		var opStr string
		switch node.Operator {
		case ast.PlusAssign:
			opStr = "+"
		case ast.MinusAssign:
			opStr = "-"
		case ast.StarAssign:
			opStr = "*"
		case ast.SlashAssign:
			opStr = "/"
		case ast.PercentAssign:
			opStr = "%"
		}

		// Reuse your existing math evaluator!
		newVal := evalInfixExpression(opStr, currentRecord.Value, val)
		if isError(newVal) {
			return newVal
		}

		// Save the newly calculated value back to memory
		result := env.Assign(node.Target.Value, newVal)
		if isError(result) {
			return result
		}

		return &object.Nil{}
	case *ast.Identifier:
		return evalIdentifier(node, env)

	// CONTROL FLOW & FUNCTIONS
	case *ast.BlockStatement:
		return evalBlockStatement(node, env)

	case *ast.IfExpression:
		return evalIfExpression(node, env)

	case *ast.CallExpression:
		// 1. Evaluate the function name (e.g., gets the Builtin object for "print")
		function := Eval(node.Function, env)
		if isError(function) {
			return function
		}

		// 2. Evaluate all the arguments passed into it
		args := evalExpressions(node.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}

		// 3. Execute!
		return applyFunction(function, args)
	}

	return nil
}

// Evaluates a list of statements in the main program
func evalProgram(program *ast.Program, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range program.Statements {
		result = Eval(statement, env)

		// ERROR HANDLING: If any statement returns an Error, stop the whole program and bubble it up!
		if err, ok := result.(*object.Error); ok {
			return err
		}
	}

	return result
}

// Prevents us from creating duplicate boolean objects in memory
func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return TRUE
	}
	return FALSE
}

// Helper function to easily throw formatted errors
func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

// Helper function to check if an object is an error
func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}
	return false
}

// Evaluates prefix expressions like !true or -5
func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	default:
		return newError("unknown operator: %s%s", operator, right.Type())
	}
}

// Evaluate the bang operator (!) for booleans and nil
func evalBangOperatorExpression(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case &object.Nil{}:
		return TRUE
	default:
		return FALSE
	}
}

// Evaluates negative numbers
func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJ && right.Type() != object.FLOAT_OBJ {
		return newError("unknown operator: -%s", right.Type())
	}

	if right.Type() == object.INTEGER_OBJ {
		value := right.(*object.Integer).Value
		return &object.Integer{Value: -value}
	}

	value := right.(*object.Float).Value
	return &object.Float{Value: -value}
}

// Evaluates infix expressions like 5 + 10 or "Hello " + "World"
func evalInfixExpression(operator string, left, right object.Object) object.Object {
	switch {
	// Integer Math
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left, right)

	// String Concatenation (e.g., "Hello " + "World")
	case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
		return evalStringInfixExpression(operator, left, right)

	// Boolean Equality (e.g., true == false)
	case operator == "==":
		return nativeBoolToBooleanObject(left == right) // Pointer comparison works because of our TRUE/FALSE singletons
	case operator == "!=":
		return nativeBoolToBooleanObject(left != right)

	// Type Mismatches
	case left.Type() != right.Type():
		return newError("type mismatch: %s %s %s", left.Type(), operator, right.Type())

	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

// Evaluates infix expressions for integers (e.g., 5 + 10 or 20 > 15)
func evalIntegerInfixExpression(operator string, left, right object.Object) object.Object {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value

	switch operator {
	case "+":
		return &object.Integer{Value: leftVal + rightVal}
	case "-":
		return &object.Integer{Value: leftVal - rightVal}
	case "*":
		return &object.Integer{Value: leftVal * rightVal}
	case "/":
		if rightVal == 0 {
			return newError("division by zero")
		}
		return &object.Integer{Value: leftVal / rightVal}
	case "%":
		if rightVal == 0 {
			return newError("modulo by zero")
		}
		return &object.Integer{Value: leftVal % rightVal}
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

// Handles string concatenation (e.g., "Hello " + "World")
func evalStringInfixExpression(operator string, left, right object.Object) object.Object {
	if operator != "+" {
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
	leftVal := left.(*object.String).Value
	rightVal := right.(*object.String).Value
	return &object.String{Value: leftVal + rightVal}
}

// Evaluates identifiers by looking them up in the environment or built-in functions
func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	// 1. Check if the identifier is a variable in the current Environment scope
	if record, ok := env.Get(node.Value); ok {
		return record.Value
	}

	// 2. Check if the identifier is actually a Built-in Function (like "print" or "input")
	if builtin, ok := builtins[node.Value]; ok {
		return builtin
	}

	// 3. If it's in neither, the user is trying to use an undeclared variable!
	return newError("identifier not found: %s", node.Value)
}

// --- Function Execution ---

func evalExpressions(exps []ast.Expression, env *object.Environment) []object.Object {
	var result []object.Object

	for _, e := range exps {
		evaluated := Eval(e, env)
		if isError(evaluated) {
			return []object.Object{evaluated} // Stop evaluating args if one fails
		}
		result = append(result, evaluated)
	}

	return result
}

func applyFunction(fn object.Object, args []object.Object) object.Object {
	switch fn := fn.(type) {
	case *object.Builtin:
		return fn.Fn(args...)
	default:
		return newError("not a function: %s", fn.Type())
	}
}

// --- Control Flow ---

func evalBlockStatement(block *ast.BlockStatement, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range block.Statements {
		result = Eval(statement, env)

		if result != nil && result.Type() == object.ERROR_OBJ {
			return result // Bubble up errors immediately
		}
	}

	return result
}

func evalIfExpression(ie *ast.IfExpression, env *object.Environment) object.Object {
	condition := Eval(ie.Condition, env)
	if isError(condition) {
		return condition
	}

	if isTruthy(condition) {
		// Execute the Consequence block in a NEW enclosed environment!
		enclosedEnv := object.NewEnclosedEnvironment(env)
		return Eval(ie.Consequence, enclosedEnv)
	} else if ie.Alternative != nil {
		// Execute the Else block in a NEW enclosed environment!
		enclosedEnv := object.NewEnclosedEnvironment(env)
		return Eval(ie.Alternative, enclosedEnv)
	} else {
		return &object.Nil{}
	}
}

func isTruthy(obj object.Object) bool {
	switch obj {
	case &object.Nil{}:
		return false
	case TRUE:
		return true
	case FALSE:
		return false
	default:
		return true
	}
}
