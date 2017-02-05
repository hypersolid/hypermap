package main

//
// import "unsafe"
//
// //go:linkname efaceHash runtime.efaceHash
// func efaceHash(i interface{}, seed uintptr) uintptr
//
// //go:linkname memhash runtime.memhash
// func memhash(p unsafe.Pointer, seed, s uintptr) uintptr
//
// func (m *Map) hashy(value interface{}) uint64 {
// 	return uint64(efaceHash(value, m.seed))
// }
