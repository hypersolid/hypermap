package main

import (
	"sync/atomic"
	"unsafe"
)

func (m *Map) pop() *element {
	if m.head == nullPtr {
		return new(element)
	}
	currentHead := (*element)(m.head)
	nextHead := currentHead.next
	if m.cas(&m.head, nextHead) {
		return currentHead
	}
	return new(element)
}

func (m *Map) cas(oldValuePointerLocation *unsafe.Pointer, newValuePointer unsafe.Pointer) bool {
	return atomic.CompareAndSwapPointer(
		oldValuePointerLocation,
		*oldValuePointerLocation,
		newValuePointer,
	)
}
