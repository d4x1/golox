package main

import (
	"fmt"
)

type parser struct {
	tokens  []token
	current int
}

func newParser(tokens []token) *parser {
	return &parser{
		tokens:  tokens,
		current: 0,
	}
}

func (p *parser) parse() ([]Stmt, error) {
	var stmts []Stmt
	for !p.isAtEnd() {
		stmt, err := p.declaration()
		if err != nil {
			return stmts, err
		}
		stmts = append(stmts, stmt)
	}
	return stmts, nil
}

func (p *parser) expression() (Expr, error) {
	return p.assignment()
}

func (p *parser) assignment() (Expr, error) {
	expr, err := p.equality()
	if err != nil {
		return nil, err
	}
	if p.match(EQUAL) {
		equalToken := p.previous()
		value, err := p.assignment()
		if err != nil {
			return nil, err
		}
		if v, ok := expr.(*VarExpr); ok {
			name := v.name
			return newAssignExpr(name, value), nil
		} else {
			return nil, fmt.Errorf("token: %s, invalid assgin target", equalToken)
		}
	}
	return expr, nil
}

func (p *parser) declaration() (Stmt, error) {
	// 这里可以单独处理下错误，如果当前语句解析出错，还可以继续解析。
	if p.match(VAR) {
		return p.varDeclaration()
	}
	return p.statement()
}

func (p *parser) statement() (Stmt, error) {
	if p.match(PRINT) {
		return p.printStatement()
	}
	if p.match(LEFT_BRACE) {
		stmts, err := p.block()
		if err != nil {
			return nil, err
		}
		return newBlockStmt(stmts), nil
	}
	return p.expressionStatement()
}

func (p *parser) block() ([]Stmt, error) {
	var stmts []Stmt
	for !p.check(RIGHT_BRACE) && !p.isAtEnd() {
		declaration, err := p.declaration()
		if err != nil {
			return nil, err
		}
		stmts = append(stmts, declaration)
	}
	name, ok := p.consume(RIGHT_BRACE)
	if !ok {
		p.parseErr(name, "expect '}' after block")
		return nil, fmt.Errorf("expect '}' after block")
	}
	return stmts, nil
}

func (p *parser) varDeclaration() (Stmt, error) {
	name, ok := p.consume(IDENTIFIER)
	if !ok {
		p.parseErr(name, "expect variable name")
		return nil, fmt.Errorf("expect variable name")
	}
	var expr Expr
	if p.match(EQUAL) {
		var err error
		expr, err = p.expression()
		if err != nil {
			return nil, err
		}
	}
	token, ok := p.consume(SEMICOLON)
	if !ok {
		p.parseErr(token, "expect ';' after value")
		return nil, fmt.Errorf("expect ';' after expression")
	}
	return newVarStmt(name, expr), nil
}

func (p *parser) printStatement() (Stmt, error) {
	value, err := p.expression()
	if err != nil {
		return nil, err
	}
	token, ok := p.consume(SEMICOLON)
	if !ok {
		p.parseErr(token, "expect ';' after value")
		return nil, fmt.Errorf("expect ';' after expression")
	}
	return newPrintStmt(value), nil
}

func (p *parser) expressionStatement() (Stmt, error) {
	expr, err := p.expression()
	if err != nil {
		return nil, err
	}
	token, ok := p.consume(SEMICOLON)
	if !ok {
		p.parseErr(token, "expect ';' after expression")
		return nil, fmt.Errorf("expect ';' after expression")
	}
	return newPrintStmt(expr), nil
}

func (p *parser) equality() (Expr, error) {
	expr, err := p.comparison()
	if err != nil {
		return nil, err
	}
	for p.match(BANG_EQUAL, EQUAL_EQUAL) {
		operator := p.previous()
		right, err := p.comparison()
		if err != nil {
			return nil, err
		}
		expr = newBinaryExpr(expr, right, operator)
	}
	return expr, nil
}

func (p *parser) comparison() (Expr, error) {
	expr, err := p.term()
	if err != nil {
		return nil, err
	}
	for p.match(LESS, LESS_EQUAL, GREATER, GREATER_EQUAL) {
		operator := p.previous()
		right, err := p.term()
		if err != nil {
			return nil, err
		}
		expr = newBinaryExpr(expr, right, operator)
	}
	return expr, nil
}

func (p *parser) term() (Expr, error) {
	expr, err := p.factor()
	if err != nil {
		return nil, err
	}
	for p.match(MINUS, PLUS) {
		operator := p.previous()
		right, err := p.factor()
		if err != nil {
			return nil, err
		}
		expr = newBinaryExpr(expr, right, operator)
	}
	return expr, nil
}

func (p *parser) factor() (Expr, error) {
	expr, err := p.unary()
	if err != nil {
		return nil, err
	}
	for p.match(SLASH, STAR) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		expr = newBinaryExpr(expr, right, operator)
	}
	return expr, nil
}

func (p *parser) unary() (Expr, error) {
	if p.match(BANG, MINUS) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		return newUnaryExpr(right, operator), nil
	}
	return p.primary()
}

func (p *parser) primary() (Expr, error) {
	if p.match(FALSE) {
		return newLiteralExpr(false), nil
	} else if p.match(TRUE) {
		return newLiteralExpr(true), nil
	} else if p.match(NIL) {
		return newLiteralExpr(nil), nil
	} else if p.match(STRING, NUMBER) {
		return newLiteralExpr(p.previous().literal), nil
	} else if p.match(IDENTIFIER) {
		return newVarExpr(p.previous()), nil
	} else if p.match(LEFT_PAREN) {
		expr, err := p.expression()
		if err != nil {
			return nil, err
		}
		token, ok := p.consume(RIGHT_PAREN)
		if !ok {
			// 这里可以暴露错误，直接 return。
			// 也可以把错误先记下，尝试用正确的方式解析。
			p.parseErr(token, "expect ')' after expression")
			return nil, fmt.Errorf("expect ')' after expression")
		}
		return newGroupingExpr(expr), nil
	}
	token := p.peek()
	return nil, fmt.Errorf("token: %+v, expect expression", token)
}

func (p *parser) consume(tokenType uint) (token, bool) {
	if p.check(tokenType) {
		token := p.advance()
		return token, true
	}
	return token{}, false
}

func (p *parser) parseErr(token token, msg string) {
	fmt.Printf("[Parse Error] token: %+v, err: %s\n", token, msg)
}

func (p *parser) match(tokenTypes ...uint) bool {
	for _, tokenType := range tokenTypes {
		if p.check(tokenType) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *parser) check(tokenType uint) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().Type == tokenType
}

func (p *parser) isAtEnd() bool {
	return p.peek().Type == EOF
}

func (p *parser) peek() token {
	return p.tokens[p.current]
}

func (p *parser) advance() token {
	if !p.isAtEnd() {
		p.current += 1
	}
	return p.previous()
}

func (p *parser) previous() token {
	return p.tokens[p.current-1]
}
