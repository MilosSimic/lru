package lru

type Cache interface {
	Get(key string) (interface{}, bool)
	Put(key string, value interface{}) (interface{}, bool)
	Clear()
	Len() int
	Contins(key string) interface{}
}
