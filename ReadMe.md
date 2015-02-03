# Go Expression

version: 0.01

go expression is a basic math expression parser and evaluator. This project is in
alpha phase and may change significantly in the near future.


## Status

**Current Support**

- supports + - * /
- supports operator precedence


**Planned Futures**

- Support for grouping ()
- Support for variables

## Basic usage

    ans:=goexpression.Eval("1+2")
	fmt.Printf("1+2=%v",ans)


The MIT License (MIT)
