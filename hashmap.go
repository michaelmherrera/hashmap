package hashmap

import (
	"encoding/binary"
	"errors"
	"fmt"
	"hash/fnv"
)

const defaultSize = 8
const loadFactor = 2

type keyValPair struct {
	key int
	val string
}

type node struct {
	key  int
	val  string
	next *node
}

// Hashmap mapping integer keys to string values
type Hashmap struct {
	buckets []node
	size    uint64
	entries uint64
}

// New Hashmap with default size of 8 and load factor of 2
func New() Hashmap {
	return Hashmap{make([]node, defaultSize), defaultSize, 0}
}

// NewWithSize Hashmap creates a hashmap of the specified initial size and a load factor of 2
func NewWithSize(size uint64) Hashmap {
	return Hashmap{make([]node, size), size, 0}
}

func (m *Hashmap) iterate(kvp chan keyValPair) {
	var idx uint64 = 0
	for idx = 0; idx < m.size; idx++ {
		curr := m.buckets[idx]
		for curr.next != nil {
			kvp <- keyValPair{curr.key, curr.val}
			curr = *curr.next
		}
	}
	close(kvp)
}

// rehash the hashmap
func (m *Hashmap) rehash() {
	newMap := NewWithSize(m.size << 1)
	kvps := make(chan keyValPair)
	go m.iterate(kvps)
	for kvp := range kvps {
		newMap.Put(kvp.key, kvp.val)
	}
	*m = newMap
}

// ToString returns a string representation of the Hashmap
func (m *Hashmap) ToString() string {
	s := "{"
	kvps := make(chan keyValPair)
	go m.iterate(kvps)
	first := true
	for kvp := range kvps {
		entry := ""
		key := kvp.key
		val := kvp.val
		if first {
			entry = fmt.Sprintf("%d: \"%s\"", key, val)
			first = false
		} else {
			entry = fmt.Sprintf(", %d: \"%s\"", key, val)
		}
		s += entry
		first = false
	}
	s += "}"
	return s
}

// Put a key-value pair into the hashmap
func (m *Hashmap) Put(key int, val string) {
	if m.size*loadFactor <= m.entries {
		m.rehash()
	}

	idx := hash(key) % m.size
	var curr *node = &m.buckets[idx]
	// If the key doesn't exist yet, traverse
	// the list until it reaches the end.
	// If the key does exist, traverse until
	// it reaches the key
	for curr.next != nil && curr.key != key {
		curr = curr.next
	}
	if curr.next == nil {
		// If it reached the end of the list (didn't encounter the key)
		curr.key = key
		curr.next = &node{0, "", nil}
	}
	curr.val = val
	m.entries++
}

// Get a value given a key, if the value for that key exists
func (m *Hashmap) Get(key int) (val string, found bool) {
	idx := hash(key) % m.size
	curr := m.buckets[idx]
	for curr.key != key {
		// If the key doesn't exist in the hashmap
		if curr.next == nil {
			return "", false
		}
		curr = *curr.next
	}
	return curr.val, true
}

// Remove the key-value pair associated with key. Return an error if the key doesn't exist
func (m *Hashmap) Remove(key int) (err error) {
	idx := hash(key) % m.size
	var curr *node = &m.buckets[idx]
	for curr.key != key && curr.next != nil {
		// The second check is necessary to prevent it from considering 0 as
		// a valid key in cases that it is not
		curr = curr.next
	}
	if curr.key != key || (curr.key == 0 && curr.next == nil) {
		// The second check is necessary to prevent it from considering 0 as
		// a valid key in cases that it is not
		msg := fmt.Sprintf("hashmap: the key %d does not exist", key)
		return errors.New(msg)
	}
	if curr.next != nil {
		/*
			If it's not the last element in the list, remove the node.
			The typical stategy of node removal is not possible, though,
			because the address of the head of the list cannot change
			due to it being part of an array. Thus, instead, we shift
			the next node to the current position and delete the next node.
		*/
		var next *node = curr.next
		curr.key = next.key
		curr.val = next.val
		curr.next = next.next
	}
	m.entries--
	return nil

}

func hash(key int) uint64 {
	binKey := make([]byte, 8)
	binary.LittleEndian.PutUint64(binKey, uint64(key))

	h := fnv.New64()
	h.Write([]byte(binKey))
	return uint64(h.Sum64())
}
