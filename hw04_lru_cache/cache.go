package hw04lrucache

import "sync"

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
	rwMutex  sync.RWMutex
}

type cacheItem struct {
	key   Key
	value interface{}
}

func (c *lruCache) Set(key Key, value interface{}) bool {
	if cacheItem, ok := c.accessItem(key); ok {
		cacheItem.value = value
		return true
	}

	c.rwMutex.Lock()
	defer c.rwMutex.Unlock()

	newCacheItem := cacheItem{key, value}
	newListItem := c.queue.PushFront(&newCacheItem)
	c.items[key] = newListItem

	if c.queue.Len() > c.capacity {
		outdatedListItem := c.queue.Back()
		c.queue.Remove(outdatedListItem)
		if cacheItem, ok := outdatedListItem.Value.(*cacheItem); ok {
			delete(c.items, cacheItem.key)
		}
	}

	return false
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	if cacheItem, ok := c.accessItem(key); ok {
		return cacheItem.value, true
	}

	return nil, false
}

func (c *lruCache) accessItem(key Key) (*cacheItem, bool) {
	c.rwMutex.RLock()
	defer c.rwMutex.RUnlock()

	if listItem, ok := c.items[key]; ok {
		c.queue.MoveToFront(listItem)
		if cacheItem, ok := listItem.Value.(*cacheItem); ok {
			return cacheItem, true
		}
	}

	return nil, false
}

func (c *lruCache) Clear() {
	c.queue = NewList()
	c.items = make(map[Key]*ListItem, c.capacity)
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}
