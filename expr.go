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
	VisitGetExpr(expr *GetExpr) interface{}
	VisitSetExpr(expr *SetExpr) interface{}
	VisitThisExpr(expr *ThisExpr) interface{}
	VisitSuperExpr(expr *SuperExpr) interface{}
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

// GetExpr represents a get (property access) expression
type GetExpr struct {
	Object Expr
	Name   Token
}

func (g *GetExpr) Accept(visitor ExprVisitor) interface{} {
	return visitor.VisitGetExpr(g)
}

// SetExpr represents a set (property assignment) expression
type SetExpr struct {
	Object Expr
	Name   Token
	Value  Expr
}

func (s *SetExpr) Accept(visitor ExprVisitor) interface{} {
	return visitor.VisitSetExpr(s)
}

// ThisExpr represents the 'this' keyword in method contexts
type ThisExpr struct {
	Keyword Token
}

func (t *ThisExpr) Accept(visitor ExprVisitor) interface{} {
	return visitor.VisitThisExpr(t)
}

// SuperExpr represents the 'super' keyword for method calls
type SuperExpr struct {
	Keyword Token
	Method  Token
}

func (s *SuperExpr) Accept(visitor ExprVisitor) interface{} {
	return visitor.VisitSuperExpr(s)
}