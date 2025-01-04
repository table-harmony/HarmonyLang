# Arrays and Slices

## Arrays

Arrays are fixed-size collections of elements of the same type.

### Declaration

```
// Array type declaration
let nums: [3]number

// Direct declaration and initialization
let arr = [3]number{42, 100, 200}
```

### Array Methods

#### len()

Returns the length of the array.

```
let arr = [3]number{42, 100, 200}
let length = arr.len()  // returns 3
```

#### get(index)

Retrieves the value at the specified index.

```
let arr = [3]number{42, 100, 200}
let value = arr.get(1)  // returns 100
```

#### set(index, value)

Sets a new value at the specified index.

```
let arr = [3]number{42, 100, 200}
arr.set(1, 300)       // arr is now [42, 300, 200]
```

#### each(fn)

Executes a function for each element in the array.

```
let arr = [42, 100, 200]
arr.each(fn(value) {
    // Process each value
    return value % 2 == 0
})
```

#### filter(fn)

Creates a new array with elements that pass the test function.

```
let arr = [42, 100, 200]
let filtered = arr
    .filter(fn(value) {
        return !value
    })
```

### Computed Members

Arrays support index-based access using square bracket notation:

```
let arr = [3]number{42, 100, 200}
arr[2] = 300           // Modifies third element
let value = arr[1]     // Retrieves second element
```

## Slices

Slices are dynamic-length views into arrays.

### Creation

#### From Arrays

```
const arr = [42, 100, 200, 300]
const slice = arr.slice(1, 3)  // Creates slice from index 1 to 3
```

#### Using Range Notation

```
// Create slice from range
const foo = 2..10        // Creates slice with numbers 2 through 10
const bar = 10..3..-2    // Creates slice with decreasing numbers
const baz = 2..10..2     // Creates slice with step of 2
```

### Slice Methods

#### append()

Adds elements to the end of the slice.

```
const slice: []number
slice = []number{42, 100, 200}
slice.append(300)        // slice is now [42, 100, 200, 300]
```

#### cap()

Returns the capacity of the slice (maximum length before reallocation).

```
const arr = []number{42, 100, 200}
const capacity = arr.cap()
```

#### len()

Returns the current length of the slice.

```
const arr = []number{42, 100, 200}
const length = arr.len()
```

## Common Patterns

### Array/Slice Iteration

```
const arr = []number{42, 100, 200}
for index, value in arr {
    // Process index and value
}
```

### Chaining Methods

```
const result = arr
    .each(fn(value) {
        return value % 2 == 0
    })
    .filter(fn(value) {
        return !value
    })
```

### Copying

```
const arr = [3]number{42, 100, 200}
const copy = arr             // Creates a copy
```

## Best Practices

1. Use arrays when the size is fixed and known
2. Use slices for dynamic collections
3. Check array/slice bounds before accessing elements
4. Use provided methods instead of direct index access when possible
5. Consider using ranges for numeric sequences
6. Be mindful of slice capacity when frequently appending

## Memory Considerations

1. Arrays are fixed size and allocated on declaration
2. Slices may grow and reallocate when capacity is exceeded
3. Copying arrays/slices creates new instances
4. Range operations create new slices
