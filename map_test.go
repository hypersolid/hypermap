package main

import (
	"regexp"
	"testing"
)

const (
	testRange = 10
)

func Test_NewMap_works(t *testing.T) {
	m := NewMap(20)
	if m == nil {
		t.Error("fail")
	}
}

func Test_String_works(t *testing.T) {
	m := NewMap(20)
	match, _ := regexp.MatchString(`HyperMap<\d+ 20 bits -> 43 bits>`, m.String())
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
	NewMap(15)
}

func Test_Map_panics_on_too_long_key(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Map should panic on short key")
		}
	}()
	NewMap(63)
}

func Test_Set_works(t *testing.T) {
	m := NewMap(20)
	for i := uint64(0); i < testRange; i++ {
		m.Set(i, i)
	}
	// for i := uint64(0); i < testRange; i++ {
	// 	value := (*m.array)[i] & m.valueMask
	// 	if value != i {
	// 		t.Errorf("Set error %d->%d", i, value)
	// 	}
	// }
}
