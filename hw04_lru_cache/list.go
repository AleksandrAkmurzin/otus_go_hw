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
	front *ListItem
	back  *ListItem
	len   int
}

func (l list) Len() int {
	return l.len
}

func (l list) Front() *ListItem {
	return l.front
}

func (l list) Back() *ListItem {
	return l.back
}

func (l *list) PushFront(v interface{}) *ListItem {
	oldFront := l.Front()
	newFront := &ListItem{v, oldFront, nil}

	if oldFront != nil {
		oldFront.Prev = newFront
	}

	l.front = newFront
	if l.Len() == 0 {
		l.back = newFront
	}
	l.len++

	return newFront
}

func (l *list) PushBack(v interface{}) *ListItem {
	oldBack := l.Back()
	newBack := &ListItem{v, nil, oldBack}

	if oldBack != nil {
		oldBack.Next = newBack
	}

	l.back = newBack
	if l.Len() == 0 {
		l.front = newBack
	}
	l.len++

	return newBack
}

func (l *list) Remove(i *ListItem) {
	prev := i.Prev
	next := i.Next

	// Update item links.
	if prev != nil {
		prev.Next = next
	}
	if next != nil {
		next.Prev = prev
	}

	// Process corner cases.
	if i == l.Front() {
		l.front = next
	}
	if i == l.Back() {
		l.back = prev
	}

	l.len--
}

func (l *list) MoveToFront(i *ListItem) {
	l.Remove(i)
	l.PushFront(i.Value)
}

func NewList() List {
	return new(list)
}
