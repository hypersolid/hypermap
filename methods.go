package main

import "unsafe"

// Set adds or replaces entry in the map
func (m *Map) Set(key, value interface{}) bool {
	retries := 0
	bucket := m.hashy(key) % m.usize
	for !m.set(key, value, bucket) {
		// atomic.AddUint64(&m.Retries, 1)
		retries++
		if retries > maxRetries {
			return false
		}
	}
	// atomic.AddUint64(&m.Load, 1)
	return true
}

func (m *Map) set(key, value interface{}, bucket uint64) bool {
	prev := &m.array[bucket]
	next := m.array[bucket]
	var entity *element

	for uintptr(next) != uintptr(0) {
		entity = (*element)(next)
		if entity != nil && entity.key == key {
			newElement := m.pop()
			newElement.key = key
			newElement.value = value
			newElement.next = entity.next
			return m.cas(prev, unsafe.Pointer(newElement))
		}
		prev = &entity.next
		next = entity.next
	}

	newElement := m.pop()
	newElement.key = key
	newElement.value = value
	newElement.next = m.array[bucket]
	return m.cas(&m.array[bucket], unsafe.Pointer(newElement))
}

// Get reads element from the Map, in case element does not exist returns nil
func (m *Map) Get(key interface{}) interface{} {
	bucket := m.hashy(key) % m.usize
	next := m.array[bucket]
	var entity *element
	for uintptr(next) != uintptr(0) {
		entity = (*element)(next)
		if entity != nil && entity.key == key {
			return entity.value
		}
		next = entity.next
	}
	return nil
}
