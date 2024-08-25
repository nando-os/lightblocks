package ordered_map

import (
	"container/list"
	"fmt"
	"sync"
)

type Map interface {
	Add(key string, value string)
	Delete(key string)
	Get(key string) (string, bool)
	Size() int
	GetAllItems() []KeyValue
}

// KeyValue pair stored in the linked list
type KeyValue struct {
	Key   string
	Value string
}

type orderedMap struct {
	items        map[string]*list.Element
	orderedItems *list.List
	mu           sync.RWMutex
}

func NewConcurrentOrderedMap() Map {
	return &orderedMap{
		items:        make(map[string]*list.Element),
		orderedItems: list.New(),
		mu:           sync.RWMutex{},
	}
}

func (o *orderedMap) Add(key, value string) {
	o.mu.Lock()
	defer o.mu.Unlock()

	if el, found := o.items[key]; found {
		el.Value.(*KeyValue).Value = value
	} else {
		kv := &KeyValue{Key: key, Value: value}
		o.items[key] = o.orderedItems.PushFront(kv)
	}
}

func (o *orderedMap) Delete(key string) {
	o.mu.Lock()
	defer o.mu.Unlock()

	if el, found := o.items[key]; found {
		o.orderedItems.Remove(el)
		delete(o.items, key)
	}
}

func (o *orderedMap) Get(key string) (string, bool) {
	o.mu.RLock()
	defer o.mu.RUnlock()
	if el, found := o.items[key]; found {
		return el.Value.(*KeyValue).Value, true
	}
	return "", false
}
func (o *orderedMap) Size() int {
	return len(o.items)
}

func (o *orderedMap) GetAllItems() []KeyValue {
	o.mu.RLock()
	defer o.mu.RUnlock()

	items := make([]KeyValue, 0, o.Size())
	for el := o.orderedItems.Front(); el != nil; el = el.Next() {
		items = append(items, *el.Value.(*KeyValue))
	}
	return items
}

func (k *KeyValue) String() string {
	return fmt.Sprintf("key=%s, value=%s", k.Key, k.Value)
}
