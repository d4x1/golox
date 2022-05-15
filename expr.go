package main

type Visitor interface {
	visitBinaryExpr(expr BinaryExpr) string
	visitUnaryExpr(expr UnaryExpr) string
	visitLiteralExpr(expr LiteralExpr) string
	visitGroupingExpr(expr GroupingExpr) string
}

type Expr interface {
	acceptStringVisitor(visitor Visitor) string
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
