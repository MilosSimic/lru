package main

import (
	"container/list"
	"errors"
	"fmt"
)

type LRU struct {
	capacity  int
	evictlist *list.List
	cache     map[string]*list.Element
	onEvict   EvictCallback
}

func NewLRU(c int, f EvictCallback) (*LRU, error) {
	if c <= 0 {
		return nil, errors.New("Capacity must be provide a positive number")
	}

	return &LRU{
		capacity:  c,
		evictlist: list.New(),
		cache:     map[string]*list.Element{},
		onEvict:   f,
	}, nil
}

func (lru *LRU) Get(key string) (interface{}, bool) {
	if val, ok := lru.cache[key]; ok {
		head := lru.evictlist.Front()
		lru.evictlist.MoveBefore(val, head)

		return val.Value, true
	}
	return nil, false
}

func (lru *LRU) Put(key string, value interface{}) (interface{}, bool) {
	if val, ok := lru.cache[key]; ok {
		val.Value.(*elem).value = value

		head := lru.evictlist.Front()
		lru.evictlist.MoveBefore(val, head)

		return val.Value, false
	}

	if lru.evictlist.Len() == lru.capacity {
		tail := lru.evictlist.Back()
		lru.evictlist.Remove(tail)

		delete(lru.cache, tail.Value.(*elem).key)

		if lru.onEvict != nil {
			lru.onEvict(tail.Value.(*elem).key, tail.Value.(*elem).value)
		}
	}

	newElem := lru.evictlist.PushFront(&elem{key, value})
	lru.cache[key] = newElem

	return newElem.Value, true
}

func (lru *LRU) Print() {
	for e := lru.evictlist.Front(); e != nil; e = e.Next() {
		fmt.Print(e.Value.(*elem).value, " ")
	}
	fmt.Println("")
}

func (lru *LRU) Clear() {
	for k, v := range lru.cache {
		if lru.onEvict != nil {
			lru.onEvict(k, v.Value.(*elem).value)
		}
		delete(lru.cache, k)
	}
	lru.evictlist.Init()
}

func (lru *LRU) Len() int {
	return lru.evictlist.Len()
}

func (lru *LRU) Contains(key string) interface{} {
	if val, ok := lru.cache[key]; ok {
		return val.Value.(*elem).value
	}
	return nil
}
