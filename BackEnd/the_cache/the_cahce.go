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
	//const MAX_CAPACITY
}

func InitCache() TheCache {
	storage := make(map[int]*Node)
	return TheCache{
		dll:     &DoublyLinkedList{},
		storage: storage,
		// MAX_Capacity:      maxCapacity,
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
	_, ok := cache.storage[data.UserId]
	if ok {
		cache.dll.moveNodeToFront(data)
		return true
	}
	if cache.dll.size() >= CACHE_CAPACITY {
		delete(cache.storage, cache.dll.tail.UserId)
		cache.dll.removeFromEnd()
	}
	cache.storage[data.UserId] = data
	cache.dll.addToFront(data)
	return len(cache.storage) == cache.dll.size()
}

func (cache *TheCache) GetKey(id int) *Node {
	node, ok := cache.storage[id]
	if ok {
		return node
	}
	return nil
}

// func main() {
// 	neg := Node{
// 		UserId:  12,
// 		UserName: "nima",
// 		Password: "159159",
// 		Name:     "bookWorn",
// 		Notes:    []Note{},
// 	}
// 	dll.AddToFront(&neg)
// }
