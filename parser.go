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

func (p *parser) parse() (Expr, error) {
	return p.expression()
}

func (p *parser) expression() (Expr, error) {
	return p.equality()
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
