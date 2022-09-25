package main

import "fmt"

type FunctionType int

const (
	FunctionTypeNone = iota
	FuntionTypeFunction
)

// 主要是为了做 semantic analysis
type resolver struct {
	interpreter         Interpreter
	scopes              *stack
	currentFunctionType FunctionType
}

func newResolver(intp Interpreter) *resolver {
	return &resolver{
		interpreter:         intp,
		scopes:              newStack(),
		currentFunctionType: FunctionTypeNone,
	}
}

func (r *resolver) beginScope() error {
	// bool 类型的 value 代表 string 类型的 key 是否完成了初始化
	scope := make(map[string]bool)
	r.scopes.Push(scope)
	return nil
}

func (r *resolver) endScope() error {
	_, err := r.scopes.Pop()
	return err
}

func (r *resolver) resolveStmts(stmts []Stmt) error {
	for _, stmt := range stmts {
		if err := r.resolveStmt(stmt); err != nil {
			return err
		}
	}
	return nil
}

func (r *resolver) resolveStmt(stmt Stmt) error {
	return stmt.acceptStmtVisitor(r)
}

func (r *resolver) resolveExpr(expr Expr) error {
	_, err := expr.acceptEvalVisitor(r)
	if err != nil {
		return err
	}
	return nil
}

func (r *resolver) declare(name token) error {
	if r.scopes.IsEmpty() {
		return nil
	}
	v, err := r.scopes.Peek()
	if err != nil {
		return err
	}
	scope, ok := v.(map[string]bool)
	if !ok {
		return errCastToMapString2Bool
	}
	// 这里用的都是指针，所以值会做相应的变化。
	if _, ok := scope[name.Lexeme]; ok {
		return fmt.Errorf("already a variable with this name %s in this scope", name.Lexeme)
	}
	scope[name.Lexeme] = false
	return nil
}

func (r *resolver) define(name token) error {
	if r.scopes.IsEmpty() {
		return nil
	}
	v, err := r.scopes.Peek()
	if err != nil {
		return err
	}
	scope, ok := v.(map[string]bool)
	if !ok {
		return errCastToMapString2Bool
	}
	// 注意值跟 declare 的区别
	scope[name.Lexeme] = true
	return nil
}

func (r *resolver) visitBinaryExpr(expr *BinaryExpr) (interface{}, error) {
	if err := r.resolveExpr(expr.left); err != nil {
		return nil, err
	}
	if err := r.resolveExpr(expr.right); err != nil {
		return nil, err
	}
	return nil, nil
}

func (r *resolver) visitUnaryExpr(expr *UnaryExpr) (interface{}, error) {
	return nil, r.resolveExpr(expr.right)
}

func (r *resolver) visitLiteralExpr(expr *LiteralExpr) (interface{}, error) {
	return nil, nil
}

func (r *resolver) visitGroupingExpr(expr *GroupingExpr) (interface{}, error) {
	return nil, r.resolveExpr(expr.expression)
}

func (r *resolver) visitVarExpr(expr *VarExpr) (interface{}, error) {
	if !r.scopes.IsEmpty() {
		e, err := r.scopes.Peek()
		if err != nil {
			return nil, err
		}
		scope, ok := e.(map[string]bool)
		if !ok {
			return nil, errCastToMapString2Bool
		}
		// 说明变量名称之前已经定义过了（但是没有初始化）
		if v, ok := scope[expr.name.Lexeme]; ok && v == false {
			return nil, fmt.Errorf("cannot read local varibale %s in its own initliazer", expr.name.Lexeme)
		}
	}
	if err := r.resolveLocal(expr, expr.name); err != nil {
		return nil, err
	}
	return nil, nil
}

func (r *resolver) visitAssignExpr(expr *AssignExpr) (interface{}, error) {
	if err := r.resolveExpr(expr.expr); err != nil {
		return nil, err
	}
	if err := r.resolveLocal(expr, expr.name); err != nil {
		return nil, err
	}
	return nil, nil
}

func (r *resolver) visitLogicalExpr(expr *LogicalExpr) (interface{}, error) {
	if err := r.resolveExpr(expr.left); err != nil {
		return nil, err
	}
	if err := r.resolveExpr(expr.right); err != nil {
		return nil, err
	}
	return nil, nil
}

func (r *resolver) visitCallExpr(expr *CallExpr) (interface{}, error) {
	if err := r.resolveExpr(expr.callee); err != nil {
		return nil, err
	}
	for _, arg := range expr.args {
		if err := r.resolveExpr(arg); err != nil {
			return nil, err
		}
	}
	return nil, nil
}

func (r *resolver) visitPrintStmt(stmt PrintStmt) error {
	return r.resolveExpr(stmt.expr)
}

func (r *resolver) visitExpressionStmt(stmt ExpressionStmt) error {
	return r.resolveExpr(stmt.expr)
}

func (r *resolver) visitVarStmt(stmt VarStmt) error {
	if err := r.declare(stmt.name); err != nil {
		return err
	}
	if stmt.expr != nil {
		if err := r.resolveExpr(stmt.expr); err != nil {
			return err
		}
	}
	if err := r.define(stmt.name); err != nil {
		return err
	}
	return nil
}

func (r *resolver) visitBlockStmt(stmt BlockStmt) error {
	if err := r.beginScope(); err != nil {
		return err
	}
	if err := r.resolveStmts(stmt.stmts); err != nil {
		return err
	}
	if err := r.endScope(); err != nil {
		return err
	}
	return nil
}

func (r *resolver) visitIFStmt(stmt IFStmt) error {
	if err := r.resolveExpr(stmt.condition); err != nil {
		return err
	}
	if err := r.resolveStmt(stmt.thenBranch); err != nil {
		return err
	}
	if stmt.elseBranch != nil {
		if err := r.resolveStmt(stmt.elseBranch); err != nil {
			return err
		}
	}
	return nil
}

func (r *resolver) visitWhileStmt(stmt WhileStmt) error {
	if err := r.resolveExpr(stmt.condition); err != nil {
		return err
	}
	if err := r.resolveStmt(stmt.body); err != nil {
		return err
	}
	return nil
}

func (r *resolver) visitFunctionStmt(stmt FunctionStmt) error {
	if err := r.declare(stmt.name); err != nil {
		return err
	}
	if err := r.define(stmt.name); err != nil {
		return err
	}
	if err := r.resolveFunction(stmt, FuntionTypeFunction); err != nil {
		return err
	}
	return nil
}

func (r *resolver) visitReturnStmt(stmt ReturnStmt) error {
	if r.currentFunctionType == FunctionTypeNone {
		return fmt.Errorf("cannot return from top-level code")
	}
	if stmt.value != nil {
		return r.resolveExpr(stmt.value)
	}
	return nil
}

func (r *resolver) resolveFunction(stmt FunctionStmt, functionType FunctionType) error {
	enclosingFunction := r.currentFunctionType
	r.currentFunctionType = functionType
	defer func() {
		r.currentFunctionType = enclosingFunction
	}()
	if err := r.beginScope(); err != nil {
		return err
	}
	for _, param := range stmt.params {
		if err := r.declare(param); err != nil {
			return err
		}
		if err := r.define(param); err != nil {
			return err
		}
	}
	if err := r.resolveStmts(stmt.stmts); err != nil {
		return err
	}
	if err := r.endScope(); err != nil {
		return err
	}
	return nil
}

func (r *resolver) resolveLocal(expr Expr, name token) error {
	// 表示当前的 scope 深度和发现变量所在的 scope 的距离
	// 本质为了解决 closure 内的变量和 global 变量名字相同的问题。
	var distance int
	for e := r.scopes.Back(); e != nil; e = e.Prev() {
		scope, ok := e.Value.(map[string]bool)
		if !ok {
			return errCastToMapString2Bool
		}
		if _, ok := scope[name.Lexeme]; ok {
			if err := r.interpreter.Resolve(expr, distance); err != nil {
				return err
			}
			return nil
		}
		distance++
	}
	return nil
}
