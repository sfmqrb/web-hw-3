package thecache

//import "fmt"

// configurables
const (
	PORT           int = 5432
	CACHE_CAPACITY int = 20
)

/*  keys and values are int32 or string
keys -> 64 char / values -> 2048 char
commands: getkey, setkey, clear
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

func (cache *TheCache) clear() {
	cache.dll = initDoublyList(CACHE_CAPACITY)
	cache.storage = make(map[int]*Node)
}

func (cache *TheCache) setKey(id int, node *Node) bool {
	_, ok := cache.storage[id]
	if ok {
		cache.dll.moveNodeToFront(node)
		return true
	}
	if cache.dll.size() >= CACHE_CAPACITY {
		delete(cache.storage, cache.dll.tail.user_id)
		cache.dll.removeFromEnd()
	}
	cache.storage[id] = node
	cache.dll.addToFront(node)
	return len(cache.storage) == cache.dll.size()
}

func (cache *TheCache) getKey(id int) *Node {
	node, ok := cache.storage[id]
	if ok {
		return node
	}
	return nil
}

// func main() {
// 	neg := Node{
// 		user_id:  12,
// 		username: "nima",
// 		password: "159159",
// 		name:     "bookWorn",
// 		notes:    []Note{},
// 	}
// 	dll.AddToFront(&neg)
// }
