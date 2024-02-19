package types

import "strings"

type ISBNBook struct {
	ISBN10  string
	ISBN13  string
	Title   string
	Binding string
	Volume  string
}

func ISBNBookCmp(A, B ISBNBook) bool {
	if len(A.ISBN13) > 0 && len(B.ISBN13) > 0 {
		return strings.EqualFold(A.ISBN13, B.ISBN13)
	}

	if len(A.ISBN10) > 0 && len(B.ISBN10) > 0 {
		return strings.EqualFold(A.ISBN10, B.ISBN10)
	}

	return false
}
