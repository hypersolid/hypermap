package main

import (
	"fmt"
	"sync/atomic"
)

// Set adds or replaces entry in the map
func (m *Map) Set(key, value uint64) {
	record := m.fuse(key, value)
	bucket := m.hashy(key) % m.size
	m.set(bucket, key, record)
}

func (m *Map) set(ptr, key, record uint64) {
	for {
		currentValue, status := m.available(ptr)
		if status && atomic.CompareAndSwapUint64(&(*m.array)[ptr], currentValue, record) {
			// fmt.Println(bitsToString((*m.array)[ptr]), "\n", bitsToString(currentValue), status, "X")
			return
		} else {
			fmt.Println((*m.array)[ptr], currentValue, status)
			return
		}
		// ptr = (ptr + 1) % m.size
	}
}

//
// // Get reads element from the Map, in case element does not exist returns nil
// func (m *Map) Get(key interface{}) interface{} {
// 	bucket := m.hashy(key) % m.usize
// 	next := m.array[bucket]
// 	var entity *element
// 	for uintptr(next) != uintptr(0) {
// 		entity = (*element)(next)
// 		if entity != nil && entity.key == key {
// 			return entity.value
// 		}
// 		next = entity.next
// 	}
// 	return nil
// }
