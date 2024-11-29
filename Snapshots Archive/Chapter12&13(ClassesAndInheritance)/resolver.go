package main

import "fmt"

type Resolver struct {
	interpreter    *Interpreter
	scopes         []map[string]bool
	currentFunction FunctionType
	currentClass   ClassType
}

func NewResolver(interpreter *Interpreter) *Resolver {
	return &Resolver{
		interpreter:    interpreter,
		scopes:         []map[string]bool{},
		currentFunction: FunctionNone,
	}
}

func (r *Resolver) Resolve(statements []Stmt) {
	for _, statement := range statements {
		r.resolveStatement(statement)
	}
}

func (r *Resolver) resolveStatement(stmt Stmt) {
	stmt.Accept(r)
}

func (r *Resolver) resolveExpression(expr Expr) {
	expr.Accept(r)
}

func (r *Resolver) beginScope() {
	r.scopes = append(r.scopes, map[string]bool{})
}

func (r *Resolver) endScope() {
	r.scopes = r.scopes[:len(r.scopes)-1]
}

func (r *Resolver) declare(name Token) {
	if len(r.scopes) == 0 {
		return
	}
	scope := r.scopes[len(r.scopes)-1]
	if _, exists := scope[name.Lexeme]; exists {
		panic(fmt.Sprintf("Variable with name '%s' already declared in this scope.", name.Lexeme))
	}
	scope[name.Lexeme] = false
}

func (r *Resolver) define(name Token) {
	if len(r.scopes) == 0 {
		return
	}
	r.scopes[len(r.scopes)-1][name.Lexeme] = true
}

func (r *Resolver) resolveLocal(expr Expr, name Token) {
	for i := len(r.scopes) - 1; i >= 0; i-- {
        if _, exists := r.scopes[i][name.Lexeme]; exists {
            r.interpreter.resolve(expr, len(r.scopes)-1-i)
            return
        }
    }
	// Not found; assume it’s global.
}

func (r *Resolver) resolveFunction(stmt *FunStmt, functionType FunctionType) {
    enclosingFunction := r.currentFunction
    r.currentFunction = functionType
    r.beginScope()
    for _, param := range stmt.Params {
        r.declare(param)
        r.define(param)
    }
    r.Resolve(stmt.Body)
    r.endScope()
    r.currentFunction = enclosingFunction
}

// Statement visitors
func (r *Resolver) VisitBlockStmt(stmt *BlockStmt) interface{} {
	r.beginScope()
	r.Resolve(stmt.Statements)
	r.endScope()
	return nil
}

func (r *Resolver) VisitVarStmt(stmt *VarStmt) interface{} {
	r.declare(stmt.Name)
	if stmt.Initializer != nil {
		r.resolveExpression(stmt.Initializer)
	}
	r.define(stmt.Name)
	return nil
}

func (r *Resolver) VisitExpressionStmt(stmt *ExpressionStmt) interface{} {
	r.resolveExpression(stmt.Expression)
	return nil
}

func (r *Resolver) VisitIfStmt(stmt *IfStmt) interface{} {
	r.resolveExpression(stmt.Condition)
	r.resolveStatement(stmt.ThenBranch)
	if stmt.ElseBranch != nil {
		r.resolveStatement(stmt.ElseBranch)
	}
	return nil
}

func (r *Resolver) VisitPrintStmt(stmt *PrintStmt) interface{} {
	r.resolveExpression(stmt.Expression)
	return nil
}

func (r *Resolver) VisitReturnStmt(stmt *ReturnStmt) interface{} {
	if r.currentFunction == FunctionNone {
		panic("Cannot return from top-level code.")
	}
	if stmt.Value != nil {
		r.resolveExpression(stmt.Value)
	}
	return nil
}

func (r *Resolver) VisitWhileStmt(stmt *WhileStmt) interface{} {
	r.resolveExpression(stmt.Condition)
	r.resolveStatement(stmt.Body)
	return nil
}

func (r *Resolver) VisitFunStmt(stmt *FunStmt) interface{} {
    r.declare(stmt.Name)
    r.define(stmt.Name)
    r.resolveFunction(stmt, FunctionFunction)
    return nil
}

func (r *Resolver) VisitBreakStmt(stmt *BreakStmt) interface{} {
	return nil
}

// Expression visitors
func (r *Resolver) VisitBinaryExpr(expr *Binary) interface{} {
	r.resolveExpression(expr.Left)
	r.resolveExpression(expr.Right)
	return nil
}

func (r *Resolver) VisitGroupingExpr(expr *Grouping) interface{} {
	r.resolveExpression(expr.Expression)
	return nil
}

func (r *Resolver) VisitLiteralExpr(expr *Literal) interface{} {
	return nil
}

func (r *Resolver) VisitUnaryExpr(expr *Unary) interface{} {
	r.resolveExpression(expr.Right)
	return nil
}

func (r *Resolver) VisitVariableExpr(expr *Variable) interface{} {
	if len(r.scopes) > 0 {
		scope := r.scopes[len(r.scopes)-1]
		if defined, exists := scope[expr.Name.Lexeme]; exists && !defined {
			panic(fmt.Sprintf("Cannot read local variable '%s' in its own initializer.", expr.Name.Lexeme))
		}
	}
	r.resolveLocal(expr, expr.Name)
	return nil
}

func (r *Resolver) VisitAssignExpr(expr *Assign) interface{} {
	r.resolveExpression(expr.Value)
	r.resolveLocal(expr, expr.Name)
	return nil
}

func (r *Resolver) VisitCallExpr(expr *Call) interface{} {
	r.resolveExpression(expr.Callee)
	for _, arg := range expr.Arguments {
		r.resolveExpression(arg)
	}
	return nil
}

type FunctionType int
const (
	FunctionNone FunctionType = iota
	FunctionFunction
	FunctionMethod
	FunctionInitializer
)

type ClassType int
const (
	ClassNone ClassType = iota
	ClassClass
	ClassSubclass
)

func (r *Resolver) VisitClassStmt(stmt *ClassStmt) interface{} {
    enclosingClass := r.currentClass
    r.currentClass = ClassClass

    r.declare(stmt.Name)
    r.define(stmt.Name)

    if stmt.Superclass != nil {
        r.currentClass = ClassSubclass
        r.resolveExpression(stmt.Superclass)
        r.beginScope()
        r.scopes[len(r.scopes)-1]["super"] = true
    }

    r.beginScope()
    r.scopes[len(r.scopes)-1]["this"] = true

    for _, method := range stmt.Methods {
        declaration := FunctionMethod
        if method.Name.Lexeme == "init" {
            declaration = FunctionInitializer
        }
        r.resolveFunction(method, declaration)
    }

    r.endScope()
    if stmt.Superclass != nil {
        r.endScope()
    }

    r.currentClass = enclosingClass
    return nil
}

func (r *Resolver) VisitThisExpr(expr *ThisExpr) interface{} {
    if r.currentClass == ClassNone {
        panic("Cannot use 'this' outside of a class.")
    }
    r.resolveLocal(expr, expr.Keyword)
    return nil
}

func (r *Resolver) VisitSuperExpr(expr *SuperExpr) interface{} {
    if r.currentClass == ClassNone {
        panic("Cannot use 'super' outside of a class.")
    } else if r.currentClass != ClassSubclass {
        panic("Cannot use 'super' in a class with no superclass.")
    }
    r.resolveLocal(expr, expr.Keyword)
    return nil
}

func (r *Resolver) VisitGetExpr(expr *GetExpr) interface{} {
    r.resolveExpression(expr.Object)
    return nil
}

func (r *Resolver) VisitSetExpr(expr *SetExpr) interface{} {
    r.resolveExpression(expr.Value)
    r.resolveExpression(expr.Object)
    return nil
}