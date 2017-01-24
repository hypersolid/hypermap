package main

import (
	"math/rand"
	"testing"
)

const defaultMapSize = 1024

func Test_Hash_of_string(t *testing.T) {
	m := NewMap(defaultMapSize ^ 2)
	if m.hashy("123") != m.hashy("123") {
		t.Error("ne equal")
	}
	if m.hashy("321") == m.hashy("123") {
		t.Error("equal")
	}
}

func Test_Hash_returns_different_values(t *testing.T) {
	m := NewMap(defaultMapSize ^ 2)
	for i := 0; i < 1000; i += 100 {
		for j := 0; j < 1000; j += 100 {
			if i == j {
				if m.hashy(i) != m.hashy(j) {
					t.Error(i, j, "ne equal")
				}
			} else {
				if m.hashy(i) == m.hashy(j) {
					t.Error(i, j, "equal")
				}
			}
		}
	}
}

func Test_Set_does_something(t *testing.T) {
	m := NewMap(defaultMapSize)
	v := rand.Int()
	swapped := m.Set(1, v)
	if !swapped {
		t.Error("not swapped")
	}
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
	m := NewMap(b.N)
	for i := 0; i < b.N; i++ {
		m.Set(i, "test")
	}
}

func Benchmark_map_write(b *testing.B) {
	m := make(map[int]string, b.N)
	for i := 0; i < b.N; i++ {
		m[i] = "test"
	}
}

func Benchmark_hashing(b *testing.B) {
	m := NewMap(defaultMapSize)
	for i := 0; i < b.N; i++ {
		m.hashy(i)
	}
}
