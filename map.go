package main

import (
	"fmt"
	"sync"
	"unsafe"
)

const (
	initialSize = 35600000
)

// Map is awesome lockfree int->int only hashtable
type Map struct {
	sync.RWMutex
	array                           *[]uint64
	seed                            uintptr
	keySize, valueSize              uint64
	deletedMask, keyMask, valueMask uint64
	maxKey, maxValue                uint64

	ptrBuffer unsafe.Pointer

	size uint64
}

// String represents Map in human readable format
func (m *Map) String() string {
	return fmt.Sprintf(
		"HyperMap<%d %d bits -> %d bits>",
		unsafe.Pointer(m),
		m.keySize,
		m.valueSize,
	)
}

// NewMap is a constructor for the Map
func NewMap(keySize uint64) *Map {
	if keySize < 16 {
		panic("For maps with key < 16 bits use simple array")
	}
	if keySize > 62 {
		panic("Maximum keySize is 62 bits")
	}

	array := make([]uint64, initialSize)
	m := &Map{
		array:     &array,
		size:      initialSize,
		keySize:   keySize,
		valueSize: 63 - keySize,
	}

	m.generateMasks()
	m.markFree()

	m.maxKey = (m.keyMask >> m.valueSize) - 1
	m.maxValue = m.valueMask

	return m
}
