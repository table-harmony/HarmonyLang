# Variables

Variables in this language are containers for storing and managing data. They come in two forms: mutable (`let`) and immutable (`const`).

## Declaration

### Basic Syntax

Variables can be declared using either `let` or `const` keywords:

```
let foo = value
const bar = value
```

### Type Annotations

Type annotations are optional but can be explicitly specified:

```
let foo: number = 42
const bar: []number = [42, 100, 200]
```

### Rules

- `let` declarations require either a type annotation or initial value:

  ```
  let foo: number       // Valid
  let foo = 42          // Valid
  let foo               // Invalid: must specify type or value
  ```

- `const` declarations must always have an initial value:
  ```
  const foo = 42        // Valid
  const foo: number     // Invalid: const requires initial value
  ```

## Mutability

### Mutable Variables (`let`)

- Can be reassigned after declaration
- Value can be modified
- Commonly used for values that need to change during program execution

```
let foo = 0
foo = 42            // Valid
foo += 100          // Valid
```

### Immutable Variables (`const`)

- Cannot be reassigned after declaration
- Value is fixed for the entire lifetime
- Used for values that should remain constant

```
const foo = 42
foo = 100          // Invalid: cannot reassign const
```

## Operators and Assignments

### Arithmetic Assignment Operators

```
let foo = 42
foo += 10              // Addition assignment (foo = foo + 10)
foo -= 5               // Subtraction assignment (foo = foo - 5)
foo *= 2               // Multiplication assignment (foo = foo * 2)
foo /= 4               // Division assignment (foo = foo / 4)
foo %= 3               // Modulo assignment (foo = foo % 3)
```

### Increment and Decrement

```
let foo = 0
foo++                 // Increment by 1
foo--                 // Decrement by 1
```

### Null Coalescing Assignment

```
foo ??= bar          // Assigns bar to foo if foo is null
```

## Scope

Variables follow block scope rules and are only accessible within their declaring block and nested blocks:

```
if true {
    let foo = 42     // foo only accessible within this block
    const bar = 100  // bar only accessible within this block
}
// foo and bar not accessible here
```

## Type Inference

The language supports automatic type inference based on the assigned value:

```
let foo = "bar"    // Inferred as string
let bar = 42       // Inferred as number
let baz = true     // Inferred as boolean
```

## Best Practices

1. Use `const` by default, and only use `let` when you need to reassign the variable
2. Include type annotations for complex types or when type inference might be ambiguous
3. Choose descriptive variable names that indicate their purpose
4. Initialize variables with meaningful values when possible
5. Keep variable scope as narrow as possible

## Common Pitfalls to Avoid

1. Forgetting to initialize variables
2. Attempting to reassign `const` variables
3. Declaring variables without either type or initial value
4. Accessing variables outside their scope

## Examples

### Basic Usage

```
let counter = 0
const maxAttempts = 3

counter++
if counter >= maxAttempts {
    // handle max attempts reached
}
```

### Scope and Shadowing

```
let foo = 42
{
    let foo = 100    // Different variable, shadows outer foo
    // foo is 100 here
}
// foo is 42 here
```
