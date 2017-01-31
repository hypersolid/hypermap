package main

import (
	"sync/atomic"
	"unsafe"
)

// grabs element from free pool or in case of high concurrency allocates new one
func (m *Map) elementFromPool() *element {
	if m.head != nullPtr {
		currentHead := (*element)(m.head)
		nextHead := currentHead.next
		if m.cas(&m.head, nextHead) {
			return currentHead
		}
	}
	atomic.AddUint64(&m.Allocations, 1)
	return new(element)
}

// executes compare and swap instruction
func (m *Map) cas(oldValuePointerLocation *unsafe.Pointer, newValuePointer unsafe.Pointer) bool {
	return atomic.CompareAndSwapPointer(
		oldValuePointerLocation,
		*oldValuePointerLocation,
		newValuePointer,
	)
}
