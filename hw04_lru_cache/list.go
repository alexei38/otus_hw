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
	lastNode  *ListItem
	firstNode *ListItem
	len       int
}

func (l *list) Len() int {
	return l.len
}

func (l *list) Front() *ListItem {
	return l.firstNode
}

func (l *list) Back() *ListItem {
	return l.lastNode
}

func (l *list) isRemoved(i *ListItem) bool {
	return i.Next == nil && i.Prev == nil && l.Len() > 1
}

func (l *list) PushFront(v interface{}) *ListItem {
	var item *ListItem
	switch res := v.(type) {
	case *ListItem:
		item = res
	default:
		item = &ListItem{Value: v}
	}
	if l.firstNode != nil {
		item.Next = l.firstNode
		item.Next.Prev = item
	}
	l.firstNode = item
	if l.lastNode == nil {
		l.lastNode = item
	}
	l.len++
	return item
}

func (l *list) PushBack(v interface{}) *ListItem {
	var item *ListItem
	switch res := v.(type) {
	case *ListItem:
		item = res
	default:
		item = &ListItem{Value: v}
	}

	if l.firstNode == nil {
		l.firstNode = item
	}
	if l.lastNode != nil {
		l.lastNode.Next = item
		item.Prev = l.lastNode
	}
	l.lastNode = item
	l.len++
	return item
}

func (l *list) Remove(i *ListItem) {
	if i == nil || l.isRemoved(i) {
		return
	}
	if i == l.lastNode && i.Prev == nil {
		l.lastNode = nil
	}
	if i == l.firstNode && i.Next == nil {
		l.firstNode = nil
	}
	if i.Prev != nil {
		i.Prev.Next = i.Next
		if i == l.lastNode {
			l.lastNode = i.Prev
		}
	}
	if i.Next != nil {
		i.Next.Prev = i.Prev
		if i == l.firstNode {
			l.firstNode = i.Next
		}
	}
	i.Prev = nil
	i.Next = nil
	l.len--
}

func (l *list) MoveToFront(i *ListItem) {
	if i != l.firstNode && !l.isRemoved(i) {
		l.Remove(i)
		l.PushFront(i)
	}
}

func NewList() List {
	var l List = &list{}
	return l
}
