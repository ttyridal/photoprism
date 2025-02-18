package txt

import "strings"

const (
	ClipDefault     = 160
	ClipKeyword     = 40
	ClipSlug        = 80
	ClipVarchar     = 255
	ClipQuery       = 1000
	ClipDescription = 16000
)

// Clip shortens a string to the given number of runes, and removes all leading and trailing white space.
func Clip(s string, size int) string {
	s = strings.TrimSpace(s)

	if s == "" || size <= 0 {
		return ""
	}

	runes := []rune(s)

	if len(runes) > size {
		s = string(runes[0 : size-1])
	}

	return s
}

func TrimLen(s string, size int) string {
	if len(s) < size || size < 4 {
		return s
	}

	return Clip(s, size-3) + "..."
}
