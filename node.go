package goexpression

import (
	"fmt"
	"strconv"
	"strings"
)

//go:generate stringer -type=NodeType
type NodeType int

const (
	NodeEmpty NodeType = iota
	NodeNumber
	NodeIdentifier
	NodeFunc
	NodeError
)

//go:generate stringer -type=NodeCategorie
type NodeCategorie int

const (
	CatValue NodeCategorie = iota
	CatFunction
	CatOther
)

type Node interface {
	NodeType() NodeType
	NodeCat() NodeCategorie
	setParent(Node)
	Parent() Node
	Items() []Node
	LastItem() Node
	Add(Node) Node
	Push(Node) Node
	Stack(Node) Node
	String() string
	Error() error
}

type NumberNode struct {
	BaseNode
	value float64
}

type IdentityNode struct {
	BaseNode
	name string
}

type FuncNode struct {
	BaseNode
	name string
}

type ErrorNode struct {
	BaseNode
}

func NewNumberNode(number string) *NumberNode {
	node := &NumberNode{BaseNode: *NewBaseNode(NodeNumber, CatValue)}
	node.me = node
	val1, err := strconv.ParseFloat(number, 64)
	if err != nil {
		node.err = err
		return node
	}
	node.value = val1
	return node
}

func (this *NumberNode) String() string {
	if this.StringContent() == "" {
		return fmt.Sprintf("Number(%v)", this.value)
	}

	return fmt.Sprintf("[Number(%v):%s]", this.value, this.StringContent())
}

func NewIdentityNode(identity string) *IdentityNode {
	node := &IdentityNode{*NewBaseNode(NodeIdentifier, CatValue), identity}
	node.me = node
	return node
}

func (this *IdentityNode) String() string {
	if this.StringContent() == "" {
		return fmt.Sprintf("Identity(%v)", this.name)
	}

	return fmt.Sprintf("[Identity(%v):%s]", this.name, this.StringContent())
}

func NewFuncNode(name string) *FuncNode {
	node := &FuncNode{*NewBaseNode(NodeFunc, CatFunction), name}
	node.me = node
	return node
}

// OperatorPrecedence return true if the operator argument is lower than the current operator.
func (this *FuncNode) OperatorPrecedence(operator string) int {
	if strings.Contains("*/", operator) && strings.Contains("+-", this.name) {
		return 1
	}
	if strings.Contains("+-", operator) && strings.Contains("*/", this.name) {
		return -1
	}
	return 0
}

func (this *FuncNode) String() string {
	if this.StringContent() == "" {
		return fmt.Sprintf("[Func(%v)]", this.name)
	}
	return fmt.Sprintf("[Func(%v):%s]", this.name, this.StringContent())
}

func NewErrorNode(err string) *ErrorNode {
	node := &ErrorNode{*NewBaseNode(NodeError, CatOther)}
	node.err = fmt.Errorf(err)
	node.me = node
	return node
}

//BaseNode
type BaseNode struct {
	me       Node
	nodeType NodeType
	nodeCat  NodeCategorie
	parent   Node
	items    []Node
	err      error
}

func NewBaseNode(nodetype NodeType, nodeCat NodeCategorie) *BaseNode {
	r1 := &BaseNode{nodeType: nodetype, nodeCat: nodeCat}
	r1.me = r1
	return r1
}

// NodeType returns the current node type.
func (b *BaseNode) NodeType() NodeType {
	return b.nodeType
}

func (b *BaseNode) NodeCat() NodeCategorie {
	return b.nodeCat
}

// Items Returns the nodes items
func (b *BaseNode) Items() []Node {
	return b.items
}

// LastItem returns the last item in the items array of the current node.
func (b *BaseNode) LastItem() Node {
	if len(b.items) > 0 {
		return b.items[len(b.items)-1]
	} else {
		return nil
	}
}

// Add a Node to the current Node and returns the node that was added.
// Ex: Take a Node A
//   [A]
//
//   x = A.Add(B)
//   [A ([B])]
//
//   x.Add(C)
//   [A ([B],[C])]
func (b *BaseNode) Add(item Node) Node {
	item.setParent(b.me)
	b.items = append(b.items, item)
	return item
}

// Add a node to the last node in the items and return the last node the node was added to.
// [A ([B],[C])]
//
// A.Stack(D) returns C
// [A ([B],[C (D)]]
//
// D is now a child node of C
func (b *BaseNode) Stack(item Node) Node {
	last := b.LastItem()
	if last == nil {
		return b.Add(item)
	}
	last.Add(item)
	return last
}

// Push replaces the current node with the item and adds the current node as a child. returns the node that was added.
// [A ([B],[C])]
// C.Push(D)
// [A ([B],[D ([C])])]
// C is now a child node of D, and D is positioned in the place of C
func (b *BaseNode) Push(item Node) Node {
	parent := b.Parent()
	if parent != nil {
		for i, v := range parent.Items() {
			if b.me == v {
				parent.Items()[i] = item
				item.setParent(parent)
				item.Add(b.me)
				return item
			}
		}
	}
	fmt.Println("PUSH ERROR", b.parent, b.me, item)
	panic("Could not find current node as child in parent node. Weird? should not be possible unless you changed stuff with the parents and children without updating everything else.")
}

// Parent return the parent node.
func (b *BaseNode) Parent() Node {
	return b.parent
}

// setParent set the current node's parent. for internal use only
func (b *BaseNode) setParent(item Node) {
	b.parent = item
}

func (b *BaseNode) StringContent() string {
	lines := make([]string, len(b.items))
	for k, v := range b.items {
		lines[k] = v.String()
	}
	if b.err != nil {
		return fmt.Sprintf("[ERROR: %s]", b.err.Error())
	} else if len(lines) > 0 {
		return fmt.Sprintf("%s", strings.Join(lines, ","))
	} else {
		return ""
	}
}

func (this *BaseNode) String() string {
	if this.StringContent() == "" {
		return "[Base()]"
	}
	return fmt.Sprintf("[Base():%s]", this.StringContent())
}

func (this *BaseNode) Error() error {
	return this.err
}
