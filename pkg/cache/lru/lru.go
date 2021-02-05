package lru

import (
	"container/list"
	"sync"
)

type entry struct {
	key   interface{}
	value interface{}
}

type Cache struct {
	capacity int
	mu       *sync.Mutex
	queue    *list.List
	htable   map[interface{}]*list.Element
}

func NewCache(capacity int) *Cache {
	cache := Cache{
		capacity: capacity,
		mu:       &sync.Mutex{},
		queue:    list.New(),
		htable:   make(map[interface{}]*list.Element, capacity),
	}

	return &cache
}

func (c *Cache) Get(key interface{}) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if elem, ok := c.htable[key]; ok {
		c.queue.MoveToFront(elem)

		return elem.Value.(*entry).value, true
	}

	return nil, false
}

func (c *Cache) Set(key, value interface{}) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	if elem, ok := c.htable[key]; ok {
		c.queue.MoveToFront(elem)

		ent := elem.Value.(*entry)
		ent.value = value

		return true
	}

	if c.queue.Len() == c.capacity {
		c.purge()
	}

	ent := &entry{
		key:   key,
		value: value,
	}
	elem := c.queue.PushFront(ent)
	c.htable[key] = elem

	return true
}

func (c *Cache) purge() {
	if lastElem := c.queue.Back(); lastElem != nil {
		c.queue.Remove(lastElem)
		delete(c.htable, lastElem.Value.(*entry).key)
	}
}
