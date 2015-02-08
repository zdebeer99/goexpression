package goexpression

import (
	"testing"
)

type mathTestValue struct {
	input    string
	haserror bool
	result   float64
}

func TestMathEval(t *testing.T) {

	var (
		x float64 = 5
		y float64 = 23
		z float64 = 12.25
	)
	context := map[string]interface{}{
		"x": x,
		"y": y,
		"z": z,
	}
	var tests []mathTestValue = []mathTestValue{
		//Test Basic Calculations
		{"1+1", false, 2},
		{"-1+2", false, 1},
		{"2-1", false, 1},
		{"1-10", false, -9},
		{"1+2*3", false, 7},
		{"2*3+1", false, 7},
		{"2*3/2", false, 2 * 3 / 2},
		{"2/2*3", false, 2 / 2 * 3},
		//Testing precedence
		{"1+2*3/2", false, 1 + 2*3/2},
		{"-3+1.5*2+5-2*2", false, -3 + 1.5*2 + 5 - 2*2},
		{"4+3-2+1", false, 4 + 3 - 2 + 1},
		{"2-3+4-2", false, 2 - 3 + 4 - 2},
		{"2.4*3+1.5*2-3.1*4-1+2", false, 2.4*3 + 1.5*2 - 3.1*4 - 1 + 2},
		//Testing brackets
		{"(1+2)*3", false, (1 + 2) * 3},
		{"3*(1+2)", false, 3 * (1 + 2)},
		{"3*(1+2)*4", false, 3 * (1 + 2) * 4},
		//Testing expressions with variables. Where {x:5}
		{"2*x", false, 2 * x},
		{"2*x+y+(z+x)*4", false, float64(2*x + y + (z+x)*4)},
		{"1+x*(50-y)/z", false, float64(1 + x*(50-y)/z)},
	}
	for i, v := range tests {
		ans := Eval(v.input, context)
		//t.Log(i, ") ", v.input, "=", ans, " Expecting:", v.result)
		if int(ans*1000) != int(v.result*1000) {
			t.Error(i, ") ", v.input, "=", ans, " Expecting:", v.result)
		}

	}
}
