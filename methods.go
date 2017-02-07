package main

import "sync/atomic"

// Set adds or replaces entry in the map
func (m *Map) Set(key, value uint64) {
	var cell *uint64
	var mptr uint64
	record := m.fuse(key, value)
	bucket := m.hashy(key) % m.size
	for {
	AGAIN:
		for ptr := bucket; ptr < bucket+m.size; ptr++ {
			mptr = ptr % m.size
			currentValue, status := m.available(mptr)
			if !status {
				continue
			}
			cell = &(*m.array)[mptr]
			if atomic.CompareAndSwapUint64(
				cell,
				currentValue,
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
	var mptr uint64
	for ptr := bucket; ptr < bucket+m.size; ptr++ {
		mptr = ptr % m.size
		if m.deleted(mptr) {
			continue
		}
		if m.key(mptr) == key {
			return m.value(mptr), true
		}
	}
	return 0, false
}
