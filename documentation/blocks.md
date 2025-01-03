# Blocks

Blocks in this language are scoped sections of code that can appear anywhere in the program. They serve as both scope boundaries and expressions.

## Basic Syntax

A block is defined by curly braces:

```
{
    // Block contents
}
```

## Block as Expressions

Blocks are expressions and will evaluate to the last expression within them. If the last statement is not an expression, the block evaluates to `nil`.

### Examples

Block returning a value:

```
let foo = {
    let bar = 42
    bar + 1           // Block evaluates to 43
}
```

Block in conditional:

```
let result = if true {
    let foo = 42
    foo + 1          // Block evaluates to 43
} else {
    100
}
```

## Scope Rules

Blocks create their own scope. Variables declared within a block are only accessible within that block and any nested blocks.

### Variable Scope Examples

```
let foo = 1
{
    let bar = 2    // bar only accessible within this block
    foo = 3        // outer variables can be accessed and modified
}
// bar is not accessible here
```

### Nested Blocks

Blocks can be nested within other blocks:

```
{
    let foo = 1
    {
        let bar = 2
        // Both foo and bar are accessible here
    }
    // Only foo is accessible here
}
// Neither foo nor bar are accessible here
```

## Block Usage in Control Structures

Blocks are used in various control structures:

### If Expressions

```
const result = if foo == "bar" {
    42
} else if foo == "baz" {
    100
} else {
    200
}
```

### Switch Expressions

At switch expressions, there can only be one default case.
However there can exists multiple cases for the same expressions at that scenrio the first case is selected.

```
const result = switch foo {
    case "bar" {
        42
    }
    case "baz", "qux" {
        100
    }
    default {
        200
    }
}
```

## Try Blocks

Blocks can be used with error handling:

```
const result = try {
    42
} catch {
    100
}
```

## Best Practices

1. Use blocks to create clear scope boundaries
2. Be mindful of variable scope within blocks
3. Remember that blocks return their last expression value
4. Use appropriate indentation to show block structure
