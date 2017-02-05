package main

import "unsafe"

//go:linkname memhash runtime.memhash
func memhash(p unsafe.Pointer, seed, s uintptr) uintptr

func (m *Map) hashy(value uint64) uint64 {
	return uint64(memhash(unsafe.Pointer(&value), m.seed, uintptr(8)))
}
