package main

// Expr is the interface for all expression types.
type Expr interface {
	Accept(visitor ExprVisitor) interface{}
}

// ExprVisitor interface for visiting each type of expression.
type ExprVisitor interface {
	VisitBinaryExpr(expr *Binary) interface{}
	VisitGroupingExpr(expr *Grouping) interface{}
	VisitLiteralExpr(expr *Literal) interface{}
	VisitUnaryExpr(expr *Unary) interface{}
	VisitVariableExpr(expr *Variable) interface{}
	VisitAssignExpr(expr *Assign) interface{}
	VisitCallExpr(expr *Call) interface{}
}

// Binary expression (e.g., a + b).
type Binary struct {
	Left     Expr
	Operator Token
	Right    Expr
}

func (b *Binary) Accept(visitor ExprVisitor) interface{} {
	return visitor.VisitBinaryExpr(b)
}

// Grouping expression (e.g., (expression)).
type Grouping struct {
	Expression Expr
}

func (g *Grouping) Accept(visitor ExprVisitor) interface{} {
	return visitor.VisitGroupingExpr(g)
}

// Literal expression (e.g., numbers, strings, nil).
type Literal struct {
	Value interface{}
}

func (l *Literal) Accept(visitor ExprVisitor) interface{} {
	return visitor.VisitLiteralExpr(l)
}

// Unary expression (e.g., -a, !a).
type Unary struct {
	Operator Token
	Right    Expr
}

func (u *Unary) Accept(visitor ExprVisitor) interface{} {
	return visitor.VisitUnaryExpr(u)
}

// Variable represents a variable usage (e.g., "a").
type Variable struct {
	Name Token
}

func (v *Variable) Accept(visitor ExprVisitor) interface{} {
	return visitor.VisitVariableExpr(v)
}

// Assign represents an assignment to a variable (e.g., "a = 10").
type Assign struct {
	Name  Token
	Value Expr
}

func (a *Assign) Accept(visitor ExprVisitor) interface{} {
	return visitor.VisitAssignExpr(a)
}

// Call represents a function call expression.
type Call struct {
    Callee    Expr
    Paren     Token
    Arguments []Expr
}

func (c *Call) Accept(visitor ExprVisitor) interface{} {
    return visitor.VisitCallExpr(c)
}
