package hypermap

import "sync/atomic"

// Set adds or replaces entry in the map
func (m *Map) Set(key, value uint64) {
	m.RLock()
	defer m.RUnlock()

	var ptr, v uint64
	var P *uint64
	record := m.fuse(key, value)
	bucket := m.hashy(key) % m.size
	for {
	AGAIN:
		for step := uint64(0); step < m.size; step++ {
			ptr = (bucket + step) % m.size
			v = m.array[ptr]
			if !m.available(ptr) {
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
	m.Lock()
	defer m.Unlock()

	// fmt.Println("Grow", m.size)

	nm := &Map{
		array:     m.array,
		size:      m.size,
		keySize:   m.keySize,
		seed:      m.seed,
		valueSize: m.valueSize,
		keyMask:   m.keyMask,
		valueMask: m.valueMask,
	}

	m.size *= 2
	m.array = make([]uint64, m.size)
	m.markFree()

	for i := uint64(0); i < nm.size; i++ {
		if nm.deleted(i) || nm.available(i) {
			continue
		}
		bucket := m.hashy(nm.key(i)) % m.size
		for step := uint64(0); step < m.size; step++ {
			ptr := (bucket + step) % m.size
			if !m.available(ptr) {
				continue
			}
			m.array[ptr] = nm.array[i]
			break
		}
	}
}
