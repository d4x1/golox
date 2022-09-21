package main

type Visitor interface {
	visitBinaryExpr(expr BinaryExpr) string
	visitUnaryExpr(expr UnaryExpr) string
	visitLiteralExpr(expr LiteralExpr) string
	visitGroupingExpr(expr GroupingExpr) string
	visitVarExpr(expr VarExpr) string
	visitAssignExpr(expr AssignExpr) string
	visitLogicalExpr(expr LogicalExpr) string
	visitCallExpr(expr CallExpr) string
}

type EvalVisitor interface {
	visitBinaryExpr(expr BinaryExpr) (interface{}, error)
	visitUnaryExpr(expr UnaryExpr) (interface{}, error)
	visitLiteralExpr(expr LiteralExpr) (interface{}, error)
	visitGroupingExpr(expr GroupingExpr) (interface{}, error)
	visitVarExpr(expr VarExpr) (interface{}, error)
	visitAssignExpr(expr AssignExpr) (interface{}, error)
	visitLogicalExpr(expr LogicalExpr) (interface{}, error)
	visitCallExpr(expr CallExpr) (interface{}, error)
}

type Expr interface {
	acceptStringVisitor(visitor Visitor) string
	acceptEvalVisitor(visitor EvalVisitor) (interface{}, error)
}
type BinaryExpr struct {
	left, right Expr
	operator    token
}

func newBinaryExpr(left, right Expr, operator token) *BinaryExpr {
	return &BinaryExpr{
		left:     left,
		right:    right,
		operator: operator,
	}
}

func (expr *BinaryExpr) acceptStringVisitor(visitor Visitor) string {
	return visitor.visitBinaryExpr(*expr)
}

func (expr *BinaryExpr) acceptEvalVisitor(visitor EvalVisitor) (interface{}, error) {
	return visitor.visitBinaryExpr(*expr)
}

type UnaryExpr struct {
	right    Expr
	operator token
}

func newUnaryExpr(right Expr, operator token) *UnaryExpr {
	return &UnaryExpr{
		right:    right,
		operator: operator,
	}
}

func (expr *UnaryExpr) acceptStringVisitor(visitor Visitor) string {
	return visitor.visitUnaryExpr(*expr)
}

func (expr *UnaryExpr) acceptEvalVisitor(visitor EvalVisitor) (interface{}, error) {
	return visitor.visitUnaryExpr(*expr)
}

type LiteralExpr struct {
	value interface{}
}

func newLiteralExpr(value interface{}) *LiteralExpr {
	return &LiteralExpr{
		value: value,
	}
}

func (expr *LiteralExpr) acceptStringVisitor(visitor Visitor) string {
	return visitor.visitLiteralExpr(*expr)
}

func (expr *LiteralExpr) acceptEvalVisitor(visitor EvalVisitor) (interface{}, error) {
	return visitor.visitLiteralExpr(*expr)
}

type GroupingExpr struct {
	expression Expr
}

func newGroupingExpr(expr Expr) *GroupingExpr {
	return &GroupingExpr{
		expression: expr,
	}
}

func (expr *GroupingExpr) acceptStringVisitor(visitor Visitor) string {
	return visitor.visitGroupingExpr(*expr)
}

func (expr *GroupingExpr) acceptEvalVisitor(visitor EvalVisitor) (interface{}, error) {
	return visitor.visitGroupingExpr(*expr)
}

type VarExpr struct {
	name token
}

func newVarExpr(name token) *VarExpr {
	return &VarExpr{
		name: name,
	}
}

func (expr *VarExpr) acceptStringVisitor(visitor Visitor) string {
	return visitor.visitVarExpr(*expr)
}

func (expr *VarExpr) acceptEvalVisitor(visitor EvalVisitor) (interface{}, error) {
	return visitor.visitVarExpr(*expr)
}

type AssignExpr struct {
	name token
	expr Expr
}

func newAssignExpr(name token, value Expr) *AssignExpr {
	return &AssignExpr{
		name: name,
		expr: value,
	}
}

func (expr *AssignExpr) acceptStringVisitor(visitor Visitor) string {
	return visitor.visitAssignExpr(*expr)
}

func (expr *AssignExpr) acceptEvalVisitor(visitor EvalVisitor) (interface{}, error) {
	return visitor.visitAssignExpr(*expr)
}

type LogicalExpr struct {
	operator token
	left     Expr
	right    Expr
}

func newLogicalExpr(operator token, left, right Expr) *LogicalExpr {
	return &LogicalExpr{
		operator: operator,
		left:     left,
		right:    right,
	}
}

func (expr *LogicalExpr) acceptStringVisitor(visitor Visitor) string {
	return visitor.visitLogicalExpr(*expr)
}

func (expr *LogicalExpr) acceptEvalVisitor(visitor EvalVisitor) (interface{}, error) {
	return visitor.visitLogicalExpr(*expr)
}

type CallExpr struct {
	callee Expr
	paren  token
	args   []Expr
}

func newCallExpr(callee Expr, paren token, args []Expr) *CallExpr {
	return &CallExpr{
		callee: callee,
		paren:  paren,
		args:   args,
	}
}

func (expr *CallExpr) acceptStringVisitor(visitor Visitor) string {
	return visitor.visitCallExpr(*expr)
}

func (expr *CallExpr) acceptEvalVisitor(visitor EvalVisitor) (interface{}, error) {
	return visitor.visitCallExpr(*expr)
}
