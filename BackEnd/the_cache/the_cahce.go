package thecache

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

//import "fmt"

// configurables
//todo config file
var CACHE_CAPACITY int

const (

	// action types of cache Note request
	Save   = 1
	Del    = 2
	Get    = 3
	GetAll = 5
	Edit   = 4

	// action types of cache login request
	Login  = 1
	SignUp = 2

	// data types of cache_data
)

type CacheData struct {
	CommandType int
	UserId      int
	UserName    string
	Password    string
	Name        string
	Notes       []Note
}

/*  keys and values are int32 or string
keys -> 64 char / values -> 2048 char
commands: getkey, setkey, Clear
*/

type TheCache struct {
	dll     *DoublyLinkedList
	storage map[int]*Node
}

//////////////	CHANGE PATH	\\\\\\\\\\\\\\\\\\\\\\\\\\

func InitCache() TheCache {
	storage := make(map[int]*Node)
	configFile, err := ioutil.ReadFile("the_cache/cache_config.json")
	var configuration map[string]int
	json.Unmarshal([]byte(configFile), &configuration)
	//json.Unmarshal(configFile, &configuration)
	CACHE_CAPACITY = configuration["MAX_CAPACITY"]
	if err != nil {
		CACHE_CAPACITY = 16
	}
	fmt.Print("max capacity ")
	fmt.Print(CACHE_CAPACITY)
	return TheCache{
		dll:     &DoublyLinkedList{},
		storage: storage,
	}
}

func (cache *TheCache) Clear() {
	cache.dll = initDoublyList(CACHE_CAPACITY)
	cache.storage = make(map[int]*Node)
}

func (cache *TheCache) SetKey(node *Node) bool {
	_, ok := cache.storage[node.UserId]
	if ok {
		cache.dll.moveNodeToFront(node)
		return true
	}
	if cache.dll.size() >= CACHE_CAPACITY {
		delete(cache.storage, cache.dll.tail.UserId)
		cache.dll.removeFromEnd()
	}
	cache.storage[node.UserId] = node
	cache.dll.addToFront(node)
	return len(cache.storage) == cache.dll.size()
}
func (cache *TheCache) SetExistingKey(data *CacheData) bool {
	node, ok := cache.storage[data.UserId]
	if ok {
		switch data.CommandType {
		case Save:
			if data.UserName != "" {
				node.UserName = data.UserName
			}
			if data.Password != "" {
				node.Password = data.Password
			}
			if data.Name != "" {
				node.Name = data.Name
			}
			if len(data.Notes) > 0 {
				node.Notes = append(node.Notes, data.Notes...)
			}
		case Del:
			if len(data.Notes) > 0 {
				for index, note := range node.Notes {
					if note.NoteId == data.Notes[0].NoteId {
						node.Notes[index] = node.Notes[len(node.Notes)-1]
						node.Notes[len(node.Notes)-1] = Note{}
						node.Notes = node.Notes[:len(node.Notes)-1]
						break
					}
				}
			}
		case Edit:
			if data.UserName != "" {
				node.UserName = data.UserName
			}
			if data.Password != "" {
				node.Password = data.Password
			}
			if data.Name != "" {
				node.Name = data.Name
			}
			if len(data.Notes) > 0 {
				for index, note := range node.Notes {
					if note.NoteId == data.Notes[0].NoteId {
						node.Notes[index] = data.Notes[0]
						break
					}
				}
			}

		}
		cache.dll.moveNodeToFront(node)
		return true
	}
	return false
}

func (cache *TheCache) GetKey(id int) *Node {
	node, ok := cache.storage[id]
	if ok {
		cache.dll.moveNodeToFront(node)
		return node
	}
	return nil
}

func (cache *TheCache) GetUserKey(username string, password string) *Node {
	for _, node := range cache.storage {
		if node.UserName == username && node.Password == password {
			cache.dll.moveNodeToFront(node)
			return node
		}
	}
	return nil
}

//////////////////////////////////////////////////////	test functions

func PrintCache(cache TheCache) {
	fmt.Println("\n|||||||||||||||||\t Printing Cache Content \t|||||||||||||||||")
	for k, v := range cache.storage {
		fmt.Printf("%d:  %s\nNote:\t", k, v.UserName)
		for _, note := range v.Notes {
			fmt.Printf("%d -> %s(%s) \t ", note.NoteId, note.NoteTitle, note.NoteType)
		}
		fmt.Println("\n______________________________________________________________")
	}
	fmt.Println("\n|||||||||||||||||\t\t The End \t\t|||||||||||||||||\n\n")
}
