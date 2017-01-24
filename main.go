package main

import (
	"math/rand"
	"sync/atomic"
	"unsafe"
)

//go:linkname memhash runtime.memhash
func memhash(p unsafe.Pointer, seed, s uintptr) uintptr

// Map is awesome lockfree hashtable
type Map struct {
	arr  []unsafe.Pointer
	seed uintptr
}

type entry struct {
	key, value interface{}
	deleted    bool
}

// NewMap is a constructor for the Map
func NewMap(size int) *Map {
	m := Map{seed: uintptr(rand.Int63())}
	m.arr = make([]unsafe.Pointer, size)
	return &m
}

// Set creates or replaces entry in the Map
func (m *Map) Set(key, value interface{}) bool {
	position := m.hashy(key)
	rp := &entry{key: key, value: value}
	result := atomic.CompareAndSwapPointer(
		&m.arr[position],
		unsafe.Pointer(m.arr[position]),
		unsafe.Pointer(rp),
	)
	return result
}

// Get reads entry from the Map, in case entry does not exist returns nil
func (m *Map) Get(key interface{}) interface{} {
	position := m.hashy(key)
	entryPointer := m.arr[position]
	if entryPointer == nil {
		return nil
	}
	entryStruct := (*entry)(entryPointer)
	return entryStruct.value
}

func (m *Map) hashy(value interface{}) uint64 {
	modulo := uint64(len(m.arr))
	switch value.(type) {
	case int:
		vv := value.(int)
		return uint64(memhash(unsafe.Pointer(&vv), m.seed, 2)) % modulo
	case string:
		vv := value.(string)
		return uint64(memhash(unsafe.Pointer(&vv), m.seed, 2)) % modulo
	default:
		panic("unknown")
	}
	return 0
}
