package goexpression

import (
	"container/list"
	s "github.com/zdebeer99/goexpression/scanner"
)

type expression struct {
	ast *s.TreeNode
}

// Bug(zdebeer): functions is eval from right to left instead from left to right.
func Eval(input string) float64 {
	expr := &expression{Parse(input)}
	return expr.eval()
}

func (this *expression) eval() float64 {
	for el := this.ast.FirstChild(); el != nil; el = el.Next() {
		node, ok := el.Value.(*s.TreeNode)
		if !ok {
			panic("Invalid Type stored in tree")
		}
		switch node.Value.Category() {
		case s.CatFunction:
			return this.switchFunction(node)
		case s.CatValue:
			return this.getNumber(node)
		}
	}
	panic("eval failed. f.")
}

func (this *expression) switchFunction(node *s.TreeNode) float64 {
	val1 := node.Value.(*s.FuncToken)
	switch val1.Name {
	case "+":
		return this.evalMathOperator(this.evalMathPlus, node.FirstChild())
	case "-":
		return this.evalMathOperator(this.evalMathMinus, node.FirstChild())
	case "*":
		return this.evalMathOperator(this.evalMathMultiply, node.FirstChild())
	case "/":
		return this.evalMathOperator(this.evalMathDevide, node.FirstChild())
	default:
		panic("Function not supported")
	}

}

func (this *expression) getNumber(val *s.TreeNode) float64 {
	switch v := val.Value.(type) {
	case *s.NumberToken:
		return v.Value
	case *s.FuncToken:
		return this.switchFunction(val)
	default:
		panic(":(")
	}
}

func (this *expression) evalMathOperator(fn func(float64, float64) float64, args *list.Element) float64 {

	arg1 := args.Value.(*s.TreeNode)
	args = args.Next()
	if args == nil {
		panic("Operator Missing Arguments.")
	}
	arg2 := args.Value.(*s.TreeNode)
	args = args.Next()
	if args == nil {
		return fn(this.getNumber(arg1), this.getNumber(arg2))
	}

	answ := fn(this.getNumber(arg1), this.getNumber(arg2))
	for ; args != nil; args = args.Next() {
		arg3 := args.Value.(*s.TreeNode)
		answ = fn(answ, this.getNumber(arg3))
	}
	return answ
}

func (this *expression) evalMathPlus(val1, val2 float64) float64 {
	return val1 + val2
}

func (this *expression) evalMathMinus(val1, val2 float64) float64 {
	return val1 - val2
}

func (this *expression) evalMathMultiply(val1, val2 float64) float64 {
	return val1 * val2
}

func (this *expression) evalMathDevide(val1, val2 float64) float64 {
	return val1 / val2
}
