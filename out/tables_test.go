package out_test

import (
	"bytes"
	"testing"

	"github.com/kevinanthony/collection-keep-updater/out"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNewTable(t *testing.T) {
	t.Parallel()

	Convey("NewTable", t, func() {
		w := bytes.NewBuffer([]byte{})
		Convey("should return table", func() {
			table := out.NewTable(w)

			So(table, ShouldNotBeNil)
			So(table.Render(), ShouldResemble, "")
		})
	})
}
