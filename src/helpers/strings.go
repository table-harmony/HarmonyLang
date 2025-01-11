package helpers

import "strings"

func Capitalize(input string) string {
	if input == "" {
		return ""
	}

	return strings.ToUpper(string(input[0])) + input[1:]
}

func ProcessEscapes(input string) string {
	// First replace literal `\\` with a placeholder to prevent it from being confused with `\n`
	placeholder := "\x00" // Temporary placeholder for `\\`
	input = strings.ReplaceAll(input, `\\`, placeholder)

	// Replace `\n`, `\t`, `\r` with their actual meanings
	input = strings.ReplaceAll(input, `\n`, "\n") // Newline
	input = strings.ReplaceAll(input, `\t`, "\t") // Tab
	input = strings.ReplaceAll(input, `\r`, "\r") // Carriage return

	// Restore literal `\\` from the placeholder
	input = strings.ReplaceAll(input, placeholder, `\`)

	return input
}
