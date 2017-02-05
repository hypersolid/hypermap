package main

//
// // grabs element from free pool or in case of high concurrency allocates new one
// func (m *Map) elementFromPool() *element {
//
// 	// retries := 0
// 	// for {
// 	// hd :=
//
// 	if m.head != nullPtr {
// 		m.mutex.Lock()
// 		hd := (*element)(m.head)
// 		m.head = hd.next
// 		m.mutex.Unlock()
// 		return hd
// 		//
// 		// currentHead :=
// 		// nextHead := currentHead.next
// 		// m.mutex.Lock()
// 		// if atomic.CompareAndSwapPointer(
// 		// 	&m.head,
// 		// 	hd,
// 		// 	nextHead,
// 		// ) {
// 		// 	m.mutex.Unlock()
// 		// 	// atomic.AddUint64(&m.Swaps, 1)
// 		// 	return currentHead
// 		// }
// 		// m.mutex.Unlock()
// 		// // if retries < maxRetries {
// 		// // 	retries++
// 		// // 	continue
// 		// }
// 	}
// 	// atomic.AddUint64(&m.Allocations, 1)
// 	return new(element)
// 	// }
// }
