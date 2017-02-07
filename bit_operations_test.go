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

func Test_deleted_works(t *testing.T) {
	m := NewMap(20)
	(*m.array)[1] |= m.deletedMask
	if m.deleted(0) {
		t.Errorf("bucket #0 must not be deleted")
	}
	if !m.deleted(1) {
		t.Errorf("bucket #1 must be deleted")
	}
}

func Test_key_works(t *testing.T) {
	m := NewMap(20)
	(*m.array)[0] = uint64(79212312312321321)
	if bitsToString(m.key(0)) != "0000000000000000000000000000000000000000000000000010001100101101" {
		t.Errorf("error key is %s", bitsToString(m.key(0)))
	}
}

func Test_value_works(t *testing.T) {
	m := NewMap(20)
	(*m.array)[0] = uint64(79212312312321321)
	if bitsToString(m.value(0)) != "0000000000000000000000110010110110101001001101101001010100101001" {
		t.Errorf("error value is %s", bitsToString(m.value(0)))
	}
}
