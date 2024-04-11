package kodansha_test

import (
	"testing"

	"github.com/kevinanthony/collection-keep-updater/ctxu"
	"github.com/kevinanthony/collection-keep-updater/source/kodansha"
	"github.com/kevinanthony/collection-keep-updater/types"
	"github.com/kevinanthony/gorps/v2/http"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSetFlags(t *testing.T) {
	t.Parallel()

	Convey("SetFlags", t, func() {
		cmdMock := types.NewICommandMock(t)

		kodansha.SetFlags(cmdMock)
	})
}

func TestSettingsHelper_SourceSettingFromConfig(t *testing.T) {
	t.Parallel()

	Convey("SourceSettingFromConfig", t, func() {
		source := getSource(t)

		Convey("should return source setting", func() {
			settings := source.SourceSettingFromConfig(nil)

			So(settings, ShouldNotBeNil)
		})
	})
}

func TestSettingsHelper_SourceSettingFromFlags(t *testing.T) {
	t.Parallel()

	Convey("SourceSettingFromFlags", t, func() {
		cmdMock := types.NewICommandMock(t)
		source := getSource(t)

		Convey("should return sources settings if", func() {
			Convey("source settings is nil", func() {
				actual, err := source.SourceSettingFromFlags(cmdMock, nil)

				So(err, ShouldBeNil)
				So(actual, ShouldNotBeNil)
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

func TestSettingsHelper_GetIDFromURL(t *testing.T) {
	t.Parallel()

	Convey("GetIDFromURL", t, func() {
		source := getSource(t)

		Convey("should return valid id when passed valid URL", func() {
			id, err := source.GetIDFromURL("https://kodansha.us/series/initial-d-omnibus")

			So(id, ShouldResemble, "initial-d-omnibus")
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

func getSource(t *testing.T) types.ISource {
	cmd := types.NewICommandMock(t)
	ctx := ctxu.NewContextMock(t)
	client := http.NewClientMock(t)

	cmd.On("Context").Return(ctx)
	ctx.On("Value", ctxu.ContextKey("http_ctx_key")).Return(client)

	source := kodansha.New(cmd)

	return source
}
