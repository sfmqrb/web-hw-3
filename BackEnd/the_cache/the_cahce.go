package thecache

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

// var dll *DoublyLinkedList = initDoublyList(CACHE_CAPACITY)

type TheCache struct {
	dll     *DoublyLinkedList
	storage map[int]*Node
}

func InitCache() TheCache {
	storage := make(map[int]*Node)
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
	// todo update data of node instead
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
	// todo update data of data instead
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
			//////////////////////////
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
				node.Notes = append(node.Notes, data.Notes...)
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
