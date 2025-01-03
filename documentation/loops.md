# Loops

The language supports several types of loop constructs, including C-style loops, infinite loops, expression loops, and iteration over data structures.

## C-Style For Loops

The traditional C-style for loop consists of three components: initialization, condition, and iteration expression.

### Basic Syntax

```
for let i = 0; i < condition; i++ {
    // loop body
}
```

### Examples

```
// Loop with array length
let sum = 0
for let i = 0; i < arr.len(); i++ {
    sum += arr[i]
}
```

## Expression Loop

A loop that continues while an expression evaluates to true.

### Syntax

```
let foo = true
for foo {
    // loop body
}
```

## Infinite Loop

A loop that runs indefinitely until explicitly broken.

### Syntax

```
for {
    // loop body
}
```

## Data Structure Iteration

### Array/Slice Iteration

Iterate over arrays or slices with index and value:

```
let arr = []number{42, 100, 200}
for index, value in arr {
    // loop body
}
```

### Map Iteration

Iterate over maps with key and value:

```
const foo = map[string -> number]{
    "bar" -> 42,
    "baz" -> 100,
}

for key, value in foo {
    // loop body
}
```

### Range Iteration

Iterate over numeric ranges:

```
for index, value in 1..10 {
    // loop body
}
```

## Loop Components

### Loop Body

- Contained within curly braces `{}`
- Can contain any valid statements or expressions
- Can be a single line or multiple lines

### Iteration Variables

- Can be used to track position (index)
- Can access current value in iteration
- Scope is limited to the loop body

## Common Patterns

### Accumulation

```
let sum = 0
for let i = 0; i < arr.len(); i++ {
    sum += arr[i]
}
```

### Array Modification

```
for let i = 0; i < arr.len(); i++ {
    arr[i] = arr[i] * 2
}
```

### Concurrent Key-Value Access

```
for key, value in foo {
    // Process both key and value
}
```

## Best Practices

1. Choose the appropriate loop type for your use case
2. Use meaningful variable names for indices and values
3. Keep loop bodies focused and concise
4. Consider using data structure iteration instead of C-style loops when possible
5. Be careful with infinite loops - ensure there's a way to exit
