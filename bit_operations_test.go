package hypermap

import (
	"testing"
)

func Test_bitsToString_works(t *testing.T) {
	m := test_create_map()
	if bitsToString(m.keyMask) != "1111111111111111111100000000000000000000000000000000000000000000" {
		t.Errorf("deletedMask should not be %s", bitsToString(m.keyMask))
	}
	if bitsToString(m.valueMask) != "0000000000000000000011111111111111111111111111111111111111111111" {
		t.Errorf("valueMask should not be %s", bitsToString(m.valueMask))
	}
}

func Test_fuse_works(t *testing.T) {
	m := test_create_map()
	result := m.fuse(99123, 978)
	if bitsToString(result) != "0001100000110011001100000000000000000000000000000000001111010010" {
		t.Errorf("fuse error %s", bitsToString(result))
	}
}

func Test_available_works(t *testing.T) {
	m := test_create_map()

	m.array[0] = m.fuse(1, 2)
	m.array[1] = m.fuse(1, 2) | m.valueMask
	m.array[2] = m.keyMask
	if ok := m.available(0); ok {
		t.Errorf("bucket 0 should not be available")
	}
	if ok := m.available(1); !ok {
		t.Errorf("bucket 1 should be available")
	}
	if ok := m.available(2); !ok {
		t.Errorf("bucket 2 should be available")
	}
}

func Test_deleted_works(t *testing.T) {
	m := test_create_map()
	m.array[1] |= m.valueMask
	if m.deleted(0) {
		t.Errorf("bucket #0 must not be deleted")
	}
	if !m.deleted(1) {
		t.Errorf("bucket #1 must be deleted")
	}
}

func Test_key_works(t *testing.T) {
	m := test_create_map()
	m.array[0] = uint64(3 << m.valueSize)
	if bitsToString(m.key(0)) != "0000000000000000000000000000000000000000000000000000000000000011" {
		t.Errorf("error key is %s", bitsToString(m.key(0)))
	}
}

func Test_value_works(t *testing.T) {
	m := test_create_map()
	m.array[0] = uint64(3)
	if bitsToString(m.value(0)) != "0000000000000000000000000000000000000000000000000000000000000011" {
		t.Errorf("error value is %s", bitsToString(m.value(0)))
	}
}
