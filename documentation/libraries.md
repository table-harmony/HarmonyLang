# Standard Libraries Documentation

## Math Library

The math library provides common mathematical operations and constants.

### Constants

- `PI`: Mathematical constant π (3.141592653589793)
- `E`: Mathematical constant e (2.718281828459045)

### Functions

```go
abs(x: number) -> number
Purpose: Returns absolute value of x
Example: abs(-42) // Returns 42
```

```go
floor(x: number) -> number
Purpose: Returns largest integer less than or equal to x
Example: floor(3.7) // Returns 3
```

```go
pow(base: number, exponent: number) -> number
Purpose: Returns base raised to power of exponent
Example: pow(2, 3) // Returns 8
```

```go
sqrt(x: number) -> number
Purpose: Returns square root of x
Example: sqrt(16) // Returns 4
```

```go
min(a: number, b: number) -> number
Purpose: Returns smaller of two numbers
Example: min(2, 5) // Returns 2
```

```go
max(a: number, b: number) -> number
Purpose: Returns larger of two numbers
Example: max(2, 5) // Returns 5
```

```go
round(x: number) -> number
Purpose: Rounds number to nearest integer
Example: round(3.7) // Returns 4
```

```go
sin(x: number) -> number
Purpose: Returns sine of angle x (in radians)
Example: sin(PI/2) // Returns 1
```

```go
cos(x: number) -> number
Purpose: Returns cosine of angle x (in radians)
Example: cos(PI) // Returns -1
```

```go
tan(x: number) -> number
Purpose: Returns tangent of angle x (in radians)
Example: tan(PI/4) // Returns 1
```

```go
clamp(num: number, min: number, max: number) -> number
Purpose: Constrains number between minimum and maximum values
Example: clamp(15, 0, 10) // Returns 10
```

```go
sign(x: number) -> number
Purpose: Returns sign of number (-1, 0, or 1)
Example: sign(-42) // Returns -1
```

```go
ln(x: number) -> number
Purpose: Returns natural logarithm of x
Example: ln(E) // Returns 1
```

```go
log2(x: number) -> number
Purpose: Returns base-2 logarithm of x
Example: log2(8) // Returns 3
```

```go
log10(x: number) -> number
Purpose: Returns base-10 logarithm of x
Example: log10(100) // Returns 2
```

## Time Library

The time library provides functionality for working with dates, times, and durations.

### Functions

```go
now() -> string
Purpose: Returns current time in RFC3339 format
Example: now() // Returns "2024-01-04T15:04:05Z07:00"
```

```go
timestamp() -> number
Purpose: Returns current Unix timestamp in seconds
Example: timestamp() // Returns 1704384245
```

```go
sleep(milliseconds: number) -> nil
Purpose: Pauses execution for specified milliseconds
Example: sleep(1000) // Sleeps for 1 second
```

```go
format(timestamp: number, layout: string) -> string
Purpose: Formats Unix timestamp according to layout
Example: format(1704384245, "2006-01-02") // Returns "2024-01-04"
```

```go
parse(layout: string, timeString: string) -> number
Purpose: Parses time string to Unix timestamp
Example: parse("2006-01-02", "2024-01-04")
```

```go
add(timestamp: number, seconds: number) -> number
Purpose: Adds seconds to Unix timestamp
Example: add(timestamp(), 3600) // Adds 1 hour
```

```go
subtract(timestamp: number, seconds: number) -> number
Purpose: Subtracts seconds from Unix timestamp
Example: subtract(timestamp(), 3600) // Subtracts 1 hour
```

```go
difference(start: number, end: number) -> number
Purpose: Returns difference in seconds between timestamps
Example: difference(t1, t2)
```

```go
day() -> number
Purpose: Returns current day of month (1-31)
Example: day() // Returns current day
```

```go
month() -> number
Purpose: Returns current month (1-12)
Example: month() // Returns current month
```

```go
year() -> number
Purpose: Returns current year
Example: year() // Returns current year
```

```go
is_leap_year(year: number) -> boolean
Purpose: Checks if given year is a leap year
Example: is_leap_year(2024) // Returns true
```

## Random Library

The random library provides various random number generation functions.

### Functions

```go
int(min: number, max: number) -> number
Purpose: Returns random integer between min (inclusive) and max (exclusive)
Example: int(1, 10) // Returns random number between 1 and 9
```

```go
float() -> number
Purpose: Returns random float between 0.0 and 1.0
Example: float() // Returns random decimal like 0.7315
```

```go
bool() -> boolean
Purpose: Returns random boolean value
Example: bool() // Returns true or false randomly
```

```go
string(length: number) -> string
Purpose: Returns random boolean value
Example: string(10) // Returns a random string with the given length
```

```go
shuffle(slice: []any) -> []any
Purpose: Randomly shuffles elements in slice
Example: shuffle([1, 2, 3, 4, 5])
```

```go
choice(slice: []any) -> any
Purpose: Returns random element from non-empty slice
Example: choice(["a", "b", "c"])
```

## OS Library

The OS library provides operating system functionality for file and environment operations.

### File System Operations

```go
read_file(path: string) -> string
Purpose: Reads entire file content as string
Example: read_file("data.txt")
```

```go
write_file(path: string, data: string) -> nil
Purpose: Writes string data to file
Example: write_file("output.txt", "Hello World")
```

```go
remove(path: string) -> nil
Purpose: Removes file or empty directory
Example: remove("old_file.txt")
```

```go
mkdir(path: string) -> nil
Purpose: Creates directory and any necessary parent directories
Example: mkdir("new/folder/path")
```

```go
list_dir(path: string) -> []string
Purpose: Lists contents of directory
Example: list_dir("./documents")
```

```go
abs_path(path: string) -> string
Purpose: Returns absolute path for given path
Example: abs_path("./relative/path")
```

### Environment Operations

```go
getenv(key: string) -> string
Purpose: Gets value of environment variable
Example: getenv("HOME")
```

```go
setenv(key: string, value: string) -> nil
Purpose: Sets value of environment variable
Example: setenv("DEBUG", "true")
```

### Process Operations

```go
pid() -> number
Purpose: Returns current process ID
Example: pid()
```

```go
exit(code: number) -> nil
Purpose: Terminates process with specified exit code
Example: exit(1)
```

### Best Practices

1. Math Library

   - Check for domain errors (e.g., sqrt of negative numbers)
   - Handle potential infinity/NaN results
   - Use appropriate precision for calculations

2. Time Library

   - Use appropriate time formats for your use case
   - Consider timezone implications
   - Handle invalid time strings in parse function

3. Random Library

   - Don't use for cryptographic purposes
   - Seed random number generator appropriately
   - Validate input ranges for int function

4. OS Library
   - Handle file operation errors appropriately
   - Check file permissions before operations
   - Use absolute paths when possible
   - Clean up resources after use
