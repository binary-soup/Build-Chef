package common

import (
	"strings"
	"unicode"
)

func SelectPlural(singular string, plural string, count int) string {
	if count == 1 {
		return singular
	} else {
		return plural
	}
}

func ReplaceChar(str string, chars string, replace rune) string {
	result := make([]rune, len(str))

	for i, char := range str {
		if strings.ContainsRune(chars, char) {
			result[i] = replace
		} else {
			result[i] = char
		}
	}

	return string(result)
}

func ToUpper(str string) string {
	result := make([]rune, len(str))

	for i, char := range str {
		result[i] = unicode.ToUpper(char)
	}

	return string(result)
}
