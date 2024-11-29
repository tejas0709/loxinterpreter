package main

import "fmt"

// Interpreter evaluates expressions.
type Interpreter struct{}

// NewInterpreter creates a new instance of the Interpreter.
func NewInterpreter() *Interpreter {
	return &Interpreter{}
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
