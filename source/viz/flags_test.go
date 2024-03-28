package viz_test

import (
	"encoding/json"
	"testing"

	"github.com/kevinanthony/collection-keep-updater/source/viz"
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
			cmdMock.On("PersistentFlags").Twice().Return(flags)

			viz.SetFlags(cmdMock)

			backlog, err := flags.GetInt("viz-max-backlog")
			So(err, ShouldBeNil)
			So(backlog, ShouldResemble, 0)

			delay, err := flags.GetString("viz-get-delay")
			So(err, ShouldBeNil)
			So(delay, ShouldResemble, "")
		})
	})
}

func TestSettingsHelper_GetIDFromURL(t *testing.T) {
	t.Parallel()

	Convey("GetIDFromURL", t, func() {
		source := getSource(t)

		Convey("should return valid id when passed valid URL", func() {
			id, err := source.GetIDFromURL("https://www.viz.com/read/manga/chainsaw-man/all")

			So(id, ShouldResemble, "chainsaw-man")
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
		flag.String("viz-get-delay", "", "")
		flag.Int("viz-max-backlog", 0, "")

		source := getSource(t)

		cmdMock.On("PersistentFlags").Return(flag)

		Convey("should return sources settings if", func() {
			Convey("source settings is nil", func() {
				So(flag.Set("viz-get-delay", "81ms"), ShouldBeNil)
				So(flag.Set("viz-max-backlog", "5"), ShouldBeNil)

				actual, err := source.SourceSettingFromFlags(cmdMock, nil)

				So(err, ShouldBeNil)
				So(toJson(t, actual), ShouldResemble, "{\"maximum_backlog\":5,\"delay_between\":81000000}")
			})
			Convey("source settings is valid", func() {
				temp, err := source.SourceSettingFromFlags(cmdMock, nil)
				So(err, ShouldBeNil)

				actual, err := source.SourceSettingFromFlags(cmdMock, temp)

				So(err, ShouldBeNil)
				So(actual, ShouldEqual, temp)
			})
			Convey("should return error when delay is unparsable", func() {
				So(flag.Set("viz-get-delay", "junk"), ShouldBeNil)

				actual, err := source.SourceSettingFromFlags(cmdMock, nil)

				So(err, ShouldBeError, "viz: cannot parse delay junk: time: invalid duration \"junk\"")
				So(toJson(t, actual), ShouldResemble, "{\"maximum_backlog\":0,\"delay_between\":null}")
			})
		})
	})
}

func TestSettingsHelper_SourceSettingFromConfig(t *testing.T) {
	t.Parallel()

	Convey("SourceSettingFromConfig", t, func() {
		data := map[string]interface{}{
			"maximum_backlog": 3,
			"delay_between":   "50ms",
		}
		source := getSource(t)

		Convey("should set from map", func() {
			settings := source.SourceSettingFromConfig(data)

			So(toJson(t, settings), ShouldResemble, `{"maximum_backlog":3,"delay_between":50000000}`)
		})
		Convey("should return empty settings when passed no data", func() {
			settings := source.SourceSettingFromConfig(nil)

			So(toJson(t, settings), ShouldResemble, `{"maximum_backlog":null,"delay_between":null}`)
		})
		Convey("should not set delay when string is un-parseable", func() {
			data["delay_between"] = "junk"
			settings := source.SourceSettingFromConfig(data)

			So(toJson(t, settings), ShouldResemble, `{"maximum_backlog":3,"delay_between":null}`)
		})
	})
}

func toJson(t *testing.T, settings types.ISourceSettings) string {
	t.Helper()

	data, err := json.Marshal(settings)
	So(err, ShouldBeNil)

	return string(data)
}
