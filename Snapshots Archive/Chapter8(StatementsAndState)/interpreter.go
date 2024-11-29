package main

import "fmt"

// Interpreter evaluates expressions.
type Interpreter struct{
	environment *Environment
}

// NewInterpreter creates a new instance of the Interpreter.
func NewInterpreter() *Interpreter {
	return &Interpreter{environment: NewEnvironment()}
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
	// Remove the error suppression defer
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
	// This will now re-panic any error from Get, rather than catching it
	return i.environment.Get(expr.Name)
}

// VisitAssignExpr assigns a value to a variable in the environment.
func (i *Interpreter) VisitAssignExpr(expr *Assign) interface{} {
	value := i.evaluate(expr.Value)
	i.environment.Assign(expr.Name, value)
	return value
}

// Execute a block with its own environment.
func (i *Interpreter) executeBlock(statements []Stmt, environment *Environment) {
	previous := i.environment
	defer func() { i.environment = previous }() // Restore the previous environment.

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

