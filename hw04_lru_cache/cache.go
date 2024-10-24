package hw04lrucache

import "sync"

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	mu       sync.Mutex
	capacity int
	queue    List
	items    map[Key]*ListItem
}

type cacheValue struct {
	key   Key
	value interface{}
}

func (l *lruCache) Set(key Key, value interface{}) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	cacheVal := cacheValue{key, value}

	if item, ok := l.items[key]; ok {
		item.Value = cacheVal
		l.queue.MoveToFront(item)
		l.items[key] = item
		return true
	}

	if l.queue.Len() == l.capacity {
		back := l.queue.Back()
		l.queue.Remove(back)

		if backVal, ok := back.Value.(cacheValue); ok {
			delete(l.items, backVal.key)
		}
	}

	item := l.queue.PushFront(cacheVal)
	l.items[key] = item

	return false
}

func (l *lruCache) Get(key Key) (interface{}, bool) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if item, ok := l.items[key]; ok {
		l.queue.MoveToFront(item)
		l.items[key] = item

		if cacheVal, ok := item.Value.(cacheValue); ok {
			return cacheVal.value, true
		}
	}

	return nil, false
}

func (l *lruCache) Clear() {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.items = make(map[Key]*ListItem)
	l.queue = NewList()
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}
