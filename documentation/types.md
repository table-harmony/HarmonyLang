# Types

Types in this language can be declared and aliased using the `type` keyword. This allows for creation of custom types, type aliases, and function types.

## Basic Type Aliases

### Syntax

```
type TypeName BaseType
```

### Examples

```
type Foo number          // Number type alias
type Bar string          // String type alias
type Baz bool           // Boolean type alias

let foo: Foo = 42
let bar: Bar = "qux"
let baz: Baz = true
```

## Function Types

### Basic Function Type

```
type Handler fn() -> any

let foo: Handler = fn() {
    return 42
}
```

### Function Type with Parameters

```
type Operation fn(foo: number, bar: number) -> number

let add: Operation = fn(x: number, y: number) -> number {
    return x + y
}

let result = add(42, 100)
```

### Function Type with Multiple Returns

```
type Parser fn(foo: string) -> (number, bool)

let parse: Parser = fn(input: string) -> (number, bool) {
    return (42, true)
}
```

## Type Usage

### Variable Declaration

```
type ID number
type Status string

let foo: ID = 42
let bar: Status = "active"
```

### Function Parameters

```
type Filter fn(foo: number) -> bool

fn process(callback: Filter) {
    let result = callback(42)
}
```

### Arrays with Custom Types

```
type Score number

let foo: [3]Score = [3]Score{42, 100, 200}
let bar: []Score = []Score{42, 100, 200}
```

### Maps with Custom Types

```
type UserID number
type UserName string

let foo: map[UserID -> UserName] = map[UserID -> UserName]{
    1 -> "bar",
    2 -> "baz"
}
```

## Best Practices

1. Use type aliases to provide semantic meaning
2. Create function types for common function signatures
3. Use descriptive names for custom types
4. Keep struct definitions focused and cohesive
5. Consider using type aliases for better code readability

## Common Patterns

### Function Type Callbacks

```
type Predicate fn(foo: number) -> bool
type Transform fn(foo: number) -> number

fn process(value: number, pred: Predicate, trans: Transform) {
    if pred(value) {
        return trans(value)
    }
    return value
}
```

### Generic Container Types

```
type Collection []number
type Dictionary map[string -> any]

let foo: Collection = []number{42, 100}
let bar: Dictionary = map[string -> any]{
    "baz" -> 42,
    "qux" -> "value"
}
```

## Type Safety

1. Type aliases create distinct types
2. Type checking is enforced at compile time
3. Type conversions must be explicit when required
4. Function types must match exactly in signature
5. Struct field types are strictly enforced
