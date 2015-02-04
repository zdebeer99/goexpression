package goexpression

import (
	"fmt"
	s "github.com/zdebeer99/goexpression/scanner"
	"strings"
)

//TODO: Implement Stack to handle drill into functions.

type parser struct {
	scan  *s.Scanner
	root  *s.TreeNode
	curr  *s.TreeNode
	err   error
	state func() bool
}

func Parse(input string) *s.TreeNode {
	root := s.NewTreeNode(s.NewEmptyToken())
	parse := &parser{s.NewScanner(input), root, root, nil, nil}
	parse.parse()
	return root
}

func (this *parser) getValue() s.Token {
	if this.curr != nil {
		return this.curr.Value
	}
	return nil
}

func (this *parser) parse() {
	this.parseExpresion()
}

func (this *parser) add(token s.Token) *s.TreeNode {
	return this.curr.Add(token)
}

func (this *parser) push(token s.Token) *s.TreeNode {
	return this.curr.Push(token)
}

func (this *parser) lastNode() *s.TreeNode {
	return this.curr.LastElement()
}

func (this *parser) parentNode() *s.TreeNode {
	return this.curr.Parent()
}

func (l *parser) error(err interface{}) {
	var text string
	if val, ok := err.(error); ok {
		text = val.Error()
	} else {
		text = err.(string)
	}
	debug := fmt.Errorf("Line: %v, Error: %s", l.scan.LineNumber(), text)
	l.add(s.NewErrorToken(debug.Error()))
	l.err = debug
}

func (this *parser) parseExpresion() bool {
	scan := this.scan
	var hasExpression bool
	this.state = this.parseValue
	for this.state != nil {
		if this.err != nil {
			return hasExpression
		}
		if scan.IsEOF() {
			return hasExpression
		}
		if this.state() {
			hasExpression = true
			continue
		}

		r := scan.Next()
		if s.IsSpace(r) {
			scan.Ignore()
			continue
		}
		if scan.IsEOF() {
			return hasExpression
		}
		this.error(fmt.Sprintf("Unexpected character %q", scan.Commit()))
		return hasExpression
	}
	return hasExpression
}

func (this *parser) parseValue() bool {
	scan := this.scan
	if scan.ScanNumber() {
		this.add(s.NewNumberToken(scan.Commit()))
		this.state = this.parseLRFunc
		return true
	}
	if scan.ScanWord() {
		this.add(s.NewIdentityToken(scan.Commit()))
		this.state = this.parseLRFunc
		return true
	}
	return false
}

func (this *parser) parseLRFunc() bool {
	scan := this.scan
	if scan.Accept("+-/*") {
		operator := scan.Commit()
		lastnode := this.lastNode()
		onode, ok := this.getValue().(*s.FuncToken)
		//push excisting operator
		if ok {
			//operator is the same as the previous one.
			if onode.Name == operator {
				this.state = this.parseValue
				return true
			}
			//change order for */ presedence
			if onode.OperatorPrecedence(operator) > 0 {
				if lastnode != nil && lastnode.Value.Category() == s.CatValue {
					this.curr = lastnode.Push(s.NewFuncToken(operator))
					this.state = this.parseValue
					return true
				}
			}
			//after */ presedence continue pushing +- operators from the bottom.
			if onode.OperatorPrecedence(operator) < 0 {
				for {
					v1, ok := this.curr.Parent().Value.(*s.FuncToken)
					if ok && strings.Index("+-", v1.Name) >= 0 {
						this.curr = this.curr.Parent()
					} else {
						break
					}
				}
			}
			//standard operator push
			this.curr = this.push(s.NewFuncToken(operator))
			this.state = this.parseValue
			return true
		}
		//push as operator argument
		if lastnode != nil {
			this.curr = lastnode.Push(s.NewFuncToken(operator))
			this.state = this.parseValue
		} else {
			this.error(fmt.Sprintf("Expecting a value before operator %q", operator))
			this.state = nil
		}
		return true
	}
	return false
}
