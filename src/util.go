package engine

import (
	"errors"
	"math"
	"strconv"
	"strings"
)

// ParseAndExec Top level function
// Analytical expression and execution
// err is not nil if an error occurs (including arithmetic runtime errors)
func ParseAndExec(s string, params map[string]float64) (r float64, err error) {
	toks, err := Parse(s)
	if err != nil {
		return 0, err
	}
	ast := NewAST(toks, s)
	if ast.Err != nil {
		return 0, ast.Err
	}
	ar := ast.ParseExpression()
	if ast.Err != nil {
		return 0, ast.Err
	}
	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
		}
	}()
	return ExprASTResult(ar, params), err
}

func ErrPos(s string, pos int) string {
	r := strings.Repeat("-", len(s)) + "\n"
	s += "\n"
	for i := 0; i < pos; i++ {
		s += " "
	}
	s += "^\n"
	return r + s + r
}

func expr2Radian(expr ExprNode, params map[string]float64) float64 {
	r := ExprASTResult(expr, params)
	if TrigonometricMode == AngleMode {
		r = r / 180 * math.Pi
	}
	return r
}

// Float64ToStr float64 -> string
func Float64ToStr(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
}

// RegFunction is Top level function
// register a new function to use in expressions
// name: be register function name. the same function name only needs to be registered once.
// argc: this is a number of parameter signatures. should be -1, 0, or a positive integer
//
//	-1 variable-length argument; >=0 fixed numbers argument
//
// fun:  function handler
func RegFunction(name string, argc int, fun func(map[string]float64, ...ExprNode) float64, funLaTex func(...ExprNode) string) error {
	if len(name) == 0 {
		return errors.New("RegFunction name is not empty")
	}
	if argc < -1 {
		return errors.New("RegFunction argc should be -1, 0, or a positive integer")
	}
	if _, ok := defFunc[name]; ok {
		return errors.New("RegFunction name is already exist")
	}
	if funLaTex == nil {
		defFunc[name] = DefineFunc{argc, fun, defaultLaTexFunc}
	} else {
		defFunc[name] = DefineFunc{argc, fun, funLaTex}
	}
	return nil
}

func RegConst(name string, value float64) error {
	if len(name) == 0 {
		return errors.New("RegConst name is not empty")
	}
	if _, ok := defConst[name]; ok {
		return errors.New("RegConst name is already exist")
	}
	defConst[name] = value
	return nil
}

func RegConstLaTex(name string, value string) error {
	if len(name) == 0 {
		return errors.New("RegConstLaTex name is not empty")
	}
	if _, ok := defConstLaTex[name]; ok {
		return errors.New("RegConstLaTex name is already exist")
	}
	defConstLaTex[name] = value
	return nil
}

// ExprASTResult is a Top level function
// AST traversal
// if an arithmetic runtime error occurs, a panic exception is thrown
func ExprASTResult(expr ExprNode, params map[string]float64) float64 {
	var l, r float64
	switch expr.(type) {
	case OperatorExprNode:
		ast := expr.(OperatorExprNode)
		l = ExprASTResult(ast.Lhs, params)
		r = ExprASTResult(ast.Rhs, params)
		return operators[ast.Op[0]].Result(l, r)
	case NumberExprNode:
		return expr.(NumberExprNode).Val
	case ConstExprNode:
		return expr.(ConstExprNode).Val
	case VariableExprNode:
		val := expr.(VariableExprNode).Val
		return params[val]
	case FunCallerExprNode:
		f := expr.(FunCallerExprNode)
		def := defFunc[f.Name]
		return def.fun(params, f.Arg...)
	}

	return 0.0
}

func ExprASTLaTex(expr ExprNode) string {
	var l, r string
	switch expr.(type) {
	case OperatorExprNode:
		ast := expr.(OperatorExprNode)
		l = ExprASTLaTex(ast.Lhs)
		r = ExprASTLaTex(ast.Rhs)
		return operators[ast.Op[0]].ToLaTex(l, r)
	case NumberExprNode:
		return expr.(NumberExprNode).Str
	case ConstExprNode:
		node := expr.(ConstExprNode)
		if defConstLaTex[node.Name] != "" {
			return defConstLaTex[node.Name]
		}
		return expr.(ConstExprNode).Name
	case VariableExprNode:
		return expr.(VariableExprNode).Val[1:]
	case FunCallerExprNode:
		f := expr.(FunCallerExprNode)
		def := defFunc[f.Name]
		return def.funLaTex(f.Arg...)
	}

	return ""
}
