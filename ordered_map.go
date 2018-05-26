package orderedmap

import (
	"fmt"
	"container/list"
)

type OrderedMap struct {
	store  map[interface{}]interface{}
	mapper map[interface{}]*list.Element
	list   *list.List
}

func NewOrderedMap() *OrderedMap {
	om := &OrderedMap{
		store:  make(map[interface{}]interface{}),
		mapper: make(map[interface{}]*list.Element),
		list:   list.New(),
	}
	return om
}

func NewOrderedMapWithArgs(args []*KVPair) *OrderedMap {
	om := NewOrderedMap()
	om.update(args)
	return om
}

func (om *OrderedMap) update(args []*KVPair) {
	for _, pair := range args {
		om.Set(pair.Key, pair.Value)
	}
}

func (om *OrderedMap) Set(key interface{}, value interface{}) {
	if _, ok := om.store[key]; ok == false {
		om.list.PushBack(key)
		//last := root.Prev
		//last.Next = newNode(last, root, key)
		//root.Prev = last.Next
		om.mapper[key] = om.list.Back()
	}
	om.store[key] = value
}

func (om *OrderedMap) Get(key interface{}) (interface{}, bool) {
	val, ok := om.store[key]
	return val, ok
}

func (om *OrderedMap) Delete(key interface{}) {
	_, ok := om.store[key]
	if ok {
		delete(om.store, key)
	}
	node, found := om.mapper[key]
	if found {
		om.list.Remove(node)
		delete(om.mapper, key)
	}
}

func (om *OrderedMap) String() string {
	builder := make([]string, len(om.store))

	var index int = 0
	iter := om.IterFunc()
	for kv, ok := iter(); ok; kv, ok = iter() {
		val, _ := om.Get(kv.Key)
		builder[index] = fmt.Sprintf("%v:%v, ", kv.Key, val)
		index++
	}
	return fmt.Sprintf("OrderedMap%v", builder)
}

func (om *OrderedMap) Iter() <-chan *KVPair {
	println("Iter() method is deprecated!. Use IterFunc() instead.")
	return om.UnsafeIter()
}

/*
Beware, Iterator leaks goroutines if we do not fully traverse the map.
For most cases, `IterFunc()` should work as an iterator.
 */
func (om *OrderedMap) UnsafeIter() <-chan *KVPair {
	keys := make(chan *KVPair)
	go func() {
		defer close(keys)
		for it := om.list.Front(); it != nil; it = it.Next() {
			v, _ := om.store[it.Value]
			keys <- &KVPair{it.Value, v}
		}
	}()
	return keys
}

func (om *OrderedMap) IterFunc() func() (*KVPair, bool) {
	it := om.list.Front()
	return func() (*KVPair, bool) {
		if it != nil {
			ret := &KVPair{it.Value, om.store[it.Value]}
			it = it.Next()
			return ret, true
		}
		return nil, false
	}
}

func (om *OrderedMap) Len() int {
	return len(om.store)
}
