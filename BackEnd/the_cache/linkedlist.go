package thecache

import (
	"github.com/uptrace/bun"
)

type Note struct {
	bun.BaseModel `bun:"table:notes,alias:u"`
	NoteId        int    `bun:"note_id,pk,autoincrement"`
	Note          string `bun:"note,notnull"`
	NoteTitle     string `bun:"title,notnull"`
	NoteType      string `bun:"type,notnull"`
	AuthorId      int    `bun:"author_id"`
}

//Notes could make error 'cause can't "append"

type Node struct {
	UserId   int
	UserName string
	Password string
	Name     string
	Notes    []Note
	Prev     *Node
	Next     *Node
}

type DoublyLinkedList struct {
	length int
	tail   *Node
	head   *Node
}

func initDoublyList(capacity int) *DoublyLinkedList {
	d := DoublyLinkedList{
		length: capacity}
	return &d
}

// changed RemoveFromFront to remove from tail
func (d *DoublyLinkedList) removeFromEnd() {
	if d.tail != nil {
		if d.head == d.tail {
			d.head = nil
			d.tail = nil
		} else {
			d.tail = d.tail.Next
			d.tail.Prev = nil
		}
		d.length--
	}
}

// Changed AddToEnd to AddToFront(*node)
func (d *DoublyLinkedList) addToFront(node *Node) {
	if d.head == nil {
		d.head = node
		d.tail = node
	} else {
		node.Prev = d.head
		d.head.Next = node
		d.head = node
	}
	d.length++
}

// Changed MoveNodeToEnd to MoveNodeToFront
func (d *DoublyLinkedList) moveNodeToFront(node *Node) {
	prev := node.Prev
	next := node.Next

	if next != nil {
		if prev != nil {
			next.Prev = prev
			prev.Next = next
		} else {
			d.tail = next
			next.Prev = nil
		}
		d.head.Next = node
		node.Prev = d.head
		node.Next = nil
		d.head = node
	}
}

func (d *DoublyLinkedList) front() *Node {
	return d.head
}

func (d *DoublyLinkedList) end() *Node {
	return d.tail
}

func (d *DoublyLinkedList) size() int {
	return d.length
}
