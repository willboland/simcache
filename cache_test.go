package simcache

import (
	"strconv"
	"testing"
	"time"
)

type pair[T any] struct {
	key   string
	value T
}

func makePairs[T any](length int) []pair[T] {
	var pairs []pair[T]
	for i := 0; i < length; i++ {
		p := pair[T]{key: strconv.Itoa(i), value: *new(T)}
		pairs = append(pairs, p)
	}
	return pairs
}

func TestNew(t *testing.T) {
	c := New[int](time.Second)
	if c == nil {
		t.Fatal("cache should not be empty when calling New")
	}
}

func TestCache_Add(t *testing.T) {
	type unitTest struct {
		name     string
		c        *Cache[int]
		pairs    []pair[int]
		expected bool
	}

	tests := []unitTest{
		{
			name:     "One Item",
			c:        New[int](time.Hour),
			pairs:    makePairs[int](1),
			expected: true,
		},
		{
			name:     "Two Items",
			c:        New[int](time.Hour),
			pairs:    makePairs[int](2),
			expected: true,
		},
		{
			name:     "Five Items",
			c:        New[int](time.Hour),
			pairs:    makePairs[int](5),
			expected: true,
		},
	}

	for _, test := range tests {
		for _, p := range test.pairs {
			actual := test.c.Add(p.key, p.value)
			if test.expected != actual {
				t.Fatalf("%s FAILED - expected %t but got %t", test.name, test.expected, actual)
			}
		}
		expectedLength := len(test.pairs)
		length := len(test.c.Keys())
		if expectedLength != length {
			t.Fatalf("%s FAILED - expected %d but got %d", test.name, expectedLength, length)
		}
	}
}

func TestCache_Set(t *testing.T) {
	c := New[int](time.Hour)
	c.Set("a", 1)
	c.Set("a", 2)
	a, f := c.Get("a")
	if !f || a != 2 {
		t.Fatalf(`FAILED to update value key "a" from Set`)
	}

	type unitTest struct {
		name  string
		c     *Cache[int]
		pairs []pair[int]
	}

	tests := []unitTest{
		{
			name:  "One Item",
			c:     New[int](time.Hour),
			pairs: makePairs[int](1),
		},
		{
			name:  "Two Items",
			c:     New[int](time.Hour),
			pairs: makePairs[int](2),
		},
		{
			name:  "Five Items",
			c:     New[int](time.Hour),
			pairs: makePairs[int](5),
		},
	}

	for _, test := range tests {
		for _, p := range test.pairs {
			test.c.Set(p.key, p.value)
			val, found := test.c.Get(p.key)
			if !found {
				t.Fatalf("%s FAILED - expected %t but got %t", test.name, true, found)
			}
			if p.value != val {
				t.Fatalf("%s FAILED - expected %d but got %d", test.name, p.value, val)
			}
		}
	}
}

func TestCache_Get(t *testing.T) {
	c := New[int](time.Hour)
	_, f := c.Get("a")
	if f {
		t.Fatalf("FAILED - found item when no items were added to cache")
	}

	type unitTest struct {
		name  string
		c     *Cache[int]
		pairs []pair[int]
	}

	tests := []unitTest{
		{
			name:  "One Item",
			c:     New[int](time.Hour),
			pairs: makePairs[int](1),
		},
		{
			name:  "Two Items",
			c:     New[int](time.Hour),
			pairs: makePairs[int](2),
		},
		{
			name:  "Five Items",
			c:     New[int](time.Hour),
			pairs: makePairs[int](5),
		},
	}

	for _, test := range tests {
		for _, p := range test.pairs {
			_ = test.c.Add(p.key, p.value)
			val, found := test.c.Get(p.key)
			if !found {
				t.Fatalf("%s FAILED - expected %t but got %t", test.name, true, found)
			}
			if p.value != val {
				t.Fatalf("%s FAILED - expected %d but got %d", test.name, p.value, val)
			}
		}
	}
}

func TestCache_Delete(t *testing.T) {
	type unitTest struct {
		name  string
		c     *Cache[int]
		pairs []pair[int]
	}

	tests := []unitTest{
		{
			name:  "One Item",
			c:     New[int](time.Hour),
			pairs: makePairs[int](1),
		},
		{
			name:  "Two Items",
			c:     New[int](time.Hour),
			pairs: makePairs[int](2),
		},
		{
			name:  "Five Items",
			c:     New[int](time.Hour),
			pairs: makePairs[int](5),
		},
	}

	for _, test := range tests {
		for _, p := range test.pairs {
			_ = test.c.Add(p.key, p.value)
		}
		for _, p := range test.pairs {
			test.c.Delete(p.key)
			_, found := test.c.Get(p.key)
			if found {
				t.Fatalf("%s FAILED - expected %t but got %t", test.name, false, found)
			}
		}
	}
}

func TestCache_Items(t *testing.T) {
	type unitTest struct {
		name  string
		c     *Cache[int]
		pairs []pair[int]
	}

	tests := []unitTest{
		{
			name:  "One Item",
			c:     New[int](time.Hour),
			pairs: makePairs[int](1),
		},
		{
			name:  "Two Items",
			c:     New[int](time.Hour),
			pairs: makePairs[int](2),
		},
		{
			name:  "Five Items",
			c:     New[int](time.Hour),
			pairs: makePairs[int](5),
		},
	}

	for _, test := range tests {
		for _, p := range test.pairs {
			_ = test.c.Add(p.key, p.value)
		}
		items := test.c.Items()
		for _, p := range test.pairs {
			val, found := items[p.key]
			if !found {
				t.Fatalf("%s FAILED - expected %t but got %t", test.name, true, found)
			}
			if p.value != val {
				t.Fatalf("%s FAILED - expected %d but got %d", test.name, p.value, val)
			}
		}
	}
}

func TestCache_Keys(t *testing.T) {
	type unitTest struct {
		name  string
		c     *Cache[int]
		pairs []pair[int]
	}

	tests := []unitTest{
		{
			name:  "One Item",
			c:     New[int](time.Hour),
			pairs: makePairs[int](1),
		},
		{
			name:  "Two Items",
			c:     New[int](time.Hour),
			pairs: makePairs[int](2),
		},
		{
			name:  "Five Items",
			c:     New[int](time.Hour),
			pairs: makePairs[int](5),
		},
	}

	for _, test := range tests {
		for _, p := range test.pairs {
			_ = test.c.Add(p.key, p.value)
		}
		keys := test.c.Keys()
		for _, p := range test.pairs {
			found := contains(p.key, keys)
			if !found {
				t.Fatalf("%s FAILED - expected %t but got %t", test.name, true, found)
			}
		}
	}
}

func TestCache_Values(t *testing.T) {
	type unitTest struct {
		name  string
		c     *Cache[int]
		pairs []pair[int]
	}

	tests := []unitTest{
		{
			name:  "One Item",
			c:     New[int](time.Hour),
			pairs: makePairs[int](1),
		},
		{
			name:  "Two Items",
			c:     New[int](time.Hour),
			pairs: makePairs[int](2),
		},
		{
			name:  "Five Items",
			c:     New[int](time.Hour),
			pairs: makePairs[int](5),
		},
	}

	for _, test := range tests {
		for _, p := range test.pairs {
			_ = test.c.Add(p.key, p.value)
		}
		values := test.c.Values()
		for _, p := range test.pairs {
			found := contains(p.value, values)
			if !found {
				t.Fatalf("%s FAILED - expected %t but got %t", test.name, true, found)
			}
		}
	}
}

func TestCache_Purge(t *testing.T) {
	c := New[int](time.Hour)
	c.Set("one", 1, time.Nanosecond)
	c.Set("two", 2)
	time.Sleep(time.Nanosecond * 2)

	count := c.Purge()
	if count != 1 {
		t.Fatal("incorrect number of items deleted: ", count)
	}
	_, found := c.Get("one")
	if found {
		t.Fatal(`"one" was found in cache when it should have been deleted`)
	}
}

func contains[T comparable](target T, s []T) bool {
	for _, actual := range s {
		if actual == target {
			return true
		}
	}
	return false
}
