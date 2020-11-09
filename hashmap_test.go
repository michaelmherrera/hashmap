package hashmap

import (
	"math/rand"
	"testing"
)

const randSeed = 42
const testIterations = 1 << 20 //Approx 1 million
const keyRange = 1 << 8

func randomWriteInRange(gomap *map[int]string, m *Hashmap, r int) {
	key := rand.Intn(r)
	val := getRandomString()
	(*m).Put(key, val)
	(*gomap)[key] = val
}

func randomWrite(gomap *map[int]string, m *Hashmap) {
	key := rand.Int()
	val := getRandomString()
	(*m).Put(key, val)
	(*gomap)[key] = val
}

func randomRemoveInRange(gomap *map[int]string, m *Hashmap, r int, t *testing.T) {
	key := rand.Intn(r)
	gomapVal, found := (*gomap)[key]
	mVal, _ := (*m).Get(key)
	err := (*m).Remove(key)
	if found && (err != nil) {
		t.Errorf("Expected to find a value in m, but found none. Key: %d, GomapVal: %s", key, gomapVal)
	} else if (!found) && (err == nil) {
		t.Errorf("Expected to find no value, but found one. Key: %d, mVal: %s", key, mVal)
	}
	delete(*gomap, key)
}

func checkAgainstGomap(gomap *map[int]string, m *Hashmap, t *testing.T) {
	for key, gomapVal := range *gomap {
		val, _ := (*m).Get(key)
		if gomapVal != val {
			t.Errorf("Key: %d CorrectVal: %s, ActualVal: %s\n", key, gomapVal, val)
		}
	}
}

func TestWriteDelete(t *testing.T) {
	gomap := make(map[int]string)
	m := New()
	rand.Seed(randSeed)
	for i := 0; i < testIterations; i++ {
		randomWriteInRange(&gomap, &m, keyRange)
		randomRemoveInRange(&gomap, &m, keyRange, t)
	}
	checkAgainstGomap(&gomap, &m, t)
}

func getRandomString() string {
	byteLen := rand.Intn(20)
	bytes := make([]byte, byteLen)
	for i := 0; i < byteLen; i++ {
		bytes[i] = byte(32 + rand.Intn(94))
	}
	return string(bytes)
}

func TestOverwrite(t *testing.T) {
	gomap := make(map[int]string)
	m := New()
	rand.Seed(randSeed)
	for i := 0; i < testIterations; i++ {
		randomWriteInRange(&gomap, &m, keyRange)
	}
	checkAgainstGomap(&gomap, &m, t)
}

func TestHighLoad(t *testing.T) {
	gomap := make(map[int]string)
	m := New()
	rand.Seed(randSeed)
	for i := 0; i < testIterations; i++ {
		randomWrite(&gomap, &m)
	}
	checkAgainstGomap(&gomap, &m, t)
}
