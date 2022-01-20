package main

import (
	"fmt"
	. "hw3/BackEnd/the_cache"
)

var cache TheCache

func main() {
	cache = InitCache()
	cache.SetKey(&Node{
		UserId:   0,
		UserName: "amm",
		Password: "amir1234",
		Name:     "amir",
		Notes:    nil,
		Prev:     nil,
		Next:     nil,
	})
	fmt.Println("aaa")
}
