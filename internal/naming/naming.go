// Package naming provides string case conversion utilities.
package naming

import (
	"strings"
	"unicode"
)

// Common Go acronyms that should be treated as single words.
var commonAcronyms = map[string]bool{
	"ID": true, "URL": true, "HTTP": true, "HTTPS": true,
	"API": true, "JSON": true, "XML": true, "SQL": true,
	"HTML": true, "CSS": true, "URI": true, "UUID": true,
	"IP": true, "TCP": true, "UDP": true, "DNS": true,
	"TLS": true, "SSL": true, "SSH": true, "EOF": true,
	"ACL": true, "TTL": true, "CPU": true, "GPU": true,
	"OS": true, "IO": true, "DB": true, "OK": true,
}

// ToSnakeCase converts a string to snake_case.
func ToSnakeCase(s string) string {
	words := splitWords(s)
	for i, w := range words {
		words[i] = strings.ToLower(w)
	}
	return strings.Join(words, "_")
}

// ToCamelCase converts a string to camelCase.
func ToCamelCase(s string) string {
	words := splitWords(s)
	for i, w := range words {
		if i == 0 {
			words[i] = strings.ToLower(w)
		} else {
			words[i] = titleWord(w)
		}
	}
	return strings.Join(words, "")
}

// ToPascalCase converts a string to PascalCase.
func ToPascalCase(s string) string {
	words := splitWords(s)
	for i, w := range words {
		words[i] = titleWord(w)
	}
	return strings.Join(words, "")
}

// ToKebabCase converts a string to kebab-case.
func ToKebabCase(s string) string {
	words := splitWords(s)
	for i, w := range words {
		words[i] = strings.ToLower(w)
	}
	return strings.Join(words, "-")
}

// Convert applies the named convention to the input string.
// Recognized conventions: "snake_case", "camelCase", "PascalCase", "kebab-case".
// Returns the input unchanged if the convention is not recognized.
func Convert(s, convention string) string {
	switch convention {
	case "snake_case":
		return ToSnakeCase(s)
	case "camelCase":
		return ToCamelCase(s)
	case "PascalCase":
		return ToPascalCase(s)
	case "kebab-case":
		return ToKebabCase(s)
	default:
		return s
	}
}

// MatchesConvention checks if a string already follows the named convention.
func MatchesConvention(s, convention string) bool {
	return s == Convert(s, convention)
}

// titleWord capitalizes the first letter and lowercases the rest,
// unless it's a known acronym.
func titleWord(w string) string {
	upper := strings.ToUpper(w)
	if commonAcronyms[upper] {
		return upper
	}
	if len(w) == 0 {
		return w
	}
	runes := []rune(w)
	runes[0] = unicode.ToUpper(runes[0])
	for i := 1; i < len(runes); i++ {
		runes[i] = unicode.ToLower(runes[i])
	}
	return string(runes)
}

// splitWords breaks a string into words by detecting boundaries:
// - underscores and hyphens (delimiters)
// - transitions from lowercase to uppercase (camelCase boundaries)
// - transitions from uppercase run to uppercase+lowercase (acronym end)
func splitWords(s string) []string {
	var words []string
	var current []rune

	runes := []rune(s)
	for i := 0; i < len(runes); i++ {
		r := runes[i]

		// Delimiter — flush current word
		if r == '_' || r == '-' || r == ' ' {
			if len(current) > 0 {
				words = append(words, string(current))
				current = nil
			}
			continue
		}

		if len(current) == 0 {
			current = append(current, r)
			continue
		}

		prev := current[len(current)-1]

		// lowercase -> uppercase: new word
		if unicode.IsLower(prev) && unicode.IsUpper(r) {
			words = append(words, string(current))
			current = []rune{r}
			continue
		}

		// uppercase -> uppercase -> lowercase: the last uppercase starts a new word
		// e.g., "HTMLParser" -> "HTML", "Parser"
		if unicode.IsUpper(prev) && unicode.IsUpper(r) && i+1 < len(runes) && unicode.IsLower(runes[i+1]) {
			words = append(words, string(current))
			current = []rune{r}
			continue
		}

		current = append(current, r)
	}

	if len(current) > 0 {
		words = append(words, string(current))
	}

	return words
}
