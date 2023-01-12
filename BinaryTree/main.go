package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/fatih/color"
)

type BinaryTree struct {
	root *TreeNode
}

type TreeNode struct {
	left  *TreeNode
	right *TreeNode
	value int64
}

func (b *BinaryTree) insert(value int64) {
	if b.root == nil {
		b.root = &TreeNode{value: value}
	} else {
		b.root.insert(value)
	}
}

func (t *TreeNode) insert(value int64) {
	if value <= t.value { // Goes left
		if t.left == nil {
			t.left = &TreeNode{value: value}
		} else {
			t.left.insert(value)
		}
	} else {
		if t.right == nil {
			t.right = &TreeNode{value: value}
		} else {
			t.right.insert(value)
		}
	}
}

func (b *BinaryTree) render() {
	color.Cyan("[%d]", b.root.value)

	b.root.renderNodes(2)
}

func (b *TreeNode) renderNodes(offset int) {
	// Left
	if b.left != nil {
		for i := 0; i < offset; i++ {
			fmt.Print(" ")
		}

		color.Blue("[L %d]", b.left.value)
		b.left.renderNodes(offset + 2)
	}

	// Right
	if b.right != nil {
		for i := 0; i < offset; i++ {
			fmt.Print(" ")
		}
		color.Green("[R %d]", b.right.value)
		b.right.renderNodes(offset + 2)
	}
}

func main() {
	// Make rand great again
	rand.Seed(time.Now().UnixMilli())

	tree := &BinaryTree{}

	// Create values
	nodes := 20
	var values []int64

	for nodes > 0 {
		nodes--

		value := rand.Int63n(100)
		values = append(values, value)

		tree.insert(value)
	}

	fmt.Println("Values: ")
	tree.render()

}
