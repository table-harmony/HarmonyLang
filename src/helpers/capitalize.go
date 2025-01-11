package helpers

import "strings"

func Capitalize(input string) string {
	if input == "" {
		return ""
	}

	return strings.ToUpper(string(input[0])) + input[1:]
}
