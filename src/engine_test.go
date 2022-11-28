package engine

import (
	"log"
	"math"
	"testing"
)

func TestParse(t *testing.T) {
	tokens, err := Parse("1+2-3*4sin($x , $y )+  $c*1")
	if err != nil {
		log.Panicln(err)
	}
	for _, token := range tokens {
		log.Println(token)
	}
}

func TestSub(t *testing.T) {
	str := "1+2-3*4sin($x)"
	log.Println(str[11:13])
}

func TestNewAST(t *testing.T) {
	result, err := ParseAndExec("12*2+sin($a)", map[string]float64{"a": math.Pi / 2, "b": 24, "c": 1})
	if err != nil {
		log.Panicln(err)
	}

	log.Println("result = ", result)
}
