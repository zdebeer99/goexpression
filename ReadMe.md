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

**Planned**

- variables Ex: 1+x where x is passed as a variable to Eval.
- Boolean types and expressions like '==', '<', '>', '!', 'and', 'or'
- Special keywords like. if, each, etc
- Calling go functions from the expression


## Basic usage

    ans:=goexpression.Eval("1+2")
	fmt.Printf("1+2=%v",ans)

