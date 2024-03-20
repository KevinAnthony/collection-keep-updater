package utils_test

import (
	"testing"

	"github.com/kevinanthony/collection-keep-updater/utils"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGet(t *testing.T) {
	t.Parallel()

	Convey("Get", t, func() {
		dict := getDict()
		Convey("should return value if value is found and of correct type", func() {
			actual := utils.Get[int](dict, "int")

			So(actual, ShouldResemble, 1)
		})
		Convey("should return default value if", func() {
			Convey("type is int and ", func() {
				Convey("key is missing from map", func() {
					actual := utils.Get[int](dict, "invalid")

					So(actual, ShouldResemble, 0)
				})
				Convey("value in map does not match type", func() {
					actual := utils.Get[int](dict, "string")

					So(actual, ShouldResemble, 0)
				})
			})
			Convey("type is string and ", func() {
				Convey("key is missing from map", func() {
					actual := utils.Get[string](dict, "invalid")

					So(actual, ShouldBeEmpty)
				})
				Convey("value in map does not match type", func() {
					actual := utils.Get[string](dict, "int")

					So(actual, ShouldBeEmpty)
				})
			})
		})
	})
}

func TestGetPtr(t *testing.T) {
	t.Parallel()

	Convey("GetPtr", t, func() {
		dict := getDict()
		Convey("should return value if value is found and of correct type", func() {
			actual := utils.GetPtr[int](dict, "int")

			So(*actual, ShouldResemble, 1)
		})
		Convey("should return default value if", func() {
			Convey("type is int and ", func() {
				Convey("key is missing from map", func() {
					actual := utils.GetPtr[int](dict, "invalid")

					So(actual, ShouldBeNil)
				})
				Convey("value in map does not match type", func() {
					actual := utils.GetPtr[int](dict, "string")

					So(actual, ShouldBeNil)
				})
			})
			Convey("type is string and ", func() {
				Convey("key is missing from map", func() {
					actual := utils.GetPtr[string](dict, "invalid")

					So(actual, ShouldBeNil)
				})
				Convey("value in map does not match type", func() {
					actual := utils.GetPtr[string](dict, "int")

					So(actual, ShouldBeNil)
				})
			})
		})
	})
}

func TestGetArray(t *testing.T) {
	t.Parallel()

	Convey("GetArray", t, func() {
		dict := getDict()
		Convey("should return value if value is found and of correct type", func() {
			actual := utils.GetArray[int](dict, "int_array")

			So(actual, ShouldResemble, []int{1, 2, 3, 4})
		})
		Convey("should return default value if", func() {
			Convey("array is of mixed type", func() {
				actual := utils.GetArray[int](dict, "mixed_array")

				So(actual, ShouldBeNil)
			})
			Convey("type is int and ", func() {
				Convey("key is missing from map", func() {
					actual := utils.GetArray[int](dict, "invalid")

					So(actual, ShouldBeNil)
				})
				Convey("value is not an array", func() {
					actual := utils.GetArray[int](dict, "int")

					So(actual, ShouldBeNil)
				})
				Convey("value in map does not match type", func() {
					actual := utils.GetArray[int](dict, "str_array")

					So(actual, ShouldBeNil)
				})
			})
			Convey("type is string and ", func() {
				Convey("key is missing from map", func() {
					actual := utils.GetArray[string](dict, "invalid")

					So(actual, ShouldBeNil)
				})
				Convey("value is not an array", func() {
					actual := utils.GetArray[int](dict, "string")

					So(actual, ShouldBeNil)
				})
				Convey("value in map does not match type", func() {
					actual := utils.GetArray[string](dict, "int_array")

					So(actual, ShouldBeNil)
				})
			})
		})
	})
}

func getDict() map[string]interface{} {
	oneI := 1
	oneS := "one"
	return map[string]interface{}{
		"int_array":   []interface{}{1, 2, 3, 4},
		"str_array":   []interface{}{"one", "two", "three", "four"},
		"mixed_array": []interface{}{"one", 2, "three", 4},
		"int_ptr":     &oneI,
		"str_ptr":     &oneS,
		"int":         oneI,
		"str":         oneS,
	}
}
