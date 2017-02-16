package hypermap

import (
	"regexp"
	"testing"
)

const (
	testRange        = 10
	testMapSmallSize = 256
)

func test_create_map() *Map {
	return NewMap(20, testMapSmallSize)
}

func Test_NewMap_works(t *testing.T) {
	m := test_create_map()
	if m == nil {
		t.Error("fail")
	}
}

func Test_String_works(t *testing.T) {
	m := test_create_map()
	match, _ := regexp.MatchString(`HyperMap<20 bits -> 44 bits>`, m.String())
	if !match {
		t.Errorf("String() should not be %s", m.String())
	}
}

func Test_Map_panics_on_too_short_key(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Map should panic on short key")
		}
	}()
	NewMap(15, testMapSmallSize)
}

func Test_Map_panics_on_too_long_key(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Map should panic on long key")
		}
	}()
	NewMap(64, testMapSmallSize)
}

func Test_Set_works(t *testing.T) {
	m := test_create_map()
	for i := uint64(0); i < testRange; i++ {
		m.Set(i, i)
	}
}

func Test_Delete_works(t *testing.T) {
	m := test_create_map()
	for i := uint64(0); i < testRange; i++ {
		m.Set(i, i)
	}

	m.Delete(2)

	if _, ok := m.Get(2); ok {
		t.Errorf("Key wasn't deleted")
	}
}

func Test_Set_over_capacity_causes_grow(t *testing.T) {
	m := test_create_map()

	for i := uint64(0); i < testMapSmallSize+1; i++ {
		m.Set(i, i)
	}
	if m.size <= testMapSmallSize {
		t.Errorf("Map should grow when it's overflown")
	}
}

func Test_Get_works(t *testing.T) {
	m := test_create_map()
	for i := uint64(0); i < testRange; i++ {
		m.Set(i, i)
	}
	for i := uint64(0); i < testRange; i++ {
		v, ok := m.Get(i)
		if !ok {
			t.Errorf("Get didn't found the value on step %d", i)
		}
		if v != i {
			t.Errorf("Get returned %d at step %d", v, i)
		}
	}
}

func Test_Grow_works(t *testing.T) {
	m := test_create_map()
	for i := uint64(0); i < testRange; i++ {
		m.Set(i, i)
	}
	m.grow(256)
	for i := uint64(0); i < testRange; i++ {
		v, ok := m.Get(i)
		if !ok {
			t.Errorf("[grow] Get didn't found the value on step %d", i)
		}
		if v != i {
			t.Errorf("[grow] Get returned %d at step %d", v, i)
		}
	}
}

func Test_Get_fails_to_find_value(t *testing.T) {
	m := test_create_map()
	v, ok := m.Get(44)
	if ok {
		t.Errorf("Get should not find anything")
	}
	if v != 0 {
		t.Errorf("empty result should be 0")
	}
}

func Test_Get_probes_array_until_finds_value(t *testing.T) {
	m := test_create_map()
	key := uint64(44)
	value := uint64(177)
	bucket := m.hashy(key) % m.size
	m.Set(key, value)

	m.array[(bucket+1)%m.size] = m.array[bucket]
	m.array[bucket] |= m.valueMask

	v, ok := m.Get(key)
	if !ok {
		t.Errorf("Get should probe until finds the value")
	}
	if v != value {
		t.Errorf("Get probing returned wrong value")
	}
}
