package main

import (
	"fmt"
)

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
	if c.Index >= len(c.Parent.Children) {
		panic(fmt.Sprintf(
			"Index=%d, len(Parent.Children)=%d",
			c.Index,
			len(c.Parent.Children),
		))
	}
	return c.Parent.Children[c.Index]
}

func (c *Cursor) Set(x *Node) {
	c.Parent.Children[c.Index] = x
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
	fmt.Printf("------- (%v).AppendChild: %#v\n", t, x)
	if t.Cursor == nil {
		t.Root.Children = append(t.Root.Children, x)
		t.Cursor = &Cursor{
			Parent: t.Root,
			Index: len(t.Root.Children) - 1,
		}
		return
	}
	cur := t.Cursor.Get()
	cur.Children = append(cur.Children, x)
	t.Cursor.Index = len(cur.Children) - 1
}

func (t *Tree) InsertParent(p *Node) {
	if t.Cursor == nil {
		panic("Cursor == nil")
	}
	c := t.Cursor
	node := c.Get()
	if node == p {
		panic("node == p")
	}
	p.Children = []*Node{node}
	c.Set(p)
}

func NewTree() *Tree {
	return &Tree{
		Root: &Node{},
	}
}
