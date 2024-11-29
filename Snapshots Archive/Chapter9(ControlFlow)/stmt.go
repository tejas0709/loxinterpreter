package main

// Stmt is the interface for all statement types.
type Stmt interface {
	Accept(visitor StmtVisitor) interface{}
}

// StmtVisitor is the interface for visiting each statement type.
type StmtVisitor interface {
	VisitExpressionStmt(stmt *ExpressionStmt) interface{}
	VisitPrintStmt(stmt *PrintStmt) interface{}
	VisitVarStmt(stmt *VarStmt) interface{}
	VisitBlockStmt(stmt *BlockStmt) interface{}
	VisitIfStmt(stmt *IfStmt) interface{}
	VisitWhileStmt(stmt *WhileStmt) interface{}
	VisitBreakStmt(stmt *BreakStmt) interface{}

}

// ExpressionStmt represents an expression as a statement.
type ExpressionStmt struct {
	Expression Expr
}

func (stmt *ExpressionStmt) Accept(visitor StmtVisitor) interface{} {
	return visitor.VisitExpressionStmt(stmt)
}

// PrintStmt represents a print statement.
type PrintStmt struct {
	Expression Expr
}

func (stmt *PrintStmt) Accept(visitor StmtVisitor) interface{} {
	return visitor.VisitPrintStmt(stmt)
}

// VarStmt represents a variable declaration.
type VarStmt struct {
	Name        Token
	Initializer Expr
}

func (stmt *VarStmt) Accept(visitor StmtVisitor) interface{} {
	return visitor.VisitVarStmt(stmt)
}

// BlockStmt represents a block of statements.
type BlockStmt struct {
	Statements []Stmt
}

func (stmt *BlockStmt) Accept(visitor StmtVisitor) interface{} {
	return visitor.VisitBlockStmt(stmt)
}

type IfStmt struct {
    Condition  Expr
    ThenBranch Stmt
    ElseBranch Stmt
}
func (stmt *IfStmt) Accept(visitor StmtVisitor) interface{} {
    return visitor.VisitIfStmt(stmt)
}

type WhileStmt struct {
    Condition Expr
    Body      Stmt
}
func (stmt *WhileStmt) Accept(visitor StmtVisitor) interface{} {
    return visitor.VisitWhileStmt(stmt)
}

// BreakStmt represents a break statement.
type BreakStmt struct {}

func (stmt *BreakStmt) Accept(visitor StmtVisitor) interface{} {
    return visitor.VisitBreakStmt(stmt)
}
