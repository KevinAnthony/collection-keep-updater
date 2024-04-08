package out_test

import (
	"testing"

	"github.com/kevinanthony/collection-keep-updater/out"

	. "github.com/smartystreets/goconvey/convey"
)

func TestPartial(t *testing.T) {
	t.Parallel()

	Convey("Partial", t, func() {
		str := "this is a long string, look at how long my string is!"
		Convey("should return partial string of n characters", func() {
			actual := out.Partial(str, 11)

			So(actual, ShouldHaveLength, 14)
			So(actual, ShouldResemble, "this is a l...")
		})
		Convey("should return whole when string is > n", func() {
			actual := out.Partial(str, len(str)+1)

			So(actual, ShouldResemble, str)
		})
	})
}

func TestIntSliceToStrOrEmpty(t *testing.T) {
	t.Parallel()

	Convey("IntSliceToStrOrEmpty", t, func() {
		slice := []int{1, 2, 3, 4, 12}
		Convey("should return a int slice as a string", func() {
			actual := out.IntSliceToStrOrEmpty(slice)

			So(actual, ShouldResemble, "1, 2, 3, 4, 12")
		})
		Convey("should return empty when int is empty", func() {
			actual := out.IntSliceToStrOrEmpty([]int{})

			So(actual, ShouldBeEmpty)
		})
	})
}

func TestValueOrEmpty(t *testing.T) {
	t.Parallel()

	Convey("should return value or empty for type", t, func() {
		Convey("of int", func() {
			Convey("with value", func() {
				i := 81

				actual := out.ValueOrEmpty(&i)

				So(actual, ShouldResemble, i)
			})
			Convey("with nil", func() {
				actual := out.ValueOrEmpty[int](nil)

				So(actual, ShouldResemble, 0)
			})
		})
		Convey("of string", func() {
			Convey("with value", func() {
				i := "81"

				actual := out.ValueOrEmpty(&i)

				So(actual, ShouldResemble, i)
			})
			Convey("with nil", func() {
				actual := out.ValueOrEmpty[string](nil)

				So(actual, ShouldResemble, "")
			})
		})
		Convey("of bool", func() {
			Convey("with value", func() {
				i := true

				actual := out.ValueOrEmpty(&i)

				So(actual, ShouldResemble, i)
			})
			Convey("with nil", func() {
				actual := out.ValueOrEmpty[bool](nil)

				So(actual, ShouldResemble, false)
			})
		})
	})
}
