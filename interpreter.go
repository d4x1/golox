package main

import (
	"fmt"
	"reflect"
)

type interpreter struct {
	env *Env
}

func newInterpreter() *interpreter {
	i := &interpreter{}
	i.env = newEnv()
	return i
}

func (i *interpreter) interpret(stmts []Stmt) {
	for _, stmt := range stmts {
		if err := i.execute(stmt); err != nil {
			fmt.Printf("Execute stmt: %s failed, reason: %v\n", stmt, err)
			return
		}
	}
	fmt.Println("Execute stmts success!")
}

func (i *interpreter) execute(stmt Stmt) error {
	return stmt.acceptStmtVisitor(i)
}

func (i *interpreter) evaluate(expr Expr) (interface{}, error) {
	return expr.acceptEvalVisitor(i)
}

func (i *interpreter) isTruthy(obj interface{}) bool {
	if obj == nil {
		return false
	}
	switch obj.(type) {
	case bool:
		return obj.(bool)
	default:
		return true
	}
}

func (i *interpreter) isEqual(obj1, obj2 interface{}) bool {
	if obj1 == nil && obj2 == nil {
		return true
	}
	if obj1 == nil || obj2 == nil {
		return false
	}
	return reflect.DeepEqual(obj1, obj2)
}

func (i *interpreter) checkNumber(obj interface{}) (float64, error) {
	switch obj.(type) {
	case uint:
		return float64(obj.(uint)), nil
	case uint8:
		return float64(obj.(uint8)), nil
	case uint16:
		return float64(obj.(uint16)), nil
	case uint32:
		return float64(obj.(uint32)), nil
	case uint64:
		return float64(obj.(uint64)), nil
	case int:
		return float64(obj.(int)), nil
	case int8:
		return float64(obj.(int8)), nil
	case int16:
		return float64(obj.(int16)), nil
	case int32:
		return float64(obj.(int32)), nil
	case int64:
		return float64(obj.(int64)), nil
	case float32:
		return float64(obj.(float32)), nil
	case float64:
		return obj.(float64), nil
	default:
		return 0, fmt.Errorf("%v is not a number", obj)
	}
}

func (i *interpreter) checkNumbers(obj1, obj2 interface{}) (float64, float64, error) {
	obj1Num, err := i.checkNumber(obj1)
	if err != nil {
		return 0, 0, err
	}
	obj2Num, err := i.checkNumber(obj2)
	if err != nil {
		return 0, 0, err
	}
	return obj1Num, obj2Num, nil
}

func (i interpreter) checkString(obj interface{}) (string, error) {
	switch obj.(type) {
	case string:
		return obj.(string), nil
	default:
		return "", fmt.Errorf("%v is not a string", obj)
	}
}

func (i *interpreter) checkStrings(obj1, obj2 interface{}) (string, string, error) {
	obj1Str, err := i.checkString(obj1)
	if err != nil {
		return "", "", err
	}
	obj2Str, err := i.checkString(obj2)
	if err != nil {
		return "", "", err
	}
	return obj1Str, obj2Str, nil
}

func (i *interpreter) visitBinaryExpr(expr BinaryExpr) (interface{}, error) {
	left, err := i.evaluate(expr.left)
	if err != nil {
		return nil, err
	}
	right, err := i.evaluate(expr.right)
	if err != nil {
		return nil, err
	}
	switch expr.operator.Type {
	case GREATER:
		leftNum, rightNum, err := i.checkNumbers(left, right)
		if err != nil {
			return nil, err
		}
		return leftNum > rightNum, nil
	case GREATER_EQUAL:
		leftNum, rightNum, err := i.checkNumbers(left, right)
		if err != nil {
			return nil, err
		}
		return leftNum >= rightNum, nil
	case LESS:
		leftNum, rightNum, err := i.checkNumbers(left, right)
		if err != nil {
			return nil, err
		}
		return leftNum < rightNum, nil
	case LESS_EQUAL:
		leftNum, rightNum, err := i.checkNumbers(left, right)
		if err != nil {
			return nil, err
		}
		return leftNum <= rightNum, nil
	case BANG_EQUAL:
		return !i.isEqual(left, right), nil
	case EQUAL_EQUAL:
		return i.isEqual(left, right), nil
	case MINUS:
		leftNum, rightNum, err := i.checkNumbers(left, right)
		if err != nil {
			return nil, err
		}
		return leftNum - rightNum, nil
	case PLUS:
		leftNum, rightNum, err := i.checkNumbers(left, right)
		if err == nil {
			return leftNum + rightNum, nil
		}
		leftStr, rightStr, err := i.checkStrings(left, right)
		if err == nil {
			return leftStr + rightStr, nil
		}
		return nil, fmt.Errorf("left: %v, right: %v are not the same type(float or string)", left, right)
	case SLASH:
		leftNum, rightNum, err := i.checkNumbers(left, right)
		if err != nil {
			return nil, err
		}
		return leftNum / rightNum, nil
	case STAR:
		leftNum, rightNum, err := i.checkNumbers(left, right)
		if err != nil {
			return nil, err
		}
		return leftNum * rightNum, nil
	default:
		return nil, fmt.Errorf("unkown operator: %v between %v and %v", expr.operator, expr.left, expr.right)
	}
}

func (i *interpreter) visitUnaryExpr(expr UnaryExpr) (interface{}, error) {
	right, err := i.evaluate(expr.right)
	if err != nil {
		return nil, err
	}
	switch expr.operator.Type {
	case BANG:
		return i.isTruthy(right), nil
	case MINUS:
		v, ok := right.(float64)
		if !ok {
			customPanic("not a float")
		}
		return -v, nil
	}
	return nil, fmt.Errorf("cannot eval -(%v)", right)
}

func (i *interpreter) visitLiteralExpr(expr LiteralExpr) (interface{}, error) {
	return expr.value, nil
}

func (i *interpreter) visitGroupingExpr(expr GroupingExpr) (interface{}, error) {
	return i.evaluate(expr.expression)
}

func (i *interpreter) visitVarExpr(expr VarExpr) (interface{}, error) {
	return i.env.Get(expr.name.Lexeme)
}

func (i *interpreter) visitPrintStmt(stmt PrintStmt) error {
	value, err := i.evaluate(stmt.expr)
	if err != nil {
		return err
	}
	fmt.Println(value)
	return nil
}

func (i *interpreter) visitExpressionStmt(stmt ExpressionStmt) error {
	i.evaluate(stmt.expr)
	return nil
}

func (i *interpreter) visitVarStmt(stmt VarStmt) error {
	var value interface{}
	if stmt.expr != nil {
		var err error
		value, err = i.evaluate(stmt.expr)
		if err != nil {
			return err
		}
	}
	i.env.Define(stmt.name.Lexeme, value)
	return nil
}

func (i *interpreter) visitBlockStmt(stmt BlockStmt) error {
	return i.executeBlock(stmt.stmts, newEnvWithEnclosing(i.env))
}

func (i *interpreter) executeBlock(stmts []Stmt, env *Env) error {
	preEnv := i.env
	i.env = env
	for _, stmt := range stmts {
		if err := i.execute(stmt); err != nil {
			return err
		}
	}
	i.env = preEnv
	return nil
}

func (i *interpreter) visitWhileStmt(stmt WhileStmt) error {
	for {
		condition, err := i.evaluate(stmt.condition)
		if err != nil {
			return err
		}
		if i.isTruthy(condition) {
			if err := i.execute(stmt.body); err != nil {
				return err
			}
		} else {
			break
		}
	}
	return nil
}

func (i *interpreter) visitIFStmt(stmt IFStmt) error {
	condition, err := i.evaluate(stmt.condition)
	if err != nil {
		return err
	}
	if i.isTruthy(condition) {
		if i.execute(stmt.thenBranch); err != nil {
			return err
		}
	} else if stmt.elseBranch != nil {
		if i.execute(stmt.elseBranch); err != nil {
			return err
		}
	}
	return nil
}

func (i *interpreter) visitAssignExpr(expr AssignExpr) (interface{}, error) {
	value, err := i.evaluate(expr.expr)
	if err != nil {
		return nil, err
	}
	if err := i.env.Assign(expr.name.Lexeme, value); err != nil {
		return nil, err
	}
	return value, nil
}

func (i *interpreter) visitLogicalExpr(expr LogicalExpr) (interface{}, error) {
	left, err := i.evaluate(expr.left)
	if err != nil {
		return nil, err
	}
	if expr.operator.Type == OR {
		if i.isTruthy(left) {
			return left, nil
		}
	} else {
		if !i.isTruthy(left) {
			return left, nil
		}
	}
	right, err := i.evaluate(expr.right)
	if err != nil {
		return nil, err
	}
	return right, nil
}
