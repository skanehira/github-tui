package utils

import (
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"
)

// Replace is customized for this project
// https://github.com/golang/go/blob/a8942d2cffd80c68febe1c908a0eb464d2f5bb40/src/strings/strings.go#L924
func Replace(s, old, new string, n int) ([]string, string) {
	if old == new || n == 0 {
		return nil, s // avoid allocation
	}

	// Compute number of replacements.
	if m := strings.Count(s, old); m == 0 {
		return nil, s // avoid allocation
	} else if n < 0 || m < n {
		n = m
	}

	// Apply replacements to buffer.
	t := make([]byte, len(s)+n*(len(new)-len(old)))
	w := 0
	start := 0
	var regionIDs []string
	for i := 0; i < n; i++ {
		j := start
		if len(old) == 0 {
			if i > 0 {
				_, wid := utf8.DecodeRuneInString(s[start:])
				j += wid
			}
		} else {
			j += strings.Index(s[start:], old)
		}
		w += copy(t[w:], s[start:j])
		w += copy(t[w:], fmt.Sprintf(new, i))
		regionIDs = append(regionIDs, strconv.Itoa(i))
		start = j + len(old)
	}
	w += copy(t[w:], s[start:])
	return regionIDs, string(t[0:w])
}
