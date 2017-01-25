package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"testing"
)

const defaultMapSize = 1024

func Test_Hash_of_string(t *testing.T) {
	m := NewMap(defaultMapSize ^ 2)
	if m.hashy("123") != m.hashy("123") {
		t.Error("'123' is not equal to '123'")
	}
	if m.hashy("321") == m.hashy("123") {
		t.Error("'321' is equal to '123'")
	}
}

func Test_Hash_returns_different_values(t *testing.T) {
	v1 := rand.Int()
	v2 := rand.Int()
	m := NewMap(defaultMapSize ^ 2)
	if m.hashy(v1) != m.hashy(v1) {
		t.Error(v1, v1, "ne equal")
	}
	if m.hashy(v1) == m.hashy(v2) {
		t.Error(v1, v2, "equal")
	}
}

func Test_Set_does_something(t *testing.T) {
	m := NewMap(defaultMapSize ^ 2)
	v := rand.Int()
	swapped := m.Set(1, v)
	if !swapped {
		t.Error("not swapped")
	}
}

func Test_Set_parallel(t *testing.T) {
	runtime.GOMAXPROCS(4)
	m := NewMap(defaultMapSize)
	var wg sync.WaitGroup
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func() {
			for k := 0; k < 10000; k++ {
				v := rand.Intn(defaultMapSize)
				swapped := m.Set(v, v)
				if !swapped {
					t.Error("not swapped retries:", m.retries, "collisions", m.collisions)
				}
			}
			wg.Done()
		}()
	}
	wg.Wait()
	fmt.Println(m.retries)
	fmt.Println(m.collisions)
}

func Test_Get_does_something(t *testing.T) {
	m := NewMap(defaultMapSize)
	v := rand.Int()
	m.Set(1, v)
	result := m.Get(1).(int)
	if result != v {
		t.Error("got:", result)
	}
}

// benchmarks

func Benchmark_hypermap_write(b *testing.B) {
	b.StopTimer()
	m := NewMap(b.N * 2)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		m.Set(i, "test")
	}
}

func Benchmark_map_write(b *testing.B) {
	b.StopTimer()
	m := make(map[int]string, b.N*2)
	mtx := new(sync.RWMutex)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		mtx.Lock()
		m[i] = "test"
		mtx.Unlock()
	}
}

func Benchmark_hypermap_read(b *testing.B) {
	b.StopTimer()
	m := NewMap(b.N * 2)
	for i := 0; i < b.N; i++ {
		m.Set(i, i)
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		result := m.Get(i)
		if result == nil || result.(int) != i {
			b.Error(i, result, m.load, m.collisions, m.retries, "busted")
		}
	}
}

func Benchmark_map_read(b *testing.B) {
	b.StopTimer()
	m := make(map[int]int, b.N*2)
	for i := 0; i < b.N; i++ {
		m[i] = i
	}
	mtx := new(sync.RWMutex)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		mtx.RLock()
		if m[i] != i {
			b.Error(i, "busted")
		}
		mtx.RUnlock()
	}
}

func Benchmark_hashing(b *testing.B) {
	m := NewMap(defaultMapSize)
	for i := 0; i < b.N; i++ {
		m.hashy(i)
	}
}
