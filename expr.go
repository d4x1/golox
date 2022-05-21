package main

type Visitor interface {
	visitBinaryExpr(expr BinaryExpr) string
	visitUnaryExpr(expr UnaryExpr) string
	visitLiteralExpr(expr LiteralExpr) string
	visitGroupingExpr(expr GroupingExpr) string
}

type EvalVisitor interface {
	visitBinaryExpr(expr BinaryExpr) (interface{}, error)
	visitUnaryExpr(expr UnaryExpr) (interface{}, error)
	visitLiteralExpr(expr LiteralExpr) (interface{}, error)
	visitGroupingExpr(expr GroupingExpr) (interface{}, error)
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
