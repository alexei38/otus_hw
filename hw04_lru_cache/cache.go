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

type lruCache struct {
	mutex    sync.Mutex
	capacity int
	queue    List
	items    map[Key]*ListItem
}

func (c *lruCache) Set(key Key, value interface{}) bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	old, ok := c.items[key]

	if ok {
		c.queue.Remove(old)
	} else if c.queue.Len() >= c.capacity && c.queue.Len() > 0 {
		removeItem := c.queue.Back()
		removeCacheItem, ok := removeItem.Value.(*cacheItem)
		if !ok {
			panic("cache is corrupted")
		}
		delete(c.items, removeCacheItem.key)
		c.queue.Remove(removeItem)
	}
	cItem := newCacheItem(key, value)
	item := c.queue.PushFront(cItem)
	c.items[key] = item
	return ok
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	item, ok := c.items[key]
	if ok {
		c.queue.MoveToFront(item)
		cItem, ok := item.Value.(*cacheItem)
		if !ok {
			panic("cache is corrupted")
		}
		return cItem.value, ok
	}
	return nil, false
}

func (c *lruCache) Clear() {
	c.items = make(map[Key]*ListItem, c.capacity)
	c.queue = NewList()
}

type cacheItem struct {
	key   Key
	value interface{}
}

func newCacheItem(key Key, value interface{}) *cacheItem {
	return &cacheItem{
		key:   key,
		value: value,
	}
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}
