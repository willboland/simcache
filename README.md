# SimCache
A simple cache for Go using generics.

SimCache allows the creation of a thread-safe cache for any type, without the need to cast the returned value.
The cache **clears expired items upon any retrieval operation**, such as `Get` or `Items`.

This is essentially https://github.com/patrickmn/go-cache but with generics. Was made because of [go-cache #149](https://github.com/patrickmn/go-cache/pull/149).

## Installation
Assuming that [you have Go installed](https://go.dev/doc/install):  
`go get -u github.com/willboland/simcache`

## Examples
SimCache allows you to pull values from the cache without the need to cast them to a specific type.
It works with primitive types:
```go
func main() {
    // Cache holds int's with a default TTL of one minute
    cache := simcache.New[int](time.Minute)
	
    cache.Set("one", 1)
    cache.Set("two", 2)
    one, _ := cache.Get("one")
    two, _ := cache.Get("two")
	
    fmt.Print(one) // 1
    fmt.Print(two) // 2
}
```
It also works with user-defined types:
```go
type Message struct {
    Author    string
    Content   string
    CreatedAt time.Time
}

func main() {
    m := Message{
        Author: "Will Boland",
        Content: "SimCache is easy to use",
    }
	
    cache := simcache.New[Message](time.Hour)
    cache.Set("wb", m)
    message, _ := cache.Get("wb")
	
    fmt.Print(message.Author)  // Will Boland
    fmt.Print(message.Content) // SimCache is easy to use
}
```

### Adding an item - `Set`
An item can be added or updated using the `Set` method. This method accepts a key/value pair, with an optional TTL.
If the TTL is not specified, the default one set when the cache was created is used. If more than one TTL is specified, it uses the first given.
```go
cache := simcache.New[int](time.Minute)
cache.Set("one", 1) // TTL is one minute
cache.Set("two", 2, time.Hour) // TTL is one hour
cache.Set("three", 3, time.Second, time.Hour) // TTL is one second
```

### Getting an item - `Get`
An item can be retrieved using the `Get` method. It returns the value in the cache for a given key and if it was found. 
If no such key exists, the returned bool will be false.

It is recommended to check if the value was found using the returned bool before accessing the value.
```go
cache := New[int](time.Minute)
cache.Set("one", 1)

one, found := cache.Get("one")

fmt.Print(found) // true
fmt.Print(one)   // 1
```

### Removing an item - `Delete`
The `Delete` method removes the item for the given key from the cache.
```go
cache.Delete("one") // Removes the item that had key "one" from cache
```

### Getting all key-value pairs - `Items`
All key-value pairs in the cache can be retrieved using the `Items` method. It returns a map of values that hold type T.
```go
cache := New[int](time.Minute)
cache.Set("one", 1)
cache.Set("two", 2)

items := cache.Items() // map[string]int{"one": 1, "two": 2}
```

### Getting all keys - `Keys`
The `Keys` method returns all keys in the cache.
```go
cache.Keys() // []string{"key1", "key2"}
```

### Getting all values - `Values`
The `Values` method returns all values in the cache.
```go
cache := New[int](time.Minute)
cache.Set("one", 1)
cache.Set("two", 2)

items := cache.Values() // []int{1, 2}
```

### Deleting all expired items - `Purge`
Items can be deleted from the cache before the next retrieval operation by calling the `Purge` method.
It returns the number of items deleted from the cache.
```go
cache.Set("one", 1) // Expired
cache.Set("two", 2) // Expired
cache.Purge() // 2
```
