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

func (t *Tree) GetCursor() *Node {
	if t.Cursor == nil {
		return t.Root
	}
	return t.Cursor.Get()
}

func (t *Tree) AppendChild(x *Node) {
	if t.Cursor == nil {
		t.Root.Children = append(t.Root.Children, x)
		t.Cursor = &Cursor{Parent: t.Root, Index: len(t.Root.Children) - 1}
		return
	}
	cur := t.Cursor.Get()
	cur.Children = append(cur.Children, x)
	t.Cursor.Index = len(cur.Children) - 1
}

func NewTree() *Tree {
	return &Tree{
		Root: &Node{},
	}
}
