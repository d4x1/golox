package main

import (
	"fmt"
)

const (
	maxArgsCount = 128

	typeFunction = "function"
)

// syntactic analysis
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

func (p *parser) or() (Expr, error) {
	expr, err := p.and()
	if err != nil {
		return nil, err
	}
	for {
		ok := p.match(OR)
		if !ok {
			break
		}
		operator := p.previous()
		right, err := p.and()
		if err != nil {
			return nil, err
		}
		expr = newLogicalExpr(operator, expr, right)
	}
	return expr, nil
}

func (p *parser) and() (Expr, error) {
	expr, err := p.equality()
	if err != nil {
		return nil, err
	}
	for {
		ok := p.match(AND)
		if !ok {
			break
		}
		operator := p.previous()
		right, err := p.equality()
		if err != nil {
			return nil, err
		}
		expr = newLogicalExpr(operator, expr, right)
	}
	return expr, nil
}

func (p *parser) assignment() (Expr, error) {
	expr, err := p.or()
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
			return nil, fmt.Errorf("token: %s, invalid assign target", equalToken)
		}
	}
	return expr, nil
}

func (p *parser) declaration() (Stmt, error) {
	// 这里可以单独处理下错误，如果当前语句解析出错，还可以继续解析。
	if p.match(FUN) {
		return p.function(typeFunction)
	}
	if p.match(VAR) {
		return p.varDeclaration()
	}
	return p.statement()
}

func (p *parser) function(kind string) (Stmt, error) {
	name, ok := p.consume(IDENTIFIER)
	if !ok {
		p.parseErr(name, fmt.Sprintf("expect '%s' name", kind))
		return nil, fmt.Errorf("expect '%s' name", kind)
	}
	if token, ok := p.consume(LEFT_PAREN); !ok {
		p.parseErr(token, fmt.Sprintf("expect '(' after %s name", kind))
		return nil, fmt.Errorf("expect '(' after %s name", kind)
	}
	var args []token
	ok = p.check(RIGHT_PAREN)
	if !ok {
		for {
			if len(args) > maxArgsCount {
				return nil, fmt.Errorf("token: %v cannot have more than %d args", p.peek(), maxArgsCount)
			}
			if token, ok := p.consume(IDENTIFIER); !ok {
				p.parseErr(token, "expect parameter name")
				return nil, fmt.Errorf("expect parameter name")
			} else {
				args = append(args, token)
			}
			if !p.match(COMMA) {
				break
			}
		}
	}

	if token, ok := p.consume(RIGHT_PAREN); !ok {
		p.parseErr(token, fmt.Sprintf("expect ')' after %s name", kind))
		return nil, fmt.Errorf("expect ')' after %s name", kind)
	}

	if token, ok := p.consume(LEFT_BRACE); !ok {
		p.parseErr(token, fmt.Sprintf("expect '{' after %s name", kind))
		return nil, fmt.Errorf("expect ')' after %s name", kind)
	}
	block, err := p.block()
	if err != nil {
		return nil, err
	}
	// block 中已经检查过 } 了，所以这里不需要再检查。
	return newFunctionStmt(name, args, block), nil
}

func (p *parser) statement() (Stmt, error) {
	if p.match(PRINT) {
		return p.printStatement()
	}
	if p.match(RETURN) {
		return p.returnStatement()
	}
	if p.match(IF) {
		return p.ifStatement()
	}
	if p.match(WHILE) {
		return p.whileStatement()
	}
	if p.match(FOR) {
		return p.forStatement()
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

func (p *parser) returnStatement() (Stmt, error) {
	keyword := p.previous()
	var value Expr
	if !p.check(SEMICOLON) {
		var err error
		value, err = p.expression()
		if err != nil {
			return nil, err
		}
	}
	token, ok := p.consume(SEMICOLON)
	if !ok {
		p.parseErr(token, "expect ';' after return value")
		return nil, fmt.Errorf("expect ';' after return value")
	}
	return newReturnStmt(keyword, value), nil
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

// 这里其实做了一个 `de-sugaring` 的操作，for loop 复用了 while 底层实现。
func (p *parser) forStatement() (Stmt, error) {
	token, ok := p.consume(LEFT_PAREN)
	if !ok {
		p.parseErr(token, "expect '(' after expression")
		return nil, fmt.Errorf("expect '(' after expression")
	}
	var initializer Stmt
	if p.match(SEMICOLON) {

	} else if p.match(VAR) {
		var err error
		initializer, err = p.varDeclaration()
		if err != nil {
			return nil, err
		}
	} else {
		var err error
		initializer, err = p.expressionStatement()
		if err != nil {
			return nil, err
		}
	}
	var condition Expr
	if p.check(SEMICOLON) {
		// condition 不存在
	} else {
		var err error
		condition, err = p.expression()
		if err != nil {
			return nil, err
		}
	}
	token, ok = p.consume(SEMICOLON)
	if !ok {
		p.parseErr(token, "expect ';' after expression")
		return nil, fmt.Errorf("expect ';' after expression")
	}

	var increment Expr
	if p.check(RIGHT_PAREN) {
		// increment 不存在
	} else {
		var err error
		increment, err = p.expression()
		if err != nil {
			return nil, err
		}
	}
	token, ok = p.consume(RIGHT_PAREN)
	if !ok {
		p.parseErr(token, "expect ')' after expression")
		return nil, fmt.Errorf("expect ')' after expression")
	}
	body, err := p.statement()
	if err != nil {
		return nil, err
	}
	if increment != nil {
		body = newBlockStmt([]Stmt{body, newExpressionStmt(increment)})
	}
	if condition == nil {
		condition = newLiteralExpr(true)
	}
	body = newWhileStmt(condition, body)
	if initializer != nil {
		body = newBlockStmt([]Stmt{initializer, body})
	}
	return body, nil
}

func (p *parser) whileStatement() (Stmt, error) {
	token, ok := p.consume(LEFT_PAREN)
	if !ok {
		p.parseErr(token, "expect '(' after expression")
		return nil, fmt.Errorf("expect '(' after expression")
	}
	condition, err := p.expression()
	if err != nil {
		return nil, err
	}
	token, ok = p.consume(RIGHT_PAREN)
	if !ok {
		p.parseErr(token, "expect ')' after expression")
		return nil, fmt.Errorf("expect ')' after expression")
	}
	body, err := p.statement()
	if err != nil {
		return nil, err
	}
	return newWhileStmt(condition, body), nil
}

func (p *parser) ifStatement() (Stmt, error) {
	token, ok := p.consume(LEFT_PAREN)
	if !ok {
		p.parseErr(token, "expect '(' after expression")
		return nil, fmt.Errorf("expect '(' after expression")
	}
	condition, err := p.expression()
	if err != nil {
		return nil, err
	}
	token, ok = p.consume(RIGHT_PAREN)
	if !ok {
		p.parseErr(token, "expect ')' after expression")
		return nil, fmt.Errorf("expect ')' after expression")
	}
	thenBranch, err := p.statement()
	if err != nil {
		return nil, err
	}
	var elseBranch Stmt
	if p.match(ELSE) {
		var err error
		elseBranch, err = p.statement()
		if err != nil {
			return nil, err
		}
	}
	return newIFStmt(condition, thenBranch, elseBranch), nil
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
	return p.call()
}

func (p *parser) call() (Expr, error) {
	expr, err := p.primary()
	if err != nil {
		return nil, err
	}
	for {
		ok := p.match(LEFT_PAREN)
		if ok {
			var err error
			expr, err = p.finishCall(expr)
			if err != nil {
				return nil, err
			}
		} else {
			break
		}
	}
	return expr, nil
}

func (p *parser) finishCall(callee Expr) (Expr, error) {
	var args []Expr
	ok := p.check(RIGHT_PAREN)
	// 这里实际上处理了空参数的 case
	if !ok {
		for {
			if len(args) > maxArgsCount {
				return nil, fmt.Errorf("token: %v cannot have more than %d args", p.peek(), maxArgsCount)
			}
			expr, err := p.expression()
			if err != nil {
				return nil, err
			}
			args = append(args, expr)
			if !p.match(COMMA) {
				break
			}
		}
	}
	paren, ok := p.consume(RIGHT_PAREN)
	if !ok {
		p.parseErr(paren, "expect ')' after arguments")
		return nil, fmt.Errorf("expect ')' after arguments")
	}
	return newCallExpr(callee, paren, args), nil
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
