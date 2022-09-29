package main

import (
	"fmt"
	"reflect"
	"strings"
)

type interpreter struct {
	globals *Env
	env     *Env
	locals  map[Expr]int // 即使是同一个 name 的 var，实际上也是不同的 Expr 对象。如果 Expr 实现的 receiver 不是 pointer 的话，就不满足这个约束了。
}

func newInterpreter() *interpreter {
	i := &interpreter{}
	i.globals = newEnv()
	i.env = i.globals
	// init native functions
	i.globals.Define("clock", newNativeFunctionClock())
	// record variables' distance to current env
	i.locals = make(map[Expr]int)
	return i
}

func (i *interpreter) GetGlobalEnv() *Env {
	return i.globals
}

func (i *interpreter) debugEnv() {
	if i == nil {
		return
	}
	if i.env == nil {
		return
	}
	var idx int
	env := i.env
	for {
		fmt.Printf("[DEBUG ENV]idx: %d, env: %+v\n", idx, env)
		idx++
		if env.enclosing == nil {
			break
		}
		env = env.enclosing
	}
}

func (i *interpreter) interpret(stmts []Stmt) {
	for _, stmt := range stmts {
		if err := i.execute(stmt); err != nil {
			fmt.Printf("Execute stmt: %s failed, reason: %v\n", stmt, err)
			return
		}
	}
	fmt.Println(strings.ToUpper("Execute stmts success!"))
}

func (i *interpreter) execute(stmt Stmt) error {
	return stmt.acceptStmtVisitor(i)
}

func (i *interpreter) evaluate(expr Expr) (interface{}, error) {
	return expr.acceptEvalVisitor(i)
}

func (i *interpreter) Resolve(expr Expr, distance int) error {
	return i.resolve(expr, distance)
}

func (i *interpreter) resolve(expr Expr, distance int) error {
	// fmt.Printf("put %s to locals, distane: %d\n", expr, distance)
	i.locals[expr] = distance
	return nil
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

func (i *interpreter) visitBinaryExpr(expr *BinaryExpr) (interface{}, error) {
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

func (i *interpreter) visitUnaryExpr(expr *UnaryExpr) (interface{}, error) {
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

func (i *interpreter) visitLiteralExpr(expr *LiteralExpr) (interface{}, error) {
	return expr.value, nil
}

func (i *interpreter) visitGroupingExpr(expr *GroupingExpr) (interface{}, error) {
	return i.evaluate(expr.expression)
}

func (i *interpreter) visitVarExpr(expr *VarExpr) (interface{}, error) {
	return i.lookupVariable(expr.name, expr)
}

func (i *interpreter) visitGetExpr(expr *GetExpr) (interface{}, error) {
	object, err := i.evaluate(expr.object)
	if err != nil {
		return nil, err
	}
	v, ok := object.(*LoxInstance)
	if !ok {
		return nil, fmt.Errorf("%s is not a LoxInstance", object)
	}
	return v.Get(expr.name)
}

func (i *interpreter) visitSetExpr(expr *SetExpr) (interface{}, error) {
	object, err := i.evaluate(expr.object)
	if err != nil {
		return nil, err
	}
	v, ok := object.(*LoxInstance)
	if !ok {
		return nil, fmt.Errorf("%s is not a LoxInstance, only LoxInstance has fields", object)
	}
	value, err := i.evaluate(expr.value)
	if err != nil {
		return nil, err
	}
	return nil, v.Set(expr.name, value)
}

func (i *interpreter) visitThisExpr(expr *ThisExpr) (interface{}, error) {
	return i.lookupVariable(expr.keyword, expr)
}

func (i *interpreter) lookupVariable(exprName token, expr Expr) (interface{}, error) {
	distance, ok := i.locals[expr]
	if ok {
		return i.env.GetAtByVarName(distance, exprName.Lexeme)
	}
	return i.globals.Get(exprName)
}

func (i *interpreter) visitPrintStmt(stmt PrintStmt) error {
	value, err := i.evaluate(stmt.expr)
	if err != nil {
		return err
	}
	if value != nil {
		fmt.Println(value)
	}
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

func (i *interpreter) ExecuteBlock(stmts []Stmt, env *Env) error {
	return i.executeBlock(stmts, env)
}

// 这里牵扯到 nesting 和 shadowing 的问题。不同的 block 中使用的 env 是不一样的。
func (i *interpreter) executeBlock(stmts []Stmt, env *Env) error {
	preEnv := i.env
	i.env = env
	defer func() {
		i.env = preEnv
	}()
	for _, stmt := range stmts {
		err := i.execute(stmt)
		if err != nil {
			return err
		}
	}
	return nil
}

// 先 Define 后 Assign 的好处是：可以在当前 class 中使用自身。
func (i *interpreter) visitClassStmt(stmt ClassStmt) error {
	i.env.Define(stmt.name.Lexeme, nil)
	methods := make(map[string]*LoxFunction) // 这里是指针会不会有问题？
	for _, method := range stmt.methods {
		var isInitializar bool
		if method.name.Lexeme == "init" {
			isInitializar = true
		}
		function := newLoxFunction(method, i.env, isInitializar)
		methods[method.name.Lexeme] = function
	}
	loxClass := newLoxClassWithMethods(stmt.name.Lexeme, methods)
	return i.env.Assign(stmt.name, loxClass)
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

// if 语句存在一个问题： 如果两个 if 之后，出现了一个 else ，那么 else 属于哪个 if ？
// 这里实际上是认为 else 跟最近的 if 搭配。
// 不同的编程语言解决这个问题都不一样，实际操作很复杂。
func (i *interpreter) visitIFStmt(stmt IFStmt) error {
	condition, err := i.evaluate(stmt.condition)
	if err != nil {
		return err
	}
	if i.isTruthy(condition) {
		if err := i.execute(stmt.thenBranch); err != nil {
			return err
		}
	} else if stmt.elseBranch != nil {
		if err := i.execute(stmt.elseBranch); err != nil {
			return err
		}
	}
	return nil
}

// assgin expr 返回了 value，这个行为值得商榷。
// 会产生很多副作用：比如执行 `a=b;` 的时候，会把赋值之后的值也打印出来。
// 但是移出这个副作用，也比较复杂。这个在当初设计的时候就需要考虑到。
func (i *interpreter) visitAssignExpr(expr *AssignExpr) (interface{}, error) {
	value, err := i.evaluate(expr.expr)
	if err != nil {
		return nil, err
	}
	distance, ok := i.locals[expr]
	if ok {
		if err := i.env.AssignAt(distance, expr.name, value); err != nil {
			return nil, err
		}
	} else {
		if err := i.globals.Assign(expr.name, value); err != nil {
			return nil, err
		}
	}
	return nil, nil
}

func (i *interpreter) visitLogicalExpr(expr *LogicalExpr) (interface{}, error) {
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

func (i *interpreter) visitCallExpr(expr *CallExpr) (interface{}, error) {
	callee, err := i.evaluate(expr.callee)
	if err != nil {
		return nil, err
	}
	var argsList []interface{}
	for _, args := range expr.args {
		arg, err := i.evaluate(args)
		if err != nil {
			return nil, err
		}
		argsList = append(argsList, arg)
	}
	if v, ok := callee.(Callable); ok {
		if len(argsList) != v.Arity() {
			return nil, fmt.Errorf("callable: %s, Expected: %d arguments but got: %d", v, v.Arity(), len(argsList))
		}
		return v.Call(i, argsList)
	}
	return nil, fmt.Errorf("%v is not callable", callee)
}

func (i *interpreter) visitFunctionStmt(stmt FunctionStmt) error {
	function := newLoxFunction(stmt, i.env, false)
	i.env.Define(stmt.name.Lexeme, function)
	return nil
}

func (i *interpreter) visitReturnStmt(stmt ReturnStmt) error {
	var value interface{}
	if stmt.value != nil {
		var err error
		value, err = i.evaluate(stmt.value)
		if err != nil {
			return err
		}
	}
	// 这里使用错误来传递值到适当的调用方
	return NewReturn(value)
}
