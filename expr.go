package main

import "fmt"

type Visitor interface {
	visitBinaryExpr(expr *BinaryExpr) string
	visitUnaryExpr(expr *UnaryExpr) string
	visitLiteralExpr(expr *LiteralExpr) string
	visitGroupingExpr(expr *GroupingExpr) string
	visitVarExpr(expr *VarExpr) string
	visitAssignExpr(expr *AssignExpr) string
	visitLogicalExpr(expr *LogicalExpr) string
	visitCallExpr(expr *CallExpr) string
	visitGetExpr(expr *GetExpr) string
	visitSetExpr(expr *SetExpr) string
	visitThisExpr(expr *ThisExpr) string
	visitSuperExpr(expr *SuperExpr) string
}

type EvalVisitor interface {
	visitBinaryExpr(expr *BinaryExpr) (interface{}, error)
	visitUnaryExpr(expr *UnaryExpr) (interface{}, error)
	visitLiteralExpr(expr *LiteralExpr) (interface{}, error)
	visitGroupingExpr(expr *GroupingExpr) (interface{}, error)
	visitVarExpr(expr *VarExpr) (interface{}, error)
	visitAssignExpr(expr *AssignExpr) (interface{}, error)
	visitLogicalExpr(expr *LogicalExpr) (interface{}, error)
	visitCallExpr(expr *CallExpr) (interface{}, error)
	visitGetExpr(expr *GetExpr) (interface{}, error)
	visitSetExpr(expr *SetExpr) (interface{}, error)
	visitThisExpr(expr *ThisExpr) (interface{}, error)
	visitSuperExpr(expr *SuperExpr) (interface{}, error)
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
	return visitor.visitBinaryExpr(expr)
}

func (expr *BinaryExpr) acceptEvalVisitor(visitor EvalVisitor) (interface{}, error) {
	return visitor.visitBinaryExpr(expr)
}

func (expr *BinaryExpr) String() string {
	return fmt.Sprintf("binary expr, left:%s operand:%s right:%s", expr.left, expr.operator, expr.right)
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
	return visitor.visitUnaryExpr(expr)
}

func (expr *UnaryExpr) acceptEvalVisitor(visitor EvalVisitor) (interface{}, error) {
	return visitor.visitUnaryExpr(expr)
}

func (expr *UnaryExpr) String() string {
	return fmt.Sprintf("unary expr: operand:%s right:%s", expr.operator, expr.right)
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
	return visitor.visitLiteralExpr(expr)
}

func (expr *LiteralExpr) acceptEvalVisitor(visitor EvalVisitor) (interface{}, error) {
	return visitor.visitLiteralExpr(expr)
}

func (expr *LiteralExpr) String() string {
	return fmt.Sprintf("literal expr, value:%v", expr.value)
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
	return visitor.visitGroupingExpr(expr)
}

func (expr *GroupingExpr) acceptEvalVisitor(visitor EvalVisitor) (interface{}, error) {
	return visitor.visitGroupingExpr(expr)
}

func (expr *GroupingExpr) String() string {
	return fmt.Sprintf("group expr, expression:%s )", expr.expression)
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
	return visitor.visitVarExpr(expr)
}

func (expr *VarExpr) acceptEvalVisitor(visitor EvalVisitor) (interface{}, error) {
	return visitor.visitVarExpr(expr)
}

func (expr VarExpr) String() string {
	return fmt.Sprintf("var expr, var:%s", expr.name)
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
	return visitor.visitAssignExpr(expr)
}

func (expr *AssignExpr) acceptEvalVisitor(visitor EvalVisitor) (interface{}, error) {
	return visitor.visitAssignExpr(expr)
}

func (expr *AssignExpr) String() string {
	return fmt.Sprintf("assign expr, name:%s = expr:%s", expr.name, expr.expr)
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
	return visitor.visitLogicalExpr(expr)
}

func (expr *LogicalExpr) acceptEvalVisitor(visitor EvalVisitor) (interface{}, error) {
	return visitor.visitLogicalExpr(expr)
}

func (expr *LogicalExpr) String() string {
	return fmt.Sprintf("logical expr, left:%s operator:%s right:%s", expr.left, expr.operator, expr.right)
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
	return visitor.visitCallExpr(expr)
}

func (expr *CallExpr) acceptEvalVisitor(visitor EvalVisitor) (interface{}, error) {
	return visitor.visitCallExpr(expr)
}

func (expr *CallExpr) String() string {
	return fmt.Sprintf("call expr, callee: %s args:%s", expr.callee, expr.args)
}

type GetExpr struct {
	object Expr
	name   token
}

func newGetExpr(object Expr, name token) *GetExpr {
	return &GetExpr{
		object: object,
		name:   name,
	}
}

func (expr *GetExpr) acceptStringVisitor(visitor Visitor) string {
	return visitor.visitGetExpr(expr)
}

func (expr *GetExpr) acceptEvalVisitor(visitor EvalVisitor) (interface{}, error) {
	return visitor.visitGetExpr(expr)
}

func (expr *GetExpr) String() string {
	return fmt.Sprintf("get expr, object: %s name:%s", expr.object, expr.name)
}

type SetExpr struct {
	object Expr
	name   token
	value  Expr
}

func newSetExpr(object Expr, name token, value Expr) *SetExpr {
	return &SetExpr{
		object: object,
		name:   name,
		value:  value,
	}
}

func (expr *SetExpr) acceptStringVisitor(visitor Visitor) string {
	return visitor.visitSetExpr(expr)
}

func (expr *SetExpr) acceptEvalVisitor(visitor EvalVisitor) (interface{}, error) {
	return visitor.visitSetExpr(expr)
}

func (expr *SetExpr) String() string {
	return fmt.Sprintf("set expr, object: %s name:%s value:%s", expr.object, expr.name, expr.value)
}

type ThisExpr struct {
	keyword token
}

func newThisExpr(keyword token) *ThisExpr {
	return &ThisExpr{
		keyword: keyword,
	}
}

func (expr *ThisExpr) acceptStringVisitor(visitor Visitor) string {
	return visitor.visitThisExpr(expr)
}

func (expr *ThisExpr) acceptEvalVisitor(visitor EvalVisitor) (interface{}, error) {
	return visitor.visitThisExpr(expr)
}

func (expr *ThisExpr) String() string {
	return fmt.Sprintf("this expr, keyword: %s", expr.keyword)
}

type SuperExpr struct {
	keyword token
	method  token
}

func newSuperExpr(keyword token, method token) *SuperExpr {
	return &SuperExpr{
		keyword: keyword,
		method:  method,
	}
}

func (expr *SuperExpr) acceptStringVisitor(visitor Visitor) string {
	return visitor.visitSuperExpr(expr)
}

func (expr *SuperExpr) acceptEvalVisitor(visitor EvalVisitor) (interface{}, error) {
	return visitor.visitSuperExpr(expr)
}

func (expr *SuperExpr) String() string {
	return fmt.Sprintf("super expr, keyword: %s, method: %s", expr.keyword, expr.method)
}
