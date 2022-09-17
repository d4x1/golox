package main

type StmtVisitor interface {
	visitPrintStmt(PrintStmt) error
	visitExpressionStmt(ExpressionStmt) error
	visitVarStmt(VarStmt) error
	visitBlockStmt(BlockStmt) error
}

type Stmt interface {
	acceptStmtVisitor(StmtVisitor) error
}

type PrintStmt struct {
	expr Expr
}

func newPrintStmt(expr Expr) Stmt {
	return PrintStmt{expr: expr}
}

func (stmt PrintStmt) acceptStmtVisitor(visitor StmtVisitor) error {
	return visitor.visitPrintStmt(stmt)
}

type ExpressionStmt struct {
	expr Expr
}

func newExpressionStmt(expr Expr) Stmt {
	return ExpressionStmt{
		expr: expr,
	}
}

func (stmt ExpressionStmt) acceptStmtVisitor(visitor StmtVisitor) error {
	return visitor.visitExpressionStmt(stmt)
}

type VarStmt struct {
	expr Expr
	name token
}

func newVarStmt(name token, expr Expr) Stmt {
	return VarStmt{
		expr: expr,
		name: name,
	}
}

func (stmt VarStmt) acceptStmtVisitor(visitor StmtVisitor) error {
	return visitor.visitVarStmt(stmt)
}

type BlockStmt struct {
	stmts []Stmt
}

func newBlockStmt(stmts []Stmt) Stmt {
	return BlockStmt{
		stmts: stmts,
	}
}

func (stmt BlockStmt) acceptStmtVisitor(visitor StmtVisitor) error {
	return visitor.visitBlockStmt(stmt)
}
