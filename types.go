package lru

type EvictCallback func(key string, value interface{})

type Elem struct {
	Key   string
	Value interface{}
}
