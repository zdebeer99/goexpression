package goexpression

import (
	"fmt"
	"strings"
)

//TODO: Implement Stack to handle drill into functions.

type parser struct {
	scan  *Scanner
	root  Node
	curr  Node
	err   error
	state func() bool
}

func Parse(input string) Node {
	root := NewBaseNode(NodeEmpty, CatOther)
	parse := &parser{NewScanner(input), root, root, nil, nil}
	parse.parse()
	return root
}

func (this *parser) parse() {
	this.parseExpresion()
}

func (this *parser) add(node Node) Node {
	return this.curr.Add(node)
}

func (this *parser) stack(node Node) Node {
	return this.curr.Stack(node)
}

func (this *parser) push(node Node) Node {
	return this.curr.Push(node)
}

func (this *parser) lastNode() Node {
	return this.curr.LastItem()
}

func (this *parser) parentNode() Node {
	return this.curr.Parent()
}

func (l *parser) error(err interface{}) {
	var text string
	if val, ok := err.(error); ok {
		text = val.Error()
	} else {
		text = err.(string)
	}
	debug := fmt.Errorf("Line: %v, Error: %s", l.scan.lineNumber(), text)
	l.add(NewErrorNode(debug.Error()))
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
		if IsSpace(r) {
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
		this.add(NewNumberNode(scan.Commit()))
		this.state = this.parseLRFunc
		return true
	}
	if scan.ScanWord() {
		this.add(NewIdentityNode(scan.Commit()))
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
		onode, ok := this.curr.(*FuncNode)
		//push excisting operator
		if ok {
			//operator is the same as the previous one.
			if onode.name == operator {
				this.state = this.parseValue
				return true
			}
			//change order for */ presedence
			if onode.OperatorPrecedence(operator) > 0 {
				if lastnode != nil && lastnode.NodeCat() == CatValue {
					this.curr = lastnode.Push(NewFuncNode(operator))
					this.state = this.parseValue
					return true
				}
			}
			//after */ presedence continue pushing +- operators from the bottom.
			if onode.OperatorPrecedence(operator) < 0 {
				for {
					v1, ok := this.curr.Parent().(*FuncNode)
					if ok && strings.Index("+-", v1.name) >= 0 {
						this.curr = v1
					} else {
						break
					}
				}
			}
			//standard operator push
			this.curr = this.push(NewFuncNode(operator))
			this.state = this.parseValue
			return true
		}
		//push as operator argument
		if lastnode != nil {
			this.curr = lastnode.Push(NewFuncNode(operator))
			this.state = this.parseValue
		} else {
			this.error(fmt.Sprintf("Expecting a value before operator %q", operator))
			this.state = nil
		}
		return true
	}
	return false
}
