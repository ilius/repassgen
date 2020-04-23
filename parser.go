package main

type NodeType uint8

const (
	ROOT NodeType = iota + 1
	STATIC
	CHARCLASS
	REPEAT
	GROUP
	FUNC_CALL
)

type Node struct {
	Type     NodeType
	Args     []interface{}
	Gen      generatorIface
	Children []*Node
}

type Cursor struct {
	Parent *Node
	Index  int
}

func (c *Cursor) Get() *Node {
	return c.Parent.Children[c.Index]
}

func (c *Cursor) Set(x *Node) {
	c.Parent.Children[c.Index] = x
}

func (c *Cursor) InsertParent(p *Node) {
	p.Children = []*Node{c.Get()}
	c.Set(p)
}

type Tree struct {
	Root   *Node
	Cursor *Cursor
}
