# goexpression technical overview

goexpression is meant to be easy to understand and taken apart to be reused in other
projects. This project can form the base of your own parser.

This document gives a overview of the workings of the code.

# Components

The Expression Parser just like most parser have 3 main components. a Scanner, Parser and
evaluator.

**Scanner**

In goexpression the scanner is actually a string iterator, removing the step
of first tokenizing the string before parsing it. The Scanner code is a simplified version of the Scanner code in golang/text/template. We re-purposed it to be a String Iterator by removing the tokenizing (goroutine) functions and only keeping the iterations functions like Next, Backup and scan helper functions.

**Parser**

The parser builds the ast tree. we use two function types, parseXxx functions and branchXxx functions.

parseXxx functions builds on the ast tree

branchXxx functions decides what parse functions to call next based on the current token.

**Evaluator**

The Evaluator runs the tree and returns the result.

**Node**

The Node is used to build the tree structure of code. The current implementation is not golike as I used OO techniquas in golang (Bad Idea). The Node will be rewritten in more golike langauge.

# golang Review

As one of the reasons I am coding the project is to test the language go, I decided that it is fair to add my five cents of opinion about golang here.

# Todo Next
2015-02-06 Declare a variable or set a variable. Ex: x=6
2015-02-06 Test Expressions parsing, parsing must only parse the first Expression. Ex: "x y" (parse to) "x"


# Change Log
2015-02-07 Refactored parser functions in parseXxx and branchXxx Functions
2015-02-06 Fixed Bug with operator Presedence not working after brackets.
2015-02-05 Busy with variables. Found Presedence bug.
2015-02-05 Added Brackets '()' parsing support, brackets eval not included yet.
2015-02-04 Changed the LinkedList in TreeNode to an Array. (Not having generics makes the Linked List very confusing to use as you need to continues convert types.)
2015-02-03 Changed the Tree Structure from a OO to a more golike (urg) tree structure





