package hw04lrucache

import (
	"sync"
)

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type itemValue struct {
	key   Key
	value interface{}
}

type lruCache struct {
	capacity int
	mx       sync.Mutex
	queue    List
	items    map[Key]*ListItem
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

func (list *lruCache) Set(key Key, value interface{}) bool {
	list.mx.Lock()
	defer list.mx.Unlock()
	if v, ok := list.items[key]; ok {
		v.Value = &itemValue{key: key, value: value}
		list.queue.MoveToFront(v)
		return true
	}
	v := list.queue.PushFront(&itemValue{key: key, value: value})
	list.items[key] = v

	if list.queue.Len() > list.capacity {
		back := list.queue.Back()
		list.queue.Remove(back)
		delete(list.items, back.Value.(*itemValue).key)
	}
	return false
}

func (list *lruCache) Get(key Key) (interface{}, bool) {
	list.mx.Lock()
	defer list.mx.Unlock()
	if v, ok := list.items[key]; ok {
		list.queue.MoveToFront(v)
		return v.Value.(*itemValue).value, true
	}
	return nil, false
}

func (list *lruCache) Clear() {
	list.mx.Lock()
	defer list.mx.Unlock()
	list.queue = NewList()
	list.items = make(map[Key]*ListItem, list.capacity)
}
