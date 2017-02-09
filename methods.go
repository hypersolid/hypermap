package main

import "sync/atomic"

// Set adds or replaces entry in the map
func (m *Map) Set(key, value uint64) {
	var ptr, v uint64
	var P *uint64
	record := m.fuse(key, value)
	bucket := m.hashy(key) % m.size
	for {
	AGAIN:
		for step := uint64(0); step < m.size; step++ {
			ptr = (bucket + step) % m.size
			v = m.array[ptr]
			if !m.available(v) {
				continue
			}
			P = &m.array[ptr]
			if atomic.CompareAndSwapUint64(
				P,
				v,
				record,
			) {
				return
			}
			goto AGAIN
		}
		panic("map is 100% full")
	}
}

// Get grabs the value for the key
func (m *Map) Get(key uint64) (uint64, bool) {
	bucket := m.hashy(key) % m.size
	var ptr uint64
	for step := uint64(0); step < m.size; step++ {
		ptr = (bucket + step) % m.size
		if m.deleted(ptr) {
			continue
		}
		if m.key(ptr) == key {
			return m.value(ptr), true
		}
	}
	return 0, false
}
