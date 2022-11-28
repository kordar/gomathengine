package engine

import (
	"log"
	"math"
	"testing"
)

func TestParse(t *testing.T) {
	s := "sum(2, 10, sin(#i)) + log(10, 6) + 10^5.038552608606513 "
	tokens, err := Parse(s)
	if err != nil {
		log.Panicln(err)
	}
	for _, token := range tokens {
		log.Println(token)
	}
	ast := NewAST(tokens, s)
	expression := ast.ParseExpression()
	if ast.Err != nil {
		log.Panicln(ast.Err)
	}
	tex := ExprASTLaTex(expression)
	log.Println(tex)
}

func TestSub(t *testing.T) {
	str := "1+2-3*4sin($x)"
	log.Println(str[11:13])
}

func TestNewAST(t *testing.T) {
	result, err := ParseAndExec("10^0.038552608606513", map[string]float64{"a": math.Pi / 2, "b": 24, "c": 1})
	if err != nil {
		log.Panicln(err)
	}

	log.Println("result = ", result)
}
