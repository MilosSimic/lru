package lru

import (
	"container/list"
	"errors"
	"fmt"
	"sync"
)

type LRU struct {
	capacity  int
	evictlist *list.List
	cache     map[string]*list.Element
	onEvict   EvictCallback
	lock      *sync.Mutex
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
	lru.lock.Lock()
	defer lru.lock.Unlock()

	if val, ok := lru.cache[key]; ok {
		head := lru.evictlist.Front()
		lru.evictlist.MoveBefore(val, head)

		return val.Value, true
	}
	return nil, false
}

func (lru *LRU) Put(key string, value interface{}) (interface{}, bool) {
	lru.lock.Lock()
	defer lru.lock.Unlock()

	if val, ok := lru.cache[key]; ok {
		val.Value.(*Elem).Value = value

		head := lru.evictlist.Front()
		lru.evictlist.MoveBefore(val, head)

		return val.Value, false
	}

	if lru.evictlist.Len() == lru.capacity {
		tail := lru.evictlist.Back()
		lru.evictlist.Remove(tail)

		delete(lru.cache, tail.Value.(*Elem).Key)

		if lru.onEvict != nil {
			lru.onEvict(tail.Value.(*Elem).Key, tail.Value.(*Elem).Value)
		}
	}

	newElem := lru.evictlist.PushFront(&Elem{key, value})
	lru.cache[key] = newElem

	return newElem.Value, true
}

func (lru *LRU) Remove(key string) bool {
	lru.lock.Lock()
	defer lru.lock.Unlock()

	if e, ok := lru.cache[key]; ok {
		lru.evictlist.Remove(e)
		delete(lru.cache, key)

		if lru.onEvict != nil {
			lru.onEvict(e.Value.(*Elem).Key, e.Value.(*Elem).Value)
		}
		return true
	}
	return false
}

func (lru *LRU) Print() {
	lru.lock.Lock()
	defer lru.lock.Unlock()

	for e := lru.evictlist.Front(); e != nil; e = e.Next() {
		fmt.Print(e.Value.(*Elem).Value, " ")
	}
	fmt.Println("")
}

func (lru *LRU) Clear() {
	lru.lock.Lock()
	defer lru.lock.Unlock()

	for k, v := range lru.cache {
		if lru.onEvict != nil {
			lru.onEvict(k, v.Value.(*Elem).Value)
		}
		delete(lru.cache, k)
	}
	lru.evictlist.Init()
}

func (lru *LRU) Len() int {
	lru.lock.Lock()
	defer lru.lock.Unlock()

	return lru.evictlist.Len()
}

func (lru *LRU) All() []*Elem {
	lru.lock.Lock()
	defer lru.lock.Unlock()

	s := []*Elem{}
	for e := lru.evictlist.Front(); e != nil; e = e.Next() {
		s = append(s, e.Value.(*Elem))
	}
	return s
}

func (lru *LRU) Init(data []*Elem) bool {
	for _, e := range data {
		_, ok := lru.Put(e.Key, e.Value)
		if !ok {
			return false
		}
	}
	return true
}

func (lru *LRU) Contains(key string) interface{} {
	if val, ok := lru.cache[key]; ok {
		return val.Value.(*Elem).Value
	}
	return nil
}
