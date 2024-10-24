package hw04lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	len  int
	head *ListItem
	tail *ListItem
}

func (l *list) Len() int {
	return l.len
}

func (l *list) Front() *ListItem {
	return l.head
}

func (l *list) Back() *ListItem {
	return l.tail
}

func (l *list) PushFront(v interface{}) *ListItem {
	item := &ListItem{Value: v, Next: nil, Prev: nil}

	if l.len == 0 {
		l.head = item
		l.tail = item
	} else {
		l.head.Prev = item
		item.Next = l.head
		l.head = item
	}

	l.len++

	return item
}

func (l *list) PushBack(v interface{}) *ListItem {
	item := &ListItem{Value: v, Next: nil, Prev: nil}

	if l.len == 0 {
		l.head = item
		l.tail = item
	} else {
		l.tail.Next = item
		item.Prev = l.tail
		l.tail = item
	}

	l.len++

	return item
}

func (l *list) Remove(i *ListItem) {
	if l.head == l.tail {
		l.head = nil
		l.tail = nil
		l.len = 0
		return
	}

	switch {
	case l.head == i:
		i.Next.Prev = nil
		l.head = i.Next
	case l.tail == i:
		i.Prev.Next = nil
		l.tail = i.Prev
	default:
		prev := i.Prev
		next := i.Next
		prev.Next = next
		next.Prev = prev
	}

	l.len--
}

func (l *list) MoveToFront(i *ListItem) {
	if l.head == i || l.len < 2 {
		return
	}

	l.Remove(i)

	l.head.Prev = i
	i.Next = l.head
	l.head = i

	l.len++
}

func NewList() List {
	return new(list)
}
