package main

import (
	"fmt"
	"reflect"
)

type interpreter struct {
}

func (i interpreter) interpret(expr Expr) {
	value, err := i.evaluate(expr)
	if err != nil {
		fmt.Printf("Eval expression failed, reason: %v\n", err)
	} else {
		fmt.Printf("Expr value is: %v\n", value)
	}
}

func (i interpreter) evaluate(expr Expr) (interface{}, error) {
	return expr.acceptEvalVisitor(i)
}

func (i interpreter) isTruthy(obj interface{}) bool {
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

func (i interpreter) isEqual(obj1, obj2 interface{}) bool {
	if obj1 == nil && obj2 == nil {
		return true
	}
	if obj1 == nil || obj2 == nil {
		return false
	}
	return reflect.DeepEqual(obj1, obj2)
}

func (i interpreter) checkNumber(obj interface{}) (float64, error) {
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

func (i interpreter) checkNumbers(obj1, obj2 interface{}) (float64, float64, error) {
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

func (i interpreter) checkStrings(obj1, obj2 interface{}) (string, string, error) {
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

func (i interpreter) visitBinaryExpr(expr BinaryExpr) (interface{}, error) {
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

func (i interpreter) visitUnaryExpr(expr UnaryExpr) (interface{}, error) {
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

func (i interpreter) visitLiteralExpr(expr LiteralExpr) (interface{}, error) {
	return expr.value, nil
}

func (i interpreter) visitGroupingExpr(expr GroupingExpr) (interface{}, error) {
	return i.evaluate(expr.expression)
}
