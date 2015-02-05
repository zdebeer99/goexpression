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

func Parse(input string) *TreeNode {
	root := NewTreeNode(NewEmptyToken())
	parse := &parser{s.NewScanner(input), root, root, nil, nil}
	parse.parse()
	return root
}

func (this *parser) getValue() Token {
	if this.curr != nil {
		return this.curr.Value
	}
	return nil
}

func (this *parser) parse() {
	this.parseExpression()
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

func (l *parser) error(err interface{}) {
	var errortxt string
	if val, ok := err.(error); ok {
		errortxt = val.Error()
	} else {
		errortxt = err.(string)
	}
	lasttoken := l.scan.Commit()
	debug := fmt.Errorf("Line: %v, Near %q, Error: %s", l.scan.LineNumber(), lasttoken, errortxt)
	l.add(NewErrorToken(debug.Error()))
	l.err = debug
}

func (this *parser) parseExpression() bool {
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
		if this.parseOpenBracket() {
			hasExpression = true
			continue
		}
		if this.parseCloseBracket() {
			hasExpression = true
			continue
		}

		r := scan.Next()
		if s.IsSpace(r) {
			scan.Ignore()
			continue
		}

		this.error(fmt.Sprintf("Unexpected character %q", scan.Commit()))
		return hasExpression
	}
	return hasExpression
}

//parseOpenBracket
func (this *parser) parseOpenBracket() bool {
	scan := this.scan
	r := scan.Next()
	if r == '(' {
		this.curr = this.add(NewGroupToken("()"))
		scan.Commit()
		return true
	}
	scan.Backup()
	return false
}

//parseCloseBracket
func (this *parser) parseCloseBracket() bool {
	scan := this.scan
	r := scan.Next()
	if r == ')' {
		for {
			if this.curr.Parent() == nil {
				this.error("Brackets not closed.")
				return true
			}
			v1, ok := this.curr.Parent().Value.(*GroupToken)
			this.curr = this.curr.Parent()
			if ok && v1.GroupType == "()" {
				scan.Commit()
				this.curr = this.curr.Parent()
				return true
			}
		}
		panic("Should be impossible to reach this point.")
	}
	scan.Backup()
	return false
}

// parseValue
func (this *parser) parseValue() bool {
	scan := this.scan
	if scan.ScanNumber() {
		this.add(NewNumberToken(scan.Commit()))
		this.state = this.parseLRFunc
		return true
	}
	if scan.ScanWord() {
		this.add(NewIdentityToken(scan.Commit()))
		this.state = this.parseLRFunc
		return true
	}
	return false
}

//parseLRFunc
func (this *parser) parseLRFunc() bool {
	scan := this.scan
	if scan.Accept("+-/*") {
		operator := scan.Commit()
		lastnode := this.lastNode()
		onode, ok := this.getValue().(*FuncToken)
		//push excisting operator up in tree structure
		if ok {
			//operator is the same current operator ignore
			if onode.Name == operator {
				this.state = this.parseValue
				return true
			}
			//change order for */ presedence
			if onode.OperatorPrecedence(operator) > 0 {
				if lastnode != nil && lastnode.Value.Category() == CatValue {
					this.curr = lastnode.Push(NewFuncToken(operator))
					this.state = this.parseValue
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
			this.state = this.parseValue
			return true
		}
		//set previous found value as argument of the operator
		if lastnode != nil {
			this.curr = lastnode.Push(NewFuncToken(operator))
			this.state = this.parseValue
		} else {
			this.error(fmt.Sprintf("Expecting a value before operator %q", operator))
			this.state = nil
		}
		return true
	}
	return false
}
