package utils

import "strings"

func ISBNNormalize(isbn string) string {
	return strings.TrimSpace(strings.ReplaceAll(isbn, "-", ""))
}
