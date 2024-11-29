package main

import "fmt"

// Interpreter evaluates expressions.
type Interpreter struct {
	environment *Environment
	globals     *Environment
	locals      map[Expr]int
}

// NewInterpreter creates a new instance of the Interpreter.
func NewInterpreter() *Interpreter {
	globals := NewEnvironment()
	return &Interpreter{
		environment: globals,
		globals:     globals,
		locals:      make(map[Expr]int),
	}
}

// Interpret evaluates an expression and prints the result.
func (i *Interpreter) Interpret(expr Expr) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Runtime error:", r)
		}
	}()

	value := i.evaluate(expr)
	fmt.Println(stringify(value))
}

func (i *Interpreter) evaluate(expr Expr) interface{} {
	return expr.Accept(i)
}

// VisitLiteralExpr evaluates a literal expression.
func (i *Interpreter) VisitLiteralExpr(expr *Literal) interface{} {
	return expr.Value
}

// VisitGroupingExpr evaluates a grouping expression.
func (i *Interpreter) VisitGroupingExpr(expr *Grouping) interface{} {
	return i.evaluate(expr.Expression)
}

// VisitUnaryExpr evaluates a unary expression.
func (i *Interpreter) VisitUnaryExpr(expr *Unary) interface{} {
	right := i.evaluate(expr.Right)

	switch expr.Operator.TokenType {
	case TokenMinus:
		return -toFloat64(right)
	case TokenBang:
		return !isTruthy(right)
	}

	return nil
}

// VisitBinaryExpr evaluates a binary expression.
func (i *Interpreter) VisitBinaryExpr(expr *Binary) interface{} {
	left := i.evaluate(expr.Left)
	right := i.evaluate(expr.Right)

	switch expr.Operator.TokenType {
	case TokenPlus:
		if isString(left) && isString(right) {
			return left.(string) + right.(string)
		} else if isNumber(left) && isNumber(right) {
			return toFloat64(left) + toFloat64(right)
		}
		panic("Operands must be two numbers or two strings.")

	case TokenMinus:
		checkNumberOperands(expr.Operator, left, right)
		return toFloat64(left) - toFloat64(right)

	case TokenStar:
		checkNumberOperands(expr.Operator, left, right)
		return toFloat64(left) * toFloat64(right)

	case TokenSlash:
		checkNumberOperands(expr.Operator, left, right)
		if toFloat64(right) == 0 {
			panic("Division by zero.")
		}
		return toFloat64(left) / toFloat64(right)

	case TokenGreater:
		checkNumberOperands(expr.Operator, left, right)
		return toFloat64(left) > toFloat64(right)

	case TokenGreaterEqual:
		checkNumberOperands(expr.Operator, left, right)
		return toFloat64(left) >= toFloat64(right)

	case TokenLess:
		checkNumberOperands(expr.Operator, left, right)
		return toFloat64(left) < toFloat64(right)

	case TokenLessEqual:
		checkNumberOperands(expr.Operator, left, right)
		return toFloat64(left) <= toFloat64(right)

	case TokenEqualEqual:
		return isEqual(left, right)

	case TokenBangEqual:
		return !isEqual(left, right)
	}

	return nil
}

// Helper functions

func isTruthy(value interface{}) bool {
	if value == nil {
		return false
	}
	if b, ok := value.(bool); ok {
		return b
	}
	return true
}

func isEqual(a, b interface{}) bool {
	return a == b
}

func toFloat64(value interface{}) float64 {
	if num, ok := value.(float64); ok {
		return num
	}
	panic("Operand must be a number.")
}

func checkNumberOperands(operator Token, left, right interface{}) {
	if isNumber(left) && isNumber(right) {
		return
	}
	panic(fmt.Sprintf("Operands for %s must be numbers.", operator.Lexeme))
}

func isNumber(value interface{}) bool {
	_, ok := value.(float64)
	return ok
}

func isString(value interface{}) bool {
	_, ok := value.(string)
	return ok
}

func stringify(value interface{}) string {
	if value == nil {
		return "nil"
	}
	return fmt.Sprintf("%v", value)
}

// InterpretStatements evaluates statements.
func (i *Interpreter) InterpretStatements(statements []Stmt) {
	defer func() {
		if r := recover(); r != nil {
			if returnValue, ok := r.(ReturnValue); ok {
				fmt.Println(stringify(returnValue.Value))
				return
			}
			panic(r)
		}
	}()

	for _, stmt := range statements {
		i.execute(stmt)
	}
}

func (i *Interpreter) execute(stmt Stmt) {
	stmt.Accept(i)
}

// Statement visitors
func (i *Interpreter) VisitExpressionStmt(stmt *ExpressionStmt) interface{} {
	i.evaluate(stmt.Expression)
	return nil
}

func (i *Interpreter) VisitPrintStmt(stmt *PrintStmt) interface{} {
	value := i.evaluate(stmt.Expression)
	fmt.Println(stringify(value))
	return nil
}

func (i *Interpreter) VisitVarStmt(stmt *VarStmt) interface{} {
	var value interface{}
	if stmt.Initializer != nil {
		value = i.evaluate(stmt.Initializer)
	}
	i.environment.Define(stmt.Name.Lexeme, value)
	return nil
}

// VisitVariableExpr retrieves a variable's value from the environment.
func (i *Interpreter) VisitVariableExpr(expr *Variable) interface{} {
	return i.lookupVariable(expr.Name, expr)
}

// VisitAssignExpr assigns a value to a variable in the environment.
func (i *Interpreter) VisitAssignExpr(expr *Assign) interface{} {
	value := i.evaluate(expr.Value)
	if distance, found := i.locals[expr]; found {
		i.environment.assignAt(distance, expr.Name, value)
	} else {
		i.environment.Assign(expr.Name, value)
	}
	return value
}

// Execute a block with its own environment.
func (i *Interpreter) executeBlock(statements []Stmt, environment *Environment) {
	previous := i.environment
	defer func() {
		i.environment = previous
		if r := recover(); r != nil {
			if returnValue, ok := r.(ReturnValue); ok {
				panic(returnValue)
			}
			panic(r)
		}
	}()

	i.environment = environment
	for _, stmt := range statements {
		i.execute(stmt)
	}
}

// VisitBlockStmt executes a block with a new environment.
func (i *Interpreter) VisitBlockStmt(stmt *BlockStmt) interface{} {
	i.executeBlock(stmt.Statements, NewEnclosedEnvironment(i.environment))
	return nil
}

func (i *Interpreter) VisitIfStmt(stmt *IfStmt) interface{} {
	if isTruthy(i.evaluate(stmt.Condition)) {
		i.execute(stmt.ThenBranch)
	} else if stmt.ElseBranch != nil {
		i.execute(stmt.ElseBranch)
	}
	return nil
}

func (i *Interpreter) VisitWhileStmt(stmt *WhileStmt) interface{} {
	for isTruthy(i.evaluate(stmt.Condition)) {
		defer func() {
			if r := recover(); r != nil {
				if _, isBreak := r.(BreakException); isBreak {
					return // Exit loop cleanly
				}
				panic(r) // Re-panic for other errors
			}
		}()
		i.execute(stmt.Body)
	}
	return nil
}

type BreakException struct{}

func (e BreakException) Error() string {
	return "Break statement executed"
}

func (i *Interpreter) VisitBreakStmt(stmt *BreakStmt) interface{} {
	panic(BreakException{})
}

// LoxFunction represents a user-defined function.
type LoxFunction struct {
	declaration   *FunStmt
	closure       *Environment
	isInitializer bool
}

func (f *LoxFunction) Call(interpreter *Interpreter, arguments []interface{}) interface{} {
	// Create a new environment enclosing the closure
	environment := NewEnclosedEnvironment(f.closure)

	// Bind arguments to parameter names in the new environment
	for i, param := range f.declaration.Params {
		environment.Define(param.Lexeme, arguments[i])
	}

	// Execute the function body
	defer func() {
		if r := recover(); r != nil {
			if returnValue, ok := r.(ReturnValue); ok {
				panic(returnValue) // Propagate return value
			}
			panic(r) // Re-panic for other errors
		}
	}()

	interpreter.executeBlock(f.declaration.Body, environment)

	// If no return statement was executed, return nil
	return nil
}

func (f *LoxFunction) Arity() int {
	return len(f.declaration.Params)
}

func (i *Interpreter) VisitFunStmt(stmt *FunStmt) interface{} {
	// Create a function that captures the current environment as its closure
	function := &LoxFunction{
		declaration: stmt,
		closure:     i.environment,
	}
	i.environment.Define(stmt.Name.Lexeme, function)
	return nil
}

// VisitCallExpr handles function calls
func (i *Interpreter) VisitCallExpr(expr *Call) interface{} {
	callee := i.evaluate(expr.Callee)

	var arguments []interface{}
	for _, arg := range expr.Arguments {
		arguments = append(arguments, i.evaluate(arg))
	}

	function, ok := callee.(Callable)
	if !ok {
		// Special handling for returned functions
		if returnedFunc, isFunc := callee.(*LoxFunction); isFunc {
			function = returnedFunc
		} else {
			panic(fmt.Sprintf("Can only call functions and classes, got %T.", callee))
		}
	}

	if len(arguments) != function.Arity() {
		panic(fmt.Sprintf("Expected %d arguments but got %d.", function.Arity(), len(arguments)))
	}

	// Capture the return value
	var returnValue interface{}
	func() {
		defer func() {
			if r := recover(); r != nil {
				if retVal, ok := r.(ReturnValue); ok {
					returnValue = retVal.Value
					return
				}
				panic(r)
			}
		}()
		returnValue = function.Call(i, arguments)
	}()

	return returnValue
}

// ReturnValue represents a function's return value
type ReturnValue struct {
	Value interface{}
}

func (r ReturnValue) Error() string {
	return fmt.Sprintf("Return value: %v", r.Value)
}

// VisitReturnStmt handles return statements
func (i *Interpreter) VisitReturnStmt(stmt *ReturnStmt) interface{} {
	var value interface{} = nil
	if stmt.Value != nil {
		value = i.evaluate(stmt.Value)
	}
	panic(ReturnValue{Value: value})
}

type Callable interface {
	Arity() int
	Call(interpreter *Interpreter, arguments []interface{}) interface{}
}

// resolve stores the resolution depth of a variable.
func (i *Interpreter) resolve(expr Expr, depth int) {
	i.locals[expr] = depth
}

// Look up a variable's value using its resolution depth.
func (i *Interpreter) lookupVariable(name Token, expr Expr) interface{} {
	if distance, found := i.locals[expr]; found {
		return i.environment.getAt(distance, name.Lexeme)
	}
	// If not found in locals, check the global environment.
	return i.environment.Get(name)
}

func (i *Interpreter) VisitClassStmt(stmt *ClassStmt) interface{} {
	var superclass *LoxClass
	if stmt.Superclass != nil {
		superValue := i.evaluate(stmt.Superclass)
		var ok bool
		superclass, ok = superValue.(*LoxClass)
		if !ok {
			panic("Superclass must be a class.")
		}
	}

	i.environment.Define(stmt.Name.Lexeme, nil)

	if stmt.Superclass != nil {
		i.environment = NewEnclosedEnvironment(i.environment)
		i.environment.Define("super", superclass)
	}

	methods := make(map[string]*LoxFunction)
	for _, method := range stmt.Methods {
		function := &LoxFunction{
			declaration:   method,
			closure:       i.environment,
			isInitializer: method.Name.Lexeme == "init",
		}
		methods[method.Name.Lexeme] = function
	}

	class := &LoxClass{
		name:       stmt.Name.Lexeme,
		superclass: superclass,
		methods:    methods,
	}

	if superclass != nil {
		i.environment = i.environment.parent
	}

	i.environment.Assign(stmt.Name, class)
	return nil
}

func (c *LoxClass) Call(interpreter *Interpreter, arguments []interface{}) interface{} {
	instance := &LoxInstance{class: c, fields: make(map[string]interface{})}
	if initializer := c.findMethod("init"); initializer != nil {
		initializer.Bind(instance).Call(interpreter, arguments)
	}
	return instance
}

func (c *LoxClass) Arity() int {
	if initializer := c.findMethod("init"); initializer != nil {
		return initializer.Arity()
	}
	return 0
}

func (c *LoxClass) findMethod(name string) *LoxFunction {
	if method, ok := c.methods[name]; ok {
		return method
	}
	if c.superclass != nil {
		return c.superclass.findMethod(name)
	}
	return nil
}

func (i *Interpreter) VisitGetExpr(expr *GetExpr) interface{} {
    object := i.evaluate(expr.Object)
    if instance, ok := object.(*LoxInstance); ok {
        return instance.Get(expr.Name)
    }
    panic(RuntimeError{expr.Name, "Only instances have properties."})
}

func (i *Interpreter) VisitSetExpr(expr *SetExpr) interface{} {
    object := i.evaluate(expr.Object)
    if instance, ok := object.(*LoxInstance); ok {
        value := i.evaluate(expr.Value)
        instance.Set(expr.Name, value)
        return value
    }
    panic(RuntimeError{expr.Name, "Only instances have fields."})
}

func (i *Interpreter) VisitThisExpr(expr *ThisExpr) interface{} {
	return i.lookupVariable(expr.Keyword, expr)
}

func (i *Interpreter) VisitSuperExpr(expr *SuperExpr) interface{} {
    distance := i.locals[expr]
    superclass := i.environment.getAt(distance, "super").(*LoxClass)
    object := i.environment.getAt(distance-1, "this").(*LoxInstance)
    method := superclass.FindMethod(expr.Method.Lexeme)
    if method == nil {
        panic(RuntimeError{expr.Method, "Undefined property '" + expr.Method.Lexeme + "'."})
    }
    return method.Bind(object)
}

func (c *LoxClass) FindMethod(name string) *LoxFunction {
	if method, found := c.methods[name]; found {
		return method
	}

	if c.superclass != nil {
		return c.superclass.FindMethod(name)
	}

	return nil
}

func (f *LoxFunction) Bind(instance *LoxInstance) *LoxFunction {
    environment := NewEnclosedEnvironment(f.closure)
    environment.Define("this", instance)
    return &LoxFunction{declaration: f.declaration, closure: environment, isInitializer: f.isInitializer}
}

// LoxClass represents a runtime class
type LoxClass struct {
	name       string
	superclass *LoxClass
	methods    map[string]*LoxFunction
}

// LoxInstance represents an instance of a class
type LoxInstance struct {
	class  *LoxClass
	fields map[string]interface{}
}

func (i *LoxInstance) Get(name Token) interface{} {
    if value, ok := i.fields[name.Lexeme]; ok {
        return value
    }

    method := i.class.FindMethod(name.Lexeme)
    if method != nil {
        return method.Bind(i)
    }

    panic(RuntimeError{name, "Undefined property '" + name.Lexeme + "'."})
}

func (i *LoxInstance) Set(name Token, value interface{}) {
    i.fields[name.Lexeme] = value
}

type RuntimeError struct {
    token Token
    message string
}

func (e RuntimeError) Error() string {
    return fmt.Sprintf("%s\n[line %d]", e.message, e.token.Line)
}