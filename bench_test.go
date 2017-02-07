package main

import (
	"math/rand"
	"runtime"
	"testing"
)

const (
	testProcs = 1
)

func Benchmark_hypermap_parallel_write(b *testing.B) {
	b.StopTimer()
	procs := runtime.GOMAXPROCS(testProcs)
	start := make(chan struct{}, testProcs)
	done := make(chan struct{}, testProcs)
	m := NewMap(40)
	for k := 0; k < testProcs; k++ {
		go func() {
			<-start
			value := uint64(rand.Intn(int(m.maxValue)))
			for i := uint64(0); i < uint64(b.N)/testProcs; i++ {
				m.Set(i, value)
			}
			done <- struct{}{}
		}()
	}
	b.StartTimer()

	for k := 0; k < testProcs; k++ {
		start <- struct{}{}
	}
	for k := 0; k < testProcs; k++ {
		<-done
	}
	runtime.GOMAXPROCS(procs)
}
