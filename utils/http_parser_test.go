package utils_test

import (
	"testing"

	"github.com/kevinanthony/collection-keep-updater/utils"

	. "github.com/smartystreets/goconvey/convey"
	"golang.org/x/net/html"
)

func TestAttrContains(t *testing.T) {
	t.Parallel()

	Convey("AttrContains", t, func() {
		attr := []html.Attribute{{Key: "one", Val: "test one"}, {Key: "two", Val: "test two"}}

		Convey("should return attribute and true if list contains attribute", func() {
			value, found := utils.AttrContains(attr, "one")

			So(found, ShouldBeTrue)
			So(value, ShouldResemble, "test one")
		})
		Convey("should return false and empty if list contains attribute", func() {
			value, found := utils.AttrContains(attr, "three")

			So(found, ShouldBeFalse)
			So(value, ShouldBeEmpty)
		})
	})
}

func TestAttrEquals(t *testing.T) {
	t.Parallel()

	Convey("AttrEquals", t, func() {
		attr := []html.Attribute{{Key: "one", Val: "test one"}, {Key: "two", Val: "test two"}}

		Convey("should return true if list contains attribute", func() {
			found := utils.AttrEquals(attr, "one", "test one")

			So(found, ShouldBeTrue)
		})
		Convey("should return false if", func() {
			Convey("attribute does not contain key", func() {
				found := utils.AttrEquals(attr, "three", "test three")

				So(found, ShouldBeFalse)
			})
			Convey("attribute contains key but value is different", func() {
				found := utils.AttrEquals(attr, "two", "test three")

				So(found, ShouldBeFalse)
			})
		})
	})
}
