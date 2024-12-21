package main

import (
	"fmt"
	"os"
)

func main() {
	bytes, err := os.ReadFile("examples/00.harmony")

	if err != nil {
		panic(err)
	}

	source := string(bytes)

	fmt.Printf("Code: %s\n", source)
}
