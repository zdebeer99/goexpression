package goexpression

type expression struct {
	ast *TreeNode
}

// Bug(zdebeer): functions is eval from right to left instead from left to right.
func Eval(input string) float64 {
	expr := &expression{Parse(input)}
	return expr.eval(expr.ast)
}

func (this *expression) eval(basenode *TreeNode) float64 {
	for _, node := range basenode.items {
		switch node.Value.Category() {
		case CatFunction:
			return this.switchFunction(node)
		case CatValue:
			return this.getNumber(node)
		case CatOther:
			this.switchOther(node)
		}
	}
	panic("eval failed. f.")
}

func (this *expression) switchOther(node *TreeNode) {
	switch v1 := node.Value.(type) {
	case *GroupToken:
		if v1.GroupType == "()" {
			this.eval(node)
			return
		}
	}
	panic("Invalid Node " + node.String())
}

func (this *expression) switchFunction(node *TreeNode) float64 {
	val1 := node.Value.(*FuncToken)
	switch val1.Name {
	case "+":
		return this.evalMathOperator(this.evalMathPlus, node.Items())
	case "-":
		return this.evalMathOperator(this.evalMathMinus, node.Items())
	case "*":
		return this.evalMathOperator(this.evalMathMultiply, node.Items())
	case "/":
		return this.evalMathOperator(this.evalMathDevide, node.Items())
	default:
		panic("Function not supported")
	}

}

func (this *expression) getNumber(node *TreeNode) float64 {
	switch v := node.Value.(type) {
	case *NumberToken:
		return v.Value
	case *FuncToken:
		return this.switchFunction(node)
	case *GroupToken:
		if v.GroupType == "()" {
			return this.eval(node)
		}
		panic("Unexpected grouping type: " + node.String())
	default:
		panic("Unexpected value: " + node.String())
	}
}

func (this *expression) evalMathOperator(fn func(float64, float64) float64, args []*TreeNode) float64 {
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
