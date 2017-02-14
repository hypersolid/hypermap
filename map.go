package hypermap

import (
	"fmt"
	"sync"
)

// Map is awesome lockfree int->int only hashtable
type Map struct {
	sync.RWMutex
	array              []uint64
	seed               uintptr
	keySize, valueSize uint64
	keyMask, valueMask uint64
	maxKey, maxValue   uint64
	size               uint64
}

// String represents Map in human readable format
func (m *Map) String() string {
	return fmt.Sprintf(
		"HyperMap<%d bits -> %d bits>",
		m.keySize,
		m.valueSize,
	)
}

// NewMap is a constructor for the Map
func NewMap(keySize uint64, size uint64) *Map {
	if keySize < 16 {
		panic("For maps with key < 16 bits use simple array")
	}
	if keySize > 63 {
		panic("Maximum keySize is 62 bits")
	}

	m := &Map{
		array:     make([]uint64, size),
		size:      size,
		keySize:   keySize,
		valueSize: 64 - keySize,
	}

	m.generateMasks()
	m.markFree()

	m.maxKey = (m.keyMask >> m.valueSize) - 1
	m.maxValue = m.valueMask - 1

	return m
}
