package main

import (
	"fmt"
	"strings"
)

type PrettyPrinter struct {
}

func (p *PrettyPrinter) parenthesize(name string, exprs ...Expr) string {
	sb := strings.Builder{}
	sb.Write([]byte("("))
	sb.Write([]byte(name))
	for _, expr := range exprs {
		sb.Write([]byte(" "))
		sb.Write([]byte(expr.acceptStringVisitor(p)))
	}
	sb.Write([]byte(")"))
	return sb.String()
}

func (p *PrettyPrinter) visitBinaryExpr(expr *BinaryExpr) string {
	return p.parenthesize(expr.operator.Lexeme, expr.left, expr.right)
}

func (p *PrettyPrinter) visitUnaryExpr(expr *UnaryExpr) string {
	return p.parenthesize(expr.operator.Lexeme, expr.right)
}

func (p *PrettyPrinter) visitLiteralExpr(expr *LiteralExpr) string {
	if expr.value == nil {
		return "nil"
	}
	return fmt.Sprintf("%v", expr.value)
}

func (p *PrettyPrinter) visitGroupingExpr(expr *GroupingExpr) string {
	return p.parenthesize("group", expr.expression)
}

func (p *PrettyPrinter) visitVarExpr(expr *VarExpr) string {
	return fmt.Sprintf("var %s", expr.name)
}

func (p *PrettyPrinter) visitAssignExpr(expr *AssignExpr) string {
	return fmt.Sprint("%s = %v", expr.name.Lexeme, expr.expr)
}

func (p *PrettyPrinter) visitLogicalExpr(expr *LogicalExpr) string {
	return fmt.Sprint("%v %s %v", expr.left, expr.operator.Lexeme, expr.right)
}

func (p *PrettyPrinter) visitCallExpr(expr *CallExpr) string {
	return fmt.Sprint("%v %v %v", expr.callee, expr.paren, expr.args)
}

func (p *PrettyPrinter) visitGetExpr(expr *GetExpr) string {
	return fmt.Sprint("%v %v", expr.object, expr.name)
}

func (p *PrettyPrinter) visitSetExpr(expr *SetExpr) string {
	return fmt.Sprint("%v %v %v", expr.object, expr.name, expr.value)
}

func (p *PrettyPrinter) visitThisExpr(expr *ThisExpr) string {
	return fmt.Sprint("%v", expr.keyword)
}

func (p *PrettyPrinter) visitSuperExpr(expr *SuperExpr) string {
	return fmt.Sprint("%v %v", expr.keyword, expr.method)
}
