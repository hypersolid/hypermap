package main

func (m *Map) hashy(key uint64) uint64 {
	key += ^(key << 32)
	key ^= (key >> 22)
	key += ^(key << 13)
	key ^= (key >> 8)
	key += (key << 3)
	key ^= (key >> 15)
	key += ^(key << 27)
	key ^= (key >> 31)
	return key
}
