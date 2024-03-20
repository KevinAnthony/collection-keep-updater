package types_test

import (
	"bytes"
	"testing"

	"github.com/kevinanthony/collection-keep-updater/types"

	. "github.com/smartystreets/goconvey/convey"
)

func TestISBNBooks_Contains(t *testing.T) {
	t.Parallel()

	var (
		bookA = types.ISBNBook{ISBN10: i10A, ISBN13: i13A}
		bookB = types.ISBNBook{ISBN10: i10B, ISBN13: i13B}
		bookC = types.ISBNBook{ISBN10: i10C, ISBN13: i13C}
		bookD = types.ISBNBook{ISBN10: i10D, ISBN13: i13D}
	)

	Convey("ISBNBooks_Contains", t, func() {
		books := types.ISBNBooks{bookA, bookB, bookC}

		Convey(" should return true when", func() {
			Convey("item is first in the list", func() {
				So(books.Contains(bookA), ShouldBeTrue)
			})
			Convey("item is in the middle in the list", func() {
				So(books.Contains(bookB), ShouldBeTrue)
			})
			Convey("item is last in the list", func() {
				So(books.Contains(bookC), ShouldBeTrue)
			})
		})

		Convey("should return false when book is not in the list", func() {
			So(books.Contains(bookD), ShouldBeFalse)
		})
	})
}

func TestISBNBooks_Diff(t *testing.T) {
	t.Parallel()

	var (
		bookA = types.ISBNBook{ISBN10: i10A, ISBN13: i13A}
		bookB = types.ISBNBook{ISBN10: i10B, ISBN13: i13B}
		bookC = types.ISBNBook{ISBN10: i10C, ISBN13: i13C}
		bookD = types.ISBNBook{ISBN10: i10D, ISBN13: i13D}
	)

	Convey("ISBNBooks_Diff", t, func() {
		Convey("should return list containing everything in list A not in list B", func() {
			booksA := types.ISBNBooks{bookD, bookC, bookB}
			booksB := types.ISBNBooks{bookA, bookC}
			expected := types.ISBNBooks{bookD, bookB}

			actual := booksA.Diff(booksB)

			So(actual, ShouldResemble, expected)
		})
	})
}

func TestISBNBooks_FindIndexByISBN(t *testing.T) {
	t.Parallel()

	var (
		bookA = types.ISBNBook{ISBN10: i10A, ISBN13: i13A}
		bookB = types.ISBNBook{ISBN10: i10B, ISBN13: i13B}
		bookC = types.ISBNBook{ISBN10: i10C, ISBN13: i13C}
	)

	Convey("ISBNBooks_FindIndexByISBN", t, func() {
		books := types.ISBNBooks{bookA, bookB, bookC}

		Convey("sanity check", func() {
			So(len(i10A), ShouldEqual, 10)
			So(len(i10B), ShouldEqual, 10)
			So(len(i10C), ShouldEqual, 10)
			So(len(i10D), ShouldEqual, 10)

			So(len(i13A), ShouldEqual, 13)
			So(len(i13B), ShouldEqual, 13)
			So(len(i13C), ShouldEqual, 13)
			So(len(i13D), ShouldEqual, 13)
		})
		Convey("should return index if searching by", func() {
			Convey("ISBN10", func() {
				So(books.FindIndexByISBN(i10B), ShouldEqual, 1)
			})
			Convey("ISBN13", func() {
				So(books.FindIndexByISBN(i13A), ShouldEqual, 0)
			})
		})
		Convey("should return -1 if", func() {
			Convey("searching by invalid ISBN", func() {
				So(books.FindIndexByISBN("invalid"), ShouldEqual, -1)
			})
			Convey("isbn is not in list", func() {
				Convey("and if valid ISBN10", func() {
					So(books.FindIndexByISBN(i10D), ShouldEqual, -1)
				})
				Convey("and if valid ISBN13", func() {
					So(books.FindIndexByISBN(i13D), ShouldEqual, -1)
				})
			})
		})
	})
}

func TestISBNBooks_Print(t *testing.T) {
	t.Parallel()

	isbnTable := `┌───────┬────────┬────────────┬───────────────┬────────┐
│ TITLE │ VOLUME │ ISBN 10    │ ISBN 13       │ SOURCE │
├───────┼────────┼────────────┼───────────────┼────────┤
│       │        │ 0123456789 │ 0123456789ABC │        │
│       │        │ 1234567890 │ 123456789ABC0 │        │
│       │        │ 2345678901 │ 23456789ABC01 │        │
│       │        │ 3456789012 │ 3456789ABC012 │        │
└───────┴────────┴────────────┴───────────────┴────────┘
`
	var (
		bookA = types.ISBNBook{ISBN10: i10A, ISBN13: i13A}
		bookB = types.ISBNBook{ISBN10: i10B, ISBN13: i13B}
		bookC = types.ISBNBook{ISBN10: i10C, ISBN13: i13C}
		bookD = types.ISBNBook{ISBN10: i10D, ISBN13: i13D}
	)

	Convey("ISBNBooks_Print", t, func() {
		books := types.ISBNBooks{bookA, bookB, bookC, bookD}

		writer := bytes.NewBufferString("")

		cmdMock := types.NewICommandMock(t)
		cmdMock.On("OutOrStdout").Once().Return(writer)

		Convey("should print out expected table", func() {
			books.Print(cmdMock)

			So(writer.String(), ShouldEqual, isbnTable)
		})
	})
}

func TestISBNBooks_RemoveAt(t *testing.T) {
	t.Parallel()

	var (
		bookA = types.ISBNBook{ISBN10: i10A, ISBN13: i13A}
		bookB = types.ISBNBook{ISBN10: i10B, ISBN13: i13B}
		bookC = types.ISBNBook{ISBN10: i10C, ISBN13: i13C}
		bookD = types.ISBNBook{ISBN10: i10D, ISBN13: i13D}
	)

	Convey("ISBNBooks_RemoveAt", t, func() {
		books := types.ISBNBooks{bookA, bookB, bookC, bookD}

		Convey("should return list with index removed", func() {
			Convey("if index is first", func() {
				So(books.RemoveAt(0), ShouldResemble, types.ISBNBooks{bookB, bookC, bookD})
			})
			Convey("if index is in the middle", func() {
				So(books.RemoveAt(2), ShouldResemble, types.ISBNBooks{bookA, bookB, bookD})
			})
			Convey("if index is last", func() {
				So(books.RemoveAt(3), ShouldResemble, types.ISBNBooks{bookA, bookB, bookC})
			})
			Convey("and should not modify original list", func() {
				So(books, ShouldHaveLength, 4)

				So(books.RemoveAt(0).RemoveAt(0).RemoveAt(0).RemoveAt(0), ShouldBeEmpty)

				So(books, ShouldHaveLength, 4)
			})
		})
		Convey("should return original list if", func() {
			Convey("index is negative", func() {
				So(books.RemoveAt(-1), ShouldResemble, books)
			})
			Convey("index is larger then list", func() {
				So(books.RemoveAt(4), ShouldResemble, books)
			})
		})
	})
}

func TestNewISBNBooks(t *testing.T) {
	t.Parallel()

	Convey("NewISBNBooks", t, func() {
		Convey("should return list with length 0 and capacity of 5", func() {
			const booksCap = 5
			books := types.NewISBNBooks(booksCap)

			So(len(books), ShouldEqual, 0)
			So(cap(books), ShouldEqual, booksCap)
		})
	})
}
