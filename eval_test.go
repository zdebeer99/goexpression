package goexpression

import (
	"testing"
)

type mathTestValue struct {
	input    string
	haserror bool
	result   float64
}

func _TestMathEval(t *testing.T) {
	var tests []mathTestValue = []mathTestValue{
		{"1+1", false, 2},
		{"-1+2", false, 1},
		{"2-1", false, 1},
		{"1-10", false, -9},
		{"1+2*3", false, 7},
		{"2*3+1", false, 7},
		{"2*3/2", false, 2 * 3 / 2},
		{"2/2*3", false, 2 / 2 * 3},
		{"1+2*3/2", false, 1 + 2*3/2},
		{"-3+1.5*2+5-2*2", false, -3 + 1.5*2 + 5 - 2*2},
		{"4+3-2+1", false, 4 + 3 - 2 + 1},
		{"2-3+4-2", false, 2 - 3 + 4 - 2},
		{"2.4*3+1.5*2-3.1*4-1+2", false, 2.4*3 + 1.5*2 - 3.1*4 - 1 + 2},
	}

	for i, v := range tests {
		ans := Eval(v.input)
		//t.Log(i, ") ", v.input, "=", ans, " Expecting:", v.result)
		if int(ans*1000) != int(v.result*1000) {
			t.Error(i, ") ", v.input, "=", ans, " Expecting:", v.result)
		}

	}
}
