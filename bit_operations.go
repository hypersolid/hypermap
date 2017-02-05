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

func bitsToString(value uint64) string {
	result := ""
	for i := uint64(0); i < 64; i++ {
		result += fmt.Sprintf("%d", (value>>(63-i))&1)
	}
	return result
}
