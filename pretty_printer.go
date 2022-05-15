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

func (p *PrettyPrinter) visitBinaryExpr(expr BinaryExpr) string {
	return p.parenthesize(expr.operator.Lexeme, expr.left, expr.right)
}

func (p *PrettyPrinter) visitUnaryExpr(expr UnaryExpr) string {
	return p.parenthesize(expr.operator.Lexeme, expr.right)
}

func (p *PrettyPrinter) visitLiteralExpr(expr LiteralExpr) string {
	if expr.value == nil {
		return "nil"
	}
	return fmt.Sprintf("%v", expr.value)
}

func (p *PrettyPrinter) visitGroupingExpr(expr GroupingExpr) string {
	return p.parenthesize("group", expr.expression)
}