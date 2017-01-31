package main

import (
	"fmt"
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

var nullPtr unsafe.Pointer

// Map is awesome lockfree hashtable
type Map struct {
	array []unsafe.Pointer
	seed  uintptr

	Retries     uint64
	Collisions  uint64
	Load        uint64
	Allocations uint64
	MaxChain    uint64

	pool []element
	head unsafe.Pointer

	usize uint64
	isize int
}

func (m *Map) String() string {
	return fmt.Sprintf(
		"Retries:%d Collisions:%d Chain:%d Load:%.2f Allocations:%d | %d",
		m.Retries,
		m.Collisions,
		m.MaxChain,
		float64(m.Load)/float64(m.usize),
		m.Allocations,
		m.usize,
	)
}

// element is internal entity that is used in the map
type element struct {
	key, value interface{}
	next       unsafe.Pointer
}

// NewMap is a constructor for the Map
func NewMap(size int) *Map {
	rand.Seed(time.Now().UTC().UnixNano())
	m := &Map{
		seed:  uintptr(rand.Intn(256)),
		array: make([]unsafe.Pointer, size),
		usize: uint64(size),
		isize: size,
	}
	m.pool = make([]element, size*poolCoefficient)
	m.head = unsafe.Pointer(&m.pool[0])
	for i := 1; i < size*poolCoefficient; i++ {
		m.pool[i-1].next = unsafe.Pointer(&m.pool[i])
	}
	return m
}
