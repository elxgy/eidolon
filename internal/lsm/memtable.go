package lsm

import (
	"bytes"
	"sort"
)

type entry struct {
	key   []byte
	value []byte
}

type MemTable struct {
	entries []entry
}

func NewMemTable() *MemTable {
	return &MemTable{
		entries: make([]entry, 0),
	}
}

func (m *MemTable) Put(newKey, newValue []byte) {

	i := sort.Search(len(m.entries), func(j int) bool {
		return bytes.Compare(m.entries[j].key, newKey) >= 0
	})

	if i < len(m.entries) && bytes.Equal(m.entries[i].key, newKey) {
		m.entries[i].value = newValue
		return
	}

	m.entries = append(m.entries, entry{})
	copy(m.entries[i+1:], m.entries[i:])
	m.entries[i] = entry{key: newKey, value: newValue}
}

func (m *MemTable) Get(key []byte) ([]byte, bool) {

	i := sort.Search(len(m.entries), func(j int) bool {
		return bytes.Compare(m.entries[j].key, key) >= 0
	})

	if i < len(m.entries) && bytes.Equal(m.entries[i].key, key) {
		return m.entries[i].value, true
	}

	return nil, false
}

func (m *MemTable) Len() int {
	return len(m.entries)
}
