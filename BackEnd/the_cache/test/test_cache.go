package main

import (
	"fmt"
	. "hw3/BackEnd/the_cache"
)

var cache TheCache

func main() {
	fmt.Println(CACHE_CAPACITY)
	cache = InitCache()
	fmt.Println(CACHE_CAPACITY)
	cache.SetKey(&Node{
		UserId:   1,
		UserName: "amm",
		Password: "amir1234",
		Name:     "amir",
		Notes:    nil,
		Prev:     nil,
		Next:     nil,
	})
	fmt.Println("aaa")
	cache.SetKey(&Node{
		UserId:   2,
		UserName: "momo",
		Password: "momo12",
		Name:     "momo",
		Notes:    nil,
		Prev:     nil,
		Next:     nil,
	})
	fmt.Println("momo")
	cache.SetKey(&Node{
		UserId:   3,
		UserName: "Aomo",
		Password: "momo12",
		Name:     "Aomo",
		Notes:    nil,
		Prev:     nil,
		Next:     nil,
	})
	amir_note := Note{
		NoteId:    1,
		AuthorId:  1,
		Note:      "salam momo",
		NoteTitle: "Hello",
		NoteType:  "family",
	}
	thedata := CacheData{
		CommandType: 1,
		UserId:      1,
		Notes:       []Note{amir_note},
	}
	cache.SetExistingKey(&thedata)
	fmt.Println("first note added")
	cache.SetKey(&Node{
		UserId:   4,
		UserName: "Bomo",
		Password: "momo12",
		Name:     "Bomo",
		Notes:    nil,
		Prev:     nil,
		Next:     nil,
	})

	momo_note := Note{
		NoteId:    1,
		AuthorId:  2,
		Note:      "Hi amir",
		NoteTitle: "Greeting",
		NoteType:  "soldier",
	}
	thedata = CacheData{
		CommandType: 1,
		UserId:      2,
		Notes:       []Note{momo_note},
	}
	cache.SetExistingKey(&thedata)
	cache.SetKey(&Node{
		UserId:   5,
		UserName: "Como",
		Password: "momo12",
		Name:     "nini",
		Notes:    nil,
		Prev:     nil,
		Next:     nil,
	})
	amir_note2 := Note{
		NoteId:    2,
		AuthorId:  1,
		Note:      "dfgdfg",
		NoteTitle: "B.S.",
		NoteType:  "random",
	}
	thedata = CacheData{
		CommandType: 1,
		UserId:      1,
		Notes:       []Note{amir_note2},
	}
	cache.SetExistingKey(&thedata)
	cache.SetKey(&Node{
		UserId:   6,
		UserName: "Komo",
		Password: "momo12",
		Name:     "nini",
		Notes:    nil,
		Prev:     nil,
		Next:     nil,
	})
	PrintCache(cache)
	thedata = CacheData{
		CommandType: 2,
		UserId:      1,
		Notes:       []Note{amir_note},
	}
	cache.SetExistingKey(&thedata)
	PrintCache(cache)

}
