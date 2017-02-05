package main

import "testing"

func Test_bitsToString_works(t *testing.T) {
	m := NewMap(20)
	if bitsToString(m.deletedMask) != "1000000000000000000000000000000000000000000000000000000000000000" {
		t.Errorf("deletedMask should not be %s", bitsToString(m.deletedMask))
	}
	if bitsToString(m.keyMask) != "0111111111111111111110000000000000000000000000000000000000000000" {
		t.Errorf("deletedMask should not be %s", bitsToString(m.keyMask))
	}
	if bitsToString(m.valueMask) != "0000000000000000000001111111111111111111111111111111111111111111" {
		t.Errorf("deletedMask should not be %s", bitsToString(m.valueMask))
	}
}

func Test_fuse_panics_on_too_big_key(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Map should panic on big key passed to fuse")
		}
	}()
	m := NewMap(16)
	m.fuse(66000, 1)
}

func Test_fuse_panics_on_too_big_value(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Map should panic on big value passed to fuse")
		}
	}()
	m := NewMap(60)
	m.fuse(1, 10)
}

func Test_fuse_works(t *testing.T) {
	m := NewMap(20)
	result := m.fuse(99123, 978)
	if bitsToString(result) != "0000110000011001100110000000000000000000000000000000001111010010" {
		t.Errorf("fuse error %s", bitsToString(result))
	}
}

func Test_available_works(t *testing.T) {
	m := NewMap(20)
	(*m.array)[0] = m.fuse(1, 2)
	(*m.array)[1] = m.fuse(1, 2) | m.deletedMask

	if _, ok := m.available(0); ok {
		t.Errorf("bucket 0 should not be available")
	}
	if _, ok := m.available(1); !ok {
		t.Errorf("bucket 1 should be available")
	}
	if _, ok := m.available(2); !ok {
		t.Errorf("bucket 2 should be available")
	}
}
