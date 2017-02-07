package main

import "sync/atomic"

// Set adds or replaces entry in the map
func (m *Map) Set(key, value uint64) {
	record := m.fuse(key, value)
	bucket := m.hashy(key) % m.size
	m.set(bucket, key, record)
}

func (m *Map) set(bucket, key, record uint64) {
	var ptr uint64
	for offset := uint64(0); offset < m.size; offset++ {
		ptr = (bucket + offset) % m.size
		currentValue, status := m.available(ptr)
		if status {
			if atomic.CompareAndSwapUint64(&(*m.array)[ptr], currentValue, record) {
				return
			}
			offset--
		}
	}
	panic("map is 100% full")
}

// Get grabs the value for the key
func (m *Map) Get(key uint64) (uint64, bool) {
	bucket := m.hashy(key) % m.size
	var ptr uint64
	for offset := uint64(0); offset < m.size; offset++ {
		ptr = (bucket + offset) % m.size
		if m.deleted(ptr) {
			continue
		}
		if m.key(ptr) == key {
			return m.value(ptr), true
		}
	}
	return 0, false
}
