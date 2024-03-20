package utils

import "strings"

func ISBNNormalize(isbn string) string {
	return strings.ReplaceAll(strings.ReplaceAll(isbn, "-", ""), " ", "")
}
