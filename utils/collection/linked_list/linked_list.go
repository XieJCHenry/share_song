package linked_list

var empty any

type listNode[T any] struct {
	Val  T
	Next *listNode[T]
}

func newNodeWithVal[T any](v T) *listNode[T] {
	return &listNode[T]{
		Val:  v,
		Next: nil,
	}
}

func newNode[T any]() *listNode[T] {
	return &listNode[T]{
		Next: nil,
	}
}

type List[T any] struct {
	head *listNode[T] // first = head->next
	tail *listNode[T] // end = tail
	size int
}

func NewList[T any]() *List[T] {
	head := newNode[T]()

	return &List[T]{
		head: head,
		tail: nil,
		size: 0,
	}
}

func NewListWith[T any](elems ...T) *List[T] {
	lst := NewList[T]()
	for _, elem := range elems {
		lst.Append(elem)
	}
	return lst
}

func (l *List[T]) Append(v T) {
	node := newNodeWithVal[T](v)
	if l.head.Next == nil {
		l.tail = node
		l.head.Next = l.tail
	} else {
		l.tail.Next = node
		l.tail = node
	}
	l.size++
}

func (l *List[T]) Insert(index int, v T) {
	if index >= l.size {
		l.Append(v)
	} else if index <= 0 {
		node := newNodeWithVal(v)
		node.Next = l.head.Next
		l.head.Next = node
	} else {
		var pos int = 0
		var prev, next *listNode[T]
		next = l.head.Next
		for ; pos < index && next != nil && next.Next != nil; pos++ {
			prev = next
			next = next.Next
		}
		node := newNodeWithVal(v)
		node.Next = next
		prev.Next = node
	}
}

func (l *List[T]) First() (v T) {
	if l.head.Next != nil {
		v = l.head.Next.Val
	}
	return
}

func (l *List[T]) Last() (v T) {
	if l.tail != nil {
		v = l.tail.Val
	}
	return
}

func (l *List[T]) Size() int {
	return l.size
}
