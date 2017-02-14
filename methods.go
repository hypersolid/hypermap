package hypermap

import (
	"fmt"
	"sync/atomic"
)

// Set adds or replaces entry in the map
func (m *Map) Set(key, value uint64) {
	m.RLock()
	var ptr, v uint64
	var P *uint64
	bucket := m.hashy(key) % m.size
AGAIN:
	for step := uint64(0); step < m.size*9/10; step++ {
		ptr = (bucket + step) % m.size
		v = m.array[ptr]
		if !m.available(ptr) {
			continue
		}
		P = &m.array[ptr]
		if atomic.CompareAndSwapUint64(
			P,
			v,
			m.fuse(key, value),
		) {
			m.RUnlock()
			return
		}
		goto AGAIN
	}
	m.RUnlock()

	m.Lock()
	m.grow()
	m.syncSet(key, value)
	m.Unlock()
}

func (m *Map) syncSet(key, value uint64) {
	bucket := m.hashy(key) % m.size
	for step := uint64(0); step < m.size; step++ {
		ptr := (bucket + step) % m.size
		if !m.available(ptr) && !m.deleted(ptr) {
			continue
		}
		m.array[ptr] = m.fuse(key, value)
		return
	}
}

// Get grabs the value for the key
func (m *Map) Get(key uint64) (uint64, bool) {
	m.RLock()
	defer m.RUnlock()

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

func (m *Map) grow() {
	fmt.Println("grow to ", m.size*2)
	nm := &Map{
		array:     m.array,
		size:      m.size,
		seed:      m.seed,
		keySize:   m.keySize,
		valueSize: m.valueSize,
		keyMask:   m.keyMask,
		valueMask: m.valueMask,
	}
	m.size *= 2
	m.array = make([]uint64, m.size)
	m.markFree()
	for i := uint64(0); i < nm.size; i++ {
		if nm.available(i) || nm.deleted(i) {
			continue
		}
		m.syncSet(nm.key(i), nm.value(i))
	}
}
