package utils_test

import (
	"testing"

	"github.com/kevinanthony/collection-keep-updater/utils"

	. "github.com/smartystreets/goconvey/convey"
)

func TestISBNNormalize(t *testing.T) {
	t.Parallel()

	Convey("ISBNNormalize", t, func() {
		Convey("should return ISBN that has been normalized", func() {
			isbn := utils.ISBNNormalize(" this-is-a test-isbn 12345   ")

			So(isbn, ShouldResemble, "thisisatestisbn12345")
		})
	})
}
