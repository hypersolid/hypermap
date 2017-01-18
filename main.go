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

func (m *Map) set(position int, key, value interface{}) (bool, error) {
	rp := &record{key: key, value: value}
	addr := (*unsafe.Pointer)(&m.arr[position])
	oldPtr := unsafe.Pointer(m.arr[position])
	newPtr := unsafe.Pointer(rp)
	result := atomic.CompareAndSwapPointer(addr, oldPtr, newPtr)
	return result, nil
}

func (m *Map) get(position int, key interface{}) interface{} {
	res := m.arr[position]
	if res == nil {
		return nil
	}
	ff := (*record)(res)
	return ff.value
}

func main() {
	m := NewMap(mapSize)

	for i := 0; i < mapSize/2; i++ {
		bts, _ := GetBytes(i)
		position := Hash(bts)
		res, _ := m.set(int(position), i, "test")
		fmt.Println("set:", res, m.arr[position], position)
	}

	for i := 0; i < mapSize; i++ {
		res := m.get(i, i)
		if res != nil {
			fmt.Println("get:", res.(string), i)
		}
	}
}

func GetBytes(key interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(key)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func Hash(data []byte) uint32 {
	return crc32.ChecksumIEEE(data) % mapSize
}
