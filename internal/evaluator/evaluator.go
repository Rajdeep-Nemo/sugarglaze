package evaluator

import (
	"fmt"
	"glaze/internal/ast"
	"glaze/internal/object"
	"math"
)

// Recursively traverses and evaluates an Abstract Syntax Tree (AST) node into a runtime object
func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {

	// STATEMENTS
	case *ast.Program:
		return evalProgram(node, env)
	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)
	case *ast.IntegerLiteral:
		if node.Value > math.MaxInt32 {
			return &object.Int64{Value: int64(node.Value)}
		}
		return &object.Int32{Value: int32(node.Value)}
	case *ast.FloatLiteral:
		return &object.Float32{Value: float32(node.Value)}
	case *ast.StringLiteral:
		return &object.String{Value: node.Value}
	case *ast.CharLiteral:
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
		var val object.Object
		var expectedType string
		if node.TypeHint != nil {
			expectedType = node.TypeHint.Name
		}

		isInit := false
		if node.Value != nil {
			val = Eval(node.Value, env)
			if isError(val) {
				return val
			}
			if expectedType != "" {
				val = enforceType(val, expectedType)
				if isError(val) {
					return val
				}
			} else {
				expectedType = string(val.Type())
			}
			isInit = true
		} else {
			val = &object.Nil{}
		}
		env.Define(node.Name.Value, val, false, expectedType, isInit)
		return &object.Nil{}

	case *ast.ConstStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		expectedType := ""
		if node.TypeHint != nil {
			expectedType = node.TypeHint.Name
			val = enforceType(val, expectedType)
			if isError(val) {
				return val
			}
		} else {
			expectedType = string(val.Type())
		}
		env.Define(node.Name.Value, val, true, expectedType, true)
		return &object.Nil{}

	case *ast.AssignStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}

		if node.Operator == ast.Assign {
			result := env.Assign(node.Target.Value, val)
			if isError(result) {
				return result
			}
			return &object.Nil{}
		}

		currentRecord, exists := env.Get(node.Target.Value)
		if !exists {
			return newError("cannot assign to undefined variable '%s'", node.Target.Value)
		}

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

		newVal := evalInfixExpression(opStr, currentRecord.Value, val)
		if isError(newVal) {
			return newVal
		}

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

// Ensures an object matches a target type by performing implicit casting or returning out-of-bounds errors
func enforceType(obj object.Object, targetType string) object.Object {
	if string(obj.Type()) == targetType {
		return obj
	}

	var floatVal float64
	var intVal int64
	var isFloat, isInt bool

	switch v := obj.(type) {
	case *object.Int8:
		intVal, isInt = int64(v.Value), true
	case *object.Int16:
		intVal, isInt = int64(v.Value), true
	case *object.Int32:
		intVal, isInt = int64(v.Value), true
	case *object.Int64:
		intVal, isInt = v.Value, true
	case *object.Uint8:
		intVal, isInt = int64(v.Value), true
	case *object.Uint16:
		intVal, isInt = int64(v.Value), true
	case *object.Uint32:
		intVal, isInt = int64(v.Value), true
	case *object.Uint64:
		if v.Value > math.MaxInt64 {
			return newError("value out of bounds for casting")
		}
		intVal, isInt = int64(v.Value), true
	case *object.Float32:
		floatVal, isFloat = float64(v.Value), true
	case *object.Float64:
		floatVal, isFloat = v.Value, true
	default:
		return newError("cannot cast %s to %s", obj.Type(), targetType)
	}

	if isInt {
		switch targetType {
		case object.I8_OBJ:
			if intVal < math.MinInt8 || intVal > math.MaxInt8 {
				return newError("value %d out of bounds for i8", intVal)
			}
			return &object.Int8{Value: int8(intVal)}
		case object.I16_OBJ:
			if intVal < math.MinInt16 || intVal > math.MaxInt16 {
				return newError("value %d out of bounds for i16", intVal)
			}
			return &object.Int16{Value: int16(intVal)}
		case object.I32_OBJ:
			if intVal < math.MinInt32 || intVal > math.MaxInt32 {
				return newError("value %d out of bounds for i32", intVal)
			}
			return &object.Int32{Value: int32(intVal)}
		case object.I64_OBJ:
			return &object.Int64{Value: intVal}
		case object.U8_OBJ:
			if intVal < 0 || intVal > math.MaxUint8 {
				return newError("value %d out of bounds for u8", intVal)
			}
			return &object.Uint8{Value: uint8(intVal)}
		case object.U16_OBJ:
			if intVal < 0 || intVal > math.MaxUint16 {
				return newError("value %d out of bounds for u16", intVal)
			}
			return &object.Uint16{Value: uint16(intVal)}
		case object.U32_OBJ:
			if intVal < 0 || intVal > math.MaxUint32 {
				return newError("value %d out of bounds for u32", intVal)
			}
			return &object.Uint32{Value: uint32(intVal)}
		case object.U64_OBJ:
			if intVal < 0 {
				return newError("value %d out of bounds for u64", intVal)
			}
			return &object.Uint64{Value: uint64(intVal)}
		case object.F32_OBJ:
			return &object.Float32{Value: float32(intVal)}
		case object.F64_OBJ:
			return &object.Float64{Value: float64(intVal)}
		}
	} else if isFloat {
		switch targetType {
		case object.F32_OBJ:
			if floatVal < -math.MaxFloat32 || floatVal > math.MaxFloat32 {
				return newError("value out of bounds for f32")
			}
			return &object.Float32{Value: float32(floatVal)}
		case object.F64_OBJ:
			return &object.Float64{Value: floatVal}
		default:
			return newError("implicit down-casting from float to int not allowed")
		}
	}
	return newError("type mismatch: cannot assign '%s' to '%s'", obj.Type(), targetType)
}

// Evaluates a sequence of statements forming a program, bubbling up errors to halt execution
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

// Two boolean objects to represent true and false values
var (
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

// Converts a native Go boolean into the interpreter's singleton TRUE or FALSE object
func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return TRUE
	}
	return FALSE
}

// Constructs and returns a formatted runtime error object
func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

// Checks whether a given evaluated object is an error
func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}
	return false
}

// Evaluates prefix expressions by delegating to specific functions for operators like ! or -
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

// Evaluates the logical NOT (!) operator by returning the opposite of an object's truthiness
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

// Evaluates the arithmetic negation (-) operator for integer and float objects
func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	switch obj := right.(type) {
	case *object.Int8:
		return &object.Int8{Value: -obj.Value}
	case *object.Int16:
		return &object.Int16{Value: -obj.Value}
	case *object.Int32:
		return &object.Int32{Value: -obj.Value}
	case *object.Int64:
		return &object.Int64{Value: -obj.Value}
	case *object.Float32:
		return &object.Float32{Value: -obj.Value}
	case *object.Float64:
		return &object.Float64{Value: -obj.Value}
	case *object.Uint8, *object.Uint16, *object.Uint32, *object.Uint64:
		return newError("cannot negate unsigned integer type %s", right.Type())
	default:
		return newError("unknown operator: -%s", right.Type())
	}
}

// Routes binary expressions to their specific evaluation logic based on operand types
func evalInfixExpression(operator string, left, right object.Object) object.Object {
	switch {
	case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
		return evalStringInfixExpression(operator, left, right)
	case operator == "==" && left.Type() == object.BOOL_OBJ && right.Type() == object.BOOL_OBJ:
		return nativeBoolToBooleanObject(left == right)
	case operator == "!=" && left.Type() == object.BOOL_OBJ && right.Type() == object.BOOL_OBJ:
		return nativeBoolToBooleanObject(left != right)
	case left.Type() != right.Type():
		return newError("type mismatch: %s %s %s", left.Type(), operator, right.Type())
	default:
		return evalNumericInfixExpression(operator, left, right)
	}
}

// Evaluates infix expressions for numeric types by performing the appropriate arithmetic or comparison operation
func evalNumericInfixExpression(operator string, left, right object.Object) object.Object {
	switch leftVal := left.(type) {
	case *object.Int8:
		return evalIntMath(operator, int64(leftVal.Value), int64(right.(*object.Int8).Value), object.I8_OBJ)
	case *object.Int16:
		return evalIntMath(operator, int64(leftVal.Value), int64(right.(*object.Int16).Value), object.I16_OBJ)
	case *object.Int32:
		return evalIntMath(operator, int64(leftVal.Value), int64(right.(*object.Int32).Value), object.I32_OBJ)
	case *object.Int64:
		return evalIntMath(operator, leftVal.Value, right.(*object.Int64).Value, object.I64_OBJ)
	case *object.Uint8:
		return evalUintMath(operator, uint64(leftVal.Value), uint64(right.(*object.Uint8).Value), object.U8_OBJ)
	case *object.Uint16:
		return evalUintMath(operator, uint64(leftVal.Value), uint64(right.(*object.Uint16).Value), object.U16_OBJ)
	case *object.Uint32:
		return evalUintMath(operator, uint64(leftVal.Value), uint64(right.(*object.Uint32).Value), object.U32_OBJ)
	case *object.Uint64:
		return evalUintMath(operator, leftVal.Value, right.(*object.Uint64).Value, object.U64_OBJ)
	case *object.Float32:
		return evalFloatMath(operator, float64(leftVal.Value), float64(right.(*object.Float32).Value), object.F32_OBJ)
	case *object.Float64:
		return evalFloatMath(operator, leftVal.Value, right.(*object.Float64).Value, object.F64_OBJ)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

// Evaluates infix expressions for integer types, performing arithmetic and comparison operations
func evalIntMath(operator string, leftVal, rightVal int64, targetType object.ObjectType) object.Object {
	var result int64
	switch operator {
	case "+":
		result = leftVal + rightVal
	case "-":
		result = leftVal - rightVal
	case "*":
		result = leftVal * rightVal
	case "/":
		if rightVal == 0 {
			return newError("division by zero")
		}
		result = leftVal / rightVal
	case "%":
		if rightVal == 0 {
			return newError("modulo by zero")
		}
		result = leftVal % rightVal
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	case "<=":
		return nativeBoolToBooleanObject(leftVal <= rightVal)
	case ">=":
		return nativeBoolToBooleanObject(leftVal >= rightVal)
	default:
		return newError("unknown operator: %s %s %s", targetType, operator, targetType)
	}

	switch targetType {
	case object.I8_OBJ:
		return &object.Int8{Value: int8(result)}
	case object.I16_OBJ:
		return &object.Int16{Value: int16(result)}
	case object.I32_OBJ:
		return &object.Int32{Value: int32(result)}
	case object.I64_OBJ:
		return &object.Int64{Value: result}
	}
	return newError("internal error during integer math")
}

// Evaluates infix expressions for unsigned integer types, performing arithmetic and comparison operations
func evalUintMath(operator string, leftVal, rightVal uint64, targetType object.ObjectType) object.Object {
	var result uint64
	switch operator {
	case "+":
		result = leftVal + rightVal
	case "-":
		result = leftVal - rightVal
	case "*":
		result = leftVal * rightVal
	case "/":
		if rightVal == 0 {
			return newError("division by zero")
		}
		result = leftVal / rightVal
	case "%":
		if rightVal == 0 {
			return newError("modulo by zero")
		}
		result = leftVal % rightVal
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	case "<=":
		return nativeBoolToBooleanObject(leftVal <= rightVal)
	case ">=":
		return nativeBoolToBooleanObject(leftVal >= rightVal)
	default:
		return newError("unknown operator: %s %s %s", targetType, operator, targetType)
	}

	switch targetType {
	case object.U8_OBJ:
		return &object.Uint8{Value: uint8(result)}
	case object.U16_OBJ:
		return &object.Uint16{Value: uint16(result)}
	case object.U32_OBJ:
		return &object.Uint32{Value: uint32(result)}
	case object.U64_OBJ:
		return &object.Uint64{Value: result}
	}
	return newError("internal error during unsigned integer math")
}

// Evaluates infix expressions for float types, performing arithmetic and comparison operations
func evalFloatMath(operator string, leftVal, rightVal float64, targetType object.ObjectType) object.Object {
	var result float64
	switch operator {
	case "+":
		result = leftVal + rightVal
	case "-":
		result = leftVal - rightVal
	case "*":
		result = leftVal * rightVal
	case "/":
		result = leftVal / rightVal
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	case "<=":
		return nativeBoolToBooleanObject(leftVal <= rightVal)
	case ">=":
		return nativeBoolToBooleanObject(leftVal >= rightVal)
	default:
		return newError("unknown operator: %s %s %s", targetType, operator, targetType)
	}

	switch targetType {
	case object.F32_OBJ:
		return &object.Float32{Value: float32(result)}
	case object.F64_OBJ:
		return &object.Float64{Value: result}
	}
	return newError("internal error during float math")
}

// Performs string concatenation when the + operator is used with string objects
func evalStringInfixExpression(operator string, left, right object.Object) object.Object {
	if operator != "+" {
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
	leftVal := left.(*object.String).Value
	rightVal := right.(*object.String).Value
	return &object.String{Value: leftVal + rightVal}
}

// Resolves an identifier by looking it up in the environment or built-in function registry
func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	if record, ok := env.Get(node.Value); ok {
		if !record.IsInitialized {
			return newError("cannot access uninitialized variable '%s'", node.Value)
		}
		return record.Value
	}
	if builtin, ok := builtins[node.Value]; ok {
		return builtin
	}
	return newError("identifier not found: %s", node.Value)
}

// --- Function Execution ---
// Evaluates a list of expressions sequentially, returning early if any evaluation fails
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

// Executes a resolved function object using the provided arguments
func applyFunction(fn object.Object, args []object.Object) object.Object {
	switch fn := fn.(type) {
	case *object.Builtin:
		return fn.Fn(args...)
	default:
		return newError("not a function: %s", fn.Type())
	}
}

// --- Control Flow ---
// Evaluates a list of statements within a block, halting early if an error occurs
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

// Evaluates conditional logic and executes the corresponding block in a newly enclosed scope
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

// Determines if an object is logically true or false for conditional evaluation
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
