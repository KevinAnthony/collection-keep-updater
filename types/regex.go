package types

import "regexp"

var (
	ISBN10regex = regexp.MustCompile(`([0-9]{1,5}[-\\ ]?[0-9]+[-\\ ]?[0-9]+[-\\ ]?[0-9X])`)
	ISBN13regex = regexp.MustCompile(`(97[89][- ]?[0-9]{1,5}[- ]?[0-9]+[- ]?[0-9]+[- ]?[0-9Xx])`)
)
