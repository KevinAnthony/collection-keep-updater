package yen_test

import (
	"testing"

	"github.com/kevinanthony/collection-keep-updater/source/yen"
	"github.com/kevinanthony/gorps/v2/http"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNew(t *testing.T) {
	t.Parallel()

	Convey("New", t, func() {
		client := http.NewClientMock(t)
		Convey("should return isource when http client is valid", func() {
			source, err := yen.New(client)

			So(source, ShouldNotBeNil)
			So(err, ShouldBeNil)
		})
		Convey("should return error when http client is nil", func() {
			source, err := yen.New(nil)

			So(err, ShouldBeError, "http client is nil")
			So(source, ShouldBeNil)
		})
	})
}
