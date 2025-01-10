# Structs

Structs are user-defined types that group together related data and methods. They support both instance and static members, and can contain fields of any type including functions.

## Declaration

### Basic Syntax

```
struct Rectangle {
    width: number
    height: number
}
```

### With Methods and Static Members

```
struct Rectangle {
    // Static field
    static counter = 0

    // Instance fields
    width: number
    height: number

    // Instance method
    fn get_area() -> number {
        return self.width * self.height
    }

    // Static method
    static fn increment_counter() {
        Rectangle.counter++
    }
}
```

## Fields

### Field Types

Fields can be declared with explicit types or initialized with values:

```
struct Point {
    x: number           // Explicit type
    const y = 0        // Initialized constant
    z = 42             // Initialized with type inference
}
```

### Constant Fields

Fields can be declared as constants using the `const` keyword:

```
struct Config {
    const version = "1.0"
    const max_retries = 3
}
```

### Function Fields

Fields can be function types:

```
struct Handler {
    callback: fn() -> void
    process: fn(data: string) -> boolean
}
```

## Methods

### Instance Methods

Instance methods have access to struct fields through the `self` keyword:

```
struct Circle {
    radius: number

    fn area() -> number {
        return 3.14159 * self.radius * self.radius
    }

    fn set_radius(r: number) {
        self.radius = r
    }
}
```

### Static Methods

Static methods belong to the struct type itself and don't have access to instance fields:

```
struct MathHelper {
    static fn square(x: number) -> number {
        return x * x
    }
}
```

## Instantiation

### Creating Instances

```
const rect = Rectangle{
    width: 10,
    height: 20
}
```

### Field Access

Fields can be accessed using dot notation:

```
rect.width = 30
let area = rect.get_area()
```

### Computed Member Access

Fields can also be accessed using square bracket notation:

```
rect["width"] = 40
let height = rect["height"]
```

## Static Members

### Accessing Static Fields

Static fields are accessed through the struct type:

```
Rectangle.counter++
let count = Rectangle["counter"]
```

### Calling Static Methods

```
Rectangle.increment_counter()
MathHelper.square(5)
```

## Nested Structs

Structs can contain other struct instances:

```
struct OuterStruct {
    inner: InnerStruct
}

struct InnerStruct {
    value: number
}

const outer = OuterStruct{
    inner: InnerStruct{
        value: 42
    }
}
```

## Best Practices

1. Use meaningful names for structs and their fields
2. Group related data and functionality together
3. Make fields constant when they shouldn't change
4. Use static members for functionality that doesn't depend on instance state
5. Consider using methods instead of function fields for better encapsulation
6. Document complex struct behavior and requirements

## Common Patterns

### Builder Pattern

```
struct HttpRequest {
    method: string
    url: string
    headers: map[string -> string]

    fn with_header(key: string, value: string) -> *HttpRequest {
        self.headers.set(key, value)
        return self
    }
}
```

### Factory Methods

```
struct Database {
    connection: string

    static fn new(host: string, port: number) -> Database {
        return Database{
            connection: host + ":" + port
        }
    }
}
```

## Memory Considerations

1. Structs are passed by value unless referenced
2. Large structs should be passed by reference to avoid copying
3. Static fields are shared across all instances
4. Instance fields are unique to each instance
5. Methods don't consume additional memory per instance
