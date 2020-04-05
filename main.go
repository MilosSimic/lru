package main

import "fmt"

func main() {
	lru, _ := NewLRU(3, func(key string, value interface{}) {
		fmt.Println("Evicted")
		fmt.Println("Key ", key)
		fmt.Println("Value", value)
	})

	lru.Put("a", 1)
	lru.Print()

	lru.Put("b", 2)
	lru.Print()

	lru.Put("c", 3)
	lru.Print()

	lru.Put("d", 4)
	lru.Print()

	lru.Get("b")
	lru.Print()

	lru.Clear()
}
