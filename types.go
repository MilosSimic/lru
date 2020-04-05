package main

type EvictCallback func(key string, value interface{})

type elem struct {
	key   string
	value interface{}
}
