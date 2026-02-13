package main

import "fmt"

// 1
type Doubleable interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~float32 | ~float64
}

func Double[T Doubleable](in T) T {
	return in * 2
}

type Printable interface {
	~int | ~float64
	fmt.Stringer
}

// 2
type printableInt int

func (p printableInt) String() string {
	return fmt.Sprintf("%d", p)
}

type printableFloat float64

func (p printableFloat) String() string {
	return fmt.Sprintf("%f", p)
}

func rePrint[T Printable](in T) {
	fmt.Println(in)
}

// 3 - Add(T); Insert(T, int); Index(T) int
type constraintedNode[T comparable] struct {
	val      T
	nextNode *constraintedNode[T]
}

type linkedList[T comparable] struct {
	head *constraintedNode[T]
	tail *constraintedNode[T]
}

func createLinkedList[T comparable](in T) *linkedList[T] {
	firstNode := constraintedNode[T]{val: in, nextNode: nil}
	return &linkedList[T]{
		head: &firstNode,
		tail: &firstNode,
	}
}

func (l *linkedList[T]) Add(in T) {
	if l == nil {
		return
	}

	newNode := constraintedNode[T]{
		val:      in,
		nextNode: nil,
	}

	l.tail.nextNode = &newNode
	l.tail = &newNode
}

func (l *linkedList[T]) Insert(in T, pos int) {
	if l == nil {
		return
	}

	if pos == 0 {
		l.head.val = in
		return
	}

	curNode := l.head.nextNode
	for i := 1; i != pos; i++ {
		if curNode == l.tail {
			return
		}

		curNode = curNode.nextNode
	}

	curNode.val = in
}

func (l linkedList[T]) Index(in T) int {
	var index int
	curNode := l.head
	for {
		if curNode.val == in {
			break
		}

		if curNode == l.tail {
			return -1
		}

		curNode = curNode.nextNode
		index++
	}
	return index
}

func main() { // 1
	someVar := 2.0
	fmt.Println(Double(someVar))
	anotherVar := 4
	fmt.Println(Double(anotherVar))
	var yetAnotherVar uint8 = 8
	fmt.Println(Double(yetAnotherVar))
	// 2
	var pInt printableInt = 3
	rePrint(pInt)
	var pFloat printableFloat = 4.0
	rePrint(pFloat)
	// 3
	newLL := createLinkedList(2)
	fmt.Println(newLL, newLL.head.val, newLL.head.nextNode)
	newLL.Add(3)
	fmt.Println(newLL, newLL.head.val, newLL.head.nextNode, newLL.tail.val, newLL.tail.nextNode)
	newLL.Insert(1, 0)
	newLL.Insert(5, 1)
	fmt.Println(newLL, newLL.head.val, newLL.head.nextNode, newLL.tail.val, newLL.tail.nextNode)
	newLL.Add(25)
	fmt.Println(newLL, newLL.head.val, newLL.head.nextNode, newLL.tail, newLL.tail.val, newLL.tail.nextNode)
	newLL.Add(50)
	newLL.Add(100)
	newLL.Add(250)
	fmt.Println(newLL.Index(100))
	fmt.Println(newLL.Index(1))
	fmt.Println(newLL.Index(250))
	fmt.Println(newLL.Index(1000))
}
