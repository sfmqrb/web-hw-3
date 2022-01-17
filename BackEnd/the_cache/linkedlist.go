package thecache

import "fmt"

type Note struct {
	//bun.BaseModel `bun:"table:notes,alias:u"`
	NoteId    int    `bun:"note_id,pk,autoincrement"`
	Note      string `bun:"note,notnull"`
	NoteTitle string `bun:"title,notnull"`
	AuthorId  int    `bun:"author_id"`
}

// notes could make error 'cause can't "append" 
type Node struct {
	user_id  int
	username string
	password string
	name     string
	notes    []Note
	prev     *Node
	next     *Node
}

type DoublyLinkedList struct {
	limit int
	tail  *Node
	head  *Node
}

func initDoublyList(capacity int) *DoublyLinkedList {
	d := DoublyLinkedList{
		limit: capacity}
	return &d
}

// changed RemoveFromFront to remove from tail
func (d *DoublyLinkedList) removeFromEnd() {
	if d.tail != nil {
		if d.head == d.tail {
			d.head = nil
			d.tail = nil
		} else {
			d.tail = d.tail.prev
			d.tail.prev = nil
		}
		d.limit--
	}
}

// Changed AddToEnd to AddToFront(*node)
func (d *DoublyLinkedList) addToFront(node *Node) {
	if d.head == nil {
		d.head = node
		d.tail = node
	} else {
		node.prev = d.head
		d.head.next = node
		d.head = node
	}
	d.limit++
}

// Changed MoveNodeToEnd to MoveNodeToFront
func (d *DoublyLinkedList) moveNodeToFront(node *Node) {
	prev := node.prev
	next := node.next

	if next != nil {
		if prev != nil {
			next.prev = prev
			prev.next = next
		} else {
			d.tail = next
			next.prev = nil
		}
		d.head.next = node
		node.prev = d.head
		node.next = nil
		d.head = node
	}
}

func (d *DoublyLinkedList) traverseForward() error {
	if d.head == nil {
		return fmt.Errorf("TraverseError: List is empty")
	}
	temp := d.head
	for temp != nil {
		fmt.Printf("user_id = %v, name = %v, prev = %v, next = %v\n", temp.user_id, temp.name, temp.prev, temp.next)
		temp = temp.prev
	}
	fmt.Println()
	return nil
}

func (d *DoublyLinkedList) front() *Node {
	return d.head
}

func (d *DoublyLinkedList) end() *Node {
	return d.tail
}

func (d *DoublyLinkedList) size() int {
	return d.limit
}
