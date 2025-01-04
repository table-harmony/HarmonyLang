package interpreter

import (
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

type ModuleType struct {
}

func (m ModuleType) String() string { return "module" }
func (m ModuleType) Equals(other Type) bool {
	_, ok := other.(ModuleType)
	return ok
}
func (m ModuleType) DefaultValue() Value {
	return NewModule()
}

type Module struct {
	exports map[string]Value
}

func NewModule() *Module {
	return &Module{
		exports: make(map[string]Value),
	}
}

// Module implements the Value interface
func (m Module) Type() Type {
	return ModuleType{}
}
func (m Module) Clone() Value {
	module := NewModule()
	for key, value := range m.exports {
		module.exports[key] = value.Clone()
	}
	return module
}
func (m Module) String() string {
	str := "module {\n"
	for key, value := range m.exports {
		str += "  " + key + ": " + value.String() + "\n"
	}
	str += "}"
	return str
}

var standard_modules = map[string]Module{
	"math":   init_math_module(),
	"time":   init_time_module(),
	"random": init_random_module(),
	"os":     init_os_module(),
	//"http": TODO,
}

func init_math_module() Module {
	module := NewModule()

	// Mathematical constant π
	module.exports["PI"] = NewNumber(3.141592653589793)

	// Mathematical constant e
	module.exports["E"] = NewNumber(2.718281828459045)

	// Mathematical functions

	// abs(number): number
	// Purpose: Returns the absolute value of a number
	module.exports["abs"] = NewNativeFunction(func(args ...Value) Value {
		arg := args[0].(Number)
		return NewNumber(math.Abs(arg.Value()))
	}, []Type{PrimitiveType{NumberType}}, PrimitiveType{NumberType})

	// floor(number): number
	// Purpose: Returns the largest integer less than or equal to a number
	module.exports["floor"] = NewNativeFunction(func(args ...Value) Value {
		arg := args[0].(Number)
		return NewNumber(math.Floor(arg.Value()))
	}, []Type{PrimitiveType{NumberType}}, PrimitiveType{NumberType})

	// pow(base: number, exponent: number): number
	// Purpose: Returns base raised to the power of exponent
	module.exports["pow"] = NewNativeFunction(func(args ...Value) Value {
		base := args[0].(Number)
		exponent := args[1].(Number)
		return NewNumber(math.Pow(base.Value(), exponent.Value()))
	}, []Type{PrimitiveType{NumberType}, PrimitiveType{NumberType}}, PrimitiveType{NumberType})

	// sqrt(number): number
	// Purpose: Returns the square root of a number
	module.exports["sqrt"] = NewNativeFunction(func(args ...Value) Value {
		arg := args[0].(Number)
		return NewNumber(math.Sqrt(arg.Value()))
	}, []Type{PrimitiveType{NumberType}}, PrimitiveType{NumberType})

	// min(a: number, b: number): number
	// Purpose: Returns the smaller of two numbers
	module.exports["min"] = NewNativeFunction(func(args ...Value) Value {
		a := args[0].(Number)
		b := args[1].(Number)
		if a.Value() < b.Value() {
			return a
		}
		return b
	}, []Type{PrimitiveType{NumberType}, PrimitiveType{NumberType}}, PrimitiveType{NumberType})

	// max(a: number, b: number): number
	// Purpose: Returns the bigger of two numbers
	module.exports["max"] = NewNativeFunction(func(args ...Value) Value {
		a := args[0].(Number)
		b := args[1].(Number)
		if a.Value() > b.Value() {
			return a
		}
		return b
	}, []Type{PrimitiveType{NumberType}, PrimitiveType{NumberType}}, PrimitiveType{NumberType})

	// round(number): number
	// Purpose: Returns the nearest integer to a number
	module.exports["round"] = NewNativeFunction(func(args ...Value) Value {
		arg := args[0].(Number)
		return NewNumber(math.Round(arg.Value()))
	}, []Type{PrimitiveType{NumberType}}, PrimitiveType{NumberType})

	// sin(number): number
	// Purpose: Returns the sine of a number (angle in radians)
	module.exports["sin"] = NewNativeFunction(func(args ...Value) Value {
		arg := args[0].(Number)
		return NewNumber(math.Sin(arg.Value()))
	}, []Type{PrimitiveType{NumberType}}, PrimitiveType{NumberType})

	// cos(number): number
	// Purpose: Returns the tangent of a number (angle in radians)
	module.exports["cos"] = NewNativeFunction(func(args ...Value) Value {
		arg := args[0].(Number)
		return NewNumber(math.Cos(arg.Value()))
	}, []Type{PrimitiveType{NumberType}}, PrimitiveType{NumberType})

	// tan(number): number
	// Purpose: Returns the tangent of a number (angle in radians)
	module.exports["tan"] = NewNativeFunction(func(args ...Value) Value {
		arg := args[0].(Number)
		return NewNumber(math.Tan(arg.Value()))
	}, []Type{PrimitiveType{NumberType}}, PrimitiveType{NumberType})

	// clamp(num: number, min: number, max: number): number
	// Purpose: Constrains a number between minimum and maximum values
	module.exports["clamp"] = NewNativeFunction(func(args ...Value) Value {
		num := args[0].(Number)
		min := args[1].(Number)
		max := args[2].(Number)
		if num.Value() < min.Value() {
			return min
		}
		if num.Value() > max.Value() {
			return max
		}
		return num
	}, []Type{PrimitiveType{NumberType}, PrimitiveType{NumberType}, PrimitiveType{NumberType}}, PrimitiveType{NumberType})

	// sign(number): number
	// Purpose: Returns -1 for negative numbers, 0 for zero, and 1 for positive numbers
	module.exports["sign"] = NewNativeFunction(func(args ...Value) Value {
		arg := args[0].(Number)
		if arg.Value() > 0 {
			return NewNumber(1)
		} else if arg.Value() < 0 {
			return NewNumber(-1)
		}
		return NewNumber(0)
	}, []Type{PrimitiveType{NumberType}}, PrimitiveType{NumberType})

	// ln(number): number
	// Purpose: Returns the natural logarithm of a number
	module.exports["ln"] = NewNativeFunction(func(args ...Value) Value {
		arg := args[0].(Number)
		return NewNumber(math.Log(arg.Value()))
	}, []Type{PrimitiveType{NumberType}}, PrimitiveType{NumberType})

	// log2(number): number
	// Purpose: Returns the base-2 logarithm of a number
	module.exports["log2"] = NewNativeFunction(func(args ...Value) Value {
		arg := args[0].(Number)
		return NewNumber(math.Log2(arg.Value()))
	}, []Type{PrimitiveType{NumberType}}, PrimitiveType{NumberType})

	// log10(number): number
	// Purpose: Returns the base-10 logarithm of a number
	module.exports["log10"] = NewNativeFunction(func(args ...Value) Value {
		arg := args[0].(Number)
		return NewNumber(math.Log10(arg.Value()))
	}, []Type{PrimitiveType{NumberType}}, PrimitiveType{NumberType})

	return *module
}

func init_time_module() Module {
	module := NewModule()

	// now(): string
	// Purpose: Returns the current time in RFC3339 format
	module.exports["now"] = NewNativeFunction(func(args ...Value) Value {
		return NewString(time.Now().Format(time.RFC3339))
	}, []Type{}, PrimitiveType{StringType})

	//timestamp(): number
	// Purpose: Returns the current Unix timestamp in seconds
	module.exports["timestamp"] = NewNativeFunction(func(args ...Value) Value {
		return NewNumber(float64(time.Now().Unix()))
	}, []Type{}, PrimitiveType{NumberType})

	// sleep(milliseconds: number): nil
	// Purpose: Pauses execution for the specified number of milliseconds
	module.exports["sleep"] = NewNativeFunction(func(args ...Value) Value {
		duration := args[0].(Number)
		time.Sleep(time.Duration(duration.Value()) * time.Millisecond)
		return NewNil()
	}, []Type{PrimitiveType{NumberType}}, PrimitiveType{NilType})

	// format(timestamp: number, layout: string): string
	// Purpose: Formats a Unix timestamp according to the specified layout
	module.exports["format"] = NewNativeFunction(func(args ...Value) Value {
		timestamp := args[0].(Number)
		layout := args[1].(String)

		t := time.Unix(int64(timestamp.Value()), 0)
		return NewString(t.Format(layout.Value()))
	}, []Type{PrimitiveType{NumberType}, PrimitiveType{StringType}}, PrimitiveType{StringType})

	// parse(layout: string, timeString: string): number
	// Purpose: Parses a time string according to the layout and returns Unix timestamp
	module.exports["parse"] = NewNativeFunction(func(args ...Value) Value {
		layout := args[0].(String)
		timeString := args[1].(String)

		t, err := time.Parse(layout.Value(), timeString.Value())
		if err != nil {
			panic("Invalid time format")
		}

		return NewNumber(float64(t.Unix()))
	}, []Type{PrimitiveType{StringType}, PrimitiveType{StringType}}, PrimitiveType{NumberType})

	module.exports["add"] = NewNativeFunction(func(args ...Value) Value {
		timestamp := args[0].(Number)
		seconds := args[1].(Number)

		t := time.Unix(int64(timestamp.Value()), 0).Add(time.Duration(seconds.Value()) * time.Second)
		return NewNumber(float64(t.Unix()))
	}, []Type{PrimitiveType{NumberType}, PrimitiveType{NumberType}}, PrimitiveType{NumberType})

	// subtract(timestamp: number, seconds: number): number
	// Purpose: Subtracts seconds from a Unix timestamp
	module.exports["subtract"] = NewNativeFunction(func(args ...Value) Value {
		timestamp := args[0].(Number)
		seconds := args[1].(Number)

		t := time.Unix(int64(timestamp.Value()), 0).Add(-time.Duration(seconds.Value()) * time.Second)
		return NewNumber(float64(t.Unix()))
	}, []Type{PrimitiveType{NumberType}, PrimitiveType{NumberType}}, PrimitiveType{NumberType})

	// difference(start: number, end: number): number
	// Purpose: Returns the difference in seconds between two timestamps
	module.exports["difference"] = NewNativeFunction(func(args ...Value) Value {
		start := args[0].(Number)
		end := args[1].(Number)

		diff := end.Value() - start.Value()
		return NewNumber(diff)
	}, []Type{PrimitiveType{NumberType}, PrimitiveType{NumberType}}, PrimitiveType{NumberType})

	// day(): number
	// Purpose: Returns the current day of the month
	module.exports["day"] = NewNativeFunction(func(args ...Value) Value {
		return NewNumber(float64(time.Now().Day()))
	}, []Type{}, PrimitiveType{NumberType})

	// month(): number
	// Purpose: Returns the current month of the year
	module.exports["month"] = NewNativeFunction(func(args ...Value) Value {
		return NewNumber(float64(time.Now().Month()))
	}, []Type{}, PrimitiveType{NumberType})

	// year(): number
	// Purpose: Returns the current year
	module.exports["year"] = NewNativeFunction(func(args ...Value) Value {
		return NewNumber(float64(time.Now().Year()))
	}, []Type{}, PrimitiveType{NumberType})

	// is_leap_year(year: number): boolean
	// Purpose: Determines if the given year is a leap year
	module.exports["is_leap_year"] = NewNativeFunction(func(args ...Value) Value {
		year := args[0].(Number)

		isLeap := year.Value() > 0 && (int(year.Value())%4 == 0 && (int(year.Value())%100 != 0 || int(year.Value())%400 == 0))
		return NewBoolean(isLeap)
	}, []Type{PrimitiveType{NumberType}}, PrimitiveType{BooleanType})

	return *module
}

func init_random_module() Module {
	rand.Seed(time.Now().UnixNano())

	module := NewModule()

	// int(min: number, max: number): number
	// Purpose: Returns a random integer between min (inclusive) and max (exclusive)
	module.exports["int"] = NewNativeFunction(
		func(args ...Value) Value {
			min := int(args[0].(Number).Value())
			max := int(args[1].(Number).Value())
			return NewNumber(float64(rand.Intn(max-min) + min))
		},
		[]Type{PrimitiveType{NumberType}, PrimitiveType{NumberType}},
		PrimitiveType{NumberType},
	)

	// float(): number
	// Purpose: Returns a random floating-point number between 0.0 and 1.0
	module.exports["float"] = NewNativeFunction(
		func(args ...Value) Value {
			return NewNumber(rand.Float64())
		},
		[]Type{},
		PrimitiveType{NumberType},
	)

	// bool(): boolean
	// Purpose: Returns a random boolean value
	module.exports["bool"] = NewNativeFunction(
		func(args ...Value) Value {
			return NewBoolean(rand.Float64() < 0.5)
		},
		[]Type{},
		PrimitiveType{BooleanType},
	)

	// string(): string
	// Purpose: Returns a random string value
	module.exports["string"] = NewNativeFunction(
		func(args ...Value) Value {
			const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
			b := make([]byte, int(args[0].(Number).Value()))
			for i := range b {
				b[i] = letterBytes[rand.Intn(len(letterBytes))]
			}
			return NewString(string(b))
		},
		[]Type{PrimitiveType{NumberType}},
		PrimitiveType{StringType},
	)

	// shuffle(slice: []any): []any
	// Purpose: Randomly shuffles the elements of a slice
	module.exports["shuffle"] = NewNativeFunction(
		func(args ...Value) Value {
			slice := args[0].(Slice)
			elements := *slice.elements
			rand.Shuffle(len(elements), func(i, j int) {
				elements[i], elements[j] = elements[j], elements[i]
			})
			return NewSlice(elements, slice._type.elementType)
		},
		[]Type{NewSliceType(PrimitiveType{AnyType})},
		NewSliceType(PrimitiveType{AnyType}),
	)

	// choice(slice: []any): any
	// Purpose: Returns a random element from a non-empty slice
	module.exports["choice"] = NewNativeFunction(
		func(args ...Value) Value {
			slice := args[0].(Slice)
			elements := *slice.elements
			if len(elements) == 0 {
				panic("Cannot choose from empty slice")
			}
			return elements[rand.Intn(len(elements))]
		},
		[]Type{NewSliceType(PrimitiveType{AnyType})},
		PrimitiveType{AnyType},
	)

	return *module
}

func init_os_module() Module {
	module := NewModule()

	// Environment variables

	// getenv(key: string): string
	// Purpose: Gets the value of an environment variable
	module.exports["getenv"] = NewNativeFunction(
		func(args ...Value) Value {
			key := args[0].(String).Value()
			return NewString(os.Getenv(key))
		},
		[]Type{PrimitiveType{StringType}},
		PrimitiveType{StringType},
	)

	//setenv(key: string, value: string): nil
	// Purpose: Sets the value of an environment variable
	module.exports["setenv"] = NewNativeFunction(
		func(args ...Value) Value {
			key := args[0].(String).Value()
			value := args[1].(String).Value()
			err := os.Setenv(key, value)
			if err != nil {
				panic(err.Error())
			}
			return NewNil()
		},
		[]Type{PrimitiveType{StringType}, PrimitiveType{StringType}},
		PrimitiveType{NilType},
	)

	// File system operations

	// read_file(path: string): string
	// Purpose: Reads the entire contents of a file as a string
	module.exports["read_file"] = NewNativeFunction(
		func(args ...Value) Value {
			path := args[0].(String).Value()
			data, err := os.ReadFile(path)
			if err != nil {
				panic(err.Error())
			}
			return NewString(string(data))
		},
		[]Type{PrimitiveType{StringType}},
		PrimitiveType{StringType},
	)

	// write_file(path: string, data: string): nil
	// Purpose: Writes data to a file, creating it if it doesn't exist
	module.exports["write_file"] = NewNativeFunction(
		func(args ...Value) Value {
			path := args[0].(String).Value()
			data := args[1].(String).Value()
			err := os.WriteFile(path, []byte(data), 0644)
			if err != nil {
				panic(err.Error())
			}
			return NewNil()
		},
		[]Type{PrimitiveType{StringType}, PrimitiveType{StringType}},
		PrimitiveType{NilType},
	)

	// remove(path: string): nil
	// Purpose: Removes a file or empty directory
	module.exports["remove"] = NewNativeFunction(
		func(args ...Value) Value {
			path := args[0].(String).Value()
			err := os.Remove(path)
			if err != nil {
				panic(err.Error())
			}
			return NewNil()
		},
		[]Type{PrimitiveType{StringType}},
		PrimitiveType{NilType},
	)

	// mkdir(path: string): nil
	// Purpose: Creates a directory and any necessary parent directories
	module.exports["mkdir"] = NewNativeFunction(
		func(args ...Value) Value {
			path := args[0].(String).Value()
			err := os.MkdirAll(path, 0755)
			if err != nil {
				panic(err.Error())
			}
			return NewNil()
		},
		[]Type{PrimitiveType{StringType}},
		PrimitiveType{NilType},
	)

	//list_dir(path: string): []string
	// Purpose: Lists the contents of a directory
	module.exports["list_dir"] = NewNativeFunction(
		func(args ...Value) Value {
			path := args[0].(String).Value()
			entries, err := os.ReadDir(path)
			if err != nil {
				panic(err.Error())
			}

			fileNames := make([]Value, 0)
			for _, entry := range entries {
				fileNames = append(fileNames, NewString(entry.Name()))
			}

			return NewSlice(fileNames, PrimitiveType{StringType})
		},
		[]Type{PrimitiveType{StringType}},
		NewSliceType(PrimitiveType{StringType}),
	)

	// abs_path(path: string): string
	// Purpose: Returns the absolute path for a given path
	module.exports["abs_path"] = NewNativeFunction(
		func(args ...Value) Value {
			path := args[0].(String).Value()
			absPath, err := filepath.Abs(path)
			if err != nil {
				panic(err.Error())
			}
			return NewString(absPath)
		},
		[]Type{PrimitiveType{StringType}},
		PrimitiveType{StringType},
	)

	// exit(code: number): nil
	// Purpose: Terminates the current process with the specified exit code
	module.exports["exit"] = NewNativeFunction(
		func(args ...Value) Value {
			code := int(args[0].(Number).Value())
			os.Exit(code)
			return NewNil()
		},
		[]Type{PrimitiveType{NumberType}},
		PrimitiveType{NilType},
	)

	return *module
}
