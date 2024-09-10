package utils

import (
	"errors"
	"strings"
	"unicode"
)

func CapitalizeFirstLetter(s string) string {
	if s == "" {
		return s
	}
	return strings.ToUpper(string(s[0])) + s[1:]
}

func StartsWithLowerCase(s string) bool {
	if len(s) == 0 {
		return false // or true, depending on your definition for an empty string
	}
	return unicode.IsLower(rune(s[0]))
}

func IsValidIdentifier(s string) bool {
	// Iterate over each character in the string
	for _, r := range s {
		// Check if the character is not a letter, digit, or underscore
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '_' {
			return false
		}
	}
	return true
}

func isUnderscope(r rune) bool {
	return r == '_'
}

func ContainsUnderscore(s string) bool {
	for _, r := range s {
		if isUnderscope(r) {
			return true
		}
	}
	return false
}

func RemoveUnderscore(s string) string {
	return strings.ReplaceAll(s, string("_"), "")
}

func GetLetter(index int) (string, error) {
	if index < 0 || index > 25 {
		return "", errors.New("index out of range, must be between 0 and 25")
	}
	letter := string('a' + index)
	return letter, nil
}

func CamelToSnake(camel string) string {
	var result strings.Builder

	for i, r := range camel {
		// Check if the rune is uppercase
		if unicode.IsUpper(r) {
			// If it's not the first character, prepend an underscore
			if i > 0 {
				result.WriteRune('_')
			}
			// Convert the uppercase letter to lowercase and add to result
			result.WriteRune(unicode.ToLower(r))
		} else {
			// If it's a lowercase letter, just add it to result
			result.WriteRune(r)
		}
	}

	return result.String()
}
