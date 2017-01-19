package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"hash/crc32"
	"sync/atomic"
	"unsafe"
)

const (
	mapSize = 10
)

type Map struct {
	arr []unsafe.Pointer
}

type record struct {
	key, value interface{}
	deleted    bool
}

func NewMap(size int) *Map {
	m := Map{}
	m.arr = make([]unsafe.Pointer, size)
	return &m
}

func (m *Map) Set(key, value interface{}) (bool, error) {
	position := m.Hash(key)
	rp := &record{key: key, value: value}
	addr := (*unsafe.Pointer)(&m.arr[position])
	oldPtr := unsafe.Pointer(m.arr[position])
	newPtr := unsafe.Pointer(rp)
	result := atomic.CompareAndSwapPointer(addr, oldPtr, newPtr)
	return result, nil
}

func (m *Map) Get(key interface{}) interface{} {
	position := m.Hash(key)
	res := m.arr[position]
	if res == nil {
		return nil
	}
	ff := (*record)(res)
	return ff.value
}

func (m *Map) Hash(value interface{}) int {
	pos := int(crc32.ChecksumIEEE(GetBytes(value).Bytes())) % len(m.arr)
	return pos
}

func main() {
	m := NewMap(mapSize * 2)

	for i := 0; i < mapSize; i++ {
		m.Set(i, "test")
	}

	for i := 0; i < mapSize; i++ {
		res := m.Get(i)
		if res != nil {
			fmt.Println("get:", res.(string), i)
		}
	}
}

func GetBytes(key interface{}) *bytes.Buffer {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	enc.Encode(key)
	return &buf
}
