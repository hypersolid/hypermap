package main

import "testing"

func Test_set_does_something(t *testing.T) {
	m := NewMap(10)
	swapped, err := m.Set(1, 101)
	if !swapped {
		t.Error("not swapped")
	}
	if err != nil {
		t.Error("err", err)
	}
}

func Test_get_does_something(t *testing.T) {
	m := NewMap(10)
	m.Set(1, 101)
	result := m.Get(1).(int)
	if result != 101 {
		t.Error("got:", result)
	}
}

func Benchmark_hypermap_write(b *testing.B) {
	m := NewMap(b.N)
	for i := 0; i < b.N; i++ {
		m.Set(i, i)
	}
}

func Benchmark_map_write(b *testing.B) {
	m := make(map[int]int, b.N)
	for i := 0; i < b.N; i++ {
		m[i] = i
	}
}

func Benchmark_hashing(b *testing.B) {
	m := NewMap(1)
	for i := 0; i < b.N; i++ {
		m.Hash(i)
	}
}
