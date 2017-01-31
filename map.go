package main

import (
	"math/rand"
	"time"
	"unsafe"
)

const (
	maxRetries          = 3
	maxElementsInBucket = 8
	poolCoefficient     = 1
	maxLoad             = 0.9
)

var nullPtr unsafe.Pointer = unsafe.Pointer(uintptr(0))

// Map is awesome lockfree hashtable
type Map struct {
	array []unsafe.Pointer
	seed  uintptr

	Retries    uint64
	Collisions uint64
	Load       uint64

	pool []element
	head unsafe.Pointer

	usize uint64
	isize int
}

type element struct {
	key, value interface{}
	next       unsafe.Pointer
}

// NewMap is a constructor for the Map
func NewMap(size int) *Map {
	rand.Seed(time.Now().UTC().UnixNano())
	m := &Map{
		seed:  uintptr(rand.Int63()),
		array: make([]unsafe.Pointer, size),
		usize: uint64(size),
		isize: size,
	}

	m.pool = make([]element, size*maxElementsInBucket)
	m.head = unsafe.Pointer(&m.pool[0])
	for i := 1; i < size*maxElementsInBucket; i++ {
		m.pool[i-1].next = unsafe.Pointer(&m.pool[i])
	}

	return m
}
