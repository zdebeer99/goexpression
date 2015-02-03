# Go Expression

version: 0.01

go expresion is a very basic expression parser and evaluator. This project is in
alpha phase and my change significantly in the near future.

## Status

**Current Support**

- supports + - * /
- supports operator precedence


**Planned Futures**

- Support for grouping ()
- Support for variables

## Basic useage

    ans:=goexpression.Eval("1+2")
	fmt.Printf("1+2=%v",ans)

