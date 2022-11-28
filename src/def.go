package engine

import (
	"errors"
	"fmt"
	"math"
	"reflect"
)

const (
	RadianMode = iota
	AngleMode
)

type DefineFunc struct {
	argc     int
	fun      func(params map[string]float64, args ...ExprNode) float64
	funLaTex func(args ...ExprNode) string
}

var defaultLaTexFunc = func(args ...ExprNode) string {
	return ""
}

// TrigonometricMode enum "RadianMode", "AngleMode"
var TrigonometricMode = RadianMode

var defConst = map[string]float64{
	"pi":    math.Pi,
	"e":     math.E,
	"infty": 0,
}

var defConstLaTex = map[string]string{
	"pi":    "π",
	"e":     "e",
	"infty": "\\infty",
}

var defFunc map[string]DefineFunc

func init() {
	defFunc = map[string]DefineFunc{
		"sin": {1, defSin, defSinLaTex},
		"cos": {1, defCos, defCosLaTex},
		"tan": {1, defTan, defTanLaTex},
		"cot": {1, defCot, defCotLaTex},
		"sec": {1, defSec, defSecLaTex},
		"csc": {1, defCsc, defCscLaTex},

		"abs":   {1, defAbs, defAbsLaTex},
		"ceil":  {1, defCeil, defaultLaTexFunc},
		"floor": {1, defFloor, defaultLaTexFunc},
		"round": {1, defRound, defaultLaTexFunc},
		"sqrt":  {1, defSqrt, defSqrtLaTex},
		"cbrt":  {1, defCbrt, defaultLaTexFunc},

		"noerr": {1, defNoerr, defaultLaTexFunc},

		"max": {-1, defMax, defaultLaTexFunc},
		"min": {-1, defMin, defaultLaTexFunc},

		"sum": {-1, defSum, defSumLaTex},

		"log": {2, defLog, defLogLaTex},
	}
}

// sin(pi/2) = 1

func defSin(params map[string]float64, expr ...ExprNode) float64 {
	return math.Sin(expr2Radian(expr[0], params))
}

func defSinLaTex(args ...ExprNode) string {
	return fmt.Sprintf("sin(%s)", ExprASTLaTex(args[0]))
}

// cos(0) = 1

func defCos(params map[string]float64, expr ...ExprNode) float64 {
	return math.Cos(expr2Radian(expr[0], params))
}

func defCosLaTex(args ...ExprNode) string {
	return fmt.Sprintf("cos(%s)", ExprASTLaTex(args[0]))
}

// tan(pi/4) = 1

func defTan(params map[string]float64, expr ...ExprNode) float64 {
	return math.Tan(expr2Radian(expr[0], params))
}

func defTanLaTex(args ...ExprNode) string {
	return fmt.Sprintf("tan(%s)", ExprASTLaTex(args[0]))
}

// cot(pi/4) = 1

func defCot(params map[string]float64, expr ...ExprNode) float64 {
	return 1 / defTan(params, expr...)
}

func defCotLaTex(args ...ExprNode) string {
	return fmt.Sprintf("cot(%s)", ExprASTLaTex(args[0]))
}

// sec(0) = 1

func defSec(params map[string]float64, expr ...ExprNode) float64 {
	return 1 / defCos(params, expr...)
}

func defSecLaTex(args ...ExprNode) string {
	return fmt.Sprintf("sec(%s)", ExprASTLaTex(args[0]))
}

// csc(pi/2) = 1

func defCsc(params map[string]float64, expr ...ExprNode) float64 {
	return 1 / defSin(params, expr...)
}

func defCscLaTex(args ...ExprNode) string {
	return fmt.Sprintf("csc(%s)", ExprASTLaTex(args[0]))
}

// abs(-2) = 2

func defAbs(params map[string]float64, expr ...ExprNode) float64 {
	return math.Abs(ExprASTResult(expr[0], params))
}

func defAbsLaTex(args ...ExprNode) string {
	return fmt.Sprintf("|%s|", ExprASTLaTex(args[0]))
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

func defSqrtLaTex(args ...ExprNode) string {
	return fmt.Sprintf("\\sqrt{%s}", ExprASTLaTex(args[0]))
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

// sum(0) = 1

func defSum(params map[string]float64, expr ...ExprNode) float64 {
	if len(expr) < 2 {
		panic(errors.New("calling function `sum` must have at least two parameter."))
	}

	name := reflect.TypeOf(expr[1]).Name()
	if name != "NumberExprNode" {
		panic(errors.New("calling function `sum` cannot be computed efficiently"))
	}

	if len(expr) == 2 {
		sumV := 0.0
		start := expr[0].(NumberExprNode)
		end := expr[1].(NumberExprNode)
		for i := start.Val; i <= end.Val; i++ {
			sumV = sumV + i
		}
		return sumV
	}
	sumV := 0.0
	start := expr[0].(NumberExprNode)
	end := expr[1].(NumberExprNode)
	for i := int(start.Val); i <= int(end.Val); i++ {
		params["#i"] = float64(i)
		v := ExprASTResult(expr[2], params)
		sumV = sumV + v
	}
	delete(params, "#i")
	return sumV
}

func defSumLaTex(args ...ExprNode) string {
	if len(args) < 2 {
		panic(errors.New("calling function `sum` must have at least two parameter."))
	}

	if len(args) == 2 {
		return fmt.Sprintf("\\sum_{i=%s}^{%s} k", ExprASTLaTex(args[0]), ExprASTLaTex(args[1]))
	}

	return fmt.Sprintf("\\sum_{i=%s}^{%s} %s", ExprASTLaTex(args[0]), ExprASTLaTex(args[1]), ExprASTLaTex(args[2]))
}

// log
func defLog(params map[string]float64, expr ...ExprNode) float64 {
	if len(expr) != 2 {
		panic(errors.New("calling function `log` must have two parameter."))
	}

	a := ExprASTResult(expr[0], params)
	b := ExprASTResult(expr[1], params)
	return math.Log10(b) / math.Log10(a)
}

func defLogLaTex(args ...ExprNode) string {
	if len(args) != 2 {
		panic(errors.New("calling function `log` must have two parameter."))
	}

	return fmt.Sprintf("\\log_{%s}^{%s}", ExprASTLaTex(args[0]), ExprASTLaTex(args[1]))
}
