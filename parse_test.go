package goexpression

import (
	"testing"
)

type TestValue struct {
	value    string
	haserror bool
	result   string
}

func TestParseExpression(t *testing.T) {
	//These tests only test parsing scenarios, eval tests can better tests for algorithmic correctness.
	var testValues []TestValue = []TestValue{
		{"", false, "[Base()]"},
		{"1", false, "[Base():Number(1)]"},                                                                       //1 test basic parse
		{"-1", false, "[Base():Number(-1)]"},                                                                     //2 test negative parse
		{"1+2", false, "[Base():[Func(+):Number(1),Number(2)]]"},                                                 //3 test basic expression
		{"-1+2", false, "[Base():[Func(+):Number(-1),Number(2)]]"},                                               //4 test expression starting with negative number
		{"1-2 ", false, "[Base():[Func(-):Number(1),Number(2)]]"},                                                //5 test negative expression.
		{"1- 2 ", false, "[Base():[Func(-):Number(1),Number(2)]]"},                                               //6 test with different combination of spaces.
		{" 1 -2 ", false, "[Base():[Func(-):Number(1),Number(2)]]"},                                              //7 test with different combination of spaces.
		{" 1 - 2 ", false, "[Base():[Func(-):Number(1),Number(2)]]"},                                             //8 test with different combination of spaces.
		{"1+2*2", false, "[Base():[Func(+):Number(1),[Func(*):Number(2),Number(2)]]]"},                           //9 test operator presedence
		{"1*2+2", false, "[Base():[Func(+):[Func(*):Number(1),Number(2)],Number(2)]]"},                           //10 test operator presedence
		{"1*2+2-4", false, "[Base():[Func(-):[Func(+):[Func(*):Number(1),Number(2)],Number(2)],Number(4)]]"},     //11 test operator presedence
		{"1+2+3", false, "[Base():[Func(+):Number(1),Number(2),Number(3)]]"},                                     //12 Test Grouping the same operator under on func.
		{"4+3-2+1", false, "[Base():[Func(+):[Func(-):[Func(+):Number(4),Number(3)],Number(2)],Number(1)]]"},     //13
		{"2-3+4-1", false, "[Base():[Func(-):[Func(+):[Func(-):Number(2),Number(3)],Number(4)],Number(1)]]"},     //14
		{"2*3+2.5*2", false, "[Base():[Func(+):[Func(*):Number(2),Number(3)],[Func(*):Number(2.5),Number(2)]]]"}, //15
		{"2+3*4+5", false, "[Base():[Func(+):[Func(+):Number(2),[Func(*):Number(3),Number(4)]],Number(5)]]"},     //16
		{"-3+1.5*2+5-2*2", false, "[Base():[Func(-):[Func(+):[Func(+):Number(-3),[Func(*):Number(1.5),Number(2)]],Number(5)],[Func(*):Number(2),Number(2)]]]"},
		{"1+2*3/6", false, "[Base():[Func(+):Number(1),[Func(/):[Func(*):Number(2),Number(3)],Number(6)]]]"}, //18
		{"2*3/5+1", false, "[Base():[Func(+):[Func(/):[Func(*):Number(2),Number(3)],Number(5)],Number(1)]]"}, //19
		{"1+2-3+2*4/8+5-6", false, "[Base():[Func(-):[Func(+):[Func(+):[Func(-):[Func(+):Number(1),Number(2)],Number(3)],[Func(/):[Func(*):Number(2),Number(4)],Number(8)]],Number(5)],Number(6)]]"},
		{"2.4*3+1.5*2-3.1*4-1+2", false, "[Base():[Func(+):[Func(-):[Func(-):[Func(+):[Func(*):Number(2.4),Number(3)],[Func(*):Number(1.5),Number(2)]],[Func(*):Number(3.1),Number(4)]],Number(1)],Number(2)]]"}, //21
		{"1+2-3+4", false, "[Base():[Func(+):[Func(-):[Func(+):Number(1),Number(2)],Number(3)],Number(4)]]"},                                                                                                     //22
		{"1-2+3*4/6", false, "[Base():[Func(+):[Func(-):Number(1),Number(2)],[Func(/):[Func(*):Number(3),Number(4)],Number(6)]]]"},                                                                               //23
		{"1-2+3*4/6+7-8", false, "[Base():[Func(-):[Func(+):[Func(+):[Func(-):Number(1),Number(2)],[Func(/):[Func(*):Number(3),Number(4)],Number(6)]],Number(7)],Number(8)]]"},                                   //24
		{"1+2+3*4*5*6/7/8-9-10-11-12", false, "[Base():[Func(-):[Func(+):Number(1),Number(2),[Func(/):[Func(*):Number(3),Number(4),Number(5),Number(6)],Number(7),Number(8)]],Number(9),Number(10),Number(11),Number(12)]]"},
	}

	for i, v := range testValues {
		node := Parse(v.value)
		if node.String() != v.result {
			t.Errorf("%v. %q:\nparsed to:%q\nexpected :%q\n\n", i, v.value, node, v.result)
		}
	}
}

func TestScanNumber(t *testing.T) {
	var testValues []TestValue = []TestValue{
		{"0", false, ""},
		{"1", false, ""},
		{"-1", false, ""},
		{"1.5", false, ""},
		{".5", false, ""},
		{"-2.56", false, ""},
		{"225486.5645", false, ""},
		{"-225486.5645", false, ""},
		{"+1", true, "1"},
	}

	for _, v := range testValues {
		if v.haserror {
			//ignore error inputs for now, need to test them in the future.
			continue
		}
		scan := NewScanner(v.value)
		r1 := scan.ScanNumber()
		if r1 {
			res1 := scan.Commit()
			//fmt.Println(res1)
			if v.value != res1 {
				t.Errorf("Scan Failed input: %q result: %q \n", v.value, res1)
			}
		} else {
			t.Errorf("Could not Scan Value %q", v.value)
		}
	}
}
