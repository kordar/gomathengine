package engine

import "fmt"

// ExprNode 抽象语法树
type ExprNode interface {
	toStr() string
}

// NumberExprNode 数值节点
type NumberExprNode struct {
	Val float64
	Str string
}

func (n NumberExprNode) toStr() string {
	return fmt.Sprintf(
		"NumberExprNode:%s",
		n.Str,
	)
}

// OperatorExprNode 操作(二叉树)节点
type OperatorExprNode struct {
	Op  string
	Lhs ExprNode
	Rhs ExprNode
}

func (o OperatorExprNode) toStr() string {
	return fmt.Sprintf(
		"OperatorExprNode: (%s %s %s)",
		o.Op,
		o.Lhs.toStr(),
		o.Rhs.toStr(),
	)
}

// FunCallerExprNode 函数表达式节点
type FunCallerExprNode struct {
	Name string
	Arg  []ExprNode
}

func (f FunCallerExprNode) toStr() string {
	return fmt.Sprintf(
		"FunCallerExprNode:%s",
		f.Name,
	)
}

// VariableExprNode 数值节点
type VariableExprNode struct {
	Val string
}

func (v VariableExprNode) toStr() string {
	return fmt.Sprintf(
		"VariableExprNode:%s",
		v.Val,
	)
}

// ConstExprNode 常量节点
type ConstExprNode struct {
	Name string
	Str  string
	Val  float64
}

func (c ConstExprNode) toStr() string {
	return fmt.Sprintf(
		"ConstExprNode:%s=%s",
		c.Name,
		c.Str,
	)
}
