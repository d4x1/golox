package main

import "fmt"

type StmtVisitor interface {
	visitPrintStmt(PrintStmt) error
	visitExpressionStmt(ExpressionStmt) error
	visitVarStmt(VarStmt) error
	visitBlockStmt(BlockStmt) error
	visitIFStmt(IFStmt) error
	visitWhileStmt(WhileStmt) error
	visitFunctionStmt(FunctionStmt) error
	visitReturnStmt(ReturnStmt) error
	visitClassStmt(ClassStmt) error
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

func (stmt PrintStmt) String() string {
	return fmt.Sprintf("print stmt, expr: %s", stmt.expr)
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

func (stmt ExpressionStmt) String() string {
	return fmt.Sprintf("expression stmt, expr: %s", stmt.expr)
}

type VarStmt struct {
	expr Expr // initializer
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

func (stmt VarStmt) String() string {
	if stmt.expr == nil {
		return fmt.Sprintf("var stmt, %s ;", stmt.name)
	}
	return fmt.Sprintf("var stmt, %s = %s;", stmt.name, stmt.expr)
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

func (stmt BlockStmt) String() string {
	return fmt.Sprintf("block stmt, { %s }", stmt.stmts)
}

type IFStmt struct {
	condition  Expr
	thenBranch Stmt
	elseBranch Stmt
}

func newIFStmt(condition Expr, thenBranch Stmt, elseBranch Stmt) Stmt {
	return IFStmt{
		condition:  condition,
		thenBranch: thenBranch,
		elseBranch: elseBranch,
	}
}

func (stmt IFStmt) acceptStmtVisitor(visitor StmtVisitor) error {
	return visitor.visitIFStmt(stmt)
}

func (stmt IFStmt) String() string {
	return fmt.Sprintf("if stmt, condition: %s, then: %s, else: %s", stmt.condition, stmt.thenBranch, stmt.elseBranch)
}

type WhileStmt struct {
	condition Expr
	body      Stmt
}

func newWhileStmt(condition Expr, body Stmt) Stmt {
	return WhileStmt{
		condition: condition,
		body:      body,
	}
}

func (stmt WhileStmt) acceptStmtVisitor(visitor StmtVisitor) error {
	return visitor.visitWhileStmt(stmt)
}

func (stmt WhileStmt) String() string {
	return fmt.Sprintf("while stmt, condition:(%s), body:{%s}", stmt.condition, stmt.body)
}

type FunctionStmt struct {
	name   token
	params []token
	stmts  []Stmt // body
}

func newFunctionStmt(name token, params []token, body []Stmt) Stmt {
	return FunctionStmt{
		name:   name,
		params: params,
		stmts:  body,
	}
}

func (stmt FunctionStmt) acceptStmtVisitor(visitor StmtVisitor) error {
	return visitor.visitFunctionStmt(stmt)
}

func (stmt FunctionStmt) String() string {
	return fmt.Sprintf("funtion stmt, name: %s, params: %s, body: %s", stmt.name, stmt.params, stmt.stmts)
}

type ReturnStmt struct {
	keyword token
	value   Expr
}

func newReturnStmt(keyword token, value Expr) Stmt {
	return ReturnStmt{
		keyword: keyword,
		value:   value,
	}
}

func (stmt ReturnStmt) acceptStmtVisitor(visitor StmtVisitor) error {
	return visitor.visitReturnStmt(stmt)
}

func (stmt ReturnStmt) String() string {
	return fmt.Sprintf("return stmt, value: %s", stmt.value)
}

type ClassStmt struct {
	name      token
	functions []FunctionStmt
}

func newClassStmt(name token, functions []FunctionStmt) Stmt {
	return ClassStmt{
		name:      name,
		functions: functions,
	}
}

func (stmt ClassStmt) acceptStmtVisitor(visitor StmtVisitor) error {
	return visitor.visitClassStmt(stmt)
}

func (stmt ClassStmt) String() string {
	return fmt.Sprintf("class stmt, name: %s, functions: %s", stmt.name, stmt.functions)
}
