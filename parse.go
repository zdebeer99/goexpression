package goexpression

import (
	"fmt"
	s "github.com/zdebeer99/goexpression/scanner"
	"strings"
)

//TODO: Implement Stack to handle drill into functions.

type parser struct {
	scan  *s.Scanner
	root  *TreeNode
	curr  *TreeNode
	err   error
	state func() bool
}

func Parse(input string) (*TreeNode, error) {
	root := NewTreeNode(NewEmptyToken())
	parse := &parser{s.NewScanner(input), root, root, nil, nil}
	parse.parse()
	return root, parse.err
}

func (this *parser) getCurr() Token {
	if this.curr != nil {
		return this.curr.Value
	}
	return nil
}

func (this *parser) parse() {
	this.pumpExpression()
}

func (this *parser) add(token Token) *TreeNode {
	return this.curr.Add(token)
}

func (this *parser) push(token Token) *TreeNode {
	return this.curr.Push(token)
}

func (this *parser) lastNode() *TreeNode {
	return this.curr.LastElement()
}

func (this *parser) parentNode() *TreeNode {
	return this.curr.Parent()
}

func (this *parser) error(err interface{}) {
	var errortxt string
	if val, ok := err.(error); ok {
		errortxt = val.Error()
	} else {
		errortxt = err.(string)
	}
	lasttoken := this.commit()
	debug := fmt.Errorf("Line: %v, Near %q, Error: %s", this.scan.LineNumber(), lasttoken, errortxt)
	this.add(NewErrorToken(debug.Error()))
	this.err = debug
}

func (this *parser) commit() string {
	return this.scan.Commit()
}

//parseOpenBracket
func (this *parser) parseOpenBracket() bool {
	this.curr = this.add(NewGroupToken("()"))
	this.commit()
	return true
}

//parseCloseBracket
func (this *parser) parseCloseBracket() bool {
	for {
		if this.curr.Parent() == nil {
			this.error("Brackets not closed.")
			return true
		}
		v1, ok := this.curr.Parent().Value.(*GroupToken)
		this.curr = this.curr.Parent()
		if ok && v1.GroupType == "()" {
			this.commit()
			this.curr = this.curr.Parent()
			return true
		}
	}
	panic("Should be impossible to reach this point.")
}

//parseOperator
func (this *parser) parseOperator() bool {
	operator := this.commit()
	lastnode := this.lastNode()
	onode, ok := this.getCurr().(*FuncToken)
	//push excisting operator up in tree structure
	if ok {
		//operator is the same current operator ignore
		if onode.Name == operator {
			return true
		}
		//change order for */ presedence
		if onode.OperatorPrecedence(operator) > 0 {
			if lastnode != nil {
				this.curr = lastnode.Push(NewFuncToken(operator))
				return true
			}
		}
		//after */ presedence fallback and continue pushing +- operators from the bottom.
		if onode.OperatorPrecedence(operator) < 0 {
			for {
				v1, ok := this.curr.Parent().Value.(*FuncToken)
				if ok && strings.Index("+-", v1.Name) >= 0 {
					this.curr = this.curr.Parent()
				} else {
					break
				}
			}
		}
		//standard operator push
		this.curr = this.push(NewFuncToken(operator))
		return true
	}
	//set previous found value as argument of the operator
	if lastnode != nil {
		this.curr = lastnode.Push(NewFuncToken(operator))
	} else {
		this.error(fmt.Sprintf("Expecting a value before operator %q", operator))
		this.state = nil
	}
	return true
}

//parseLRFunc
func (this *parser) parseLRFunc() bool {
	operator := this.commit()
	lastnode := this.lastNode()
	if lastnode != nil {
		this.curr = lastnode.Push(NewFuncToken(operator))
	} else {
		this.error(fmt.Sprintf("Expecting a value before operator %q", operator))
		this.state = nil
	}
	return false
}
