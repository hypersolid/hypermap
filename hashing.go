package main

import "unsafe"

//go:linkname memhash runtime.memhash
func memhash(p unsafe.Pointer, seed, s uintptr) uintptr

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
