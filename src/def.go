package engine

import (
	"errors"
	"math"
)

const (
	RadianMode = iota
	AngleMode
)

type DefineFunc struct {
	argc int
	fun  func(params map[string]float64, args ...ExprNode) float64
}

//	type defS struct {
//		argc int
//		args []ExprNode
//		fun  func(params map[string]float64, expr ...ExprNode) float64
//	}
//
// TrigonometricMode enum "RadianMode", "AngleMode"
var TrigonometricMode = RadianMode

var defConst = map[string]float64{
	"pi": math.Pi,
	"π":  math.Pi,
}

var defFunc map[string]DefineFunc

func init() {
	defFunc = map[string]DefineFunc{
		"sin": {1, defSin},
		"cos": {1, defCos},
		"tan": {1, defTan},
		"cot": {1, defCot},
		"sec": {1, defSec},
		"csc": {1, defCsc},

		"abs":   {1, defAbs},
		"ceil":  {1, defCeil},
		"floor": {1, defFloor},
		"round": {1, defRound},
		"sqrt":  {1, defSqrt},
		"cbrt":  {1, defCbrt},

		"noerr": {1, defNoerr},

		"max": {-1, defMax},
		"min": {-1, defMin},
	}
}

// sin(pi/2) = 1

func defSin(params map[string]float64, expr ...ExprNode) float64 {
	return math.Sin(expr2Radian(expr[0], params))
}

// cos(0) = 1

func defCos(params map[string]float64, expr ...ExprNode) float64 {
	return math.Cos(expr2Radian(expr[0], params))
}

// tan(pi/4) = 1

func defTan(params map[string]float64, expr ...ExprNode) float64 {
	return math.Tan(expr2Radian(expr[0], params))
}

// cot(pi/4) = 1

func defCot(params map[string]float64, expr ...ExprNode) float64 {
	return 1 / defTan(params, expr...)
}

// sec(0) = 1

func defSec(params map[string]float64, expr ...ExprNode) float64 {
	return 1 / defCos(params, expr...)
}

// csc(pi/2) = 1

func defCsc(params map[string]float64, expr ...ExprNode) float64 {
	return 1 / defSin(params, expr...)
}

// abs(-2) = 2

func defAbs(params map[string]float64, expr ...ExprNode) float64 {
	return math.Abs(ExprASTResult(expr[0], params))
}

// ceil(4.2) = ceil(4.8) = 5

func defCeil(params map[string]float64, expr ...ExprNode) float64 {
	return math.Ceil(ExprASTResult(expr[0], params))
}

// floor(4.2) = floor(4.8) = 4

func defFloor(params map[string]float64, expr ...ExprNode) float64 {
	return math.Floor(ExprASTResult(expr[0], params))
}

// round(4.2) = 4
// round(4.6) = 5

func defRound(params map[string]float64, expr ...ExprNode) float64 {
	return math.Round(ExprASTResult(expr[0], params))
}

// sqrt(4) = 2
// sqrt(4) = abs(sqrt(4))
// returns only the absolute value of the result

func defSqrt(params map[string]float64, expr ...ExprNode) float64 {
	return math.Sqrt(ExprASTResult(expr[0], params))
}

// cbrt(27) = 3

func defCbrt(params map[string]float64, expr ...ExprNode) float64 {
	return math.Cbrt(ExprASTResult(expr[0], params))
}

// max(2) = 2
// max(2, 3) = 3
// max(2, 3, 1) = 3

func defMax(params map[string]float64, expr ...ExprNode) float64 {
	if len(expr) == 0 {
		panic(errors.New("calling function `max` must have at least one parameter."))
	}
	if len(expr) == 1 {
		return ExprASTResult(expr[0], params)
	}
	maxV := ExprASTResult(expr[0], params)
	for i := 1; i < len(expr); i++ {
		v := ExprASTResult(expr[i], params)
		maxV = math.Max(maxV, v)
	}
	return maxV
}

// min(2) = 2
// min(2, 3) = 2
// min(2, 3, 1) = 1
func defMin(params map[string]float64, expr ...ExprNode) float64 {
	if len(expr) == 0 {
		panic(errors.New("calling function `min` must have at least one parameter."))
	}
	if len(expr) == 1 {
		return ExprASTResult(expr[0], params)
	}
	maxV := ExprASTResult(expr[0], params)
	for i := 1; i < len(expr); i++ {
		v := ExprASTResult(expr[i], params)
		maxV = math.Min(maxV, v)
	}
	return maxV
}

// noerr(1/0) = 0
// noerr(2.5/(1-1)) = 0
func defNoerr(params map[string]float64, expr ...ExprNode) (r float64) {
	defer func() {
		if e := recover(); e != nil {
			r = 0
		}
	}()
	return ExprASTResult(expr[0], params)
}
