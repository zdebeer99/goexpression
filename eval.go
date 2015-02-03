package goexpression

type expression struct {
	ast Node
}

// Bug(zdebeer): functions is eval from right to left instead from left to right.
func Eval(input string) float64 {
	expr := &expression{Parse(input)}
	return expr.eval()
}

func (this *expression) eval() float64 {
	for _, node := range this.ast.Items() {
		switch node.NodeCat() {
		case CatFunction:
			return this.switchFunction(node.(*FuncNode))
		case CatValue:
			return this.getNumber(node)
		}
	}
	panic(":(")
}

func (this *expression) switchFunction(fnode *FuncNode) float64 {
	switch fnode.name {
	case "+":
		return this.evalMathOperator(this.evalMathPlus, fnode.Items())
	case "-":
		return this.evalMathOperator(this.evalMathMinus, fnode.Items())
	case "*":
		return this.evalMathOperator(this.evalMathMultiply, fnode.Items())
	case "/":
		return this.evalMathOperator(this.evalMathDevide, fnode.Items())
	default:
		panic("Function not supported")
	}

}

func (this *expression) getNumber(val Node) float64 {
	switch v := val.(type) {
	case *NumberNode:
		return v.value
	case *FuncNode:
		return this.switchFunction(v)
	default:
		panic(":(")
	}
}

func (this *expression) evalMathOperator(fn func(float64, float64) float64, args []Node) float64 {
	cnt := len(args)
	switch {
	case cnt < 2:
		panic("Operator Missing Arguments.")
	case cnt == 2:
		return fn(this.getNumber(args[0]), this.getNumber(args[1]))
	default:
		answ := fn(this.getNumber(args[0]), this.getNumber(args[1]))
		for i := 2; i < cnt; i++ {
			answ = fn(answ, this.getNumber(args[i]))
		}
		return answ
	}
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
