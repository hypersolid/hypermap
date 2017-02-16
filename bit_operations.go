package hypermap

import "fmt"

func (m *Map) generateMasks() {
	for i := uint64(0); i < m.keySize; i++ {
		m.keyMask |= 1 << (i + m.valueSize)
	}
	for i := uint64(0); i < m.valueSize; i++ {
		m.valueMask |= 1 << i
	}
}

func (m *Map) markFree() {
	for i := uint64(0); i < m.size; i++ {
		m.array[i] |= m.keyMask
	}
}

func (m *Map) fuse(key, value uint64) uint64 {
	return (key << m.valueSize) | value
}

func (m *Map) available(bucket uint64) bool {
	if (m.array[bucket]&m.keyMask) == m.keyMask || (m.array[bucket]&m.valueMask) == m.valueMask {
		return true
	}
	return false
}

func (m *Map) deleted(bucket uint64) bool {
	return (m.array[bucket] & m.valueMask) == m.valueMask
}

func (m *Map) key(bucket uint64) uint64 {
	return (m.array[bucket] & m.keyMask) >> m.valueSize
}

func (m *Map) value(bucket uint64) uint64 {
	return m.array[bucket] & m.valueMask
}

func bitsToString(value uint64) string {
	result := ""
	var tmp uint64
	for i := uint64(0); i < 64; i++ {
		tmp = (value >> (63 - i)) & 1
		result += fmt.Sprintf("%d", tmp)
	}
	return result
}
