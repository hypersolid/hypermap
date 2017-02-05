package main

import "fmt"

func (m *Map) generateMasks() {
	m.deletedMask = 1 << 63
	for i := uint64(0); i < m.keySize; i++ {
		m.keyMask |= 1 << (i + m.valueSize)
	}
	for i := uint64(0); i < m.valueSize; i++ {
		m.valueMask |= 1 << i
	}
}

func (m *Map) markFree() {
	for i := uint64(0); i < m.size; i++ {
		(*m.array)[i] |= m.keyMask
	}
}

func (m *Map) fuse(key, value uint64) uint64 {
	if key > m.maxKey {
		panic("key too large")
	}
	if value > m.maxValue {
		panic("value too large")
	}
	return (key << m.valueSize) | value
}

func (m *Map) available(bucket uint64) (uint64, bool) {
	current := (*m.array)[bucket]
	if (current & m.deletedMask) == m.deletedMask {
		return current, true
	}
	if (current & m.keyMask) == m.keyMask {
		return current, true
	}
	return 0, false
}

func bitsToString(value uint64) string {
	result := ""
	for i := uint64(0); i < 64; i++ {
		result += fmt.Sprintf("%d", (value>>(63-i))&1)
	}
	return result
}
