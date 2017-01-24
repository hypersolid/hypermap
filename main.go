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
	arr    []unsafe.Pointer
	seed   uintptr
	filled int
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

func (m *Map) entityAt(position int) *entry {
	circularPosition := position % len(m.arr)
	entryPointer := m.arr[circularPosition]
	if entryPointer == nil {
		return nil
	}
	return (*entry)(entryPointer)
}

func (m *Map) cas(position int, rp *entry) bool {
	circularPosition := position % len(m.arr)
	return atomic.CompareAndSwapPointer(
		&m.arr[circularPosition],
		unsafe.Pointer(m.arr[circularPosition]),
		unsafe.Pointer(rp),
	)
}

// Set creates or replaces entry in the Map
func (m *Map) Set(key, value interface{}) bool {
	position := int(m.hashy(key) % uint64(len(m.arr)))
	for i := 0; i < len(m.arr); i++ {
		entity := m.entityAt(position + i)
		if entity != nil && entity.key != key {
			continue
		}
		m.filled++
		rp := &entry{key: key, value: value}
		return m.cas(position+i, rp)
	}

	return false
}

// Get reads entry from the Map, in case entry does not exist returns nil
func (m *Map) Get(key interface{}) interface{} {
	position := int(m.hashy(key) % uint64(len(m.arr)))
	for i := 0; i < len(m.arr); i++ {
		entity := m.entityAt(position + i)
		if entity != nil && entity.key == key {
			return entity.value
		}
	}
	return nil
}

func (m *Map) hashy(value interface{}) uint64 {
	switch value.(type) {
	case int:
		vv := value.(int)
		return uint64(memhash(unsafe.Pointer(&vv), m.seed, 4))
	case string:
		vv := value.(string)
		return uint64(memhash(unsafe.Pointer(&vv), m.seed, 4))
	default:
		panic("unknown")
	}
}
