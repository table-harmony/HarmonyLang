# Functions

Functions in this language can be declared in several ways, including regular functions, anonymous functions, and functions with type annotations.

## Function Declaration

### Basic Function Syntax

```
fn foo(bar) {
    // function body
    return bar
}
```

### Type Annotations

Functions can have type annotations for parameters and return values:

```
fn foo(bar: number) -> number {
    return bar
}
```

## Function Types

### Function Type Declaration

Function types can be declared using the `fn` keyword with parameter and return types:

```
type handler fn() -> any

// Function variable with type annotation
let foo: fn(bar: number) -> any = nil
```

### Function Type Assignment

```
foo = fn(bar: number) -> number {
    return 42
}
```

## Anonymous Functions

Anonymous functions can be created using the `fn` keyword:

```
const foo = fn(bar: number) -> number {
    return bar
}

const baz = fn() {
}
```

## Higher-Order Functions

Functions that return other functions (closures):

```
const generator = fn(foo: number) {
    return fn(bar: number) { return foo + bar }
}

const add42 = generator(42)
add42(100)
```

## Recursive Functions

Functions can call themselves:

```
fn sum(foo) {
    if foo == 0 {
        return 0
    }
    return recursive(foo - 1) + foo
}
```

## Return Values

- Functions automatically return the value of their last expression if it's not a statement
- The `return` keyword can be used to explicitly return a value
- If no return value is specified and the last line is a statement, `nil` is returned

## Best Practices

1. Use clear and descriptive function names
2. Include type annotations for better code clarity
3. Keep functions focused on a single responsibility
4. Document complex function parameters and return values
5. Consider using named functions instead of anonymous functions for reusability
