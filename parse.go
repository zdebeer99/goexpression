package goexpression

import (
	"fmt"
	"github.com/zdebeer99/goexpression/scanner"
)

type stateFn func(*parser) stateFn

type parser struct {
	scan  *scanner.Scanner
	root  *TreeNode
	curr  *TreeNode
	err   error
	state stateFn
}

func Parse(input string) (*TreeNode, error) {
	root := NewTreeNode(NewEmptyToken())
	parse := &parser{scanner.NewScanner(input), root, root, nil, nil}
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

func (this *parser) error(err interface{}, a ...interface{}) {
	var errortxt string
	if val, ok := err.(error); ok {
		errortxt = val.Error()
	} else {
		if len(a) > 0 {
			errortxt = fmt.Sprintf(err.(string), a)
		} else {
			errortxt = err.(string)
		}
	}
	lasttoken := this.commit()
	if len(lasttoken) < 10 {
		for i := len(lasttoken); i < 10 && !this.scan.IsEOF(); i++ {
			this.scan.Next()
		}
		lasttoken = lasttoken + this.commit()
	}
	debug := fmt.Errorf("Line: %v, near %q, Error: %s", this.scan.LineNumber(), lasttoken, errortxt)
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
func (this *parser) parseCloseBracket() stateFn {
	for {
		v1, ok := this.curr.Value.(*GroupToken)
		if ok && v1.GroupType == "()" {
			this.commit()
			this.curr = this.curr.Parent()
			return branchExpressionOperatorPart
		}
		if ok && v1.GroupType == "" {
			//must be a bracket part of a parent loop, exit this sub loop.
			this.scan.Backup()
			return nil
		}
		if this.curr.Parent() == nil {
			this.error("Brackets not closed.")
			return nil
		}
		this.curr = this.curr.Parent()
	}
	panic("Should be impossible to reach this point.")
}

func (this *parser) AcceptOperator() bool {	
	scan:=this.scan
	for _,op:= range operatorList{
		if scan.Prefix(op) {
			return true
		}
	}
	return false
}

//parseOperator
func (this *parser) parseOperator() bool {
	operator := this.commit()
	lastnode := this.lastNode()
	onode, ok := this.getCurr().(*OperatorToken)
	//push excisting operator up in tree structure
	if ok {
		//operator is the same current operator ignore
		if onode.Operator == operator {
			return true
		}
		//change order for */ presedence
		//fmt.Println(onode, operator, onode.Precedence(operator))		
		if onode.Precedence(operator) > 0 {
			if lastnode != nil {
				this.curr = lastnode.Push(NewOperatorToken(operator))
				return true
			}
		}
		//after */ presedence fallback and continue pushing +- operators from the bottom.
		if onode.Precedence(operator) < 0 {
			for {
				v1, ok := this.curr.Parent().Value.(*OperatorToken)
				//if ok && strings.Index("+-", v1.Name) >= 0 {
				if ok && operators.Level(v1.Operator)>=0{
					this.curr = this.curr.Parent()
				} else {
					break
				}
			}
		}
		//standard operator push
		this.curr = this.push(NewOperatorToken(operator))
		return true
	}
	//set previous found value as argument of the operator
	if lastnode != nil {
		this.curr = lastnode.Push(NewOperatorToken(operator))
	} else {
		this.error(fmt.Sprintf("Expecting a value before operator %q", operator))
		this.state = nil
	}
	return true
}

//parseLRFunc
func (this *parser) parseLRFunc() bool {
	lrfunc := this.commit()
	lastnode := this.lastNode()
	if lastnode != nil {
		this.curr = lastnode.Push(NewLRFuncToken(lrfunc))
	} else {
		this.error(fmt.Sprintf("Expecting a value before operator %q", lrfunc))
		this.state = nil
	}
	return false
}

func (this *parser) ParseText() string {
	scan := this.scan
	r := scan.Next()
	if r == '"' || r == '\'' {
		scan.Ignore()
		endqoute := r
		for {
			r = scan.Next()
			if r == endqoute {
				scan.Backup()
				txt := scan.Commit()
				scan.Next()
				scan.Ignore()
				return txt
			}
			if scan.IsEOF() {
				this.error("Missing Qoute and end of text.")
				return "Error"
			}
		}
	}
	return ""
}
