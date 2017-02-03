package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
	"unsafe"
)

const (
	maxRetries          = 0
	maxElementsInBucket = 8
	poolCoefficient     = 4
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
	Swaps       uint64

	pool []element
	head unsafe.Pointer

	usize uint64
	isize int

	mutex sync.RWMutex
}

func (m *Map) String() string {

	ratio := -1.0
	if m.Swaps < m.Allocations {
		ratio = float64(m.Swaps) / float64(m.Allocations)
	}
	return fmt.Sprintf(
		"Retries:%d Collisions:%d Chain:%d Load:%.2f Picks/Allocations: %.2f | %d",
		m.Retries,
		m.Collisions,
		m.MaxChain,
		float64(m.Load)/float64(m.usize),
		ratio,
		m.usize,
	)
}

// element is internal entity that is used in the map
type element struct {
	next       unsafe.Pointer
	key, value interface{}
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
