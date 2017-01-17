package main

import (
  "fmt"
  "unsafe"
  "sync/atomic"
)

const (
  mapSize = 5
)

type Map struct {
  arr []unsafe.Pointer
}

type record struct {
  key, value interface{}
  deleted bool
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
  return (*record)(m.arr[position]).value
}

func main() {
  m := NewMap(mapSize)

  for i := 0; i< mapSize; i++ {
    res, _ := m.set(i, i, "test")
    fmt.Println("set:", res, m.arr[i])
  }

  for i := 0; i< mapSize; i++ {
    fmt.Println("get:", m.get(i, i).(string))
  }
}
