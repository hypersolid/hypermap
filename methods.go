package hypermap

import "sync/atomic"

// Set adds or replaces entry in the map
func (m *Map) Set(key, value uint64) {
	if key > m.maxKey {
		panic("key too large " + bitsToString(key))
	}
	if value > m.maxValue {
		panic("value too large " + bitsToString(value))
	}
	for !m.parallelSet(key, value) {
		// we need to check if failure was caused by high map load
		if float64(m.load)/float64(m.size) > 0.95 {
			newSize := m.size * 3
			m.Lock()
			m.grow(newSize)
			m.Unlock()
		}
	}
}

func (m *Map) parallelSet(key, value uint64) bool {
	m.RLock()

	var lockedValue, ptr uint64
	bucket := m.hashy(key) % m.size
	for step := uint64(0); step < m.size; step++ {
		ptr = (bucket + step) % m.size
		lockedValue = m.array[ptr]
		replaceable := m.key(ptr) == key
		if m.available(ptr) || replaceable {
			status := atomic.CompareAndSwapUint64(
				&m.array[ptr],
				lockedValue,
				m.fuse(key, value),
			)
			if status && !replaceable {
				atomic.AddUint64(&m.load, 1)
			}
			m.RUnlock()
			return status
		}
	}
	m.RUnlock()
	return false
}

func (m *Map) syncSet(key, value uint64) {
	bucket := m.hashy(key) % m.size
	for step := uint64(0); step < m.size; step++ {
		ptr := (bucket + step) % m.size
		if m.key(ptr) == key {
			m.array[ptr] = m.fuse(key, value)
			return
		}
		if m.available(ptr) {
			atomic.AddUint64(&m.load, 1)
			m.array[ptr] = m.fuse(key, value)
			return
		}
	}
}

// Deletes the key
func (m *Map) Delete(key uint64) {
	m.parallelSet(key, m.valueMask)
}

// Get grabs the value for the key
func (m *Map) Get(key uint64) (uint64, bool) {
	m.RLock()

	bucket := m.hashy(key) % m.size
	var ptr uint64
	for step := uint64(0); step < m.size; step++ {
		ptr = (bucket + step) % m.size
		if m.deleted(ptr) {
			continue
		}
		if m.key(ptr) == key {
			m.RUnlock()
			return m.value(ptr), true
		}
	}
	m.RUnlock()
	return 0, false
}

func (m *Map) grow(newSize uint64) {
	if m.size >= newSize {
		// map was already grown by an other thread
		// thus do nothing
		return
	}
	oldMap := &Map{
		array:     m.array,
		size:      m.size,
		seed:      m.seed,
		keySize:   m.keySize,
		valueSize: m.valueSize,
		keyMask:   m.keyMask,
		valueMask: m.valueMask,
	}
	m.load = 0
	m.size = newSize
	m.array = make([]uint64, newSize)
	m.markFree()
	for i := uint64(0); i < oldMap.size; i++ {
		if oldMap.available(i) {
			continue
		}
		m.syncSet(oldMap.key(i), oldMap.value(i))
	}
}
