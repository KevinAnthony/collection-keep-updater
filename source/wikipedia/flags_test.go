package wikipedia_test

import (
	"encoding/json"
	"testing"

	"github.com/kevinanthony/collection-keep-updater/source/wikipedia"
	"github.com/kevinanthony/collection-keep-updater/types"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/spf13/pflag"
)

func TestSetFlags(t *testing.T) {
	t.Parallel()

	Convey("SetFlags", t, func() {
		cmdMock := types.NewICommandMock(t)

		Convey("should set two flags", func() {
			flags := &pflag.FlagSet{}
			cmdMock.On("PersistentFlags").Times(4).Return(flags)

			wikipedia.SetFlags(cmdMock)

			vheader, err := flags.GetString("wiki-volume-header")
			So(err, ShouldBeNil)
			So(vheader, ShouldResemble, "")

			iheader, err := flags.GetString("wiki-isbn-header")
			So(err, ShouldBeNil)
			So(iheader, ShouldResemble, "")

			theader, err := flags.GetString("wiki-title-header")
			So(err, ShouldBeNil)
			So(theader, ShouldResemble, "")

			tables, err := flags.GetIntSlice("wiki-table-numbers")
			So(err, ShouldBeNil)
			So(tables, ShouldResemble, []int{})
		})
	})
}

func TestSettingsHelper_GetIDFromURL(t *testing.T) {
	t.Parallel()

	Convey("GetIDFromURL", t, func() {
		source := getSource(t)

		Convey("should return valid id when passed valid URL", func() {
			id, err := source.GetIDFromURL("https://en.wikipedia.org/wiki/Lists_of_One_Piece_chapters")

			So(id, ShouldResemble, "Lists_of_One_Piece_chapters")
			So(err, ShouldBeNil)
		})
		Convey("should return error when", func() {
			Convey("url is empty", func() {
				id, err := source.GetIDFromURL("")

				So(id, ShouldBeEmpty)
				So(err, ShouldBeError, "unknown/unset url. url is required")
			})
			Convey("url is malformed", func() {
				id, err := source.GetIDFromURL("???")

				So(id, ShouldBeEmpty)
				So(err, ShouldBeError, "url is malformed")
			})
		})
	})
}

func TestSettingsHelper_SourceSettingFromFlags(t *testing.T) {
	t.Parallel()

	Convey("SourceSettingFromFlags", t, func() {
		cmdMock := types.NewICommandMock(t)
		flag := &pflag.FlagSet{}
		flag.String("wiki-volume-header", "", "")
		flag.String("wiki-isbn-header", "", "")
		flag.String("wiki-title-header", "", "")
		flag.IntSlice("wiki-table-numbers", []int{}, "")

		source := getSource(t)

		cmdMock.On("PersistentFlags").Return(flag)

		Convey("should return sources settings if", func() {
			Convey("source settings is nil", func() {
				So(flag.Set("wiki-volume-header", "volume"), ShouldBeNil)
				So(flag.Set("wiki-isbn-header", "isbn"), ShouldBeNil)
				So(flag.Set("wiki-title-header", "table"), ShouldBeNil)
				So(flag.Set("wiki-table-numbers", "1,2,3"), ShouldBeNil)

				actual, err := source.SourceSettingFromFlags(cmdMock, nil)

				So(err, ShouldBeNil)
				So(toJson(t, actual), ShouldResemble, "{\"volume_header\":\"volume\",\"title_header\":\"table\",\"isbn_header\":\"isbn\",\"tables\":[1,2,3]}")
			})
			Convey("source settings is valid", func() {
				temp, err := source.SourceSettingFromFlags(cmdMock, nil)
				So(err, ShouldBeNil)

				actual, err := source.SourceSettingFromFlags(cmdMock, temp)

				So(err, ShouldBeNil)
				So(actual, ShouldEqual, temp)
			})
		})
	})
}

func TestSettingsHelper_SourceSettingFromConfig(t *testing.T) {
	t.Parallel()

	Convey("SourceSettingFromConfig", t, func() {
		data := map[string]interface{}{
			"volume_header": "v_header",
			"title_header":  "t_header",
			"isbn_header":   "i_header",
			"tables":        []interface{}{1, 2, 3, 4},
		}
		source := getSource(t)

		Convey("should set from map", func() {
			settings := source.SourceSettingFromConfig(data)

			So(toJson(t, settings), ShouldResemble,
				`{"volume_header":"v_header","title_header":"t_header","isbn_header":"i_header","tables":[1,2,3,4]}`)
		})
		Convey("should return empty settings when passed no data", func() {
			settings := source.SourceSettingFromConfig(nil)

			So(toJson(t, settings), ShouldResemble,
				`{"volume_header":null,"title_header":null,"isbn_header":null,"tables":null}`)
		})
		Convey("should not set delay when string is un-parseable", func() {
			data["tables"] = "junk"
			settings := source.SourceSettingFromConfig(data)

			So(toJson(t, settings), ShouldResemble,
				`{"volume_header":"v_header","title_header":"t_header","isbn_header":"i_header","tables":null}`)
		})
		Convey("should return empty settings if isbn header is nil", func() {
			data["isbn_header"] = ""

			settings := source.SourceSettingFromConfig(data)

			So(toJson(t, settings), ShouldResemble,
				`{"volume_header":null,"title_header":null,"isbn_header":null,"tables":null}`)
		})
	})
}

func toJson(t *testing.T, settings types.ISourceSettings) string {
	t.Helper()

	data, err := json.Marshal(settings)
	So(err, ShouldBeNil)

	return string(data)
}
