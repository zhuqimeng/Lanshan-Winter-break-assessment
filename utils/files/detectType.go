package files

import "strings"

func IsMarkdown(filename string) bool {
	f := strings.ToLower(filename)
	return strings.HasSuffix(f, ".md") ||
		strings.HasSuffix(f, ".markdown")
}
