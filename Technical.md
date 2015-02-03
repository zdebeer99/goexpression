# goexpression technical overview

goexpression is meant to be easy to understand and taken apart to be reused in other
projects. This project can form the base of your own parser.

This document gives a overview of the workings of the code.

# Components

The Expression Parser just like most parser have 3 main components. a Scanner, Parser and
evaluator.

**Scanner**

In goexpression the scanner is actually a string iterator, removing the step
of first tokenizing the string before parsing it. The Scanner code is a simplified version of the Scanner code in golang/text/template. We re-purposed it to be a String Iterator by removing all the tokenizing functions and only keeping the iterations functions like Next, Backup, etc. 

The Idea of replacing the traditional scanner with a string iterator is that if your string iterator can simplify navigating through the string easy enough the tokenize step should not be needed. This decision will propably not work with more complex scanners but for simple scanners this is a nice win as your script syntax rules is only required in the parser not in two separate places. Making it easily to adapt for other customs scripts.

**Parser**

The parser build a tree with nodes.


**Evaluator**

The Evaluator runs the tree and returns the result.

**Node**

The Node is used to build the tree structure of code. The current implementation is not golike as I used OO techniquas in golang (Bad Idea). The Node will be rewritten in more golike langauge.

# golang Review

As one of the reasons I am coding the project is to test the language go, I decided that it is fair to add my five cents of opinion about golang here.





