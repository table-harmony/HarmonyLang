# References and Values

In this language, everything is passed by value by default. To pass by reference, you must explicitly use the reference operator (`&`) and dereference operator (`*`).

## Basic Types

### Pass by Value

Basic types are always passed by value unless explicitly referenced:

```
fn modify(foo: number) {
    foo = 42        // Only modifies local copy
}

let bar = 100
modify(bar)         // bar is still 100
```

### Pass by Reference

To modify the original value, use references:

```
fn modify(foo: *number) {
    *foo = 42       // Modifies original value
}

let bar = 100
modify(&bar)        // bar is now 42
```

## Arrays

### Pass by Value

Arrays are copied when passed:

```
fn modify(arr: [3]number) {
    arr[0] = 42     // Only modifies local copy
}

let foo = [3]number{100, 200, 300}
modify(foo)         // foo[0] is still 100
```

### Pass by Reference

To modify the original array:

```
fn modify(arr: *[3]number) {
    (*arr)[0] = 42  // Modifies original array
}

let foo = [3]number{100, 200, 300}
modify(&foo)        // foo[0] is now 42
```

## Slices

### Pass by Value

Even though slices contain a reference to the underlying array, the slice struct itself is copied:

```
fn modify(arr: []number) {
    arr = []number{42, 100}    // Only modifies local copy
}

let foo = []number{100, 200, 300}
modify(foo)                    // foo is unchanged
```

### Pass by Reference

To modify the slice itself:

```
fn modify(arr: *[]number) {
    *arr = []number{42, 100}   // Modifies original slice
}

let foo = []number{100, 200, 300}
modify(&foo)                   // foo is now [42, 100]
```

## Maps

### Pass by Value

Maps are copied when passed:

```
fn modify(dict: map[string -> number]) {
    dict["foo"] = 42           // Only modifies local copy
}

let bar = map[string -> number]{
    "foo" -> 100,
    "bar" -> 200
}
modify(bar)                    // bar["foo"] is still 100
```

### Pass by Reference

To modify the original map:

```
fn modify(dict: *map[string -> number]) {
    (*dict)["foo"] = 42        // Modifies original map
}

let bar = map[string -> number]{
    "foo" -> 100,
    "bar" -> 200
}
modify(&bar)                   // bar["foo"] is now 42
```

## Common Patterns

### Swapping Values

```
fn swap(a: *number, b: *number) {
    let temp = *a
    *a = *b
    *b = temp
}

let foo = 42
let bar = 100
swap(&foo, &bar)              // foo is 100, bar is 42
```

## Best Practices

1. Use pass by value by default
2. Use references only when you need to modify the original value
3. Document functions that take references clearly
4. Be careful with multiple levels of references
5. Consider using return values instead of reference parameters when possible

## Common Pitfalls

1. Forgetting to dereference when modifying values
2. Forgetting to use address-of operator when passing references
3. Creating references to temporary values
4. Not properly handling nil references
5. Creating circular references

## Safety Considerations

1. References cannot be nil unless explicitly typed as optional
2. Dereferencing a nil reference causes a panic
3. References cannot outlive their referred values
4. The language prevents dangling references
