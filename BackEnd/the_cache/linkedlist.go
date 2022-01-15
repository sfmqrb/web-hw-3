package thecache

import "fmt"

type node struct {
	key   string
	value string
	prev  *node
	next  *node
}

type doublyLinkedList struct {
	limit  int
	tail *node
	head *node
}

func initDoublyList() *doublyLinkedList {
	return &doublyLinkedList{}
}

func (d *doublyLinkedList) __AddToFront__(key, value string) { // removable
	newNode := &node{
		key:   key,
		value: value,
	}
	if d.head == nil {
		d.head = newNode
		d.tail = newNode
	} else {
		newNode.prev = d.head
		d.head.next = newNode
		d.head = newNode
	}
	d.limit++
	//return
}

// changed RemoveFromFront to remove from tail
func (d *doublyLinkedList) RemoveFromEnd() {
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
func (d *doublyLinkedList) AddToFront(node *node) {
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
func (d *doublyLinkedList) MoveNodeToFront(node *node) {
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

func (d *doublyLinkedList) TraverseForward() error {
	if d.head == nil {
		return fmt.Errorf("TraverseError: List is empty")
	}
	temp := d.head
	for temp != nil {
		fmt.Printf("key = %v, value = %v, prev = %v, next = %v\n", temp.key, temp.value, temp.prev, temp.next)
		temp = temp.prev
	}
	fmt.Println()
	return nil
}

func (d *doublyLinkedList) Front() *node {
	return d.head
}

func (d *doublyLinkedList) End() *node {
	return d.tail
}

func (d *doublyLinkedList) Size() int {
	return d.limit
}
