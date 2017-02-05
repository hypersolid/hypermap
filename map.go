package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
	"unsafe"
)

const (
	initialSize = 256
)

// Map is awesome lockfree int->int only hashtable
type Map struct {
	sync.RWMutex
	array                           *[]uint64
	seed                            uintptr
	keySize, valueSize              uint64
	deletedMask, keyMask, valueMask uint64

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

	rand.Seed(time.Now().UTC().UnixNano())
	array := make([]uint64, initialSize)
	m := &Map{
		seed:      uintptr(rand.Int()),
		array:     &array,
		size:      initialSize,
		keySize:   keySize,
		valueSize: 63 - keySize,
	}

	m.generateMasks()
	m.markFree()

	return m
}
