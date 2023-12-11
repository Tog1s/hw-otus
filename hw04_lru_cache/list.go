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

func NewList() List {
	return new(list)
}

func (list *list) Len() int {
	return list.len
}

func (list *list) Front() *ListItem {
	return list.head
}

func (list *list) Back() *ListItem {
	return list.tail
}

func (list *list) Remove(item *ListItem) {
	if item == nil {
		return
	}

	if item.Prev != nil {
		item.Prev.Next = item.Next
	} else {
		list.head = item.Next
	}

	if item.Next != nil {
		item.Next.Prev = item.Prev
	} else {
		list.tail = item.Prev
	}
	list.len--
}

func (list *list) MoveToFront(item *ListItem) {
	if list.head == item {
		return
	}
	list.Remove(item)
	list.PushFront(item.Value)
}

func (list *list) PushFront(value interface{}) *ListItem {
	newListItem := &ListItem{
		Value: value,
	}

	if list.head == nil {
		list.head = newListItem
		list.tail = newListItem
	} else {
		newListItem.Next = list.head
		list.head.Prev = newListItem
		list.head = newListItem
	}
	list.len++
	return list.head
}

func (list *list) PushBack(value interface{}) *ListItem {
	newListItem := &ListItem{
		Value: value,
	}
	if list.head == nil {
		list.head = newListItem
		list.tail = newListItem
	} else {
		newListItem.Prev = list.tail
		list.tail.Next = newListItem
		list.tail = newListItem
	}
	list.len++
	return list.tail
}
