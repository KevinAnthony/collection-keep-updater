package types_test

import (
	"testing"

	"github.com/kevinanthony/collection-keep-updater/types"

	. "github.com/smartystreets/goconvey/convey"
)

const (
	i10A = "0123456789"
	i13A = "0123456789ABC"
	i10B = "1234567890"
	i13B = "123456789ABC0"
	i10C = "2345678901"
	i13C = "23456789ABC01"
	i10D = "3456789012"
	i13D = "3456789ABC012"
)

func TestISBNBook_Equals(t *testing.T) {
	t.Parallel()

	Convey("ISBNBook_Equals", t, func() {

		bookA := types.ISBNBook{}
		bookB := types.ISBNBook{}

		Convey("should equal when", func() {
			Convey("ISBN10 and ISBN13 are both equal", func() {
				bookA.ISBN10 = i10A
				bookA.ISBN13 = i13A
				bookB.ISBN10 = i10A
				bookB.ISBN13 = i13A

				So(bookA.Equals(bookB), ShouldBeTrue)
				So(bookB.Equals(bookA), ShouldBeTrue)
			})
			Convey("ISBN10 is equal and ISBN13", func() {
				Convey("book A is empty", func() {
					bookA.ISBN10 = i10A
					bookB.ISBN10 = i10A
					bookB.ISBN13 = i13A

					So(bookA.Equals(bookB), ShouldBeTrue)
					So(bookB.Equals(bookA), ShouldBeTrue)
				})
				Convey("book B is empty", func() {
					bookA.ISBN10 = i10A
					bookA.ISBN13 = i13A
					bookB.ISBN10 = i10A

					So(bookA.Equals(bookB), ShouldBeTrue)
					So(bookB.Equals(bookA), ShouldBeTrue)
				})
			})
			Convey("ISBN13 is equal and ISBN10 is empty", func() {
				Convey("book A is empty", func() {
					bookA.ISBN13 = i13A
					bookB.ISBN13 = i13A
					bookB.ISBN10 = i10A

					So(bookA.Equals(bookB), ShouldBeTrue)
					So(bookB.Equals(bookA), ShouldBeTrue)
				})
				Convey("book B is empty", func() {
					bookA.ISBN13 = i13A
					bookA.ISBN10 = i10A
					bookB.ISBN13 = i13A

					So(bookA.Equals(bookB), ShouldBeTrue)
					So(bookB.Equals(bookA), ShouldBeTrue)
				})
			})
			Convey("ISBN13 is equal and ISBN10 are not equal", func() {
				bookA.ISBN10 = i10A
				bookA.ISBN13 = i13A
				bookB.ISBN10 = i10B
				bookB.ISBN13 = i13A

				So(bookA.Equals(bookB), ShouldBeTrue)
				So(bookB.Equals(bookA), ShouldBeTrue)
			})
		})
		Convey("should not equal when", func() {
			Convey("ISBN10 and ISBN13 are both not equal", func() {
				bookA.ISBN10 = i10A
				bookA.ISBN13 = i13A
				bookB.ISBN10 = i10B
				bookB.ISBN13 = i13B

				So(bookA.Equals(bookB), ShouldBeFalse)
				So(bookB.Equals(bookA), ShouldBeFalse)
			})
			Convey("ISBN10 is not equal and ISBN13", func() {
				Convey("book A is empty", func() {
					bookA.ISBN10 = i10A
					bookB.ISBN10 = i10B
					bookB.ISBN13 = i13A

					So(bookA.Equals(bookB), ShouldBeFalse)
					So(bookB.Equals(bookA), ShouldBeFalse)
				})
				Convey("book B is empty", func() {
					bookA.ISBN10 = i10A
					bookA.ISBN13 = i13A
					bookB.ISBN10 = i10B

					So(bookA.Equals(bookB), ShouldBeFalse)
					So(bookB.Equals(bookA), ShouldBeFalse)
				})
			})
			Convey("ISBN13 is not equal and ISBN10 is empty", func() {
				Convey("book A is empty", func() {
					bookA.ISBN13 = i13A
					bookB.ISBN13 = i13B
					bookB.ISBN10 = i10A

					So(bookA.Equals(bookB), ShouldBeFalse)
					So(bookB.Equals(bookA), ShouldBeFalse)
				})
				Convey("book B is empty", func() {
					bookA.ISBN13 = i13A
					bookA.ISBN10 = i10A
					bookB.ISBN13 = i13B

					So(bookA.Equals(bookB), ShouldBeFalse)
					So(bookB.Equals(bookA), ShouldBeFalse)
				})
			})
			Convey("book A and Book B are both empty", func() {
				So(bookA.Equals(bookB), ShouldBeFalse)
				So(bookB.Equals(bookA), ShouldBeFalse)
			})
			Convey("book A ISBN13 and ISBN10 are both empty", func() {
				bookB.ISBN10 = i10B
				bookB.ISBN13 = i13B

				So(bookA.Equals(bookB), ShouldBeFalse)
				So(bookB.Equals(bookA), ShouldBeFalse)
			})
			Convey("book B ISBN13 and ISBN10 are both empty", func() {
				bookA.ISBN10 = i10A
				bookA.ISBN13 = i13A

				So(bookA.Equals(bookB), ShouldBeFalse)
				So(bookB.Equals(bookA), ShouldBeFalse)
			})
			Convey("book A ISBN13 and book B ISBN10 are both empty", func() {
				bookA.ISBN10 = i10B
				bookB.ISBN13 = i13B

				So(bookA.Equals(bookB), ShouldBeFalse)
				So(bookB.Equals(bookA), ShouldBeFalse)
			})
			Convey("book B ISBN13 and book A ISBN10 are both empty", func() {
				bookB.ISBN10 = i10B
				bookA.ISBN13 = i13B

				So(bookA.Equals(bookB), ShouldBeFalse)
				So(bookB.Equals(bookA), ShouldBeFalse)
			})
		})
	})
}
