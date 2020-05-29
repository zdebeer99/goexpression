package goexpression

import (
	"testing"

	s "github.com/zdebeer99/goexpression/scanner"
)

type TestValue struct {
	value    string
	haserror bool
	result   string
}

func TestParseExpression(t *testing.T) {
	//These tests only test parsing scenarios, eval tests can better tests for algorithmic correctness.
	var testValues []TestValue = []TestValue{
		{"", false, "Base()"},
		{"1", false, "[Base():Number(1)]"},                                                                       //1 test basic parse
		{"-1", false, "[Base():Number(-1)]"},                                                                     //2 test negative parse
		{"1+2", false, "[Base():[Func(+):Number(1),Number(2)]]"},                                                 //3 test basic expression
		{"-1+2", false, "[Base():[Func(+):Number(-1),Number(2)]]"},                                               //4 test expression starting with negative number
		{"1-2 ", false, "[Base():[Func(-):Number(1),Number(2)]]"},                                                //5 test negative expression.
		{"1- 2 ", false, "[Base():[Func(-):Number(1),Number(2)]]"},                                               //6 test with different combination of spaces.
		{" 1 -2 ", false, "[Base():[Func(-):Number(1),Number(2)]]"},                                              //7 test with different combination of spaces.
		{" 1 - 2 ", false, "[Base():[Func(-):Number(1),Number(2)]]"},                                             //8 test with different combination of spaces.
		{"1+2*2", false, "[Base():[Func(+):Number(1),[Func(*):Number(2),Number(2)]]]"},                           //9 test operator presedence
		{"1+2^2", false, "[Base():[Func(+):Number(1),[Func(^):Number(2),Number(2)]]]"},                           //9 test operator presedence
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
		//26 Test Parsing of functions with brackets.
		{"3*(1+2)", false, "[Base():[Func(*):Number(3),[Group(()):[Func(+):Number(1),Number(2)]]]]"},             //26
		{"(1+2)*3", false, "[Base():[Func(*):[Group(()):[Func(+):Number(1),Number(2)]],Number(3)]]"},             //27
		{"4*(1+2)*3", false, "[Base():[Func(*):Number(4),[Group(()):[Func(+):Number(1),Number(2)]],Number(3)]]"}, //28
		//29 Test Parsing of functions with variables.
		{"1+x", false, "[Base():[Func(+):Number(1),Identity(x)]]"}, //29 Test parsing of variabe
		{"2*x+y+(z+x)*4", false, "[Base():[Func(+):[Func(*):Number(2),Identity(x)],Identity(y),[Func(*):[Group(()):[Func(+):Identity(z),Identity(x)]],Number(4)]]]"}, //30
		{"x+1 y", true, "[Base():[Func(+):Identity(x),Number(1),[Base():[ERROR: Line: 1, near \" y\", Error: Unexpected end of expression. '' not parsed. ]]]]"},     //31 Test expression parsing stopping after end of expression.
		{"x=y*3", false, "[Base():[Func(=):Identity(x),[Group():[Func(*):Identity(y),Number(3)]]]]"},                                                                 //32
		{"'hello'", false, "[Base():\"hello\"]"}, //33
		{"function(var1,var2+3,'Hello')", false, "[Base():Func function([[Group():Identity(var1)] [Group():[Func(+):Identity(var2),Number(3)]] [Group():\"Hello\"]])]"},                                             //34
		{"function(var1,var2+3,'Hello')+3*4", false, "[Base():[Func(+):Func function([[Group():Identity(var1)] [Group():[Func(+):Identity(var2),Number(3)]] [Group():\"Hello\"]]),[Func(*):Number(3),Number(4)]]]"}, //35
		{"1==1", false, "[Base():[Func(==):Number(1),Number(1)]]"},
		{"1==1 and 2==2", false, "[Base():[Func(and):[Func(==):Number(1),Number(1)],[Func(==):Number(2),Number(2)]]]"},
		{"1==1 && 2==2 || 5!=6", false, "[Base():[Func(||):[Func(&&):[Func(==):Number(1),Number(1)],[Func(==):Number(2),Number(2)]],[Func(!=):Number(5),Number(6)]]]"},
		{"1*2+3==3+1*2", false, "[Base():[Func(==):[Func(+):[Func(*):Number(1),Number(2)],Number(3)],[Func(+):Number(3),[Func(*):Number(1),Number(2)]]]]"},
		{"1*2+3==3+1*2 && (5!=6 || 5==6)", false, "[Base():[Func(&&):[Func(==):[Func(+):[Func(*):Number(1),Number(2)],Number(3)],[Func(+):Number(3),[Func(*):Number(1),Number(2)]]],[Group(()):[Func(||):[Func(!=):Number(5),Number(6)],[Func(==):Number(5),Number(6)]]]]]"},
	}

	for i, v := range testValues {
		node, err := Parse(v.value)
		if err != nil && !v.haserror {
			t.Errorf("%v. %q:\nError:%s\nparsed:%s", i, v.value, err.Error(), node)
			continue
		}
		if err == nil && v.haserror {
			t.Errorf("%v. %q:\nNo Error Raised. Expecting and Error but none was returned.:\nparsed:%s", i, v.value, node)
			continue
		}
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
		scan := s.NewScanner(v.value)
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
