package hypermap

import (
	"math/rand"
	"runtime"
	"sync"
	"testing"
)

const (
	testCores = 8
	testProcs = 10000
	testSize  = 35600000
)

func Benchmark_hypermap_parallel_write(b *testing.B) {
	procs := runtime.GOMAXPROCS(testCores)
	start := make(chan struct{}, testProcs)
	done := make(chan struct{}, testProcs)
	m := NewMap(40, testSize)
	keys := make([]uint64, b.N)
	values := make([]uint64, b.N)
	for i := 0; i < b.N; i++ {
		keys[i] = uint64(rand.Intn(1 << 40))
		values[i] = uint64(rand.Intn(1 << 24))
	}
	for k := 0; k < testProcs; k++ {
		btt := k
		go func() {
			<-start

			for i := uint64(0); i < uint64(b.N)/testProcs; i++ {
				j := uint64(btt)*(uint64(b.N)/testProcs) + i
				m.Set(keys[j], values[j])
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
	b.StopTimer()
	runtime.GOMAXPROCS(procs)
}

func Benchmark_map_parallel_write(b *testing.B) {
	procs := runtime.GOMAXPROCS(testCores)
	start := make(chan struct{}, testProcs)
	done := make(chan struct{}, testProcs)
	mutex := sync.RWMutex{}
	m := make(map[uint64]uint64, testSize)
	keys := make([]uint64, b.N)
	values := make([]uint64, b.N)
	for i := 0; i < b.N; i++ {
		keys[i] = uint64(rand.Intn(1 << 40))
		values[i] = uint64(rand.Intn(1 << 24))
	}
	for k := 0; k < testProcs; k++ {
		btt := k
		go func() {
			<-start

			for i := uint64(0); i < uint64(b.N)/testProcs; i++ {
				j := uint64(btt)*(uint64(b.N)/testProcs) + i
				mutex.Lock()
				m[keys[j]] = values[j]
				mutex.Unlock()
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
	b.StopTimer()
	runtime.GOMAXPROCS(procs)
}
