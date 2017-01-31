package main

import (
	"math/rand"
	"sync/atomic"
	"time"
	"unsafe"
)

const (
	maxRetries          = 3
	maxElementsInBucket = 8
)

var nullPtr unsafe.Pointer = unsafe.Pointer(uintptr(0))

//go:linkname memhash runtime.memhash
func memhash(p unsafe.Pointer, seed, s uintptr) uintptr

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

func (m *Map) pop() *element {
	if m.head == nullPtr {
		return new(element)
	}
	currentHead := (*element)(m.head)
	nextHead := currentHead.next
	if m.cas(&m.head, nextHead) {
		return currentHead
	}
	return new(element)
}

func (m *Map) cas(oldValuePointerLocation *unsafe.Pointer, newValuePointer unsafe.Pointer) bool {
	return atomic.CompareAndSwapPointer(
		oldValuePointerLocation,
		*oldValuePointerLocation,
		newValuePointer,
	)
}

// Set adds or replaces entry in the map
func (m *Map) Set(key, value interface{}) bool {
	retries := 0
	bucket := m.hashy(key) % m.usize
	for !m.set(key, value, bucket) {
		// atomic.AddUint64(&m.Retries, 1)
		retries++
		if retries > maxRetries {
			return false
		}
	}
	// atomic.AddUint64(&m.Load, 1)
	return true
}

func (m *Map) set(key, value interface{}, bucket uint64) bool {
	prev := &m.array[bucket]
	next := m.array[bucket]
	var entity *element

	for uintptr(next) != uintptr(0) {
		entity = (*element)(next)
		if entity != nil && entity.key == key {
			newElement := m.pop()
			newElement.key = key
			newElement.value = value
			newElement.next = entity.next
			return m.cas(prev, unsafe.Pointer(newElement))
		}
		prev = &entity.next
		next = entity.next
	}

	newElement := m.pop()
	newElement.key = key
	newElement.value = value
	newElement.next = m.array[bucket]
	return m.cas(&m.array[bucket], unsafe.Pointer(newElement))
}

// Get reads element from the Map, in case element does not exist returns nil
func (m *Map) Get(key interface{}) interface{} {
	bucket := m.hashy(key) % m.usize
	next := m.array[bucket]
	var entity *element
	for uintptr(next) != uintptr(0) {
		entity = (*element)(next)
		if entity != nil && entity.key == key {
			return entity.value
		}
		next = entity.next
	}
	return nil
}

func (m *Map) hashy(value interface{}) uint64 {
	switch value.(type) {
	case int:
		realValue := value.(int)
		return uint64(memhash(unsafe.Pointer(&realValue), m.seed, 4))
	case string:
		realValue := value.(string)
		return uint64(memhash(unsafe.Pointer(&realValue), m.seed, 4))
	default:
		panic("unknown type")
	}
}
