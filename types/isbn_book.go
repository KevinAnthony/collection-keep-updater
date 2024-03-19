package types

import (
	"strings"
)

type ISBNBook struct {
	ISBN10 string
	ISBN13 string
	Title  string
	Volume string
	Source string
}

func (a ISBNBook) Equals(b ISBNBook) bool {
	if len(a.ISBN13) > 0 && len(b.ISBN13) > 0 {
		return strings.EqualFold(a.ISBN13, b.ISBN13)
	}

	if len(a.ISBN10) > 0 && len(b.ISBN10) > 0 {
		return strings.EqualFold(a.ISBN10, b.ISBN10)
	}

	return false
}
