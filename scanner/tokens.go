package scanner

import (
	"fmt"
	"strconv"
	"strings"
)

type Token interface {
	Category() TokenCategory
	SetError(err error)
	Error() error
	String() string
}

type TokenCategory int

const (
	CatOther TokenCategory = iota
	CatFunction
	CatValue
)

type EmptyToken struct {
	tokencat TokenCategory
	err      error
}

func NewEmptyToken() *EmptyToken {
	return &EmptyToken{CatOther, nil}
}

func (this *EmptyToken) Category() TokenCategory {
	return this.tokencat
}

func (this *EmptyToken) Error() error {
	return this.err
}

func (this *EmptyToken) SetError(err error) {
	this.err = err
}

func (this *EmptyToken) String() string {
	return "Base()"
}

type ErrorToken struct {
	EmptyToken
}

func NewErrorToken(err string) *ErrorToken {
	return &ErrorToken{EmptyToken{CatOther, fmt.Errorf(err)}}
}

type NumberToken struct {
	EmptyToken
	Value float64
}

func NewNumberToken(value string) *NumberToken {
	node := &NumberToken{EmptyToken{CatValue, nil}, 0}
	val1, err := strconv.ParseFloat(value, 64)
	if err != nil {
		panic("Number node failed to parse string to number")
		return node
	}
	node.Value = val1
	return node
}

func (this *NumberToken) String() string {
	return fmt.Sprintf("Number(%v)", this.Value)
}

type IdentityToken struct {
	EmptyToken
	Name string
}

func NewIdentityToken(name string) *IdentityToken {
	return &IdentityToken{EmptyToken{CatValue, nil}, name}
}

func (this *IdentityToken) String() string {
	return fmt.Sprintf("Identity(%s)", this.Name)
}

type FuncToken struct {
	EmptyToken
	Name string
}

func NewFuncToken(name string) *FuncToken {
	return &FuncToken{EmptyToken{CatFunction, nil}, name}
}

func (this *FuncToken) String() string {
	return fmt.Sprintf("Func(%s)", this.Name)
}

// OperatorPrecedence return true if the operator argument is lower than the current operator.
func (this *FuncToken) OperatorPrecedence(operator string) int {
	if strings.Contains("*/", operator) && strings.Contains("+-", this.Name) {
		return 1
	}
	if strings.Contains("+-", operator) && strings.Contains("*/", this.Name) {
		return -1
	}
	return 0
}

type GroupToken struct {
	EmptyToken
	GroupType string
}

func NewGroupToken(group string) *GroupToken {
	return &GroupToken{EmptyToken{CatOther, nil}, group}
}

func (this *GroupToken) String() string {
	return fmt.Sprintf("Group(%s)", this.GroupType)
}
