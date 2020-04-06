package lru

type EvictCallback func(key string, value interface{})

type elem struct {
	Key   string
	Value interface{}
}
