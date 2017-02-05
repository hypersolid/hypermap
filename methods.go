package main

//
// import (
// 	"sync/atomic"
// 	"unsafe"
// )
//
// // Set adds or replaces entry in the map
// func (m *Map) Set(key, value interface{}) bool {
// 	bucket := m.hashy(key) % m.usize
// 	el := m.elementFromPool()
// 	el.key = key
// 	el.value = value
// 	for !m.set(key, value, bucket, el) {
// 		// atomic.AddUint64(&m.Retries, 1)
// 		// runtime.Gosched()
// 	}
// 	return true
// }
//
// func (m *Map) set(key, value interface{}, bucket uint64, el *element) bool {
// 	var entity *element
// 	parentPtrLocation := &m.array[bucket]
// 	currentElement := m.array[bucket]
//
// 	var steps uint64
// 	for currentElement != nullPtr {
// 		steps++
// 		entity = (*element)(currentElement)
// 		if entity.key == key {
// 			// atomic.AddUint64(&m.Collisions, 1)
// 			el.next = entity.next
// 			return atomic.CompareAndSwapPointer(
// 				parentPtrLocation,
// 				currentElement,
// 				unsafe.Pointer(el),
// 			)
// 		}
// 		parentPtrLocation = &entity.next
// 		currentElement = entity.next
// 	}
// 	// if steps > 0 {
// 	// 	if steps > m.MaxChain {
// 	// 		m.MaxChain = steps
// 	// 	}
// 	// } else {
// 	// 	// atomic.AddUint64(&m.Load, 1)
// 	// }
// 	el.next = currentElement
// 	return atomic.CompareAndSwapPointer(
// 		parentPtrLocation,
// 		currentElement,
// 		unsafe.Pointer(el),
// 	)
// }
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
