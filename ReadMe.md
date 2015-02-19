# Go Expression

Version: 0.01 Alpha
MIT License (MIT)

go expression is a basic math expression parser and evaluator. This project is in
alpha phase and may change significantly in the near future.

The purpose is to learn go and use the expression parser in other projects.

## Status

**Supported**

- Basic Math Operators like '+', '-', '*', '/'
- Operator precedence, Ex: 1+2*3 = 1+6 = 7
- grouping () Ex: (1+2)*3
- variables Ex: 1+x where x is passed as a variable to Eval.
- only parsing works for;
- Define a variable Ex: x=6

**Parsing only**
- Text values inclused in qoutes.
- Functions with arguments, ex: myfunc(1,2)
- Boolean types and expressions like '==', '<', '>', '!', 'and', 'or'

**Planned**


- Special keywords like. if, each, etc
- Calling go functions from the expression


## Basic usage

	context := map[string]interface{}{
		"x": 5,
		"y": 21,
		"z": 12.5,
	}
    ans:=goexpression.Eval("1+x*(50-y)/z", context)
	fmt.Printf("=%v",ans)

