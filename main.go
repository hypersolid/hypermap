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
	arr        []unsafe.Pointer
	seed       uintptr
	load       uint64
	writes     uint64
	hits       uint64
	misses     uint64
	reads      uint64
	collisions uint64
	retries    uint64
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

func (m *Map) wrap(position uint64) int {
	return int(position % uint64(len(m.arr)))
}

func (m *Map) entityAt(position int) *entry {
	entryPointer := m.arr[position]
	if entryPointer == nil {
		return nil
	}
	return (*entry)(entryPointer)
}

func (m *Map) probe(position uint64, step int) uint64 {
	return position + uint64(step)
}

func (m *Map) cas(position int, rp *entry) bool {
	circularPosition := position % len(m.arr)
	return atomic.CompareAndSwapPointer(
		&m.arr[circularPosition],
		unsafe.Pointer(m.arr[circularPosition]),
		unsafe.Pointer(rp),
	)
}

func (m *Map) Set(key, value interface{}) bool {
	r := 0
	h := m.hashy(key)
	for !m.set(key, value, h) {
		atomic.AddUint64(&(m.retries), 1)
		r++
		if r > 10 {
			return false
		}
	}
	return true
}

// Set creates or replaces entry in the Map
func (m *Map) set(key, value interface{}, position uint64) bool {
	rp := &entry{key: key, value: value}
	for i := 0; i < len(m.arr); i++ {
		pos := m.wrap(m.probe(position, i))
		entity := m.entityAt(pos)
		if entity != nil && entity.key != key {
			atomic.AddUint64(&m.collisions, 1)
			continue
		}
		atomic.AddUint64(&m.load, 1)
		return m.cas(pos, rp)
	}
	return false
}

// Get reads entry from the Map, in case entry does not exist returns nil
func (m *Map) Get(key interface{}) interface{} {
	position := m.hashy(key)
	for i := 0; i < len(m.arr); i++ {
		pos := m.wrap(m.probe(position, i))
		entity := m.entityAt(pos)
		if entity != nil && entity.key == key {
			return entity.value
		}
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
