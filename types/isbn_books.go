package types

import (
	"github.com/kevinanthony/collection-keep-updater/out"

	"github.com/jedib0t/go-pretty/v6/table"
)

type ISBNBooks []ISBNBook

func NewISBNBooks(length int) ISBNBooks {
	return make(ISBNBooks, 0, length)
}

func (b ISBNBooks) Diff(other ISBNBooks) ISBNBooks {
	diff := NewISBNBooks(0)

	for _, book := range b {
		if other.Contains(book) {
			continue
		}

		diff = append(diff, book)
	}

	return diff
}

func (b ISBNBooks) Contains(book ISBNBook) bool {
	for _, l := range b {
		if l.Equals(book) {
			return true
		}
	}

	return false
}

func (b ISBNBooks) FindIndexByISBN(isbn string) int {
	var book ISBNBook
	if len(isbn) == 10 {
		book.ISBN10 = isbn
	} else if len(isbn) == 13 {
		book.ISBN13 = isbn
	} else {
		return -1
	}

	for i, l := range b {
		if l.Equals(book) {
			return i
		}
	}

	return -1
}

func (b ISBNBooks) RemoveAt(i int) ISBNBooks {
	if i < 0 || i >= len(b) {
		return b
	}

	return append(b[:i], b[i+1:]...)
}

func (b ISBNBooks) Print(cmd ICommand) {
	t := out.NewTable(cmd.OutOrStdout())
	t.AppendHeader(table.Row{"Title", "Volume", "ISBN 10", "ISBN 13", "Source"})
	for _, book := range b {
		t.AppendRow(table.Row{book.Title, book.Volume, book.ISBN10, book.ISBN13, book.Source})
	}

	t.Render()
}
