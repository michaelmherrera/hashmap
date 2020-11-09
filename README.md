# Go Hashmap

> An implementation of an `<int, string>` hashmap in go.

## Functions

```go
// New Hashmap with default size of 8 and load factor of 2
func New() Hashmap {}

// NewWithSize Hashmap creates a hashmap of the specified initial size and a load factor of 2
func NewWithSize(size uint64) Hashmap {}

// Put a key-value pair into the hashmap
func (m *Hashmap) Put(key int, val string) {}

// Get a value given a key, if the value for that key exists
func (m *Hashmap) Get(key int) (val string, found bool) {}

// Remove the key-value pair associated with key. Return an error if the key doesn't exist
func (m *Hashmap) Remove(key int) (err error) {}

// ToString returns a string representation of the Hashmap, formatted similar to a Python dictionary
func (m *Hashmap) ToString() string {}
```
