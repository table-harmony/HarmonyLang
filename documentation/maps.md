# Maps

Maps are key-value collections where each key must be unique. They provide efficient lookups and modifications of values based on their associated keys.

## Declaration

```
// Map type declaration
let prices: map[string -> number]

// Direct declaration and initialization
let scores = map[string -> number]{
    "foo" -> 42,
    "bar" -> 100
}
```

## Map Methods

#### get(key)

Retrieves the value associated with the specified key. Panics if the key doesn't exist.

```
let scores = map[string -> number]{"foo" -> 42}
let value = scores.get("foo")  // returns 42
let panic = scores.get("baz")  // panics: key doesn't exist
```

#### set(key, value)

Sets a value for the specified key. Creates a new entry if the key doesn't exist, or updates the existing value if it does.

```
let scores = map[string -> number]{"foo" -> 42}
scores.set("bar", 100)    // adds new entry
scores.set("foo", 50)     // updates existing entry
```

#### pop(key)

Removes and returns the value associated with the specified key.

```
let scores = map[string -> number]{
    "foo" -> 42,
    "bar" -> 100
}
let value = scores.pop("foo")  // removes entry and returns 42
```

#### intersect(other)

Creates a new map containing only the entries whose keys exist in both maps. Values from the other map take precedence.

```
let map1 = map[string -> number]{
    "foo" -> 42,
    "bar" -> 100
}
let map2 = map[string -> number]{
    "foo" -> 50,
    "baz" -> 200
}
let intersection = map1.intersect(map2)  // {"foo" -> 50}
```

#### union(other)

Creates a new map containing all entries from both maps. For keys that exist in both maps, values from the other map take precedence.

```
let map1 = map[string -> number]{
    "foo" -> 42,
    "bar" -> 100
}
let map2 = map[string -> number]{
    "foo" -> 50,
    "baz" -> 200
}
let combined = map1.union(map2)  // {"foo" -> 50, "bar" -> 100, "baz" -> 200}
```

#### keys()

Returns an array containing all keys in the map.

```
let scores = map[string -> number]{
    "foo" -> 42,
    "bar" -> 100
}
let keys = scores.keys()  // returns ["foo", "bar"]
```

#### values()

Returns an array containing all values in the map.

```
let scores = map[string -> number]{
    "foo" -> 42,
    "bar" -> 100
}
let values = scores.values()  // returns [42, 100]
```

#### exists(key)

Checks if a key exists in the map. Returns a boolean.

```
let scores = map[string -> number]{"foo" -> 42}
let exists = scores.exists("foo")  // returns true
let missing = scores.exists("bar") // returns false
```

## Computed Members

Maps support key-based access using square bracket notation:

```
let scores = map[string -> number]{"foo" -> 42}
scores["bar"] = 100           // Adds new entry
let value = scores["foo"]     // Retrieves value
```

## Common Patterns

### Map Iteration

```
let scores = map[string -> number]{
    "foo" -> 42,
    "bar" -> 100
}
for key, value in scores {
    // Process key and value
}
```

### Map Copying

```
let original = map[string -> number]{"foo" -> 42}
let copy = original          // Creates a copy
```

## Best Practices

1. Use descriptive key names
2. Check for key existence before using get() if panic is undesirable
3. Use set() instead of direct assignment for better clarity
4. Consider using exists() before operations on optional keys
5. Use intersect() and union() for set-like operations
6. Be mindful of key uniqueness

## Memory Considerations

1. Maps dynamically resize as needed
2. Copying maps creates new instances
3. Large maps may impact memory usage
4. Removed entries are garbage collected
